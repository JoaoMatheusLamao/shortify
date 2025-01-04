package main

import (
	"log"
	"shortify/pkg/config"
	"shortify/pkg/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.DebugMode)

	log.Println("Starting server...")

	engine := config.SetupServer()

	log.Println("Connecting to Redis...")
	rd, err := config.NewRedis()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	log.Println("Initializing routes...")
	routes.InitiateRoutes(engine, rd)

	engine.Run(":8080")

	log.Println("Server started on port 8080")
}
