package harvest_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/becoded/go-harvest/harvest"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Create(t *testing.T) {
	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	createdOne := time.Date(
		2020, 1, 25, 19, 20, 46, 0, time.UTC)
	updatedOne := time.Date(
		2020, 1, 25, 19, 20, 57, 0, time.UTC)

	roles := []string{"Project Manager"}

	tests := []struct {
		name       string
		args       *harvest.UserCreateRequest
		method     string
		path       string
		formValues values
		body       string
		response   string
		want       *harvest.User
		wantErr    error
	}{
		{
			name: "create user",
			args: &harvest.UserCreateRequest{
				FirstName:        harvest.String("George"),
				LastName:         harvest.String("Frank"),
				Email:            harvest.String("george@example.com"),
				IsProjectManager: harvest.Bool(true),
			},
			method:     "POST",
			path:       "/users",
			formValues: values{},
			body:       "user/create/body_1.json",
			response:   "user/create/response_1.json",
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
				AvatarURL:                    harvest.String("https://{ACCOUNT_SUBDOMAIN}.harvestapp.com/assets/profile_images/big_ben.png?1485372046"),
			},
			wantErr: nil,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.method)
				testFormValues(t, r, tt.formValues)
				testBody(t, r, tt.body)
				testWriteResponse(t, w, tt.response)
			})

			user, _, err := service.User.Create(context.Background(), tt.args)

			if tt.wantErr != nil {
				assert.Nil(t, user)
				assert.EqualError(t, err, tt.wantErr.Error())

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, user)
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	type args struct {
		userID int64
	}

	tests := []struct {
		name       string
		args       args
		method     string
		path       string
		formValues values
		body       string
		response   string
		wantErr    error
	}{
		{
			name: "delete 1",
			args: args{
				userID: 1,
			},
			method:     "DELETE",
			path:       "/users/1",
			formValues: values{},
			body:       "user/delete/body_1.json",
			response:   "user/delete/response_1.json",
			wantErr:    nil,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.method)
				testFormValues(t, r, tt.formValues)
				testBody(t, r, tt.body)
				testWriteResponse(t, w, tt.response)
			})

			_, err := service.User.Delete(context.Background(), tt.args.userID)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())

				return
			}

			assert.NoError(t, err)
			assert.Nil(t, err)
		})
	}
}

func TestUserService_Get(t *testing.T) {
	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	type args struct {
		userID int64
	}

	createdOne := time.Date(
		2020, 5, 1, 22, 34, 41, 0, time.UTC)
	updatedOne := time.Date(
		2020, 5, 1, 22, 34, 52, 0, time.UTC)

	roles := []string{"Developer"}

	tests := []struct {
		name       string
		args       args
		method     string
		path       string
		formValues values
		body       string
		response   string
		want       *harvest.User
		wantErr    error
	}{
		{
			name: "get 3230547",
			args: args{
				userID: 3230547,
			},
			method:     "GET",
			path:       "/users/3230547",
			formValues: values{},
			body:       "user/get/body_1.json",
			response:   "user/get/response_1.json",
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
				AvatarURL:                    harvest.String("https://cache.harvestapp.com/assets/profile_images/abraj_albait_towers.png?1498516481"),

				CreatedAt: &createdOne,
				UpdatedAt: &updatedOne,
			},
			wantErr: nil,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.method)
				testFormValues(t, r, tt.formValues)
				testBody(t, r, tt.body)
				testWriteResponse(t, w, tt.response)
			})

			user, _, err := service.User.Get(context.Background(), tt.args.userID)

			if tt.wantErr != nil {
				assert.Nil(t, user)
				assert.EqualError(t, err, tt.wantErr.Error())

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, user)
		})
	}
}

func TestUserService_Current(t *testing.T) {
	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	type args struct {
		userID int64
	}

	createdOne := time.Date(
		2020, 5, 1, 20, 41, 0o0, 0, time.UTC)
	updatedOne := time.Date(
		2020, 5, 1, 20, 42, 25, 0, time.UTC)

	roles := []string{"Founder", "CEO"}

	tests := []struct {
		name       string
		args       args
		method     string
		path       string
		formValues values
		body       string
		response   string
		want       *harvest.User
		wantErr    error
	}{
		{
			name: "get 1",
			args: args{
				userID: 1,
			},
			method:     "GET",
			path:       "/users/me",
			formValues: values{},
			body:       "user/current/body_1.json",
			response:   "user/current/response_1.json",
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
				AvatarURL:                    harvest.String("https://cache.harvestapp.com/assets/profile_images/allen_bradley_clock_tower.png?1498509661"),
			},
			wantErr: nil,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.method)
				testFormValues(t, r, tt.formValues)
				testBody(t, r, tt.body)
				testWriteResponse(t, w, tt.response)
			})

			user, _, err := service.User.Current(context.Background())

			if tt.wantErr != nil {
				assert.Nil(t, user)
				assert.EqualError(t, err, tt.wantErr.Error())

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, user)
		})
	}
}

