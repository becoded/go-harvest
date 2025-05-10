package harvest_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/becoded/go-harvest/harvest"
	"github.com/stretchr/testify/assert"
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
