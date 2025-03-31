package handlers

import "github.com/gofiber/fiber/v2"

// HealthCheckHandler handles health check requests
type HealthCheckHandler struct{}

// NewHealthCheckHandler initializes a new HealthCheckHandler
func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

// Handle returns the health status of the service
func (h *HealthCheckHandler) Handle(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "OK"})
}
