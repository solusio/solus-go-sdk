package solus

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fakeISOImage = ISOImage{
	ID:                 1,
	Name:               "Ubuntu Rescue",
	Icon:               fakeIcon,
	User:               fakeUser,
	Visibility:         ISOImageVisibilityPublic,
	OSType:             ISOImageOSTypeLinux,
	ISOURL:             "https://example.com/ubuntu.iso",
	UseTLS:             true,
	ISOChecksumMethod:  ISOImageChecksumMethodSHA256,
	ISOChecksum:        "abc123",
	Size:               2048,
	ShowURLAndChecksum: true,
	ShowTLS:            true,
}

func TestISOImagesService_List(t *testing.T) {
	expected := ISOImagesResponse{
		Data: []ISOImage{fakeISOImage},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/iso_images", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assertRequestQuery(t, r, url.Values{
			"filter[search]": []string{"Ubuntu"},
		})

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	f := (&FilterISOImages{}).ByName("Ubuntu")

	actual, err := createTestClient(t, s.URL).ISOImages.List(context.Background(), f)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}
