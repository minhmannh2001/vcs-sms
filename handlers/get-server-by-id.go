package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/sms/cache"
	"github.com/minhmannh2001/sms/controller"
	"github.com/minhmannh2001/sms/database"
	"github.com/minhmannh2001/sms/entity"
	"github.com/minhmannh2001/sms/helper"
	"github.com/minhmannh2001/sms/service"
)

var (
	serverDatabase   database.SMSDatabase        = database.NewSMSDatabase()
	serverService    service.ServerService       = service.NewServerService(serverDatabase)
	serverCache      cache.ServerCache           = cache.NewRedisCache("redis:6379", 0, 20)
	serverController controller.ServerController = controller.NewServerController(serverService, serverCache)
)

// Get server by ID godoc
// @Summary Get server by ID
// @Description Using ID to check server's existence
//
//	@Tags			Server CRUD
//
// @Accept json
// @Produce json
//
//	@Param			id	path		int	true	"Server ID"
//
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/server/{id} [get]
// @Security apiKey
func GetServerById(ctx *gin.Context) {
	var server entity.Server
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Server doesn't exist", err.Error(), helper.EmptyObj{}))
		return
	}

	server, err = serverController.ViewServer(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Server doesn't exist", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Server exist", server))
}
