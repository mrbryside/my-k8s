// Package ketcd is a package that provides an etcd client for storing and retrieving data.
package ketcd

import (
	"context"
	"errors"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdClient wraps the etcd client for key-value operations
type EtcdClient struct {
	Client *clientv3.Client
}

// PodInfo represents the essential information about a pod
type PodInfo struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Status    string            `json:"status"`
	PodIP     string            `json:"podIP"`
	Labels    map[string]string `json:"labels,omitempty"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

// NewEtcdClient creates a new etcd client
func NewEtcdClient(host string) (ec EtcdClient, err error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{host},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return ec, errors.Join(err, fmt.Errorf("failed to create etcd client"))
	}
	// Initialize the EtcdClient struct
	return EtcdClient{
		Client: client,
	}, nil
}

// setKey sets a key-value pair in etcd
func (c EtcdClient) setKey(ctx context.Context, key, value string) error {
	_, err := c.Client.Put(ctx, key, value)
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

// watchKey watches a key for changes
// func (c EtcdClient) watchKey(ctx context.Context, key string) clientv3.WatchChan {
// 	return c.Client.Watch(ctx, key)
// }

// Close closes the etcd client
func (c EtcdClient) Close() error {
	return c.Client.Close()
}
