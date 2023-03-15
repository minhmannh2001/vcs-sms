package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/sms/entity"
	"github.com/minhmannh2001/sms/helper"
	"github.com/xuri/excelize/v2"
)

// Import Servers godoc
// @Summary Import servers
// @Description Import servers within a excel file
//
//	@Tags			Server CRUD
//
// @Accept json
// @Produce json
//
//	@Param			server	formData	file		true	"Update server"
//
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/servers [post]
// @Security apiKey
func ImportServers(ctx *gin.Context) {
	upload_file, _, err := ctx.Request.FormFile("file")
	path := filepath.Join(".", "upload-files")
	err = os.MkdirAll(path, os.ModePerm)
	fullPath := path + "/servers-to-import.xlsx"
	new_file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot import servers", err.Error(), helper.EmptyObj{}))
		return
	}
	defer new_file.Close()
	_, err = io.Copy(new_file, upload_file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot import servers", err.Error(), helper.EmptyObj{}))
		return
	}

	servers_file, err := excelize.OpenFile(fullPath)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot import servers", err.Error(), helper.EmptyObj{}))
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := servers_file.Close(); err != nil {
			ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot import servers", err.Error(), helper.EmptyObj{}))
			return
		}
	}()
	var created_successful_servers []entity.Server
	var created_unsuccessful_servers []entity.Server
	// Get all the rows in the Sheet1.
	rows, err := servers_file.GetRows("Sheet1")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Cannot import servers", err.Error(), helper.EmptyObj{}))
		return
	}
	skip_row := true
	for _, row := range rows {
		if skip_row == true {
			skip_row = false
			continue
		}
		new_server := entity.Server{}
		count := 1
		for _, colCell := range row {
			switch count {
			case 1:
				new_server.Name = colCell
			case 2:
				new_server.Ipv4 = colCell
			case 3:
				new_server.User = colCell
			case 4:
				new_server.Password = colCell
			}
			count += 1
		}
		log.Println(new_server)
		existed_ip := serverController.CheckServerExistence(new_server.Ipv4)
		existed_name := serverController.CheckServerExistence(new_server.Name)
		if existed_ip || existed_name {
			created_unsuccessful_servers = append(created_unsuccessful_servers, new_server)
			continue
		} else {
			err = serverController.CreateServer(&new_server)
			if err != nil {
				created_unsuccessful_servers = append(created_unsuccessful_servers, new_server)
				continue
			}
			created_successful_servers = append(created_successful_servers, new_server)
		}
	}
	serverController.ImportServers()

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Import task is finished", gin.H{"created_servers": created_successful_servers, "uncreated_servers": created_unsuccessful_servers}))
}
