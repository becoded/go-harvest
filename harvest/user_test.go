package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestUserService_Create(t *testing.T) {
	t.Parallel()

	createdOne := time.Date(
		2020, 1, 25, 19, 20, 46, 0, time.UTC)
	updatedOne := time.Date(
		2020, 1, 25, 19, 20, 57, 0, time.UTC)

	roles := []string{"Project Manager"}

	tests := []struct {
		name      string
		data      *harvest.UserCreateRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.User
		wantErr   bool
	}{
		{
			name: "create user",
			data: &harvest.UserCreateRequest{
				FirstName:        harvest.String("George"),
				LastName:         harvest.String("Frank"),
				Email:            harvest.String("george@example.com"),
				IsProjectManager: harvest.Bool(true),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testFormValues(t, r, values{})
					testBody(t, r, "user/create/body_1.json")
					testWriteResponse(t, w, "user/create/response_1.json")
				})
			},
			want: &harvest.User{
				ID:                           harvest.Int64(3),
				FirstName:                    harvest.String("Gary"),
				LastName:                     harvest.String("Frank"),
				Email:                        harvest.String("george@example.com"),
				Timezone:                     harvest.String("Eastern Time (US & Canada)"),
				HasAccessToAllFutureProjects: harvest.Bool(false),
				IsContractor:                 harvest.Bool(false),
				IsAdmin:                      harvest.Bool(false),
				IsActive:                     harvest.Bool(true),
				CanSeeRates:                  harvest.Bool(false),
				CanCreateProjects:            harvest.Bool(false),
				CanCreateInvoices:            harvest.Bool(false),
				Telephone:                    harvest.String(""),
				IsProjectManager:             harvest.Bool(true),
				CreatedAt:                    &createdOne,
				UpdatedAt:                    &updatedOne,
				WeeklyCapacity:               harvest.Int(126000),
				DefaultHourlyRate:            harvest.Float64(0),
				CostRate:                     harvest.Float64(0),
				Roles:                        &roles,
				AvatarURL: harvest.String(
					"https://{ACCOUNT_SUBDOMAIN}.harvestapp.com/assets/profile_images/big_ben.png?1485372046",
				),
			},
			wantErr: false,
		},
		{
			name: "error creating user",
			data: &harvest.UserCreateRequest{
				FirstName:        harvest.String("George"),
				LastName:         harvest.String("Frank"),
				Email:            harvest.String("george@example.com"),
				IsProjectManager: harvest.Bool(true),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
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

			user, _, err := service.User.Create(context.Background(), tt.data)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, user)
			}
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		userID    int64
		setupMock func(mux *http.ServeMux)
		wantErr   bool
	}{
		{
			name:   "delete 1",
			userID: 1,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/users/1", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					testFormValues(t, r, values{})
					testBody(t, r, "user/delete/body_1.json")
					testWriteResponse(t, w, "user/delete/response_1.json")
				})
			},
			wantErr: false,
		},
		{
			name:   "error deleting user",
			userID: 1,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/users/1", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					http.Error(w, `{"message":"Internal Server Error"}`, http.StatusInternalServerError)
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

			_, err := service.User.Delete(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_Get(t *testing.T) {
	t.Parallel()

	createdOne := time.Date(
		2020, 5, 1, 22, 34, 41, 0, time.UTC)
	updatedOne := time.Date(
		2020, 5, 1, 22, 34, 52, 0, time.UTC)

	roles := []string{"Developer"}

	tests := []struct {
		name      string
		userID    int64
		setupMock func(mux *http.ServeMux)
		want      *harvest.User
		wantErr   bool
	}{
		{
			name:   "get 3230547",
			userID: 3230547,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/users/3230547", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "user/get/body_1.json")
					testWriteResponse(t, w, "user/get/response_1.json")
				})
			},
			want: &harvest.User{
				ID:                           harvest.Int64(3230547),
				IsActive:                     harvest.Bool(true),
				FirstName:                    harvest.String("Jim"),
				LastName:                     harvest.String("Allen"),
				Email:                        harvest.String("jimallen@example.com"),
				Telephone:                    harvest.String(""),
				Timezone:                     harvest.String("Mountain Time (US & Canada)"),
				HasAccessToAllFutureProjects: harvest.Bool(false),
				IsContractor:                 harvest.Bool(false),
				IsAdmin:                      harvest.Bool(false),
				IsProjectManager:             harvest.Bool(false),
				CanSeeRates:                  harvest.Bool(false),
				CanCreateInvoices:            harvest.Bool(false),
				CanCreateProjects:            harvest.Bool(false),
				Roles:                        &roles,
				CostRate:                     harvest.Float64(50),
				DefaultHourlyRate:            harvest.Float64(100),
				WeeklyCapacity:               harvest.Int(126000),
				AvatarURL: harvest.String(
					"https://cache.harvestapp.com/assets/profile_images/abraj_albait_towers.png?1498516481",
				),

				CreatedAt: &createdOne,
				UpdatedAt: &updatedOne,
			},
			wantErr: false,
		},
		{
			name:   "error getting user",
			userID: 3230547,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/users/3230547", func(w http.ResponseWriter, r *http.Request) {
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

			user, _, err := service.User.Get(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, user)
			}
		})
	}
}

