package main

import (
	"net/http"

	"github.com/minhmannh2001/sms/cache"
	"github.com/minhmannh2001/sms/controller"
	"github.com/minhmannh2001/sms/database"
	_ "github.com/minhmannh2001/sms/docs"
	"github.com/minhmannh2001/sms/handlers"
	"github.com/minhmannh2001/sms/middlewares"
	"github.com/minhmannh2001/sms/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	serverDatabase database.SMSDatabase  = database.NewSMSDatabase()
	serverService  service.ServerService = service.NewServerService(serverDatabase)
	loginService   service.LoginService  = service.NewLoginService()
	jwtService     service.JWTService    = service.NewJWTService()
	serverCache    cache.ServerCache     = cache.NewRedisCache("redis:6379", 0, 20)

	serverController     controller.ServerController     = controller.NewServerController(serverService, serverCache)
	loginController      controller.LoginController      = controller.NewLoginController(loginService, jwtService)
	connectionController controller.ConnectionController = controller.NewConnectionController()
)

// @title     Server Management System
// @version         1.0
// @description     A server management service API in Go using Gin framework.

// @contact.name   Nguyen Minh Manh
// @contact.url    https://www.facebook.com/minhmannh2001/
// @contact.email  nguyenminhmannh2001@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
//@securityDefinitions.apikey apiKey
//@in header
//@name Authorization

func main() {
	defer serverDatabase.CloseDBConnection()

	router := gin.Default()
	go connectionController.SendCheckConnection()
	go connectionController.ReceiveCheckConnection()
	go serverController.AutoReportServerStatus()

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Login Endpoint: Authentication + Token creation
	router.POST("/login", middlewares.DefaultStructuredLogger(), handlers.LoginHandler)

	router.GET("/checkExistence/:ip", middlewares.DefaultStructuredLogger(), func(ctx *gin.Context) {
		ip := ctx.Param("ip")
		existed := serverController.CheckServerExistence(ip)
		ctx.JSON(http.StatusOK, gin.H{"existed": existed})
	})

	apiRouters := router.Group("/api", middlewares.AuthorizeJWT(), middlewares.DefaultStructuredLogger())
	{
		apiRouters.GET("/server/:id", handlers.GetServerById)

		apiRouters.POST("/server", handlers.CreateNewServer)

		apiRouters.PUT("/server/:id", handlers.UpdateServer)

		apiRouters.DELETE("/server/:id", handlers.DeleteServer)

		apiRouters.GET("/servers/report", handlers.ReportServer)

		apiRouters.GET("/servers", handlers.GetOrExportServers)

		apiRouters.POST("/servers", handlers.ImportServers)
	}

	router.Run()

}

// https://dev.to/santosh/how-to-integrate-swagger-ui-in-go-backend-gin-edition-2cbd
// https://github.com/swaggo/gin-swagger/blob/master/example/basic/api/api.go
// https://github.com/swaggo/swag/blob/master/example/celler/controller/accounts.go
// Link Bao bao: https://docs.google.com/document/d/1leXYrD0SMG9F9McZM5hV9o3bBAcRctotoms4_3Qnryc/edit#
