package utils

import "github.com/gin-gonic/gin"

// GetCurrentProtocolAndHost returns the current protocol and host
func GetCurrentProtocolAndHost(c *gin.Context) string {
	protocol := "http"
	if c.Request.TLS != nil {
		protocol = "https"
	}
	host := c.Request.Host

	return protocol + "://" + host
}
