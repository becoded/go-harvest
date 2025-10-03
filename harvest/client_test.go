package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestClientService_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(mux *http.ServeMux)
		want      *harvest.ClientList
		wantErr   bool
	}{
		{
			name: "Valid Client List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "client/list/body_1.json")
					testWriteResponse(t, w, "client/list/response_1.json")
				})
			},
			want: &harvest.ClientList{
				Clients: []*harvest.Client{
					{
						ID:        harvest.Int64(1),
						Name:      harvest.String("Client 1"),
						IsActive:  harvest.Bool(true),
						Address:   harvest.String("Address line 1"),
						Currency:  harvest.String("EUR"),
						CreatedAt: harvest.TimeTimeP(time.Date(2018, 1, 31, 20, 34, 30, 0, time.UTC)),
						UpdatedAt: harvest.TimeTimeP(time.Date(2018, 5, 31, 21, 34, 30, 0, time.UTC)),
					},
					{
						ID:        harvest.Int64(2),
						Name:      harvest.String("Client 2"),
						IsActive:  harvest.Bool(false),
						Address:   harvest.String("Address line 2"),
						Currency:  harvest.String("EUR"),
						CreatedAt: harvest.TimeTimeP(time.Date(2018, 3, 2, 10, 12, 13, 0, time.UTC)),
						UpdatedAt: harvest.TimeTimeP(time.Date(2018, 4, 30, 12, 13, 14, 0, time.UTC)),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/clients?page=1&per_page=100"),
						Next:     nil,
						Previous: nil,
						Last:     harvest.String("https://api.harvestapp.com/v2/clients?page=1&per_page=100"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Error Fetching Client List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					http.Error(w, `{"message":"Internal Server Error"}`, http.StatusInternalServerError)
				})
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			tt.setupMock(mux)

			got, _, err := service.Client.List(context.Background(), &harvest.ClientListOptions{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestClientService_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		clientID  int64
		setupMock func(mux *http.ServeMux)
		want      *harvest.Client
		wantErr   bool
	}{
		{
			name:     "Valid Client Retrieval",
			clientID: 1,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/clients/1", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "client/get/body_1.json")
					testWriteResponse(t, w, "client/get/response_1.json")
				})
			},
			want: &harvest.Client{
				ID:        harvest.Int64(1),
				Name:      harvest.String("Client 1"),
				IsActive:  harvest.Bool(true),
				Address:   harvest.String("Address line 1"),
				Currency:  harvest.String("EUR"),
				CreatedAt: harvest.TimeTimeP(time.Date(2018, 1, 31, 20, 34, 30, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2018, 5, 31, 21, 34, 30, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name:     "Client Not Found",
			clientID: 999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/clients/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					http.Error(w, `{"message":"Client not found"}`, http.StatusNotFound)
				})
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			tt.setupMock(mux)

			got, _, err := service.Client.Get(context.Background(), tt.clientID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestClientService_CreateClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		request   *harvest.ClientCreateRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.Client
		wantErr   bool
	}{
		{
			name: "Valid Client Creation",
			request: &harvest.ClientCreateRequest{
				Name:     harvest.String("Client new"),
				IsActive: harvest.Bool(true),
				Address:  harvest.String("Address line 1"),
				Currency: harvest.String("EUR"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testFormValues(t, r, values{})
					testBody(t, r, "client/create/body_1.json")
					testWriteResponse(t, w, "client/create/response_1.json")
				})
			},
			want: &harvest.Client{
				ID:        harvest.Int64(1),
				Name:      harvest.String("Client 1"),
				IsActive:  harvest.Bool(true),
				Address:   harvest.String("Address line 1"),
				Currency:  harvest.String("EUR"),
				CreatedAt: harvest.TimeTimeP(time.Date(2018, 1, 31, 20, 34, 30, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2018, 5, 31, 21, 34, 30, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name: "Invalid Client Creation - Missing Name",
			request: &harvest.ClientCreateRequest{
				IsActive: harvest.Bool(true),
				Address:  harvest.String("Address line 1"),
				Currency: harvest.String("EUR"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/clients", func(w http.ResponseWriter, _ *http.Request) {
					http.Error(w, `{"message":"Name is required"}`, http.StatusBadRequest)
				})
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			tt.setupMock(mux)

			got, _, err := service.Client.Create(context.Background(), tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestClientService_UpdateClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		clientID  int64
		request   *harvest.ClientUpdateRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.Client
		wantErr   bool
	}{
		{
			name:     "Valid Client Update",
			clientID: 1,
			request: &harvest.ClientUpdateRequest{
				Name:     harvest.String("Client updated"),
				IsActive: harvest.Bool(true),
				Address:  harvest.String("Address line 1"),
				Currency: harvest.String("EUR"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/clients/1", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					testFormValues(t, r, values{})
					testBody(t, r, "client/update/body_1.json")
					testWriteResponse(t, w, "client/update/response_1.json")
				})
			},
			want: &harvest.Client{
				ID:        harvest.Int64(1),
				Name:      harvest.String("Client updated"),
				IsActive:  harvest.Bool(true),
				Address:   harvest.String("Address line 1"),
				Currency:  harvest.String("EUR"),
				CreatedAt: harvest.TimeTimeP(time.Date(2018, 1, 31, 20, 34, 30, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2018, 5, 31, 21, 34, 30, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name:     "Client Not Found",
			clientID: 999,
			request: &harvest.ClientUpdateRequest{
				Name:     harvest.String("Non-existent client"),
				IsActive: harvest.Bool(false),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/clients/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					http.Error(w, `{"message":"Client not found"}`, http.StatusNotFound)
				})
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			tt.setupMock(mux)

			got, _, err := service.Client.Update(context.Background(), tt.clientID, tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestClientService_DeleteClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		clientID  int64
		setupMock func(mux *http.ServeMux)
		wantErr   bool
	}{
		{
			name:     "Valid Client Deletion",
			clientID: 1,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/clients/1", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					testFormValues(t, r, values{})
					w.WriteHeader(http.StatusNoContent)
				})
			},
			wantErr: false,
		},
		{
			name:     "Client Not Found",
			clientID: 999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/clients/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					http.Error(w, `{"message":"Client not found"}`, http.StatusNotFound)
				})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			tt.setupMock(mux)

			_, err := service.Client.Delete(context.Background(), tt.clientID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.Client
		want string
	}{
		{
			name: "Client with all fields",
			in: harvest.Client{
				ID:        harvest.Int64(123),
				Name:      harvest.String("Test Client"),
				IsActive:  harvest.Bool(true),
				Address:   harvest.String("123 Main St"),
				Currency:  harvest.String("USD"),
				CreatedAt: harvest.TimeTimeP(time.Date(2018, 1, 31, 20, 34, 30, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2018, 5, 31, 21, 34, 30, 0, time.UTC)),
			},
			want: `harvest.Client{ID:123, Name:"Test Client", IsActive:true, Address:"123 Main St", Currency:"USD", CreatedAt:time.Time{2018-01-31 20:34:30 +0000 UTC}, UpdatedAt:time.Time{2018-05-31 21:34:30 +0000 UTC}}`, //nolint: lll
		},
		{
			name: "Client with minimal fields",
			in: harvest.Client{
				ID:   harvest.Int64(456),
				Name: harvest.String("Minimal Client"),
			},
			want: `harvest.Client{ID:456, Name:"Minimal Client"}`,
		},
		{
			name: "Client with nil fields",
			in: harvest.Client{
				ID:       harvest.Int64(789),
				IsActive: harvest.Bool(false),
			},
			want: `harvest.Client{ID:789, IsActive:false}`,
		},
		{
			name: "Empty Client",
			in:   harvest.Client{},
			want: `harvest.Client{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.in.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClientList_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.ClientList
		want string
	}{
		{
			name: "ClientList with multiple clients",
			in: harvest.ClientList{
				Clients: []*harvest.Client{
					{
						ID:       harvest.Int64(1),
						Name:     harvest.String("Client One"),
						IsActive: harvest.Bool(true),
					},
					{
						ID:       harvest.Int64(2),
						Name:     harvest.String("Client Two"),
						IsActive: harvest.Bool(false),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.ClientList{Clients:[harvest.Client{ID:1, Name:"Client One", IsActive:true} harvest.Client{ID:2, Name:"Client Two", IsActive:false}], Pagination:harvest.Pagination{PerPage:100, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
		},
		{
			name: "ClientList with single client",
			in: harvest.ClientList{
				Clients: []*harvest.Client{
					{
						ID:   harvest.Int64(999),
						Name: harvest.String("Solo Client"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(50),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.ClientList{Clients:[harvest.Client{ID:999, Name:"Solo Client"}], Pagination:harvest.Pagination{PerPage:50, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
		},
		{
			name: "Empty ClientList",
			in: harvest.ClientList{
				Clients: []*harvest.Client{},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(0),
					TotalEntries: harvest.Int(0),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.ClientList{Clients:[], Pagination:harvest.Pagination{PerPage:100, TotalPages:0, TotalEntries:0, Page:1}}`, //nolint: lll
		},
		{
			name: "ClientList with Links",
			in: harvest.ClientList{
				Clients: []*harvest.Client{
					{
						ID:   harvest.Int64(100),
						Name: harvest.String("Test Client"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(25),
					TotalPages:   harvest.Int(2),
					TotalEntries: harvest.Int(50),
					Page:         harvest.Int(1),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/clients?page=1&per_page=25"),
						Next:     harvest.String("https://api.harvestapp.com/v2/clients?page=2&per_page=25"),
						Previous: nil,
						Last:     harvest.String("https://api.harvestapp.com/v2/clients?page=2&per_page=25"),
					},
				},
			},
			want: `harvest.ClientList{Clients:[harvest.Client{ID:100, Name:"Test Client"}], Pagination:harvest.Pagination{PerPage:25, TotalPages:2, TotalEntries:50, Page:1, Links:harvest.PageLinks{First:"https://api.harvestapp.com/v2/clients?page=1&per_page=25", Next:"https://api.harvestapp.com/v2/clients?page=2&per_page=25", Last:"https://api.harvestapp.com/v2/clients?page=2&per_page=25"}}}`, //nolint: lll
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.in.String()
			assert.Equal(t, tt.want, got)
		})
	}
}
