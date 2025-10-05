package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type ProjectTaskAssignment struct {
	// Unique ID for the task assignment.
	ID *int64 `json:"id,omitempty"`
	// An object containing the id, name, and code of the associated project.
	Project *Project `json:"project,omitempty"`
	// An object containing the id and name of the associated task.
	Task *Task `json:"task,omitempty"`
	// Whether the task assignment is active or archived.
	IsActive *bool `json:"is_active,omitempty"`
	// Whether the task assignment is billable or not.
	// For example: if set to true, all time tracked on this project for the associated task will be marked as billable.
	Billable *bool `json:"billable,omitempty"`
	// Rate used when the project’s bill_by is Tasks.
	HourlyRate *float64 `json:"hourly_rate,omitempty"`
	// Budget used when the project’s budget_by is task or task_fees.
	Budget *float64 `json:"budget,omitempty"`
	// Date and time the task assignment was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Date and time the task assignment was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type ProjectTaskAssignmentCreateRequest struct {
	// required The ID of the task to associate with the project.
	TaskID *int64 `json:"task_id"`
	// optional Whether the task assignment is active or archived. Defaults to true
	IsActive *bool `json:"is_active,omitempty"`
	// optional Whether the task assignment is billable or not. Defaults to false.
	Billable *bool `json:"billable,omitempty"`
	// optional Rate used when the project’s bill_by is Tasks.
	// Defaults to null when billing by task hourly rate, otherwise 0.
	HourlyRate *float64 `json:"hourly_rate,omitempty"`
	// optional Budget used when the project’s budget_by is task or task_fees.
	Budget *float64 `json:"budget,omitempty"`
}

type ProjectTaskAssignmentList struct {
	TaskAssignments []*ProjectTaskAssignment `json:"task_assignments"`

	Pagination
}

func (p ProjectTaskAssignment) String() string {
	return Stringify(p)
}

func (p ProjectTaskAssignmentList) String() string {
	return Stringify(p)
}

type ProjectTaskAssignmentListOptions struct {
	// Pass true to only return active projects and false to return inactive projects.
	IsActive bool `url:"is_active,omitempty"`
	// Only return projects belonging to the client with the given ID.
	ClientID int64 `url:"client_id,omitempty"`
	// Only return projects that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

// ListTaskAssignments returns a list of your task assignments.
func (s *ProjectService) ListTaskAssignments(
	ctx context.Context,
	projectID int64,
	opt *ProjectTaskAssignmentListOptions,
) (*ProjectTaskAssignmentList, *http.Response, error) {
	u := fmt.Sprintf("projects/%d/task_assignments", projectID)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	projectTaskAssignmentList := new(ProjectTaskAssignmentList)

	resp, err := s.client.Do(ctx, req, projectTaskAssignmentList)
	if err != nil {
		return nil, resp, err
	}

	return projectTaskAssignmentList, resp, nil
}

// GetTaskAssignment retrieves the task assignment with the given ID.
func (s *ProjectService) GetTaskAssignment(
	ctx context.Context,
	projectID int64,
	taskAssignmentID int64,
) (*ProjectTaskAssignment, *http.Response, error) {
	u := fmt.Sprintf("projects/%d/task_assignments/%d", projectID, taskAssignmentID)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	projectTaskAssignment := new(ProjectTaskAssignment)

	resp, err := s.client.Do(ctx, req, projectTaskAssignment)
	if err != nil {
		return nil, resp, err
	}

	return projectTaskAssignment, resp, nil
}

// CreateTaskAssignment creates a new task assignment object.
func (s *ProjectService) CreateTaskAssignment(
	ctx context.Context,
	projectID int64,
	data *ProjectTaskAssignmentCreateRequest,
) (*ProjectTaskAssignment, *http.Response, error) {
	u := fmt.Sprintf("projects/%d/task_assignments", projectID)

	req, err := s.client.NewRequest(ctx, "POST", u, data)
	if err != nil {
		return nil, nil, err
	}

	projectTaskAssignment := new(ProjectTaskAssignment)

	resp, err := s.client.Do(ctx, req, projectTaskAssignment)
	if err != nil {
		return nil, resp, err
	}

	return projectTaskAssignment, resp, nil
}
