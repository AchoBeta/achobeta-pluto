package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/internal/logic"
	"tgwp/internal/response"
	"tgwp/internal/types"
	"tgwp/log/zlog"
)

// CheckAutoLogin
//
//	@Description: 验证是否可以自动登录
//	@param c
func CheckAutoLogin(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.TokenReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "CheckAutoLogin request: %v", req)
	resp, err := logic.NewTokenLogic().ReflashRtoken(ctx, req)
	response.Response(c, resp, err)
}
