package solus

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplicationsService_List(t *testing.T) {
	expected := ApplicationsResponse{
		Data: []Application{
			fakeApplication,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/applications", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Applications.List(context.Background())
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestApplicationsService_Create(t *testing.T) {
	data := ApplicationCreateRequest{
		Name:             "name",
		URL:              "url",
		IconID:           1,
		CloudInitVersion: "cloud init version",
		UserDataTemplate: "user data template",
		JSONSchema:       "json schema",
		IsVisible:        true,
		LoginLink: LoginLink{
			Type:    LoginLinkTypeURL,
			Content: "login link content",
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/applications", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakeApplication)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Applications.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, fakeApplication, actual)
}

func TestApplicationsService_Update(t *testing.T) {
	data := ApplicationCreateRequest{
		Name:             "updated-name",
		URL:              "https://example.com/app",
		IconID:           10,
		CloudInitVersion: "v3",
		UserDataTemplate: "template",
		JSONSchema:       "{\"type\":\"object\"}",
		IsVisible:        false,
		LoginLink: LoginLink{
			Type:    LoginLinkTypeURL,
			Content: "https://example.com/login",
		},
		AvailablePlans: []int{1, 2},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/applications/42", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusOK, fakeApplication)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Applications.Update(context.Background(), 42, data)
	require.NoError(t, err)
	require.Equal(t, fakeApplication, actual)
}

func TestApplicationFullRequestPreservesExplicitValues(t *testing.T) {
	legacy, err := json.Marshal(ApplicationCreateRequest{AvailablePlans: []int{}})
	require.NoError(t, err)
	assert.NotContains(t, string(legacy), "available_plans")
	full, err := json.Marshal(ApplicationWriteRequest{ApplicationCreateRequest: ApplicationCreateRequest{}, AvailablePlans: []int{}})
	require.NoError(t, err)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(full, &got))
	assert.Equal(t, []interface{}{}, got["available_plans"])
	assert.Equal(t, false, got["is_visible"])
	assert.Contains(t, got, "vs_disk_driver")
	assert.Nil(t, got["vs_disk_driver"])
	assert.Contains(t, got, "external_integration_id")
	assert.Nil(t, got["external_integration_id"])
	assert.Equal(t, LoginLinkType("js_code"), LoginLinkTypeJSCode)
	assert.Equal(t, LoginLinkType("js-code"), LoginLinkTypeJSCodeCanonical)
}

func TestApplicationsServiceFullLifecycle(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodGet, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, method, r.Method)
				if method == http.MethodPost {
					assert.Equal(t, "/applications", r.URL.Path)
				} else {
					assert.Equal(t, "/applications/7", r.URL.Path)
				}
				if method == http.MethodDelete {
					w.WriteHeader(204)
					return
				}
				if method == http.MethodPost || method == http.MethodPut {
					var body map[string]interface{}
					require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
					assert.Equal(t, float64(13), body["external_integration_id"])
					assert.Equal(t, []interface{}{}, body["available_plans"])
				}
				code := 200
				if method == http.MethodPost {
					code = 201
				}
				w.WriteHeader(code)
				_, _ = w.Write([]byte(`{"data":{"id":7,"vs_disk_driver":null,"vs_firmware":"efi","is_deprecated":true,"external_integration":{"id":13}}}`))
			})
			defer s.Close()
			service := createTestClient(t, s.URL).Applications
			id := 13
			req := ApplicationWriteRequest{ExternalIntegrationID: &id, AvailablePlans: []int{}}
			var app Application
			var err error
			switch method {
			case http.MethodPost:
				app, err = service.CreateFull(context.Background(), req)
			case http.MethodPut:
				app, err = service.UpdateFull(context.Background(), 7, req)
			case http.MethodGet:
				app, err = service.Get(context.Background(), 7)
			case http.MethodDelete:
				err = service.Delete(context.Background(), 7)
			}
			require.NoError(t, err)
			if method != http.MethodDelete {
				assert.Equal(t, 7, app.ID)
				assert.True(t, app.HasVSDiskDriver)
				assert.Nil(t, app.VSDiskDriver)
				assert.True(t, app.HasVSFirmware)
				require.NotNil(t, app.VSFirmware)
				assert.Equal(t, "efi", *app.VSFirmware)
				assert.True(t, app.IsDeprecated)
				require.NotNil(t, app.ExternalIntegration)
				assert.Equal(t, 13, app.ExternalIntegration.ID)
			}
		})
	}
}

