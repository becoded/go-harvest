package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// ExpenseService handles communication with the issue related
// methods of the Harvest API.
//
// Harvest API docs: https://help.getharvest.com/api-v2/expenses-api/expenses/expenses/
type ExpenseService service

type Expense struct {
	// Unique ID for the expense.
	ID *int64 `json:"id,omitempty"`
	// An object containing the expense’s client id, name, and currency.
	Client *Client `json:"client,omitempty"`
	// An object containing the expense’s project id, name, and code.
	Project *Project `json:"project,omitempty"`
	// An object containing the expense’s expense category id, name, unit_price, and unit_name.
	ExpenseCategory *ExpenseCategory `json:"expense_category,omitempty"`
	// An object containing the id and name of the user that recorded the expense.
	User *User `json:"user,omitempty"`
	// A user assignment object of the user that recorded the expense.
	UserAssignment *ProjectUserAssignment `json:"user_assignment,omitempty"`
	// An object containing the expense’s receipt URL and file name.
	Receipt *Receipt `json:"receipt,omitempty"`
	// Once the expense has been invoiced, this field will include the associated invoice’s id and number.
	Invoice *Invoice `json:"invoice,omitempty"`
	// Textual notes used to describe the expense.
	Notes *string `json:"notes,omitempty"`
	// Whether the expense is billable or not.
	Billable *bool `json:"billable,omitempty"`
	// Whether the expense has been approved or closed for some other reason.
	IsClosed *bool `json:"is_closed,omitempty"`
	// Whether the expense has been been invoiced, approved, or the project or person related to the expense is archived.
	IsLocked *bool `json:"is_locked,omitempty"`
	// Whether or not the expense has been marked as invoiced.
	IsBilled *bool `json:"is_billed,omitempty"`
	// An explanation of why the expense has been locked.
	LockedReason *string `json:"locked_reason,omitempty"`
	// Date the expense occurred.
	SpentDate *Date `json:"spent_date,omitempty"`
	// Date and time the expense was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Date and time the expense was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type Receipt struct {
	URL         *string `json:"url,omitempty"`
	FileName    *string `json:"file_name,omitempty"`
	FileSize    *int64  `json:"file_size,omitempty"`
	ContentType *string `json:"content_type,omitempty"`
}

type ExpenseList struct {
	Expenses []*Expense `json:"expenses"`

	Pagination
}

type ExpenseCreateRequest struct {
	// optional	The ID of the user associated with this expense. Defaults to the ID of the currently authenticated user.
	UserID *int64 `json:"user_id,omitempty"`
	// required	The ID of the project associated with this expense.
	ProjectID *int64 `json:"project_id"`
	// required	The ID of the expense category this expense is being tracked against.
	ExpenseCategoryID *int64 `json:"expense_category_id"`
	// required	Date the expense occurred.
	SpentDate *Date `json:"spent_date"`
	// *optional	The quantity of units to use in calculating the total_cost of the expense.
	Units *int64 `json:"units,omitempty"`
	// *optional	The total amount of the expense.
	TotalCost *float64 `json:"total_cost,omitempty"`
	// optional	Textual notes used to describe the expense.
	Notes *string `json:"notes,omitempty"`
	// optional	Whether this expense is billable or not. Defaults to true.
	Billable *bool `json:"billable,omitempty"`
	// TO DO add receipt file
	// optional	A receipt file to attach to the expense.
	// If including a receipt, you must submit a multipart/form-data request.
	// Receipt *file `json:"receipt,omitempty"`
}

type ExpenseUpdateRequest struct {
	// The ID of the project associated with this expense.
	ProjectID *int64 `json:"project_id,omitempty"`
	// The ID of the expense category this expense is being tracked against.
	ExpenseCategoryID *int64 `json:"expense_category_id,omitempty"`
	// Date the expense occurred.
	SpentDate *Date `json:"spent_date,omitempty"`
	// The quantity of units to use in calculating the total_cost of the expense.
	Units *int64 `json:"units,omitempty"`
	// The total amount of the expense.
	TotalCost *float64 `json:"total_cost,omitempty"`
	// Textual notes used to describe the expense.
	Notes *string `json:"notes,omitempty"`
	// Whether this expense is billable or not. Defaults to true.
	Billable *bool `json:"billable,omitempty"`
	// TO DO add receipt file
	// A receipt file to attach to the expense. If including a receipt, you must submit a multipart/form-data request.
	// Receipt *file `json:"receipt,omitempty"`
	// Whether an attached expense receipt should be deleted. Pass true to delete the expense receipt.
	DeleteReceipt *bool `json:"delete_receipt,omitempty"`
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
	// Only return expenses belonging to the user with the given ID.
	UserID *int64 `url:"user_id,omitempty"`
	// Only return expenses belonging to the client with the given ID.
	ClientID *int64 `url:"client_id,omitempty"`
	// Only return expenses belonging to the project with the given ID.
	ProjectID *int64 `url:"project_id,omitempty"`
	// Pass true to only return expenses that have been invoiced and
	// false to return expenses that have not been invoiced.
	IsBilled *bool `url:"is_billed,omitempty"`
	// Only return expenses that have been updated since the given date and time.
	UpdatedSince *time.Time `url:"updated_since,omitempty"`
	// Only return expenses with a spent_date on or after the given date.
	From *Date `url:"from,omitempty"`
	// Only return expenses with a spent_date on or before the given date.
	To *Date `url:"to,omitempty"`

	ListOptions
}

// List returns a list of your expenses.
func (s *ExpenseService) List(ctx context.Context, opt *ExpenseListOptions) (*ExpenseList, *http.Response, error) {
	u := "expenses"

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
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

// Get retrieves the expense with the given ID.
func (s *ExpenseService) Get(ctx context.Context, expenseID int64) (*Expense, *http.Response, error) {
	u := fmt.Sprintf("expenses/%d", expenseID)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
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

// Create creates a new expense object.
func (s *ExpenseService) Create(ctx context.Context, data *ExpenseCreateRequest) (*Expense, *http.Response, error) {
	u := "expenses"

	req, err := s.client.NewRequest(ctx, "POST", u, data)
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

// Update Updates the specific expense by setting the values of the parameters passed.
func (s *ExpenseService) Update(
	ctx context.Context,
	expenseID int64,
	data *ExpenseUpdateRequest,
) (*Expense, *http.Response, error) {
	u := fmt.Sprintf("expenses/%d", expenseID)

	req, err := s.client.NewRequest(ctx, "PATCH", u, data)
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

// Delete deletes an expense.
func (s *ExpenseService) Delete(ctx context.Context, expenseID int64) (*http.Response, error) {
	u := fmt.Sprintf("expenses/%d", expenseID)

	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
