package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"strings"
	"tgwp/global"
	"tgwp/internal/logic"
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

		//正式使用，    测试时需注释掉
		UserId := logic.GetUserId(c)

		// 绑定 team_id，根据请求方法选择解析方式
		var req any
		var err error

		// 在请求处理中读取并缓存请求体
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			zlog.CtxErrorf(ctx, "读取请求体失败: %v", err)
		} else {
			// 保存请求体内容供后续使用
			cachedRequestBody := string(bodyBytes)
			zlog.CtxInfof(ctx, "请求体内容: %s", cachedRequestBody)

			// 将读取的请求体内容重新设置到 c.Request.Body，供后续处理使用
			c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		// 使用 cachedRequestBody 保存的内容来做日志记录或其它操作

		switch c.Request.Method {
		case "GET":
			req, err = types.BindReq[types.ParamsRuleCheck](c)
		case "POST", "PUT":
			req, err = types.BindReq[types.JsonRuleCheck](c)
		case "DELETE":
			req, err = types.BindReq[types.UriRuleCheck](c)
		}

		if err != nil {
			zlog.CtxErrorf(ctx, "PermissionMiddleware 参数绑定失败: %v", err)
			zlog.CtxInfof(ctx, "请求详情: Method=%s, Headers=%v, Query=%v", c.Request.Method, c.Request.Header, c.Request.URL.Query())
			response.NewResponse(c).Error(response.PARAM_NOT_VALID)
			c.Abort()
			return
		}
		zlog.CtxInfof(ctx, "PermissionMiddleware middleware: %v", req)

		// 重要!! 将读取的请求体内容重新设置到 c.Request.Body，供后续处理使用
		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		// 类型断言为包含 UserId 和 TeamId 的结构体
		//var userId int64 //正式使用时需删除
		var teamId int64

		switch v := req.(type) {
		case types.ParamsRuleCheck:
			//userId = v.UserId
			teamId = v.TeamId
		case types.JsonRuleCheck:
			//userId = v.UserId
			teamId = v.TeamId
		case types.UriRuleCheck:
			//userId = v.UserId
			teamId = v.TeamId
		default:
			zlog.CtxErrorf(ctx, "无效的绑定类型: %T", req)
			response.NewResponse(c).Error(response.PARAM_TYPE_ERROR)
			c.Abort()
			return
		}

		var Url string
		//如果是delete请求，需要截取 url的前面部分
		if c.Request.Method == "DELETE" {
			url := c.Request.URL.Path
			index := strings.Index(url, "/delete")
			if index != -1 {
				Url = url[:index+len("/delete")]
				zlog.CtxInfof(ctx, "Base URL:%v", Url)
			}
		} else {
			Url = c.Request.URL.Path
			zlog.CtxInfof(ctx, "Base URL:%v", Url)
		}

		// CheckUserPermissions 检查用户权限
		exist, err := repo.NewCasbinRepo(global.DB).CheckUserPermission(Url, UserId, teamId) //正式使用
		//exist, err := repo.NewCasbinRepo(global.DB).CheckUserPermission(Url, userId, teamId) //测试时使用

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
