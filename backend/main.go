package main

import (
    "github.com/gin-gonic/gin"
    "backend/database"
    "backend/routes"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    _ "backend/docs"
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
    routes.UserRoutes(router)
    // swagger endpoint
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    // start server
    router.Run(":8080")
}
