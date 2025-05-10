package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

/** https://help.getharvest.com/api-v2/projects-api/projects/user-assignments/#the-user-assignment-object **/

type ProjectUserAssignment struct {
	// Unique ID for the user assignment.
	ID *int64 `json:"id,omitempty"`
	// An object containing the id, name, and code of the associated project.
	Project *Project `json:"project,omitempty"`
	// An object containing the id and name of the associated user.
	User *User `json:"user,omitempty"`
	// Whether the user assignment is active or archived.
	IsActive *bool `json:"is_active,omitempty"`
	// Determines if the user has project manager permissions for the project.
	IsProjectManager *bool `json:"is_project_manager,omitempty"`
	// Rate used when the project’s bill_by is People.
	HourlyRate *float64 `json:"hourly_rate,omitempty"`
	// Budget used when the project’s budget_by is person.
	Budget *float64 `json:"budget,omitempty"`
	// Date and time the user assignment was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Date and time the user assignment was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type ProjectUserAssignmentList struct {
	UserAssignments []*ProjectUserAssignment `json:"user_assignments"`

	Pagination
}

func (p ProjectUserAssignment) String() string {
	return Stringify(p)
}

func (p ProjectUserAssignmentList) String() string {
	return Stringify(p)
}

type ProjectUserAssignmentListOptions struct {
	// Pass true to only return active projects and false to return inactive projects.
	IsActive bool `url:"is_active,omitempty"`
	// Only return projects belonging to the client with the given ID.
	ClientID int64 `url:"client_id,omitempty"`
	// Only return projects that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

// ListUserAssignments returns a list of your projects user assignments.
func (s *ProjectService) ListUserAssignments(
	ctx context.Context,
	projectID int64,
	opt *ProjectUserAssignmentListOptions,
) (*ProjectUserAssignmentList, *http.Response, error) {
	u := fmt.Sprintf("projects/%d/user_assignments", projectID)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	projectUserAssignmentList := new(ProjectUserAssignmentList)

	resp, err := s.client.Do(ctx, req, &projectUserAssignmentList)
	if err != nil {
		return nil, resp, err
	}

	return projectUserAssignmentList, resp, nil
}

// GetUserAssignment retrieves the user assignment with the given ID.
func (s *ProjectService) GetUserAssignment(
	ctx context.Context,
	projectID int64,
	userAssignmentID int64,
) (*ProjectUserAssignment, *http.Response, error) {
	u := fmt.Sprintf("projects/%d/user_assignments/%d", projectID, userAssignmentID)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	projectUserAssignment := new(ProjectUserAssignment)

	resp, err := s.client.Do(ctx, req, projectUserAssignment)
	if err != nil {
		return nil, resp, err
	}

	return projectUserAssignment, resp, nil
}
