package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

/** https://help.getharvest.com/api-v2/users-api/users/project-assignments/ **/

type UserProjectAssignment struct {
	// Unique ID for the project assignment.
	ID *int64 `json:"id,omitempty"`
	// Whether the project assignment is active or archived.
	IsActive *bool `json:"is_active,omitempty"`
	// Determines if the user has project manager permissions for the project.
	IsProjectManager *bool `json:"is_project_manager,omitempty"`
	// Determines which billable rate(s) will be used on the project for this user when bill_by is People.
	// When true, the project will use the user’s default billable rates.
	// When false, the project will use the custom rate defined on this user assignment.
	UseDefaultRates *bool `json:"use_default_rates,omitempty"`
	// Rate used when the project’s bill_by is People.
	HourlyRate *float64 `json:"hourly_rate,omitempty"`
	// Budget used when the project’s budget_by is person.
	Budget *float64 `json:"budget,omitempty"`
	// Date and time the project assignment was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Date and time the project assignment was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	// An object containing the assigned project id, name, and code.
	Project *Project `json:"project,omitempty"`
	// An object containing the project’s client id and name.
	Client *Client `json:"client,omitempty"`
	// Array of task assignment objects associated with the project.
	TaskAssignments *[]ProjectTaskAssignment `json:"task_assignments,omitempty"`
}

type UserProjectAssignmentList struct {
	ProjectAssignments []*UserProjectAssignment `json:"project_assignments"`

	Pagination
}

func (p UserProjectAssignment) String() string {
	return Stringify(p)
}

func (p UserProjectAssignmentList) String() string {
	return Stringify(p)
}

type UserProjectAssignmentListOptions struct {
	// Only return projects that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

type MyProjectAssignmentListOptions struct {
	ListOptions
}

// ListProjectAssignments returns a list of active project assignments for the user.
func (s *UserService) ListProjectAssignments(
	ctx context.Context,
	userID int64,
	opt *UserProjectAssignmentListOptions,
) (*UserProjectAssignmentList, *http.Response, error) {
	u := fmt.Sprintf("users/%d/project_assignments", userID)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	list := new(UserProjectAssignmentList)

	resp, err := s.client.Do(ctx, req, &list)
	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}

// GetMyProjectAssignments returns a list of your active project assignments.
func (s *UserService) GetMyProjectAssignments(
	ctx context.Context,
	opt *MyProjectAssignmentListOptions,
) (*UserProjectAssignmentList, *http.Response, error) {
	u := "users/me/project_assignments"

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	list := new(UserProjectAssignmentList)

	resp, err := s.client.Do(ctx, req, &list)
	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}
