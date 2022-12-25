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
		"https://app.onenus.link",
	}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers", "Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization")
	corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE")
	return corsConfig
}

func main() {
	router := gin.Default()

	cmd := flag.String("cmd", "", "")
	flag.Parse()
	str_cmd := string(*cmd)

	router.Use(cors.New(CORSConfig()))

	routes.SetupRoutes(router)

	if str_cmd == "reset" {
		seed.DeleteAll()
	} else if str_cmd == "seed" {
		seed.GenerateData()

	} else if str_cmd == "update" {
		seed.UpdateUsers()
		seed.UpdatePosts()
	}

	router.Run()
}
