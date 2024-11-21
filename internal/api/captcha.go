package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/internal/logic"
	"tgwp/internal/response"
	"tgwp/internal/types"
	"tgwp/log/zlog"
)

// GetCode
//
//	@Description: 获取验证码
//	@param c
func GetCode(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.PhoneReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "GetCode request: %v", req)
	err = logic.NewCodeLogic().GenCode(ctx, req)
	response.Response(c, nil, err)
	return
}

// LoginWithCode
//
//	@Description: 用验证码登录
//	@param c
func LoginWithCode(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.PhoneReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "LoginWithCode request: %v", req)
	resp, err := logic.NewCodeLogic().GenLoginData(ctx, req, c.ClientIP(), c.Request.UserAgent())
	response.Response(c, resp, err)
	return
}
