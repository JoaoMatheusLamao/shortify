package routes

import (
	"shortify/pkg/controlers"
	"shortify/pkg/persistence"

	"github.com/gin-gonic/gin"
)

// InitiateRoutes is a function that initializes the routes for the application
func InitiateRoutes(engine *gin.Engine, db *persistence.InMemoryDB) {
	engine.GET("/healthcheck", controlers.HealthCheck)

}
