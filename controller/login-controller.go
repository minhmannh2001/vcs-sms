package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/sms/entity"
	"github.com/minhmannh2001/sms/service"
)

type LoginController interface {
	Login(ctx *gin.Context) (string, int64)
}

type loginController struct {
	loginService service.LoginService
	jWtService   service.JWTService
}

func NewLoginController(loginService service.LoginService,
	jWtService service.JWTService) LoginController {
	return &loginController{
		loginService: loginService,
		jWtService:   jWtService,
	}
}

func (controller *loginController) Login(ctx *gin.Context) (string, int64) {
	var credentials entity.Credentials
	err := ctx.ShouldBind(&credentials)
	if err != nil {
		return "", 0
	}
	isAuthenticated := controller.loginService.Login(credentials.Username, credentials.Password)
	if isAuthenticated {
		return controller.jWtService.GenerateToken(credentials.Username, true)
	}
	return "", 0
}
