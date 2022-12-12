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

func CORSConfig() cors.Config {
    corsConfig := cors.DefaultConfig()
    corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
		"http://192.168.0.100:3000",
		"http://onenus.s3-website-ap-southeast-1.amazonaws.com",
	}
    corsConfig.AllowCredentials = true
    corsConfig.AddAllowHeaders("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers", "Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization")
    corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE")
    return corsConfig
}

func main() {
	router := gin.Default()
	// router.Use(cors.Default())
	router.Use(cors.New(CORSConfig()))
	routes.SetupRoutes(router)
	router.Run()
}