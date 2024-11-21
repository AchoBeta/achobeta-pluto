package logic

import (
	"context"
	"fmt"
	"tgwp/internal/types"
	"tgwp/log/zlog"
	"tgwp/util"
)

type FeiShuLogic struct {
}

func NewFeiShuLogic() *FeiShuLogic {
	return &FeiShuLogic{}
}

// GetFeiShuList 获取飞书多维表格
func (l *FeiShuLogic) GetFeiShuList(ctx context.Context, UserID int64, ForceUpdate bool) (resp types.GetFeiShuListResp, err error) {
	// 获取用户飞书open_id
	openID, err := util.GetFeiShuUserOpenID("18300156621") // 由于通过UserID获取手机号需要个人信息模块那边的函数，此处暂时使用测试手机号
	if err != nil {
		zlog.Errorf("get feishu open_id error:%v", err)
		return
	}
	fmt.Println(openID)

	// 获取多维表格数据
	resp, err = util.GetFeiShuList(ctx, openID, ForceUpdate)

	if err != nil {
		zlog.Errorf("get message error:%v", err)
		return
	} else {
		zlog.Infof("get message success")
	}

	return
}
