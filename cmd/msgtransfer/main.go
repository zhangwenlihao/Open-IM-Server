package main

import (
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/cmd"
)

func main() {
	msgTransferCmd := cmd.NewMsgTransferCmd()
	msgTransferCmd.AddPrometheusPortFlag()
	if err := msgTransferCmd.Exec(); err != nil {
		panic(err.Error())
	}
}
