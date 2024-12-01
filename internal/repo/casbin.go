package repo

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"tgwp/internal/model"
	"tgwp/log/zlog"
	"tgwp/util"
	"time"
)

type CasbinRepo struct {
	DB *gorm.DB
}

func NewCasbinRepo(db *gorm.DB) *CasbinRepo {
	return &CasbinRepo{DB: db}
}

const NoPermission = 1

type RoleWithLevel struct {
	Role  int64
	Level int
}

// GetCasbin
//
//	@Description: 获取权限组
//	@receiver r
//	@param userid
//	@param teamid
//	@return int
//	@return []string
//	@return error
func (r CasbinRepo) GetCasbin(userid, teamid int64) (int, []string, error) {
	// 查询用户的角色，包括默认团队（team_id = 1）
	roles, err := r.queryRoles(userid, teamid)
	if err != nil {
		zlog.Errorf("查询用户 %d 在团队 %d 的角色失败：%v", userid, teamid, err)
		return NoPermission, nil, err
	}
	if len(roles) == 0 {
		return NoPermission, nil, nil
	}

	// 查询用户权限等级
	level, err := r.queryLevel(userid, teamid)
	if err != nil {
		zlog.Errorf("查询用户 %d 在团队 %d 的权限等级失败：%v", userid, teamid, err)
		return NoPermission, nil, err
	}

	// 查询用户拥有的 URL
	urls, err := r.queryURLs(roles, teamid)
	if err != nil {
		zlog.Errorf("查询用户 %d 在团队 %d 拥有的 URL 失败：%v", userid, teamid, err)
		return NoPermission, nil, err
	}

	return level, urls, nil
}

// 查询用户的角色
func (r CasbinRepo) queryRoles(userid, teamid int64) ([]int64, error) {
	var roles []int64
	err := r.DB.Model(&model.Casbin{}).
		Joins("JOIN user_power AS up1 ON up1.member_id = casbin.v0").
		Select("DISTINCT casbin.v2").
		Where("(casbin.ptype = ? AND casbin.v0 = ? AND casbin.v1 = ?) OR (casbin.ptype = ? AND casbin.v0 = ? AND casbin.v1 = ?)",
			"g", userid, teamid, // 指定团队角色
			"g", userid, 1). // 默认团队角色（team_id = 1）
		Pluck("casbin.v2", &roles).Error
	return roles, err
}

// 查询用户权限等级
func (r CasbinRepo) queryLevel(userid, teamid int64) (int, error) {
	var level int
	err := r.DB.Model(&model.User_Power{}).
		Select("level").
		Where("member_id = ? AND team_id = ?", userid, teamid).
		Scan(&level).Error
	return level, err
}

// 查询用户拥有的 URL
func (r CasbinRepo) queryURLs(roles []int64, teamid int64) ([]string, error) {
	var urls []string
	err := r.DB.Model(&model.Casbin{}).
		Select("DISTINCT casbin.v2").
		Where("casbin.ptype = ? AND casbin.v0 IN ? AND casbin.v1 IN ?",
			"p", roles, []int64{teamid, 1}). // 同时查询 teamid 和默认团队（team_id = 1）
		Pluck("casbin.v2", &urls).Error
	if err != nil {
		return nil, err
	}
	return urls, nil
}

// CheckUserPermission
//
//	@Description: 检查用户权限
//	@receiver r
//	@param url
//	@param userId
//	@param teamId
//	@return bool
//	@return error
const SUPERMANGER_ID = 33333

func (r CasbinRepo) CheckUserPermission(url string, userId, teamId int64) (bool, error) {
	defer util.RecordTime(time.Now())()

	//新建团队时不传团队Id，则赋给默认值
	if teamId == 0 {
		teamId = ROOTTEAM
	}

	var roles []string
	err := r.DB.Model(&model.Casbin{}).
		Select(RoleOrUrl). // 获取 g 规则中的 roleid
		Where("(casbin.ptype = ? AND casbin.v0 = ? AND casbin.v1 = ?) OR (casbin.ptype = ? AND casbin.v0 = ? AND casbin.v1 = ?)",
			"g", userId, teamId, // 指定团队角色
			"g", userId, 1). // 默认团队角色（team_id = 1）
		Find(&roles).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		zlog.Errorf("查询用户对应的用户组失败：%v", err)
		return false, err
	}

	var managers []int64
	for _, role := range roles {
		//将string类型转化为 int64
		manager, err := strconv.ParseInt(role, 10, 64)
		if err != nil {
			//转换失败
			return false, err
		}
		if manager == SUPERMANGER_ID {
			//说明 这个用户 是超级管理员，直接验证成功
			return true, nil
		}
		managers = append(managers, manager)
	}

	var res string
	// 使用 roleid 和 teamid 查询拥有的 URL
	err = r.DB.Model(&model.Casbin{}).
		Select(RoleOrUrl). // 获取 p 规则中的 url
		Where(fmt.Sprintf("%s = 'p' AND %s in (?) AND %s = ? AND %s = ?", C_Type, C_User, C_Team, RoleOrUrl), managers, teamId, url).
		First(&res).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 查询不到记录，无权限
			return false, nil
		}
		zlog.Errorf("权限验证出错：%v", err)
		return false, err
	}

	zlog.Infof("权限验证成功: managers=%v, teamId=%d, url=%s", managers, teamId, url)

	// 查询成功，存在记录，有权限
	return true, nil
}
