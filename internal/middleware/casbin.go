package middleware

import (
	"github.com/gin-gonic/gin"
	"tgwp/global"
	"tgwp/internal/repo"
	"tgwp/internal/response"
	"tgwp/internal/types"
	"tgwp/log/zlog"
)

// PermissionMiddleware
//
//	@Description:
//	@return gin.HandlerFunc
//
// 权限校验中间件：检查用户是否有权限访问某个资源
func PermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := zlog.GetCtxFromGin(c)
		if userid, exists := c.Get(global.TOKEN_USER_ID); exists {
			c.Set("user_id", userid)
		}
		//绑定 user_id 和 team_id
		req, err := types.BindReq[types.RuleCheck](c)

		if err != nil {
			zlog.CtxErrorf(ctx, "PermissionMiddleware err:%v", err)
			response.NewResponse(c).Error(response.PARAM_NOT_VALID)
			c.Abort()
			return
		}
		zlog.CtxInfof(ctx, "PermissionMiddleware middleware: %v", req)

		Url := c.Request.URL.Path

		// CheckUserPermissions 检查用户权限
		exist, err := repo.NewPermissionRepo(global.DB).CheckUserPermission(Url, req.UserId, req.TeamId)

		if err != nil {
			response.NewResponse(c).Error(response.PARAM_NOT_VALID)
			c.Abort()
			return
		}
		if exist == false {
			response.NewResponse(c).Error(response.INSUFFICENT_PERMISSIONS)
			c.Abort()
			return
		}
		c.Next() // 继续处理请求
	}
}
