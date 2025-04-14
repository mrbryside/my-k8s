// Package main
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mrbryside/ketcd"
	"github.com/mrbryside/kserver/server"
)

var (
	etcdClient ketcd.EtcdClient
	srv        server.Server
)

// NewEtcdClient creates a new etcd client
func initEtcdClient(endpoint string) {
	client, err := ketcd.NewEtcdClient(endpoint)
	if err != nil {
		panic(fmt.Sprintf("Failed to create etcd client: %v", err))
	}
	etcdClient = client
}

func initServer() {
	srv = server.
		NewServer(etcdClient).
		SetupRoutes()
}

func main() {
	// Initialize etcd client
	initEtcdClient(os.Getenv("ETCD_HOST"))
	initServer()

	defer func() {
		_ = etcdClient.Close()
	}()

	// Start the server
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
