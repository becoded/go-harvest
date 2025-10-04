package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestInvoiceService_CreateInvoice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		data      *harvest.InvoiceCreateRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.Invoice
		wantErr   bool
	}{
		{
			name: "Create a free-form invoice",
			data: func() *harvest.InvoiceCreateRequest {
				dueDate := time.Date(2017, 7, 27, 0, 0, 0, 0, time.UTC)

				lineItems := []harvest.InvoiceLineItemRequest{
					{
						Kind:        harvest.String("Service"),
						Description: harvest.String("ABC Project"),
						UnitPrice:   harvest.Float64(5000.0),
					},
				}

				return &harvest.InvoiceCreateRequest{
					ClientID: harvest.Int64(5735774),
					Subject:  harvest.String("ABC Project Quote"),
					DueDate: &harvest.Date{
						Time: dueDate,
					},
					LineItems: &lineItems,
				}
			}(),
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testFormValues(t, r, values{})
					testBody(t, r, "invoice/create/body_1.json")
					testWriteResponse(t, w, "invoice/create/response_1.json")
				})
			},
			want: func() *harvest.Invoice {
				issueDate := time.Date(
					2017, 6, 27, 0, 0, 0, 0, time.UTC)
				dueDate := time.Date(
					2017, 7, 27, 0, 0, 0, 0, time.UTC)
				createdOne := time.Date(
					2017, 6, 27, 16, 34, 24, 0, time.UTC)
				updatedOne := time.Date(
					2017, 6, 27, 16, 34, 24, 0, time.UTC)

				paymentOptions := []string{"credit_card"}

				lineItems := []harvest.InvoiceLineItem{
					{
						ID:          harvest.Int64(53341928),
						Kind:        harvest.String("Service"),
						Description: harvest.String("ABC Project"),
						Quantity:    harvest.Float64(1.0),
						UnitPrice:   harvest.Float64(5000.0),
						Amount:      harvest.Float64(5000.0),
						Taxed:       harvest.Bool(false),
						Taxed2:      harvest.Bool(false),
					},
				}

				return &harvest.Invoice{
					ID: harvest.Int64(13150453),
					Client: &harvest.Client{
						ID:   harvest.Int64(5735774),
						Name: harvest.String("ABC Corp"),
					},
					LineItems: &lineItems,
					Creator: &harvest.User{
						ID:   harvest.Int64(1782884),
						Name: harvest.String("Bob Powell"),
					},
					ClientKey:      harvest.String("8b86437630b6c260c1bfa289f0154960f83b606d"),
					Number:         harvest.String("1002"),
					Amount:         harvest.Float64(5000.0),
					DueAmount:      harvest.Float64(5000.0),
					TaxAmount:      harvest.Float64(0.0),
					Tax2Amount:     harvest.Float64(0.0),
					DiscountAmount: harvest.Float64(0.0),
					Currency:       harvest.String("USD"),
					Subject:        harvest.String("ABC Project Quote"),
					State:          harvest.String("draft"),
					IssueDate: &harvest.Date{
						Time: issueDate,
					},
					DueDate: &harvest.Date{
						Time: dueDate,
					},
					PaymentTerm:    harvest.String("custom"),
					PaymentOptions: &paymentOptions,
					CreatedAt:      &createdOne,
					UpdatedAt:      &updatedOne,
				}
			}(),
		},
		{
			name: "Create an invoice based on tracked time and expenses",
			data: func() *harvest.InvoiceCreateRequest {
				fromDate := time.Date(2017, 3, 1, 0, 0, 0, 0, time.UTC)
				toDate := time.Date(2017, 3, 31, 0, 0, 0, 0, time.UTC)
				projectIDs := []int64{14307913}
				lineItemsImport := harvest.InvoiceLineItemImportRequest{
					ProjectIDs: &projectIDs,
					Time: &harvest.InvoiceLineItemImportTimeRequest{
						SummaryType: harvest.String("task"),
						From: &harvest.Date{
							Time: fromDate,
						},
						To: &harvest.Date{
							Time: toDate,
						},
					},
					Expenses: &harvest.InvoiceLineItemImportExpenseRequest{
						SummaryType: harvest.String("category"),
					},
				}

				return &harvest.InvoiceCreateRequest{
					ClientID:        harvest.Int64(5735774),
					Subject:         harvest.String("ABC Project Quote"),
					PaymentTerm:     harvest.String("upon receipt"),
					LineItemsImport: &lineItemsImport,
				}
			}(),
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testFormValues(t, r, values{})
					testBody(t, r, "invoice/create/body_2.json")
					testWriteResponse(t, w, "invoice/create/response_2.json")
				})
			},
			want: func() *harvest.Invoice {
				issueDate := time.Date(
					2018, 2, 12, 0, 0, 0, 0, time.UTC)
				periodStartDate := time.Date(
					2017, 3, 1, 0, 0, 0, 0, time.UTC)
				periodEndDate := time.Date(
					2017, 3, 31, 0, 0, 0, 0, time.UTC)

				dueDate := time.Date(
					2018, 2, 12, 0, 0, 0, 0, time.UTC)
				createdOne := time.Date(
					2018, 2, 12, 21, 2, 37, 0, time.UTC)
				updatedOne := time.Date(
					2018, 2, 12, 21, 2, 37, 0, time.UTC)

				paymentOptions := []string{"credit_card"}

				lineItems := []harvest.InvoiceLineItem{
					{
						ID:          harvest.Int64(64957723),
						Kind:        harvest.String("Service"),
						Description: harvest.String("[MW] Marketing Website: Graphic Design (03/01/2017 - 03/31/2017)"),
						Quantity:    harvest.Float64(2.0),
						UnitPrice:   harvest.Float64(100.0),
						Amount:      harvest.Float64(200.0),
						Taxed:       harvest.Bool(false),
						Taxed2:      harvest.Bool(false),
						Project: &harvest.Project{
							ID:   harvest.Int64(14307913),
							Name: harvest.String("Marketing Website"),
							Code: harvest.String("MW"),
						},
					},
					{
						ID:          harvest.Int64(64957724),
						Kind:        harvest.String("Product"),
						Description: harvest.String("[MW] Marketing Website: Meals "),
						Quantity:    harvest.Float64(1.0),
						UnitPrice:   harvest.Float64(133.35),
						Amount:      harvest.Float64(133.35),
						Taxed:       harvest.Bool(false),
						Taxed2:      harvest.Bool(false),
						Project: &harvest.Project{
							ID:   harvest.Int64(14307913),
							Name: harvest.String("Marketing Website"),
							Code: harvest.String("MW"),
						},
					},
				}

				return &harvest.Invoice{
					ID: harvest.Int64(15340591),
					Client: &harvest.Client{
						ID:   harvest.Int64(5735774),
						Name: harvest.String("ABC Corp"),
					},
					LineItems: &lineItems,
					Creator: &harvest.User{
						ID:   harvest.Int64(1782884),
						Name: harvest.String("Bob Powell"),
					},
					ClientKey:      harvest.String("16173155e0a01542b8c7f689888cb3eaeda0dc94"),
					Number:         harvest.String("1002"),
					PurchaseOrder:  harvest.String(""),
					Amount:         harvest.Float64(333.35),
					DueAmount:      harvest.Float64(333.35),
					TaxAmount:      harvest.Float64(0.0),
					Tax2Amount:     harvest.Float64(0.0),
					DiscountAmount: harvest.Float64(0.0),
					Currency:       harvest.String("USD"),
					Subject:        harvest.String("ABC Project Quote"),
					Notes:          harvest.String(""),
					State:          harvest.String("draft"),
					PeriodStart: &harvest.Date{
						Time: periodStartDate,
					},
					PeriodEnd: &harvest.Date{
						Time: periodEndDate,
					},
					IssueDate: &harvest.Date{
						Time: issueDate,
					},
					DueDate: &harvest.Date{
						Time: dueDate,
					},
					PaymentTerm:    harvest.String("upon receipt"),
					PaymentOptions: &paymentOptions,
					CreatedAt:      &createdOne,
					UpdatedAt:      &updatedOne,
				}
			}(),
		},
		{
			name: "error creating invoice",
			data: &harvest.InvoiceCreateRequest{
				ClientID: harvest.Int64(5735774),
				Subject:  harvest.String("ABC Project Quote"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices", func(w http.ResponseWriter, r *http.Request) {
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

			invoice, _, err := service.Invoice.Create(context.Background(), tt.data)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, invoice)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, invoice)
			}
		})
	}
}

