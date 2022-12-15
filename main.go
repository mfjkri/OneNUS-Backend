package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mfjkri/One-NUS-Backend/database"
	"github.com/mfjkri/One-NUS-Backend/routes"
	"github.com/mfjkri/One-NUS-Backend/utils"
)

func init() {
	if os.Getenv("DEPLOYED_MODE") == "" {
		utils.LoadEnv()
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

func SimulateLatency() gin.HandlerFunc {
	return func(c *gin.Context) {
		time.Sleep(time.Second * time.Duration(rand.Float64() * 4))
	}
}

func main() {
	router := gin.Default()

	// Middleware functions
	router.Use(cors.New(CORSConfig()))
	if os.Getenv("SIMULATE_LATENCY") == "true" {
		fmt.Println("Simulating latency for this server instance...")
		router.Use(SimulateLatency())
	}

	routes.SetupRoutes(router)
	router.Run()
}