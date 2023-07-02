package rpcclient

import (
	"context"

	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/config"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/discoveryregistry"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/proto/auth"
	"google.golang.org/grpc"
)

func NewAuth(discov discoveryregistry.SvcDiscoveryRegistry) *Auth {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImAuthName)
	if err != nil {
		panic(err)
	}
	client := auth.NewAuthClient(conn)
	return &Auth{discov: discov, conn: conn, Client: client}
}

type Auth struct {
	conn   grpc.ClientConnInterface
	Client auth.AuthClient
	discov discoveryregistry.SvcDiscoveryRegistry
}
