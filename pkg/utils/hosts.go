package utils

import (
	"os"

	"github.com/gin-gonic/gin"
)

// GetCurrentProtocolAndHost returns the current protocol and host
func GetCurrentProtocolAndHost(c *gin.Context) string {
	protocol := "http"
	if os.Getenv("SERVER_TLS") == "true" {
		protocol = "https"
	}
	host := c.Request.Host

	return protocol + "://" + host
}
