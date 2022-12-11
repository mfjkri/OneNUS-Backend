package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mfjkri/One-NUS-Backend/database"
	"github.com/mfjkri/One-NUS-Backend/routes"
	"github.com/mfjkri/One-NUS-Backend/run"
)

func init() {
	run.LoadEnv()
	database.Connect()
	database.Migrate()

}

func main() {
	r := gin.Default()
	routes.SetupRoutes(r)
	r.Run()
}