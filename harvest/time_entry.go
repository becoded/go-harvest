package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// TimeEntryService handles communication with the issue related
// methods of the Harvest API.
//
// Harvest API docs: https://help.getharvest.com/api-v2/timesheets-api/timesheets/time-entries/
type TimesheetService service

const basePathTimeEntries = "time_entries"

type TimeEntry struct {
	// Unique ID for the time entry.
	ID *int64 `json:"id,omitempty"`
	// Date of the time entry.
	SpentDate *Date `json:"spent_date,omitempty"`
	// An object containing the id and name of the associated user.
	User *User `json:"user,omitempty"`
	// A user assignment object of the associated user.
	UserAssignment *ProjectUserAssignment `json:"user_assignment,omitempty"`
	// An object containing the id and name of the associated client.
	Client *Client `json:"client,omitempty"`
	// An object containing the id and name of the associated project.
	Project *Project `json:"project,omitempty"`
	// An object containing the id and name of the associated task.
	Task *Task `json:"task,omitempty"`
	// A task assignment object of the associated task.
	TaskAssignment *ProjectTaskAssignment `json:"task_assignment,omitempty"`
	// An object containing the id, group_id, permalink, service, and service_icon_url
	// of the associated external reference.
	ExternalReference *ExternalReference `json:"external_reference,omitempty"`
	// Once the time entry has been invoiced, this field will include the associated invoice’s id and number.
	Invoice *Invoice `json:"invoice,omitempty"`
	// Number of (decimal time) hours tracked in this time entry.
	Hours *float64 `json:"hours,omitempty"`
	// Number of (decimal time) hours tracked in this time entry used in summary reports and invoices.
	// This value is rounded according to the Time Rounding setting in your Preferences.
	RoundedHours *float64 `json:"rounded_hours,omitempty"`
	// Notes attached to the time entry.
	Notes *string `json:"notes,omitempty"`
	// Whether or not the time entry has been locked.
	IsLocked *bool `json:"is_locked,omitempty"`
	// Why the time entry has been locked.
	LockedReason *string `json:"locked_reason,omitempty"`
	// Whether or not the time entry has been approved via Timesheet Approval.
	IsClosed *bool `json:"is_closed,omitempty"`
	// Whether or not the time entry has been marked as invoiced.
	IsBilled *bool `json:"is_billed,omitempty"`
	// Date and time the timer was started (if tracking by duration).
	TimerStartedAt *time.Time `json:"timer_started_at,omitempty"`
	// Time the time entry was started (if tracking by start/end times).
	StartedTime *Time `json:"started_time,omitempty"`
	// Time the time entry was ended (if tracking by start/end times).
	EndedTime *Time `json:"ended_time,omitempty"`
	// Whether or not the time entry is currently running.
	IsRunning *bool `json:"is_running,omitempty"`
	// Whether or not the time entry is billable.
	Billable *bool `json:"billable,omitempty"`
	// Whether or not the time entry counts towards the project budget.
	Budgeted *bool `json:"budgeted,omitempty"`
	// The billable rate for the time entry.
	BillableRate *float64 `json:"billable_rate,omitempty"`
	// The cost rate for the time entry.
	CostRate *float64 `json:"cost_rate,omitempty"`
	// Date and time the time entry was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Date and time the time entry was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type ExternalReference struct {
	ID             *string `json:"id,omitempty"`
	GroupID        *string `json:"group_id,omitempty"`
	Permalink      *string `json:"permalink,omitempty"`
	Service        *string `json:"service,omitempty"`
	ServiceIconURL *string `json:"service_icon_url,omitempty"`
}

type TimeEntryList struct {
	TimeEntries []*TimeEntry `json:"time_entries"`

	Pagination
}

func (p TimeEntry) String() string {
	return Stringify(p)
}

func (p TimeEntryList) String() string {
	return Stringify(p)
}

