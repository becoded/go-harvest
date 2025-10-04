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

func TestRole_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.Role
		want string
	}{
		{
			name: "Role with all fields",
			in: harvest.Role{
				ID:        harvest.Int64(1),
				Name:      harvest.String("Role 1"),
				UserIDs:   harvest.Ints64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
				CreatedAt: harvest.TimeTimeP(time.Date(2018, 1, 31, 20, 34, 30, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2018, 5, 31, 21, 34, 30, 0, time.UTC)),
			},
			want: `harvest.Role{ID:1, Name:"Role 1", UserIDs:[1 2 3 4 5 6 7 8 9 10], CreatedAt:time.Time{2018-01-31 20:34:30 +0000 UTC}, UpdatedAt:time.Time{2018-05-31 21:34:30 +0000 UTC}}`, //nolint: lll
		},
		{
			name: "Role with minimal fields",
			in: harvest.Role{
				ID:   harvest.Int64(2),
				Name: harvest.String("Role 2"),
			},
			want: `harvest.Role{ID:2, Name:"Role 2"}`,
		},
		{
			name: "Role with users only",
			in: harvest.Role{
				ID:      harvest.Int64(3),
				Name:    harvest.String("Role 3"),
				UserIDs: harvest.Ints64([]int64{1, 2, 3}),
			},
			want: `harvest.Role{ID:3, Name:"Role 3", UserIDs:[1 2 3]}`,
		},
		{
			name: "Role with timestamps",
			in: harvest.Role{
				ID:        harvest.Int64(4),
				Name:      harvest.String("Role 4"),
				CreatedAt: harvest.TimeTimeP(time.Date(2020, 6, 15, 10, 30, 0, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2020, 6, 20, 15, 45, 0, 0, time.UTC)),
			},
			want: `harvest.Role{ID:4, Name:"Role 4", CreatedAt:time.Time{2020-06-15 10:30:00 +0000 UTC}, UpdatedAt:time.Time{2020-06-20 15:45:00 +0000 UTC}}`, //nolint: lll
		},
		{
			name: "Empty Role",
			in:   harvest.Role{},
			want: `harvest.Role{}`,
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

func TestRoleList_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.RoleList
		want string
	}{
		{
			name: "RoleList with multiple roles",
			in: harvest.RoleList{
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
			want: `harvest.RoleList{Roles:[harvest.Role{ID:1, Name:"Role 1", UserIDs:[1 2 3 4 5], CreatedAt:time.Time{2018-01-31 20:34:30 +0000 UTC}, UpdatedAt:time.Time{2018-05-31 21:34:30 +0000 UTC}} harvest.Role{ID:2, Name:"Role 2", UserIDs:[6 7 8 9 10], CreatedAt:time.Time{2018-03-02 10:12:13 +0000 UTC}, UpdatedAt:time.Time{2018-04-30 12:13:14 +0000 UTC}}], Pagination:harvest.Pagination{PerPage:100, TotalPages:1, TotalEntries:2, Page:1, Links:harvest.PageLinks{First:"https://api.harvestapp.com/v2/roles?page=1&per_page=100", Last:"https://api.harvestapp.com/v2/roles?page=1&per_page=100"}}}`, //nolint: lll
		},
		{
			name: "RoleList with single role",
			in: harvest.RoleList{
				Roles: []*harvest.Role{
					{
						ID:   harvest.Int64(1),
						Name: harvest.String("Single Role"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.RoleList{Roles:[harvest.Role{ID:1, Name:"Single Role"}], Pagination:harvest.Pagination{PerPage:100, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
		},
		{
			name: "Empty RoleList",
			in: harvest.RoleList{
				Roles: []*harvest.Role{},
			},
			want: `harvest.RoleList{Roles:[], Pagination:harvest.Pagination{}}`,
		},
		{
			name: "RoleList with Links",
			in: harvest.RoleList{
				Roles: []*harvest.Role{
					{
						ID:   harvest.Int64(1),
						Name: harvest.String("Role with Links"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(50),
					TotalPages:   harvest.Int(3),
					TotalEntries: harvest.Int(150),
					NextPage:     harvest.Int(2),
					PreviousPage: nil,
					Page:         harvest.Int(1),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/roles?page=1"),
						Next:     harvest.String("https://api.harvestapp.com/v2/roles?page=2"),
						Previous: nil,
						Last:     harvest.String("https://api.harvestapp.com/v2/roles?page=3"),
					},
				},
			},
			want: `harvest.RoleList{Roles:[harvest.Role{ID:1, Name:"Role with Links"}], Pagination:harvest.Pagination{PerPage:50, TotalPages:3, TotalEntries:150, NextPage:2, Page:1, Links:harvest.PageLinks{First:"https://api.harvestapp.com/v2/roles?page=1", Next:"https://api.harvestapp.com/v2/roles?page=2", Last:"https://api.harvestapp.com/v2/roles?page=3"}}}`, //nolint: lll
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
