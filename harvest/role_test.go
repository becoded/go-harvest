package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestRoleService_CreateRole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		request   *harvest.RoleCreateRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.Role
		wantErr   bool
	}{
		{
			name: "Valid Role Creation",
			request: &harvest.RoleCreateRequest{
				Name:    harvest.String("Role new"),
				UserIDs: harvest.Ints64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/roles", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testFormValues(t, r, values{})
					testBody(t, r, "role/create/body_1.json")
					testWriteResponse(t, w, "role/create/response_1.json")
				})
			},
			want: &harvest.Role{
				ID:        harvest.Int64(1),
				Name:      harvest.String("Role new"),
				UserIDs:   harvest.Ints64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
				CreatedAt: harvest.TimeTimeP(time.Date(2018, 1, 31, 20, 34, 30, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2018, 5, 31, 21, 34, 30, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name: "Invalid Role Creation - Missing Name",
			request: &harvest.RoleCreateRequest{
				UserIDs: harvest.Ints64([]int64{1, 2, 3}),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/roles", func(w http.ResponseWriter, _ *http.Request) {
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

			got, _, err := service.Role.Create(context.Background(), tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestRoleService_DeleteRole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		roleID    int64
		setupMock func(mux *http.ServeMux)
		wantErr   bool
	}{
		{
			name:   "Valid Role Deletion",
			roleID: 1,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/roles/1", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					testFormValues(t, r, values{})
					w.WriteHeader(http.StatusNoContent)
				})
			},
			wantErr: false,
		},
		{
			name:   "Role Not Found",
			roleID: 999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/roles/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					http.Error(w, `{"message":"Role not found"}`, http.StatusNotFound)
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

			_, err := service.Role.Delete(context.Background(), tt.roleID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRoleService_GetRole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		roleID    int64
		setupMock func(mux *http.ServeMux)
		want      *harvest.Role
		wantErr   bool
	}{
		{
			name:   "Valid Role Retrieval",
			roleID: 1,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/roles/1", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "role/get/body_1.json")
					testWriteResponse(t, w, "role/get/response_1.json")
				})
			},
			want: &harvest.Role{
				ID:        harvest.Int64(1),
				Name:      harvest.String("Role 1"),
				UserIDs:   harvest.Ints64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
				CreatedAt: harvest.TimeTimeP(time.Date(2018, 1, 31, 20, 34, 30, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2018, 5, 31, 21, 34, 30, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name:   "Role Not Found",
			roleID: 999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/roles/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					http.Error(w, `{"message":"Role not found"}`, http.StatusNotFound)
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

			got, _, err := service.Role.Get(context.Background(), tt.roleID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestRoleService_ListRoles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(mux *http.ServeMux)
		want      *harvest.RoleList
		wantErr   bool
	}{
		{
			name: "Valid Role List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/roles", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "role/list/body_1.json")
					testWriteResponse(t, w, "role/list/response_1.json")
				})
			},
			want: &harvest.RoleList{
				Roles: []*harvest.Role{
					{
						ID:        harvest.Int64(1),
						Name:      harvest.String("Role 1"),
						UserIDs:   harvest.Ints64([]int64{1, 2, 3, 4, 5}),
						CreatedAt: harvest.TimeTimeP(time.Date(2018, 1, 31, 20, 34, 30, 0, time.UTC)),
						UpdatedAt: harvest.TimeTimeP(time.Date(2018, 5, 31, 21, 34, 30, 0, time.UTC)),
					},
					{
						ID:        harvest.Int64(2),
						Name:      harvest.String("Role 2"),
						UserIDs:   harvest.Ints64([]int64{6, 7, 8, 9, 10}),
						CreatedAt: harvest.TimeTimeP(time.Date(2018, 3, 2, 10, 12, 13, 0, time.UTC)),
						UpdatedAt: harvest.TimeTimeP(time.Date(2018, 4, 30, 12, 13, 14, 0, time.UTC)),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					NextPage:     nil,
					PreviousPage: nil,
					Page:         harvest.Int(1),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/roles?page=1&per_page=100"),
						Next:     nil,
						Previous: nil,
						Last:     harvest.String("https://api.harvestapp.com/v2/roles?page=1&per_page=100"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Error Fetching Role List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/roles", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Role.List(context.Background(), &harvest.RoleListOptions{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestRoleService_UpdateRole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		roleID    int64
		request   *harvest.RoleUpdateRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.Role
		wantErr   bool
	}{
		{
			name:   "Valid Role Update",
			roleID: 1,
			request: &harvest.RoleUpdateRequest{
				Name:    harvest.String("Role update"),
				UserIDs: harvest.Ints64([]int64{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/roles/1", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					testFormValues(t, r, values{})
					testBody(t, r, "role/update/body_1.json")
					testWriteResponse(t, w, "role/update/response_1.json")
				})
			},
			want: &harvest.Role{
				ID:        harvest.Int64(1),
				Name:      harvest.String("Role update"),
				UserIDs:   harvest.Ints64([]int64{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}),
				CreatedAt: harvest.TimeTimeP(time.Date(2018, 1, 31, 20, 34, 30, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2018, 5, 31, 21, 34, 30, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name:   "Role Not Found",
			roleID: 999,
			request: &harvest.RoleUpdateRequest{
				Name:    harvest.String("Non-existent role"),
				UserIDs: harvest.Ints64([]int64{1, 2, 3}),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/roles/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					http.Error(w, `{"message":"Role not found"}`, http.StatusNotFound)
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

			got, _, err := service.Role.Update(context.Background(), tt.roleID, tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
