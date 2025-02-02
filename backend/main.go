package main

import (
    "github.com/gin-gonic/gin"
    "backend/database"
    "backend/routes"
)

func main() {
    // initialize database
	database.InitDatabase()
    // initialize router
    router := gin.Default()
    routes.UserRoutes(router)
    // start server
    router.Run(":8080")
}
