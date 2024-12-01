package logic

import (
	"context"
	"errors"
	"fmt"
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
	InvalidData           = response.MsgCode{10001, "无效数据"}
	codeUsersFoundField   = response.MsgCode{40031, "用户列表查询失败"}
	codeMemberCreateField = response.MsgCode{40032, "团队成员新增失败"}
	codeMemberDeleteField = response.MsgCode{40033, "团队成员删除失败"}
	codeMemberChangeField = response.MsgCode{40034, "团队成员修改失败"}
)

type MemberLogic struct {
}

func NewMemberLogic() *MemberLogic {
	return &MemberLogic{}
}

// GetMemberList
//
//	@Description: 查看用户基本信息
//	@receiver l
//	@param ctx
//	@param req
//	@return *types.MemberlistResp
//	@return error
func (l *MemberLogic) GetMemberList(ctx context.Context, req types.MemberlistReq) (*types.MemberlistResp, error) {
	defer util.RecordTime(time.Now())()

	users, err := repo.NewMemberRepo(global.DB).GetMemberlistRepo(req.TeamID, req.Page, req.Perpage)
	if err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			zlog.CtxWarnf(ctx, "无效数据: %v", err)
			return nil, response.ErrResp(err, InvalidData)
		} else {
			zlog.Errorf("用户列表查询失败: %v", err)
			return nil, response.ErrResp(err, codeUsersFoundField)
		}
	}

	return &users, nil
}

// GetMemberDetail
//
//	@Description: 查看用户详细信息
//	@receiver l
//	@param ctx
//	@param req
//	@return resp
//	@return err
func (l *MemberLogic) GetMemberDetail(ctx context.Context, UserId, MemberId int64) (resp *types.GetMemberDetailResp, err error) {
	defer util.RecordTime(time.Now())()

	resp, err = repo.NewMemberRepo(global.DB).GetMemberDetail(UserId, MemberId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "user not found: %v", err)
			return nil, response.ErrResp(err, codeUserNotFound)
		} else {
			zlog.Errorf("get user error: %v", err)
			return nil, response.ErrResp(err, codeUserFoundField)
		}

	}
	return
}

var (
	codeServerError     = response.MsgCode{Code: 60003, Msg: "上锁失败"}
	codeOperationLocked = response.MsgCode{Code: 20013, Msg: "该操作被加了锁"}
)

type LikeCountLogic struct{}

func NewLikeCountLogic() *LikeCountLogic {
	return &LikeCountLogic{}
}

// PutLikeCount
//
//	@Description: 用户点赞/取消赞
//	@receiver l
//	@param ctx
//	@param UserId
//	@param MemberId
//	@return resp
//	@return err
func (l *LikeCountLogic) PutLikeCount(ctx context.Context, UserId, MemberId int64) (resp *types.LikeCountResp, err error) {
	defer util.RecordTime(time.Now())()

	// 用 redis 加锁
	lockKey := fmt.Sprintf("like:lock:user:%d:member:%d", UserId, MemberId)
	// 尝试获取锁   返回布尔值 true:成功获取锁  false:上锁失败
	locked, err := global.Rdb.SetNX(ctx, lockKey, 1, 1*time.Second).Result()
	if err != nil {
		zlog.CtxErrorf(ctx, "Redis 上锁失败: %v", err)
		return nil, response.ErrResp(err, codeServerError)
	}
	if !locked {
		// 未获取到锁，说明该操作正在被其他请求处理
		zlog.CtxInfof(ctx, "点赞/取消赞操作正被 user: %d, member: %d  使用，请稍等 1 s", UserId, MemberId)
		return nil, response.ErrResp(err, codeOperationLocked)
	}
	//释放锁
	defer global.Rdb.Del(ctx, lockKey)

	resp, err = repo.NewLikeCountRepo(global.DB).PutLikeCount(UserId, MemberId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.CtxWarnf(ctx, "找不到用户: %v", err)
			return nil, response.ErrResp(err, codeUserNotFound)
		} else {
			zlog.Errorf("用户查询失败: %v", err)
			return nil, response.ErrResp(err, codeUserFoundField)
		}
	}
	return
}

// CreateMember
//
//	@Description: 新增团队成员
//	@receiver l
//	@param ctx
//	@param req
//	@return *types.CreateMembersResp
//	@return error
func (l *MemberLogic) CreateMember(ctx context.Context, req types.CreateMemberReq) (*types.CreateMembersResp, error) {
	defer util.RecordTime(time.Now())()

	err := repo.NewMemberRepo(global.DB).CreateMember(req)
	if err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			zlog.CtxWarnf(ctx, "无效数据: %v", err)
			return nil, response.ErrResp(err, InvalidData)
		} else {
			zlog.Errorf("创建团队成员失败: %v", err)
			return nil, response.ErrResp(err, codeMemberCreateField)
		}

	}
	return &types.CreateMembersResp{}, nil

}

// DeleteMember
//
//	@Description: 删除该成员在这个团队的关系
//	@receiver l
//	@param ctx
//	@param req
//	@return *types.DeleteMembersResp
//	@return error
func (l *MemberLogic) DeleteMember(ctx context.Context, req types.DeleteMemberReq) (*types.DeleteMembersResp, error) {
	defer util.RecordTime(time.Now())()

	err := repo.NewMemberRepo(global.DB).DeleteMember(req.MemberId, req.TeamId)
	if err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			zlog.CtxWarnf(ctx, "无效数据: %v", err)
			return nil, response.ErrResp(err, InvalidData)
		} else {
			zlog.Errorf("删除团队成员失败: %v", err)
			return nil, response.ErrResp(err, codeMemberDeleteField)
		}

	}
	return &types.DeleteMembersResp{}, nil

}

// PutMember
//
//	@Description: 更改成员信息
//	@receiver l
//	@param ctx
//	@param req
//	@return *types.PutTeamMemberResp
//	@return error
func (l *MemberLogic) PutMember(ctx context.Context, req types.PutTeamMemberReq) (*types.PutTeamMemberResp, error) {
	defer util.RecordTime(time.Now())()

	err := repo.NewMemberRepo(global.DB).PutMember(req)
	if err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			zlog.CtxWarnf(ctx, "无效数据: %v", err)
			return nil, response.ErrResp(err, InvalidData)
		} else {
			zlog.Errorf("修改团队成员失败: %v", err)
			return nil, response.ErrResp(err, codeMemberChangeField)
		}

	}
	return &types.PutTeamMemberResp{}, nil
}