type TimeEntryListOptions struct {
	// Only return time entries belonging to the user with the given ID.
	UserID *int64 `url:"user_id,omitempty"`
	// Only return time entries belonging to the client with the given ID.
	ClientID *int64 `url:"client_id,omitempty"`
	// Only return time entries belonging to the project with the given ID.
	ProjectID *int64 `url:"project_id,omitempty"`
	// Only return time entries belonging to the task with the given ID.
	TaskID *int64 `url:"task_id,omitempty"`
	// Pass true to only return time entries that have been invoiced and false to return time
	// entries that have not been invoiced.
	IsBilled *bool `url:"is_billed,omitempty"`
	// Pass true to only return running time entries and false to return non-running time entries.
	IsRunning *bool `url:"is_running,omitempty"`
	// Only return time entries that have been updated since the given date and time.*/
	UpdatedSince *time.Time `url:"updated_since,omitempty"`
	// Only return time entries with a spent_date on or after the given date.
	From *Date `url:"from,omitempty"`
	// Only return time entries with a spent_date on or before the given date.
	To *Date `url:"to,omitempty"`

	ListOptions
}

type TimeEntryCreateViaDuration struct {
	// optional	The ID of the user to associate with the time entry. Defaults to the currently authenticated user’s ID.
	UserID *int64 `json:"user_id,omitempty"`
	// required	The ID of the project to associate with the time entry.
	ProjectID *int64 `json:"project_id"`
	// required	The ID of the task to associate with the time entry.
	TaskID *int64 `json:"task_id"`
	// required	The ISO 8601 formatted date the time entry was spent.
	SpentDate *Date `json:"spent_date"`
	// optional	The current amount of time tracked.
	// If provided, the time entry will be created with the specified hours and is_running will be set to false.
	// If not provided, hours will be set to 0.0 and is_running will be set to true.
	Hours *float64 `json:"hours,omitempty"`
	// optional	Any notes to be associated with the time entry.
	Notes *string `json:"notes,omitempty"`
	// TO DO
	// optional	An object containing the id, group_id, and permalink of the external reference.
	// ExternalReference *object `json:"external_reference,omitempty"`
}

type TimeEntryCreateViaStartEndTime struct {
	// optional	The ID of the user to associate with the time entry. Defaults to the currently authenticated user’s ID.
	UserID *int64 `json:"user_id,omitempty"`
	// required	The ID of the project to associate with the time entry.
	ProjectID *int64 `json:"project_id"`
	// required	The ID of the task to associate with the time entry.
	TaskID *int64 `json:"task_id"`
	// required	The ISO 8601 formatted date the time entry was spent.
	SpentDate *Date `json:"spent_date"`
	// optional	The time the entry started. Defaults to the current time. Example: “8:00am”.
	StartedTime *Time `json:"started_time,omitempty"`
	// optional	The time the entry ended. If provided, is_running will be set to false.
	// If not provided, is_running will be set to true.
	EndedTime *Time `json:"ended_time,omitempty"`
	// optional	Any notes to be associated with the time entry.
	Notes *string `json:"notes,omitempty"`
	// To do
	// optional	An object containing the id, group_id, and permalink of the external reference.
	// External_reference *object `json:"external_reference,omitempty"`
}

type TimeEntryUpdate struct {
	// required	The ID of the project to associate with the time entry.
	ProjectID *int64 `json:"project_id"`
	// required	The ID of the task to associate with the time entry.
	TaskID *int64 `json:"task_id"`
	// required	The ISO 8601 formatted date the time entry was spent.
	SpentDate *Date `json:"spent_date"`
	// optional	The time the entry started. Defaults to the current time. Example: “8:00am”.
	StartedTime *Time `json:"started_time,omitempty"`
	// optional	The time the entry ended. If provided, is_running will be set to false.
	// If not provided, is_running will be set to true.
	EndedTime *Time `json:"ended_time,omitempty"`
	// optional	The current amount of time tracked.
	// If provided, the time entry will be created with the specified hours and is_running will be set to false.
	// If not provided, hours will be set to 0.0 and is_running will be set to true.
	Hours *float64 `json:"hours,omitempty"`
	// optional	Any notes to be associated with the time entry.
	Notes *string `json:"notes,omitempty"`
	// TO DO
	// optional	An object containing the id, group_id, and permalink of the external reference.
	// ExternalReference *object `json:"external_reference,omitempty"`
}

