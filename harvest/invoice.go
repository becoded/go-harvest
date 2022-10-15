package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// InvoiceService handles communication with the issue related
// methods of the Harvest API.
//
// Harvest API docs: https://help.getharvest.com/api-v2/invoices-api/invoices/invoices/
type InvoiceService service

type Invoice struct {
	// Unique ID for the invoice.
	ID *int64 `json:"id,omitempty"`
	// An object containing invoice’s client id and name.
	Client *Client `json:"client,omitempty"`
	// Array of invoice line items.
	LineItems *[]InvoiceLineItem `json:"line_items,omitempty"`
	// An object containing the associated estimate’s id.
	Estimate *Estimate `json:"estimate,omitempty"`
	// retainer *object `json:"retainer,omitempty"` // An object containing the associated retainer’s id.
	// An object containing the id and name of the person that created the invoice.
	Creator *User `json:"creator,omitempty"`
	// Used to build a URL to the public web invoice for your client:
	// https://{ACCOUNT_SUBDOMAIN}.harvestapp.com/client/invoices/abc123456
	ClientKey *string `json:"client_key,omitempty"`
	// If no value is set, the number will be automatically generated.
	Number *string `json:"number,omitempty"`
	// The purchase order number.
	PurchaseOrder *string `json:"purchase_order,omitempty"`
	// The total amount for the invoice, including any discounts and taxes.
	Amount *float64 `json:"amount,omitempty"`
	// The total amount due at this time for this invoice.
	DueAmount *float64 `json:"due_amount,omitempty"`
	// This percentage is applied to the subtotal, including line items and discounts.
	Tax *float64 `json:"tax,omitempty"`
	// The first amount of tax included, calculated from tax. If no tax is defined, this value will be null.
	TaxAmount *float64 `json:"tax_amount,omitempty"`
	// This percentage is applied to the subtotal, including line items and discounts.
	Tax2 *float64 `json:"tax2,omitempty"`
	// The amount calculated from tax2.
	Tax2Amount *float64 `json:"tax2_amount,omitempty"`
	// This percentage is subtracted from the subtotal.
	Discount *float64 `json:"discount,omitempty"`
	// The amount calcuated from discount.
	DiscountAmount *float64 `json:"discount_amount,omitempty"`
	// The invoice subject.
	Subject *string `json:"subject,omitempty"`
	// Any additional notes included on the invoice.
	Notes *string `json:"notes,omitempty"`
	// The currency code associated with this invoice.
	Currency *string `json:"currency,omitempty"`
	// Start of the period during which time entries and expenses were added to this invoice.
	PeriodStart *Date `json:"period_start,omitempty"`
	// End of the period during which time entries and expenses were added to this invoice.
	PeriodEnd *Date `json:"period_end,omitempty"`
	// Date the invoice was issued.
	IssueDate *Date `json:"issue_date,omitempty"`
	// Date the invoice is due.
	DueDate *Date `json:"due_date,omitempty"`
	// Date and time the invoice was sent.
	SentAt *time.Time `json:"sent_at,omitempty"`
	// Date and time the invoice was paid.
	PaidAt *time.Time `json:"paid_at,omitempty"`
	// Date and time the invoice was closed.
	ClosedAt *time.Time `json:"closed_at,omitempty"`
	// Date and time the invoice was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Date and time the invoice was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type InvoiceLineItem struct {
	// Unique ID for the line item.
	ID *int64 `json:"id,omitempty"`
	// An object containing the associated project’s id, name, and code.
	Project *Project `json:"project,omitempty"`
	// The name of an invoice item category.
	Kind *string `json:"kind,omitempty"`
	// Text description of the line item.
	Description *string `json:"description,omitempty"`
	// The unit quantity of the item.
	Quantity *float64 `json:"quantity,omitempty"`
	// The individual price per unit.
	UnitPrice *float64 `json:"unit_price,omitempty"`
	// The line item subtotal (quantity * unit_price).
	Amount *float64 `json:"amount,omitempty"`
	// Whether the invoice’s tax percentage applies to this line item.
	Taxed *bool `json:"taxed,omitempty"`
	// Whether the invoice’s tax2 percentage applies to this line item.
	Taxed2 *bool `json:"taxed2,omitempty"`
}

type InvoiceList struct {
	Invoices []*Invoice `json:"invoices"`

	Pagination
}

func (p Invoice) String() string {
	return Stringify(p)
}

