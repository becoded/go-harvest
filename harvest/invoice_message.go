package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Harvest API docs: https://help.getharvest.com/api-v2/invoices-api/invoices/invoice-messages/

type InvoiceMessage struct {
	Id                         *int64                     `json:"id,omitempty"`                             // Unique ID for the message.
	SentBy                     *string                    `json:"sent_by,omitempty"`                        // Name of the user that created the message.
	SentByEmail                *string                    `json:"sent_by_email,omitempty"`                  // Email of the user that created the message.
	SentFrom                   *string                    `json:"sent_from,omitempty"`                      // Name of the user that the message was sent from.
	SentFromEmail              *string                    `json:"sent_from_email,omitempty"`                // Email of the user that message was sent from.
	Recipients                 *[]InvoiceMessageRecipient `json:"recipients,omitempty"`                     // Array of invoice message recipients.
	Subject                    *string                    `json:"subject,omitempty"`                        // The message subject.
	Body                       *string                    `json:"body,omitempty"`                           // The message body.
	IncludeLinkToClientInvoice *bool                      `json:"include_link_to_client_invoice,omitempty"` // Whether to include a link to the client invoice in the message body. Not used when thank_you is true.
	AttachPdf                  *bool                      `json:"attach_pdf,omitempty"`                     // Whether to attach the invoice PDF to the message email.
	SendMeACopy                *bool                      `json:"send_me_a_copy,omitempty"`                 // Whether to email a copy of the message to the current user.
	ThankYou                   *bool                      `json:"thank_you,omitempty"`                      // Whether this is a thank you message.
	EventType                  *bool                      `json:"event_type,omitempty"`                     // The type of invoice event that occurred with the message: send, close, draft, re-open, or view.
	Reminder                   *bool                      `json:"reminder,omitempty"`                       // Whether this is a reminder message.
	SendReminderOn             *Date                      `json:"send_reminder_on,omitempty"`               // The date the reminder email will be sent.
	CreatedAt                  *time.Time                 `json:"created_at,omitempty"`                     // Date and time the message was created.
	UpdatedAt                  *time.Time                 `json:"updated_at,omitempty"`                     // Date and time the message was last updated.
}

type InvoiceMessageRecipient struct {
	Name  *string `json:"name,omitempty"` // Name of the message recipient.
	Email *string `json:"email"`          // Email of the message recipient.
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
	Recipients                 *[]InvoiceMessageRecipient `json:"recipients"`                               // required	Array of recipient parameters. See below for details.
	Subject                    *string                    `json:"subject,omitempty"`                        // optional	The message subject.
	Body                       *string                    `json:"body,omitempty"`                           // optional	The message body.
	IncludeLinkToClientInvoice *bool                      `json:"include_link_to_client_invoice,omitempty"` // optional	If set to true, a link to the client invoice URL will be included in the message email. Defaults to false. Ignored when thank_you is set to true.
	AttachPdf                  *bool                      `json:"attach_pdf,omitempty"`                     // optional	If set to true, a PDF of the invoice will be attached to the message email. Defaults to false.
	SendMeACopy                *bool                      `json:"send_me_a_copy,omitempty"`                 // optional	If set to true, a copy of the message email will be sent to the current user. Defaults to false.
	ThankYou                   *bool                      `json:"thank_you,omitempty"`                      // optional	If set to true, a thank you message email will be sent. Defaults to false.
	EventType                  *bool                      `json:"event_type,omitempty"`                     // optional	If provided, runs an event against the invoice. Options: close, draft, re-open, or send.
}

func (s *InvoiceService) ListMessages(ctx context.Context, invoiceId int64, opt *InvoiceMessageListOptions) (*InvoiceList, *http.Response, error) {
	u := fmt.Sprintf("invoices/%d/messages", invoiceId)
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

func (s *InvoiceService) CreateInvoiceMessage(ctx context.Context, invoiceId int64, data *InvoiceMessageCreateRequest) (*InvoiceMessage, *http.Response, error) {
	u := fmt.Sprintf("invoices/%d/messages", invoiceId)

	req, err := s.client.NewRequest("POST", u, data)
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

func (s *InvoiceService) DeleteInvoiceMessage(ctx context.Context, invoiceId, invoiceMessageId int64) (*http.Response, error) {
	u := fmt.Sprintf("invoices/%d/messages/%d", invoiceId, invoiceMessageId)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *InvoiceService) MarkAsSent(ctx context.Context, invoiceId int64) (*InvoiceMessage, *http.Response, error) {
	return s.SendEvent(ctx, invoiceId, &EventTypeRequest{EventType: "send"})

}

func (s *InvoiceService) MarkAsDraft(ctx context.Context, invoiceId int64) (*InvoiceMessage, *http.Response, error) {
	return s.SendEvent(ctx, invoiceId, &EventTypeRequest{EventType: "draft"})
}

func (s *InvoiceService) MarkAsClosed(ctx context.Context, invoiceId int64) (*InvoiceMessage, *http.Response, error) {
	return s.SendEvent(ctx, invoiceId, &EventTypeRequest{EventType: "close"})
}

func (s *InvoiceService) MarkAsReopen(ctx context.Context, invoiceId int64) (*InvoiceMessage, *http.Response, error) {
	return s.SendEvent(ctx, invoiceId, &EventTypeRequest{EventType: "re-open"})
}

func (s *InvoiceService) SendEvent(ctx context.Context, invoiceId int64, data *EventTypeRequest) (*InvoiceMessage, *http.Response, error) {
	u := fmt.Sprintf("invoices/%d/messages", invoiceId)
	req, err := s.client.NewRequest("POST", u, data)
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
