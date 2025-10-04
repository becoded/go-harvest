package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// TaskService handles communication with the task related
// methods of the Harvest API.
//
// Harvest API docs: https://help.getharvest.com/api-v2/tasks-api/tasks/tasks/
type TaskService service

type Task struct {
	// Unique ID for the task.
	ID *int64 `json:"id,omitempty"`
	// The name of the task.
	Name *string `json:"name,omitempty"`
	// Used in determining whether default tasks should be marked billable when creating a new project.
	BillableByDefault *bool `json:"billable_by_default,omitempty"`
	// The hourly rate to use for this task when it is added to a project.
	DefaultHourlyRate *float64 `json:"default_hourly_rate,omitempty"`
	// Whether this task should be automatically added to future projects.
	IsDefault *bool `json:"is_default,omitempty"`
	// Whether this task is active or archived.
	IsActive *bool `json:"is_active,omitempty"`
	// Date and time the task was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Date and time the task was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type TaskList struct {
	Tasks []*Task `json:"tasks"`

	Pagination
}

func (c Task) String() string {
	return Stringify(c)
}

func (c TaskList) String() string {
	return Stringify(c)
}

type TaskListOptions struct {
	// Pass true to only return active projects and false to return inactive projects.
	IsActive bool `url:"is_active,omitempty"`
	// Only return projects that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

type TaskCreateRequest struct {
	// required	The name of the task.
	Name *string `json:"name"`
	// optional	Used in determining whether default tasks should be marked billable when creating a new project.
	// Defaults to true.
	BillableByDefault *bool `json:"billable_by_default,omitempty"`
	// optional	The default hourly rate to use for this task when it is added to a project. Defaults to 0.
	DefaultHourlyRate *float64 `json:"default_hourly_rate,omitempty"`
	// optional	Whether this task should be automatically added to future projects. Defaults to false.
	IsDefault *bool `json:"is_default,omitempty"`
	// optional	Whether this task is active or archived. Defaults to true.
	IsActive *bool `json:"is_active,omitempty"`
}

type TaskUpdateRequest struct {
	// The name of the task.
	Name *string `json:"name,omitempty"`
	// Used in determining whether default tasks should be marked billable when creating a new project.
	BillableByDefault *bool `json:"billable_by_default,omitempty"`
	// The default hourly rate to use for this task when it is added to a project.
	DefaultHourlyRate *float64 `json:"default_hourly_rate,omitempty"`
	// Whether this task should be automatically added to future projects.
	IsDefault *bool `json:"is_default,omitempty"`
	// Whether this task is active or archived.
	IsActive *bool `json:"is_active,omitempty"`
}

// List returns a list of your tasks.
func (s *TaskService) List(ctx context.Context, opt *TaskListOptions) (*TaskList, *http.Response, error) {
	u := "tasks"

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	taskList := new(TaskList)

	resp, err := s.client.Do(ctx, req, taskList)
	if err != nil {
		return nil, resp, err
	}

	return taskList, resp, nil
}

// Get retrieves the task with the given ID.
func (s *TaskService) Get(ctx context.Context, taskID int64) (*Task, *http.Response, error) {
	u := fmt.Sprintf("tasks/%d", taskID)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	task := new(Task)

	resp, err := s.client.Do(ctx, req, task)
	if err != nil {
		return nil, resp, err
	}

	return task, resp, nil
}

// Create creates a new task object.
func (s *TaskService) Create(ctx context.Context, data *TaskCreateRequest) (*Task, *http.Response, error) {
	u := "tasks"

	req, err := s.client.NewRequest(ctx, "POST", u, data)
	if err != nil {
		return nil, nil, err
	}

	task := new(Task)

	resp, err := s.client.Do(ctx, req, task)
	if err != nil {
		return nil, resp, err
	}

	return task, resp, nil
}

// Update updates the specific task.
func (s *TaskService) Update(
	ctx context.Context,
	taskID int64,
	data *TaskUpdateRequest,
) (*Task, *http.Response, error) {
	u := fmt.Sprintf("tasks/%d", taskID)

	req, err := s.client.NewRequest(ctx, "PATCH", u, data)
	if err != nil {
		return nil, nil, err
	}

	task := new(Task)

	resp, err := s.client.Do(ctx, req, task)
	if err != nil {
		return nil, resp, err
	}

	return task, resp, nil
}

// Delete deletes a task.
func (s *TaskService) Delete(ctx context.Context, taskID int64) (*http.Response, error) {
	u := fmt.Sprintf("tasks/%d", taskID)

	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
