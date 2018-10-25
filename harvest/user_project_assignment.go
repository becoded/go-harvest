package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

/** https://help.getharvest.com/api-v2/users-api/users/project-assignments/#the-project-assignment-object **/

type UserProjectAssignment struct {
	Id               *int64                   `json:"id,omitempty"`                 // Unique ID for the project assignment.
	IsActive         *bool                    `json:"is_active,omitempty"`          // Whether the project assignment is active or archived.
	IsProjectManager *bool                    `json:"is_project_manager,omitempty"` // Determines if the user has project manager permissions for the project.
	HourlyRate       *float64                 `json:"hourly_rate,omitempty"`        // Rate used when the project’s bill_by is People.
	Budget           *float64                 `json:"budget,omitempty"`             // Budget used when the project’s budget_by is person.
	CreatedAt        *time.Time               `json:"created_at,omitempty"`         // Date and time the project assignment was created.
	UpdatedAt        *time.Time               `json:"updated_at,omitempty"`         // Date and time the project assignment was last updated.
	Project          *Project                 `json:"project,omitempty"`            // An object containing the assigned project id, name, and code.
	Client           *Client                  `json:"client,omitempty"`             // An object containing the project’s client id and name.
	TaskAssignments  []*ProjectTaskAssignment `json:"task_assignments,omitempty"`   // Array of task assignment objects associated with the project.
}

type UserProjectAssignmentList struct {
	UserAssignments []*UserProjectAssignment `json:"project_assignments"`

	Pagination
}

type UserProjectAssignmentListOptions struct {
	// Only return project assignments that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

func (s *ProjectService) ListProjectAssignments(ctx context.Context, userId int64, opt *UserProjectAssignmentListOptions) (*UserProjectAssignmentList, *http.Response, error) {
	u := fmt.Sprintf("users/%d/project_assignments", userId)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	userProjectAssignmentList := new(UserProjectAssignmentList)
	resp, err := s.client.Do(ctx, req, &userProjectAssignmentList)
	if err != nil {
		return nil, resp, err
	}

	return userProjectAssignmentList, resp, nil
}
