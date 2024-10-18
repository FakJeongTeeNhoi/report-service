package router

import (
	"github.com/FakJeongTeeNhoi/report-system/controller"
	"github.com/gin-gonic/gin"
)

func ReportRouterGroup(server *gin.RouterGroup) {
	report := server.Group("/report")
	{
		report.GET("/:space", controller.GetSpaceStatistic)
	}
}
