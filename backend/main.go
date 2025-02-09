package main

import (
	"backend/database"
	_ "backend/docs"
	"backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func main() {
	// Initialize the database connection
	database.InitDB()

	// Create a new Gin router
	r := gin.Default()

	// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register routes
	routes.RegisterRoutes(r)

	// Start the server
	r.Run(":8080")
}
