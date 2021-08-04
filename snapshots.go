// Copyright 1999-2021. Plesk International GmbH. All rights reserved.

package solus

import (
	"context"
	"fmt"
)

type SnapshotsService service

type Snapshot struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	// Size a size of snapshot in Gb.
	Size      float64        `json:"size"`
	Status    SnapshotStatus `json:"status"`
	CreatedAt string         `json:"created_at"`
}

type SnapshotStatus string

const (
	SnapshotStatusAvailable  SnapshotStatus = "available"
	SnapshotStatusProcessing SnapshotStatus = "processing"
	SnapshotStatusFailed     SnapshotStatus = "failed"
)

type SnapshotRequest struct {
	Name string `json:"name"`
}

type snapshotResponse struct {
	Data Snapshot `json:"data"`
}

func (s *SnapshotsService) Get(ctx context.Context, id int) (Snapshot, error) {
	var resp snapshotResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("snapshots/%d", id), &resp)
}

func (s *SnapshotsService) Revert(ctx context.Context, id int) (Task, error) {
	return s.client.asyncPost(ctx, fmt.Sprintf("snapshots/%d/revert", id))
}

func (s *SnapshotsService) Delete(ctx context.Context, id int) (Task, error) {
	return s.client.asyncDelete(ctx, fmt.Sprintf("snapshots/%d", id))
}
