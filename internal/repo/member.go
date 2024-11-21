package repo

import (
	"errors"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"tgwp/global"
	"tgwp/internal/model"
	"tgwp/internal/types"
	"tgwp/log/zlog"
	"tgwp/util"
	"time"
)

type MemberDetailRepo struct {
	DB *gorm.DB
}

func NewMemberDetailRepo(db *gorm.DB) *MemberDetailRepo {
	return &MemberDetailRepo{
		DB: db,
	}
}

// GetMemberDetail
//
//	@Description: 查询 团队成员详细信息
//	@receiver r
//	@param userID
//	@return *types.GetMemberDetailResp
//	@return error
func (r *MemberDetailRepo) GetMemberDetail(userID int64) (*types.GetMemberDetailResp, error) {
	var resp types.GetMemberDetailResp

	var positions []struct {
		TeamID   int64  `json:"team_id"`
		TeamName string `json:"team_name"`
		NodeID   int64  `json:"node_id"`
		NodeName string `json:"node_name"`
		Level    int    `json:"level"`
	}

	// 查询用户的基本信息
	err := r.DB.Model(&model.Member{}).
		Select(`member.name, member.sex, member.create_date, member.id_card, 
				member.phone_num, member.email, member.grade, member.major, 
				member.student_id, member.experience, member.status, member.like_count`).
		Where("id = ?", userID).
		First(&resp).
		Error
	if err != nil {
		zlog.Errorf("查询用户详细信息失败: %v", err)
		return nil, err
	}

	// 查询用户的团队、职位和权限
	err = r.DB.Model(&model.Team_Member_Structure{}).
		Select(`team.id AS team_id, team.name AS team_name, 
				structure.id AS node_id, structure.node_name AS node_name, 
				user_power.level AS level`).
		Joins("JOIN team ON team.id = team_member_structure.team_id").
		Joins("JOIN structure ON structure.id = team_member_structure.structure_id").
		Joins("JOIN user_power ON user_power.member_id = team_member_structure.member_id AND user_power.team_id = team_member_structure.team_id").
		Where("team_member_structure.member_id = ?", userID).
		Find(&positions).
		Error
	if err != nil {
		zlog.Errorf("查询用户团队职位信息失败: %w", err)
		return nil, err
	}

	// 整理结果：按团队分组职位
	positionMap := make(map[int64]*types.MemberPositions)
	for _, pos := range positions {
		//将职位数据按团队分组
		if _, exists := positionMap[pos.TeamID]; !exists {
			positionMap[pos.TeamID] = &types.MemberPositions{
				TeamId:      pos.TeamID,
				TeamName:    pos.TeamName,
				Level:       pos.Level,
				StructNodes: []types.StructureNodes{},
			}
		}
		positionMap[pos.TeamID].StructNodes = append(positionMap[pos.TeamID].StructNodes, types.StructureNodes{
			ID:   pos.NodeID,
			Name: pos.NodeName,
		})
	}

	// 转换为切片
	for _, pos := range positionMap {
		resp.MemberPosition = append(resp.MemberPosition, *pos)
	}

	return &resp, nil
}

type MemberlistRepo struct {
	DB *gorm.DB
}

func NewMemberlistRepo(db *gorm.DB) *MemberlistRepo {
	return &MemberlistRepo{DB: db}
}

