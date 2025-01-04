package routes

import (
	"shortify/pkg/config"
	"shortify/pkg/services/healthcheck"
	"shortify/pkg/services/redirect"
	"shortify/pkg/services/shorten"

	"github.com/gin-gonic/gin"
)

// InitiateRoutes is a function that initializes the routes for the application
func InitiateRoutes(engine *gin.Engine, rd *config.Redis) {
	engine.GET("/healthcheck", healthcheck.HealthCheckControler)

	engine.POST("/shorten", shorten.ShortenControler(rd))

	engine.GET("/r/:shortURL", redirect.RedirectControler(rd))

}
