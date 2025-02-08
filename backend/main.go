package main

import (
	"backend/database"
	_ "backend/docs"
	"backend/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title    The Grid Backend API
// @version  0.0
// @host     localhost:8080
// @BasePath /
func main() {
	// initialize database
	database.InitDatabase()
	// initialize router
	router := gin.Default()
	routes.AuthRoutes(router)
	routes.UsersRoutes(router)
	// swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// start server
	router.Run(":8080")
}
