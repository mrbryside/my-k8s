package server

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/mrbryside/ketcd"
)

// test

// EtcdClient interface defines the methods needed from the etcd client
type EtcdClient interface {
	ListPods(ctx context.Context, namespace string) ([]ketcd.PodInfo, error)
	StorePod(ctx context.Context, pod ketcd.PodInfo) error
}

// Server represents the HTTP server with Fiber
type Server struct {
	app        *fiber.App
	etcdClient EtcdClient
}

// NewServer creates a new server instance with the given etcd client
func NewServer(etcdClient EtcdClient) Server {
	// Initialize Fiber with default config
	app := fiber.New(fiber.Config{
		ServerHeader: "KMaster Server",
		AppName:      "myown-k8s",
	})

	return Server{
		app:        app,
		etcdClient: etcdClient,
	}
}

// SetupRoutes configures all routes for the application
func (s Server) SetupRoutes() Server {
	// Root route
	s.app.Get("/", s.Root)

	// Pod management routes
	podGroup := s.app.Group("/api/pods")
	podGroup.Get("/", s.ListPods)
	podGroup.Post("/", s.CreatePod)

	// Health check route
	s.app.Get("/health", s.HealthCheck)

	return s
}

// Start initializes routes and starts the server
func (s Server) Start(addr string) error {
	// Start the server
	return s.app.Listen(addr)
}
