package harvest


import (
"context"
"fmt"
"time"
	"net/http"
)

/** https://help.getharvest.com/api-v2/projects-api/projects/user-assignments/#the-user-assignment-object **/

type ProjectUserAssignment struct {
	Id *int64 `json:"id,omitempty"` // Unique ID for the user assignment.
	User *User `json:"user,omitempty"` // An object containing the id and name of the associated user.
	IsActive *bool `json:"is_active,omitempty"` // Whether the user assignment is active or archived.
	IsProjectManager *bool `json:"is_project_manager,omitempty"` // Determines if the user has project manager permissions for the project.
	HourlyRate *float64 `json:"hourly_rate,omitempty"` // Rate used when the project’s bill_by is People.
	Budget *float64 `json:"budget,omitempty"` // Budget used when the project’s budget_by is person.
	CreatedAt *time.Time `json:"created_at,omitempty"` // Date and time the user assignment was created.
	UpdatedAt *time.Time `json:"updated_at,omitempty"` // Date and time the user assignment was last updated.
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
	IsActive	bool `url:"is_active,omitempty"`
	// Only return projects belonging to the client with the given ID.
	ClientId	int64 `url:"client_id,omitempty"`
	// Only return projects that have been updated since the given date and time.
	UpdatedSince	time.Time `url:"updated_since,omitempty"`

	ListOptions
}

func (s *ProjectService) ListUserAssignments(ctx context.Context, projectId int64, opt *ProjectUserAssignmentListOptions) (*ProjectUserAssignmentList, *http.Response, error) {
	u := fmt.Sprintf("projects/%d/user_assignments", projectId)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
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

func (s *ProjectService) GetUserAssignment(ctx context.Context, projectId int64, userAssignmentId int64) (*ProjectUserAssignment, *http.Response, error) {
	u := fmt.Sprintf("projects/%d/user_assignments/%d", projectId, userAssignmentId)
	req, err := s.client.NewRequest("GET", u, nil)
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
