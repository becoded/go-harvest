package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type ProjectTaskAssignment struct {
	Id         *int64     `json:"id,omitempty"`          // Unique ID for the task assignment.
	Task       *Task      `json:"task,omitempty"`        // An object containing the id and name of the associated task.
	IsActive   *bool      `json:"is_active,omitempty"`   // Whether the task assignment is active or archived.
	Billable   *bool      `json:"billable,omitempty"`    // Whether the task assignment is billable or not. For example: if set to true, all time tracked on this project for the associated task will be marked as billable.
	HourlyRate *float64   `json:"hourly_rate,omitempty"` // Rate used when the project’s bill_by is Tasks.
	Budget     *float64   `json:"budget,omitempty"`      // Budget used when the project’s budget_by is task or task_fees.
	CreatedAt  *time.Time `json:"created_at,omitempty"`  // Date and time the task assignment was created.
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`  // Date and time the task assignment was last updated.
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
	ClientId int64 `url:"client_id,omitempty"`
	// Only return projects that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

func (s *ProjectService) ListTaskAssignments(ctx context.Context, projectId int64, opt *ProjectTaskAssignmentListOptions) (*ProjectTaskAssignmentList, *http.Response, error) {
	u := fmt.Sprintf("projects/%d/task_assignments", projectId)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	projectTaskAssignmentList := new(ProjectTaskAssignmentList)
	resp, err := s.client.Do(ctx, req, &projectTaskAssignmentList)
	if err != nil {
		return nil, resp, err
	}

	return projectTaskAssignmentList, resp, nil
}

func (s *ProjectService) GetTaskAssignment(ctx context.Context, projectId int64, taskAssignmentId int64) (*ProjectTaskAssignment, *http.Response, error) {
	u := fmt.Sprintf("projects/%d/task_assignments/%d", projectId, taskAssignmentId)
	req, err := s.client.NewRequest("GET", u, nil)
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
