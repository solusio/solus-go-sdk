package solus

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRetry(t *testing.T) {
	var counter int

	err := retry(func(attempt int) (bool, error) {
		counter++
		require.Equal(t, counter, attempt)
		return true, errors.New("fake error")
	})
	require.EqualError(t, err, errMaxRetriesReached.Error())
	require.Equal(t, maxRetries, counter)
}
