package auth

import (
	"context"

	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/config"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/constant"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/db/cache"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/db/controller"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/mcontext"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/tokenverify"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/discoveryregistry"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/errs"
	pbAuth "github.com/xiaoyiEdu/Open-IM-Server/pkg/proto/auth"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/proto/msggateway"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/rpcclient"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/utils"
	"google.golang.org/grpc"
)

type authServer struct {
	authDatabase   controller.AuthDatabase
	userRpcClient  *rpcclient.UserRpcClient
	RegisterCenter discoveryregistry.SvcDiscoveryRegistry
}

func Start(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	userRpcClient := rpcclient.NewUserRpcClient(client)
	pbAuth.RegisterAuthServer(server, &authServer{
		userRpcClient:  &userRpcClient,
		RegisterCenter: client,
		authDatabase:   controller.NewAuthDatabase(cache.NewMsgCacheModel(rdb), config.Config.Secret, config.Config.TokenPolicy.Expire),
	})
	return nil
}

func (s *authServer) UserToken(ctx context.Context, req *pbAuth.UserTokenReq) (*pbAuth.UserTokenResp, error) {
	resp := pbAuth.UserTokenResp{}
	if req.Secret != config.Config.Secret {
		return nil, errs.ErrIdentity.Wrap("secret invalid")
	}
	if _, err := s.userRpcClient.GetUserInfo(ctx, req.UserID); err != nil {
		return nil, err
	}
	token, err := s.authDatabase.CreateToken(ctx, req.UserID, int(req.PlatformID))
	if err != nil {
		return nil, err
	}
	resp.Token = token
	resp.ExpireTimeSeconds = config.Config.TokenPolicy.Expire * 24 * 60 * 60
	return &resp, nil
}

func (s *authServer) parseToken(ctx context.Context, tokensString string) (claims *tokenverify.Claims, err error) {
	claims, err = tokenverify.GetClaimFromToken(tokensString)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	m, err := s.authDatabase.GetTokensWithoutError(ctx, claims.UserID, claims.PlatformID)
	if err != nil {
		return nil, err
	}
	if len(m) == 0 {
		return nil, errs.ErrTokenNotExist.Wrap()
	}
	if v, ok := m[tokensString]; ok {
		switch v {
		case constant.NormalToken:
			return claims, nil
		case constant.KickedToken:
			return nil, errs.ErrTokenKicked.Wrap()
		default:
			return nil, utils.Wrap(errs.ErrTokenUnknown, "")
		}
	}
	return nil, errs.ErrTokenNotExist.Wrap()
}

func (s *authServer) ParseToken(ctx context.Context, req *pbAuth.ParseTokenReq) (resp *pbAuth.ParseTokenResp, err error) {
	resp = &pbAuth.ParseTokenResp{}
	claims, err := s.parseToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	resp.UserID = claims.UserID
	resp.Platform = constant.PlatformIDToName(claims.PlatformID)
	resp.ExpireTimeSeconds = claims.ExpiresAt.Unix()
	return resp, nil
}

func (s *authServer) ForceLogout(ctx context.Context, req *pbAuth.ForceLogoutReq) (*pbAuth.ForceLogoutResp, error) {
	resp := pbAuth.ForceLogoutResp{}
	if err := tokenverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if err := s.forceKickOff(ctx, req.UserID, req.PlatformID, mcontext.GetOperationID(ctx)); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *authServer) forceKickOff(ctx context.Context, userID string, platformID int32, operationID string) error {
	conns, err := s.RegisterCenter.GetConns(ctx, config.Config.RpcRegisterName.OpenImMessageGatewayName)
	if err != nil {
		return err
	}
	for _, v := range conns {
		client := msggateway.NewMsgGatewayClient(v)
		kickReq := &msggateway.KickUserOfflineReq{KickUserIDList: []string{userID}, PlatformID: platformID}
		_, err := client.KickUserOffline(ctx, kickReq)
		s.RegisterCenter.CloseConn(v)
		return utils.Wrap(err, "")
	}
	return errs.ErrInternalServer.Wrap()
}
