package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// HealthCheck handles the health check endpoint
func (s *Server) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "healthy",
		"time":   time.Now(),
	})
}
