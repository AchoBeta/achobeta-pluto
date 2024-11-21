package manager

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"tgwp/internal/middleware"
)

//主要管理路由组和中间件的注册

// PathHandler 是一个用于注册路由组的函数类型
type PathHandler func(rg *gin.RouterGroup)

// Middleware 是一个用于生成中间件的函数类型
type Middleware func() gin.HandlerFunc

// RouteManager 管理不同的路由组，按业务功能分组
type RouteManager struct {
	LoginRoutes   *gin.RouterGroup // 登录相关的路由组
	ProfileRoutes *gin.RouterGroup // 个人信息相关的路由组
	TeamRoutes    *gin.RouterGroup // 团队信息相关的路由组
	CommonRoutes  *gin.RouterGroup //特殊功能相关的路由组
	MessageRoutes *gin.RouterGroup // 消息相关的路由组
	DevicesRoutes *gin.RouterGroup // 展示常用设备页面相关操作路由
}

// NewRouteManager 创建一个新的 RouteManager 实例，包含各业务功能的路由组
func NewRouteManager(router *gin.Engine) *RouteManager {
	return &RouteManager{
		LoginRoutes:   router.Group("/api/login"),   // 初始化登录路由组
		ProfileRoutes: router.Group("/api/profile"), // 初始化个人信息路由组
		TeamRoutes:    router.Group("/api/team"),    // 初始化团队信息路由组
		CommonRoutes:  router.Group("/api/common"),  //通用功能相关的路由组
		MessageRoutes: router.Group("/api/message"), //通用功能相关的路由组
		DevicesRoutes: router.Group("/api/devices"), // 展示常用设备页面相关操作路由
	}
}

// RegisterLoginRoutes 注册登录相关的路由处理函数
func (rm *RouteManager) RegisterLoginRoutes(handler PathHandler) {
	handler(rm.LoginRoutes)
}

// RegisterProfileRoutes 注册个人信息相关的路由处理函数
func (rm *RouteManager) RegisterProfileRoutes(handler PathHandler) {
	handler(rm.ProfileRoutes)
}

// RegisterTeamRoutes 注册团队信息相关的路由处理函数
func (rm *RouteManager) RegisterTeamRoutes(handler PathHandler) {
	handler(rm.TeamRoutes)
}

// RegisterCommonRoutes通用功能相关的路由组
func (rm *RouteManager) RegisterCommonRoutes(handler PathHandler) {
	handler(rm.CommonRoutes)
}

// 展示常用设备页面相关操作路由
func (rm *RouteManager) RegisterDevicesRoutes(handler PathHandler) {
	handler(rm.DevicesRoutes)
}

// RegisterMiddleware 根据组名为对应的路由组注册中间件
// group 参数为 "login"、"profile"、"team"或"Common"，分别对应不同的路由组
func (rm *RouteManager) RegisterMiddleware(group string, middleware Middleware) {
	switch group {
	case "login":
		rm.LoginRoutes.Use(middleware())
	case "profile":
		rm.ProfileRoutes.Use(middleware())
	case "team":
		rm.TeamRoutes.Use(middleware())
	case "common":
		rm.CommonRoutes.Use(middleware())
	case "devices":
		rm.DevicesRoutes.Use(middleware())
	}
}

// RequestGlobalMiddleware 注册全局中间件，应用于所有路由
func RequestGlobalMiddleware(r *gin.Engine) {
	r.Use(requestid.New())
	r.Use(middleware.AddTraceId())
	r.Use(middleware.Cors())
}

// HandleMessageRoutes 处理消息相关的路由处理函数
func (rm *RouteManager) HandleMessageRoutes(handler PathHandler) {
	handler(rm.MessageRoutes)
}