func TestUserService_Current(t *testing.T) {
	t.Parallel()

	createdOne := time.Date(
		2020, 5, 1, 20, 41, 0o0, 0, time.UTC)
	updatedOne := time.Date(
		2020, 5, 1, 20, 42, 25, 0, time.UTC)

	roles := []string{"Founder", "CEO"}

	tests := []struct {
		name      string
		setupMock func(mux *http.ServeMux)
		want      *harvest.User
		wantErr   bool
	}{
		{
			name: "get current user",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "user/current/body_1.json")
					testWriteResponse(t, w, "user/current/response_1.json")
				})
			},
			want: &harvest.User{
				ID:                           harvest.Int64(1782884),
				FirstName:                    harvest.String("Bob"),
				LastName:                     harvest.String("Powell"),
				Email:                        harvest.String("bobpowell@example.com"),
				Telephone:                    harvest.String(""),
				Timezone:                     harvest.String("Mountain Time (US & Canada)"),
				HasAccessToAllFutureProjects: harvest.Bool(false),
				IsContractor:                 harvest.Bool(false),
				IsAdmin:                      harvest.Bool(true),
				IsProjectManager:             harvest.Bool(false),
				CanSeeRates:                  harvest.Bool(true),
				CanCreateProjects:            harvest.Bool(true),
				CanCreateInvoices:            harvest.Bool(true),
				IsActive:                     harvest.Bool(true),
				CreatedAt:                    &createdOne,
				UpdatedAt:                    &updatedOne,
				WeeklyCapacity:               harvest.Int(126000),
				DefaultHourlyRate:            harvest.Float64(100.0),
				CostRate:                     harvest.Float64(75.0),
				Roles:                        &roles,
				AvatarURL: harvest.String(
					"https://cache.harvestapp.com/assets/profile_images/allen_bradley_clock_tower.png?1498509661",
				),
			},
			wantErr: false,
		},
		{
			name: "error getting current user",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
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

			user, _, err := service.User.Current(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, user)
			}
		})
	}
}

