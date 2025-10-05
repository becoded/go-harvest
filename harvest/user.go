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
	// Unique ID for the user.
	ID *int64 `json:"id,omitempty"`
	// The first name of the user.
	FirstName *string `json:"first_name,omitempty"`
	// The last name of the user.
	LastName *string `json:"last_name,omitempty"`
	// The full name of the user - this is populated when listing time entries
	Name *string `json:"name,omitempty"`
	// The email address of the user.
	Email *string `json:"email,omitempty"`
	// The telephone number for the user.
	Telephone *string `json:"telephone,omitempty"`
	// The user's timezone.
	Timezone *string `json:"timezone,omitempty"`
	// Whether the user should be automatically added to future projects.
	HasAccessToAllFutureProjects *bool `json:"has_access_to_all_future_projects,omitempty"`
	// Whether the user is a contractor or an employee.
	IsContractor *bool `json:"is_contractor,omitempty"`
	// Whether the user has admin permissions.
	IsAdmin *bool `json:"is_admin,omitempty"`
	// Whether the user has project manager permissions.
	IsProjectManager *bool `json:"is_project_manager,omitempty"`
	// Whether the user can see billable rates on projects. Only applicable to project managers.
	CanSeeRates *bool `json:"can_see_rates,omitempty"`
	// Whether the user can create projects. Only applicable to project managers.
	CanCreateProjects *bool `json:"can_create_projects,omitempty"`
	// Whether the user can create invoices. Only applicable to project managers.
	CanCreateInvoices *bool `json:"can_create_invoices,omitempty"`
	// Whether the user is active or archived.
	IsActive *bool `json:"is_active,omitempty"`
	// The number of hours per week this person is available to work in seconds.
	// For example, if a person's capacity is 35 hours, the API will return 126000 seconds.
	WeeklyCapacity *int `json:"weekly_capacity,omitempty"`
	// The billable rate to use for this user when they are added to a project.
	DefaultHourlyRate *float64 `json:"default_hourly_rate,omitempty"`
	// The cost rate to use for this user when calculating a project's costs vs billable amount.
	CostRate *float64 `json:"cost_rate,omitempty"`
	// of strings	The role names assigned to this person.
	Roles *[]string `json:"roles,omitempty"`
	// The URL to the user's avatar image.
	AvatarURL *string `json:"avatar_url,omitempty"`
	// Date and time the user was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Date and time the user was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
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
	// required	The first name of the user.
	FirstName *string `json:"first_name"`
	// required	The last name of the user.
	LastName *string `json:"last_name"`
	// required	The email address of the user.
	Email *string `json:"email"`
	// optional	The telephone number for the user.
	Telephone *string `json:"telephone,omitempty"`
	// optional	The user's timezone. Defaults to the company's timezone. See a list of supported time zones.
	Timezone *string `json:"timezone,omitempty"`
	// optional	Whether the user should be automatically added to future projects. Defaults to false.
	HasAccessToAllFutureProjects *bool `json:"has_access_to_all_future_projects,omitempty"`
	// optional	Whether the user is a contractor or an employee. Defaults to false.
	IsContractor *bool `json:"is_contractor,omitempty"`
	// optional	Whether the user has admin permissions. Defaults to false.
	IsAdmin *bool `json:"is_admin,omitempty"`
	// optional	Whether the user has project manager permissions. Defaults to false.
	IsProjectManager *bool `json:"is_project_manager,omitempty"`
	// optional	Whether the user can see billable rates on projects.
	// Only applicable to project managers. Defaults to false.
	CanSeeRates *bool `json:"can_see_rates,omitempty"`
	// optional	Whether the user can create projects. Only applicable to project managers. Defaults to false.
	CanCreateProjects *bool `json:"can_create_projects,omitempty"`
	// optional	Whether the user can create invoices. Only applicable to project managers. Defaults to false.
	CanCreateInvoices *bool `json:"can_create_invoices,omitempty"`
	// optional	Whether the user is active or archived. Defaults to true.
	IsActive *bool `json:"is_active,omitempty"`
	// optional	The number of hours per week this person is available to work in seconds.
	// Defaults to 126000 seconds (35 hours).
	WeeklyCapacity *int `json:"weekly_capacity,omitempty"`
	// optional	The billable rate to use for this user when they are added to a project. Defaults to 0.
	DefaultHourlyRate *float64 `json:"default_hourly_rate,omitempty"`
	// optional	The cost rate to use for this user when calculating a project's costs vs billable amount. Defaults to 0.
	CostRate *float64 `json:"cost_rate,omitempty"`
	// of strings	optional	The role names assigned to this person.
	Roles []*string `json:"roles,omitempty"`
}

