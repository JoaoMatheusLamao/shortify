package healthcheck

import "github.com/gin-gonic/gin"

// HealthCheck is a service that returns a map with a status key
func HealthCheck() map[string]string {
	return map[string]string{"status": "ok"}
}

// HealthCheckControler is a function that returns a 200 status code
func HealthCheckControler(c *gin.Context) {
	c.JSON(200, HealthCheck())
}
