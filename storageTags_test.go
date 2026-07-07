package solus

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageTagsService_Delete(t *testing.T) {
	t.Run("no_content", func(t *testing.T) {
		s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/storage_tags/10", r.URL.Path)
			assert.Equal(t, http.MethodDelete, r.Method)
			w.WriteHeader(http.StatusNoContent)
		})
		defer s.Close()

		err := createTestClient(t, s.URL).StorageTags.Delete(context.Background(), 10)
		require.NoError(t, err)
	})

	t.Run("ok_with_body", func(t *testing.T) {
		s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/storage_tags/10", r.URL.Path)
			assert.Equal(t, http.MethodDelete, r.Method)
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"status":"ok"}`))
			require.NoError(t, err)
		})
		defer s.Close()

		err := createTestClient(t, s.URL).StorageTags.Delete(context.Background(), 10)
		require.NoError(t, err)
	})
}
