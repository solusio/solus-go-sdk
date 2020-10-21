package solus

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
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
