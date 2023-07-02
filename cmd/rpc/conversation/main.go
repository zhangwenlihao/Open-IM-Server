package main

import (
	"github.com/xiaoyiEdu/Open-IM-Server/internal/rpc/conversation"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/cmd"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/config"
)

func main() {
	rpcCmd := cmd.NewRpcCmd("conversation")
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(); err != nil {
		panic(err.Error())
	}
	if err := rpcCmd.StartSvr(config.Config.RpcRegisterName.OpenImConversationName, conversation.Start); err != nil {
		panic(err.Error())
	}
}
