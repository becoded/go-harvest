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
	Id *int64 `json:"id,omitempty"` // Unique ID for the expense.
	Client *Client `json:"client,omitempty"` // An object containing the expense’s client id, name, and currency.
	Project *Project `json:"project,omitempty"` // An object containing the expense’s project id, name, and code.
	ExpenseCategory *ExpenseCategory `json:"expense_category,omitempty"` // An object containing the expense’s expense category id, name, unit_price, and unit_name.
	User *User `json:"user,omitempty"` // An object containing the id and name of the user that recorded the expense.
	UserAssignment *ProjectUserAssignment `json:"user_assignment,omitempty"` // A user assignment object of the user that recorded the expense.
	Receipt *Receipt `json:"receipt,omitempty"` // An object containing the expense’s receipt URL and file name.
	Invoice *Invoice `json:"invoice,omitempty"` // Once the expense has been invoiced, this field will include the associated invoice’s id and number.
	Notes *string `json:"notes,omitempty"` // Textual notes used to describe the expense.
	Billable *bool `json:"billable,omitempty"` // Whether the expense is billable or not.
	IsClosed *bool `json:"is_closed,omitempty"` // Whether the expense has been approved or closed for some other reason.
	IsLocked *bool `json:"is_locked,omitempty"` // Whether the expense has been been invoiced, approved, or the project or person related to the expense is archived.
	IsBilled *bool `json:"is_billed,omitempty"` // Whether or not the expense has been marked as invoiced.
	LockedReason *string `json:"locked_reason,omitempty"` // An explanation of why the expense has been locked.
	SpentDate *Date `json:"spent_date,omitempty"` // Date the expense occurred.
	CreatedAt *time.Time `json:"created_at,omitempty"` // Date and time the expense was created.
	UpdatedAt *time.Time `json:"updated_at,omitempty"` // Date and time the expense was last updated.
}

type Receipt struct {
	Url *string `json:"url,omitempty"`
	FileName *string `json:"file_name,omitempty"`
	FileSize *int64 `json:"file_size,omitempty"`
	ContentType *string `json:"content_type,omitempty"`
}

type ExpenseList struct {
	Expenses []*Expense `json:"expenses"`

	Pagination
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
	UserId *int64 `url:"user_id,omitempty"` // Only return expenses belonging to the user with the given ID.
	ClientId *int64 `url:"client_id,omitempty"` // Only return expenses belonging to the client with the given ID.
	ProjectId *int64 `url:"project_id,omitempty"` // Only return expenses belonging to the project with the given ID.
	IsBilled *bool `url:"is_billed,omitempty"` // Pass true to only return expenses that have been invoiced and false to return expenses that have not been invoiced.
	UpdatedSince *time.Time `url:"updated_since,omitempty"` // Only return expenses that have been updated since the given date and time.
	From *Date `url:"from,omitempty"` // Only return expenses with a spent_date on or after the given date.
	To *Date `url:"to,omitempty"` // Only return expenses with a spent_date on or before the given date.

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
