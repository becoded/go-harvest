package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// UserService handles communication with the user related
// methods of the Harvest API.
//
// Harvest API docs: https://help.getharvest.com/api-v2/users-api/users/users/
type UserService service

type User struct {
	Id                           *int64     `json:"id,omitempty"`                                // Unique ID for the user.
	FirstName                    *string    `json:"first_name,omitempty"`                        // The first name of the user.
	LastName                     *string    `json:"last_name,omitempty"`                         // The last name of the user.
	Email                        *string    `json:"email,omitempty"`                             // The email address of the user.
	Telephone                    *string    `json:"telephone,omitempty"`                         // The telephone number for the user.
	Timezone                     *string    `json:"timezone,omitempty"`                          // The user's timezone.
	HasAccessToAllFutureProjects *bool      `json:"has_access_to_all_future_projects,omitempty"` // Whether the user should be automatically added to future projects.
	IsContractor                 *bool      `json:"is_contractor,omitempty"`                     // Whether the user is a contractor or an employee.
	IsAdmin                      *bool      `json:"is_admin,omitempty"`                          // Whether the user has admin permissions.
	IsProjectManager             *bool      `json:"is_project_manager,omitempty"`                // Whether the user has project manager permissions.
	CanSeeRates                  *bool      `json:"can_see_rates,omitempty"`                     // Whether the user can see billable rates on projects. Only applicable to project managers.
	CanCreateProjects            *bool      `json:"can_create_projects,omitempty"`               // Whether the user can create projects. Only applicable to project managers.
	CanCreateInvoices            *bool      `json:"can_create_invoices,omitempty"`               // Whether the user can create invoices. Only applicable to project managers.
	IsActive                     *bool      `json:"is_active,omitempty"`                         // Whether the user is active or archived.
	WeeklyCapacity               *int       `json:"weekly_capacity,omitempty"`                   // The number of hours per week this person is available to work in seconds. For example, if a person's capacity is 35 hours, the API will return 126000 seconds.
	DefaultHourlyRate            *float64   `json:"default_hourly_rate,omitempty"`               // The billable rate to use for this user when they are added to a project.
	CostRate                     *float64   `json:"cost_rate,omitempty"`                         // The cost rate to use for this user when calculating a project's costs vs billable amount.
	Roles                        *[]string  `json:"roles,omitempty"`                             // of strings	The role names assigned to this person.
	AvatarUrl                    *string    `json:"avatar_url,omitempty"`                        // The URL to the user's avatar image.
	CreatedAt                    *time.Time `json:"created_at,omitempty"`                        // Date and time the user was created.
	UpdatedAt                    *time.Time `json:"updated_at,omitempty"`                        // Date and time the user was last updated.
}

type UserList struct {
	Users []*User `json:"users"`

	Pagination
}

func (c User) String() string {
	return Stringify(c)
}

func (c UserList) String() string {
	return Stringify(c)
}

type UserCreateRequest struct {
	FirstName                    *string   `json:"first_name"`                                  // required	The first name of the user.
	LastName                     *string   `json:"last_name"`                                   // required	The last name of the user.
	Email                        *string   `json:"email"`                                       // required	The email address of the user.
	Telephone                    *string   `json:"telephone,omitempty"`                         // optional	The telephone number for the user.
	Timezone                     *string   `json:"timezone,omitempty"`                          // optional	The user's timezone. Defaults to the company's timezone. See a list of supported time zones.
	HasAccessToAllFutureProjects *bool     `json:"has_access_to_all_future_projects,omitempty"` // optional	Whether the user should be automatically added to future projects. Defaults to false.
	IsContractor                 *bool     `json:"is_contractor,omitempty"`                     // optional	Whether the user is a contractor or an employee. Defaults to false.
	IsAdmin                      *bool     `json:"is_admin,omitempty"`                          // optional	Whether the user has admin permissions. Defaults to false.
	IsProjectManager             *bool     `json:"is_project_manager,omitempty"`                // optional	Whether the user has project manager permissions. Defaults to false.
	CanSeeRates                  *bool     `json:"can_see_rates,omitempty"`                     // optional	Whether the user can see billable rates on projects. Only applicable to project managers. Defaults to false.
	CanCreateProjects            *bool     `json:"can_create_projects,omitempty"`               // optional	Whether the user can create projects. Only applicable to project managers. Defaults to false.
	CanCreateInvoices            *bool     `json:"can_create_invoices,omitempty"`               // optional	Whether the user can create invoices. Only applicable to project managers. Defaults to false.
	IsActive                     *bool     `json:"is_active,omitempty"`                         // optional	Whether the user is active or archived. Defaults to true.
	WeeklyCapacity               *int      `json:"weekly_capacity,omitempty"`                   // optional	The number of hours per week this person is available to work in seconds. Defaults to 126000 seconds (35 hours).
	DefaultHourlyRate            *float64  `json:"default_hourly_rate,omitempty"`               // optional	The billable rate to use for this user when they are added to a project. Defaults to 0.
	CostRate                     *float64  `json:"cost_rate,omitempty"`                         // optional	The cost rate to use for this user when calculating a project's costs vs billable amount. Defaults to 0.
	Roles                        []*string `json:"roles,omitempty"`                             // of strings	optional	The role names assigned to this person.
}

