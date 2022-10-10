/*
 * @Author: jlz
 * @Date: 2022-08-11 09:51:06
 * @LastEditTime: 2022-08-31 16:30:24
 * @LastEditors: jlz
 * @Description:
 */

package router

import (
	_ "dnds_go/docs"
	"dnds_go/global"
	"dnds_go/middleware"
	"dnds_go/provider"
	constantService "dnds_go/src/service/constant"
	gatewayService "dnds_go/src/service/gateway"
	loginService "dnds_go/src/service/login"
	testService "dnds_go/src/service/test"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routers(r *gin.RouterGroup) {
	// 测试路由组
	TestRouter := r.Group("/test")
	{
		TestRouter.GET("/conf", testService.TestConf)
		TestRouter.POST("/add", testService.Create)
		TestRouter.GET("/hello", testService.Create)
		TestRouter.GET("/mqSend", testService.MqSend)
		TestRouter.GET("/exportExcel", testService.ExportEexcel)
		TestRouter.POST("/testValidator", testService.TestValidator)
		TestRouter.POST("/beeValidator", testService.TestBeeValidator)
		TestRouter.GET("/goroutineTest", testService.GoroutineTest)
	}
	// 登录
	LogintRouter := r.Group("/")
	{
		LogintRouter.POST("/admin/login", loginService.Login)
		LogintRouter.POST("/createUser", loginService.CreateUser)
	}
	// 常量
	ConstantRouter := r.Group("/constant")
	{
		ConstantRouter.POST("/addConstant", constantService.CreateConstant)
		ConstantRouter.POST("/editConstant", constantService.UpdateConstant)
		ConstantRouter.GET("/getConstantList", constantService.GetConstant)
		ConstantRouter.POST("/deleteConstant", constantService.DeleteConstant)
		// 常量子项
		ConstantRouter.POST("/addDetail", constantService.CreateCstDetail)
		ConstantRouter.POST("/editDetail", constantService.UpdateCstDetail)
		ConstantRouter.POST("/deleteDetail", constantService.DeleteCstDetail)
		ConstantRouter.GET("/getDetailList", constantService.GetDetailAllByCoding)

	}
	// 事件常量
	EventCstRouter := r.Group("/alarmEvent")
	{
		EventCstRouter.POST("/addEvent", constantService.CreateEventCst)
		EventCstRouter.POST("/setEventDetail", constantService.UpdateEventCst)
		EventCstRouter.GET("/getEventList", constantService.GetEventCstList)
		EventCstRouter.POST("/deleteEvent", constantService.DeleteEventCst)
		EventCstRouter.POST("/testData", constantService.TestData)
		EventCstRouter.GET("/getSensorType", constantService.GetSensorTypeByPoint)
	}
	GatewayRouter := r.Group("/gateway")
	{
		GatewayRouter.POST("/addGateway", gatewayService.Create)
		GatewayRouter.POST("/editGateway", gatewayService.Update)
		GatewayRouter.GET("/getAvailableGateway", gatewayService.GetList)
	}
	// websocket路径
	WsGroup := r.Group("/ws")
	{
		// channel: websocket分组名称
		WsGroup.GET("/:channel", provider.WebsocketManager.WsClient)
	}
	// swagger服务
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
func InitRouter() *gin.Engine {
	if viper.GetString("Server.Mode") == "production" {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
	}

	var Router = gin.New()

	// 引入全局中间件 日志、错误恢复、跨域
	Router.Use(
		middleware.GinLogger(global.Logger),
		middleware.GinRecovery(global.Logger, true),
		middleware.CorsAuth(), // 允许跨域
	)
	// // 公共路由
	PublicGroup := Router.Group("/")
	{
		Routers(PublicGroup)
	}

	return Router
}
