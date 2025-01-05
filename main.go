package main

import (
	"log"
	"os"
	"shortify/pkg/config"
	"shortify/pkg/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	gin.SetMode(gin.ReleaseMode)

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	defer cfg.CloseAll()

	engine := config.SetupServer(cfg)

	log.Println("Initializing routes...")
	routes.InitiateRoutes(engine, cfg)

	engine.Run(":" + getPort())

	log.Println("Server started on port " + getPort())
}

// getPort returns the port from the environment variable PORT
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