func TestInvoiceService_DeleteInvoice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		invoiceID int64
		setupMock func(mux *http.ServeMux)
		wantErr   bool
	}{
		{
			name:      "Delete an invoice",
			invoiceID: 13150453,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150453", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					testFormValues(t, r, values{})
					testBody(t, r, "invoice/delete/body_1.json")
					testWriteResponse(t, w, "invoice/delete/response_1.json")
				})
			},
			wantErr: false,
		},
		{
			name:      "error deleting invoice",
			invoiceID: 13150453,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150453", func(w http.ResponseWriter, r *http.Request) {
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

			_, err := service.Invoice.Delete(context.Background(), tt.invoiceID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInvoiceService_GetInvoice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		invoiceID int64
		setupMock func(mux *http.ServeMux)
		want      *harvest.Invoice
		wantErr   bool
	}{
		{
			name:      "get an invoice",
			invoiceID: 13150378,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150378", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "invoice/get/body_1.json")
					testWriteResponse(t, w, "invoice/get/response_1.json")
				})
			},
			want: func() *harvest.Invoice {
				issueDate := time.Date(
					2017, 2, 1, 0, 0, 0, 0, time.UTC)
				dueDate := time.Date(
					2017, 3, 3, 0, 0, 0, 0, time.UTC)
				createdOne := time.Date(
					2017, 6, 27, 16, 24, 30, 0, time.UTC)
				updatedOne := time.Date(
					2017, 6, 27, 16, 24, 57, 0, time.UTC)
				sentAt := time.Date(
					2017, 2, 1, 7, 0, 0, 0, time.UTC)
				paidAt := time.Date(
					2017, 2, 21, 0, 0, 0, 0, time.UTC)
				paymentOptions := []string{"credit_card"}

				lineItems := []harvest.InvoiceLineItem{
					{
						ID:          harvest.Int64(53341450),
						Kind:        harvest.String("Service"),
						Description: harvest.String("50% of Phase 1 of the Online Store"),
						Quantity:    harvest.Float64(100.0),
						UnitPrice:   harvest.Float64(100.0),
						Amount:      harvest.Float64(10000.0),
						Taxed:       harvest.Bool(true),
						Taxed2:      harvest.Bool(true),
						Project: &harvest.Project{
							ID:   harvest.Int64(14308069),
							Name: harvest.String("Online Store - Phase 1"),
							Code: harvest.String("OS1"),
						},
					},
				}

				return &harvest.Invoice{
					ID: harvest.Int64(13150378),
					Client: &harvest.Client{
						ID:   harvest.Int64(5735776),
						Name: harvest.String("123 Industries"),
					},
					LineItems: &lineItems,
					Creator: &harvest.User{
						ID:   harvest.Int64(1782884),
						Name: harvest.String("Bob Powell"),
					},
					ClientKey:      harvest.String("9e97f4a65c5b83b1fc02f54e5a41c9dc7d458542"),
					Number:         harvest.String("1000"),
					PurchaseOrder:  harvest.String("1234"),
					Amount:         harvest.Float64(10700.0),
					DueAmount:      harvest.Float64(0.0),
					Tax:            harvest.Float64(5.0),
					TaxAmount:      harvest.Float64(500.0),
					Tax2:           harvest.Float64(2.0),
					Tax2Amount:     harvest.Float64(200.0),
					DiscountAmount: harvest.Float64(0.0),
					Currency:       harvest.String("USD"),
					Subject:        harvest.String("Online Store - Phase 1"),
					Notes:          harvest.String("Some notes about the invoice."),
					State:          harvest.String("paid"),
					IssueDate: &harvest.Date{
						Time: issueDate,
					},
					DueDate: &harvest.Date{
						Time: dueDate,
					},
					Estimate: &harvest.Estimate{
						ID: harvest.Int64(1439814),
					},
					PaymentTerm:    harvest.String("custom"),
					PaymentOptions: &paymentOptions,
					SentAt:         &sentAt,
					PaidAt:         &paidAt,
					CreatedAt:      &createdOne,
					UpdatedAt:      &updatedOne,
				}
			}(),
		},
		{
			name:      "error getting invoice",
			invoiceID: 13150378,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150378", func(w http.ResponseWriter, r *http.Request) {
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

			invoice, _, err := service.Invoice.Get(context.Background(), tt.invoiceID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, invoice)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, invoice)
			}
		})
	}
}

