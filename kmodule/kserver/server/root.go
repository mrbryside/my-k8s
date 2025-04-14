package server

import "github.com/gofiber/fiber/v2"

// Root handles the root endpoint
func (s *Server) Root(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Welcome to myown-k8s API server",
		"status":  "running",
	})
}
