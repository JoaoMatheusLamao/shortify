package routes

import (
	"shortify/pkg/controlers"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// InitiateRoutes is a function that initializes the routes for the application
func InitiateRoutes(engine *gin.Engine, rd *redis.Client) {
	engine.GET("/healthcheck", controlers.HealthCheck)

}
