package harvest

import (
	"context"
	"fmt"
	"time"
	"net/http"
)

// EstimateService handles communication with the issue related
// methods of the Harvest API.
//
// Harvest API docs: https://help.getharvest.com/api-v2/estimates-api/estimates/estimates/
type EstimateService service

type Estimate struct {
	Id             *int64              `json:"id,omitempty"`              // Unique ID for the estimate.
	Client         *Client             `json:"client,omitempty"`          // An object containing estimate’s client id and name.
	LineItems      *[]EstimateLineItem `json:"line_items,omitempty"`      // Array of estimate line items.
	Creator        *User               `json:"creator,omitempty"`         // An object containing the id and name of the person that created the estimate.
	ClientKey      *string             `json:"client_key,omitempty"`      // Used to build a URL to the public web invoice for your client: https://{ACCOUNT_SUBDOMAIN}.harvestapp.com/client/invoices/abc123456
	Number         *string             `json:"number,omitempty"`          // If no value is set, the number will be automatically generated.
	PurchaseOrder  *string             `json:"purchase_order,omitempty"`  // The purchase order number.
	Amount         *float64            `json:"amount,omitempty"`          // The total amount for the estimate, including any discounts and taxes.
	Tax            *float64            `json:"tax,omitempty"`             // This percentage is applied to the subtotal, including line items and discounts.
	TaxAmount      *float64            `json:"tax_amount,omitempty"`      // The first amount of tax included, calculated from tax. If no tax is defined, this value will be null.
	Tax2           *float64            `json:"tax2,omitempty"`            // This percentage is applied to the subtotal, including line items and discounts.
	Tax2Amount     *float64            `json:"tax2_amount,omitempty"`     // The amount calculated from tax2.
	Discount       *float64            `json:"discount,omitempty"`        // This percentage is subtracted from the subtotal.
	DiscountAmount *float64            `json:"discount_amount,omitempty"` // The amount calcuated from discount.
	Subject        *string             `json:"subject,omitempty"`         // The estimate subject.
	Notes          *string             `json:"notes,omitempty"`           // Any additional notes included on the estimate.
	Currency       *string             `json:"currency,omitempty"`        // The currency code associated with this estimate.
	IssueDate      *Date               `json:"issue_date,omitempty"`      // Date the estimate was issued.
	SentAt         *time.Time          `json:"sent_at,omitempty"`         // Date and time the estimate was sent.
	AcceptedAt     *time.Time          `json:"accepted_at,omitempty"`     // Date and time the estimate was accepted.
	DeclinedAt     *time.Time          `json:"declined_at,omitempty"`     // Date and time the estimate was declined.
	CreatedAt      *time.Time          `json:"created_at,omitempty"`      // Date and time the estimate was created.
	UpdatedAt      *time.Time          `json:"updated_at,omitempty"`      // Date and time the estimate was last updated.
}

type EstimateLineItem struct {
	Id          *int64   `json:"id,omitempty"`          // Unique ID for the line item.
	Kind        *string  `json:"kind,omitempty"`        // The name of an estimate item category.
	Description *string  `json:"description,omitempty"` // Text description of the line item.
	Quantity    *int64   `json:"quantity,omitempty"`    // The unit quantity of the item.
	UnitPrice   *float64 `json:"unit_price,omitempty"`  // The individual price per unit.
	Amount      *float64 `json:"amount,omitempty"`      // The line item subtotal (quantity * unit_price).
	Taxed       *bool    `json:"taxed,omitempty"`       // Whether the estimate’s tax percentage applies to this line item.
	Taxed2      *bool    `json:"taxed2,omitempty"`      // Whether the estimate’s tax2 percentage applies to this line item.
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
	ClientId int64 `url:"client_id,omitempty"`
	// Only return estimates that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

func (s *EstimateService) List(ctx context.Context, opt *EstimateListOptions) (*EstimateList, *http.Response, error) {
	u := "estimates"
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

func (s *EstimateService) Get(ctx context.Context, estimateId int64) (*Estimate, *http.Response, error) {
	u := fmt.Sprintf("estimates/%d", estimateId)
	req, err := s.client.NewRequest("GET", u, nil)
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
