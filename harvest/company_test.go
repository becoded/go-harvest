package harvest_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestCompanyService_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(mux *http.ServeMux)
		want      *harvest.Company
		wantErr   bool
	}{
		{
			name: "Valid Company Retrieval",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/company", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "company/get/body_1.json")
					testWriteResponse(t, w, "company/get/response_1.json")
				})
			},
			want: &harvest.Company{
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
			},
			wantErr: false,
		},
		{
			name: "Error Fetching Company",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/company", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					http.Error(w, `{"message":"Internal Server Error"}`, http.StatusInternalServerError)
				})
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			tt.setupMock(mux)

			got, _, err := service.Company.Get(context.Background())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCompany_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.Company
		want string
	}{
		{
			name: "Company with all fields",
			in: harvest.Company{
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
			},
			want: `harvest.Company{BaseURI:"https://organisation.harvestapp.com", FullDomain:"organisation.harvestapp.com", Name:"organisation", IsActive:true, WeekStartDay:"Monday", WantsTimestampTimers:false, TimeFormat:"hours_minutes", PlanType:"free", Clock:"24h", DecimalSymbol:",", ThousandsSeparator:".", ColorScheme:"blue", ExpenseFeature:true, InvoiceFeature:true, EstimateFeature:true, ApprovalFeature:false}`, //nolint: lll
		},
		{
			name: "Company with minimal fields",
			in: harvest.Company{
				Name:     harvest.String("Test Company"),
				IsActive: harvest.Bool(true),
			},
			want: `harvest.Company{Name:"Test Company", IsActive:true}`,
		},
		{
			name: "Company with features enabled",
			in: harvest.Company{
				Name:            harvest.String("Feature Company"),
				ExpenseFeature:  harvest.Bool(true),
				InvoiceFeature:  harvest.Bool(true),
				EstimateFeature: harvest.Bool(true),
				ApprovalFeature: harvest.Bool(true),
			},
			want: `harvest.Company{Name:"Feature Company", ExpenseFeature:true, InvoiceFeature:true, EstimateFeature:true, ApprovalFeature:true}`, //nolint: lll
		},
		{
			name: "Company with time settings",
			in: harvest.Company{
				Name:                 harvest.String("Time Company"),
				WeekStartDay:         harvest.String("Sunday"),
				WantsTimestampTimers: harvest.Bool(true),
				TimeFormat:           harvest.String("decimal"),
				Clock:                harvest.String("12h"),
			},
			want: `harvest.Company{Name:"Time Company", WeekStartDay:"Sunday", WantsTimestampTimers:true, TimeFormat:"decimal", Clock:"12h"}`, //nolint: lll
		},
		{
			name: "Company with number formatting",
			in: harvest.Company{
				Name:               harvest.String("Format Company"),
				DecimalSymbol:      harvest.String("."),
				ThousandsSeparator: harvest.String(","),
			},
			want: `harvest.Company{Name:"Format Company", DecimalSymbol:".", ThousandsSeparator:","}`,
		},
		{
			name: "Empty Company",
			in:   harvest.Company{},
			want: `harvest.Company{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.in.String()
			assert.Equal(t, tt.want, got)
		})
	}
}