type UserUpdateRequest struct {
	FirstName                    *string   `json:"first_name,omitempty"`                        // The first name of the user. Can't be updated if the user is inactive.
	LastName                     *string   `json:"last_name,omitempty"`                         // The last name of the user. Can't be updated if the user is inactive.
	Email                        *string   `json:"email,omitempty"`                             // The email address of the user. Can't be updated if the user is inactive.
	Telephone                    *string   `json:"telephone,omitempty"`                         // The telephone number for the user.
	Timezone                     *string   `json:"timezone,omitempty"`                          // The user's timezone. Defaults to the company's timezone. See a list of supported time zones.
	HasAccessToAllFutureProjects *bool     `json:"has_access_to_all_future_projects,omitempty"` // Whether the user should be automatically added to future projects.
	IsContractor                 *bool     `json:"is_contractor,omitempty"`                     // Whether the user is a contractor or an employee.
	IsAdmin                      *bool     `json:"is_admin,omitempty"`                          // Whether the user has admin permissions.
	IsProjectManager             *bool     `json:"is_project_manager,omitempty"`                // Whether the user has project manager permissions.
	CanSeeRates                  *bool     `json:"can_see_rates,omitempty"`                     // Whether the user can see billable rates on projects. Only applicable to project managers.
	CanCreateProjects            *bool     `json:"can_create_projects,omitempty"`               // Whether the user can create projects. Only applicable to project managers.
	CanCreateInvoices            *bool     `json:"can_create_invoices,omitempty"`               // Whether the user can create invoices. Only applicable to project managers.
	IsActive                     *bool     `json:"is_active,omitempty"`                         // Whether the user is active or archived.
	WeeklyCapacity               *int      `json:"weekly_capacity,omitempty"`                   // The number of hours per week this person is available to work in seconds.
	DefaultHourlyRate            *float64  `json:"default_hourly_rate,omitempty"`               // The billable rate to use for this user when they are added to a project.
	CostRate                     *float64  `json:"cost_rate,omitempty"`                         // The cost rate to use for this user when calculating a project's costs vs billable amount.
	Roles                        []*string `json:"roles,omitempty"`                             // of strings	The role names assigned to this person.
}

type UserListOptions struct {
	// Pass true to only return active projects and false to return inactive projects.
	IsActive bool `url:"is_active,omitempty"`
	// Only return projects that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

func (s *UserService) ListUsers(ctx context.Context, opt *UserListOptions) (*UserList, *http.Response, error) {
	u := "users"
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	userList := new(UserList)
	resp, err := s.client.Do(ctx, req, &userList)
	if err != nil {
		return nil, resp, err
	}

	return userList, resp, nil
}

func (s *UserService) GetUser(ctx context.Context, userId int64) (*User, *http.Response, error) {
	u := fmt.Sprintf("users/%d", userId)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(ctx, req, &user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, nil
}

func (s *UserService) Current(ctx context.Context, userId int64) (*User, *http.Response, error) {
	u := "users/me"

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(ctx, req, &user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, nil
}

func (s *UserService) CreateUser(ctx context.Context, data *UserCreateRequest) (*User, *http.Response, error) {
	u := "users"

	req, err := s.client.NewRequest("POST", u, data)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(ctx, req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userId int64, data *UserUpdateRequest) (*User, *http.Response, error) {
	u := fmt.Sprintf("users/%d", userId)

	req, err := s.client.NewRequest("PATCH", u, data)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(ctx, req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userId int64) (*http.Response, error) {
	u := fmt.Sprintf("users/%d", userId)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
