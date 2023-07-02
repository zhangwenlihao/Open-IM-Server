package main

import (
	"github.com/xiaoyiEdu/Open-IM-Server/internal/tools"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/cmd"
)

func main() {
	cronTaskCmd := cmd.NewCronTaskCmd()
	if err := cronTaskCmd.Exec(tools.StartCronTask); err != nil {
		panic(err.Error())
	}
}
