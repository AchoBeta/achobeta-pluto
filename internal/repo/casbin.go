package repo

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"tgwp/internal/model"
	"tgwp/util"
	"time"
)

type CasbinRepo struct {
	DB *gorm.DB
}

const Nothing = 0

func NewCasbinRepo(db *gorm.DB) *CasbinRepo {
	return &CasbinRepo{DB: db}
}

// GetCasbin
//
//	@Description:  获取权限组
//	@receiver r
//	@param userid
//	@param teamid
//	@return []string
//	@return error
func (r CasbinRepo) GetCasbin(userid, teamid int64) (int, []string, error) {
	// 根据 UserId 查询用户对应的角色
	var Power []struct {
		role  int64
		level int
	}

	err := r.DB.Model(&model.Casbin{}).
		Joins("JOIN user_power ON user_power.member_id = casbin.v0").
		Joins("JOIN user_power ON user_power.team_id = casbin.v1").
		Select(RoleOrUrl, C_Level). // 获取 g 规则中的 roleid
		Where(&model.Casbin{
			Ptype: "g",
			V0:    userid,
			V1:    teamid,
		}).
		Find(&Power).Error
	if err != nil {
		return Nothing, nil, err
	}
	Level := Power[0].level

	var roles []int64
	for _, power := range Power {
		roles = append(roles, power.role)
	}

	// 使用 roleid 和 teamid 查询拥有的 URL
	var urls []string
	err = r.DB.Model(&model.Casbin{}).
		Select(RoleOrUrl). // 获取 p 规则中的 url
		Where(fmt.Sprintf("%s = 'p' AND %s in (?) AND %s = ?", C_Type, C_User, C_Team), roles, teamid).
		Find(&urls).Error
	if err != nil {
		return Nothing, nil, err
	}

	return Level, urls, nil
}

// 查询用户权限
type PermissionRepo struct {
	DB *gorm.DB
}

func NewPermissionRepo(db *gorm.DB) *PermissionRepo {
	return &PermissionRepo{DB: db}
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
func (r PermissionRepo) CheckUserPermission(url string, userId, teamId int64) (bool, error) {
	defer util.RecordTime(time.Now())()

	var roles []string
	err := r.DB.Model(&model.Casbin{}).
		Select(RoleOrUrl). // 获取 g 规则中的 roleid
		Where(&model.Casbin{
			Ptype: "g",
			V0:    userId,
			V1:    teamId,
		}).
		Find(&roles).Error
	if err != nil {
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
		managers = append(managers, manager)
	}

	var res string
	// 使用 roleid 和 teamid 查询拥有的 URL
	err = r.DB.Model(&model.Casbin{}).
		Select(RoleOrUrl). // 获取 p 规则中的 url
		Where(fmt.Sprintf("%s = 'p' AND %s in (?) AND %s = ? AND %s = ?", C_Type, C_User, C_Team), managers, teamId, url).
		First(&res).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 查询不到记录，无权限
			return false, nil
		}
		// 查询出错，返回错误
		return false, err
	}

	// 查询成功，存在记录，有权限
	return true, nil
}
