package logic

import (
	"context"
	"tgwp/global"
	"tgwp/internal/repo"
	"tgwp/internal/types"
	"tgwp/log/zlog"
	"tgwp/util"
	"time"
)

type CasbinLogic struct {
}

func NewCasbinLogic() *CasbinLogic {
	return &CasbinLogic{}
}

// CasbinLogic
//
//	@Description:获取权限组
//	@receiver l
//	@param ctx
//	@param req
//	@return resp
//	@return err
func (l *CasbinLogic) CasbinLogic(ctx context.Context, UserId, TeamId int64) (resp types.RuleResp, err error) {
	defer util.RecordTime(time.Now())()
	// 前端没有传团队id 时 仅返回第一个团队ID，所有的团队ID，状态码code，信息获取（成功/失败）msg
	if TeamId != 0 {
		//获取 该用户在该团队所拥有的 特殊 url
		Level, urls, err := repo.NewCasbinRepo(global.DB).GetCasbin(UserId, TeamId)
		if err != nil {
			zlog.CtxErrorf(ctx, "%v", err)
		}
		//把找出来的url给出参 Url
		resp.Url = urls
		resp.Level = Level
	}

	//获取 第一个团队ID 和 名称，所有的团队ID 和名称
	FirstTeam, Team, err := repo.NewTeamIdRepo(global.DB).GetTeamId(UserId)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v", err)
	}
	resp.FirstTeamID = FirstTeam.TeamId
	resp.FirstTeamName = FirstTeam.TeamName
	resp.Team = Team

	return
}
