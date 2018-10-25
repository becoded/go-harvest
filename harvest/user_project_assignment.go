package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

/** https://help.getharvest.com/api-v2/users-api/users/project-assignments/ **/

type UserProjectAssignment struct {
	Id *int64 `json:"id,omitempty"` // Unique ID for the project assignment.
	IsActive *bool `json:"is_active,omitempty"` // Whether the project assignment is active or archived.
	IsProjectManager *bool `json:"is_project_manager,omitempty"` // Determines if the user has project manager permissions for the project.
	HourlyRate *float64 `json:"hourly_rate,omitempty"` // Rate used when the project’s bill_by is People.
	Budget *float64 `json:"budget,omitempty"` // Budget used when the project’s budget_by is person.
	CreatedAt *time.Time `json:"created_at,omitempty"` // Date and time the project assignment was created.
	UpdatedAt *time.Time `json:"updated_at,omitempty"` // Date and time the project assignment was last updated.
	Project *Project `json:"project,omitempty"` // An object containing the assigned project id, name, and code.
	Client *Client `json:"client,omitempty"` // An object containing the project’s client id and name.
	TaskAssignments *[]ProjectTaskAssignment `json:"task_assignments,omitempty"` // Array of task assignment objects associated with the project.
}

type UserProjectAssignmentList struct {
	UserAssignments []*UserProjectAssignment `json:"user_assignments"`

	Pagination
}

func (p UserProjectAssignment) String() string {
	return Stringify(p)
}

func (p UserProjectAssignmentList) String() string {
	return Stringify(p)
}

type UserProjectAssignmentListOptions struct {
	// Pass true to only return active projects and false to return inactive projects.
	IsActive bool `url:"is_active,omitempty"`
	// Only return projects belonging to the client with the given ID.
	ClientId int64 `url:"client_id,omitempty"`
	// Only return projects that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

func (s *ProjectService) ListUserAssignments(ctx context.Context, projectId int64, opt *UserProjectAssignmentListOptions) (*UserProjectAssignmentList, *http.Response, error) {
	u := fmt.Sprintf("projects/%d/user_assignments", projectId)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
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

func (s *ProjectService) GetUserAssignment(ctx context.Context, projectId int64, userAssignmentId int64) (*UserProjectAssignment, *http.Response, error) {
	u := fmt.Sprintf("projects/%d/user_assignments/%d", projectId, userAssignmentId)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	UserProjectAssignment := new(UserProjectAssignment)
	resp, err := s.client.Do(ctx, req, UserProjectAssignment)
	if err != nil {
		return nil, resp, err
	}

	return UserProjectAssignment, resp, nil
}