func TestInvoiceService_ListInvoices(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		options   *harvest.InvoiceListOptions
		setupMock func(mux *http.ServeMux)
		want      *harvest.InvoiceList
		wantErr   bool
	}{
		{
			name:    "get invoice list",
			options: &harvest.InvoiceListOptions{},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "invoice/list/body_1.json")
					testWriteResponse(t, w, "invoice/list/response_1.json")
				})
			},
			want: func() *harvest.InvoiceList {
				issueDateOne := time.Date(
					2017, 4, 1, 0, 0, 0, 0, time.UTC)
				dueDateOne := time.Date(
					2017, 4, 1, 0, 0, 0, 0, time.UTC)
				periodStartOne := time.Date(
					2017, 3, 1, 0, 0, 0, 0, time.UTC)
				periodEndOne := time.Date(
					2017, 3, 1, 0, 0, 0, 0, time.UTC)
				createdOne := time.Date(
					2017, 6, 27, 16, 27, 16, 0, time.UTC)
				updatedOne := time.Date(
					2017, 8, 23, 22, 25, 59, 0, time.UTC)
				sentAtOne := time.Date(
					2017, 8, 23, 22, 25, 59, 0, time.UTC)

				paymentOptionsOne := []string{"credit_card"}

				lineItemsOne := []harvest.InvoiceLineItem{
					{
						ID:          harvest.Int64(53341602),
						Kind:        harvest.String("Service"),
						Description: harvest.String("03/01/2017 - Project Management: [9:00am - 11:00am] Planning meetings"),
						Quantity:    harvest.Float64(2.0),
						UnitPrice:   harvest.Float64(100.0),
						Amount:      harvest.Float64(200.0),
						Taxed:       harvest.Bool(true),
						Taxed2:      harvest.Bool(true),
						Project: &harvest.Project{
							ID:   harvest.Int64(14308069),
							Name: harvest.String("Online Store - Phase 1"),
							Code: harvest.String("OS1"),
						},
					},
					{
						ID:          harvest.Int64(53341603),
						Kind:        harvest.String("Service"),
						Description: harvest.String("03/01/2017 - Programming: [1:00pm - 2:00pm] Importing products"),
						Quantity:    harvest.Float64(1.0),
						UnitPrice:   harvest.Float64(100.0),
						Amount:      harvest.Float64(100.0),
						Taxed:       harvest.Bool(true),
						Taxed2:      harvest.Bool(true),
						Project: &harvest.Project{
							ID:   harvest.Int64(14308069),
							Name: harvest.String("Online Store - Phase 1"),
							Code: harvest.String("OS1"),
						},
					},
				}

				issueDateTwo := time.Date(
					2017, 2, 1, 0, 0, 0, 0, time.UTC)
				dueDateTwo := time.Date(
					2017, 3, 3, 0, 0, 0, 0, time.UTC)
				createdTwo := time.Date(
					2017, 6, 27, 16, 24, 30, 0, time.UTC)
				updatedTwo := time.Date(
					2017, 6, 27, 16, 24, 57, 0, time.UTC)
				sentAtTwo := time.Date(
					2017, 2, 1, 7, 0, 0, 0, time.UTC)
				paidAtTwo := time.Date(
					2017, 2, 21, 0, 0, 0, 0, time.UTC)

				lineItemsTwo := []harvest.InvoiceLineItem{
					{
						ID:          harvest.Int64(53341450),
						Kind:        harvest.String("Service"),
						Description: harvest.String("50% of Phase 1 of the Online Store"),
						Quantity:    harvest.Float64(100.0),
						UnitPrice:   harvest.Float64(100.0),
						Amount:      harvest.Float64(10000.0),
						Taxed:       harvest.Bool(true),
						Taxed2:      harvest.Bool(true),
						Project: &harvest.Project{
							ID:   harvest.Int64(14308069),
							Name: harvest.String("Online Store - Phase 1"),
							Code: harvest.String("OS1"),
						},
					},
				}

				return &harvest.InvoiceList{
					Invoices: []*harvest.Invoice{
						{
							ID: harvest.Int64(13150403),
							Client: &harvest.Client{
								ID:   harvest.Int64(5735776),
								Name: harvest.String("123 Industries"),
							},
							LineItems: &lineItemsOne,
							Creator: &harvest.User{
								ID:   harvest.Int64(1782884),
								Name: harvest.String("Bob Powell"),
							},
							ClientKey:      harvest.String("21312da13d457947a217da6775477afee8c2eba8"),
							Number:         harvest.String("1001"),
							PurchaseOrder:  harvest.String(""),
							Amount:         harvest.Float64(288.9),
							DueAmount:      harvest.Float64(288.9),
							Tax:            harvest.Float64(5),
							TaxAmount:      harvest.Float64(13.5),
							Tax2:           harvest.Float64(2),
							Tax2Amount:     harvest.Float64(5.4),
							Discount:       harvest.Float64(10.0),
							DiscountAmount: harvest.Float64(30.0),
							Currency:       harvest.String("EUR"),
							Subject:        harvest.String("Online Store - Phase 1"),
							Notes:          harvest.String("Some notes about the invoice."),
							State:          harvest.String("open"),
							IssueDate: &harvest.Date{
								Time: issueDateOne,
							},
							DueDate: &harvest.Date{
								Time: dueDateOne,
							},
							PeriodStart: &harvest.Date{
								Time: periodStartOne,
							},
							PeriodEnd: &harvest.Date{
								Time: periodEndOne,
							},
							PaymentTerm:    harvest.String("upon receipt"),
							PaymentOptions: &paymentOptionsOne,
							CreatedAt:      &createdOne,
							UpdatedAt:      &updatedOne,
							SentAt:         &sentAtOne,
						},
						{
							ID: harvest.Int64(13150378),
							Client: &harvest.Client{
								ID:   harvest.Int64(5735776),
								Name: harvest.String("123 Industries"),
							},
							LineItems: &lineItemsTwo,
							Creator: &harvest.User{
								ID:   harvest.Int64(1782884),
								Name: harvest.String("Bob Powell"),
							},
							ClientKey:      harvest.String("9e97f4a65c5b83b1fc02f54e5a41c9dc7d458542"),
							Number:         harvest.String("1000"),
							PurchaseOrder:  harvest.String("1234"),
							Amount:         harvest.Float64(10700.0),
							DueAmount:      harvest.Float64(0.0),
							Tax:            harvest.Float64(5.0),
							TaxAmount:      harvest.Float64(500.0),
							Tax2:           harvest.Float64(2.0),
							Tax2Amount:     harvest.Float64(200.0),
							DiscountAmount: harvest.Float64(0.0),
							Currency:       harvest.String("USD"),
							Subject:        harvest.String("Online Store - Phase 1"),
							Notes:          harvest.String("Some notes about the invoice."),
							State:          harvest.String("paid"),
							IssueDate: &harvest.Date{
								Time: issueDateTwo,
							},
							DueDate: &harvest.Date{
								Time: dueDateTwo,
							},
							Estimate: &harvest.Estimate{
								ID: harvest.Int64(1439814),
							},
							PaymentTerm: harvest.String("custom"),
							SentAt:      &sentAtTwo,
							PaidAt:      &paidAtTwo,
							CreatedAt:   &createdTwo,
							UpdatedAt:   &updatedTwo,
						},
					},
					Pagination: harvest.Pagination{
						PerPage:      harvest.Int(2000),
						TotalPages:   harvest.Int(1),
						TotalEntries: harvest.Int(2),
						NextPage:     nil,
						PreviousPage: nil,
						Page:         harvest.Int(1),
						Links: &harvest.PageLinks{
							First:    harvest.String("https://api.harvestapp.com/v2/invoices?page=1&per_page=2000"),
							Next:     nil,
							Previous: nil,
							Last:     harvest.String("https://api.harvestapp.com/v2/invoices?page=1&per_page=2000"),
						},
					},
				}
			}(),
		},
		{
			name:    "error listing invoices",
			options: &harvest.InvoiceListOptions{},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices", func(w http.ResponseWriter, r *http.Request) {
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

			invoices, _, err := service.Invoice.List(context.Background(), tt.options)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, invoices)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, invoices)
			}
		})
	}
}

