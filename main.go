package main

import (
	"flag"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mfjkri/One-NUS-Backend/database"
	"github.com/mfjkri/One-NUS-Backend/routes"
	"github.com/mfjkri/One-NUS-Backend/seed"
	"github.com/mfjkri/One-NUS-Backend/utils"
)

func init() {
	utils.LoadEnv()
	database.Connect()
	database.Migrate()
}

func CORSConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
		"http://192.168.0.100:3000",
		"https://app.onenus.link",
	}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers", "Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization")
	corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE")
	return corsConfig
}

func main() {
	// Some command utilities
	cmd := flag.String("cmd", "", "")
	flag.Parse()
	str_cmd := string(*cmd)

	// Set Gin mode based on env var
	gin.SetMode(os.Getenv("GIN_MODE"))

	// Create a new router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(CORSConfig()))

	// Initialize the Routes
	routes.SetupRoutes(router)

	// Check for any command parameters used
	if str_cmd == "reset" {
		seed.DeleteAll()
	} else if str_cmd == "seed" {
		seed.GenerateData()

	} else if str_cmd == "update" {
		seed.UpdateUsers()
		seed.UpdatePosts()
	}

	// Start listening
	router.Run()
}
