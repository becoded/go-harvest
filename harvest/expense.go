package harvest

import (
	"context"
	"fmt"
	"time"
	"net/http"
)

// ExpenseService handles communication with the issue related
// methods of the Harvest API.
//
// Harvest API docs: https://help.getharvest.com/api-v2/expenses-api/expenses/expenses/
type ExpenseService service

type Expense struct {
	Id              *int64                 `json:"id,omitempty"`               // Unique ID for the expense.
	Client          *Client                `json:"client,omitempty"`           // An object containing the expense’s client id, name, and currency.
	Project         *Project               `json:"project,omitempty"`          // An object containing the expense’s project id, name, and code.
	ExpenseCategory *ExpenseCategory       `json:"expense_category,omitempty"` // An object containing the expense’s expense category id, name, unit_price, and unit_name.
	User            *User                  `json:"user,omitempty"`             // An object containing the id and name of the user that recorded the expense.
	UserAssignment  *ProjectUserAssignment `json:"user_assignment,omitempty"`  // A user assignment object of the user that recorded the expense.
	Receipt         *Receipt               `json:"receipt,omitempty"`          // An object containing the expense’s receipt URL and file name.
	Invoice         *Invoice               `json:"invoice,omitempty"`          // Once the expense has been invoiced, this field will include the associated invoice’s id and number.
	Notes           *string                `json:"notes,omitempty"`            // Textual notes used to describe the expense.
	Billable        *bool                  `json:"billable,omitempty"`         // Whether the expense is billable or not.
	IsClosed        *bool                  `json:"is_closed,omitempty"`        // Whether the expense has been approved or closed for some other reason.
	IsLocked        *bool                  `json:"is_locked,omitempty"`        // Whether the expense has been been invoiced, approved, or the project or person related to the expense is archived.
	IsBilled        *bool                  `json:"is_billed,omitempty"`        // Whether or not the expense has been marked as invoiced.
	LockedReason    *string                `json:"locked_reason,omitempty"`    // An explanation of why the expense has been locked.
	SpentDate       *Date                  `json:"spent_date,omitempty"`       // Date the expense occurred.
	CreatedAt       *time.Time             `json:"created_at,omitempty"`       // Date and time the expense was created.
	UpdatedAt       *time.Time             `json:"updated_at,omitempty"`       // Date and time the expense was last updated.
}

type Receipt struct {
	Url         *string `json:"url,omitempty"`
	FileName    *string `json:"file_name,omitempty"`
	FileSize    *int64  `json:"file_size,omitempty"`
	ContentType *string `json:"content_type,omitempty"`
}

type ExpenseList struct {
	Expenses []*Expense `json:"expenses"`

	Pagination
}

type ExpenseCreateRequest struct {
	UserId *int64 `json:"user_id,omitempty"` // optional	The ID of the user associated with this expense. Defaults to the ID of the currently authenticated user.
	ProjectId *int64 `json:"project_id"` // required	The ID of the project associated with this expense.
	ExpenseCategoryId *int64 `json:"expense_category_id"` // required	The ID of the expense category this expense is being tracked against.
	SpentDate *Date `json:"spent_date"` // required	Date the expense occurred.
	Units *int64 `json:"units,omitempty"` // *optional	The quantity of units to use in calculating the total_cost of the expense.
	TotalCost *float64 `json:"total_cost,omitempty"` // *optional	The total amount of the expense.
	Notes *string `json:"notes,omitempty"` // optional	Textual notes used to describe the expense.
	Billable *bool `json:"billable,omitempty"` // optional	Whether this expense is billable or not. Defaults to true.
	// TODO add receipt file
	// Receipt *file `json:"receipt,omitempty"` // optional	A receipt file to attach to the expense. If including a receipt, you must submit a multipart/form-data request.
}

