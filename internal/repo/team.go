package repo

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"tgwp/global"
	"tgwp/internal/model"
	"tgwp/internal/types"
	"tgwp/log/zlog"
	"tgwp/util"
	"time"
)

type TeamRepo struct {
	DB *gorm.DB
}

func NewTeamRepo(db *gorm.DB) *TeamRepo {
	return &TeamRepo{
		DB: db,
	}
}

// CreateTeam
//
//	@Description: 新增团队 ，初始化团队架构
//	@receiver r
//	@param TeamName
//	@return types.CreateTeamResp
//	@return error
const ROOTTEAM = 1

func (r TeamRepo) CreateTeam(TeamName string) (*types.CreateTeamResp, error) {

	// 没有根节点才创建，有就不用
	// 禁用外键约束
	err := r.DB.Exec("SET FOREIGN_KEY_CHECKS = 0").Error
	if err != nil {
		zlog.Errorf("临时禁用外键约束失败: %v", err)
	}
	//创建团队根节点
	err = r.DB.Model(&model.Team{}).
		FirstOrCreate(&model.Team{
			CommonModel: model.CommonModel{
				ID: ROOTTEAM,
			},
			Name: "超级管理员所在的团队",
		}).
		Error
	if err != nil {
		zlog.Errorf("创建团队根节点 失败：%v", err)
		return &types.CreateTeamResp{}, err
	}

	//创建团队架构新节点
	err = r.DB.Model(&model.Structure{}).
		FirstOrCreate(&model.Structure{
			CommonModel: model.CommonModel{
				ID: global.ROOT_ID,
			},
			FatherId: global.ROOT_ID,
			NodeName: "所有团队的根节点",
			TeamId:   ROOTTEAM,
		}).
		Error
	if err != nil {
		zlog.Errorf("创建团队架构新节点：%v", err)
		return &types.CreateTeamResp{}, err
	}

	NormalManger := &model.Member{
		CommonModel: model.CommonModel{
			ID: 22222,
		},
		Name:       "普通管理员",
		CreateDate: time.Now(),
		PhoneNum:   "22222",
	}
	//创建普通管理员
	err = r.DB.Model(&model.Member{}).
		FirstOrCreate(NormalManger).
		Error
	if err != nil {
		zlog.Errorf("创建普通管理员失败：%v", err)
		return &types.CreateTeamResp{}, err
	}

	SuperManger := &model.Member{
		CommonModel: model.CommonModel{
			ID: 33333,
		},
		Name:       "超级管理员",
		CreateDate: time.Now(),
		PhoneNum:   "33333",
	}
	//创建超级管理员
	err = r.DB.Model(&model.Member{}).
		FirstOrCreate(SuperManger).
		Error
	if err != nil {
		zlog.Errorf("创建超级管理员失败：%v", err)
		return &types.CreateTeamResp{}, err
	}

	// 启用外键约束
	err = r.DB.Exec("SET FOREIGN_KEY_CHECKS = 1").Error
	if err != nil {
		log.Fatalf("临时启动外键约束失败: %v", err)
	}

	//创建新团队
	err = r.DB.Model(&model.Team{}).
		Create(&model.Team{
			Name: TeamName,
		}).
		Error
	if err != nil {
		zlog.Errorf("生成新团队id 失败：%v", err)
		return &types.CreateTeamResp{}, err
	}

	//找到新创建的团队ID
	var TeamId int64
	err = r.DB.Model(&model.Team{}).
		Select(C_Id).
		Where(&model.Team{Name: TeamName}).
		First(&TeamId).
		Error
	if err != nil {
		zlog.Errorf("未查询到新团队：%v", err)
		return &types.CreateTeamResp{}, err
	}

	//初始化团队架构
	err = r.DB.Model(&model.Structure{}).
		Create(&model.Structure{
			FatherId: global.ROOT_ID,
			NodeName: TeamName,
			TeamId:   TeamId,
		}).
		Error
	if err != nil {
		zlog.Errorf("新团队架构初始化失败：%v", err)
		return &types.CreateTeamResp{}, err
	}

	//初始化团队权限组
	Rules := []*model.Casbin{}

	// 普通管理员
	for _, url := range global.NORMAL_ADMIN_URLS {
		rule := &model.Casbin{
			Ptype: "p",
			V0:    global.NORMAL_ADMINISTRATOR,
			V1:    TeamId,
			V2:    url,
		}
		Rules = append(Rules, rule)
	}
	// 超级管理员
	for _, url := range global.SUPER_ADMIN_URLS {
		rule := &model.Casbin{
			Ptype: "p",
			V0:    global.SUPERL_ADMINISTRATOR, // 超级管理员
			V1:    1,
			V2:    url,
		}
		Rules = append(Rules, rule)
	}

	for _, Rule := range Rules {
		err = r.DB.Model(&model.Casbin{}).
			Create(Rule).
			Error
	}

	if err != nil {
		zlog.Errorf("新团队权限组初始化失败：%v", err)
		return &types.CreateTeamResp{}, err
	}
	return &types.CreateTeamResp{}, nil
}

// GetTeamId
//
//	@Description:  获得团队id,name
//	@receiver r
//	@param userid
//	@return first_team
//	@return team
//	@return err
const c_id = "id"
const c_teamid = "team_id"
const c_teamname = "name"

func (r TeamRepo) GetTeamId(userid int64) (types.Team, []types.Team, error) {
	defer util.RecordTime(time.Now())()

	var first_team types.Team
	var team []types.Team

	//查询第一个团队
	err := r.DB.Model(&model.Team{}).
		Select("team.id, team.name").
		Joins("JOIN team_member_structure AS tms ON tms.team_id = team.id").
		Where("tms.member_id = ?", userid).
		First(&first_team).
		Error
	if err != nil {
		zlog.Warnf("未找到用户的第一个团队：%v", userid)
		zlog.Errorf("用户所在的第一个团队信息获取失败：%v", err)
		return types.Team{}, nil, err
	}

	//查询所有团队
	err = r.DB.Model(&model.Team{}).
		Select(fmt.Sprintf("%s, %s", c_id, c_teamname)).
		Where("team.deleted_at IS NULL").
		Find(&team).
		Error
	if err != nil {
		zlog.Errorf("团队信息获取失败：%v", err)
		return types.Team{}, nil, err
	}
	return first_team, team, nil
}
