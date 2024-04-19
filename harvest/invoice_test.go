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

	type fields struct {
		path     string
		method   string
		body     string
		response string
	}

	type args struct {
		ctx  context.Context
		data *harvest.InvoiceCreateRequest
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *harvest.Invoice
		wantError error
	}{
		{
			name: "Create a free-form invoice",
			fields: fields{
				path:     "/invoices",
				method:   "POST",
				body:     "invoice/create/body_1.json",
				response: "invoice/create/response_1.json",
			},
			args: args{
				ctx: context.Background(),
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
			fields: fields{
				path:     "/invoices",
				method:   "POST",
				body:     "invoice/create/body_2.json",
				response: "invoice/create/response_2.json",
			},
			args: args{
				ctx: context.Background(),
				data: func() *harvest.InvoiceCreateRequest {
					fromDate := time.Date(2017, 3, 1, 0, 0, 0, 0, time.UTC)
					toDate := time.Date(2017, 3, 31, 0, 0, 0, 0, time.UTC)
					projectIds := []int64{14307913}
					lineItemsImport := harvest.InvoiceLineItemImportRequest{
						ProjectIds: &projectIds,
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			mux.HandleFunc(tt.fields.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.fields.method)
				testFormValues(t, r, values{})
				testBody(t, r, tt.fields.body)
				testWriteResponse(t, w, tt.fields.response)
			})

			invoice, _, err := service.Invoice.Create(tt.args.ctx, tt.args.data)
			assert.Equal(t, tt.want, invoice)

			if tt.wantError != nil {
				assert.EqualError(t, err, tt.wantError.Error())

				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestInvoiceService_DeleteInvoice(t *testing.T) {
	t.Parallel()

	type fields struct {
		path     string
		method   string
		body     string
		response string
	}

	type args struct {
		ctx  context.Context
		data int64
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *harvest.Invoice
		wantError error
	}{
		{
			name: "Delete an invoice",
			fields: fields{
				path:     "/invoices/13150453",
				method:   "DELETE",
				body:     "invoice/delete/body_1.json",
				response: "invoice/delete/response_1.json",
			},
			args: args{
				ctx:  context.Background(),
				data: 13150453,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			mux.HandleFunc(tt.fields.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.fields.method)
				testFormValues(t, r, values{})
				testBody(t, r, tt.fields.body)
				testWriteResponse(t, w, tt.fields.response)
			})

			_, err := service.Invoice.Delete(tt.args.ctx, tt.args.data)

			if tt.wantError != nil {
				assert.EqualError(t, err, tt.wantError.Error())

				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestInvoiceService_GetInvoice(t *testing.T) {
	t.Parallel()

	type fields struct {
		path     string
		method   string
		body     string
		response string
	}

	type args struct {
		ctx  context.Context
		data int64
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *harvest.Invoice
		wantError error
	}{
		{
			name: "get an invoice",
			fields: fields{
				path:     "/invoices/13150378",
				method:   "GET",
				body:     "invoice/get/body_1.json",
				response: "invoice/get/response_1.json",
			},
			args: args{
				ctx:  context.Background(),
				data: 13150378,
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			mux.HandleFunc(tt.fields.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.fields.method)
				testFormValues(t, r, values{})
				testBody(t, r, tt.fields.body)
				testWriteResponse(t, w, tt.fields.response)
			})

			invoice, _, err := service.Invoice.Get(tt.args.ctx, tt.args.data)
			assert.Equal(t, tt.want, invoice)

			if tt.wantError != nil {
				assert.EqualError(t, err, tt.wantError.Error())

				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestInvoiceService_ListInvoices(t *testing.T) {
	t.Parallel()

	type fields struct {
		path     string
		method   string
		body     string
		response string
	}

	type args struct {
		ctx  context.Context
		data *harvest.InvoiceListOptions
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *harvest.InvoiceList
		wantError error
	}{
		{
			name: "get invoice list",
			fields: fields{
				path:     "/invoices",
				method:   "GET",
				body:     "invoice/list/body_1.json",
				response: "invoice/list/response_1.json",
			},
			args: args{
				ctx:  context.Background(),
				data: &harvest.InvoiceListOptions{},
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			mux.HandleFunc(tt.fields.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.fields.method)
				testFormValues(t, r, values{})
				testBody(t, r, tt.fields.body)
				testWriteResponse(t, w, tt.fields.response)
			})

			invoices, _, err := service.Invoice.List(tt.args.ctx, tt.args.data)
			assert.Equal(t, tt.want, invoices)

			if tt.wantError != nil {
				assert.EqualError(t, err, tt.wantError.Error())

				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestInvoiceService_UpdateInvoice(t *testing.T) {
	t.Parallel()

	type fields struct {
		path     string
		method   string
		body     string
		response string
	}

	type args struct {
		ctx  context.Context
		id   int64
		data *harvest.InvoiceUpdateRequest
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *harvest.Invoice
		wantError error
	}{
		{
			name: "update an invoice",
			fields: fields{
				path:     "/invoices/13150453",
				method:   "PATCH",
				body:     "invoice/update/body_1.json",
				response: "invoice/update/response_1.json",
			},
			args: args{
				ctx: context.Background(),
				id:  13150453,
				data: func() *harvest.InvoiceUpdateRequest {
					return &harvest.InvoiceUpdateRequest{
						PurchaseOrder: harvest.String("2345"),
					}
				}(),
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			mux.HandleFunc(tt.fields.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.fields.method)
				testFormValues(t, r, values{})
				testBody(t, r, tt.fields.body)
				testWriteResponse(t, w, tt.fields.response)
			})

			invoice, _, err := service.Invoice.Update(tt.args.ctx, tt.args.id, tt.args.data)
			assert.Equal(t, tt.want, invoice)

			if tt.wantError != nil {
				assert.EqualError(t, err, tt.wantError.Error())

				return
			}

			assert.NoError(t, err)
		})
	}
}