// MemberlistRepo
//
//	@Description: 分页查询  用户基本信息
//	@receiver r
//	@param TeamId
//	@param Page
//	@param Perpage
//	@return types.MemberlistResp
//	@return error
func (r *MemberlistRepo) MemberlistRepo(TeamId int64, Page, Perpage int) (types.MemberlistResp, error) {

	var Members []types.Members

	Offset := (Page - 1) * Perpage
	if Offset < 0 {
		Offset = 0
	}

	//先这样，要美观的话之后更改
	err := r.DB.Table("member").
		Select(`member.id AS member_id, member.name, member.grade, member.major, member.status, 
			member.phone_num, GROUP_CONCAT(structure.node_name) AS positions`).
		Joins(`JOIN team_member_structure ON team_member_structure.member_id = member.id`).
		Joins(`JOIN structure ON structure.id = team_member_structure.structure_id`).
		Where("team_member_structure.team_id = ?", TeamId).
		Group("member.id"). // 按 Member 分组
		Offset(Offset).
		Limit(Perpage).
		Find(&Members).
		Error

	//在查询后将 Positions 的逗号分隔值拆分为切片： 查询 当前用户所处团队的每个职位名称
	for i, member := range Members {
		Members[i].PositionList = strings.Split(member.Positions, ",")
	}

	if err != nil {
		zlog.Errorf("查询用户简单信息失败")
		return types.MemberlistResp{}, err
	}
	return types.MemberlistResp{Members: Members}, nil
}

const (
	init_likecount = 0
	team_id        = "id"
	init_teamname  = "未分配团队"
)

// 将 int64 转化成 string
var (
	NormalManger = strconv.FormatInt(global.NORMAL_ADMINISTRATOR, 10)
	SuperManger  = strconv.FormatInt(global.SUPERL_ADMINISTRATOR, 10)
)

type CreateMemberRepo struct {
	DB *gorm.DB
}

func NewCreateMemberRepo(db *gorm.DB) *CreateMemberRepo {
	return &CreateMemberRepo{
		DB: db,
	}
}

// CreateMember
//
//	@Description: 新增团队成员
//	@receiver r
//	@param req
//	@return error
func (r *CreateMemberRepo) CreateMember(req types.CreateMemberReq) error {
	defer util.RecordTime(time.Now())()
	err := r.DB.Model(&model.Member{}).
		Create(&model.Member{
			CommonModel: model.CommonModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name:       req.Name,
			Sex:        req.Sex,
			CreateDate: req.CreateDate,
			IdCard:     req.IdCard,
			PhoneNum:   req.PhoneNum,
			Email:      req.Email,
			Grade:      req.Grade,
			Major:      req.Major,
			StudentID:  req.StudentID,
			Experience: req.Experience,
			Status:     req.Status,
			LikeCount:  init_likecount,
		}).
		Error
	if err != nil {
		zlog.Errorf("新增团队成员失败：: %v", err)
		return err
	}

	//获取新成员id
	var UserID int64

	err = r.DB.Model(&model.Member{}).
		Select(C_Id).
		Where(&model.Member{
			PhoneNum: req.PhoneNum,
		}).
		First(&UserID).
		Error
	if err != nil {
		zlog.Errorf("获取新成员id失败: %v", err)
		return err
	}

	if req.MemberPositions == nil {
		//说明没有给他分配任何团队，那么把他放入“未分配团队”的团队中
		var InitTeamId int64

		err = r.DB.Model(&model.Team{}).
			Select(team_id).
			Where("name = ?", init_teamname).
			First(&InitTeamId).
			Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				//说明没有“未分配团队”的团队id不存在，重新创建
				_, err = NewCreateTeamRepo(global.DB).CreateTeam(init_teamname)
				if err != nil {
					zlog.Errorf("创建未分配团队失败: %v", err)
					return err
				}
			} else {
				zlog.Errorf("查询未分配团队失败：: %v", err)
				return err
			}
			//重新获取
			err = r.DB.Model(&model.Team{}).
				Select(team_id).
				Where("name = ?", init_teamname).
				First(&InitTeamId).
				Error
			if err != nil {
				zlog.Errorf("查询未分配团队失败：: %v", err)
				return err
			}
		}

		// 新成员id 和 未分配团队id 在 Team_Member_Structure 做一个简单关联
		err = r.DB.Model(&model.Team_Member_Structure{}).
			Create(&model.Team_Member_Structure{
				MemberId: UserID,
				TeamId:   InitTeamId,
			}).Error
		if err != nil {
			zlog.Errorf("新成员放入未分配团队失败：: %v", err)
			return err
		}
	}

	for _, member_position := range req.MemberPositions {
		for _, structure_node := range member_position.StructNodes {
			// 新成员id 和 团队id 对 Team_Member_Structure表 初始化
			// 多职位id关联
			err = r.DB.Model(&model.Team_Member_Structure{}).
				Create(&model.Team_Member_Structure{
					MemberId:    UserID,
					TeamId:      member_position.TeamId,
					StructureId: structure_node.ID,
				}).Error
			if err != nil {
				zlog.Errorf("对Team_Member_Structure表初始化失败: %v", err)
				return err
			}
		}
		// 对 User_Power表 初始化
		err = r.DB.Model(&model.User_Power{}).
			Create(&model.User_Power{
				MemberId: UserID,
				TeamId:   member_position.TeamId,
				Level:    member_position.Level,
			}).Error
		if err != nil {
			zlog.Errorf("对User_Power表初始化失败: %v", err)
			return err
		}

		// 根据权限级别 对 Casbin表 初始化
		if member_position.Level >= 2 {
			err = r.DB.Model(&model.Casbin{}).
				Create(&model.Casbin{
					Ptype: "g",
					V0:    UserID,
					V1:    member_position.TeamId,
					V2:    NormalManger,
				}).Error
			if err != nil {
				zlog.Errorf("对权限等级大于等于2的casbin表初始化失败: %v", err)
				return err
			}
		}
		if member_position.Level == 3 {
			err = r.DB.Model(&model.Casbin{}).
				Create(&model.Casbin{
					Ptype: "g",
					V0:    UserID,
					V1:    member_position.TeamId,
					V2:    SuperManger,
				}).Error
			if err != nil {
				zlog.Errorf("对权限等级等于3的casbin表初始化失败: %v", err)
				return err
			}
		}

	}

	return nil
}

