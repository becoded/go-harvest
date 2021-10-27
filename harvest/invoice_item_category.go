package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Harvest API docs: https://help.getharvest.com/api-v2/invoices-api/invoices/invoice-item-categories/

type InvoiceItemCategory struct {
	// Unique ID for the invoice item category.
	ID *int64 `json:"id,omitempty"`
	// The name of the invoice item category.
	Name *string `json:"name,omitempty"`
	// Whether this invoice item category is used for billable hours when generating an invoice.
	UseAsService *bool `json:"use_as_service,omitempty"`
	// Whether this invoice item category is used for expenses when generating an invoice.
	UseAsExpense *bool `json:"use_as_expense,omitempty"`
	// Date and time the invoice item category was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Date and time the invoice item category was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type InvoiceItemCategoryList struct {
	InvoiceItemCategories []*InvoiceItemCategory `json:"invoice_item_categories"`

	Pagination
}

type InvoiceItemCategoryRequest struct {
	// required	The name of the invoice item category.
	Name *string `json:"name,omitempty"`
}

func (p InvoiceItemCategory) String() string {
	return Stringify(p)
}

func (p InvoiceItemCategoryList) String() string {
	return Stringify(p)
}

type InvoiceItemCategoryListOptions struct {
	// Only return invoice item categories that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

func (s *InvoiceService) ListItemCategories(
	ctx context.Context,
	opt *InvoiceItemCategoryListOptions,
) (*InvoiceItemCategoryList, *http.Response, error) {
	u := "invoice_item_categories"

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	invoiceItemCategoryList := new(InvoiceItemCategoryList)

	resp, err := s.client.Do(ctx, req, &invoiceItemCategoryList)
	if err != nil {
		return nil, resp, err
	}

	return invoiceItemCategoryList, resp, nil
}

func (s *InvoiceService) GetItemCategory(
	ctx context.Context,
	invoiceItemCategoryID int64,
) (*InvoiceItemCategory, *http.Response, error) {
	u := fmt.Sprintf("invoice_item_categories/%d", invoiceItemCategoryID)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	invoiceItemCategory := new(InvoiceItemCategory)

	resp, err := s.client.Do(ctx, req, invoiceItemCategory)
	if err != nil {
		return nil, resp, err
	}

	return invoiceItemCategory, resp, nil
}

func (s *InvoiceService) CreateItemCategory(
	ctx context.Context,
	data *InvoiceItemCategoryRequest,
) (*InvoiceItemCategory, *http.Response, error) {
	u := "invoice_item_categories"

	req, err := s.client.NewRequest(ctx, "POST", u, data)
	if err != nil {
		return nil, nil, err
	}

	invoiceItemCategory := new(InvoiceItemCategory)

	resp, err := s.client.Do(ctx, req, invoiceItemCategory)
	if err != nil {
		return nil, resp, err
	}

	return invoiceItemCategory, resp, nil
}

func (s *InvoiceService) UpdateItemCategory(
	ctx context.Context,
	invoiceItemCategoryID int64,
	data *InvoiceItemCategoryRequest,
) (*InvoiceItemCategory, *http.Response, error) {
	u := fmt.Sprintf("invoice_item_categories/%d", invoiceItemCategoryID)

	req, err := s.client.NewRequest(ctx, "PATCH", u, data)
	if err != nil {
		return nil, nil, err
	}

	invoiceItemCategory := new(InvoiceItemCategory)

	resp, err := s.client.Do(ctx, req, invoiceItemCategory)
	if err != nil {
		return nil, resp, err
	}

	return invoiceItemCategory, resp, nil
}

func (s *InvoiceService) DeleteItemCategory(ctx context.Context, invoiceItemCategoryID int64) (*http.Response, error) {
	u := fmt.Sprintf("invoice_item_categories/%d", invoiceItemCategoryID)

	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
