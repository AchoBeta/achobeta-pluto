package repo

import (
	"fmt"
	"gorm.io/gorm"
	"tgwp/internal/model"
	"tgwp/internal/types"
	"tgwp/log/zlog"
	"time"
)

const c_NodeName = "node_name"

type StructureRepo struct {
	DB *gorm.DB
}

func NewStructureRepo(db *gorm.DB) *StructureRepo {
	return &StructureRepo{DB: db}
}

type MyNode struct {
	MyselfId int64  `gorm:"column:id"`
	NodeName string `gorm:"column:node_name"`
}

// GetNode
//
//	@Description:  获取架构节点
//	@receiver r
//	@param fatherid
//	@param teamid
//	@return []MyNode
//	@return error
func (r StructureRepo) GetNode(fatherid, teamid int64) ([]MyNode, error) {

	var mynode []MyNode
	err := r.DB.Model(&model.Structure{}).
		Select(C_Id, c_NodeName).
		Where(&model.Structure{
			FatherId: fatherid,
			TeamId:   teamid,
		}).
		Find(&mynode).
		Error
	if err != nil {
		zlog.Errorf("获取当前节点失败：%v", err)
		return nil, err
	}
	return mynode, nil
}

// InsertNode
//
//	@Description: 新增架构节点
//	@receiver r
//	@param Node
//	@return error
func (r StructureRepo) CreateNode(Node types.TeamStructure) error {

	var node = model.Structure{
		CommonModel: model.CommonModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		TeamId:   Node.TeamId,
		FatherId: Node.FatherId,
		NodeName: Node.NodeName,
	}

	err := r.DB.Model(&model.Structure{}).
		Create(&node).
		Error
	return err
}

// DeleteNode
//
//	@Description: 删除架构节点
//	@receiver r
//	@param Node
//	@return error
func (r StructureRepo) DeleteNode(Node types.TeamStructure) error {

	var node = model.Structure{
		CommonModel: model.CommonModel{
			ID: Node.MyselfId},
		TeamId:   Node.TeamId,
		FatherId: Node.FatherId,
		NodeName: Node.NodeName,
	}

	err := r.DB.Model(&model.Structure{}).
		Where(fmt.Sprintf("%s = ?", C_Id), Node.MyselfId).
		Delete(&node).
		Error
	return err
}
