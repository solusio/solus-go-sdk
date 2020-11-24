// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package solus

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServersMigrationsService_Create(t *testing.T) {
	data := ServersMigrationRequest{
		IsLive:                       true,
		PreserveIPs:                  true,
		DestinationComputeResourceID: 1,
		Servers:                      []int{2, 3},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers_migrations", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakeServersMigration)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ServersMigrations.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, fakeServersMigration, actual)
}
