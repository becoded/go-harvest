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
		{
			name:      "Error Creating Payment",
			invoiceID: 13150378,
			request: &harvest.InvoicePaymentRequest{
				Amount: harvest.Float64(1575.86),
				PaidAt: harvest.TimeTimeP(time.Date(2017, 7, 24, 13, 32, 18, 0, time.UTC)),
				Notes:  harvest.String("Paid by phone"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150378/payments", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
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
		{
			name:      "Error Fetching Payment List",
			invoiceID: 13150378,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150378/payments", func(w http.ResponseWriter, r *http.Request) {
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
		{
			name:             "Error Deleting Payment",
			invoiceID:        13150378,
			invoicePaymentID: 10112854,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150378/payments/10112854", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					http.Error(w, `{"message":"Internal Server Error"}`, http.StatusInternalServerError)
				})
			},
			wantErr: true,
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

func TestInvoicePayment_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.InvoicePayment
		want string
	}{
		{
			name: "InvoicePayment with all fields",
			in: harvest.InvoicePayment{
				ID:              harvest.Int64(10336386),
				Amount:          harvest.Float64(1575.86),
				PaidAt:          harvest.TimeTimeP(time.Date(2017, 7, 24, 13, 32, 18, 0, time.UTC)),
				RecordedBy:      harvest.String("Jane Bar"),
				RecordedByEmail: harvest.String("jane@example.com"),
				Notes:           harvest.String("Paid by phone"),
				TransactionID:   harvest.String("txn_12345"),
				CreatedAt:       harvest.TimeTimeP(time.Date(2017, 7, 28, 14, 42, 44, 0, time.UTC)),
				UpdatedAt:       harvest.TimeTimeP(time.Date(2017, 7, 28, 14, 42, 44, 0, time.UTC)),
				PaymentGateway: &harvest.PaymentGateway{
					ID:   harvest.Int64(1234),
					Name: harvest.String("Stripe"),
				},
			},
			want: `harvest.InvoicePayment{ID:10336386, Amount:1575.86, PaidAt:time.Time{2017-07-24 13:32:18 +0000 UTC}, RecordedBy:"Jane Bar", RecordedByEmail:"jane@example.com", Notes:"Paid by phone", TransactionID:"txn_12345", PaymentGateway:harvest.PaymentGateway{ID:1234, Name:"Stripe"}, CreatedAt:time.Time{2017-07-28 14:42:44 +0000 UTC}, UpdatedAt:time.Time{2017-07-28 14:42:44 +0000 UTC}}`, //nolint: lll
		},
		{
			name: "InvoicePayment with minimal fields",
			in: harvest.InvoicePayment{
				ID:     harvest.Int64(999),
				Amount: harvest.Float64(100.00),
			},
			want: `harvest.InvoicePayment{ID:999, Amount:100}`,
		},
		{
			name: "InvoicePayment with nil PaymentGateway fields",
			in: harvest.InvoicePayment{
				ID:             harvest.Int64(10112854),
				Amount:         harvest.Float64(10700),
				PaidAt:         harvest.TimeTimeP(time.Date(2017, 2, 21, 0, 0, 0, 0, time.UTC)),
				PaymentGateway: &harvest.PaymentGateway{ID: nil, Name: nil},
			},
			want: `harvest.InvoicePayment{ID:10112854, Amount:10700, PaidAt:time.Time{2017-02-21 00:00:00 +0000 UTC}, PaymentGateway:harvest.PaymentGateway{}}`, //nolint: lll
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

func TestInvoicePaymentList_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.InvoicePaymentList
		want string
	}{
		{
			name: "InvoicePaymentList with multiple payments",
			in: harvest.InvoicePaymentList{
				InvoicePayments: []*harvest.InvoicePayment{
					{
						ID:              harvest.Int64(10112854),
						Amount:          harvest.Float64(10700),
						PaidAt:          harvest.TimeTimeP(time.Date(2017, 2, 21, 0, 0, 0, 0, time.UTC)),
						RecordedBy:      harvest.String("Alice Doe"),
						RecordedByEmail: harvest.String("alice@example.com"),
						Notes:           harvest.String("Paid via check #4321"),
					},
					{
						ID:              harvest.Int64(10336386),
						Amount:          harvest.Float64(1575.86),
						PaidAt:          harvest.TimeTimeP(time.Date(2017, 7, 24, 13, 32, 18, 0, time.UTC)),
						RecordedBy:      harvest.String("Jane Bar"),
						RecordedByEmail: harvest.String("jane@example.com"),
						Notes:           harvest.String("Paid by phone"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(2000),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.InvoicePaymentList{InvoicePayments:[harvest.InvoicePayment{ID:10112854, Amount:10700, PaidAt:time.Time{2017-02-21 00:00:00 +0000 UTC}, RecordedBy:"Alice Doe", RecordedByEmail:"alice@example.com", Notes:"Paid via check #4321"} harvest.InvoicePayment{ID:10336386, Amount:1575.86, PaidAt:time.Time{2017-07-24 13:32:18 +0000 UTC}, RecordedBy:"Jane Bar", RecordedByEmail:"jane@example.com", Notes:"Paid by phone"}], Pagination:harvest.Pagination{PerPage:2000, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
		},
		{
			name: "InvoicePaymentList with single payment",
			in: harvest.InvoicePaymentList{
				InvoicePayments: []*harvest.InvoicePayment{
					{
						ID:     harvest.Int64(999),
						Amount: harvest.Float64(500.00),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.InvoicePaymentList{InvoicePayments:[harvest.InvoicePayment{ID:999, Amount:500}], Pagination:harvest.Pagination{PerPage:100, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
		},
		{
			name: "Empty InvoicePaymentList",
			in: harvest.InvoicePaymentList{
				InvoicePayments: []*harvest.InvoicePayment{},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(2000),
					TotalPages:   harvest.Int(0),
					TotalEntries: harvest.Int(0),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.InvoicePaymentList{InvoicePayments:[], Pagination:harvest.Pagination{PerPage:2000, TotalPages:0, TotalEntries:0, Page:1}}`, //nolint: lll
		},
		{
			name: "InvoicePaymentList with Links",
			in: harvest.InvoicePaymentList{
				InvoicePayments: []*harvest.InvoicePayment{
					{
						ID:     harvest.Int64(100),
						Amount: harvest.Float64(250.00),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(50),
					TotalPages:   harvest.Int(3),
					TotalEntries: harvest.Int(150),
					Page:         harvest.Int(2),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/invoices/13150378/payments?page=1&per_page=50"),
						Next:     harvest.String("https://api.harvestapp.com/v2/invoices/13150378/payments?page=3&per_page=50"),
						Previous: harvest.String("https://api.harvestapp.com/v2/invoices/13150378/payments?page=1&per_page=50"),
						Last:     harvest.String("https://api.harvestapp.com/v2/invoices/13150378/payments?page=3&per_page=50"),
					},
				},
			},
			want: `harvest.InvoicePaymentList{InvoicePayments:[harvest.InvoicePayment{ID:100, Amount:250}], Pagination:harvest.Pagination{PerPage:50, TotalPages:3, TotalEntries:150, Page:2, Links:harvest.PageLinks{First:"https://api.harvestapp.com/v2/invoices/13150378/payments?page=1&per_page=50", Next:"https://api.harvestapp.com/v2/invoices/13150378/payments?page=3&per_page=50", Previous:"https://api.harvestapp.com/v2/invoices/13150378/payments?page=1&per_page=50", Last:"https://api.harvestapp.com/v2/invoices/13150378/payments?page=3&per_page=50"}}}`, //nolint: lll
		},
		{
			name: "InvoicePaymentList with PaymentGateway",
			in: harvest.InvoicePaymentList{
				InvoicePayments: []*harvest.InvoicePayment{
					{
						ID:     harvest.Int64(1),
						Amount: harvest.Float64(100.00),
						PaymentGateway: &harvest.PaymentGateway{
							ID:   harvest.Int64(10),
							Name: harvest.String("Stripe"),
						},
					},
					{
						ID:     harvest.Int64(2),
						Amount: harvest.Float64(200.00),
						PaymentGateway: &harvest.PaymentGateway{
							ID:   harvest.Int64(20),
							Name: harvest.String("PayPal"),
						},
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(10),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.InvoicePaymentList{InvoicePayments:[harvest.InvoicePayment{ID:1, Amount:100, PaymentGateway:harvest.PaymentGateway{ID:10, Name:"Stripe"}} harvest.InvoicePayment{ID:2, Amount:200, PaymentGateway:harvest.PaymentGateway{ID:20, Name:"PayPal"}}], Pagination:harvest.Pagination{PerPage:10, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
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

func TestPaymentGateway_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.PaymentGateway
		want string
	}{
		{
			name: "PaymentGateway with all fields",
			in: harvest.PaymentGateway{
				ID:   harvest.Int64(1234),
				Name: harvest.String("Stripe"),
			},
			want: `harvest.PaymentGateway{ID:1234, Name:"Stripe"}`,
		},
		{
			name: "PaymentGateway with ID only",
			in: harvest.PaymentGateway{
				ID: harvest.Int64(5678),
			},
			want: `harvest.PaymentGateway{ID:5678}`,
		},
		{
			name: "PaymentGateway with Name only",
			in: harvest.PaymentGateway{
				Name: harvest.String("PayPal"),
			},
			want: `harvest.PaymentGateway{Name:"PayPal"}`,
		},
		{
			name: "Empty PaymentGateway",
			in:   harvest.PaymentGateway{},
			want: `harvest.PaymentGateway{}`,
		},
		{
			name: "PaymentGateway pointer with all fields",
			in: harvest.PaymentGateway{
				ID:   harvest.Int64(9999),
				Name: harvest.String("Square"),
			},
			want: `harvest.PaymentGateway{ID:9999, Name:"Square"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.in.String()
			assert.Equal(t, tt.want, got)

			// Also test with pointer to ensure coverage
			ptr := &tt.in
			gotPtr := ptr.String()
			assert.Equal(t, tt.want, gotPtr)
		})
	}
}
