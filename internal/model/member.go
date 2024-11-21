package model

import "time"

// 个人详细信息表
type Member struct {
	CommonModel
	Name       string    `gorm:"column:name; type:varchar(20); index; comment:'真实姓名'"`
	Sex        string    `gorm:"column:sex; type:char(2); comment:'性别'"`
	CreateDate time.Time `gorm:"column:create_date; type:date; comment:'加入时间'"`
	IdCard     string    `gorm:"column:idcard; type:varchar(50); comment:'身份证'"`
	PhoneNum   uint64    `gorm:"column:phone_num; type:int unsigned; comment:'手机号码'"`
	Email      string    `gorm:"column:email; type:varchar(30); comment:'邮箱'"`
	Grade      uint      `gorm:"column:grade; type:int unsigned; comment:'年级'"`
	Major      string    `gorm:"column:major; type:varchar(30); comment:'专业'"`
	StudentID  uint64    `gorm:"column:student_id; type:int unsigned; comment:'学号'"`
	Experience string    `gorm:"column:experience; type:varchar(200); comment:'实习、创业、就职经历'"`
	Status     string    `gorm:"column:status; type:varchar(10); comment:'现状'"`
	LikeCount  uint64    `gorm:"column:like_count; type:int unsigned; comment:'点赞数量'"`
}

func (t *Member) TableName() string {
	return "member"
}
