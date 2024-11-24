package global

import (
	"tgwp/util/snowflake"
	"time"
)

// 所有常量文件读取位置
const (
	DEFAULT_CONFIG_FILE_PATH      = "/config.yaml"
	REDIS_PHONE_CODE              = "Achobeta:phone.login.code:%s:string"
	REDIS_PHONE                   = "Achobeta:phone.login.effective:%s:string"
	ATOKEN_EFFECTIVE_TIME         = time.Hour * 12
	RTOKEN_EFFECTIVE_TIME         = time.Hour * 24 * 30
	AUTH_ENUMS_ATOKEN             = "atoken"
	AUTH_ENUMS_RTOKEN             = "rtoken"
	DEFAULT_NODE_ID               = 1
	TOKEN_USER_ID                 = "UserId"
	ROOT_ID                       = 1                                  // 根节点 ID
	NORMAL_ADMINISTRATOR          = 22222                              //普通管理员
	SUPERL_ADMINISTRATOR          = 33333                              //超级管理员
	FEISHU_APP_ID                 = "cli_a7a37dc364b2900c"             // 飞书自建应用ID
	FEISHU_APP_SECRET             = "cqFNuZmHaIKlMeFAr8546cHlaquXw8ep" // 飞书自建应用Secret
	FEISHU_APP_TOKEN              = "M5l2bHYEiaYq2esmVM1cTyamn5s"      // 飞书自建应用Token
	FEISHU_TASK_TABLE_ID          = "tblM1AuOpuhpxBSb"                 //飞书多维表格任务表ID
	FEISHU_LIST_UPDATE_TIME       = 60 * 5                             // 飞书任务表更新频率(秒 int格式)
	FEISHU_LIST_WILL_OVERDUE_TIME = 60 * 60 * 24                       // 飞书任务表即将过期时间(秒 int格式)
)

var Node, _ = snowflake.NewNode(DEFAULT_NODE_ID)

var NORMAL_ADMIN_URLS = []string{
	"/api/team/memberlist/delete",
	"/api/team/memberlist/create",
	"/api/team/membermsg/change",
	"/api/team/structure/collection",
}

// global 包中的定义
var SUPER_ADMIN_URLS = []string{
	"/api/team/structure/change",
	"/api/team/structure/create",
}
