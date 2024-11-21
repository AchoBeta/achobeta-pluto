package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tgwp/global"
	"tgwp/internal/repo"
	"tgwp/internal/response"
	"tgwp/internal/types"
	"tgwp/log/zlog"
	"tgwp/util"
	"time"
)

var (
	codeNodeNotFound    = response.MsgCode{Code: 40027, Msg: "未找到节点"}
	codeNodeCreateField = response.MsgCode{Code: 40028, Msg: "新增节点失败"}
	codeNodeDeleteField = response.MsgCode{Code: 40029, Msg: "删除节点失败"}
)

type TeamNodeLogic struct {
}

func NewTeamNodeLogic() *TeamNodeLogic {
	return &TeamNodeLogic{}
}

// TeamNodeLogic
//
//	@Description:  保存新增的节点，删除被删除的节点
//	@receiver l
//	@param ctx
//	@param req
//	@return types.PutTeamNodeResp
//	@return error
func (l *TeamNodeLogic) TeamNodeLogic(ctx context.Context, req types.PutTeamNodeReq) (*types.PutTeamNodeResp, error) {
	defer util.RecordTime(time.Now())()

	for _, Node := range req.TeamStructures {
		if Node.IsDeleted == false {
			// 没被删除且没有自身节点值 ：新增节点
			err := repo.NewStructureRepo(global.DB).CreateNode(Node)
			if err != nil {
				zlog.CtxErrorf(ctx, "insert Node error: %v", err)
				return nil, response.ErrResp(err, codeNodeCreateField)
			}
		} else {
			// true ：被删除的节点
			err := repo.NewStructureRepo(global.DB).DeleteNode(Node)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					zlog.CtxWarnf(ctx, "Node not found: %v", err)
					return nil, response.ErrResp(err, codeNodeNotFound)
				} else {
					zlog.CtxErrorf(ctx, "delete Node error: %v", err)
					return nil, response.ErrResp(err, codeNodeDeleteField)
				}
			}
		}
	}

	return &types.PutTeamNodeResp{}, nil
}
