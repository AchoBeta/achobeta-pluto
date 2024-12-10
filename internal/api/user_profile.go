package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/internal/logic"
	"tgwp/internal/response"
	"tgwp/log/zlog"
)

// GetUserDetail
//
//	@Description: 获取用户自己的详情信息
//	@param c
func GetUserDetail(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)

	// 获取用户id
	UserID := logic.GetUserId(c)

	// 直接使用团队用户获取详细信息的函数
	resp, err := logic.NewMemberLogic().GetMemberDetail(ctx, UserID, UserID)
	response.Response(c, resp, err)

	return
}

// PutUserLikeCount
//
//	@Description: 用户给自己点赞
//	@param c
func PutUserLikeCount(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)

	// 获取用户id
	UserID := logic.GetUserId(c)

	// 直接使用团队用户点赞的函数
	resp, err := logic.NewLikeCountLogic().PutLikeCount(ctx, UserID, UserID)
	response.Response(c, resp, err)

	return
}
