package controller

import (
	"github.com/FakJeongTeeNhoi/report-service/model"
	"github.com/FakJeongTeeNhoi/report-service/model/response"

	"github.com/evangwt/go-csv"
	"github.com/gin-gonic/gin"

	"bytes"
	"net/http"
)

func GetSpaceStatistic(c *gin.Context) {
	spaceID := c.Param("spaceID")

	reports, err := model.GetReportsBySpace(spaceID)
	if err != nil {
		response.NotFound("No reports found for the specified space").AbortWithError(c)
		return
	}

	c.JSON(200, response.CommonResponse{
		Success: true,
	}.AddInterfaces(map[string]interface{}{
		"reports": reports,
	}))
}

func DownloadSpaceStatistic(c *gin.Context) {
	spaceID := c.Param("spaceID")

	reports, err := model.GetReportsBySpace(spaceID)
	if err != nil {
		response.NotFound("No reports found for the specified space").AbortWithError(c)
		return
	}

	csvData := bytes.NewBuffer(nil)
	w := csv.NewWriter(csvData)

	if err := w.Write([]string{
		"ID",
		"Reservation ID",
		"Room ID",
		"Space ID",
		"Space Name",
		"Status",
		"Start Datetime",
		"End Datetime",
	}); err != nil {
		response.InternalServerError("Failed to write csv header").AbortWithError(c)
	}

	for _, report := range reports {
		if err := w.Write(report.ArrayOfString()); err != nil {
			response.InternalServerError("Failed to write csv data").AbortWithError(c)
		}
	}

	err = w.Flush()
	if err != nil {
		response.InternalServerError("Failed to flush csv writer").AbortWithError(c)
	}

	c.Writer.Header().Set("Content-Disposition", "attachment; filename=report.csv")
	c.Writer.Header().Set("Content-Type", "text/csv")
	c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
	c.Data(http.StatusOK, "text/csv; charset=utf-8", []byte(csvData.String()))
}