type ExpenseUpdateRequest struct {
	ProjectId *int64 `json:"project_id,omitempty"` // The ID of the project associated with this expense.
	ExpenseCategoryId *int64 `json:"expense_category_id,omitempty"` // The ID of the expense category this expense is being tracked against.
	Spent_date *Date `json:"spent_date,omitempty"` // Date the expense occurred.
	Units *int64 `json:"units,omitempty"` // The quantity of units to use in calculating the total_cost of the expense.
	TotalCost *float64 `json:"total_cost,omitempty"` // The total amount of the expense.
	Notes *string `json:"notes,omitempty"` // Textual notes used to describe the expense.
	Billable *bool `json:"billable,omitempty"` // Whether this expense is billable or not. Defaults to true.
	// TODO add receipt file
	// Receipt *file `json:"receipt,omitempty"` // A receipt file to attach to the expense. If including a receipt, you must submit a multipart/form-data request.
	DeleteReceipt *bool `json:"delete_receipt,omitempty"` // Whether an attached expense receipt should be deleted. Pass true to delete the expense receipt.
}

func (p Expense) String() string {
	return Stringify(p)
}

func (p Receipt) String() string {
	return Stringify(p)
}

func (p ExpenseList) String() string {
	return Stringify(p)
}

type ExpenseListOptions struct {
	UserId       *int64     `url:"user_id,omitempty"`       // Only return expenses belonging to the user with the given ID.
	ClientId     *int64     `url:"client_id,omitempty"`     // Only return expenses belonging to the client with the given ID.
	ProjectId    *int64     `url:"project_id,omitempty"`    // Only return expenses belonging to the project with the given ID.
	IsBilled     *bool      `url:"is_billed,omitempty"`     // Pass true to only return expenses that have been invoiced and false to return expenses that have not been invoiced.
	UpdatedSince *time.Time `url:"updated_since,omitempty"` // Only return expenses that have been updated since the given date and time.
	From         *Date      `url:"from,omitempty"`          // Only return expenses with a spent_date on or after the given date.
	To           *Date      `url:"to,omitempty"`            // Only return expenses with a spent_date on or before the given date.

	ListOptions
}

func (s *ExpenseService) List(ctx context.Context, opt *ExpenseListOptions) (*ExpenseList, *http.Response, error) {
	u := "expenses"
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	expenseList := new(ExpenseList)
	resp, err := s.client.Do(ctx, req, &expenseList)
	if err != nil {
		return nil, resp, err
	}

	return expenseList, resp, nil
}

func (s *ExpenseService) Get(ctx context.Context, expenseId int64) (*Expense, *http.Response, error) {
	u := fmt.Sprintf("expenses/%d", expenseId)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	expense := new(Expense)
	resp, err := s.client.Do(ctx, req, expense)
	if err != nil {
		return nil, resp, err
	}

	return expense, resp, nil
}

func (s *ExpenseService) CreateExpense(ctx context.Context, data *ExpenseCreateRequest) (*Expense, *http.Response, error) {
	u := "expenses"

	req, err := s.client.NewRequest("POST", u, data)
	if err != nil {
		return nil, nil, err
	}

	expense := new(Expense)
	resp, err := s.client.Do(ctx, req, expense)
	if err != nil {
		return nil, resp, err
	}

	return expense, resp, nil
}

func (s *ExpenseService) UpdateExpense(ctx context.Context, expenseId int64, data *ExpenseUpdateRequest) (*Expense, *http.Response, error) {
	u := fmt.Sprintf("expenses/%d", expenseId)

	req, err := s.client.NewRequest("PATCH", u, data)
	if err != nil {
		return nil, nil, err
	}

	expense := new(Expense)
	resp, err := s.client.Do(ctx, req, expense)
	if err != nil {
		return nil, resp, err
	}

	return expense, resp, nil
}

func (s *ExpenseService) DeleteExpense(ctx context.Context, expenseId int64) (*http.Response, error) {
	u := fmt.Sprintf("expenses/%d", expenseId)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
