// Autogenerated file. Do not edit!

package solus

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIPBlocksResponse_Next(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		page := int32(1)

		s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			p := atomic.LoadInt32(&page)

			assert.Equal(t, "/ipblocks", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, strconv.Itoa(int(p)), r.URL.Query().Get("page"))

			if p == 3 {
				writeJSON(t, w, http.StatusOK, IPBlocksResponse{
					Data: []IPBlock{
						{
							ID: int(p),
						},
					},
					paginatedResponse: paginatedResponse{
						Links: ResponseLinks{
							Next: r.URL.String(),
						},
						Meta: ResponseMeta{
							CurrentPage: int(p),
							LastPage:    3,
						},
					},
				})
				return
			}
			atomic.AddInt32(&page, 1)

			q := r.URL.Query()
			q.Set("page", strconv.Itoa(int(p)+1))
			r.URL.RawQuery = q.Encode()

			writeJSON(t, w, http.StatusOK, IPBlocksResponse{
				paginatedResponse: paginatedResponse{
					Links: ResponseLinks{
						Next: r.URL.String(),
					},
					Meta: ResponseMeta{
						CurrentPage: int(p),
						LastPage:    3,
					},
				},
				Data: []IPBlock{{ID: int(p)}},
			})
		})
		defer s.Close()

		resp := IPBlocksResponse{
			paginatedResponse: paginatedResponse{
				Links: ResponseLinks{
					Next: fmt.Sprintf("%s/ipblocks?page=1", s.URL),
				},
				Meta: ResponseMeta{
					CurrentPage: 1,
					LastPage:    3,
				},
				service: &service{createTestClient(t, s.URL)},
			},
		}

		i := 1
		for resp.Next(context.Background()) {
			require.Equal(t, []IPBlock{{ID: i}}, resp.Data)
			i++
		}
		require.NoError(t, resp.err)
		require.Equal(t, 4, i, "Expects to get 3 entity, but got less")
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("failed to make request", func(t *testing.T) {
			asserter, addr := startBrokenTestServer(t)

			resp := IPBlocksResponse{
				paginatedResponse: paginatedResponse{
					Links: ResponseLinks{
						Next: fmt.Sprintf("%s/ipblocks?page=1", addr),
					},
					Meta: ResponseMeta{
						CurrentPage: 1,
						LastPage:    3,
					},
					service: &service{createTestClient(t, addr)},
				},
			}

			resp.Next(context.Background())
			asserter(t, http.MethodGet, "/ipblocks?page=1", resp.Err())
		})

		t.Run("invalid status code", func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/ipblocks", r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, strconv.Itoa(1), r.URL.Query().Get("page"))
				w.WriteHeader(http.StatusBadRequest)
			})
			defer s.Close()

			resp := IPBlocksResponse{
				paginatedResponse: paginatedResponse{
					Links: ResponseLinks{
						Next: fmt.Sprintf("%s/ipblocks?page=1", s.URL),
					},
					Meta: ResponseMeta{
						CurrentPage: 1,
						LastPage:    3,
					},
					service: &service{createTestClient(t, s.URL)},
				},
			}

			resp.Next(context.Background())
			assert.EqualError(t, resp.Err(), fmt.Sprintf(
				"HTTP GET %s/ipblocks?page=1 returns 400 status code",
				s.URL,
			))
		})

		t.Run("failed to unmarshal", func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/ipblocks", r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, strconv.Itoa(1), r.URL.Query().Get("page"))

				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("fake"))
				require.NoError(t, err)
			})
			defer s.Close()

			resp := IPBlocksResponse{
				paginatedResponse: paginatedResponse{
					Links: ResponseLinks{
						Next: fmt.Sprintf("%s/ipblocks?page=1", s.URL),
					},
					Meta: ResponseMeta{
						CurrentPage: 1,
						LastPage:    3,
					},
					service: &service{createTestClient(t, s.URL)},
				},
			}

			resp.Next(context.Background())
			assert.EqualError(
				t,
				resp.Err(),
				"failed to decode \"fake\": invalid character 'k' in literal false (expecting 'l')",
			)
		})
	})
}
