package model

// 团队架构表
type Structure struct {
	CommonModel
	TeamId   int64  `gorm:"column:team_id;     type:bigint; index; not null; comment:'团队ID'"`
	FatherId int64  `gorm:"column:father_id;   type:bigint; index; comment:'父节点ID'"`
	NodeName string `gorm:"column:node_name;   type:varchar(50); index; not null; comment:'分组/职位名字'"`

	// 外键关联
	Parent *Structure `gorm:"foreignKey:FatherId;references:ID"`
}

// 出参放在 types的 TeamStructResp 中

func (t *Structure) TableName() string {
	return "structure"
}
