package solus

import (
	"context"
	"fmt"
)

// SnapshotsService handles all available methods with snapshots.
type StorageService service

// ShortStorage represents storage embedded in a compute resource summary.
// Credentials are intentionally excluded: the API can return an object, an empty
// array, or null for that field, and server reads do not need credential data.
type ShortStorage struct {
	ID     int         `json:"id"`
	Name   string      `json:"name"`
	Type   StorageType `json:"type"`
	Status string      `json:"status"`
}

// Storage represents a storage.
type Storage struct {
	ID                      int                    `json:"id"`
	Name                    string                 `json:"name"`
	Type                    StorageType            `json:"type"`
	Path                    string                 `json:"path"`
	Mount                   string                 `json:"mount"`
	ThinPool                string                 `json:"thin_pool"`
	IsAvailableForBalancing bool                   `json:"is_available_for_balancing"`
	ServersCount            int                    `json:"servers_count"`
	ComputeResourcesCount   int                    `json:"compute_resources_count"`
	FreeSpace               float64                `json:"free_space"`
	Credentials             map[string]interface{} `json:"credentials"`
}

type storageResponse struct {
	Data Storage `json:"data"`
}

// Get gets specified storage.
func (s *StorageService) Get(ctx context.Context, id int) (Storage, error) {
	var resp storageResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("storages/%d", id), &resp)
}

// Delete deletes specified storage.
func (s *StorageService) Delete(ctx context.Context, id int) error {
	return s.client.syncDelete(ctx, fmt.Sprintf("storages/%d", id))
}
