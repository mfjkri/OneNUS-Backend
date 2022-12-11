package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mfjkri/One-NUS-Backend/database"
	"github.com/mfjkri/One-NUS-Backend/routes"
	"github.com/mfjkri/One-NUS-Backend/run"
)

func init() {
	if os.Getenv("DEPLOYED_MODE") == "" {
		run.LoadEnv()
	}
	database.Connect()
	database.Migrate()

}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	routes.SetupRoutes(router)
	router.Run()
}