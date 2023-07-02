package main

import (
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/cmd"
)

func main() {
	msgGatewayCmd := cmd.NewMsgGatewayCmd()
	msgGatewayCmd.AddWsPortFlag()
	msgGatewayCmd.AddPortFlag()
	msgGatewayCmd.AddPrometheusPortFlag()
	if err := msgGatewayCmd.Exec(); err != nil {
		panic(err.Error())
	}
}
