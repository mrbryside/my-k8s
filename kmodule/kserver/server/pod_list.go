package server

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ListPods lists all pods or pods in a namespace
func (s *Server) ListPods(c *fiber.Ctx) error {
	// Get namespace from query parameter, empty means all namespaces
	namespace := c.Query("namespace", "")

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	pods, err := s.etcdClient.ListPods(ctx, namespace)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to list pods: %v", err),
		})
	}

	return c.JSON(pods)
}