type DeleteMemberRepo struct {
	DB *gorm.DB
}

func NewDeleteMemberRepo(db *gorm.DB) *DeleteMemberRepo {
	return &DeleteMemberRepo{
		DB: db,
	}
}

// DeleteMember
//
//	@Description: 删除团队成员
//	@receiver r
//	@param MemberId
//	@param TeamId
//	@return error
func (r *DeleteMemberRepo) DeleteMember(MemberId, TeamId int64) error {
	defer util.RecordTime(time.Now())()

	//删除 用户 和 团队 关联的职位
	err := r.DB.Model(&model.Team_Member_Structure{}).
		Delete(&model.Team_Member_Structure{}).
		Where(&model.Team_Member_Structure{
			MemberId: MemberId,
			TeamId:   TeamId,
		}).Error
	if err != nil {
		zlog.Errorf("删除当前成员和团队相关职位信息失败: %v", err)
		return err
	}

	//删除 用户 和 团队 关联的权限级别
	err = r.DB.Model(&model.User_Power{}).
		Delete(&model.User_Power{}).
		Where(&model.User_Power{
			MemberId: MemberId,
			TeamId:   TeamId,
		}).Error
	if err != nil {
		zlog.Errorf("删除当前成员和团队相关权限级别失败: %v", err)
		return err
	}

	//删除 用户 和 团队 关联的 casbin表
	err = r.DB.Model(&model.Casbin{}).
		Delete(&model.Casbin{}).
		Where(&model.Casbin{
			Ptype: "g",
			V0:    MemberId,
			V1:    TeamId,
		}).Error
	if err != nil {
		zlog.Errorf("删除当前成员和团队相关casbin信息失败: %v", err)
		return err
	}

	//删除 用户 和 团队 关联的 点赞表
	err = r.DB.
		Where("memberid_like = ? OR memberid_belike = ?", MemberId, MemberId).
		Delete(&model.Like_Status{}).
		Error
	if err != nil {
		zlog.Errorf("删除当前成员相关点赞信息信息失败: %v", err)
		return err
	}

	return nil
}

type PutMemberRepo struct {
	DB *gorm.DB
}

func NewPutMemberRepo(db *gorm.DB) *PutMemberRepo {
	return &PutMemberRepo{
		DB: db,
	}
}

