package router

import (
	"github.com/FakJeongTeeNhoi/report-service/controller"
	"github.com/gin-gonic/gin"
)

func ReportRouterGroup(server *gin.RouterGroup) {
	report := server.Group("/report")
	{
		report.GET("/:spaceID", controller.GetSpaceStatistic)
	}
}
