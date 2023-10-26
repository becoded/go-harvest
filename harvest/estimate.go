package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// EstimateService handles communication with the issue related
// methods of the Harvest API.
//
// Harvest API docs: https://help.getharvest.com/api-v2/estimates-api/estimates/estimates/
type EstimateService service

type Estimate struct {
	// Unique ID for the estimate.
	ID *int64 `json:"id,omitempty"`
	// An object containing estimate’s client id and name.
	Client *Client `json:"client,omitempty"`
	// Array of estimate line items.
	LineItems *[]EstimateLineItem `json:"line_items,omitempty"`
	// An object containing the id and name of the person that created the estimate.
	Creator *User `json:"creator,omitempty"`
	// Used to build a URL to the public web invoice for your client:
	// https://{ACCOUNT_SUBDOMAIN}.harvestapp.com/client/invoices/abc123456
	ClientKey *string `json:"client_key,omitempty"`
	// If no value is set, the number will be automatically generated.
	Number *string `json:"number,omitempty"`
	// The purchase order number.
	PurchaseOrder *string `json:"purchase_order,omitempty"`
	// The total amount for the estimate, including any discounts and taxes.
	Amount *float64 `json:"amount,omitempty"`
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
	// The estimate subject.
	Subject *string `json:"subject,omitempty"`
	// Any additional notes included on the estimate.
	Notes *string `json:"notes,omitempty"`
	// The currency code associated with this estimate.
	Currency *string `json:"currency,omitempty"`
	// The current state of the estimate: draft, sent, accepted, or declined.
	State *string `json:"state,omitempty"`
	// Date the estimate was issued.
	IssueDate *Date `json:"issue_date,omitempty"`
	// Date and time the estimate was sent.
	SentAt *time.Time `json:"sent_at,omitempty"`
	// Date and time the estimate was accepted.
	AcceptedAt *time.Time `json:"accepted_at,omitempty"`
	// Date and time the estimate was declined.
	DeclinedAt *time.Time `json:"declined_at,omitempty"`
	// Date and time the estimate was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Date and time the estimate was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type EstimateLineItem struct {
	// Unique ID for the line item.
	ID *int64 `json:"id,omitempty"`
	// The name of an estimate item category.
	Kind *string `json:"kind,omitempty"`
	// Text description of the line item.
	Description *string `json:"description,omitempty"`
	// The unit quantity of the item.
	Quantity *int64 `json:"quantity,omitempty"`
	// The individual price per unit.
	UnitPrice *float64 `json:"unit_price,omitempty"`
	// The line item subtotal (quantity * unit_price).
	Amount *float64 `json:"amount,omitempty"`
	// Whether the estimate’s tax percentage applies to this line item.
	Taxed *bool `json:"taxed,omitempty"`
	// Whether the estimate’s tax2 percentage applies to this line item.
	Taxed2 *bool `json:"taxed2,omitempty"`
}

type EstimateList struct {
	Estimates []*Estimate `json:"estimates"`

	Pagination
}

func (p Estimate) String() string {
	return Stringify(p)
}

func (p EstimateList) String() string {
	return Stringify(p)
}

type EstimateListOptions struct {
	// Only return estimates belonging to the client with the given ID.
	ClientID int64 `url:"client_id,omitempty"`
	// Only return estimates that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

// List will return a list of your estimates.
func (s *EstimateService) List(ctx context.Context, opt *EstimateListOptions) (*EstimateList, *http.Response, error) {
	u := "estimates"

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
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

// Get retrieves the estimate with the given ID.
func (s *EstimateService) Get(ctx context.Context, estimateID int64) (*Estimate, *http.Response, error) {
	u := fmt.Sprintf("estimates/%d", estimateID)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	estimate := new(Estimate)

	resp, err := s.client.Do(ctx, req, estimate)
	if err != nil {
		return nil, resp, err
	}

	return estimate, resp, nil
}
