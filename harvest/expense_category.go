package harvest

import (
	"context"
	"fmt"
	"time"
	"net/http"
)

// Harvest API docs: https://help.getharvest.com/api-v2/expenses-api/expenses/expense-categories/

type ExpenseCategory struct {
	Id *int64 `json:"id,omitempty"` // Unique ID for the expense category.
	Name *string `json:"name,omitempty"` // The name of the expense category.
	UnitName *string `json:"unit_name,omitempty"` // The unit name of the expense category.
	UnitPrice *float64 `json:"unit_price,omitempty"` // The unit price of the expense category.
	IsActive *bool `json:"is_active,omitempty"` // Whether the expense category is active or archived.
	CreatedAt *time.Time `json:"created_at,omitempty"` // Date and time the expense category was created.
	UpdatedAt *time.Time `json:"updated_at,omitempty"` // Date and time the expense category was last updated.
}

type ExpenseCategoryList struct {
	ExpenseCategories []*ExpenseCategory `json:"invoice_item_categories"`

	Pagination
}

type ExpenseCategoryRequest struct {
	Name *string `json:"name,omitempty"` // required	The name of the expense category.
	UnitName *string `json:"unit_name,omitempty"` // optional	The unit name of the expense category.
	UnitPrice *float64 `json:"unit_price,omitempty"` // optional	The unit price of the expense category.
	IsActive *bool `json:"is_active,omitempty"` // optional	Whether the expense category is active or archived. Defaults to true.
}

func (p ExpenseCategory) String() string {
	return Stringify(p)
}

func (p ExpenseCategoryList) String() string {
	return Stringify(p)
}

type ExpenseCategoryListOptions struct {
	IsActive *bool `url:"is_active,omitempty"` // Pass true to only return active expense categories and false to return inactive expense categories.
	UpdatedSince *time.Time `url:"updated_since,omitempty"` // Only return expense categories that have been updated since the given date and time.

	ListOptions
}

func (s *ExpenseService) ListExpenseCategories(ctx context.Context, opt *ExpenseCategoryListOptions) (*ExpenseCategoryList, *http.Response, error) {
	u := "expense_categories"
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	expenseCategoryList := new(ExpenseCategoryList)
	resp, err := s.client.Do(ctx, req, &expenseCategoryList)
	if err != nil {
		return nil, resp, err
	}

	return expenseCategoryList, resp, nil
}

func (s *ExpenseService) GetExpenseCategory(ctx context.Context, expenseCategoryId int64) (*ExpenseCategory, *http.Response, error) {
	u := fmt.Sprintf("expense_categories/%d", expenseCategoryId)
	req, err := s.client.NewRequest("GET", u, nil)
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

func (s *ExpenseService) CreateExpenseCategory(ctx context.Context, data *ExpenseCategoryRequest) (*ExpenseCategory, *http.Response, error) {
	u := "expense_categories"

	req, err := s.client.NewRequest("POST", u, data)
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

func (s *ExpenseService) UpdateExpenseCategory(ctx context.Context, expenseCategoryId int64, data *ExpenseCategoryRequest) (*ExpenseCategory, *http.Response, error) {
	u := fmt.Sprintf("expense_categories/%d", expenseCategoryId)

	req, err := s.client.NewRequest("PATCH", u, data)
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

func (s *ExpenseService) DeleteExpenseCategory(ctx context.Context, expenseCategoryId int64) (*http.Response, error) {
	u := fmt.Sprintf("expense_categories/%d", expenseCategoryId)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
