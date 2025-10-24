package main

import (
	"auth-service/internal/handlers"
	"auth-service/internal/service"
	"auth-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

var log logger.Log

func main() {
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	registrService := service.NewRegistrService()

	registrHandler := handlers.NewRegistrHandler(registrService)

	api := router.Group("/api")
	{
		regist := api.Group("/auth")
		{
			regist.POST("/register", registrHandler.NewRegistrService)
		}
	}

	router.Run(":8080")

}
