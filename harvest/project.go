package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// ProjectService handles communication with the issue related
// methods of the Harvest API.
//
// Harvest API docs: https://help.getharvest.com/api-v2/projects-api/projects/projects/
type ProjectService service

type Project struct {
	// Unique ID for the project.
	ID *int64 `json:"id,omitempty"`
	// An object containing the project’s client id, name, and currency.
	Client *Client `json:"client,omitempty"`
	// Unique name for the project.
	Name *string `json:"name,omitempty"`
	// The code associated with the project.
	Code *string `json:"code,omitempty"`
	// Whether the project is active or archived.
	IsActive *bool `json:"is_active,omitempty"`
	// Whether the project is billable or not.
	IsBillable *bool `json:"is_billable,omitempty"`
	// Whether the project is a fixed-fee project or not.
	IsFixedFee *bool `json:"is_fixed_fee,omitempty"`
	// The method by which the project is invoiced.
	BillBy *string `json:"bill_by,omitempty"`
	// Rate for projects billed by Project Hourly Rate.
	HourlyRate *float64 `json:"hourly_rate,omitempty"`
	// The budget in hours for the project when budgeting by time.
	Budget *float64 `json:"budget,omitempty"`
	// The method by which the project is budgeted.
	BudgetBy *string `json:"budget_by,omitempty"`
	// Whether the project budget is monthly or not.
	BudgetIsMonthly *bool `json:"budget_is_monthly,omitempty"`
	// Whether project managers should be notified when the project goes over budget.
	NotifyWhenOverBudget *bool `json:"notify_when_over_budget,omitempty"`
	// Percentage value used to trigger over budget email alerts.
	OverBudgetNotificationPercentage *float64 `json:"over_budget_notification_percentage,omitempty"`
	// Date of last over budget notification. If none have been sent, this will be null.
	OverBudgetNotificationDate *Date `json:"over_budget_notification_date,omitempty"`
	// Option to show project budget to all employees. Does not apply to Total Project Fee projects.
	ShowBudgetToAll *bool `json:"show_budget_to_all,omitempty"`
	// The monetary budget for the project when budgeting by money.
	CostBudget *float64 `json:"cost_budget,omitempty"`
	// Option for budget of Total Project Fees projects to include tracked expenses.
	CostBudgetIncludeExpenses *bool `json:"cost_budget_include_expenses,omitempty"`
	// The amount you plan to invoice for the project. Only used by fixed-fee projects.
	Fee *float64 `json:"fee,omitempty"`
	// Project notes.
	Notes *string `json:"notes,omitempty"`
	// Date the project was started.
	StartsOn *Date `json:"starts_on,omitempty"`
	// Date the project will end.
	EndsOn *Date `json:"ends_on,omitempty"`
	// Date and time the project was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Date and time the project was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type ProjectList struct {
	Projects []*Project `json:"projects"`

	Pagination
}

func (p Project) String() string {
	return Stringify(p)
}

func (p ProjectList) String() string {
	return Stringify(p)
}

type ProjectListOptions struct {
	// Pass true to only return active projects and false to return inactive projects.
	IsActive bool `url:"is_active,omitempty"`
	// Only return projects belonging to the client with the given ID.
	ClientID int64 `url:"client_id,omitempty"`
	// Only return projects that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

// List returns a list of your projects.
func (s *ProjectService) List(ctx context.Context, opt *ProjectListOptions) (*ProjectList, *http.Response, error) {
	u := "projects"

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	projectList := new(ProjectList)

	resp, err := s.client.Do(ctx, req, &projectList)
	if err != nil {
		return nil, resp, err
	}

	return projectList, resp, nil
}

// Get retrieves the project with the given ID.
func (s *ProjectService) Get(ctx context.Context, projectID int64) (*Project, *http.Response, error) {
	u := fmt.Sprintf("projects/%d", projectID)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	project := new(Project)

	resp, err := s.client.Do(ctx, req, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, nil
}
