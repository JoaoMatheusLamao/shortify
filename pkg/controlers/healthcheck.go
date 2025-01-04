package controlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck is a controller that returns a 200 status code
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
