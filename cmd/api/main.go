package main

import (
	"fmt"
	"log"
	"os"
	"shortify/internal/config"
	"shortify/internal/routes"
	"shortify/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	fmt.Println((os.Getenv("ENVIROMENT_EXEC")))

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error creating config: %v", err)
	}
	defer cfg.CloseAll()

	engine := middleware.SetupServer(cfg)

	routes.InitiateRoutes(engine, cfg)

	startServer(engine)
}

func startServer(engine *gin.Engine) {
	certFile, keyFile := getCertFiles()
	if certFile != "" && keyFile != "" {
		log.Println("Starting server with TLS...")
		if err := engine.RunTLS(":8080", certFile, keyFile); err != nil {
			log.Fatalf("Error starting TLS server: %v", err)
		}
	} else {
		log.Println("Starting server...")
		if err := engine.Run(":8080"); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}
	log.Println("Server started on port 8080")
}

// getPort returns the port from the environment variable PORT
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

// getCertFiles returns the certificate files from the environment variables
func getCertFiles() (string, string) {
	certFile := os.Getenv("CERT_FILE")
	keyFile := os.Getenv("KEY_FILE")
	return certFile, keyFile
}
