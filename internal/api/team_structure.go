package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/internal/logic"
	"tgwp/internal/response"
	"tgwp/internal/types"
	"tgwp/log/zlog"
)

// GetTeamStructure
//
//	@Description: 获取 完整团队架构
//	@param c
func GetTeamStructure(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.TeamStructReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetTeamStructure err:%v", err)
		return
	}
	zlog.CtxInfof(ctx, "GetTeamStructure request: %v", req)
	resp, err := logic.NewStructureLogic().GetStructure(ctx, req)

	response.Response(c, resp, err)

	return
}

// PutTeamNode
//
//	@Description: 保存 更改了的节点信息
//	@param c
func PutTeamNode(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.PutTeamNodeReq](c)
	zlog.CtxInfof(ctx, "PutTeamNode request: %v", req)
	resp, err := logic.NewStructureLogic().PutStructureNode(ctx, req)

	response.Response(c, resp, err)
	return
}

// CreateTeam
//
//	@Description: 新增团队
//	@param c
func CreateTeam(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.CreateTeamReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "CreateTeam err:%v", err)
		return
	}
	zlog.CtxInfof(ctx, "CreateTeam request: %v", req)
	resp, err := logic.NewTeamLogic().CreateTeam(ctx, req)

	response.Response(c, resp, err)

	return

}
