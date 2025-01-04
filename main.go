package main

import (
	"log"
	"os"
	"shortify/pkg/config"
	"shortify/pkg/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.DebugMode)

	log.Println("Starting server...")

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	defer cfg.CloseAll()

	log.Println("Setting up server...")
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
