package model

type Casbin struct {
	CommonModel
	Ptype string `gorm:"column:ptype;type:char(2);index;not null;comment:'权限类型'"`
	V0    int64  `gorm:"column:v0;type:bigint;index;not null;comment:'用户ID'"`
	V1    int64  `gorm:"column:v1;type:bigint;index;not null;comment:'团队ID'"`
	V2    string `gorm:"column:v2;type:varchar(100);index;comment:'用户的请求URL'"`

	//外键关联
	Member Member `gorm:"foreignKey:V0;references:ID"`
	Team   Team   `gorm:"foreignKey:V1;references:ID"`
}

func (t *Casbin) TableName() string {
	return "casbin"
}