func TestUserService_List(t *testing.T) {
	t.Parallel()

	createdOne := time.Date(
		2020, 5, 1, 22, 34, 41, 0, time.UTC)
	updatedOne := time.Date(
		2020, 5, 1, 22, 34, 52, 0, time.UTC)
	createdTwo := time.Date(
		2020, 5, 1, 22, 15, 45, 0, time.UTC)
	updatedTwo := time.Date(
		2020, 5, 1, 22, 32, 52, 0, time.UTC)
	createdThree := time.Date(
		2020, 5, 1, 20, 41, 0o0, 0, time.UTC)
	updatedThree := time.Date(
		2020, 5, 1, 20, 42, 25, 0, time.UTC)

	tests := []struct {
		name      string
		setupMock func(mux *http.ServeMux)
		want      *harvest.UserList
		wantErr   bool
	}{
		{
			name: "list users",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "user/list/body_1.json")
					testWriteResponse(t, w, "user/list/response_1.json")
				})
			},
			want: &harvest.UserList{
				Users: []*harvest.User{
					{
						ID:                           harvest.Int64(3230547),
						FirstName:                    harvest.String("Jim"),
						LastName:                     harvest.String("Allen"),
						Email:                        harvest.String("jimallen@example.com"),
						Telephone:                    harvest.String(""),
						Timezone:                     harvest.String("Mountain Time (US & Canada)"),
						HasAccessToAllFutureProjects: harvest.Bool(false),
						IsContractor:                 harvest.Bool(false),
						IsAdmin:                      harvest.Bool(false),
						IsProjectManager:             harvest.Bool(false),
						CanSeeRates:                  harvest.Bool(false),
						CanCreateProjects:            harvest.Bool(false),
						CanCreateInvoices:            harvest.Bool(false),
						WeeklyCapacity:               harvest.Int(126000),
						DefaultHourlyRate:            harvest.Float64(100),
						CostRate:                     harvest.Float64(50),
						Roles:                        &[]string{"Developer"},
						AvatarURL: harvest.String(
							"https://cache.harvestapp.com/assets/profile_images/abraj_albait_towers.png?1498516481",
						),
						IsActive:  harvest.Bool(true),
						CreatedAt: &createdOne,
						UpdatedAt: &updatedOne,
					},
					{
						ID:                           harvest.Int64(1782959),
						FirstName:                    harvest.String("Kim"),
						LastName:                     harvest.String("Allen"),
						Email:                        harvest.String("kimallen@example.com"),
						Telephone:                    harvest.String(""),
						Timezone:                     harvest.String("Eastern Time (US & Canada)"),
						HasAccessToAllFutureProjects: harvest.Bool(true),
						IsContractor:                 harvest.Bool(false),
						IsAdmin:                      harvest.Bool(false),
						IsProjectManager:             harvest.Bool(true),
						CanSeeRates:                  harvest.Bool(false),
						CanCreateProjects:            harvest.Bool(false),
						CanCreateInvoices:            harvest.Bool(false),
						IsActive:                     harvest.Bool(true),
						CreatedAt:                    &createdTwo,
						UpdatedAt:                    &updatedTwo,
						WeeklyCapacity:               harvest.Int(126000),
						DefaultHourlyRate:            harvest.Float64(100.0),
						CostRate:                     harvest.Float64(50.0),
						Roles:                        &[]string{"Designer"},
						AvatarURL: harvest.String(
							"https://cache.harvestapp.com/assets/profile_images/cornell_clock_tower.png?1498515345",
						),
					},
					{
						ID:                           harvest.Int64(1782884),
						FirstName:                    harvest.String("Bob"),
						LastName:                     harvest.String("Powell"),
						Email:                        harvest.String("bobpowell@example.com"),
						Telephone:                    harvest.String(""),
						Timezone:                     harvest.String("Mountain Time (US & Canada)"),
						HasAccessToAllFutureProjects: harvest.Bool(false),
						IsContractor:                 harvest.Bool(false),
						IsAdmin:                      harvest.Bool(true),
						IsProjectManager:             harvest.Bool(false),
						CanSeeRates:                  harvest.Bool(true),
						CanCreateProjects:            harvest.Bool(true),
						CanCreateInvoices:            harvest.Bool(true),
						IsActive:                     harvest.Bool(true),
						CreatedAt:                    &createdThree,
						UpdatedAt:                    &updatedThree,
						WeeklyCapacity:               harvest.Int(126000),
						DefaultHourlyRate:            harvest.Float64(100.0),
						CostRate:                     harvest.Float64(75.0),
						Roles:                        &[]string{"Founder", "CEO"},
						AvatarURL: harvest.String(
							"https://cache.harvestapp.com/assets/profile_images/allen_bradley_clock_tower.png?1498509661",
						),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(3),
					NextPage:     nil,
					PreviousPage: nil,
					Page:         harvest.Int(1),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/users?page=1&per_page=100"),
						Next:     nil,
						Previous: nil,
						Last:     harvest.String("https://api.harvestapp.com/v2/users?page=1&per_page=100"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error listing users",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
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

			userList, _, err := service.User.List(context.Background(), &harvest.UserListOptions{})

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, userList)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, userList)
			}
		})
	}
}

