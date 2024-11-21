package logic

import (
	"context"
	"tgwp/global"
	"tgwp/internal/repo"
	"tgwp/internal/response"
	"tgwp/internal/types"
	"tgwp/log/zlog"
	"tgwp/util"
	"time"
)

type TokenLogic struct {
}

func NewTokenLogic() *TokenLogic {
	return &TokenLogic{}
}

// GenAtoken
//
//	@Description: 生成atoken
//	@receiver l
//	@param ctx
//	@param data
//	@return resp
//	@return err
func (l *TokenLogic) GenAtoken(ctx context.Context, data util.TokenData) (resp types.TokenResp, err error) {
	defer util.RecordTime(time.Now())()
	resp.Atoken, err = util.GenToken(data)
	if err != nil {
		zlog.CtxErrorf(ctx, "AtokenLogic err:%v", err)
		return
	}
	return
}

// GenRtoken
//
//	@Description: 生成rtoken
//	@receiver l
//	@param ctx
//	@param data
//	@return resp
//	@return err
func (l *TokenLogic) GenRtoken(ctx context.Context, data util.TokenData) (resp types.TokenResp, err error) {
	defer util.RecordTime(time.Now())()
	resp.Rtoken, err = util.GenToken(data)
	if err != nil {
		zlog.CtxErrorf(ctx, "TokenLogic err:%v", err)
	}
	data.Time = global.ATOKEN_EFFECTIVE_TIME
	data.Class = global.AUTH_ENUMS_ATOKEN
	resp.Atoken, err = util.GenToken(data)
	if err != nil {
		zlog.CtxErrorf(ctx, "RtokenLogic err:%v", err)
		return
	}
	return
}

// ReflashRtoken
//
//	@Description: 用于rtoken刷新,同时可以用于验证是否自动登录
//	@receiver l
//	@param ctx
//	@param req
//	@return resp
//	@return err
func (l *TokenLogic) ReflashRtoken(ctx context.Context, req types.TokenReq) (resp types.TokenResp, err error) {
	//解析token是否有效，并取出上一次的值
	data, err := util.IdentifyToken(ctx, req.Token)
	if err != nil {
		//对应token无效，直接让他返回
		return resp, response.ErrResp(err, response.TOKEN_IS_EXPIRED)
	}
	//判断其是否为rtoken
	if data.Class != global.AUTH_ENUMS_RTOKEN {
		return resp, response.ErrResp(err, response.PARAM_TYPE_ERROR)
	}
	//判断rtoken的签名是否有效
	err = repo.NewSignRepo(global.DB).CompareSign(data.Issuer)
	if err != nil {
		//表明找不到issuer相等的，即rtoken是无效的
		return resp, response.ErrResp(err, response.PARAM_NOT_VALID)
	}
	//生成新的token
	resp, err = NewTokenLogic().GenRtoken(ctx, data)
	if err != nil {
		return resp, response.ErrResp(err, response.COMMON_FAIL)
	}
	return resp, nil
}
