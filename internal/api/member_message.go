package api

import (
	"github.com/gin-gonic/gin"
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

	//正式使用，测试时需注释掉
	UserId := logic.GetUserId(c)

	req, err := types.BindReq[types.GetMemberDetailReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "GetMemberDetail request: %v", req)
	resp, err := logic.NewMemberLogic().GetMemberDetail(ctx, UserId, req.MemberID)
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

	//正式使用，测试时需注释掉
	UserId := logic.GetUserId(c)

	req, err := types.BindReq[types.LikeCountReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "PutLikeCount request: %v", req)
	resp, err := logic.NewLikeCountLogic().PutLikeCount(ctx, UserId, req.MemberID) //正式使用
	//resp, err := logic.NewLikeCountLogic().PutLikeCount(ctx, req.UserID, req.MemberID) //测试时使用
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
		zlog.CtxErrorf(ctx, "绑定请求参数失败: %v", err)
		return
	}
	zlog.CtxInfof(ctx, "PutTeamMember request: %v", req)
	resp, err := logic.NewMemberLogic().PutMember(ctx, req)
	response.Response(c, resp, err)

	return
}
