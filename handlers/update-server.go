package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/sms/entity"
	"github.com/minhmannh2001/sms/helper"
)

// Update Server godoc
// @Summary Update server
// @Description Update server with provided information
//
//	@Tags			Server CRUD
//
// @Accept json
// @Produce json
//
//	@Param			server	body		entity.Server	true	"Update server"
//	@Param			id	path		int	true	"Server ID"
//
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/server/{id} [put]
// @Security apiKey
func UpdateServer(ctx *gin.Context) {
	var server entity.Server
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot update server with this information", err.Error(), helper.EmptyObj{}))
		return
	}

	if err := ctx.ShouldBind(&server); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot update server with this information", err.Error(), helper.EmptyObj{}))
		return
	}

	_ = id

	server.Id = id

	err = serverController.UpdateServer(&server)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot update server with this information", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Updated server successfully!", helper.EmptyObj{}))

}
