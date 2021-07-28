package solus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterPlans(t *testing.T) {
	f := FilterPlans{}

	f.
		ByStorageType(StorageTypeNameNFS).
		ByImageFormat(ImageFormatQCOW2).
		ByName("name").
		ByDiskSize(42)

	require.Equal(t, map[string]string{
		"filter[storage_type]": "nfs",
		"filter[search]":       "name",
		"filter[image_format]": "qcow2",
		"filter[disk]":         "42",
	}, f.data)
}
