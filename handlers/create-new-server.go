package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/sms/entity"
	"github.com/minhmannh2001/sms/helper"
)

// Create New Server godoc
// @Summary Create new server
// @Description Create new server with provided information
//
//	@Tags			Server CRUD
//
// @Accept json
// @Produce json
//
//	@Param			server	body		entity.Server	true	"Add server"
//
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/server/ [post]
// @Security apiKey
func CreateNewServer(ctx *gin.Context) {
	var server entity.Server
	if err := ctx.ShouldBind(&server); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot create new server with this information", err.Error(), helper.EmptyObj{}))
		return
	}
	err := serverController.CreateServer(&server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot create new server with this information", err.Error(), helper.EmptyObj{}))
		return
	}
	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Created successfully new server!", helper.EmptyObj{}))
}

// https://stackoverflow.com/questions/56176814/how-to-add-jwt-auth-to-swagger-go-echo-swaggo-swag
