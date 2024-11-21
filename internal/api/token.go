package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/internal/logic"
	"tgwp/internal/response"
	"tgwp/internal/types"
	"tgwp/log/zlog"
)

// ReflashRtoken
//
//	@Description: 前端用rtoken刷新token
//	@param c
func ReflashRtoken(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.TokenReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "ReflashRtoken request: %v", req)
	resp, err := logic.NewTokenLogic().ReflashRtoken(ctx, req)
	response.Response(c, resp, err)
}