func TestUserService_Update(t *testing.T) {
	t.Parallel()

	createdOne := time.Date(
		2018, 1, 1, 19, 20, 46, 0, time.UTC)
	updatedOne := time.Date(
		2019, 1, 25, 19, 20, 57, 0, time.UTC)

	tests := []struct {
		name      string
		userID    int64
		data      *harvest.UserUpdateRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.User
		wantErr   bool
	}{
		{
			name:   "update user",
			userID: 3237198,
			data: &harvest.UserUpdateRequest{
				FirstName: harvest.String("Project"),
				LastName:  harvest.String("Manager"),
				Email:     harvest.String("pm@example.com"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/users/3237198", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					testFormValues(t, r, values{})
					testBody(t, r, "user/update/body_1.json")
					testWriteResponse(t, w, "user/update/response_1.json")
				})
			},
			want: &harvest.User{
				ID:                           harvest.Int64(3237198),
				DefaultHourlyRate:            harvest.Float64(120),
				IsActive:                     harvest.Bool(true),
				CreatedAt:                    &createdOne,
				UpdatedAt:                    &updatedOne,
				FirstName:                    harvest.String("Project"),
				LastName:                     harvest.String("Manager"),
				Email:                        harvest.String("pm@example.com"),
				Telephone:                    harvest.String(""),
				Timezone:                     harvest.String("Eastern Time (US & Canada)"),
				HasAccessToAllFutureProjects: harvest.Bool(true),
				IsContractor:                 harvest.Bool(false),
				IsAdmin:                      harvest.Bool(false),
				IsProjectManager:             harvest.Bool(true),
				CanSeeRates:                  harvest.Bool(true),
				CanCreateInvoices:            harvest.Bool(true),
				CanCreateProjects:            harvest.Bool(true),
				WeeklyCapacity:               harvest.Int(126000),
				CostRate:                     harvest.Float64(50),
				Roles:                        &[]string{"Project Manager"},
				AvatarURL: harvest.String(
					"https://{ACCOUNT_SUBDOMAIN}.harvestapp.com/assets/profile_images/big_ben.png?1485372046",
				),
			},
			wantErr: false,
		},
		{
			name:   "error updating user",
			userID: 3237198,
			data: &harvest.UserUpdateRequest{
				FirstName: harvest.String("Project"),
				LastName:  harvest.String("Manager"),
				Email:     harvest.String("pm@example.com"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/users/3237198", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
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

			user, _, err := service.User.Update(context.Background(), tt.userID, tt.data)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, user)
			}
		})
	}
}

