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
	Id        *int64             `json:"id,omitempty"`         // Unique ID for the invoice.
	Client    *Client            `json:"client,omitempty"`     // An object containing invoice’s client id and name.
	LineItems *[]InvoiceLineItem `json:"line_items,omitempty"` // Array of invoice line items.
	Estimate  *Estimate          `json:"estimate,omitempty"`   // An object containing the associated estimate’s id.
	//retainer *object `json:"retainer,omitempty"` // An object containing the associated retainer’s id.
	Creator        *User      `json:"creator,omitempty"`         // An object containing the id and name of the person that created the invoice.
	ClientKey      *string    `json:"client_key,omitempty"`      // Used to build a URL to the public web invoice for your client: https://{ACCOUNT_SUBDOMAIN}.harvestapp.com/client/invoices/abc123456
	Number         *string    `json:"number,omitempty"`          // If no value is set, the number will be automatically generated.
	PurchaseOrder  *string    `json:"purchase_order,omitempty"`  // The purchase order number.
	Amount         *float64   `json:"amount,omitempty"`          // The total amount for the invoice, including any discounts and taxes.
	DueAmount      *float64   `json:"due_amount,omitempty"`      // The total amount due at this time for this invoice.
	Tax            *float64   `json:"tax,omitempty"`             // This percentage is applied to the subtotal, including line items and discounts.
	TaxAmount      *float64   `json:"tax_amount,omitempty"`      // The first amount of tax included, calculated from tax. If no tax is defined, this value will be null.
	Tax2           *float64   `json:"tax2,omitempty"`            // This percentage is applied to the subtotal, including line items and discounts.
	Tax2Amount     *float64   `json:"tax2_amount,omitempty"`     // The amount calculated from tax2.
	Discount       *float64   `json:"discount,omitempty"`        // This percentage is subtracted from the subtotal.
	DiscountAmount *float64   `json:"discount_amount,omitempty"` // The amount calcuated from discount.
	Subject        *string    `json:"subject,omitempty"`         // The invoice subject.
	Notes          *string    `json:"notes,omitempty"`           // Any additional notes included on the invoice.
	Currency       *string    `json:"currency,omitempty"`        // The currency code associated with this invoice.
	PeriodStart    *Date      `json:"period_start,omitempty"`    // Start of the period during which time entries and expenses were added to this invoice.
	PeriodEnd      *Date      `json:"period_end,omitempty"`      // End of the period during which time entries and expenses were added to this invoice.
	IssueDate      *Date      `json:"issue_date,omitempty"`      // Date the invoice was issued.
	DueDate        *Date      `json:"due_date,omitempty"`        // Date the invoice is due.
	SentAt         *time.Time `json:"sent_at,omitempty"`         // Date and time the invoice was sent.
	PaidAt         *time.Time `json:"paid_at,omitempty"`         // Date and time the invoice was paid.
	ClosedAt       *time.Time `json:"closed_at,omitempty"`       // Date and time the invoice was closed.
	CreatedAt      *time.Time `json:"created_at,omitempty"`      // Date and time the invoice was created.
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`      // Date and time the invoice was last updated.
}

type InvoiceLineItem struct {
	Id          *int64   `json:"id,omitempty"`          // Unique ID for the line item.
	Project     *Project `json:"project,omitempty"`     // An object containing the associated project’s id, name, and code.
	Kind        *string  `json:"kind,omitempty"`        // The name of an invoice item category.
	Description *string  `json:"description,omitempty"` // Text description of the line item.
	Quantity    *int16   `json:"quantity,omitempty"`    // The unit quantity of the item.
	UnitPrice   *float64 `json:"unit_price,omitempty"`  // The individual price per unit.
	Amount      *float64 `json:"amount,omitempty"`      // The line item subtotal (quantity * unit_price).
	Taxed       *bool    `json:"taxed,omitempty"`       // Whether the invoice’s tax percentage applies to this line item.
	Taxed2      *bool    `json:"taxed2,omitempty"`      // Whether the invoice’s tax2 percentage applies to this line item.
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
	ClientId int64 `url:"client_id,omitempty"`
	//Only return invoices associated with the project with the given ID.
	ProjectId int64 `url:"project_id,omitempty"`
	// Only return invoices that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

type InvoiceCreateRequest struct {
	ClientId        *int64                          `json:"client_id"`                   // required	The ID of the client this invoice belongs to.
	RetainerId      *int64                          `json:"retainer_id,omitempty"`       // optional	The ID of the retainer associated with this invoice.
	EstimateId      *int64                          `json:"estimate_id,omitempty"`       // optional	The ID of the estimate associated with this invoice.
	Number          *string                         `json:"number,omitempty"`            // optional	If no value is set, the number will be automatically generated.
	PurchaseOrder   *string                         `json:"purchase_order,omitempty"`    // optional	The purchase order number.
	Tax             *float64                        `json:"tax,omitempty"`               // optional	This percentage is applied to the subtotal, including line items and discounts. Example: use 10.0 for 10.0%.
	Tax2            *float64                        `json:"tax2,omitempty"`              // optional	This percentage is applied to the subtotal, including line items and discounts. Example: use 10.0 for 10.0%.
	Discount        *float64                        `json:"discount,omitempty"`          // optional	This percentage is subtracted from the subtotal. Example: use 10.0 for 10.0%.
	Subject         *string                         `json:"subject,omitempty"`           // optional	The invoice subject.
	Notes           *string                         `json:"notes,omitempty"`             // optional	Any additional notes to include on the invoice.
	Currency        *string                         `json:"currency,omitempty"`          // optional	The currency used by the invoice. If not provided, the client’s currency will be used. See a list of supported currencies
	IssueDate       *Date                           `json:"issue_date,omitempty"`        // optional	Date the invoice was issued. Defaults to today’s date.
	DueDate         *Date                           `json:"due_date,omitempty"`          // optional	Date the invoice is due. Defaults to the issue_date.
	LineItems       *[]InvoiceLineItemRequest       `json:"line_items,omitempty"`        // optional	Array of line item parameters
	LineItemsImport *[]InvoiceLineItemImportRequest `json:"line_items_import,omitempty"` // optional	An line items import object
}

type InvoiceLineItemRequest struct {
	Id          *int64   `json:"id,omitempty"`          // Unique ID for the line item.
	ProjectId   *int64   `json:"project_id,omitempty"`  // optional	The ID of the project associated with this line item.
	Kind        *string  `json:"kind"`                  // required	The name of an invoice item category.
	Description *string  `json:"description,omitempty"` // optional	Text description of the line item.
	Quantity    *int64   `json:"quantity,omitempty"`    // optional	The unit quantity of the item. Defaults to 1.
	UnitPrice   *float64 `json:"unit_price"`            // required	The individual price per unit.
	Taxed       *bool    `json:"taxed,omitempty"`       // optional	Whether the invoice’s tax percentage applies to this line item. Defaults to false.
	Taxed2      *bool    `json:"taxed2,omitempty"`      // optional	Whether the invoice’s tax2 percentage applies to this line item. Defaults to false.
	Destroy     *bool    `json:"_destroy,omitempty"`    // optional	DeleteRole an invoice line item
}

type InvoiceLineItemImportRequest struct {
	ProjectIds *[]int64                             `json:"project_ids"`        // required	An array of the client’s project IDs you’d like to include time/expenses from.
	Time       *InvoiceLineItemImportTimeRequest    `json:"time,omitempty"`     // optional	A time import object.
	Expenses   *InvoiceLineItemImportExpenseRequest `json:"expenses,omitempty"` // optional	An expense import object.
}

type InvoiceLineItemImportTimeRequest struct {
	Summary_type *string `json:"summary_type"`   // required	How to summarize the time entries per line item. Options: project, task, people, or detailed.
	From         *Date   `json:"from,omitempty"` // optional	Start date for included time entries. Must be provided if to is present. If neither from or to are provided, all unbilled time entries will be included.
	To           *Date   `json:"to,omitempty"`   // optional	End date for included time entries. Must be provided if from is present. If neither from or to are provided, all unbilled time entries will be included.
}

type InvoiceLineItemImportExpenseRequest struct {
	Summary_type  *string `json:"summary_type"`             // required	How to summarize the expenses per line item. Options: project, category, people, or detailed.
	From          *Date   `json:"from,omitempty"`           // optional	Start date for included expenses. Must be provided if to is present. If neither from or to are provided, all unbilled expenses will be included.
	To            *Date   `json:"to,omitempty"`             // optional	End date for included expenses. Must be provided if from is present. If neither from or to are provided, all unbilled expenses will be included.
	AttachReceipt *bool   `json:"attach_receipt,omitempty"` // optional	If set to true, a PDF containing an expense report with receipts will be attached to the invoice. Defaults to false.
}

type InvoiceUpdateRequest struct {
	ClientId      *int64                    `json:"client_id,omitempty"`      // The ID of the client this invoice belongs to.
	RetainerId    *int64                    `json:"retainer_id,omitempty"`    // The ID of the retainer associated with this invoice.
	EstimateId    *int64                    `json:"estimate_id,omitempty"`    // The ID of the estimate associated with this invoice.
	Number        *string                   `json:"number,omitempty"`         // If no value is set, the number will be automatically generated.
	PurchaseOrder *string                   `json:"purchase_order,omitempty"` // The *purchase `json:"The,omitempty"` // order number.
	Tax           *float64                  `json:"tax,omitempty"`            // This percentage is applied to the subtotal, including line items and discounts.Example: use 10.0 for 10.0%.
	Tax2          *float64                  `json:"tax2,omitempty"`           // This percentage is applied to the subtotal, including line items and discounts.Example: use 10.0 for 10.0%.
	Discount      *float64                  `json:"discount,omitempty"`       // This percentage is subtracted from the subtotal.Example: use 10.0 for 10.0%.
	Subject       *string                   `json:"subject,omitempty"`        // The *invoice `json:"The,omitempty"` // subject.
	Notes         *string                   `json:"notes,omitempty"`          // Any additional notes to include on the invoice.
	Currency      *string                   `json:"currency,omitempty"`       // The currency used by the invoice.If not provided, the client’s currency will be used.See a list of supported currencies
	IssueDate     *Date                     `json:"issue_date,omitempty"`     // Date the invoice was issued.
	DueDate       *Date                     `json:"due_date,omitempty"`       // Date the invoice is due.
	LineItems     *[]InvoiceLineItemRequest `json:"line_items,omitempty"`     // Array of line item parameters
}

func (s *InvoiceService) ListInvoices(ctx context.Context, opt *InvoiceListOptions) (*InvoiceList, *http.Response, error) {
	u := "invoices"
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
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

func (s *InvoiceService) GetInvoice(ctx context.Context, invoiceId int64) (*Invoice, *http.Response, error) {
	u := fmt.Sprintf("invoices/%d", invoiceId)
	req, err := s.client.NewRequest("GET", u, nil)
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

func (s *InvoiceService) CreateInvoice(ctx context.Context, data *InvoiceCreateRequest) (*Invoice, *http.Response, error) {
	u := "invoices"

	req, err := s.client.NewRequest("POST", u, data)
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

func (s *InvoiceService) UpdateInvoice(ctx context.Context, invoiceId int64, data *InvoiceUpdateRequest) (*Invoice, *http.Response, error) {
	u := fmt.Sprintf("invoices/%d", invoiceId)

	req, err := s.client.NewRequest("PATCH", u, data)
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

func (s *InvoiceService) DeleteInvoice(ctx context.Context, invoiceId int64) (*http.Response, error) {
	u := fmt.Sprintf("invoices/%d", invoiceId)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
