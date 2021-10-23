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
	UserAssignments []*UserProjectAssignment `json:"project_assignments"`

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

func (s *ProjectService) ListProjectAssignments(
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

	UserProjectAssignmentList := new(UserProjectAssignmentList)

	resp, err := s.client.Do(ctx, req, &UserProjectAssignmentList)
	if err != nil {
		return nil, resp, err
	}

	return UserProjectAssignmentList, resp, nil
}

func (s *ProjectService) GetMyProjectAssignments(
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

	UserProjectAssignmentList := new(UserProjectAssignmentList)

	resp, err := s.client.Do(ctx, req, &UserProjectAssignmentList)
	if err != nil {
		return nil, resp, err
	}

	return UserProjectAssignmentList, resp, nil
}
