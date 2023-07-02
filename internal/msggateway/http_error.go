package msggateway

import "github.com/xiaoyiEdu/Open-IM-Server/pkg/apiresp"

func httpError(ctx *UserConnContext, err error) {
	apiresp.HttpError(ctx.RespWriter, err)
}
