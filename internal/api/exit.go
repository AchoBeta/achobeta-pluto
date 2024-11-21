package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/internal/logic"
	"tgwp/internal/response"
	"tgwp/internal/types"
	"tgwp/log/zlog"
)

// ExitSystem
//
//	@Description: 用户自己退出接口
//	@param c
func ExitSystem(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindJson[types.TokenReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "ExitSystem request: %v", req)
	err = logic.NewExitLogic().ExitSystem(ctx, req)
	response.Response(c, nil, err)
}
