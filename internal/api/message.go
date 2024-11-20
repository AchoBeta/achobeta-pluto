package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/global"
	"tgwp/internal/logic"
	"tgwp/internal/response"
	"tgwp/internal/types"
	"tgwp/log/zlog"
)

// SetMessage 设置消息, 存入 [消息表]
func SetMessage(c *gin.Context) {
	// 解析请求参数
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.SetMessageReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "BindReq failed: %v", err)
		err = response.ErrResp(err, response.INTERNAL_ERROR)
		return
	}
	zlog.CtxInfof(ctx, "Casbin request: %v", req)

	// logic 层处理
	resp, err := logic.NewMessageLogic().SetMessage(req)

	// 响应
	if err != nil {
		response.Response(c, resp, err)
		return
	} else {
		response.NewResponse(c).Success(resp)
	}
	return
}

// JoinMessage 连接用户与消息, 存入 [用户-消息表]
func JoinMessage(c *gin.Context) {
	// 解析请求参数
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.JoinMessageReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "BindReq failed: %v", err)
		err = response.ErrResp(err, response.INTERNAL_ERROR)
		return
	}
	TempUserID, _ := c.Get(global.TOKEN_USER_ID)
	UserID := TempUserID.(int64)

	zlog.CtxInfof(ctx, "Casbin request: %v", req)

	// logic 层处理
	resp, err := logic.NewMessageLogic().JoinMessage(req, UserID)

	// 响应
	if err != nil {
		response.Response(c, resp, err)
		return
	} else {
		response.NewResponse(c).Success(resp)
	}

	return
}

// GetMessage 获取消息, 从 [消息表] 获取
func GetMessage(c *gin.Context) {
	//fmt.Println(util.GenToken(util.TokenData{Userid: 114514, Class: "atoken", Issuer: "", Time: time.Hour * 24 * 365}))
	// 解析请求参数
	ctx := zlog.GetCtxFromGin(c)
	TempUserID, _ := c.Get(global.TOKEN_USER_ID)
	UserID := TempUserID.(int64)

	pageStr := c.DefaultQuery("page", "1")
	timestampStr := c.DefaultQuery("timestamp", "0")

	zlog.CtxInfof(ctx, "Casbin request: %v %v %v", UserID, pageStr, timestampStr)

	// logic 层处理
	resp, err := logic.NewMessageLogic().GetMessage(UserID, pageStr, timestampStr)

	// 响应
	if err != nil {
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	} else {
		response.NewResponse(c).Success(resp)
	}

	return
}

// MarkReadMessage 标记已读, 更新 [用户-消息表] 的 read_at 字段
func MarkReadMessage(c *gin.Context) {
	// 解析请求参数
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.MarkReadMessageReq](c)
	if err != nil {
		zlog.Errorf("BindReq failed: %v", err)
		err = response.ErrResp(err, response.INTERNAL_ERROR)
		return
	}
	zlog.CtxInfof(ctx, "Casbin request: %v", req)

	// logic 层处理
	resp, err := logic.NewMessageLogic().MarkReadMessage(req)

	// 响应
	if err != nil {
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	} else {
		response.NewResponse(c).Success(resp)
	}

	return
}

// SendMessage 一键发送消息，SetMessage 和 JoinMessage 合并运用
func SendMessage(c *gin.Context) {
	// 解析请求参数
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.SendMessageReq](c)
	if err != nil {
		zlog.CtxErrorf(ctx, "BindReq failed: %v", err)
		err = response.ErrResp(err, response.INTERNAL_ERROR)
		return
	}
	TempUserID, _ := c.Get(global.TOKEN_USER_ID)
	UserID := TempUserID.(int64)

	zlog.CtxInfof(ctx, "Casbin request: %v", req)

	// logic 层处理
	resp, err := logic.NewMessageLogic().SendMessage(req, UserID)

	// 响应
	if err != nil {
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	} else {
		response.NewResponse(c).Success(resp)
	}

	return
}
