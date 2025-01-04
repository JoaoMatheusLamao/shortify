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
	healthGroup := engine.Group("/healthcheck")
	healthGroup.GET("/", healthcheck.HealthCheck)

	shortenGroup := engine.Group("/shorten")
	shortenGroup.POST("/", shorten.CreateShortenURL(rd))

	redirectGroup := engine.Group("/r")
	redirectGroup.GET("/:shortURL", redirect.FindOriginalURLAndRedirect(rd))

}
