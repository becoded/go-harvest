package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Harvest API docs: https://help.getharvest.com/api-v2/expenses-api/expenses/expense-categories/

type ExpenseCategory struct {
	// Unique ID for the expense category.
	ID *int64 `json:"id,omitempty"`
	// The name of the expense category.
	Name *string `json:"name,omitempty"`
	// The unit name of the expense category.
	UnitName *string `json:"unit_name,omitempty"`
	// The unit price of the expense category.
	UnitPrice *float64 `json:"unit_price,omitempty"`
	// Whether the expense category is active or archived.
	IsActive *bool `json:"is_active,omitempty"`
	// Date and time the expense category was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Date and time the expense category was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type ExpenseCategoryList struct {
	ExpenseCategories []*ExpenseCategory `json:"expense_categories"`

	Pagination
}

type ExpenseCategoryRequest struct {
	// required	The name of the expense category.
	Name *string `json:"name,omitempty"`
	// optional	The unit name of the expense category.
	UnitName *string `json:"unit_name,omitempty"`
	// optional	The unit price of the expense category.
	UnitPrice *float64 `json:"unit_price,omitempty"`
	// optional	Whether the expense category is active or archived. Defaults to true.
	IsActive *bool `json:"is_active,omitempty"`
}

func (p ExpenseCategory) String() string {
	return Stringify(p)
}

func (p ExpenseCategoryList) String() string {
	return Stringify(p)
}

type ExpenseCategoryListOptions struct {
	// Pass true to only return active expense categories and false to return inactive expense categories.
	IsActive *bool `url:"is_active,omitempty"`
	// Only return expense categories that have been updated since the given date and time.
	UpdatedSince *time.Time `url:"updated_since,omitempty"`

	ListOptions
}

// ListExpenseCategories returns a list of your expense categories.
func (s *ExpenseService) ListExpenseCategories(
	ctx context.Context,
	opt *ExpenseCategoryListOptions,
) (*ExpenseCategoryList, *http.Response, error) {
	u := "expense_categories"

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	expenseCategoryList := new(ExpenseCategoryList)

	resp, err := s.client.Do(ctx, req, expenseCategoryList)
	if err != nil {
		return nil, resp, err
	}

	return expenseCategoryList, resp, nil
}

// GetExpenseCategory retrieves the expense category with the given ID.
func (s *ExpenseService) GetExpenseCategory(
	ctx context.Context,
	expenseCategoryID int64,
) (*ExpenseCategory, *http.Response, error) {
	u := fmt.Sprintf("expense_categories/%d", expenseCategoryID)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	expenseCategory := new(ExpenseCategory)

	resp, err := s.client.Do(ctx, req, expenseCategory)
	if err != nil {
		return nil, resp, err
	}

	return expenseCategory, resp, nil
}

// CreateExpenseCategory creates a new expense category object.
func (s *ExpenseService) CreateExpenseCategory(
	ctx context.Context,
	data *ExpenseCategoryRequest,
) (*ExpenseCategory, *http.Response, error) {
	u := "expense_categories"

	req, err := s.client.NewRequest(ctx, "POST", u, data)
	if err != nil {
		return nil, nil, err
	}

	expenseCategory := new(ExpenseCategory)

	resp, err := s.client.Do(ctx, req, expenseCategory)
	if err != nil {
		return nil, resp, err
	}

	return expenseCategory, resp, nil
}

// UpdateExpenseCategory updates the specific expense category.
func (s *ExpenseService) UpdateExpenseCategory(
	ctx context.Context,
	expenseCategoryID int64,
	data *ExpenseCategoryRequest,
) (*ExpenseCategory, *http.Response, error) {
	u := fmt.Sprintf("expense_categories/%d", expenseCategoryID)

	req, err := s.client.NewRequest(ctx, "PATCH", u, data)
	if err != nil {
		return nil, nil, err
	}

	expenseCategory := new(ExpenseCategory)

	resp, err := s.client.Do(ctx, req, expenseCategory)
	if err != nil {
		return nil, resp, err
	}

	return expenseCategory, resp, nil
}

// DeleteExpenseCategory deletes an expense category.
func (s *ExpenseService) DeleteExpenseCategory(
	ctx context.Context,
	expenseCategoryID int64,
) (*http.Response, error) {
	u := fmt.Sprintf("expense_categories/%d", expenseCategoryID)

	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