func TestInvoiceService_UpdateInvoice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		invoiceID int64
		data      *harvest.InvoiceUpdateRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.Invoice
		wantErr   bool
	}{
		{
			name:      "update an invoice",
			invoiceID: 13150453,
			data: func() *harvest.InvoiceUpdateRequest {
				return &harvest.InvoiceUpdateRequest{
					PurchaseOrder: harvest.String("2345"),
				}
			}(),
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150453", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					testFormValues(t, r, values{})
					testBody(t, r, "invoice/update/body_1.json")
					testWriteResponse(t, w, "invoice/update/response_1.json")
				})
			},
			want: func() *harvest.Invoice {
				issueDate := time.Date(
					2017, 6, 27, 0, 0, 0, 0, time.UTC)
				dueDate := time.Date(
					2017, 7, 27, 0, 0, 0, 0, time.UTC)
				createdOne := time.Date(
					2017, 6, 27, 16, 34, 24, 0, time.UTC)
				updatedOne := time.Date(
					2017, 6, 27, 16, 36, 33, 0, time.UTC)

				paymentOptions := []string{"credit_card"}

				lineItems := []harvest.InvoiceLineItem{
					{
						ID:          harvest.Int64(53341928),
						Kind:        harvest.String("Service"),
						Description: harvest.String("ABC Project"),
						Quantity:    harvest.Float64(1.0),
						UnitPrice:   harvest.Float64(5000.0),
						Amount:      harvest.Float64(5000.0),
						Taxed:       harvest.Bool(false),
						Taxed2:      harvest.Bool(false),
					},
				}

				return &harvest.Invoice{
					ID: harvest.Int64(13150453),
					Client: &harvest.Client{
						ID:   harvest.Int64(5735774),
						Name: harvest.String("ABC Corp"),
					},
					PurchaseOrder: harvest.String("2345"),
					LineItems:     &lineItems,
					Creator: &harvest.User{
						ID:   harvest.Int64(1782884),
						Name: harvest.String("Bob Powell"),
					},
					ClientKey:      harvest.String("8b86437630b6c260c1bfa289f0154960f83b606d"),
					Number:         harvest.String("1002"),
					Amount:         harvest.Float64(5000.0),
					DueAmount:      harvest.Float64(5000.0),
					TaxAmount:      harvest.Float64(0.0),
					Tax2Amount:     harvest.Float64(0.0),
					DiscountAmount: harvest.Float64(0.0),
					Currency:       harvest.String("USD"),
					Subject:        harvest.String("ABC Project Quote"),
					State:          harvest.String("draft"),
					IssueDate: &harvest.Date{
						Time: issueDate,
					},
					DueDate: &harvest.Date{
						Time: dueDate,
					},
					PaymentTerm:    harvest.String("custom"),
					PaymentOptions: &paymentOptions,
					CreatedAt:      &createdOne,
					UpdatedAt:      &updatedOne,
				}
			}(),
		},
		{
			name:      "error updating invoice",
			invoiceID: 13150453,
			data: &harvest.InvoiceUpdateRequest{
				PurchaseOrder: harvest.String("2345"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150453", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
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

			invoice, _, err := service.Invoice.Update(context.Background(), tt.invoiceID, tt.data)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, invoice)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, invoice)
			}
		})
	}
}

