package main

import (
	"log"
	"net/http"
	"os"
	"shortify/pkg/config"
	"shortify/pkg/routes"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error creating config: %v", err)
	}
	defer cfg.CloseAll()

	engine := config.SetupServer(cfg)

	log.Println("Initializing routes...")
	routes.InitiateRoutes(engine, cfg)

	startServer(createServer(engine))
}

func createServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         ":" + getPort(),
		Handler:      handler,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

func startServer(server *http.Server) {
	certFile, keyFile := getCertFiles()
	if certFile != "" && keyFile != "" {
		log.Println("Starting server with TLS...")
		log.Println("Cert file: " + certFile)
		log.Println("Key file: " + keyFile)
		if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
			log.Fatalf("Error starting TLS server: %v", err)
		}
	} else {
		log.Println("Starting server...")
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}
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

// getCertFiles returns the certificate files from the environment variables
func getCertFiles() (string, string) {
	certFile := os.Getenv("CERT_FILE")
	keyFile := os.Getenv("KEY_FILE")
	return certFile, keyFile
}
