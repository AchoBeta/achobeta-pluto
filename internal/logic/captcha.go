package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tgwp/global"
	"tgwp/internal/handler"
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
	err = handler.PostCode(ctx, req.Phone)
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
	if !handler.CompareCode(ctx, req.Code, req.Phone) {
		return resp, response.ErrResp(err, response.CAPTCHA_ERROR)
	}
	resp.Ip = ip
	resp.UserAgent = user_agent
	resp.LoginId = snowflake.GenId(node)
	//这里做关于团队成员的判断，思凯会设计一个函数，我传手机号，看看他团队表内有么有
	resp.IsTeam = true //暂时认定都是团队成员
	//一个手机号对应的user_id是一样的
	//这里到时候外键关联用户表，userid就是逻辑外键，手机号也可以删掉了，但是现在不处理
	//这里生成的id现阶段是方便测试用
	user_id, err := repo.NewSignRepo(global.DB).CheckUserId(req.Phone)
	if err != nil {
		//这里的err是代表找不到对应的user_id,所以生成一个新的id
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user_id = node.Generate().Int64()
		} else {
			zlog.CtxErrorf(ctx, "CheckUserId err: %v", err)
			return resp, response.ErrResp(err, response.COMMON_FAIL)
		}
	}
	if req.AutoLogin {
		issuer := snowflake.GenId(node)
		resp.Atoken, err = util.GenToken(util.FullToken(global.AUTH_ENUMS_ATOKEN, issuer, user_id))
		resp.Rtoken, err = util.GenToken(util.FullToken(global.AUTH_ENUMS_RTOKEN, issuer, user_id))
		//将点了自动登录的用户的login_id,issuer插入签名表
		data := model.Sign{
			UserId:     user_id,
			Issuer:     issuer,
			OnlineTime: time.Now(),
			LoginId:    resp.LoginId,
			IP:         resp.Ip,
			UserAgent:  resp.UserAgent,
			Phone:      req.Phone,
		}
		err = repo.NewSignRepo(global.DB).InsertSign(data)
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
