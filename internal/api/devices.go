package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/internal/logic"
	"tgwp/internal/response"
	"tgwp/internal/types"
	"tgwp/log/zlog"
)

func ShowDevices(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.DevicesReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "ShowDevices request: %v", req)
	resp, err := logic.NewDevicesLogic().ShowDevices(ctx, req)
	response.Response(c, resp, err)
	return
}
