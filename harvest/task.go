package harvest

import (
	"context"
	"time"
	"fmt"
	"net/http"
)

// TaskService handles communication with the task related
// methods of the Harvest API.
//
// Harvest API docs: https://help.getharvest.com/api-v2/tasks-api/tasks/tasks/
type TaskService service

type Task struct {
	Id                *int64     `json:"id,omitempty"`                  // Unique ID for the task.
	Name              *string    `json:"name,omitempty"`                // The name of the task.
	BillableByDefault *bool      `json:"billable_by_default,omitempty"` // Used in determining whether default tasks should be marked billable when creating a new project.
	DefaultHourlyRate *float64   `json:"default_hourly_rate,omitempty"` // The hourly rate to use for this task when it is added to a project.
	IsDefault         *bool      `json:"is_default,omitempty"`          // Whether this task should be automatically added to future projects.
	IsActive          *bool      `json:"is_active,omitempty"`           // Whether this task is active or archived.
	CreatedAt         *time.Time `json:"created_at,omitempty"`          // Date and time the task was created.
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`          // Date and time the task was last updated.
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
	Name              *string  `json:"name"`                          // required	The name of the task.
	BillableByDefault *bool    `json:"billable_by_default,omitempty"` // optional	Used in determining whether default tasks should be marked billable when creating a new project. Defaults to true.
	DefaultHourlyRate *float64 `json:"default_hourly_rate,omitempty"` // optional	The default hourly rate to use for this task when it is added to a project. Defaults to 0.
	IsDefault         *bool    `json:"is_default,omitempty"`          // optional	Whether this task should be automatically added to future projects. Defaults to false.
	IsActive          *bool    `json:"is_active,omitempty"`           // optional	Whether this task is active or archived. Defaults to true.
}

type TaskUpdateRequest struct {
	Name                *string  `json:"name,omitempty"`                // The name of the task.
	Billable_by_default *bool    `json:"billable_by_default,omitempty"` // Used in determining whether default tasks should be marked billable when creating a new project.
	Default_hourly_rate *float64 `json:"default_hourly_rate,omitempty"` // The default hourly rate to use for this task when it is added to a project.
	Is_default          *bool    `json:"is_default,omitempty"`          // Whether this task should be automatically added to future projects.
	Is_active           *bool    `json:"is_active,omitempty"`           // Whether this task is active or archived.
}

func (s *TaskService) List(ctx context.Context, opt *TaskListOptions) (*TaskList, *http.Response, error) {
	u := "tasks"
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	taskList := new(TaskList)
	resp, err := s.client.Do(ctx, req, &taskList)
	if err != nil {
		return nil, resp, err
	}

	return taskList, resp, nil
}

func (s *TaskService) Get(ctx context.Context, taskId int64) (*Task, *http.Response, error) {
	u := fmt.Sprintf("tasks/%d", taskId)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	task := new(Task)
	resp, err := s.client.Do(ctx, req, &task)
	if err != nil {
		return nil, resp, err
	}

	return task, resp, nil
}

func (s *TaskService) CreateTask(ctx context.Context, data *TaskCreateRequest) (*Task, *http.Response, error) {
	u := "tasks"

	req, err := s.client.NewRequest("POST", u, data)
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

func (s *TaskService) UpdateTask(ctx context.Context, taskId int64, data *TaskUpdateRequest) (*Task, *http.Response, error) {
	u := fmt.Sprintf("tasks/%d", taskId)

	req, err := s.client.NewRequest("PATCH", u, data)
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

func (s *TaskService) DeleteTask(ctx context.Context, taskId int64) (*http.Response, error) {
	u := fmt.Sprintf("tasks/%d", taskId)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}