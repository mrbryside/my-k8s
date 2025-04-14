package server

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mrbryside/ketcd"
)

// CreatePod lists all pods or pods in a namespace
func (s *Server) CreatePod(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	podInfo := ketcd.PodInfo{
		Name: c.FormValue("name"),
	}

	err := s.etcdClient.StorePod(ctx, podInfo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to create pod: %v", err),
		})
	}

	return c.JSON("success")
}
