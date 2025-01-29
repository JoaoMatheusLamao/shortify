package routes

import (
	"shortify/internal/config"
	"shortify/internal/handlers/fieldmapper"
	"shortify/internal/handlers/healthcheck"
	"shortify/internal/handlers/redirect"
	"shortify/internal/handlers/shorten"

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

	mapperGroup := engine.Group("/mapper")
	mapperGroup.POST("/", fieldmapper.Mapper())

}
