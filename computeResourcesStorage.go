package solus

import (
	"context"
	"fmt"
)

// ComputeResourceStorageCreateRequest represents available properties for creating
// a new storage on a compute resource.
type ComputeResourceStorageCreateRequest struct {
	Type                    StorageTypeName `json:"type"`
	Path                    string          `json:"path"`
	ThinPool                string          `json:"thin_pool,omitempty"`
	IsAvailableForBalancing bool            `json:"is_available_for_balancing"`
}

// ComputeResourceStorageUpdateRequest represents available properties for
// updating a storage on a compute resource.
type ComputeResourceStorageUpdateRequest struct {
	IsAvailableForBalancing bool `json:"is_available_for_balancing"`
	StorageTagID            int  `json:"storage_tag_id,omitempty"`
}

// StorageCreate creates a new storage for the specified compute resource.
func (s *ComputeResourcesService) StorageCreate(
	ctx context.Context,
	id int,
	data ComputeResourceStorageCreateRequest,
) (Storage, error) {
	var resp storageResponse
	return resp.Data, s.client.create(ctx, fmt.Sprintf("compute_resources/%d/storages", id), data, &resp)
}

// StorageList lists storages for the specified compute resource.
func (s *ComputeResourcesService) StorageList(ctx context.Context, id int) ([]Storage, error) {
	var resp struct {
		Data []Storage `json:"data"`
	}
	return resp.Data, s.client.get(ctx, fmt.Sprintf("compute_resources/%d/storages", id), &resp)
}

// StorageUpdate updates a storage on the specified compute resource.
func (s *ComputeResourcesService) StorageUpdate(
	ctx context.Context,
	id int,
	storageID int,
	data ComputeResourceStorageUpdateRequest,
) ([]Storage, error) {
	var resp struct {
		Data []Storage `json:"data"`
	}
	return resp.Data, s.client.update(ctx, fmt.Sprintf("compute_resources/%d/storages/%d", id, storageID), data, &resp)
}
