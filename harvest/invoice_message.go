package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Harvest API docs: https://help.getharvest.com/api-v2/invoices-api/invoices/invoice-messages/

type InvoiceMessage struct {
	// Unique ID for the message.
	ID *int64 `json:"id,omitempty"`
	// Name of the user that created the message.
	SentBy *string `json:"sent_by,omitempty"`
	// Email of the user that created the message.
	SentByEmail *string `json:"sent_by_email,omitempty"`
	// Name of the user that the message was sent from.
	SentFrom *string `json:"sent_from,omitempty"`
	// Email of the user that message was sent from.
	SentFromEmail *string `json:"sent_from_email,omitempty"`
	// Array of invoice message recipients.
	Recipients *[]InvoiceMessageRecipient `json:"recipients,omitempty"`
	// The message subject.
	Subject *string `json:"subject,omitempty"`
	// The message body.
	Body *string `json:"body,omitempty"`
	// Whether to include a link to the client invoice in the message body. Not used when thank_you is true.
	IncludeLinkToClientInvoice *bool `json:"include_link_to_client_invoice,omitempty"`
	// Whether to attach the invoice PDF to the message email.
	AttachPdf *bool `json:"attach_pdf,omitempty"`
	// Whether to email a copy of the message to the current user.
	SendMeACopy *bool `json:"send_me_a_copy,omitempty"`
	// Whether this is a thank you message.
	ThankYou *bool `json:"thank_you,omitempty"`
	// The type of invoice event that occurred with the message: send, close, draft, re-open, or view.
	EventType *bool `json:"event_type,omitempty"`
	// Whether this is a reminder message.
	Reminder *bool `json:"reminder,omitempty"`
	// The date the reminder email will be sent.
	SendReminderOn *Date `json:"send_reminder_on,omitempty"`
	// Date and time the message was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Date and time the message was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type InvoiceMessageRecipient struct {
	// Name of the message recipient.
	Name *string `json:"name,omitempty"`
	// Email of the message recipient.
	Email *string `json:"email"`
}

type InvoiceMessageList struct {
	Invoices []*Invoice `json:"invoices"`

	Pagination
}

func (p InvoiceMessage) String() string {
	return Stringify(p)
}

func (p InvoiceMessageList) String() string {
	return Stringify(p)
}

type InvoiceMessageListOptions struct {
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

type EventTypeRequest struct {
	EventType string `json:"event_type"`
}

type InvoiceMessageCreateRequest struct {
	// required	Array of recipient parameters. See below for details.
	Recipients *[]InvoiceMessageRecipient `json:"recipients"`
	// optional	The message subject.
	Subject *string `json:"subject,omitempty"`
	// optional	The message body.
	Body *string `json:"body,omitempty"`
	// optional	If set to true, a link to the client invoice URL will be included in the message email.
	// Defaults to false. Ignored when thank_you is set to true.
	IncludeLinkToClientInvoice *bool `json:"include_link_to_client_invoice,omitempty"`
	// optional	If set to true, a PDF of the invoice will be attached to the message email. Defaults to false.
	AttachPdf *bool `json:"attach_pdf,omitempty"`
	// optional	If set to true, a copy of the message email will be sent to the current user. Defaults to false.
	SendMeACopy *bool `json:"send_me_a_copy,omitempty"`
	// optional	If set to true, a thank you message email will be sent. Defaults to false.
	ThankYou *bool `json:"thank_you,omitempty"`
	// optional	If provided, runs an event against the invoice. Options: close, draft, re-open, or send.
	EventType *bool `json:"event_type,omitempty"`
}

func (s *InvoiceService) ListInvoiceMessages(
	ctx context.Context,
	invoiceID int64,
	opt *InvoiceMessageListOptions,
) (*InvoiceList, *http.Response, error) {
	u := fmt.Sprintf("invoices/%d/messages", invoiceID)

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

func (s *InvoiceService) CreateInvoiceMessage(
	ctx context.Context,
	invoiceID int64,
	data *InvoiceMessageCreateRequest,
) (*InvoiceMessage, *http.Response, error) {
	u := fmt.Sprintf("invoices/%d/messages", invoiceID)

	req, err := s.client.NewRequest(ctx, "POST", u, data)
	if err != nil {
		return nil, nil, err
	}

	invoiceMessage := new(InvoiceMessage)

	resp, err := s.client.Do(ctx, req, invoiceMessage)
	if err != nil {
		return nil, resp, err
	}

	return invoiceMessage, resp, nil
}

func (s *InvoiceService) DeleteInvoiceMessage(
	ctx context.Context,
	invoiceID,
	invoiceMessageID int64,
) (*http.Response, error) {
	u := fmt.Sprintf("invoices/%d/messages/%d", invoiceID, invoiceMessageID)

	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *InvoiceService) MarkAsSent(
	ctx context.Context,
	invoiceID int64,
) (*InvoiceMessage, *http.Response, error) {
	return s.SendEvent(ctx, invoiceID, &EventTypeRequest{EventType: "send"})
}

func (s *InvoiceService) MarkAsDraft(ctx context.Context, invoiceID int64) (*InvoiceMessage, *http.Response, error) {
	return s.SendEvent(ctx, invoiceID, &EventTypeRequest{EventType: "draft"})
}

func (s *InvoiceService) MarkAsClosed(ctx context.Context, invoiceID int64) (*InvoiceMessage, *http.Response, error) {
	return s.SendEvent(ctx, invoiceID, &EventTypeRequest{EventType: "close"})
}

func (s *InvoiceService) MarkAsReopen(ctx context.Context, invoiceID int64) (*InvoiceMessage, *http.Response, error) {
	return s.SendEvent(ctx, invoiceID, &EventTypeRequest{EventType: "re-open"})
}

func (s *InvoiceService) SendEvent(
	ctx context.Context,
	invoiceID int64,
	data *EventTypeRequest,
) (*InvoiceMessage, *http.Response, error) {
	u := fmt.Sprintf("invoices/%d/messages", invoiceID)

	req, err := s.client.NewRequest(ctx, "POST", u, data)
	if err != nil {
		return nil, nil, err
	}

	invoiceMessage := new(InvoiceMessage)

	resp, err := s.client.Do(ctx, req, invoiceMessage)
	if err != nil {
		return nil, resp, err
	}

	return invoiceMessage, resp, nil
}
