package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/sms/helper"
)

// Delete server by ID godoc
// @Summary Delete server by ID
// @Description Using ID to delete server
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
// @Router /api/server/{id} [delete]
// @Security apiKey
func DeleteServer(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot delete server with this id", err.Error(), helper.EmptyObj{}))
		return
	}

	err = serverController.DeleteServer(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot delete server with this id", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Deleted server successfully!", helper.EmptyObj{}))

}
