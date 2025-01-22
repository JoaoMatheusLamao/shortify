package healthcheck

import "github.com/gin-gonic/gin"

// HealthCheck is a function that returns a 200 status code
func HealthCheck(c *gin.Context) {
	c.JSON(200, map[string]string{"status": "ok"})
}
