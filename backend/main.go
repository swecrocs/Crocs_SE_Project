package main

import (
	"backend/database"
	_ "backend/docs"
	"backend/routes"

	"time"

	"github.com/gin-contrib/cors"
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
	// enable CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"}, // frontend hosting port
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// register routes
	routes.AuthRoutes(router)
	routes.UsersRoutes(router)
	routes.ProjectsRoutes(router)
	// swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// start server
	router.Run(":8080")
}
