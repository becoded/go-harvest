package harvest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserService_Create(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	createdOne := time.Date(
		2020, 1, 25, 19, 20, 46, 0, time.UTC)
	updatedOne := time.Date(
		2020, 1, 25, 19, 20, 57, 0, time.UTC)

	roles := []string{"Project Manager"}

	tests := []struct {
		name       string
		args       *UserCreateRequest
		method     string
		path       string
		formValues values
		body       string
		response   string
		want       *User
		wantErr    error
	}{
		{
			name: "create user",
			args: &UserCreateRequest{
				FirstName:        String("George"),
				LastName:         String("Frank"),
				Email:            String("george@example.com"),
				IsProjectManager: Bool(true),
			},
			method:     "POST",
			path:       "/users",
			formValues: values{},
			body:       "user/create/body_1.json",
			response:   "user/create/response_1.json",
			want: &User{
				Id:                           Int64(3),
				FirstName:                    String("Gary"),
				LastName:                     String("Frank"),
				Email:                        String("george@example.com"),
				Timezone:                     String("Eastern Time (US & Canada)"),
				HasAccessToAllFutureProjects: Bool(false),
				IsContractor:                 Bool(false),
				IsAdmin:                      Bool(false),
				IsActive:                     Bool(true),
				CanSeeRates:                  Bool(false),
				CanCreateProjects:            Bool(false),
				CanCreateInvoices:            Bool(false),
				Telephone:                    String(""),
				IsProjectManager:             Bool(true),
				CreatedAt:                    &createdOne,
				UpdatedAt:                    &updatedOne,
				WeeklyCapacity:               Int(126000),
				DefaultHourlyRate:            Float64(0),
				CostRate:                     Float64(0),
				Roles:                        &roles,
				AvatarUrl:                    String("https://{ACCOUNT_SUBDOMAIN}.harvestapp.com/assets/profile_images/big_ben.png?1485372046"),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.method)
				testFormValues(t, r, tt.formValues)
				testBody(t, r, tt.body)
				response, err := os.ReadFile(filepath.Join("..", "testdata", tt.response))
				assert.NoError(t, err)
				_, err = fmt.Fprint(w, string(response))
				assert.NoError(t, err)
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
	service, mux, _, teardown := setup()
	defer teardown()

	type args struct {
		userId int64
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
				userId: 1,
			},
			method:     "DELETE",
			path:       "/users/1",
			formValues: values{},
			body:       "user/delete/body_1.json",
			response:   "user/delete/response_1.json",
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.method)
				testFormValues(t, r, tt.formValues)
				testBody(t, r, tt.body)
				response, err := os.ReadFile(filepath.Join("..", "testdata", tt.response))
				assert.NoError(t, err)
				_, err = fmt.Fprint(w, string(response))
				assert.NoError(t, err)
			})

			_, err := service.User.Delete(context.Background(), tt.args.userId)

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
	service, mux, _, teardown := setup()
	defer teardown()

	type args struct {
		userId int64
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
		want       *User
		wantErr    error
	}{
		{
			name: "get 3230547",
			args: args{
				userId: 3230547,
			},
			method:     "GET",
			path:       "/users/3230547",
			formValues: values{},
			body:       "user/get/body_1.json",
			response:   "user/get/response_1.json",
			want: &User{
				Id:                           Int64(3230547),
				IsActive:                     Bool(true),
				FirstName:                    String("Jim"),
				LastName:                     String("Allen"),
				Email:                        String("jimallen@example.com"),
				Telephone:                    String(""),
				Timezone:                     String("Mountain Time (US & Canada)"),
				HasAccessToAllFutureProjects: Bool(false),
				IsContractor:                 Bool(false),
				IsAdmin:                      Bool(false),
				IsProjectManager:             Bool(false),
				CanSeeRates:                  Bool(false),
				CanCreateInvoices:            Bool(false),
				CanCreateProjects:            Bool(false),
				Roles:                        &roles,
				CostRate:                     Float64(50),
				DefaultHourlyRate:            Float64(100),
				WeeklyCapacity:               Int(126000),
				AvatarUrl:                    String("https://cache.harvestapp.com/assets/profile_images/abraj_albait_towers.png?1498516481"),

				CreatedAt: &createdOne,
				UpdatedAt: &updatedOne,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.method)
				testFormValues(t, r, tt.formValues)
				testBody(t, r, tt.body)
				response, err := os.ReadFile(filepath.Join("..", "testdata", tt.response))
				assert.NoError(t, err)
				_, err = fmt.Fprint(w, string(response))
				assert.NoError(t, err)
			})

			user, _, err := service.User.Get(context.Background(), tt.args.userId)

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
	service, mux, _, teardown := setup()
	defer teardown()

	type args struct {
		userId int64
	}

	createdOne := time.Date(
		2020, 5, 1, 20, 41, 00, 0, time.UTC)
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
		want       *User
		wantErr    error
	}{
		{
			name: "get 1",
			args: args{
				userId: 1,
			},
			method:     "GET",
			path:       "/users/me",
			formValues: values{},
			body:       "user/current/body_1.json",
			response:   "user/current/response_1.json",
			want: &User{
				Id:                           Int64(1782884),
				FirstName:                    String("Bob"),
				LastName:                     String("Powell"),
				Email:                        String("bobpowell@example.com"),
				Telephone:                    String(""),
				Timezone:                     String("Mountain Time (US & Canada)"),
				HasAccessToAllFutureProjects: Bool(false),
				IsContractor:                 Bool(false),
				IsAdmin:                      Bool(true),
				IsProjectManager:             Bool(false),
				CanSeeRates:                  Bool(true),
				CanCreateProjects:            Bool(true),
				CanCreateInvoices:            Bool(true),
				IsActive:                     Bool(true),
				CreatedAt:                    &createdOne,
				UpdatedAt:                    &updatedOne,
				WeeklyCapacity:               Int(126000),
				DefaultHourlyRate:            Float64(100.0),
				CostRate:                     Float64(75.0),
				Roles:                        &roles,
				AvatarUrl:                    String("https://cache.harvestapp.com/assets/profile_images/allen_bradley_clock_tower.png?1498509661"),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.method)
				testFormValues(t, r, tt.formValues)
				testBody(t, r, tt.body)
				response, err := os.ReadFile(filepath.Join("..", "testdata", tt.response))
				assert.NoError(t, err)
				_, err = fmt.Fprint(w, string(response))
				assert.NoError(t, err)
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
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `{"users":[{"id":3230547,"first_name":"Jim","last_name":"Allen","email":"jimallen@example.com","telephone":"","timezone":"Mountain Time (US & Canada)","has_access_to_all_future_projects":false,"is_contractor":false,"is_admin":false,"is_project_manager":false,"can_see_rates":false,"can_create_projects":false,"can_create_invoices":false,"is_active":true,"created_at":"2020-05-01T22:34:41Z","updated_at":"2020-05-01T22:34:52Z","weekly_capacity":126000,"default_hourly_rate":100,"cost_rate":50,"roles":["Developer"],"avatar_url":"https://cache.harvestapp.com/assets/profile_images/abraj_albait_towers.png?1498516481"},{"id":1782959,"first_name":"Kim","last_name":"Allen","email":"kimallen@example.com","telephone":"","timezone":"Eastern Time (US & Canada)","has_access_to_all_future_projects":true,"is_contractor":false,"is_admin":false,"is_project_manager":true,"can_see_rates":false,"can_create_projects":false,"can_create_invoices":false,"is_active":true,"created_at":"2020-05-01T22:15:45Z","updated_at":"2020-05-01T22:32:52Z","weekly_capacity":126000,"default_hourly_rate":100,"cost_rate":50,"roles":["Designer"],"avatar_url":"https://cache.harvestapp.com/assets/profile_images/cornell_clock_tower.png?1498515345"},{"id":1782884,"first_name":"Bob","last_name":"Powell","email":"bobpowell@example.com","telephone":"","timezone":"Mountain Time (US & Canada)","has_access_to_all_future_projects":false,"is_contractor":false,"is_admin":true,"is_project_manager":false,"can_see_rates":true,"can_create_projects":true,"can_create_invoices":true,"is_active":true,"created_at":"2020-05-01T20:41:00Z","updated_at":"2020-05-01T20:42:25Z","weekly_capacity":126000,"default_hourly_rate":100,"cost_rate":75,"roles":["Founder","CEO"],"avatar_url":"https://cache.harvestapp.com/assets/profile_images/allen_bradley_clock_tower.png?1498509661"}],"per_page":100,"total_pages":1,"total_entries":3,"next_page":null,"previous_page":null,"page":1,"links":{"first":"https://api.harvestapp.com/v2/users?page=1&per_page=100","next":null,"previous":null,"last":"https://api.harvestapp.com/v2/users?page=1&per_page=100"}}`)
	})

	userList, _, err := service.User.List(context.Background(), &UserListOptions{})
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
		2020, 5, 1, 20, 41, 00, 0, time.UTC)
	updatedThree := time.Date(
		2020, 5, 1, 20, 42, 25, 0, time.UTC)

	want := &UserList{
		Users: []*User{
			{
				Id:                           Int64(3230547),
				FirstName:                    String("Jim"),
				LastName:                     String("Allen"),
				Email:                        String("jimallen@example.com"),
				Telephone:                    String(""),
				Timezone:                     String("Mountain Time (US & Canada)"),
				HasAccessToAllFutureProjects: Bool(false),
				IsContractor:                 Bool(false),
				IsAdmin:                      Bool(false),
				IsProjectManager:             Bool(false),
				CanSeeRates:                  Bool(false),
				CanCreateProjects:            Bool(false),
				CanCreateInvoices:            Bool(false),
				WeeklyCapacity:               Int(126000),
				DefaultHourlyRate:            Float64(100),
				CostRate:                     Float64(50),
				Roles:                        &[]string{"Developer"},
				AvatarUrl:                    String("https://cache.harvestapp.com/assets/profile_images/abraj_albait_towers.png?1498516481"),
				IsActive:                     Bool(true),
				CreatedAt:                    &createdOne,
				UpdatedAt:                    &updatedOne,
			}, {
				Id:                           Int64(1782959),
				FirstName:                    String("Kim"),
				LastName:                     String("Allen"),
				Email:                        String("kimallen@example.com"),
				Telephone:                    String(""),
				Timezone:                     String("Eastern Time (US & Canada)"),
				HasAccessToAllFutureProjects: Bool(true),
				IsContractor:                 Bool(false),
				IsAdmin:                      Bool(false),
				IsProjectManager:             Bool(true),
				CanSeeRates:                  Bool(false),
				CanCreateProjects:            Bool(false),
				CanCreateInvoices:            Bool(false),
				IsActive:                     Bool(true),
				CreatedAt:                    &createdTwo,
				UpdatedAt:                    &updatedTwo,
				WeeklyCapacity:               Int(126000),
				DefaultHourlyRate:            Float64(100.0),
				CostRate:                     Float64(50.0),
				Roles:                        &[]string{"Designer"},
				AvatarUrl:                    String("https://cache.harvestapp.com/assets/profile_images/cornell_clock_tower.png?1498515345"),
			},
			{
				Id:                           Int64(1782884),
				FirstName:                    String("Bob"),
				LastName:                     String("Powell"),
				Email:                        String("bobpowell@example.com"),
				Telephone:                    String(""),
				Timezone:                     String("Mountain Time (US & Canada)"),
				HasAccessToAllFutureProjects: Bool(false),
				IsContractor:                 Bool(false),
				IsAdmin:                      Bool(true),
				IsProjectManager:             Bool(false),
				CanSeeRates:                  Bool(true),
				CanCreateProjects:            Bool(true),
				CanCreateInvoices:            Bool(true),
				IsActive:                     Bool(true),
				CreatedAt:                    &createdThree,
				UpdatedAt:                    &updatedThree,
				WeeklyCapacity:               Int(126000),
				DefaultHourlyRate:            Float64(100.0),
				CostRate:                     Float64(75.0),
				Roles:                        &[]string{"Founder", "CEO"},
				AvatarUrl:                    String("https://cache.harvestapp.com/assets/profile_images/allen_bradley_clock_tower.png?1498509661"),
			}},
		Pagination: Pagination{
			PerPage:      Int(100),
			TotalPages:   Int(1),
			TotalEntries: Int(3),
			NextPage:     nil,
			PreviousPage: nil,
			Page:         Int(1),
			Links: &PageLinks{
				First:    String("https://api.harvestapp.com/v2/users?page=1&per_page=100"),
				Next:     nil,
				Previous: nil,
				Last:     String("https://api.harvestapp.com/v2/users?page=1&per_page=100"),
			},
		},
	}

	assert.Equal(t, want, userList)
}

