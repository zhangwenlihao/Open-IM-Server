package rpcclient

import (
	"context"
	"strings"
	"time"

	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/config"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/tokenverify"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/discoveryregistry"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/errs"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/proto/sdkws"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/proto/user"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/utils"
	"google.golang.org/grpc"
)

type User struct {
	conn   grpc.ClientConnInterface
	Client user.UserClient
	discov discoveryregistry.SvcDiscoveryRegistry
}

func NewUser(discov discoveryregistry.SvcDiscoveryRegistry) *User {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	conn, err := discov.GetConn(ctx, config.Config.RpcRegisterName.OpenImUserName)
	if err != nil {
		panic(err)
	}
	client := user.NewUserClient(conn)
	return &User{discov: discov, Client: client, conn: conn}
}

type UserRpcClient User

func NewUserRpcClient(client discoveryregistry.SvcDiscoveryRegistry) UserRpcClient {
	return UserRpcClient(*NewUser(client))
}

func (u *UserRpcClient) GetUsersInfo(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error) {
	resp, err := u.Client.GetDesignateUsers(ctx, &user.GetDesignateUsersReq{
		UserIDs: userIDs,
	})
	if err != nil {
		return nil, err
	}
	if ids := utils.Single(userIDs, utils.Slice(resp.UsersInfo, func(e *sdkws.UserInfo) string {
		return e.UserID
	})); len(ids) > 0 {
		return nil, errs.ErrUserIDNotFound.Wrap(strings.Join(ids, ","))
	}
	return resp.UsersInfo, nil
}

func (u *UserRpcClient) GetUserInfo(ctx context.Context, userID string) (*sdkws.UserInfo, error) {
	users, err := u.GetUsersInfo(ctx, []string{userID})
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

func (u *UserRpcClient) GetUsersInfoMap(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error) {
	users, err := u.GetUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMap(users, func(e *sdkws.UserInfo) string {
		return e.UserID
	}), nil
}

func (u *UserRpcClient) GetPublicUserInfos(ctx context.Context, userIDs []string, complete bool) ([]*sdkws.PublicUserInfo, error) {
	users, err := u.GetUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return utils.Slice(users, func(e *sdkws.UserInfo) *sdkws.PublicUserInfo {
		return &sdkws.PublicUserInfo{
			UserID:   e.UserID,
			Nickname: e.Nickname,
			FaceURL:  e.FaceURL,
			Ex:       e.Ex,
		}
	}), nil
}

func (u *UserRpcClient) GetPublicUserInfo(ctx context.Context, userID string) (*sdkws.PublicUserInfo, error) {
	users, err := u.GetPublicUserInfos(ctx, []string{userID}, true)
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

func (u *UserRpcClient) GetPublicUserInfoMap(ctx context.Context, userIDs []string, complete bool) (map[string]*sdkws.PublicUserInfo, error) {
	users, err := u.GetPublicUserInfos(ctx, userIDs, complete)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMap(users, func(e *sdkws.PublicUserInfo) string {
		return e.UserID
	}), nil
}

func (u *UserRpcClient) GetUserGlobalMsgRecvOpt(ctx context.Context, userID string) (int32, error) {
	resp, err := u.Client.GetGlobalRecvMessageOpt(ctx, &user.GetGlobalRecvMessageOptReq{
		UserID: userID,
	})
	if err != nil {
		return 0, err
	}
	return resp.GlobalRecvMsgOpt, err
}

func (u *UserRpcClient) Access(ctx context.Context, ownerUserID string) error {
	_, err := u.GetUserInfo(ctx, ownerUserID)
	if err != nil {
		return err
	}
	return tokenverify.CheckAccessV3(ctx, ownerUserID)
}
