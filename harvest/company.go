package harvest

import (
	"context"
	"net/http"
)

// CompanyService handles communication with the company related
// methods of the Harvest API.
//
// Harvest API docs: https://help.getharvest.com/api-v2/company-api/company/company/
type CompanyService service

type Company struct {
	// The Harvest URL for the company.
	BaseURI *string `json:"base_uri,omitempty"`
	// The Harvest domain for the company.
	FullDomain *string `json:"full_domain,omitempty"`
	// The name of the company.
	Name *string `json:"name,omitempty"`
	// Whether the company is active or archived.
	IsActive *bool `json:"is_active,omitempty"`
	// The week day used as the start of the week. Returns one of: Saturday, Sunday, or Monday.
	WeekStartDay *string `json:"week_start_day,omitempty"`
	// Whether time is tracked via duration or start and end times.
	WantsTimestampTimers *bool `json:"wants_timestamp_timers,omitempty"`
	// The format used to display time in Harvest. Returns either decimal or hours_minutes.
	TimeFormat *string `json:"time_format,omitempty"`
	// The type of plan the company is on. Examples: trial, free, or simple-v4
	PlanType *string `json:"plan_type,omitempty"`
	// Used to represent whether the company is using a 12-hour or 24-hour clock. Returns either 12h or 24h.
	Clock *string `json:"clock,omitempty"`
	// Symbol *used `json:"Symbol,omitempty"` //when *formatting `json:"when,omitempty"` //decimals.
	DecimalSymbol *string `json:"decimal_symbol,omitempty"`
	// Separator *used `json:"Separator,omitempty"` //when formatting numbers.
	ThousandsSeparator *string `json:"thousands_separator,omitempty"`
	// The color scheme being used in the Harvest web client.
	ColorScheme *string `json:"color_scheme,omitempty"`
	// Whether the expense module is enabled.
	ExpenseFeature *bool `json:"expense_feature,omitempty"`
	// Whether the invoice module is enabled.
	InvoiceFeature *bool `json:"invoice_feature,omitempty"`
	// Whether the estimate module is enabled.
	EstimateFeature *bool `json:"estimate_feature,omitempty"`
	// Whether *the `json:"Whether,omitempty"` //approval module is enabled.
	ApprovalFeature *bool `json:"approval_feature,omitempty"`
}

func (c Company) String() string {
	return Stringify(c)
}

// Get retrieves the company for the currently authenticated user.
func (s *CompanyService) Get(ctx context.Context) (*Company, *http.Response, error) {
	u := "company"

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	company := new(Company)

	resp, err := s.client.Do(ctx, req, &company)
	if err != nil {
		return nil, resp, err
	}

	return company, resp, nil
}
