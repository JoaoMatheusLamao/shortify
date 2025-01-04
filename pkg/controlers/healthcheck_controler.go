package controlers

import (
	"net/http"
	"shortify/pkg/services/healthcheck"

	"github.com/gin-gonic/gin"
)

// HealthCheck is a controller that returns a 200 status code
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, healthcheck.HealthCheck())
}
