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

type StructureLogic struct {
}

func NewStructureLogic() *StructureLogic {
	return &StructureLogic{}
}

var (
	rootFoundField    = response.MsgCode{Code: 40023, Msg: "根节点查询失败"}
	rootNotFound      = response.MsgCode{Code: 40024, Msg: "根节点不存在"}
	childFoundField   = response.MsgCode{Code: 40025, Msg: "孩子节点查询失败"}
	childrootNotFound = response.MsgCode{Code: 40026, Msg: "孩子节点不存在"}
)

// StructureLogic
//
//	@Description:  获取 完整团队架构
//	@receiver l
//	@param ctx
//	@param req
//	@return types.TeamStructResp
//	@return error
func (l *StructureLogic) StructureLogic(ctx context.Context, req types.TeamStructReq) (types.TeamStructResp, error) {
	defer util.RecordTime(time.Now())()

	teamStructures := []types.TeamStructure{}

	//找到该团队的根节点
	root, err := repo.NewStructureRepo(global.DB).GetNode(global.ROOT_ID, req.TeamId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "root not found: %v", err)
			return types.TeamStructResp{}, response.ErrResp(err, rootFoundField)
		} else {
			zlog.Errorf("get root error: %v", err)
			return types.TeamStructResp{}, response.ErrResp(err, rootNotFound)
		}
	}
	Root := root[0].MyselfId

	// 递归获取节点信息
	err = l.getStructure(ctx, Root, req.TeamId, &teamStructures)
	if err != nil {
		zlog.CtxErrorf(ctx, "Failed to get children for fatherid: %d, teamid: %d, error: %v", Root, req.TeamId, err)
		return types.TeamStructResp{}, err
	}

	return types.TeamStructResp{TeamStructures: teamStructures}, nil
}

// getStructure
//
//	@Description:  递归获取节点信息
//	@receiver l
//	@param ctx
//	@param fatherid
//	@param teamid
//	@param result
//	@return error
func (l *StructureLogic) getStructure(ctx context.Context, fatherid, teamid int64, result *[]types.TeamStructure) error {
	// 获取当前节点的所有子节点
	children, err := repo.NewStructureRepo(global.DB).GetNode(fatherid, teamid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "child not found: %v", err)
			return response.ErrResp(err, childFoundField)
		} else {
			zlog.Errorf("get child error: %v", err)
			return response.ErrResp(err, childrootNotFound)
		}
	}

	for _, child := range children {
		node := types.TeamStructure{
			TeamId:    teamid,
			MyselfId:  child.MyselfId,
			FatherId:  fatherid,
			NodeName:  child.NodeName,
			IsDeleted: false, // 假设在查询时过滤了已删除的节点
		}

		*result = append(*result, node)

		// 递归获取子节点的子节点
		err = l.getStructure(ctx, child.MyselfId, teamid, result)
		if err != nil {
			zlog.CtxErrorf(ctx, "Failed to get children for fatherid: %d, teamid: %d, error: %v", fatherid, teamid, err)
			return err
		}
	}
	return nil
}
