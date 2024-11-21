package repo

import (
	"gorm.io/gorm"
	"tgwp/global"
	"tgwp/internal/model"
	"tgwp/internal/types"
	"tgwp/log/zlog"
	"tgwp/util"
	"time"
)

type CreateTeamRepo struct {
	DB *gorm.DB
}

func NewCreateTeamRepo(db *gorm.DB) *CreateTeamRepo {
	return &CreateTeamRepo{
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
func (r CreateTeamRepo) CreateTeam(TeamName string) (*types.CreateTeamResp, error) {

	//创建新团队
	err := r.DB.Model(&model.Team{}).
		Create(&model.Team{Name: TeamName}).
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
			CommonModel: model.CommonModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
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
			CommonModel: model.CommonModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
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
			CommonModel: model.CommonModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Ptype: "p",
			V0:    global.SUPERL_ADMINISTRATOR, // 超级管理员
			V1:    TeamId,
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

type TeamIdRepo struct {
	DB *gorm.DB
}

const c_teamid = "team_id"
const c_teamname = "name"

func NewTeamIdRepo(db *gorm.DB) *TeamIdRepo {
	return &TeamIdRepo{DB: db}
}

// GetTeamId
//
//	@Description:
//	@receiver r
//	@param userid
//	@return first_team
//	@return team
//	@return err
func (r TeamIdRepo) GetTeamId(userid int64) (first_team types.Team, team []types.Team, err error) {
	defer util.RecordTime(time.Now())()
	err = r.DB.Model(&model.Team_Member_Structure{}).
		Select(c_teamid, c_teamname).
		Joins("JOIN team on team.id = team_member_structure.team_id").
		Where(&model.Team_Member_Structure{
			MemberId: userid,
		}).
		First(&first_team).
		Error
	if err != nil {
		zlog.Errorf("用户所在的第一个团队信息获取失败：%v", err)
		return
	}

	err = r.DB.Model(&model.Team{}).
		Joins("JOIN team on team.id = team_member_structure.team_id").
		Select(C_Id, c_teamname).
		Find(&team).
		Error
	if err != nil {
		zlog.Errorf("团队信息获取失败：%v", err)
		return
	}
	return
}
