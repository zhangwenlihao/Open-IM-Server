package main

import (
	"github.com/xiaoyiEdu/Open-IM-Server/internal/rpc/user"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/cmd"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/config"
)

func main() {
	rpcCmd := cmd.NewRpcCmd("user")
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(); err != nil {
		panic(err.Error())
	}
	if err := rpcCmd.StartSvr(config.Config.RpcRegisterName.OpenImUserName, user.Start); err != nil {
		panic(err.Error())
	}
}
