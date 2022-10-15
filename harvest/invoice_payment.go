package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Harvest API docs: https://help.getharvest.com/api-v2/invoices-api/invoices/invoice-payments/

type InvoicePayment struct {
	// Unique ID for the payment.
	ID *int64 `json:"id,omitempty"`
	// The amount of the payment.
	Amount *string `json:"amount,omitempty"`
	// Date and time the payment was made.
	PaidAt *time.Time `json:"paid_at,omitempty"`
	// The name of the person who recorded the payment.
	RecordedBy *string `json:"recorded_by,omitempty"`
	// The email of the person who recorded the payment.
	RecordedByEmail *string `json:"recorded_by_email,omitempty"`
	// Any notes associated with the payment.
	Notes *string `json:"notes,omitempty"`
	// Either the card authorization or PayPal transaction ID.
	TransactionID *string `json:"transaction_id,omitempty"`
	// The payment gateway id and name used to process the payment.
	PaymentGateway *PaymentGateway `json:"payment_gateway,omitempty"`
	// Date and time the payment was recorded.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Date and time the payment was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type PaymentGateway struct {
	ID   *int64  `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

type InvoicePaymentList struct {
	InvoicePayments []*InvoicePayment `json:"invoice_payments"`

	Pagination
}

type InvoicePaymentRequest struct {
	// required The amount of the payment.
	Amount *float64 `json:"amount"`
	// optional Date and time the payment was made. Pass either paid_at or paid_date, but not both.
	PaidAt *time.Time `json:"paid_at,omitempty"`
	// optional	Date the payment was made. Pass either paid_at or paid_date, but not both.
	PaidDate *Date `json:"paid_date,omitempty"`
	// optional Any notes to be associated with the payment.
	Notes *string `json:"notes,omitempty"`
}

func (p InvoicePayment) String() string {
	return Stringify(p)
}

func (p InvoicePaymentList) String() string {
	return Stringify(p)
}

func (p PaymentGateway) String() string {
	return Stringify(p)
}

type InvoicePaymentListOptions struct {
	// Only return invoice payments that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

// ListPayments returns a list of payments associate with a given invoice.
func (s *InvoiceService) ListPayments(
	ctx context.Context,
	invoiceID int64,
	opt *InvoicePaymentListOptions,
) (*InvoicePaymentList, *http.Response, error) {
	u := fmt.Sprintf("invoices/%d/payments", invoiceID)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	invoicePaymentList := new(InvoicePaymentList)

	resp, err := s.client.Do(ctx, req, &invoicePaymentList)
	if err != nil {
		return nil, resp, err
	}

	return invoicePaymentList, resp, nil
}

// CreatePayment creates a new invoice payment object.
func (s *InvoiceService) CreatePayment(
	ctx context.Context,
	invoiceID int64,
	data *InvoicePaymentRequest,
) (*InvoicePayment, *http.Response, error) {
	u := fmt.Sprintf("invoices/%d/payments", invoiceID)

	req, err := s.client.NewRequest(ctx, "POST", u, data)
	if err != nil {
		return nil, nil, err
	}

	invoicePayment := new(InvoicePayment)

	resp, err := s.client.Do(ctx, req, invoicePayment)
	if err != nil {
		return nil, resp, err
	}

	return invoicePayment, resp, nil
}

// DeleteInvoicePayment deletes an invoice payment.
func (s *InvoiceService) DeleteInvoicePayment(
	ctx context.Context,
	invoiceID,
	invoicePaymentID int64,
) (*http.Response, error) {
	u := fmt.Sprintf("invoices/%d/payments/%d", invoiceID, invoicePaymentID)

	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
