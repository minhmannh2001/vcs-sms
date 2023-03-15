package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/sms/helper"
)

// Report Server Information Intentionally godoc
// @Summary Report server information intentionally
// @Description Report server information
//
//	@Tags			Server CRUD
//
// @Accept json
// @Produce json
// @Success 200 {object} helper.Response
// @Router /api/servers/report [get]
// @Security apiKey
func ReportServer(ctx *gin.Context) {
	serverController.ReportServerStatus()
	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Sending server report successfully", helper.EmptyObj{}))
}
