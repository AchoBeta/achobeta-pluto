package handler

import (
	"context"
	"fmt"
	"tgwp/global"
	"tgwp/internal/response"
	"tgwp/log/zlog"
	"tgwp/util"
	"time"
)

var (
	codeHasSend = response.MsgCode{Code: 20007, Msg: "该手机号在一分钟内已经发送给验证码"}
)

// PostCode
//
//	@Description: 用作验证码处理
//	@param ctx
//	@param phone
//	@return error
func PostCode(ctx context.Context, phone string) (err error) {
	if !AccessCode(ctx, phone) {
		zlog.CtxErrorf(ctx, "Access code error")
		return response.ErrResp(err, codeHasSend)
	}
	//生成6位数字的验证码
	code := util.RandomCode()
	text := fmt.Sprintf("你的验证码是%s", code)
	//将验证码放入redis5分钟
	err = global.Rdb.Set(ctx, fmt.Sprintf(global.REDIS_PHONE_CODE, phone), code, time.Second*300).Err()
	if err != nil {
		zlog.CtxErrorf(ctx, "Store the verification code err: %v", err)
		return response.ErrResp(err, response.COMMON_FAIL)
	}
	//防刷处理
	err = global.Rdb.Set(ctx, fmt.Sprintf(global.REDIS_PHONE, phone), 0, time.Second*60).Err()
	if err != nil {
		zlog.CtxErrorf(ctx, "Restrict multiple access err: %v", err)
		return response.ErrResp(err, response.COMMON_FAIL)
	}
	//发送验证码到用户手机
	if err := PostMessage(text, phone); err != nil {
		zlog.CtxErrorf(ctx, "PostMessage err: %v", err)
		return response.ErrResp(err, response.CONNECT_PHONE_ERROR)
	}
	return
}

// PostMessage
//
//	@Description: 发送验证码到用户手机
//	@param text
//	@param phone
//	@return error
func PostMessage(text string, phone string) error {
	return nil
}

// CompareCode
//
//	@Description: 对比验证码是否有效
//	@param ctx
//	@param code
//	@param phone
//	@return bool
func CompareCode(ctx context.Context, code, phone string) bool {
	return code == global.Rdb.Get(ctx, fmt.Sprintf(global.REDIS_PHONE_CODE, phone)).Val()
}

// AccessCode
//
//	@Description: 判断手机是否在一分钟内已经发过验证码
//	@param ctx
//	@param phone
//	@return bool
func AccessCode(ctx context.Context, phone string) bool {
	//存在，即手机号一分钟内被记录
	if exist := global.Rdb.Exists(ctx, fmt.Sprintf(global.REDIS_PHONE, phone)).Val(); exist == 1 {
		return false
	}
	return true
}
