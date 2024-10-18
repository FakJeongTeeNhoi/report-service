package controller

import (
	"github.com/FakJeongTeeNhoi/report-system/model"
	"github.com/FakJeongTeeNhoi/report-system/model/response"
	"github.com/gin-gonic/gin"
)

func GetSpaceStatistic(c *gin.Context) {
	space := c.Param("space")

	reports, err := model.GetReportsBySpace(space)
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
