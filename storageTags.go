package solus

import (
	"context"
	"fmt"
)

// StorageTagsService handles all available methods with storage tags.
type StorageTagsService service

// Delete deletes specified storage tag.
func (s *StorageTagsService) Delete(ctx context.Context, id int) error {
	return s.client.syncDelete(ctx, fmt.Sprintf("storage_tags/%d", id))
}
