package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/internal/logic"
	"tgwp/internal/response"
	"tgwp/internal/types"
	"tgwp/log/zlog"
)

// GetTeamMemberlist
//
//	@Description: 查询团队列表--用户基础信息
//	@param c
func GetTeamMemberlist(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)

	req, err := types.BindReq[types.MemberlistReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "GetTeamMemberlist err:%v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "GetTeamMemberlist request: %v", req)
	resp, err := logic.NewMemberListic().GetMemberList(ctx, req)
	response.Response(c, resp, err)

	return
}

// CreateTeamMember
//
//	@Description: 新增用户
//	@param c
func CreateTeamMember(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)

	req, err := types.BindReq[types.CreateMemberReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "Create TeamMember err:%v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "Create TeamMember request: %v", req)
	resp, err := logic.NewCreateMemberLogic().CreateMember(ctx, req)
	response.Response(c, resp, err)

	return
}

// DeleteTeamMember
//
//	@Description: 删除用户
//	@param c
func DeleteTeamMember(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)

	req, err := types.BindReq[types.DeleteMemberReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "Delete TeamMember err:%v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "Delete TeamMember request: %v", req)
	resp, err := logic.NewDeleteMemberLogic().DeleteMember(ctx, req)
	response.Response(c, resp, err)

	return
}