// PutMember
//
//	@Description: 编辑团队成员
//	@receiver r
//	@param req
//	@return error
func (r *PutMemberRepo) PutMember(req types.PutTeamMemberReq) error {
	defer util.RecordTime(time.Now())()

	//更改基本信息
	err := r.DB.Model(&model.Member{}).Updates(&model.Member{
		Name:       req.Name,
		Sex:        req.Sex,
		CreateDate: req.CreateDate,
		IdCard:     req.IdCard,
		PhoneNum:   req.PhoneNum,
		Email:      req.Email,
		Grade:      req.Grade,
		Major:      req.Major,
		StudentID:  req.StudentID,
		Experience: req.Experience,
		Status:     req.Status,
	}).Where(&model.Member{
		CommonModel: model.CommonModel{
			ID: req.ID,
		},
		PhoneNum: req.PhoneNum,
	}).Error
	if err != nil {
		zlog.Errorf("修改成员信息失败: %v", err)
		return err
	}

	/// 删除旧职位信息
	err = r.DB.Where("member_id = ?", req.ID).
		Delete(&model.Team_Member_Structure{}).
		Error
	if err != nil {
		zlog.Errorf("删除成员旧职位失败: %v", err)
		return err
	}

	// 增加新职位信息
	for _, position := range req.MemberPosition {
		for _, node := range position.StructNodes {
			newStructure := model.Team_Member_Structure{
				TeamId:      position.TeamId,
				MemberId:    req.ID,
				StructureId: node.ID,
			}
			if err = r.DB.Create(&newStructure).Error; err != nil {
				zlog.Errorf("增加成员新职位失败: %v", err)
				return err
			}
		}

		//更新 user_power权限表
		err = r.DB.Model(&model.User_Power{}).
			Update("level", position.Level).
			Where(&model.User_Power{
				MemberId: req.ID,
				TeamId:   position.TeamId,
			}).Error
		if err != nil {
			zlog.Errorf("更新成员在当前团队的权限级别失败: %v", err)
			return err
		}

		//删除 用户 和 团队 关联的 旧casbin表
		err = r.DB.Model(&model.Casbin{}).
			Delete(&model.Casbin{}).
			Where(&model.Casbin{
				Ptype: "g",
				V0:    req.ID,
				V1:    position.TeamId,
			}).Error
		if err != nil {
			zlog.Errorf("删除成员在当前团队旧casbin信息失败: %v", err)
			return err
		}

		//新增 casbin 数据
		if position.Level >= 2 {
			err = r.DB.Model(&model.Casbin{}).
				Create(&model.Casbin{
					Ptype: "g",
					V0:    req.ID,
					V1:    position.TeamId,
					V2:    NormalManger,
				}).Error
			if err != nil {
				zlog.Errorf("增加成员在当前团队权限等级大于等于2的新casbin信息失败: %v", err)
				return err
			}
		}
		if position.Level == 3 {
			err = r.DB.Model(&model.Casbin{}).
				Create(&model.Casbin{
					Ptype: "g",
					V0:    req.ID,
					V1:    position.TeamId,
					V2:    SuperManger,
				}).Error
			if err != nil {
				zlog.Errorf("增加成员在当前团队权限等级等于3的新casbin信息失败: %v", err)
				return err
			}
		}

	}
	return nil
}

// 留给淳桂                    //
type JudgeUserRepo struct {
	DB *gorm.DB
}

func NewJudgeUserRepo(db *gorm.DB) *JudgeUserRepo {
	return &JudgeUserRepo{DB: db}
}

// JudgeUser
//
//	@Description: 通过手机号判断是游客还是团队成员
//	@receiver r
//	@param Phone
//	@return int64
//	@return bool
//	@return error
func (r JudgeUserRepo) JudgeUser(Phone uint64) (int64, bool, error) {
	defer util.RecordTime(time.Now())()

	var UserID int64

	err := r.DB.Model(&model.Member{}).
		Select(C_Id).
		Where(&model.Member{
			PhoneNum: Phone,
		}).
		First(&UserID).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, false, nil
		}
		return 0, false, err
	}
	return UserID, true, nil
}
