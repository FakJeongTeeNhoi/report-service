package main

import (
	"fmt"

	"github.com/FakJeongTeeNhoi/report-service/controller"
	"github.com/FakJeongTeeNhoi/report-service/router"
	"github.com/FakJeongTeeNhoi/report-service/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting server...")

	// TODO: Connect to database using gorm
	service.ConnectMongoDB()

	server := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	server.Use(cors.New(corsConfig))

	api := server.Group("/api")
	go controller.StartConsumeDataFromQueue("Receiver", []string{"reservation.*"})

	// TODO: Add routes here
	router.ReportRouterGroup(api)

	err = server.Run(":3021")
	if err != nil {
		panic(err)
	}

	// TODO: Add graceful shutdown
}
