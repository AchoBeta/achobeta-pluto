package model

import "time"

type Sign struct {
	CommonModel
	LoginId    string    `gorm:"column:login_id;type:char(19);unique;not null;comment:'登录id'"`
	Issuer     string    `gorm:"column:issuer;type:char(19);unique;not null;comment:'签发标识'"`
	UserId     int64     `gorm:"column:user_id;type:bigint;not null;comment:'用户id'"`
	OnlineTime time.Time `gorm:"column:online_time;not null;comment:'上线时间'"`
	UserAgent  string    `gorm:"column:user_agent;type:varchar(50);not null;comment:'用户代理'"`
	IP         string    `gorm:"column:ip;type:VARCHAR(45);not null;comment:'ip地址'"`
	DeviceName string    `gorm:"column:device;type:VARCHAR(50);comment:'设备名称''"`
	Phone      string    `gorm:"column:phone;type:char(11);not null;comment:'手机号'"`
}

func (t *Sign) TableName() string {
	return "sign"
}
