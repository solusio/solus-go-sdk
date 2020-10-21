package solus

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFilter_add(t *testing.T) {
	f := filter{}
	f.add("foo", "bar")
	f.add("fizz", "buzz")
	f.add("foo", "foo")

	require.Equal(t, map[string]string{
		"foo":  "foo",
		"fizz": "buzz",
	}, f.data)
}

func TestFilter_addInt(t *testing.T) {
	f := filter{}
	f.addInt("foo", 42)
	f.addInt("fizz", 1337)
	f.addInt("foo", 100)

	require.Equal(t, map[string]string{
		"foo":  "100",
		"fizz": "1337",
	}, f.data)
}
