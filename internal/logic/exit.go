package logic

import (
	"context"
	"tgwp/global"
	"tgwp/internal/repo"
	"tgwp/internal/response"
	"tgwp/internal/types"
	"tgwp/util"
)

type ExitLogic struct {
}

func NewExitLogic() *ExitLogic {
	return &ExitLogic{}
}

// ExitSystem
//
//	@Description: 用于用户自己退出
//	@receiver l
//	@param ctx
//	@param req
//	@return err
func (l *ExitLogic) ExitSystem(ctx context.Context, req types.TokenReq) (err error) {
	//解析token是否有效，并取出上一次的值
	data, err := util.IdentifyToken(ctx, req.Token)
	if err != nil {
		//对应token无效，直接让他返回
		return response.ErrResp(err, response.TOKEN_IS_EXPIRED)
	}
	//判断token的签名
	if data.Issuer == "" {
		//此时用户不是自动登录的,前端回收token
		return nil
	} else {
		//自动登录的用户就去签名表删除对应数据
		err = repo.NewSignRepo(global.DB).DeleteSign(data.Issuer)
		if err != nil {
			//表明找不到issuer相等的，即rtoken是无效的
			return response.ErrResp(err, response.COMMON_FAIL)
		}
	}
	return nil
}
