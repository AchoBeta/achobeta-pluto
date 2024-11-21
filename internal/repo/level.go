package repo

import (
	"gorm.io/gorm"
	"tgwp/internal/model"
)

// 获取权限级别
type LevelRepo struct {
	DB *gorm.DB
}

func NewLevelRepo(db *gorm.DB) *LevelRepo {
	return &LevelRepo{DB: db}
}

// GetLevel
//
//	@Description:
//	@receiver r
//	@param userid
//	@param teamid
//	@return int
//	@return error
func (r LevelRepo) GetLevel(userid, teamid int64) (int, error) {
	var level int
	err := r.DB.Model(&model.User_Power{}).
		Select(C_Level).
		Where(&model.User_Power{
			MemberId: userid,
			TeamId:   teamid,
		}).
		First(&level).Error
	if err != nil {
		return 0, err
	}

	return level, nil
}