func TestApplicationFieldPresence(t *testing.T) {
	var app Application
	require.NoError(t, json.Unmarshal([]byte(`{"id":1}`), &app))
	assert.False(t, app.HasVSDiskDriver)
	assert.False(t, app.HasVSFirmware)
	require.NoError(t, json.Unmarshal([]byte(`{"id":1,"vs_disk_driver":null,"vs_firmware":null}`), &app))
	assert.True(t, app.HasVSDiskDriver)
	assert.True(t, app.HasVSFirmware)
	assert.Nil(t, app.VSDiskDriver)
	assert.Nil(t, app.VSFirmware)
}

func TestApplicationDecodeErrorDoesNotExposeResponse(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":{"user_data_template":"sentinel-secret","id":"bad-id"}}`))
	})
	defer s.Close()
	_, err := createTestClient(t, s.URL).Applications.Get(context.Background(), 7)
	require.Error(t, err)
	assert.NotContains(t, err.Error(), "sentinel-secret")
}

func TestApplicationHTTPErrors(t *testing.T) {
	for _, code := range []int{401, 403, 404, 422, 500} {
		t.Run(http.StatusText(code), func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(code)
				_, _ = w.Write([]byte(`{"message":"sentinel-secret"}`))
			})
			defer s.Close()
			service := createTestClient(t, s.URL).Applications
			_, err := service.Get(context.Background(), 7)
			require.Error(t, err)
			assert.NotContains(t, err.Error(), "sentinel-secret")
			assert.Equal(t, code == 404, IsNotFound(err))
			require.Error(t, service.Delete(context.Background(), 7))
		})
	}
}

func TestApplicationFullWritesDoNotExposeDecodeErrors(t *testing.T) {
	for _, method := range []string{"create", "update"} {
		t.Run(method, func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if method == "create" {
					w.WriteHeader(201)
				}
				_, _ = w.Write([]byte(`{"data":{"user_data_template":"sentinel-secret","id":"bad"}}`))
			})
			defer s.Close()
			service := createTestClient(t, s.URL).Applications
			var err error
			if method == "create" {
				_, err = service.CreateFull(context.Background(), ApplicationWriteRequest{})
			} else {
				_, err = service.UpdateFull(context.Background(), 7, ApplicationWriteRequest{})
			}
			require.Error(t, err)
			assert.NotContains(t, err.Error(), "sentinel-secret")
		})
	}
}

func TestApplicationDecodeOverflowDoesNotExposeValue(t *testing.T) {
	const sensitiveNumber = "1234567890123456789012345678901234567890"
	for _, method := range []string{"get", "create", "update"} {
		t.Run(method, func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if method == "create" {
					w.WriteHeader(http.StatusCreated)
				}
				_, _ = w.Write([]byte(`{"data":{"id":` + sensitiveNumber + `}}`))
			})
			defer s.Close()
			service := createTestClient(t, s.URL).Applications
			var err error
			switch method {
			case "get":
				_, err = service.Get(context.Background(), 7)
			case "create":
				_, err = service.CreateFull(context.Background(), ApplicationWriteRequest{})
			case "update":
				_, err = service.UpdateFull(context.Background(), 7, ApplicationWriteRequest{})
			}
			require.Error(t, err)
			require.NotContains(t, err.Error(), sensitiveNumber)
			var field *json.UnmarshalTypeError
			require.True(t, errors.As(err, &field))
			assert.Equal(t, "invalid value", field.Value)
		})
	}
}