func TestInvoice_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.Invoice
		want string
	}{
		{
			name: "Invoice with all fields",
			in: harvest.Invoice{
				ID: harvest.Int64(13150378),
				Client: &harvest.Client{
					ID:   harvest.Int64(5735776),
					Name: harvest.String("123 Industries"),
				},
				LineItems: &[]harvest.InvoiceLineItem{
					{
						ID:          harvest.Int64(53341450),
						Kind:        harvest.String("Service"),
						Description: harvest.String("50% of Phase 1 of the Online Store"),
						Quantity:    harvest.Float64(100.0),
						UnitPrice:   harvest.Float64(100.0),
						Amount:      harvest.Float64(10000.0),
						Taxed:       harvest.Bool(true),
						Taxed2:      harvest.Bool(true),
						Project: &harvest.Project{
							ID:   harvest.Int64(14308069),
							Name: harvest.String("Online Store - Phase 1"),
							Code: harvest.String("OS1"),
						},
					},
				},
				Estimate: &harvest.Estimate{
					ID: harvest.Int64(1439814),
				},
				Creator: &harvest.User{
					ID:   harvest.Int64(1782884),
					Name: harvest.String("Bob Powell"),
				},
				ClientKey:      harvest.String("9e97f4a65c5b83b1fc02f54e5a41c9dc7d458542"),
				Number:         harvest.String("1000"),
				PurchaseOrder:  harvest.String("1234"),
				Amount:         harvest.Float64(10700.0),
				DueAmount:      harvest.Float64(0.0),
				Tax:            harvest.Float64(5.0),
				TaxAmount:      harvest.Float64(500.0),
				Tax2:           harvest.Float64(2.0),
				Tax2Amount:     harvest.Float64(200.0),
				DiscountAmount: harvest.Float64(0.0),
				Currency:       harvest.String("USD"),
				Subject:        harvest.String("Online Store - Phase 1"),
				Notes:          harvest.String("Some notes about the invoice."),
				State:          harvest.String("paid"),
				IssueDate: &harvest.Date{
					Time: time.Date(2017, 2, 1, 0, 0, 0, 0, time.UTC),
				},
				DueDate: &harvest.Date{
					Time: time.Date(2017, 3, 3, 0, 0, 0, 0, time.UTC),
				},
				PaymentTerm:    harvest.String("custom"),
				PaymentOptions: &[]string{"credit_card"},
				SentAt:         harvest.TimeTimeP(time.Date(2017, 2, 1, 7, 0, 0, 0, time.UTC)),
				PaidAt:         harvest.TimeTimeP(time.Date(2017, 2, 21, 0, 0, 0, 0, time.UTC)),
				CreatedAt:      harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 24, 30, 0, time.UTC)),
				UpdatedAt:      harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 24, 57, 0, time.UTC)),
			},
			want: `harvest.Invoice{ID:13150378, Client:harvest.Client{ID:5735776, Name:"123 Industries"}, LineItems:[harvest.InvoiceLineItem{ID:53341450, Project:harvest.Project{ID:14308069, Name:"Online Store - Phase 1", Code:"OS1"}, Kind:"Service", Description:"50% of Phase 1 of the Online Store", Quantity:100, UnitPrice:100, Amount:10000, Taxed:true, Taxed2:true}], Estimate:harvest.Estimate{ID:1439814}, Creator:harvest.User{ID:1782884, Name:"Bob Powell"}, ClientKey:"9e97f4a65c5b83b1fc02f54e5a41c9dc7d458542", Number:"1000", PurchaseOrder:"1234", Amount:10700, DueAmount:0, Tax:5, TaxAmount:500, Tax2:2, Tax2Amount:200, DiscountAmount:0, Subject:"Online Store - Phase 1", Notes:"Some notes about the invoice.", Currency:"USD", State:"paid", IssueDate:harvest.Date{{2017-02-01 00:00:00 +0000 UTC}}, DueDate:harvest.Date{{2017-03-03 00:00:00 +0000 UTC}}, PaymentTerm:"custom", PaymentOptions:["credit_card"], SentAt:time.Time{2017-02-01 07:00:00 +0000 UTC}, PaidAt:time.Time{2017-02-21 00:00:00 +0000 UTC}, CreatedAt:time.Time{2017-06-27 16:24:30 +0000 UTC}, UpdatedAt:time.Time{2017-06-27 16:24:57 +0000 UTC}}`, //nolint: lll
		},
		{
			name: "Invoice with minimal fields",
			in: harvest.Invoice{
				ID:     harvest.Int64(999),
				Number: harvest.String("1001"),
				Client: &harvest.Client{
					ID:   harvest.Int64(5735774),
					Name: harvest.String("ABC Corp"),
				},
			},
			want: `harvest.Invoice{ID:999, Client:harvest.Client{ID:5735774, Name:"ABC Corp"}, Number:"1001"}`,
		},
		{
			name: "Invoice with period dates",
			in: harvest.Invoice{
				ID:     harvest.Int64(15340591),
				Number: harvest.String("1002"),
				PeriodStart: &harvest.Date{
					Time: time.Date(2017, 3, 1, 0, 0, 0, 0, time.UTC),
				},
				PeriodEnd: &harvest.Date{
					Time: time.Date(2017, 3, 31, 0, 0, 0, 0, time.UTC),
				},
				Amount:    harvest.Float64(333.35),
				DueAmount: harvest.Float64(333.35),
				State:     harvest.String("draft"),
			},
			want: `harvest.Invoice{ID:15340591, Number:"1002", Amount:333.35, DueAmount:333.35, State:"draft", PeriodStart:harvest.Date{{2017-03-01 00:00:00 +0000 UTC}}, PeriodEnd:harvest.Date{{2017-03-31 00:00:00 +0000 UTC}}}`, //nolint: lll
		},
		{
			name: "Invoice with discount",
			in: harvest.Invoice{
				ID:             harvest.Int64(13150403),
				Number:         harvest.String("1001"),
				Amount:         harvest.Float64(288.9),
				Discount:       harvest.Float64(10.0),
				DiscountAmount: harvest.Float64(30.0),
				Currency:       harvest.String("EUR"),
				State:          harvest.String("open"),
			},
			want: `harvest.Invoice{ID:13150403, Number:"1001", Amount:288.9, Discount:10, DiscountAmount:30, Currency:"EUR", State:"open"}`, //nolint: lll
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

func TestInvoiceList_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.InvoiceList
		want string
	}{
		{
			name: "InvoiceList with multiple invoices",
			in: harvest.InvoiceList{
				Invoices: []*harvest.Invoice{
					{
						ID: harvest.Int64(13150403),
						Client: &harvest.Client{
							ID:   harvest.Int64(5735776),
							Name: harvest.String("123 Industries"),
						},
						Number:   harvest.String("1001"),
						Amount:   harvest.Float64(288.9),
						Currency: harvest.String("EUR"),
						State:    harvest.String("open"),
					},
					{
						ID: harvest.Int64(13150378),
						Client: &harvest.Client{
							ID:   harvest.Int64(5735776),
							Name: harvest.String("123 Industries"),
						},
						Number:   harvest.String("1000"),
						Amount:   harvest.Float64(10700.0),
						Currency: harvest.String("USD"),
						State:    harvest.String("paid"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(2000),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.InvoiceList{Invoices:[harvest.Invoice{ID:13150403, Client:harvest.Client{ID:5735776, Name:"123 Industries"}, Number:"1001", Amount:288.9, Currency:"EUR", State:"open"} harvest.Invoice{ID:13150378, Client:harvest.Client{ID:5735776, Name:"123 Industries"}, Number:"1000", Amount:10700, Currency:"USD", State:"paid"}], Pagination:harvest.Pagination{PerPage:2000, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
		},
		{
			name: "InvoiceList with single invoice",
			in: harvest.InvoiceList{
				Invoices: []*harvest.Invoice{
					{
						ID:     harvest.Int64(999),
						Number: harvest.String("1005"),
						State:  harvest.String("draft"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(50),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.InvoiceList{Invoices:[harvest.Invoice{ID:999, Number:"1005", State:"draft"}], Pagination:harvest.Pagination{PerPage:50, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
		},
		{
			name: "Empty InvoiceList",
			in: harvest.InvoiceList{
				Invoices: []*harvest.Invoice{},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(0),
					TotalEntries: harvest.Int(0),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.InvoiceList{Invoices:[], Pagination:harvest.Pagination{PerPage:100, TotalPages:0, TotalEntries:0, Page:1}}`, //nolint: lll
		},
		{
			name: "InvoiceList with Links",
			in: harvest.InvoiceList{
				Invoices: []*harvest.Invoice{
					{
						ID:     harvest.Int64(100),
						Number: harvest.String("1010"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(25),
					TotalPages:   harvest.Int(3),
					TotalEntries: harvest.Int(75),
					Page:         harvest.Int(2),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/invoices?page=1&per_page=25"),
						Next:     harvest.String("https://api.harvestapp.com/v2/invoices?page=3&per_page=25"),
						Previous: harvest.String("https://api.harvestapp.com/v2/invoices?page=1&per_page=25"),
						Last:     harvest.String("https://api.harvestapp.com/v2/invoices?page=3&per_page=25"),
					},
				},
			},
			want: `harvest.InvoiceList{Invoices:[harvest.Invoice{ID:100, Number:"1010"}], Pagination:harvest.Pagination{PerPage:25, TotalPages:3, TotalEntries:75, Page:2, Links:harvest.PageLinks{First:"https://api.harvestapp.com/v2/invoices?page=1&per_page=25", Next:"https://api.harvestapp.com/v2/invoices?page=3&per_page=25", Previous:"https://api.harvestapp.com/v2/invoices?page=1&per_page=25", Last:"https://api.harvestapp.com/v2/invoices?page=3&per_page=25"}}}`, //nolint: lll
		},
		{
			name: "InvoiceList with line items",
			in: harvest.InvoiceList{
				Invoices: []*harvest.Invoice{
					{
						ID:     harvest.Int64(13150453),
						Number: harvest.String("1002"),
						LineItems: &[]harvest.InvoiceLineItem{
							{
								ID:          harvest.Int64(53341928),
								Kind:        harvest.String("Service"),
								Description: harvest.String("ABC Project"),
								Quantity:    harvest.Float64(1.0),
								UnitPrice:   harvest.Float64(5000.0),
								Amount:      harvest.Float64(5000.0),
								Taxed:       harvest.Bool(false),
								Taxed2:      harvest.Bool(false),
							},
						},
						Amount: harvest.Float64(5000.0),
						State:  harvest.String("draft"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(10),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.InvoiceList{Invoices:[harvest.Invoice{ID:13150453, LineItems:[harvest.InvoiceLineItem{ID:53341928, Kind:"Service", Description:"ABC Project", Quantity:1, UnitPrice:5000, Amount:5000, Taxed:false, Taxed2:false}], Number:"1002", Amount:5000, State:"draft"}], Pagination:harvest.Pagination{PerPage:10, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
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
