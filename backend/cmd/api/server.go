package main

import (
	"github.com/gin-gonic/gin"
	"github.com/isa-ntana/to-do_case/internal/handler"
)

func setupServer(taskHandler *handler.TaskHandler) *gin.Engine {
	router := gin.Default()

	router.Use(corsMiddleware())

	apiGroup := router.Group("/api/v1")
	taskHandler.RegisterRoutes(apiGroup)

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		context.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(204)
			return
		}

		context.Next()
	}
}