func TestUserService_Update(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/3237198", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		testBody(t, r, "user/update/body_1.json")
		response, err := os.ReadFile(filepath.Join("..", "testdata", "user/update/response_1.json"))
		assert.NoError(t, err)
		_, err = fmt.Fprint(w, string(response))
		assert.NoError(t, err)
	})

	user, _, err := service.User.Update(context.Background(), 3237198, &UserUpdateRequest{
		FirstName: String("Project"),
		LastName:  String("Manager"),
		Email:     String("pm@example.com"),
	})
	assert.NoError(t, err)

	createdOne := time.Date(
		2018, 1, 1, 19, 20, 46, 0, time.UTC)
	updatedOne := time.Date(
		2019, 1, 25, 19, 20, 57, 0, time.UTC)

	want := &User{
		Id:                           Int64(3237198),
		DefaultHourlyRate:            Float64(120),
		IsActive:                     Bool(true),
		CreatedAt:                    &createdOne,
		UpdatedAt:                    &updatedOne,
		FirstName:                    String("Project"),
		LastName:                     String("Manager"),
		Email:                        String("pm@example.com"),
		Telephone:                    String(""),
		Timezone:                     String("Eastern Time (US & Canada)"),
		HasAccessToAllFutureProjects: Bool(true),
		IsContractor:                 Bool(false),
		IsAdmin:                      Bool(false),
		IsProjectManager:             Bool(true),
		CanSeeRates:                  Bool(true),
		CanCreateInvoices:            Bool(true),
		CanCreateProjects:            Bool(true),
		WeeklyCapacity:               Int(126000),
		CostRate:                     Float64(50),
		Roles:                        &[]string{"Project Manager"},
		AvatarUrl:                    String("https://{ACCOUNT_SUBDOMAIN}.harvestapp.com/assets/profile_images/big_ben.png?1485372046"),
	}

	assert.Equal(t, want, user)
}
