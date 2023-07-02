package cmd

import (
	"github.com/xiaoyiEdu/Open-IM-Server/internal/msggateway"
	//"github.com/xiaoyiEdu/Open-IM-Server/internal/msggateway"
	"github.com/spf13/cobra"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/constant"
)

type MsgGatewayCmd struct {
	*RootCmd
}

func NewMsgGatewayCmd() MsgGatewayCmd {
	return MsgGatewayCmd{NewRootCmd("msgGateway")}
}

func (m *MsgGatewayCmd) AddWsPortFlag() {
	m.Command.Flags().IntP(constant.FlagWsPort, "w", 0, "ws server listen port")
}

func (m *MsgGatewayCmd) getWsPortFlag(cmd *cobra.Command) int {
	port, _ := cmd.Flags().GetInt(constant.FlagWsPort)
	return port
}

func (m *MsgGatewayCmd) addRunE() {
	m.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return msggateway.RunWsAndServer(m.getPortFlag(cmd), m.getWsPortFlag(cmd), m.getPrometheusPortFlag(cmd))
	}
}

func (m *MsgGatewayCmd) Exec() error {
	m.addRunE()
	return m.Execute()
}
