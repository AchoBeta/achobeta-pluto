package logic

import (
	"context"
	"tgwp/internal/types"
	"tgwp/util"
	"time"
)

type DevicesLogic struct {
}

func NewDevicesLogic() *DevicesLogic {
	return &DevicesLogic{}
}

// ShowDevices
//
//	@Description: 展示常用设备
//	@receiver l
//	@param ctx
//	@param data
//	@return resp
//	@return err
func (l *DevicesLogic) ShowDevices(ctx context.Context, req types.DevicesReq) (resp types.PhoneResp, err error) {
	defer util.RecordTime(time.Now())()
	return
}
