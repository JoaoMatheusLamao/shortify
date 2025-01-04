package routes

import (
	"shortify/pkg/config"
	"shortify/pkg/controlers"

	"github.com/gin-gonic/gin"
)

// InitiateRoutes is a function that initializes the routes for the application
func InitiateRoutes(engine *gin.Engine, rd *config.Redis) {
	engine.GET("/healthcheck", controlers.HealthCheck)

}
