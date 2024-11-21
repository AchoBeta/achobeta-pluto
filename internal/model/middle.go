package model

// 团队-架构-用户中间表
type Team_Member_Structure struct {
	CommonModel
	TeamId      int64 `gorm:"column:team_id; type:bigint; index; not null; comment:'团队ID'"`
	MemberId    int64 `gorm:"column:member_id; type:bigint; index; not null; comment:'成员ID'"`
	StructureId int64 `gorm:"column:structure_id; type:bigint; index; not null; comment:'职位ID'"`

	// 外键关联
	Team      Team      `gorm:"foreignKey:TeamId; references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Member    Member    `gorm:"foreignKey:MemberId; references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Structure Structure `gorm:"foreignKey:StructureId; references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (t *Team_Member_Structure) TableName() string {
	return "team_member_structure"
}

// 用户权限表
type User_Power struct {
	CommonModel
	MemberId int64 `gorm:"column:member_id; type:bigint;   index;not null; comment:'成员ID'"`
	TeamId   int64 `gorm:"column:team_id;   type:bigint;   index;not null; comment:'团队ID'"`
	Level    int   `gorm:"column:level;     type:int; index;not null; default:1; comment:'权限等级'"`

	// 外键关联
	Member Member `gorm:"foreignKey:MemberId; references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Team   Team   `gorm:"foreignKey:TeamId; references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (t *User_Power) TableName() string {
	return "user_power"
}

// 点赞表
type Like_Status struct {
	CommonModel
	IsLiked         bool  `gorm:"column:is_liked; type:boolean; index; comment:'点赞情况'"`
	MemberId_Like   int64 `gorm:"column:memberid_like; type:int unsigned ; index; comment:'点赞的用户id'"`
	MemberId_BeLike int64 `gorm:"column:memberid_belike; type:int unsigned; index; comment:'被点赞的用户id'"`

	Member_Like   Member `gorm:"foreignKey:MemberId_Like;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Member_BeLike Member `gorm:"foreignKey:MemberId_BeLike;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (t *Like_Status) TableName() string {
	return "like_status"
}
