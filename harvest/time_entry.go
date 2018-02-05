package harvest


import (
	"context"
	"fmt"
	"time"
	"net/http"
)

// TimeEntryService handles communication with the issue related
// methods of the Harvest API.
//
// Harvest API docs: https://help.getharvest.com/api-v2/timesheets-api/timesheets/time-entries/
type TimesheetService service


type TimeEntry struct {
	Id *int64 `json:"id,omitempty"` // Unique ID for the time entry.
	SpentDate *Date `json:"spent_date,omitempty"` // Date of the time entry.
	User *User `json:"user,omitempty"` // An object containing the id and name of the associated user.
	UserAssignment *ProjectUserAssignment `json:"user_assignment,omitempty"` // A user assignment object of the associated user.
	Client *Client `json:"client,omitempty"` // An object containing the id and name of the associated client.
	Project *Project `json:"project,omitempty"` // An object containing the id and name of the associated project.
	Task *Task `json:"task,omitempty"` // An object containing the id and name of the associated task.
	TaskAssignment *ProjectTaskAssignment `json:"task_assignment,omitempty"` // A task assignment object of the associated task.
	ExternalReference *ExternalReference `json:"external_reference,omitempty"` // An object containing the id, group_id, permalink, service, and service_icon_url of the associated external reference.
	Invoice *Invoice `json:"invoice,omitempty"` // Once the time entry has been invoiced, this field will include the associated invoice’s id and number.
	Hours *float64 `json:"hours,omitempty"` // Number of (decimal time) hours tracked in this time entry.
	Notes *string `json:"notes,omitempty"` // Notes attached to the time entry.
	IsLocked *bool `json:"is_locked,omitempty"` // Whether or not the time entry has been locked.
	LockedReason *string `json:"locked_reason,omitempty"` // Why the time entry has been locked.
	IsClosed *bool `json:"is_closed,omitempty"` // Whether or not the time entry has been approved via Timesheet Approval.
	IsBilled *bool `json:"is_billed,omitempty"` // Whether or not the time entry has been marked as invoiced.
	TimerStartedAt *time.Time `json:"timer_started_at,omitempty"` // Date and time the timer was started (if tracking by duration).
	StartedTime *time.Time `json:"started_time,omitempty"` // Time the time entry was started (if tracking by start/end times).
	EndedTime *time.Time `json:"ended_time,omitempty"` // Time the time entry was ended (if tracking by start/end times).
	IsRunning *bool `json:"is_running,omitempty"` // Whether or not the time entry is currently running.
	Billable *bool `json:"billable,omitempty"` // Whether or not the time entry is billable.
	Budgeted *bool `json:"budgeted,omitempty"` // Whether or not the time entry counts towards the project budget.
	BillableRate *float64 `json:"billable_rate,omitempty"` // The billable rate for the time entry.
	CostRate *float64 `json:"cost_rate,omitempty"` // The cost rate for the time entry.
	CreatedAt *time.Time `json:"created_at,omitempty"` // Date and time the time entry was created.
	UpdatedAt *time.Time `json:"updated_at,omitempty"` // Date and time the time entry was last updated.
}

type ExternalReference struct {
	Id *string `json:"id,omitempty"`
	GroupId *string `json:"group_id,omitempty"`
	Permalink *string `json:"permalink,omitempty"`
	Service *string `json:"service,omitempty"`
	ServiceIconUrl *string `json:"service_icon_url,omitempty"`
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
	UserId *int64 `json:"user_id,omitempty"` // Only return time entries belonging to the user with the given ID.
	ClientId *int64 `json:"client_id,omitempty"` // Only return time entries belonging to the client with the given ID.
	ProjectId *int64 `json:"project_id,omitempty"` // Only return time entries belonging to the project with the given ID.
	IsBilled *bool `json:"is_billed,omitempty"` // Pass true to only return time entries that have been invoiced and false to return time entries that have not been invoiced.
	IsRunning *bool `json:"is_running,omitempty"` // Pass true to only return running time entries and false to return non-running time entries.
	UpdatedSince *time.Time `json:"updated_since,omitempty"` // Only return time entries that have been updated since the given date and time.
	From *Date `json:"from,omitempty"` // Only return time entries with a spent_date on or after the given date.
	To *Date `json:"to,omitempty"` // Only return time entries with a spent_date on or before the given date.

	ListOptions
}

func (s *TimesheetService) List(ctx context.Context, opt *TimeEntryListOptions) (*TimeEntryList, *http.Response, error) {
	u := "time_entries"
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
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

func (s *TimesheetService) Get(ctx context.Context, timeEntryId int64) (*TimeEntry, *http.Response, error) {
	u := fmt.Sprintf("time_entries/%d", timeEntryId)
	req, err := s.client.NewRequest("GET", u, nil)
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