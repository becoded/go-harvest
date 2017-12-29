package harvest


import (
"context"
"fmt"
"time"
	"net/http"
)

// ProjectService handles communication with the issue related
// methods of the Harvest API.
//
// Harvest API docs: https://help.getharvest.com/api-v2/projects-api/projects/projects/
type ProjectService service


type Project struct {
	Id *int `json:"id,omitempty"` //Unique ID for the project.
	Client *Client `json:"client,omitempty"` //An object containing the project’s client id, name, and currency.
	Name *string `json:"name,omitempty"` //Unique name for the project.
	Code *string `json:"code,omitempty"` //The code associated with the project.
	IsActive *bool `json:"is_active,omitempty"` //Whether the project is active or archived.
	IsBillable *bool `json:"is_billable,omitempty"` //Whether the project is billable or not.
	IsFixedFee *bool `json:"is_fixed_fee,omitempty"` //Whether the project is a fixed-fee project or not.
	BillBy *string `json:"bill_by,omitempty"` //The method by which the project is invoiced.
	HourlyRate *float64 `json:"hourly_rate,omitempty"` //Rate for projects billed by Project Hourly Rate.
	Budget *float64 `json:"budget,omitempty"` //The budget in hours for the project when budgeting by time.
	BudgetBy *string `json:"budget_by,omitempty"` //The method by which the project is budgeted.
	NotifyWhenOverBudget *bool `json:"notify_when_over_budget,omitempty"` //Whether project managers should be notified when the project goes over budget.
	OverBudgetNotificationPercentage *float64 `json:"over_budget_notification_percentage,omitempty"` //Percentage value used to trigger over budget email alerts.
	OverBudgetNotificationDate *time.Time `json:"over_budget_notification_date,omitempty"` //Date of last over budget notification. If none have been sent, this will be null.
	ShowBudgetToAll *bool `json:"show_budget_to_all,omitempty"` //Option to show project budget to all employees. Does not apply to Total Project Fee projects.
	CostBudget *float64 `json:"cost_budget,omitempty"` //The monetary budget for the project when budgeting by money.
	CostBudgetIncludeExpenses *bool `json:"cost_budget_include_expenses,omitempty"` //Option for budget of Total Project Fees projects to include tracked expenses.
	Fee *float64 `json:"fee,omitempty"` //The amount you plan to invoice for the project. Only used by fixed-fee projects.
	Notes *string `json:"notes,omitempty"` //Project notes.
	StartsOn *time.Time `json:"starts_on,omitempty"` //Date the project was started.
	EndsOn *time.Time `json:"ends_on,omitempty"` //Date the project will end.
	CreatedAt *time.Time `json:"created_at,omitempty"` //Date and time the project was created.
	UpdatedAt *time.Time `json:"updated_at,omitempty"` //Date and time the project was last updated.
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
	IsActive	bool `url:"is_active,omitempty"`
	// Only return projects belonging to the client with the given ID.
	ClientId	int `url:"client_id,omitempty"`
	// Only return projects that have been updated since the given date and time.
	UpdatedSince	time.Time `url:"updated_since,omitempty"`

	ListOptions
}

func (s *ProjectService) List(ctx context.Context, opt *ProjectListOptions) (*ProjectList, *http.Response, error) {
	u := "projects"
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
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

func (s *ProjectService) Get(ctx context.Context, projectId string) (*Project, *http.Response, error) {
	u := fmt.Sprintf("projects/%d", projectId)
	req, err := s.client.NewRequest("GET", u, nil)
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