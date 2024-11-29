package logic

import (
	"context"
	"fmt"
	"tgwp/global"
	"tgwp/internal/model"
	"tgwp/internal/repo"
	"tgwp/internal/response"
	"tgwp/internal/types"
	"tgwp/log/zlog"
	"tgwp/util"
	"tgwp/util/snowflake"
	"time"
)

type CodeLogic struct {
}

func NewCodeLogic() *CodeLogic {
	return &CodeLogic{}
}

// GenCode
//
//	@Description: 生成验证码，发送到用户手机
//	@receiver l
//	@param ctx
//	@param req
//	@return err
func (l *CodeLogic) GenCode(ctx context.Context, req types.PhoneReq) (err error) {
	defer util.RecordTime(time.Now())()
	//校验手机号
	if flag := util.IndetifyPhone(req.Phone); !flag {
		return response.ErrResp(err, response.PHONE_ERROR)
	}
	//生成随机验证码并发送到对应用户
	err = PostCode(ctx, req.Phone)
	if err != nil {
		zlog.CtxErrorf(ctx, "GenCode err: %v", err)
		return
	}
	return
}

// GenLoginData
//
//	@Description: 为登陆后的用户授予一些个人信息
//	@receiver l
//	@param ctx
//	@param AutoLogin
//	@param resp
//	@return err
func (l *CodeLogic) GenLoginData(ctx context.Context, req types.PhoneReq, ip, user_agent string) (resp types.PhoneResp, err error) {
	defer util.RecordTime(time.Now())()
	node, err := snowflake.NewNode(global.DEFAULT_NODE_ID)
	if err != nil {
		zlog.CtxErrorf(ctx, "NewNode err: %v", err)
		return resp, response.ErrResp(err, response.COMMON_FAIL)
	}
	if !CompareCode(ctx, req.Code, req.Phone) {
		return resp, response.ErrResp(err, response.CAPTCHA_ERROR)
	}
	var user_id int64
	user_id, resp.IsTeam, err = repo.NewMemberRepo(global.DB).JudgeUser(req.Phone)
	if err != nil {
		zlog.CtxErrorf(ctx, "GenLoginData err: %v", err)
		//这里的err只可能是数据库出现错误了
		return resp, response.ErrResp(err, response.COMMON_FAIL)
	}
	if !resp.IsTeam {
		//非团队人员直接返回
		return resp, err
	}
	resp.Ip = ip
	resp.UserAgent = user_agent
	if req.AutoLogin {
		issuer := snowflake.GetString12Id(node)
		resp.Atoken, err = util.GenToken(util.FullToken(global.AUTH_ENUMS_ATOKEN, issuer, user_id))
		resp.Rtoken, err = util.GenToken(util.FullToken(global.AUTH_ENUMS_RTOKEN, issuer, user_id))
		//将点了自动登录的用户的信息插入签名表
		data := model.Sign{
			UserId:     user_id,
			Issuer:     issuer,
			OnlineTime: time.Now(),
			IP:         resp.Ip,
			UserAgent:  resp.UserAgent,
		}
		//由于这个id只是用于移除常用设备，和填充常用设备的名字，所以也是只有自动登陆的有
		resp.Id, err = repo.NewSignRepo(global.DB).InsertSign(data)
		if err != nil {
			zlog.CtxErrorf(ctx, "InsertSign err: %v", err)
			return resp, response.ErrResp(err, response.COMMON_FAIL)
		}
	} else {
		issuer := "" //没有自动登录，做置空处理
		resp.Atoken, err = util.GenToken(util.FullToken(global.AUTH_ENUMS_ATOKEN, issuer, user_id))
	}
	return
}

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
	//code := util.RandomCode()
	code := 123456
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
