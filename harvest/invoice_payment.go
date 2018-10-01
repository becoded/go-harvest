package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Harvest API docs: https://help.getharvest.com/api-v2/invoices-api/invoices/invoice-payments/

type InvoicePayment struct {
	Id              *int64          `json:"id,omitempty"`                // Unique ID for the payment.
	Amount          *string         `json:"amount,omitempty"`            // The amount of the payment.
	PaidAt          *time.Time      `json:"paid_at,omitempty"`           // Date and time the payment was made.
	RecordedBy      *string         `json:"recorded_by,omitempty"`       // The name of the person who recorded the payment.
	RecordedByEmail *string         `json:"recorded_by_email,omitempty"` // The email of the person who recorded the payment.
	Notes           *string         `json:"notes,omitempty"`             // Any notes associated with the payment.
	TransactionId   *string         `json:"transaction_id,omitempty"`    // Either the card authorization or PayPal transaction ID.
	PaymentGateway  *PaymentGateway `json:"payment_gateway,omitempty"`   // The payment gateway id and name used to process the payment.
	CreatedAt       *time.Time      `json:"created_at,omitempty"`        // Date and time the payment was recorded.
	UpdatedAt       *time.Time      `json:"updated_at,omitempty"`        // Date and time the payment was last updated.
}

type PaymentGateway struct {
	Id   *int64  `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

type InvoicePaymentList struct {
	InvoicePayments []*InvoicePayment `json:"invoice_payments"`

	Pagination
}

type InvoicePaymentRequest struct {
	Amount   *float64   `json:"amount"`              // required The amount of the payment.
	PaidAt   *time.Time `json:"paid_at,omitempty"`   // optional Date and time the payment was made. Pass either paid_at or paid_date, but not both.
	PaidDate *Date      `json:"paid_date,omitempty"` // optional	Date the payment was made. Pass either paid_at or paid_date, but not both.
	Notes    *string    `json:"notes,omitempty"`     // optional Any notes to be associated with the payment.
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

func (s *InvoiceService) ListPayments(ctx context.Context, invoiceId int64, opt *InvoicePaymentListOptions) (*InvoicePaymentList, *http.Response, error) {
	u := fmt.Sprintf("invoices/%d/payments", invoiceId)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
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

func (s *InvoiceService) CreatePayment(ctx context.Context, invoiceId int64, data *InvoicePaymentRequest) (*InvoicePayment, *http.Response, error) {
	u := fmt.Sprintf("invoices/%d/payments", invoiceId)

	req, err := s.client.NewRequest("POST", u, data)
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

func (s *InvoiceService) DeleteInvoicePayment(ctx context.Context, invoiceId, invoicePaymentId int64) (*http.Response, error) {
	u := fmt.Sprintf("invoices/%d/payments/%d", invoiceId, invoicePaymentId)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