func (p InvoiceList) String() string {
	return Stringify(p)
}

type InvoiceListOptions struct {
	// Only return invoices belonging to the client with the given ID.
	ClientID int64 `url:"client_id,omitempty"`
	// Only return invoices associated with the project with the given ID.
	ProjectID int64 `url:"project_id,omitempty"`
	// Only return invoices that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

type InvoiceCreateRequest struct {
	// required	The ID of the client this invoice belongs to.
	ClientID *int64 `json:"client_id"`
	// optional	The ID of the retainer associated with this invoice.
	RetainerID *int64 `json:"retainer_id,omitempty"`
	// optional	The ID of the estimate associated with this invoice.
	EstimateID *int64 `json:"estimate_id,omitempty"`
	// optional	If no value is set, the number will be automatically generated.
	Number *string `json:"number,omitempty"`
	// optional	The purchase order number.
	PurchaseOrder *string `json:"purchase_order,omitempty"`
	// optional	This percentage is applied to the subtotal, including line items and discounts.
	// Example: use 10.0 for 10.0%.
	Tax *float64 `json:"tax,omitempty"`
	// optional	This percentage is applied to the subtotal, including line items and discounts.
	// Example: use 10.0 for 10.0%.
	Tax2 *float64 `json:"tax2,omitempty"`
	// optional	This percentage is subtracted from the subtotal. Example: use 10.0 for 10.0%.
	Discount *float64 `json:"discount,omitempty"`
	// optional	The invoice subject.
	Subject *string `json:"subject,omitempty"`
	// optional	Any additional notes to include on the invoice.
	Notes *string `json:"notes,omitempty"`
	// optional	The currency used by the invoice.
	// If not provided, the client’s currency will be used. See a list of supported currencies
	Currency *string `json:"currency,omitempty"`
	// optional	Date the invoice was issued. Defaults to today’s date.
	IssueDate *Date `json:"issue_date,omitempty"`
	// optional	Date the invoice is due. Defaults to the issue_date.
	DueDate *Date `json:"due_date,omitempty"`
	// optional	Array of line item parameters
	LineItems *[]InvoiceLineItemRequest `json:"line_items,omitempty"`
	// optional	An line items import object
	LineItemsImport *[]InvoiceLineItemImportRequest `json:"line_items_import,omitempty"`
}

type InvoiceLineItemRequest struct {
	// Unique ID for the line item.
	ID *int64 `json:"id,omitempty"`
	// optional	The ID of the project associated with this line item.
	ProjectID *int64 `json:"project_id,omitempty"`
	// required	The name of an invoice item category.
	Kind *string `json:"kind"`
	// optional	Text description of the line item.
	Description *string `json:"description,omitempty"`
	// optional	The unit quantity of the item. Defaults to 1.
	Quantity *int64 `json:"quantity,omitempty"`
	// required	The individual price per unit.
	UnitPrice *float64 `json:"unit_price"`
	// optional	Whether the invoice’s tax percentage applies to this line item. Defaults to false.
	Taxed *bool `json:"taxed,omitempty"`
	// optional	Whether the invoice’s tax2 percentage applies to this line item. Defaults to false.
	Taxed2 *bool `json:"taxed2,omitempty"`
	// optional	Delete an invoice line item
	Destroy *bool `json:"_destroy,omitempty"`
}

type InvoiceLineItemImportRequest struct {
	// required	An array of the client’s project IDs you’d like to include time/expenses from.
	ProjectIds *[]int64 `json:"project_ids"`
	// optional	A time import object.
	Time *InvoiceLineItemImportTimeRequest `json:"time,omitempty"`
	// optional	An expense import object.
	Expenses *InvoiceLineItemImportExpenseRequest `json:"expenses,omitempty"`
}

type InvoiceLineItemImportTimeRequest struct {
	// required	How to summarize the time entries per line item. Options: project, task, people, or detailed.
	SummaryType *string `json:"summary_type"`
	// optional	Start date for included time entries.
	// Must be provided if to is present. If neither from or to are provided,
	// all unbilled time entries will be included.
	From *Date `json:"from,omitempty"`
	// optional	End date for included time entries.
	// Must be provided if from is present. If neither from or to are provided,
	// all unbilled time entries will be included.
	To *Date `json:"to,omitempty"`
}

type InvoiceLineItemImportExpenseRequest struct {
	// required	How to summarize the expenses per line item. Options: project, category, people, or detailed.
	SummaryType *string `json:"summary_type"`
	// optional	Start date for included expenses.
	// Must be provided if to is present. If neither from or to are provided, all unbilled expenses will be included.
	From *Date `json:"from,omitempty"`
	// optional	End date for included expenses.
	// Must be provided if from is present. If neither from or to are provided, all unbilled expenses will be included.
	To *Date `json:"to,omitempty"`
	// optional	If set to true, a PDF containing an expense report with receipts will be attached to the invoice.
	// Defaults to false.
	AttachReceipt *bool `json:"attach_receipt,omitempty"`
}

type InvoiceUpdateRequest struct {
	// The ID of the client this invoice belongs to.
	ClientID *int64 `json:"client_id,omitempty"`
	// The ID of the retainer associated with this invoice.
	RetainerID *int64 `json:"retainer_id,omitempty"`
	// The ID of the estimate associated with this invoice.
	EstimateID *int64 `json:"estimate_id,omitempty"`
	// If no value is set, the number will be automatically generated.
	Number *string `json:"number,omitempty"`
	// The *purchase `json:"The,omitempty"` // order number.
	PurchaseOrder *string `json:"purchase_order,omitempty"`
	// This percentage is applied to the subtotal, including line items and discounts.Example: use 10.0 for 10.0%.
	Tax *float64 `json:"tax,omitempty"`
	// This percentage is applied to the subtotal, including line items and discounts.Example: use 10.0 for 10.0%.
	Tax2 *float64 `json:"tax2,omitempty"`
	// This percentage is subtracted from the subtotal.Example: use 10.0 for 10.0%.
	Discount *float64 `json:"discount,omitempty"`
	// The *invoice `json:"The,omitempty"` // subject.
	Subject *string `json:"subject,omitempty"`
	// Any additional notes to include on the invoice.
	Notes *string `json:"notes,omitempty"`
	// The currency used by the invoice.If not provided,
	// the client’s currency will be used.See a list of supported currencies
	Currency *string `json:"currency,omitempty"`
	// Date the invoice was issued.
	IssueDate *Date `json:"issue_date,omitempty"`
	// Date the invoice is due.
	DueDate *Date `json:"due_date,omitempty"`
	// Array of line item parameters
	LineItems *[]InvoiceLineItemRequest `json:"line_items,omitempty"`
}

// List returns a list of your invoices.
func (s *InvoiceService) List(ctx context.Context, opt *InvoiceListOptions) (*InvoiceList, *http.Response, error) {
	u := "invoices"

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	invoiceList := new(InvoiceList)

	resp, err := s.client.Do(ctx, req, &invoiceList)
	if err != nil {
		return nil, resp, err
	}

	return invoiceList, resp, nil
}

// Get retrieves the invoice with the given ID.
func (s *InvoiceService) Get(ctx context.Context, invoiceID int64) (*Invoice, *http.Response, error) {
	u := fmt.Sprintf("invoices/%d", invoiceID)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	invoice := new(Invoice)

	resp, err := s.client.Do(ctx, req, invoice)
	if err != nil {
		return nil, resp, err
	}

	return invoice, resp, nil
}

// Create creates a new invoice object.
func (s *InvoiceService) Create(ctx context.Context, data *InvoiceCreateRequest) (*Invoice, *http.Response, error) {
	u := "invoices"

	req, err := s.client.NewRequest(ctx, "POST", u, data)
	if err != nil {
		return nil, nil, err
	}

	invoice := new(Invoice)

	resp, err := s.client.Do(ctx, req, invoice)
	if err != nil {
		return nil, resp, err
	}

	return invoice, resp, nil
}

// Update updates the specific invoice.
func (s *InvoiceService) Update(
	ctx context.Context,
	invoiceID int64,
	data *InvoiceUpdateRequest,
) (*Invoice, *http.Response, error) {
	u := fmt.Sprintf("invoices/%d", invoiceID)

	req, err := s.client.NewRequest(ctx, "PATCH", u, data)
	if err != nil {
		return nil, nil, err
	}

	invoice := new(Invoice)

	resp, err := s.client.Do(ctx, req, invoice)
	if err != nil {
		return nil, resp, err
	}

	return invoice, resp, nil
}

// Delete deletes an invoice.
func (s *InvoiceService) Delete(ctx context.Context, invoiceID int64) (*http.Response, error) {
	u := fmt.Sprintf("invoices/%d", invoiceID)

	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
