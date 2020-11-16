package solus

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestOsImagesService_List(t *testing.T) {
	expected := OsImagesResponse{
		Data: []OsImage{
			fakeOsImage,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/os_images", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).OsImages.List(context.Background())
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestOsImagesService_Create(t *testing.T) {
	data := OsImageCreateRequest{
		Name:      "name",
		Icon:      "icon",
		IsVisible: true,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/os_images", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakeOsImage)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).OsImages.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, fakeOsImage, actual)
}

func TestOsImagesService_Delete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/os_images/10", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(204)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).OsImages.Delete(context.Background(), 10)
	require.NoError(t, err)
}

func TestOsImagesService_OsImageVersionCreate(t *testing.T) {
	data := OsImageVersionRequest{
		Position:         1,
		Version:          "version",
		Url:              "http://example.com",
		CloudInitVersion: "v2",
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/os_images/10/versions", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakeOsImageVersion)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).OsImages.OsImageVersionCreate(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, fakeOsImageVersion, actual)
}
