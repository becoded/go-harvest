package harvest


import (
"context"
"fmt"
"time"
	"net/http"
)

// Harvest API docs: https://help.getharvest.com/api-v2/estimates-api/estimates/estimate-messages/

type EstimateMessage struct {
	Id *int64 `json:"id,omitempty"` // Unique ID for the message.
	SentBy *string `json:"sent_by,omitempty"` // Name of the user that created the message.
	SentByEmail *string `json:"sent_by_email,omitempty"` // Email of the user that created the message.
	SentFrom *string `json:"sent_from,omitempty"` // Name of the user that the message was sent from.
	SentFromEmail *string `json:"sent_from_email,omitempty"` // Email of the user that message was sent from.
	Recipients *[]EstimateMessageRecipient `json:"recipients,omitempty"` // Array of estimate message recipients.
	Subject *string `json:"subject,omitempty"` // The message subject.
	Body *string `json:"body,omitempty"` // The message body.
	SendMeACopy *bool `json:"send_me_a_copy,omitempty"` // Whether to email a copy of the message to the current user.
	EventType *bool `json:"event_type,omitempty"` // The type of estimate event that occurred with the message: send, accept, decline, re-open, view, or invoice.
	CreatedAt *time.Time `json:"created_at,omitempty"` // Date and time the message was created.
	UpdatedAt *time.Time `json:"updated_at,omitempty"` // Date and time the message was last updated.
}

type EstimateMessageRecipient  struct {
	Name *string `json:"name,omitempty"` // Name of the message recipient.
	Email *string `json:"email,omitempty"` // Email of the message recipient.
}

type EstimateMessageList struct {
	Estimates []*Estimate `json:"estimates"`

	Pagination
}

func (p EstimateMessage) String() string {
	return Stringify(p)
}

func (p EstimateMessageList) String() string {
	return Stringify(p)
}

type EstimateMessageListOptions struct {
	UpdatedSince	time.Time `url:"updated_since,omitempty"`

	ListOptions
}

type EstimateEventTypeRequest struct {
	EventType string `json:"event_type"`
}

func (s *EstimateService) ListMessages(ctx context.Context, estimateId int64, opt *EstimateMessageListOptions) (*EstimateList, *http.Response, error) {
	u := fmt.Sprintf("estimates/%d/messages", estimateId)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	estimateList := new(EstimateList)
	resp, err := s.client.Do(ctx, req, &estimateList)
	if err != nil {
		return nil, resp, err
	}

	return estimateList, resp, nil
}

func (s *EstimateService) CreateEstimateMessage(ctx context.Context)  {
	// TODO
}

func (s *EstimateService) DeleteEstimateMessage(ctx context.Context, estimateId, estimateMessageId int64) (*http.Response, error) {
	u := fmt.Sprintf("estimates/%d/messages/%d", estimateId, estimateMessageId)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *EstimateService) MarkAsSent(ctx context.Context, estimateId int64) (*EstimateMessage, *http.Response, error) {
	return s.SendEvent(ctx, estimateId, &EstimateEventTypeRequest{EventType:"send"})

}

func (s *EstimateService) MarkAsAccepted(ctx context.Context, estimateId int64) (*EstimateMessage, *http.Response, error) {
	return s.SendEvent(ctx, estimateId, &EstimateEventTypeRequest{EventType:"accept"})
}

func (s *EstimateService) MarkAsDeclined(ctx context.Context, estimateId int64) (*EstimateMessage, *http.Response, error) {
	return s.SendEvent(ctx, estimateId, &EstimateEventTypeRequest{EventType:"decline"})
}

func (s *EstimateService) MarkAsReopen(ctx context.Context, estimateId int64) (*EstimateMessage, *http.Response, error) {
	return s.SendEvent(ctx, estimateId, &EstimateEventTypeRequest{EventType:"re-open"})
}

func (s *EstimateService) SendEvent(ctx context.Context, estimateId int64, data *EstimateEventTypeRequest) (*EstimateMessage, *http.Response, error) {
	u := fmt.Sprintf("estimates/%d/messages", estimateId)
	req, err := s.client.NewRequest("POST", u, data)
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