package solus

import (
	"context"
)

// ISOImagesService handles all available methods with ISO images.
type ISOImagesService service

// ISOImageVisibility represents available visibility types for ISO images.
type ISOImageVisibility string

const (
	// ISOImageVisibilityPublic indicates a public ISO image.
	ISOImageVisibilityPublic ISOImageVisibility = "public"

	// ISOImageVisibilityPrivate indicates a private ISO image.
	ISOImageVisibilityPrivate ISOImageVisibility = "private"
)

// ISOImageOSType represents available OS types for ISO images.
type ISOImageOSType string

const (
	// ISOImageOSTypeLinux indicates a Linux ISO image.
	ISOImageOSTypeLinux ISOImageOSType = "Linux"

	// ISOImageOSTypeWindows indicates a Windows ISO image.
	ISOImageOSTypeWindows ISOImageOSType = "Windows"
)

// ISOImageChecksumMethod represents available checksum methods for ISO images.
type ISOImageChecksumMethod string

const (
	// ISOImageChecksumMethodMD5 indicates MD5 checksum.
	ISOImageChecksumMethodMD5 ISOImageChecksumMethod = "md5"

	// ISOImageChecksumMethodSHA1 indicates SHA1 checksum.
	ISOImageChecksumMethodSHA1 ISOImageChecksumMethod = "sha1"

	// ISOImageChecksumMethodSHA256 indicates SHA256 checksum.
	ISOImageChecksumMethodSHA256 ISOImageChecksumMethod = "sha256"
)

// ISOImage represents an ISO image.
type ISOImage struct {
	ID                   int                    `json:"id"`
	Name                 string                 `json:"name"`
	Icon                 Icon                   `json:"icon"`
	User                 User                   `json:"user"`
	Visibility           ISOImageVisibility     `json:"visibility"`
	OSType               ISOImageOSType         `json:"os_type"`
	ISOURL               string                 `json:"iso_url"`
	UseTLS               bool                   `json:"use_tls"`
	ISOChecksumMethod    ISOImageChecksumMethod `json:"iso_checksum_method"`
	ISOChecksum          string                 `json:"iso_checksum"`
	Size                 int                    `json:"size"`
	ShowURLAndChecksum   bool                   `json:"show_url_and_checksum"`
	ShowTLS              bool                   `json:"show_tls"`
}

// ISOImagesResponse represents paginated list of ISO images.
// This cursor can be used for iterating over all available ISO images.
type ISOImagesResponse struct {
	paginatedResponse

	Data []ISOImage `json:"data"`
}

// List lists ISO images.
func (s *ISOImagesService) List(ctx context.Context, filter *FilterISOImages) (ISOImagesResponse, error) {
	resp := ISOImagesResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "iso_images", &resp, withFilter(filter.data))
}
