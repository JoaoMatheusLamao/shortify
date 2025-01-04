package main

import (
	"log"
	"shortify/pkg/config"
	"shortify/pkg/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	log.Println("Starting server...")

	gin.SetMode(gin.DebugMode)

	engine, db := config.SetupServer()

	routes.InitiateRoutes(engine, db)

	engine.Run(":8080")
}