type UserUpdateRequest struct {
	// The first name of the user. Can't be updated if the user is inactive.
	FirstName *string `json:"first_name,omitempty"`
	// The last name of the user. Can't be updated if the user is inactive.
	LastName *string `json:"last_name,omitempty"`
	// The email address of the user. Can't be updated if the user is inactive.
	Email *string `json:"email,omitempty"`
	// The telephone number for the user.
	Telephone *string `json:"telephone,omitempty"`
	// The user's timezone. Defaults to the company's timezone. See a list of supported time zones.
	Timezone *string `json:"timezone,omitempty"`
	// Whether the user should be automatically added to future projects.
	HasAccessToAllFutureProjects *bool `json:"has_access_to_all_future_projects,omitempty"`
	// Whether the user is a contractor or an employee.
	IsContractor *bool `json:"is_contractor,omitempty"`
	// Whether the user has admin permissions.
	IsAdmin *bool `json:"is_admin,omitempty"`
	// Whether the user has project manager permissions.
	IsProjectManager *bool `json:"is_project_manager,omitempty"`
	// Whether the user can see billable rates on projects. Only applicable to project managers.
	CanSeeRates *bool `json:"can_see_rates,omitempty"`
	// Whether the user can create projects. Only applicable to project managers.
	CanCreateProjects *bool `json:"can_create_projects,omitempty"`
	// Whether the user can create invoices. Only applicable to project managers.
	CanCreateInvoices *bool `json:"can_create_invoices,omitempty"`
	// Whether the user is active or archived.
	IsActive *bool `json:"is_active,omitempty"`
	// The number of hours per week this person is available to work in seconds.
	WeeklyCapacity *int `json:"weekly_capacity,omitempty"`
	// The billable rate to use for this user when they are added to a project.
	DefaultHourlyRate *float64 `json:"default_hourly_rate,omitempty"`
	// The cost rate to use for this user when calculating a project's costs vs billable amount.
	CostRate *float64 `json:"cost_rate,omitempty"`
	// of strings	The role names assigned to this person.
	Roles []*string `json:"roles,omitempty"`
}

type UserListOptions struct {
	// Pass true to only return active projects and false to return inactive projects.
	IsActive bool `url:"is_active,omitempty"`
	// Only return projects that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

// List returns a list of your users.
func (s *UserService) List(ctx context.Context, opt *UserListOptions) (*UserList, *http.Response, error) {
	u := "users"

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	userList := new(UserList)

	resp, err := s.client.Do(ctx, req, userList)
	if err != nil {
		return nil, resp, err
	}

	return userList, resp, nil
}

// Get retrieves the user with the given ID.
func (s *UserService) Get(ctx context.Context, userID int64) (*User, *http.Response, error) {
	u := fmt.Sprintf("users/%d", userID)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
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

// Current retrieves the currently authenticated user.
func (s *UserService) Current(ctx context.Context) (*User, *http.Response, error) {
	u := "users/me"

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
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

// Create creates a new user object.
func (s *UserService) Create(ctx context.Context, data *UserCreateRequest) (*User, *http.Response, error) {
	u := "users"

	req, err := s.client.NewRequest(ctx, "POST", u, data)
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

// Update updates the specific user.
func (s *UserService) Update(
	ctx context.Context,
	userID int64,
	data *UserUpdateRequest,
) (*User, *http.Response, error) {
	u := fmt.Sprintf("users/%d", userID)

	req, err := s.client.NewRequest(ctx, "PATCH", u, data)
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

// Delete deletes a user.
func (s *UserService) Delete(ctx context.Context, userID int64) (*http.Response, error) {
	u := fmt.Sprintf("users/%d", userID)

	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
