package solus

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestStorageTypesService_List(t *testing.T) {
	expected := []StorageType{
		{
			ID:      1,
			Name:    "foo",
			Formats: []ImageFormat{ImageFormatQCOW2},
		},
		{
			ID:      2,
			Name:    "bar",
			Formats: []ImageFormat{ImageFormatRaw},
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/storage_types", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).StorageTypes.List(context.Background())
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
