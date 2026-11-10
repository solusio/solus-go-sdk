package solus

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectsService_Members(t *testing.T) {
	expected := []ProjectMember{
		{
			ID:      1,
			UserID:  2,
			Email:   "user@example.com",
			Status:  "active",
			IsOwner: true,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/10/members", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		writeResponse(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Projects.Members(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestProjectsService_MemberCreate(t *testing.T) {
	data := ProjectMemberCreateRequest{
		Email: "user@example.com",
	}
	expected := ProjectMember{
		ID:     1,
		Email:  "user@example.com",
		Status: "invited",
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/10/members", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)
		writeResponse(t, w, http.StatusCreated, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Projects.MemberCreate(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestProjectsService_MemberDelete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/10/members/3", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).Projects.MemberDelete(context.Background(), 10, 3)
	require.NoError(t, err)
}

func TestProjectsService_MemberResendInvite(t *testing.T) {
	expected := ProjectMember{
		ID:     3,
		Email:  "user@example.com",
		Status: "invited",
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/10/members/3/resend_invite", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		writeResponse(t, w, http.StatusCreated, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Projects.MemberResendInvite(context.Background(), 10, 3)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestProjectsService_Join(t *testing.T) {
	expected := ProjectJoin{
		Email:        "user@example.com",
		IsRegistered: true,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/10/join", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		writeResponse(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Projects.Join(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestProjectsService_Leave(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/10/leave", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusOK)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).Projects.Leave(context.Background(), 10)
	require.NoError(t, err)
}

func TestProjectsService_Plans(t *testing.T) {
	expected := PlansResponse{
		Data: []Plan{fakePlan},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/10/plans", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assertRequestQuery(t, r, url.Values{
			"filter[search]":       []string{"name"},
			"filter[storage_type]": []string{"fb"},
		})
		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	filter := (&FilterPlans{}).ByName("name").ByStorageType(StorageTypeNameFB)
	actual, err := createTestClient(t, s.URL).Projects.Plans(context.Background(), 10, filter)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestProjectsService_SSHKeys(t *testing.T) {
	expected := []SSHKey{fakeSSHKey}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/10/ssh_keys", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		writeResponse(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Projects.SSHKeys(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestProjectsService_SSHKeyCreate(t *testing.T) {
	data := ProjectSSHKeyCreateRequest{
		Name: "example key",
		Body: "ssh-rsa AAA",
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/10/ssh_keys", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)
		writeResponse(t, w, http.StatusCreated, fakeSSHKey)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Projects.SSHKeyCreate(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, fakeSSHKey, actual)
}

func TestProjectsService_TokenPricing(t *testing.T) {
	expected := ProjectTokenPricing{
		UnitCost:     1.23,
		CurrencyCode: "USD",
		Taxes: []ProjectTokenPricingTax{
			{
				Label: "VAT",
				Rate:  10,
			},
		},
	}

	t.Run("billing enabled", func(t *testing.T) {
		s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/projects/10/token_pricing", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)
			writeResponse(t, w, http.StatusOK, expected)
		})
		defer s.Close()

		actual, err := createTestClient(t, s.URL).Projects.TokenPricing(context.Background(), 10)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("billing disabled", func(t *testing.T) {
		s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/projects/10/token_pricing", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)
			w.WriteHeader(http.StatusNoContent)
		})
		defer s.Close()

		actual, err := createTestClient(t, s.URL).Projects.TokenPricing(context.Background(), 10)
		require.NoError(t, err)
		require.Equal(t, ProjectTokenPricing{}, actual)
	})
}
