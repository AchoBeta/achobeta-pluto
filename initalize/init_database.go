package initalize

import (
	"tgwp/configs"
	"tgwp/global"
	"tgwp/internal/model"
	"tgwp/internal/pkg/database"
	"tgwp/internal/pkg/mysqlx"
	"tgwp/internal/pkg/redisx"
	"tgwp/log/zlog"
)

func InitDataBase(config configs.Config) {
	switch config.DB.Driver {
	case "mysql":
		database.InitDataBases(mysqlx.NewMySql(), config)
	default:

		zlog.Fatalf("不支持的数据库驱动：%s", config.DB.Driver)
	}
	if config.App.Env != "pro" {
		err := global.DB.AutoMigrate()

		//迁移数据库所有的表
		migrateTables()

		//自动迁移 sign 表，确保表结构存在
		global.DB.AutoMigrate(&model.Sign{})

		//自动迁移 sign 表，确保表结构存在
		global.DB.AutoMigrate(&model.Sign{})

		if err != nil {
			zlog.Fatalf("数据库迁移失败！")
		}
	}
	zlog.Infof("数据库初始化成功！")
}
func InitRedis(config configs.Config) {
	if config.Redis.Enable {
		var err error
		global.Rdb, err = redisx.GetRedisClient(config)
		if err != nil {
			zlog.Errorf("无法初始化Redis : %v", err)
		}
	} else {
		zlog.Warnf("不使用Redis")
	}

}

func migrateTables() {
	// 自动迁移 casbin 表，确保表结构存在
	global.DB.AutoMigrate(&model.Casbin{})

	// 自动迁移 team 表，确保表结构存在
	global.DB.AutoMigrate(&model.Team{})

	// 自动迁移 member 表，确保表结构存在
	global.DB.AutoMigrate(&model.Member{})

	// 自动迁移 like_status 表，确保表结构存在
	global.DB.AutoMigrate(&model.Like_Status{})

	// 自动迁移 structure 表，确保表结构存在
	global.DB.AutoMigrate(&model.Structure{})

	// 自动迁移 team_member_structure 表，确保表结构存在
	global.DB.AutoMigrate(&model.Team_Member_Structure{})

	// 自动迁移 user_power 表，确保表结构存在
	global.DB.AutoMigrate(&model.User_Power{})
}
