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

var (
	codeTeamCreateField = response.MsgCode{Code: 40030, Msg: "新增团队失败"}
)

// TeamLogic
// @Description:  创建团队
type TeamLogic struct {
}

func NewTeamLogic() *TeamLogic {
	return &TeamLogic{}
}

func (l *TeamLogic) TeamLogic(ctx context.Context, req types.CreateTeamReq) (resp *types.CreateTeamResp, err error) {
	defer util.RecordTime(time.Now())()
	resp, err = repo.NewCreateTeamRepo(global.DB).CreateTeam(req.Name)
	if err != nil {
		zlog.CtxErrorf(ctx, "创建团队失败: %v", err)
		return nil, response.ErrResp(err, codeTeamCreateField)
	}
	return resp, nil
}
