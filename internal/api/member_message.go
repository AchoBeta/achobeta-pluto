package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/global"
	"tgwp/internal/logic"
	"tgwp/internal/response"
	"tgwp/internal/types"
	"tgwp/log/zlog"
)

// GetMemberDetail
//
//	@Description: 查询用户详细信息
//	@param c
func GetMemberDetail(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)

	req, err := types.BindReq[types.GetMemberDetailReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "GetMemberDetail request: %v", req)
	resp, err := logic.NewMemberDetailLogic().GetMemberDetail(ctx, req)
	// 更加人性化的response返回，这样减少重复代码的书写
	response.Response(c, resp, err)

	return
}

// PutLikeCount
//
//	@Description: 用户点赞/取消赞
//	@param c
func PutLikeCount(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)

	userid, exists := c.Get(global.TOKEN_USER_ID)
	if !exists {
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	UserID := userid.(int64)

	req, err := types.BindReq[types.LikeCountReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "PutLikeCount request: %v", req)
	resp, err := logic.NewLikeCountLogic().PutLikeCount(ctx, UserID, req.MemberID)
	response.Response(c, resp, err)

	return
}

// PutTeamMember
//
//	@Description: 编辑成员信息
//	@param c
func PutTeamMember(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)

	req, err := types.BindReq[types.PutTeamMemberReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "PutTeamMember request: %v", req)
	resp, err := logic.NewPutMemberLogic().PutMember(ctx, req)
	response.Response(c, resp, err)

	return
}