// List returns a list of time entries.
func (s *TimesheetService) List(
	ctx context.Context,
	opt *TimeEntryListOptions,
) (*TimeEntryList, *http.Response, error) {
	u := basePathTimeEntries

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	timeEntryList := new(TimeEntryList)

	resp, err := s.client.Do(ctx, req, &timeEntryList)
	if err != nil {
		return nil, resp, err
	}

	return timeEntryList, resp, nil
}

// Get retrieves the time entry with the given ID.
func (s *TimesheetService) Get(ctx context.Context, timeEntryID int64) (*TimeEntry, *http.Response, error) {
	u := fmt.Sprintf("%s/%d", basePathTimeEntries, timeEntryID)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	timeEntry := new(TimeEntry)

	resp, err := s.client.Do(ctx, req, timeEntry)
	if err != nil {
		return nil, resp, err
	}

	return timeEntry, resp, nil
}

// CreateTimeEntryViaDuration creates a time entry object via duration.
func (s *TimesheetService) CreateTimeEntryViaDuration(
	ctx context.Context,
	data *TimeEntryCreateViaDuration,
) (*TimeEntry, *http.Response, error) {
	u := basePathTimeEntries

	req, err := s.client.NewRequest(ctx, "POST", u, data)
	if err != nil {
		return nil, nil, err
	}

	timeEntry := new(TimeEntry)

	resp, err := s.client.Do(ctx, req, timeEntry)
	if err != nil {
		return nil, resp, err
	}

	return timeEntry, resp, nil
}

// CreateTimeEntryViaStartEndTime creates a time entry object via start and end time.
func (s *TimesheetService) CreateTimeEntryViaStartEndTime(
	ctx context.Context,
	data *TimeEntryCreateViaStartEndTime,
) (*TimeEntry, *http.Response, error) {
	u := basePathTimeEntries

	req, err := s.client.NewRequest(ctx, "POST", u, data)
	if err != nil {
		return nil, nil, err
	}

	timeEntry := new(TimeEntry)

	resp, err := s.client.Do(ctx, req, timeEntry)
	if err != nil {
		return nil, resp, err
	}

	return timeEntry, resp, nil
}

// UpdateTimeEntry updates the specific time entry.
func (s *TimesheetService) UpdateTimeEntry(
	ctx context.Context,
	timeEntryID int64,
	data *TimeEntryUpdate,
) (*TimeEntry, *http.Response, error) {
	u := fmt.Sprintf("%s/%d", basePathTimeEntries, timeEntryID)

	req, err := s.client.NewRequest(ctx, "PATCH", u, data)
	if err != nil {
		return nil, nil, err
	}

	timeEntry := new(TimeEntry)

	resp, err := s.client.Do(ctx, req, timeEntry)
	if err != nil {
		return nil, resp, err
	}

	return timeEntry, resp, nil
}

// DeleteTimeEntry deletes a time entry.
func (s *TimesheetService) DeleteTimeEntry(ctx context.Context, timeEntryID int64) (*http.Response, error) {
	u := fmt.Sprintf("%s/%d", basePathTimeEntries, timeEntryID)

	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// RestartTimeEntry restarts a stopped time entry.
func (s *TimesheetService) RestartTimeEntry(
	ctx context.Context,
	timeEntryID int64,
) (*TimeEntry, *http.Response, error) {
	u := fmt.Sprintf("%s/%d/restart", basePathTimeEntries, timeEntryID)

	req, err := s.client.NewRequest(ctx, "PATCH", u, nil)
	if err != nil {
		return nil, nil, err
	}

	timeEntry := new(TimeEntry)

	resp, err := s.client.Do(ctx, req, timeEntry)
	if err != nil {
		return nil, resp, err
	}

	return timeEntry, resp, nil
}

// StopTimeEntry stops a running time entry.
func (s *TimesheetService) StopTimeEntry(ctx context.Context, timeEntryID int64) (*TimeEntry, *http.Response, error) {
	u := fmt.Sprintf("%s/%d/stop", basePathTimeEntries, timeEntryID)

	req, err := s.client.NewRequest(ctx, "PATCH", u, nil)
	if err != nil {
		return nil, nil, err
	}

	timeEntry := new(TimeEntry)

	resp, err := s.client.Do(ctx, req, timeEntry)
	if err != nil {
		return nil, resp, err
	}

	return timeEntry, resp, nil
}