func TestUser_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.User
		want string
	}{
		{
			name: "User with all fields",
			in: harvest.User{
				ID:                           harvest.Int64(1782884),
				FirstName:                    harvest.String("Bob"),
				LastName:                     harvest.String("Powell"),
				Name:                         harvest.String("Bob Powell"),
				Email:                        harvest.String("bobpowell@example.com"),
				Telephone:                    harvest.String("555-1234"),
				Timezone:                     harvest.String("Mountain Time (US & Canada)"),
				HasAccessToAllFutureProjects: harvest.Bool(false),
				IsContractor:                 harvest.Bool(false),
				IsAdmin:                      harvest.Bool(true),
				IsProjectManager:             harvest.Bool(false),
				CanSeeRates:                  harvest.Bool(true),
				CanCreateProjects:            harvest.Bool(true),
				CanCreateInvoices:            harvest.Bool(true),
				IsActive:                     harvest.Bool(true),
				WeeklyCapacity:               harvest.Int(126000),
				DefaultHourlyRate:            harvest.Float64(100.0),
				CostRate:                     harvest.Float64(75.0),
				Roles:                        &[]string{"Founder", "CEO"},
				AvatarURL:                    harvest.String("https://example.com/avatar.png"),
				CreatedAt:                    harvest.TimeTimeP(time.Date(2020, 5, 1, 20, 41, 0, 0, time.UTC)),
				UpdatedAt:                    harvest.TimeTimeP(time.Date(2020, 5, 1, 20, 42, 25, 0, time.UTC)),
			},
			want: `harvest.User{ID:1782884, FirstName:"Bob", LastName:"Powell", Name:"Bob Powell", Email:"bobpowell@example.com", Telephone:"555-1234", Timezone:"Mountain Time (US & Canada)", HasAccessToAllFutureProjects:false, IsContractor:false, IsAdmin:true, IsProjectManager:false, CanSeeRates:true, CanCreateProjects:true, CanCreateInvoices:true, IsActive:true, WeeklyCapacity:126000, DefaultHourlyRate:100, CostRate:75, Roles:["Founder" "CEO"], AvatarURL:"https://example.com/avatar.png", CreatedAt:time.Time{2020-05-01 20:41:00 +0000 UTC}, UpdatedAt:time.Time{2020-05-01 20:42:25 +0000 UTC}}`, //nolint: lll
		},
		{
			name: "User with minimal fields",
			in: harvest.User{
				ID:        harvest.Int64(999),
				FirstName: harvest.String("John"),
				LastName:  harvest.String("Doe"),
				Email:     harvest.String("john@example.com"),
			},
			want: `harvest.User{ID:999, FirstName:"John", LastName:"Doe", Email:"john@example.com"}`,
		},
		{
			name: "User with boolean flags",
			in: harvest.User{
				ID:                           harvest.Int64(3230547),
				FirstName:                    harvest.String("Jim"),
				LastName:                     harvest.String("Allen"),
				Email:                        harvest.String("jimallen@example.com"),
				HasAccessToAllFutureProjects: harvest.Bool(false),
				IsContractor:                 harvest.Bool(false),
				IsAdmin:                      harvest.Bool(false),
				IsProjectManager:             harvest.Bool(false),
				CanSeeRates:                  harvest.Bool(false),
				CanCreateProjects:            harvest.Bool(false),
				CanCreateInvoices:            harvest.Bool(false),
				IsActive:                     harvest.Bool(true),
			},
			want: `harvest.User{ID:3230547, FirstName:"Jim", LastName:"Allen", Email:"jimallen@example.com", HasAccessToAllFutureProjects:false, IsContractor:false, IsAdmin:false, IsProjectManager:false, CanSeeRates:false, CanCreateProjects:false, CanCreateInvoices:false, IsActive:true}`, //nolint: lll
		},
		{
			name: "User with roles",
			in: harvest.User{
				ID:        harvest.Int64(1782959),
				FirstName: harvest.String("Kim"),
				LastName:  harvest.String("Allen"),
				Email:     harvest.String("kimallen@example.com"),
				Roles:     &[]string{"Designer", "Developer", "Project Manager"},
			},
			want: `harvest.User{ID:1782959, FirstName:"Kim", LastName:"Allen", Email:"kimallen@example.com", Roles:["Designer" "Developer" "Project Manager"]}`, //nolint: lll
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

func TestUserList_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.UserList
		want string
	}{
		{
			name: "UserList with multiple users",
			in: harvest.UserList{
				Users: []*harvest.User{
					{
						ID:                harvest.Int64(1782884),
						FirstName:         harvest.String("Bob"),
						LastName:          harvest.String("Powell"),
						Email:             harvest.String("bobpowell@example.com"),
						IsAdmin:           harvest.Bool(true),
						IsActive:          harvest.Bool(true),
						DefaultHourlyRate: harvest.Float64(100.0),
						CostRate:          harvest.Float64(75.0),
						CreatedAt:         harvest.TimeTimeP(time.Date(2020, 5, 1, 20, 41, 0, 0, time.UTC)),
						UpdatedAt:         harvest.TimeTimeP(time.Date(2020, 5, 1, 20, 42, 25, 0, time.UTC)),
					},
					{
						ID:                harvest.Int64(3230547),
						FirstName:         harvest.String("Jim"),
						LastName:          harvest.String("Allen"),
						Email:             harvest.String("jimallen@example.com"),
						IsAdmin:           harvest.Bool(false),
						IsActive:          harvest.Bool(true),
						DefaultHourlyRate: harvest.Float64(100),
						CostRate:          harvest.Float64(50),
						CreatedAt:         harvest.TimeTimeP(time.Date(2020, 5, 1, 22, 34, 41, 0, time.UTC)),
						UpdatedAt:         harvest.TimeTimeP(time.Date(2020, 5, 1, 22, 34, 52, 0, time.UTC)),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.UserList{Users:[harvest.User{ID:1782884, FirstName:"Bob", LastName:"Powell", Email:"bobpowell@example.com", IsAdmin:true, IsActive:true, DefaultHourlyRate:100, CostRate:75, CreatedAt:time.Time{2020-05-01 20:41:00 +0000 UTC}, UpdatedAt:time.Time{2020-05-01 20:42:25 +0000 UTC}} harvest.User{ID:3230547, FirstName:"Jim", LastName:"Allen", Email:"jimallen@example.com", IsAdmin:false, IsActive:true, DefaultHourlyRate:100, CostRate:50, CreatedAt:time.Time{2020-05-01 22:34:41 +0000 UTC}, UpdatedAt:time.Time{2020-05-01 22:34:52 +0000 UTC}}], Pagination:harvest.Pagination{PerPage:100, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
		},
		{
			name: "UserList with single user",
			in: harvest.UserList{
				Users: []*harvest.User{
					{
						ID:        harvest.Int64(999),
						FirstName: harvest.String("John"),
						LastName:  harvest.String("Doe"),
						Email:     harvest.String("john@example.com"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.UserList{Users:[harvest.User{ID:999, FirstName:"John", LastName:"Doe", Email:"john@example.com"}], Pagination:harvest.Pagination{PerPage:100, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
		},
		{
			name: "Empty UserList",
			in: harvest.UserList{
				Users: []*harvest.User{},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(0),
					TotalEntries: harvest.Int(0),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.UserList{Users:[], Pagination:harvest.Pagination{PerPage:100, TotalPages:0, TotalEntries:0, Page:1}}`, //nolint: lll
		},
		{
			name: "UserList with Links",
			in: harvest.UserList{
				Users: []*harvest.User{
					{
						ID:        harvest.Int64(100),
						FirstName: harvest.String("Alice"),
						LastName:  harvest.String("Smith"),
						Email:     harvest.String("alice@example.com"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(50),
					TotalPages:   harvest.Int(3),
					TotalEntries: harvest.Int(150),
					Page:         harvest.Int(2),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/users?page=1&per_page=50"),
						Next:     harvest.String("https://api.harvestapp.com/v2/users?page=3&per_page=50"),
						Previous: harvest.String("https://api.harvestapp.com/v2/users?page=1&per_page=50"),
						Last:     harvest.String("https://api.harvestapp.com/v2/users?page=3&per_page=50"),
					},
				},
			},
			want: `harvest.UserList{Users:[harvest.User{ID:100, FirstName:"Alice", LastName:"Smith", Email:"alice@example.com"}], Pagination:harvest.Pagination{PerPage:50, TotalPages:3, TotalEntries:150, Page:2, Links:harvest.PageLinks{First:"https://api.harvestapp.com/v2/users?page=1&per_page=50", Next:"https://api.harvestapp.com/v2/users?page=3&per_page=50", Previous:"https://api.harvestapp.com/v2/users?page=1&per_page=50", Last:"https://api.harvestapp.com/v2/users?page=3&per_page=50"}}}`, //nolint: lll
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
