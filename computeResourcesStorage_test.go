package solus

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComputeResourcesService_StorageCreate(t *testing.T) {
	data := ComputeResourceStorageCreateRequest{
		Type:                    StorageTypeNameFB,
		Path:                    "fake path",
		ThinPool:                "fake thinpool",
		IsAvailableForBalancing: true,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/42/storages", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakeStorage)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ComputeResources.StorageCreate(context.Background(), 42, data)
	require.NoError(t, err)
	require.Equal(t, fakeStorage, actual)
}

func TestComputeResourcesService_StorageList(t *testing.T) {
	expected := []Storage{
		fakeStorage,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/42/storages", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ComputeResources.StorageList(context.Background(), 42)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
