package a2r

import (
	"context"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/checker"

	"github.com/xiaoyiEdu/Open-IM-Server/pkg/apiresp"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/log"
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/errs"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func Call[A, B, C any](
	rpc func(client C, ctx context.Context, req *A, options ...grpc.CallOption) (*B, error),
	client C,
	c *gin.Context,
) {
	var req A
	if err := c.BindJSON(&req); err != nil {
		log.ZWarn(c, "gin bind json error", err, "req", req)
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap()) // 参数错误
		return
	}
	if err := checker.Validate(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.Wrap(err.Error())) // 参数校验失败
		return
	}
	data, err := rpc(client, c, &req)
	if err != nil {
		apiresp.GinError(c, err) // RPC调用失败
		return
	}
	apiresp.GinSuccess(c, data) // 成功
}
