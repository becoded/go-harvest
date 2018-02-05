package harvest

import (
	"context"
	"fmt"
	"time"
	"net/http"
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

func (s *InvoiceService) List(ctx context.Context, opt *InvoiceListOptions) (*InvoiceList, *http.Response, error) {
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

func (s *InvoiceService) Get(ctx context.Context, invoiceId int64) (*Invoice, *http.Response, error) {
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
