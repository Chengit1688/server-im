package http

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
	"im/pkg/code"
	"net/http"
)

type Resp struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(ctx *gin.Context, payload ...interface{}) {
	resp := &Resp{}
	resp.Code = 0
	resp.Message = "success"

	switch len(payload) {
	case 0:
		resp.Data = struct{}{}
	default:
		resp.Data = payload[0]
	}
	ctx.JSON(http.StatusOK, resp)
}

func Failed(ctx *gin.Context, err error) {
	resp := &Resp{}

	var res *status.Status
	if err != nil {
		res, _ = status.FromError(err)
	} else {
		res, _ = status.FromError(code.ErrUnknown)

	}

	resp.Code = int32(res.Code())
	resp.Message = res.Message()
	resp.Data = struct{}{}
	ctx.JSON(http.StatusOK, resp)
}
