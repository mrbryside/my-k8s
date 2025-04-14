package ketcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// buildPodKey creates a consistent key format for pod storage
func buildPodKey(namespace, name string) string {
	return fmt.Sprintf("/pods/%s/%s", namespace, name)
}

// StorePod stores pod information in etcd
func (c EtcdClient) StorePod(ctx context.Context, pod PodInfo) error {
	pod.UpdatedAt = time.Now()
	if pod.CreatedAt.IsZero() {
		pod.CreatedAt = pod.UpdatedAt
	}

	data, err := json.Marshal(pod)
	if err != nil {
		return errors.Join(err, errors.New("unable to marshal pod data"))
	}

	return c.setKey(
		ctx,
		buildPodKey(pod.Namespace, pod.Name),
		string(data),
	)
}

// GetPod retrieves pod information from etcd
func (c EtcdClient) GetPod(ctx context.Context, namespace, name string) (*PodInfo, error) {
	key := buildPodKey(namespace, name)
	resp, err := c.Client.Get(ctx, key)
	if err != nil {
		return nil, errors.Join(err, errors.New("unable to get pod data"))
	}

	if len(resp.Kvs) == 0 {
		return nil, errors.Join(err, errors.New("pod not found"))
	}

	var pod PodInfo
	if err := json.Unmarshal(resp.Kvs[0].Value, &pod); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pod data: %w", err)
	}

	return &pod, nil
}

// ListPods lists all pods or pods in a specific namespace
func (c EtcdClient) ListPods(ctx context.Context, namespace string) ([]PodInfo, error) {
	var prefix string
	if namespace == "" {
		prefix = "/pods/"
	} else {
		prefix = fmt.Sprintf("/pods/%s/", namespace)
	}

	resp, err := c.Client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	pods := make([]PodInfo, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var pod PodInfo
		if err := json.Unmarshal(kv.Value, &pod); err != nil {
			return nil, fmt.Errorf("failed to unmarshal pod data: %w", err)
		}
		pods = append(pods, pod)
	}

	return pods, nil
}

// DeletePod removes a pod from etcd
func (c EtcdClient) DeletePod(ctx context.Context, namespace, name string) error {
	key := buildPodKey(namespace, name)
	_, err := c.Client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete pod %s/%s: %w", namespace, name, err)
	}
	return nil
}

// WatchPods watches for changes to all pods or pods in a specific namespace
func (c EtcdClient) WatchPods(ctx context.Context, namespace string) clientv3.WatchChan {
	var prefix string
	if namespace == "" {
		prefix = "/pods/"
	} else {
		prefix = fmt.Sprintf("/pods/%s/", namespace)
	}

	return c.Client.Watch(ctx, prefix, clientv3.WithPrefix())
}

// ProcessPodEvents demonstrates how to process pod events from the watch channel
func (c EtcdClient) ProcessPodEvents(_ context.Context, watchChan clientv3.WatchChan) {
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			var pod PodInfo
			if err := json.Unmarshal(event.Kv.Value, &pod); err != nil {
				fmt.Printf("Error unmarshalling pod data: %v\n", err)
				continue
			}

			switch event.Type {
			case clientv3.EventTypePut:
				if event.IsCreate() {
					fmt.Printf("Pod created: %s/%s with IP %s\n", pod.Namespace, pod.Name, pod.PodIP)
				} else {
					fmt.Printf("Pod updated: %s/%s, status: %s, IP: %s\n", pod.Namespace, pod.Name, pod.Status, pod.PodIP)
				}
			case clientv3.EventTypeDelete:
				fmt.Printf("Pod deleted: %s\n", string(event.Kv.Key))
			}
		}
	}
}
