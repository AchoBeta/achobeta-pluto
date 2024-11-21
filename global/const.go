package global

import (
	"tgwp/util/snowflake"
	"time"
)

// 所有常量文件读取位置
const (
	DEFAULT_CONFIG_FILE_PATH = "/config.yaml"
	REDIS_PHONE_CODE         = "Achobeta:phone.login.code:%s:string"
	REDIS_PHONE              = "Achobeta:phone.login.effective:%s:string"
	ATOKEN_EFFECTIVE_TIME    = time.Second * 300
	RTOKEN_EFFECTIVE_TIME    = time.Hour * 1
	AUTH_ENUMS_ATOKEN        = "atoken"
	AUTH_ENUMS_RTOKEN        = "rtoken"
	DEFAULT_NODE_ID          = 1
	TOKEN_USER_ID            = "UserId"
	ROOT_ID                  = 1     // 根节点 ID
	NORMAL_ADMINISTRATOR     = 22222 //普通管理员
	SUPERL_ADMINISTRATOR     = 33333 //超级管理员
)

var Node, _ = snowflake.NewNode(DEFAULT_NODE_ID)

var NORMAL_ADMIN_URLS = []string{
	"/api/team/memberlist/delete",
	"/api/team/memberlist/put",
	"/api/team/membermsg/save",
	"/api/team/structure/collection",
}

// global 包中的定义
var SUPER_ADMIN_URLS = []string{
	"/api/team/structure/change",
	"/api/team/structure/add",
}
