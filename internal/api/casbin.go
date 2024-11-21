package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/global"
	"tgwp/internal/logic"
	"tgwp/internal/response"
	"tgwp/internal/types"
	"tgwp/log/zlog"
)

// GetPower
//
//	@Description: 获得权限组
//	@param c
func GetPower(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)

	userid, exists := c.Get(global.TOKEN_USER_ID)
	if !exists {
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	UserId := userid.(int64)
	req, err := types.BindReq[types.RuleReq](c)

	if err != nil {
		zlog.CtxErrorf(ctx, "GetPower err:%v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "GetPower request: %v", req)
	//获取出参
	resp, err := logic.NewCasbinLogic().CasbinLogic(ctx, UserId, req.TeamId)

	if err != nil {
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	} else {
		response.NewResponse(c).Success(resp)
	}
	return
}
