package healthcheck

// HealthCheck is a service that returns a map with a status key
func HealthCheck() map[string]string {
	return map[string]string{"status": "ok"}
}
