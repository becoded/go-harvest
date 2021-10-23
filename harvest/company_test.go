package harvest_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/becoded/go-harvest/harvest"
	"github.com/stretchr/testify/assert"
)

func TestCompanyService_Get(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/company", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `{"base_uri":"https://organisation.harvestapp.com","full_domain":"organisation.harvestapp.com","name":"organisation","is_active":true,"week_start_day":"Monday","wants_timestamp_timers": false,"time_format":"hours_minutes","plan_type":"free","clock":"24h","decimal_symbol":",","thousands_separator":".","color_scheme":"blue","expense_feature":true,"invoice_feature":true,"estimate_feature":true,"approval_feature":false}`)
	})

	company, _, err := service.Company.Get(context.Background())
	assert.NoError(t, err)

	want := &harvest.Company{
		BaseURI:              harvest.String("https://organisation.harvestapp.com"),
		FullDomain:           harvest.String("organisation.harvestapp.com"),
		Name:                 harvest.String("organisation"),
		IsActive:             harvest.Bool(true),
		WeekStartDay:         harvest.String("Monday"),
		WantsTimestampTimers: harvest.Bool(false),
		TimeFormat:           harvest.String("hours_minutes"),
		PlanType:             harvest.String("free"),
		Clock:                harvest.String("24h"),
		DecimalSymbol:        harvest.String(","),
		ThousandsSeparator:   harvest.String("."),
		ColorScheme:          harvest.String("blue"),
		ExpenseFeature:       harvest.Bool(true),
		InvoiceFeature:       harvest.Bool(true),
		EstimateFeature:      harvest.Bool(true),
		ApprovalFeature:      harvest.Bool(false),
	}

	assert.Equal(t, want, company)
}
