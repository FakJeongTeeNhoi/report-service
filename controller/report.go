package controller

import (
	"github.com/FakJeongTeeNhoi/report-service/model"
	"github.com/FakJeongTeeNhoi/report-service/model/response"
	"github.com/gin-gonic/gin"
)

func GetSpaceStatistic(c *gin.Context) {
	spaceName := c.Param("spaceName")

	reports, err := model.GetReportsBySpace(spaceName)
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
