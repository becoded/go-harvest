package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Harvest API docs: https://help.getharvest.com/api-v2/estimates-api/estimates/estimate-messages/

type EstimateMessage struct {
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
	// Array of estimate message recipients.
	Recipients *[]EstimateMessageRecipient `json:"recipients,omitempty"`
	// The message subject.
	Subject *string `json:"subject,omitempty"`
	// The message body.
	Body *string `json:"body,omitempty"`
	// Whether to email a copy of the message to the current user.
	SendMeACopy *bool `json:"send_me_a_copy,omitempty"`
	// The type of estimate event that occurred with the message: send, accept, decline, re-open, view, or invoice.
	EventType *string `json:"event_type,omitempty"`
	// Date and time the message was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Date and time the message was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type EstimateMessageRecipient struct {
	// Name of the message recipient.
	Name *string `json:"name,omitempty"`
	// Email of the message recipient.
	Email *string `json:"email,omitempty"`
}

type EstimateMessageCreateRequest struct {
	// required	Array of recipient parameters. See below for details.
	Recipients *[]EstimateMessageRecipientCreateRequest `json:"recipients,omitempty"`
	// optional	The message subject.
	Subject *string `json:"subject,omitempty"`
	// optional	The message body.
	Body *string `json:"body,omitempty"`
	// optional	If set to true, a copy of the message email will be sent to the current user. Defaults to false.
	SendMeACopy *bool `json:"send_me_a_copy,omitempty"`
	// optional	If provided, runs an event against the estimate. Options: “accept”, “decline”, “re-open”, or “send”.
	EventType *string `json:"event_type,omitempty"`
}

type EstimateMessageRecipientCreateRequest struct {
	// optional	Name of the message recipient.
	Name *string `json:"name,omitempty"`
	// required	Email of the message recipient.
	Email *string `json:"email"`
}

type EstimateMessageList struct {
	EstimateMessages []*EstimateMessage `json:"estimate_messages"`

	Pagination
}

func (p EstimateMessage) String() string {
	return Stringify(p)
}

func (p EstimateMessageList) String() string {
	return Stringify(p)
}

type EstimateMessageListOptions struct {
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

type EstimateEventTypeRequest struct {
	EventType string `json:"event_type"`
}

// ListEstimateMessages returns a list of messages associated with a given estimate.
func (s *EstimateService) ListEstimateMessages(
	ctx context.Context,
	estimateID int64,
	opt *EstimateMessageListOptions,
) (*EstimateMessageList, *http.Response, error) {
	u := fmt.Sprintf("estimates/%d/messages", estimateID)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	estimateMessageList := new(EstimateMessageList)

	resp, err := s.client.Do(ctx, req, estimateMessageList)
	if err != nil {
		return nil, resp, err
	}

	return estimateMessageList, resp, nil
}

// CreateEstimateMessage creates a new estimate message object.
func (s *EstimateService) CreateEstimateMessage(
	ctx context.Context,
	estimateID int64,
	data *EstimateMessageCreateRequest,
) (*EstimateMessage, *http.Response, error) {
	u := fmt.Sprintf("estimates/%d/messages", estimateID)

	req, err := s.client.NewRequest(ctx, "POST", u, data)
	if err != nil {
		return nil, nil, err
	}

	estimateMessage := new(EstimateMessage)

	resp, err := s.client.Do(ctx, req, estimateMessage)
	if err != nil {
		return nil, resp, err
	}

	return estimateMessage, resp, nil
}

// DeleteEstimateMessage deletes an estimate message.
func (s *EstimateService) DeleteEstimateMessage(
	ctx context.Context,
	estimateID,
	estimateMessageID int64,
) (*http.Response, error) {
	u := fmt.Sprintf("estimates/%d/messages/%d", estimateID, estimateMessageID)

	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// MarkAsSent marks a draft estimate as sent.
func (s *EstimateService) MarkAsSent(
	ctx context.Context,
	estimateID int64,
) (*EstimateMessage, *http.Response, error) {
	return s.SendEvent(ctx, estimateID, &EstimateEventTypeRequest{EventType: "send"})
}

// MarkAsAccepted marks an open estimate as accepted.
func (s *EstimateService) MarkAsAccepted(
	ctx context.Context,
	estimateID int64,
) (*EstimateMessage, *http.Response, error) {
	return s.SendEvent(ctx, estimateID, &EstimateEventTypeRequest{EventType: "accept"})
}

// MarkAsDeclined marks an open estimate as declined.
func (s *EstimateService) MarkAsDeclined(
	ctx context.Context,
	estimateID int64,
) (*EstimateMessage, *http.Response, error) {
	return s.SendEvent(ctx, estimateID, &EstimateEventTypeRequest{EventType: "decline"})
}

// MarkAsReopen re-opens a closed estimate.
func (s *EstimateService) MarkAsReopen(
	ctx context.Context,
	estimateID int64,
) (*EstimateMessage, *http.Response, error) {
	return s.SendEvent(ctx, estimateID, &EstimateEventTypeRequest{EventType: "re-open"})
}

// SendEvent will send an EstimateEventType.
func (s *EstimateService) SendEvent(
	ctx context.Context,
	estimateID int64,
	data *EstimateEventTypeRequest,
) (*EstimateMessage, *http.Response, error) {
	u := fmt.Sprintf("estimates/%d/messages", estimateID)

	req, err := s.client.NewRequest(ctx, "POST", u, data)
	if err != nil {
		return nil, nil, err
	}

	estimateMessage := new(EstimateMessage)

	resp, err := s.client.Do(ctx, req, estimateMessage)
	if err != nil {
		return nil, resp, err
	}

	return estimateMessage, resp, nil
}
