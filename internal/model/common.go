package model

import (
	"gorm.io/gorm"
	"tgwp/global"
	"time"
)

// CommonModel 每张表都有的四个东西，最好不要用 gorm.model（虽然他们一模一样）
type CommonModel struct {
	ID        int64 `gorm:"primaryKey;column:id;type:bigint;<-:create"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *CommonModel) BeforeCreate() error {
	// 生成雪花ID
	b.ID = global.Node.Generate().Int64()

	return nil
}
