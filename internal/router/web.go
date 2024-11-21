package routerg

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"tgwp/configs"
	"tgwp/global"
	"tgwp/internal/api"
	"tgwp/internal/manager"
	"tgwp/internal/middleware"
	"tgwp/internal/response"
	"tgwp/log/zlog"
)

// RunServer 启动服务器 路由层
func RunServer() {
	r, err := listen()
	if err != nil {
		zlog.Errorf("Listen error: %v", err)
		panic(err.Error())
	}
	r.Run(fmt.Sprintf("%s:%d", configs.Conf.App.Host, configs.Conf.App.Port)) // 启动 Gin 服务器
}

// listen 配置 Gin 服务器
func listen() (*gin.Engine, error) {
	r := gin.Default() // 创建默认的 Gin 引擎

	// 注册全局中间件（例如获取 Trace ID）
	manager.RequestGlobalMiddleware(r)

	// 创建 RouteManager 实例
	routeManager := manager.NewRouteManager(r)

	// 注册各业务路由组的具体路由
	registerRoutes(routeManager)
	messageRoutes(routeManager)

	return r, nil
}

// registerRoutes 注册各业务路由的具体处理函数
func registerRoutes(routeManager *manager.RouteManager) {
	//通用功能相关路由
	routeManager.RegisterCommonRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/rtoken", api.ReflashRtoken)
	})
	// 登录相关路由
	routeManager.RegisterLoginRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/login", api.LoginWithCode)
		rg.POST("/code", api.GetCode)
		rg.GET("/test", middleware.ReflashAtoken(), func(c *gin.Context) {
			if token, exists := c.Get(global.AUTH_ENUMS_ATOKEN); exists {
				response.NewResponse(c).Success(token)
			}
			//告诉后面的人如何拿到token里面的数据
			if data, exists := c.Get(global.TOKEN_USER_ID); exists {
				response.NewResponse(c).Success(data)
			}
		})
		//是否可以自动登录
		rg.POST("/auto", api.CheckAutoLogin)
		//用户自行退出登录
		rg.DELETE("/exit", api.ExitSystem)
	})

	// 展示常用设备页面相关操作路由
	routeManager.RegisterDevicesRoutes(func(rg *gin.RouterGroup) {

	})
	// 个人信息相关路由
	routeManager.RegisterProfileRoutes(func(rg *gin.RouterGroup) {

		// example
		rg.GET("/info", func(c *gin.Context) {

		})
		rg.GET("test", api.Test)
	})

	// 团队信息相关路由
	routeManager.RegisterTeamRoutes(func(rg *gin.RouterGroup) {

		//解析 jwt，获取 user_id
		rg.Use(middleware.ReflashAtoken())

		//获得权限组
		rg.GET("/power", api.GetPower)

		// 团队架构管理子路由
		TeamStructure := rg.Group("/structure")
		{
			//检验权限
			TeamStructure.Use(middleware.PermissionMiddleware())
			// 获取 完整团队架构
			TeamStructure.GET("/collection", api.GetTeamStructure)
			//保存 更改了的节点信息
			TeamStructure.PUT("/change", api.PutTeamNode)
			//新增团队
			TeamStructure.POST("/add", api.CreateTeam)
		}

		// 团队成员列表子路由
		MemberList := rg.Group("/memberlist")
		{
			//查询团队列表--用户基础信息
			MemberList.GET("/get", api.GetTeamMemberlist)
			//新增用户
			MemberList.POST("/post", middleware.PermissionMiddleware(), api.CreateTeamMember)
			//删除用户
			MemberList.DELETE("/delete", middleware.PermissionMiddleware(), api.DeleteTeamMember)
		}

		// 团队成员信息管理子路由
		MemberMsg := rg.Group("/membermsg")
		{
			//查询用户详细信息
			MemberMsg.GET("/details", api.GetMemberDetail)
			//给用户点赞/取消赞
			MemberMsg.PUT("/like", api.PutLikeCount)
			//编辑用户信息
			MemberMsg.PUT("/save", middleware.PermissionMiddleware(), api.PutTeamMember)
		}
	})
}

// messageRoutes 注册消息相关路由的具体处理函数
func messageRoutes(routeManager *manager.RouteManager) {
	routeManager.HandleMessageRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/set", api.SetMessage)
		rg.POST("/join", middleware.ReflashAtoken(), api.JoinMessage)
		rg.GET("/get", middleware.ReflashAtoken(), api.GetMessage)
		rg.POST("/markread", api.MarkReadMessage)
	})
}
