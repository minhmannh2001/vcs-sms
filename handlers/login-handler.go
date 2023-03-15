package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/sms/controller"
	"github.com/minhmannh2001/sms/helper"
	"github.com/minhmannh2001/sms/service"
)

var (
	loginService    service.LoginService       = service.NewLoginService()
	jwtService      service.JWTService         = service.NewJWTService()
	loginController controller.LoginController = controller.NewLoginController(loginService, jwtService)
)

// JWT Login godoc
// @Summary Login Handler
// @Description API authentication and authorization.
//
//	@Tags			Login
//
// @Accept json
// @Produce json
//
//	@Param			credential	body		entity.Credentials	true	"Add credential"
//
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /login [post]
func LoginHandler(ctx *gin.Context) {
	token, expirationTime := loginController.Login(ctx)
	if token != "" {
		ctx.JSON(http.StatusOK, helper.BuildToken(true, token, helper.EmptyObj{}))
		_ = expirationTime
		// ctx.SetCookie("token", token, int(expirationTime), "/", "localhost", false, true)
	} else {
		ctx.JSON(http.StatusUnauthorized, helper.BuildErrorResponse("Authenticate unsuccessfully", "Check your information", helper.EmptyObj{}))
	}
}
