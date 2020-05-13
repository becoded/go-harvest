package harvest

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestCompanyService_GetCompany(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/company", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `{"base_uri":"https://organisation.harvestapp.com","full_domain":"organisation.harvestapp.com","name":"organisation","is_active":true,"week_start_day":"Monday","wants_timestamp_timers": false,"time_format":"hours_minutes","plan_type":"free","clock":"24h","decimal_symbol":",","thousands_separator":".","color_scheme":"blue","expense_feature":true,"invoice_feature":true,"estimate_feature":true,"approval_feature":false}`)
	})

	company, _, err := service.Company.GetCompany(context.Background())
	if err != nil {
		t.Errorf("Company.GetCompany returned error: %v", err)
	}

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

	if !reflect.DeepEqual(company, want) {
		t.Errorf("Company.GetCompany returned %+v, want %+v", company, want)
	}
}
