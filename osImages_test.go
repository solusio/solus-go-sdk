package solus

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestIsValidCloudInitVersion(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cc := []string{
			string(CloudInitVersionV0),
			string(CloudInitVersionCentOS6),
			string(CloudInitVersionDebian9),
			string(CloudInitVersionV2),
			string(CloudInitVersionV2Alpine),
			string(CloudInitVersionV2Centos),
			string(CloudInitVersionV2Debian10),
			string(CloudInitVersionCloudBase),
		}

		for _, c := range cc {
			t.Run(c, func(t *testing.T) {
				actual := IsValidCloudInitVersion(c)
				assert.True(t, actual)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		assert.False(t, IsValidCloudInitVersion("invalid"))
		assert.False(t, IsValidCloudInitVersion("null"))
		assert.False(t, IsValidCloudInitVersion(""))
	})
}

func TestOsImagesService_List(t *testing.T) {
	expected := OsImagesResponse{
		Data: []OsImage{
			fakeOsImage,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/os_images", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assertRequestQuery(t, r, url.Values{
			"filter[search]": []string{"name"},
		})

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	f := (&FilterOsImages{}).ByName("name")

	actual, err := createTestClient(t, s.URL).OsImages.List(context.Background(), f)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestOsImagesService_Get(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/os_images/10", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, fakeOsImage)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).OsImages.Get(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeOsImage, actual)
}

func TestOsImagesService_Update(t *testing.T) {
	data := OsImageRequest{
		Name:      "name",
		IconID:    null.IntFrom(1),
		IsVisible: true,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/os_images/10", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusOK, fakeOsImage)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).OsImages.Update(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, fakeOsImage, actual)
}

func TestOsImagesService_Create(t *testing.T) {
	data := OsImageRequest{
		Name:      "name",
		IconID:    null.IntFrom(1),
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

		w.WriteHeader(http.StatusNoContent)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).OsImages.Delete(context.Background(), 10)
	require.NoError(t, err)
}

func TestOsImagesService_ListVersion(t *testing.T) {
	expected := []OsImageVersion{
		fakeKvmOsImageVersion,
		fakeVzOsImageVersion,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/os_images/10/versions", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).OsImages.ListVersion(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestOsImagesService_CreateVersion(t *testing.T) {
	data := OsImageVersionRequest{
		Position:           1,
		Version:            "version",
		VirtualizationType: VirtualizationTypeKVM,
		URL:                "http://example.com",
		CloudInitVersion:   CloudInitVersionV2,
		IsVisible:          true,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/os_images/10/versions", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakeKvmOsImageVersion)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).OsImages.CreateVersion(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, fakeKvmOsImageVersion, actual)
}