func TestUserService_List(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `{"users":[{"id":3230547,"first_name":"Jim","last_name":"Allen","email":"jimallen@example.com","telephone":"","timezone":"Mountain Time (US & Canada)","has_access_to_all_future_projects":false,"is_contractor":false,"is_admin":false,"is_project_manager":false,"can_see_rates":false,"can_create_projects":false,"can_create_invoices":false,"is_active":true,"created_at":"2020-05-01T22:34:41Z","updated_at":"2020-05-01T22:34:52Z","weekly_capacity":126000,"default_hourly_rate":100,"cost_rate":50,"roles":["Developer"],"avatar_url":"https://cache.harvestapp.com/assets/profile_images/abraj_albait_towers.png?1498516481"},{"id":1782959,"first_name":"Kim","last_name":"Allen","email":"kimallen@example.com","telephone":"","timezone":"Eastern Time (US & Canada)","has_access_to_all_future_projects":true,"is_contractor":false,"is_admin":false,"is_project_manager":true,"can_see_rates":false,"can_create_projects":false,"can_create_invoices":false,"is_active":true,"created_at":"2020-05-01T22:15:45Z","updated_at":"2020-05-01T22:32:52Z","weekly_capacity":126000,"default_hourly_rate":100,"cost_rate":50,"roles":["Designer"],"avatar_url":"https://cache.harvestapp.com/assets/profile_images/cornell_clock_tower.png?1498515345"},{"id":1782884,"first_name":"Bob","last_name":"Powell","email":"bobpowell@example.com","telephone":"","timezone":"Mountain Time (US & Canada)","has_access_to_all_future_projects":false,"is_contractor":false,"is_admin":true,"is_project_manager":false,"can_see_rates":true,"can_create_projects":true,"can_create_invoices":true,"is_active":true,"created_at":"2020-05-01T20:41:00Z","updated_at":"2020-05-01T20:42:25Z","weekly_capacity":126000,"default_hourly_rate":100,"cost_rate":75,"roles":["Founder","CEO"],"avatar_url":"https://cache.harvestapp.com/assets/profile_images/allen_bradley_clock_tower.png?1498509661"}],"per_page":100,"total_pages":1,"total_entries":3,"next_page":null,"previous_page":null,"page":1,"links":{"first":"https://api.harvestapp.com/v2/users?page=1&per_page=100","next":null,"previous":null,"last":"https://api.harvestapp.com/v2/users?page=1&per_page=100"}}`)
	})

	userList, _, err := service.User.List(context.Background(), &harvest.UserListOptions{})
	assert.NoError(t, err)

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

	want := &harvest.UserList{
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
				AvatarURL:                    harvest.String("https://cache.harvestapp.com/assets/profile_images/abraj_albait_towers.png?1498516481"),
				IsActive:                     harvest.Bool(true),
				CreatedAt:                    &createdOne,
				UpdatedAt:                    &updatedOne,
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
				AvatarURL:                    harvest.String("https://cache.harvestapp.com/assets/profile_images/cornell_clock_tower.png?1498515345"),
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
				AvatarURL:                    harvest.String("https://cache.harvestapp.com/assets/profile_images/allen_bradley_clock_tower.png?1498509661"),
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
	}

	assert.Equal(t, want, userList)
}

func TestUserService_Update(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/users/3237198", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		testBody(t, r, "user/update/body_1.json")
		testWriteResponse(t, w, "user/update/response_1.json")
	})

	user, _, err := service.User.Update(context.Background(), 3237198, &harvest.UserUpdateRequest{
		FirstName: harvest.String("Project"),
		LastName:  harvest.String("Manager"),
		Email:     harvest.String("pm@example.com"),
	})
	assert.NoError(t, err)

	createdOne := time.Date(
		2018, 1, 1, 19, 20, 46, 0, time.UTC)
	updatedOne := time.Date(
		2019, 1, 25, 19, 20, 57, 0, time.UTC)

	want := &harvest.User{
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
		AvatarURL:                    harvest.String("https://{ACCOUNT_SUBDOMAIN}.harvestapp.com/assets/profile_images/big_ben.png?1485372046"),
	}

	assert.Equal(t, want, user)
}
