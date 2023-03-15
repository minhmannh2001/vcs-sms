package handlers

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/sms/helper"
	"github.com/xuri/excelize/v2"
)

// Export Servers godoc
// @Summary Export servers
// @Description View or export servers based on url query
//
//	@Tags			Server CRUD
//
// @Accept  json
// @Produce  json
// @Param   from     query    int     false        "From"
// @Param   to      query    int     false        "To"
// @Param   perpage      query    int     false        "Account Per Page"
// @Param   sortby      query    string     false        "Sort By"
// @Param   order      query    string     false        "Order"
// @Param   filter      query    string     false        "Filter"
// @Param   export      query    string     false        "Export"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/servers [get]
// @Security apiKey
func GetOrExportServers(ctx *gin.Context) {

	var from, to, perpage int
	if ctx.Query("from") != "" {
		from, _ = strconv.Atoi(ctx.Query("from"))
	}

	if ctx.Query("to") != "" {
		to, _ = strconv.Atoi(ctx.Query("to"))
	}

	if ctx.Query("perpage") != "" {
		perpage, _ = strconv.Atoi(ctx.Query("perpage"))
	}

	sortby := ctx.Query("sortby")
	order := ctx.Query("order")
	filter := ctx.Query("filter")

	log.Println(from, to, perpage, sortby, order, filter)
	servers, err := serverController.ViewServers(from, to, perpage, sortby, order, filter)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Something is wrong", err.Error(), helper.EmptyObj{}))
		return
	}

	export := ctx.Query("export")
	if export != "" {
		new_file := excelize.NewFile()
		defer func() {
			if err := new_file.Close(); err != nil {
				ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Something is wrong", err.Error(), helper.EmptyObj{}))
				return
			}
		}()

		new_file.SetCellValue("Sheet1", "A1", "ID")
		new_file.SetCellValue("Sheet1", "B1", "Server Name")
		new_file.SetCellValue("Sheet1", "C1", "IPv4")
		new_file.SetCellValue("Sheet1", "D1", "User")
		new_file.SetCellValue("Sheet1", "E1", "Password")
		new_file.SetCellValue("Sheet1", "F1", "Status")
		new_file.SetCellValue("Sheet1", "G1", "Created At")
		new_file.SetCellValue("Sheet1", "H1", "Updated At")
		for idx, server := range servers {
			_ = server
			row_idx := strconv.FormatInt(int64(idx+2), 10)
			cell := "A" + row_idx
			new_file.SetCellValue("Sheet1", cell, server.Id)
			cell = "B" + row_idx
			new_file.SetCellValue("Sheet1", cell, server.Name)
			cell = "C" + row_idx
			new_file.SetCellValue("Sheet1", cell, server.Ipv4)
			cell = "D" + row_idx
			new_file.SetCellValue("Sheet1", cell, server.User)
			cell = "E" + row_idx
			new_file.SetCellValue("Sheet1", cell, server.Password)
			cell = "F" + row_idx
			new_file.SetCellValue("Sheet1", cell, server.Status)
			cell = "G" + row_idx
			new_file.SetCellValue("Sheet1", cell, server.CreatedAt)
			cell = "H" + row_idx
			new_file.SetCellValue("Sheet1", cell, server.UpdatedAt)
		}
		var b bytes.Buffer
		if err := new_file.Write(&b); err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Something is wrong", err.Error(), helper.EmptyObj{}))
			return
		}
		downloadName := time.Now().UTC().Format("servers-20060102150405.xlsx")
		ctx.Header("Content-Description", "File Transfer")
		ctx.Header("Content-Disposition", "attachment; filename="+downloadName)
		ctx.Data(http.StatusOK, "application/octet-stream", b.Bytes())
		if err := new_file.SaveAs("exported-servers.xlsx"); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		ctx.JSON(http.StatusOK, helper.BuildResponse(true, "exporting servers...", helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "exporting servers...", servers))

	// serverController.ExportServers()

}
