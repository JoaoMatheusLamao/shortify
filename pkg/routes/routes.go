package routes

import (
	"shortify/pkg/config"
	"shortify/pkg/services/healthcheck"
	"shortify/pkg/services/redirect"
	"shortify/pkg/services/shorten"

	"github.com/gin-gonic/gin"
)

// InitiateRoutes is a function that initializes the routes for the application
func InitiateRoutes(engine *gin.Engine, cfg *config.Config) {
	healthGroup := engine.Group("/healthcheck")
	healthGroup.GET("/", healthcheck.HealthCheck)

	shortenGroup := engine.Group("/shorten")
	shortenGroup.POST("/", shorten.CreateShortenURL(cfg))

	redirectGroup := engine.Group("/r")
	redirectGroup.GET("/:shortURL", redirect.FindOriginalURLAndRedirect(cfg))

}
