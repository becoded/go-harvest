package harvest

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompanyService_Get(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/company", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		testBody(t, r, "company/get/body_1.json")
		testWriteResponse(t, w, "company/get/response_1.json")
	})

	company, _, err := service.Company.Get(context.Background())
	assert.NoError(t, err)

	want := &Company{
		BaseUri:              String("https://organisation.harvestapp.com"),
		FullDomain:           String("organisation.harvestapp.com"),
		Name:                 String("organisation"),
		IsActive:             Bool(true),
		WeekStartDay:         String("Monday"),
		WantsTimestampTimers: Bool(false),
		TimeFormat:           String("hours_minutes"),
		PlanType:             String("free"),
		Clock:                String("24h"),
		DecimalSymbol:        String(","),
		ThousandsSeparator:   String("."),
		ColorScheme:          String("blue"),
		ExpenseFeature:       Bool(true),
		InvoiceFeature:       Bool(true),
		EstimateFeature:      Bool(true),
		ApprovalFeature:      Bool(false),
	}

	assert.Equal(t, want, company)
}
