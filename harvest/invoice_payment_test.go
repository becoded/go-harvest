package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestInvoicePaymentService_CreatePayment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		invoiceID int64
		request   *harvest.InvoicePaymentRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.InvoicePayment
		wantErr   bool
	}{
		{
			name:      "Valid Payment Creation",
			invoiceID: 13150378,
			request: &harvest.InvoicePaymentRequest{
				Amount: harvest.Float64(1575.86),
				PaidAt: harvest.TimeTimeP(time.Date(2017, 7, 24, 13, 32, 18, 0, time.UTC)),
				Notes:  harvest.String("Paid by phone"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150378/payments", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testBody(t, r, "invoice_payment/create/body_1.json")
					testWriteResponse(t, w, "invoice_payment/create/response_1.json")
				})
			},
			want: &harvest.InvoicePayment{
				ID:              harvest.Int64(10336386),
				Amount:          harvest.Float64(1575.86),
				PaidAt:          harvest.TimeTimeP(time.Date(2017, 7, 24, 13, 32, 18, 0, time.UTC)),
				RecordedBy:      harvest.String("Jane Bar"),
				RecordedByEmail: harvest.String("jane@example.com"),
				Notes:           harvest.String("Paid by phone"),
				TransactionID:   nil,
				CreatedAt:       harvest.TimeTimeP(time.Date(2017, 7, 28, 14, 42, 44, 0, time.UTC)),
				UpdatedAt:       harvest.TimeTimeP(time.Date(2017, 7, 28, 14, 42, 44, 0, time.UTC)),
				PaymentGateway:  &harvest.PaymentGateway{ID: nil, Name: nil},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			tt.setupMock(mux)

			got, _, err := service.Invoice.CreatePayment(context.Background(), tt.invoiceID, tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestInvoicePaymentService_ListPayments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		invoiceID int64
		setupMock func(mux *http.ServeMux)
		want      *harvest.InvoicePaymentList
		wantErr   bool
	}{
		{
			name:      "Valid Payment List",
			invoiceID: 13150378,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150378/payments", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testWriteResponse(t, w, "invoice_payment/list/response_1.json")
				})
			},
			want: &harvest.InvoicePaymentList{
				InvoicePayments: []*harvest.InvoicePayment{
					{
						ID:              harvest.Int64(10112854),
						Amount:          harvest.Float64(10700),
						PaidAt:          harvest.TimeTimeP(time.Date(2017, 2, 21, 0, 0, 0, 0, time.UTC)),
						RecordedBy:      harvest.String("Alice Doe"),
						RecordedByEmail: harvest.String("alice@example.com"),
						Notes:           harvest.String("Paid via check #4321"),
						TransactionID:   nil,
						CreatedAt:       harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 24, 57, 0, time.UTC)),
						UpdatedAt:       harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 24, 57, 0, time.UTC)),
						PaymentGateway: &harvest.PaymentGateway{
							ID:   harvest.Int64(1234),
							Name: harvest.String("Linkpoint International"),
						},
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(2000),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
					Links: &harvest.PageLinks{
						First: harvest.String("https://api.harvestapp.com/v2/invoices/13150378/payments?page=1&per_page=2000"),
						Last:  harvest.String("https://api.harvestapp.com/v2/invoices/13150378/payments?page=1&per_page=2000"),
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			tt.setupMock(mux)

			got, _, err := service.Invoice.ListPayments(context.Background(), tt.invoiceID, &harvest.InvoicePaymentListOptions{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestInvoicePaymentService_DeletePayment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		invoiceID        int64
		invoicePaymentID int64
		setupMock        func(mux *http.ServeMux)
		wantErr          bool
	}{
		{
			name:             "Valid Payment Deletion",
			invoiceID:        13150378,
			invoicePaymentID: 10112854,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150378/payments/10112854", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					w.WriteHeader(http.StatusNoContent)
				})
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			tt.setupMock(mux)

			_, err := service.Invoice.DeletePayment(context.Background(), tt.invoiceID, tt.invoicePaymentID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
