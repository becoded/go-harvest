package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestInvoice_ListInvoiceMessages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		invoiceID int64
		setupMock func(mux *http.ServeMux)
		want      *harvest.InvoiceMessageList
		wantErr   bool
	}{
		{
			name:      "Valid List",
			invoiceID: 13150403,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150403/messages", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testWriteResponse(t, w, "invoice_message/list/response_1.json")
				})
			},
			want: &harvest.InvoiceMessageList{
				InvoiceMessages: []*harvest.InvoiceMessage{
					{
						ID:                         harvest.Int64(27835209),
						SentBy:                     harvest.String("Bob Powell"),
						SentByEmail:                harvest.String("bobpowell@example.com"),
						SentFrom:                   harvest.String("Bob Powell"),
						SentFromEmail:              harvest.String("bobpowell@example.com"),
						IncludeLinkToClientInvoice: harvest.Bool(false),
						SendMeACopy:                harvest.Bool(false),
						ThankYou:                   harvest.Bool(false),
						Reminder:                   harvest.Bool(false),
						AttachPdf:                  harvest.Bool(true),
						Recipients: &[]harvest.InvoiceMessageRecipient{
							{
								Name:  harvest.String("Richard Roe"),
								Email: harvest.String("richardroe@example.com"),
							},
						},
						Subject: harvest.String("Past due invoice reminder: #1001 from API Examples"),
						Body: harvest.String("Dear Customer,\r\n\r\n" +
							"This is a friendly reminder to let you know that Invoice 1001 is 144 days past due. " +
							"If you have already sent the payment, please disregard this message. " +
							"If not, we would appreciate your prompt attention to this matter.\r\n\r\n" +
							"Thank you for your business.\r\n\r\n" +
							"Cheers,\r\n" +
							"API Examples"),
						CreatedAt: harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 15, 6, 0, time.UTC)),
						UpdatedAt: harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 15, 6, 0, time.UTC)),
					},
					{
						ID:                         harvest.Int64(27835207),
						SentBy:                     harvest.String("Bob Powell"),
						SentByEmail:                harvest.String("bobpowell@example.com"),
						SentFrom:                   harvest.String("Bob Powell"),
						SentFromEmail:              harvest.String("bobpowell@example.com"),
						IncludeLinkToClientInvoice: harvest.Bool(false),
						SendMeACopy:                harvest.Bool(true),
						ThankYou:                   harvest.Bool(false),
						Reminder:                   harvest.Bool(false),
						AttachPdf:                  harvest.Bool(true),
						Recipients: &[]harvest.InvoiceMessageRecipient{
							{
								Name:  harvest.String("Richard Roe"),
								Email: harvest.String("richardroe@example.com"),
							},
							{
								Name:  harvest.String("Bob Powell"),
								Email: harvest.String("bobpowell@example.com"),
							},
						},
						Subject: harvest.String("Invoice #1001 from API Examples"),
						Body: harvest.String("---------------------------------------------\r\n" +
							"Invoice Summary\r\n" +
							"---------------------------------------------\r\n" +
							"Invoice ID: 1001\r\n" +
							"Issue Date: 04/01/2017\r\n" +
							"Client: 123 Industries\r\n" +
							"P.O. Number: \r\n" +
							"Amount: €288.90\r\n" +
							"Due: 04/01/2017 (upon receipt)\r\n\r\n" +
							"The detailed invoice is attached as a PDF.\r\n\r\n" +
							"Thank you!\r\n" +
							"---------------------------------------------"),
						CreatedAt: harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 14, 49, 0, time.UTC)),
						UpdatedAt: harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 14, 49, 0, time.UTC)),
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
						First:    harvest.String("https://api.harvestapp.com/api/v2/invoices/13150403/messages?page=1&per_page=2000"),
						Next:     nil,
						Previous: nil,
						Last:     harvest.String("https://api.harvestapp.com/v2/invoices/13150403/messages?page=1&per_page=2000"),
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

			got, _, err := service.Invoice.ListInvoiceMessages(
				context.Background(),
				tt.invoiceID,
				&harvest.InvoiceMessageListOptions{},
			)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestInvoice_CreateInvoiceMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		invoiceID int64
		request   *harvest.InvoiceMessageCreateRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.InvoiceMessage
		wantErr   bool
	}{
		{
			name:      "Valid Creation",
			invoiceID: 13150403,
			request: &harvest.InvoiceMessageCreateRequest{
				Subject:     harvest.String("Invoice #1001"),
				Body:        harvest.String("The invoice is attached below."),
				AttachPdf:   harvest.Bool(true),
				SendMeACopy: harvest.Bool(true),
				Recipients: &[]harvest.InvoiceMessageRecipient{
					{
						Name:  harvest.String("Richard Roe"),
						Email: harvest.String("richardroe@example.com"),
					},
				},
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150403/messages", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testWriteResponse(t, w, "invoice_message/create/response_1.json")
				})
			},
			want: &harvest.InvoiceMessage{
				ID:                         harvest.Int64(27835324),
				SentBy:                     harvest.String("Bob Powell"),
				SentByEmail:                harvest.String("bobpowell@example.com"),
				SentFrom:                   harvest.String("Bob Powell"),
				SentFromEmail:              harvest.String("bobpowell@example.com"),
				AttachPdf:                  harvest.Bool(true),
				SendMeACopy:                harvest.Bool(true),
				IncludeLinkToClientInvoice: harvest.Bool(false),
				ThankYou:                   harvest.Bool(false),
				Reminder:                   harvest.Bool(false),
				Recipients: &[]harvest.InvoiceMessageRecipient{
					{
						Name:  harvest.String("Richard Roe"),
						Email: harvest.String("richardroe@example.com"),
					}, {
						Name:  harvest.String("Bob Powell"),
						Email: harvest.String("bobpowell@example.com"),
					},
				},
				Subject:   harvest.String("Invoice #1001"),
				Body:      harvest.String("The invoice is attached below."),
				CreatedAt: harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 25, 59, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 25, 59, 0, time.UTC)),
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

			got, _, err := service.Invoice.CreateInvoiceMessage(context.Background(), tt.invoiceID, tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestInvoice_DeleteInvoiceMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		invoiceID int64
		messageID int64
		setupMock func(mux *http.ServeMux)
		wantErr   bool
	}{
		{
			name:      "Valid Deletion",
			invoiceID: 13150403,
			messageID: 27835324,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150403/messages/27835324", func(w http.ResponseWriter, r *http.Request) {
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

			_, err := service.Invoice.DeleteInvoiceMessage(context.Background(), tt.invoiceID, tt.messageID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInvoice_MarkAsSent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		invoiceID int64
		setupMock func(mux *http.ServeMux)
		want      *harvest.InvoiceMessage
		wantErr   bool
	}{
		{
			name:      "Valid Mark As Sent",
			invoiceID: 13150403,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150403/messages", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testWriteResponse(t, w, "invoice_message/mark_as_sent/response_1.json")
				})
			},
			want: &harvest.InvoiceMessage{
				ID:                         harvest.Int64(27835325),
				SentBy:                     harvest.String("Bob Powell"),
				SentByEmail:                harvest.String("bobpowell@example.com"),
				SentFrom:                   harvest.String("Bob Powell"),
				SentFromEmail:              harvest.String("bobpowell@example.com"),
				EventType:                  harvest.String("send"),
				IncludeLinkToClientInvoice: harvest.Bool(false),
				SendMeACopy:                harvest.Bool(false),
				ThankYou:                   harvest.Bool(false),
				Reminder:                   harvest.Bool(false),
				AttachPdf:                  harvest.Bool(false),
				Recipients:                 &[]harvest.InvoiceMessageRecipient{},
				CreatedAt:                  harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 25, 59, 0, time.UTC)),
				UpdatedAt:                  harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 25, 59, 0, time.UTC)),
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

			got, _, err := service.Invoice.MarkAsSent(context.Background(), tt.invoiceID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestInvoice_MarkAsDraft(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		invoiceID int64
		setupMock func(mux *http.ServeMux)
		want      *harvest.InvoiceMessage
		wantErr   bool
	}{
		{
			name:      "Valid Mark As Draft",
			invoiceID: 13150403,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150403/messages", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testWriteResponse(t, w, "invoice_message/mark_as_draft/response_1.json")
				})
			},
			want: &harvest.InvoiceMessage{
				ID:                         harvest.Int64(27835328),
				SentBy:                     harvest.String("Bob Powell"),
				SentByEmail:                harvest.String("bobpowell@example.com"),
				SentFrom:                   harvest.String("Bob Powell"),
				SentFromEmail:              harvest.String("bobpowell@example.com"),
				EventType:                  harvest.String("draft"),
				IncludeLinkToClientInvoice: harvest.Bool(false),
				SendMeACopy:                harvest.Bool(false),
				ThankYou:                   harvest.Bool(false),
				Reminder:                   harvest.Bool(false),
				AttachPdf:                  harvest.Bool(false),
				Recipients:                 &[]harvest.InvoiceMessageRecipient{},
				CreatedAt:                  harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 25, 59, 0, time.UTC)),
				UpdatedAt:                  harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 25, 59, 0, time.UTC)),
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

			got, _, err := service.Invoice.MarkAsDraft(context.Background(), tt.invoiceID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestInvoice_MarkAsClosed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		invoiceID int64
		setupMock func(mux *http.ServeMux)
		want      *harvest.InvoiceMessage
		wantErr   bool
	}{
		{
			name:      "Valid Mark As Closed",
			invoiceID: 13150403,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150403/messages", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testWriteResponse(t, w, "invoice_message/mark_as_closed/response_1.json")
				})
			},
			want: &harvest.InvoiceMessage{
				ID:                         harvest.Int64(27835326),
				SentBy:                     harvest.String("Bob Powell"),
				SentByEmail:                harvest.String("bobpowell@example.com"),
				SentFrom:                   harvest.String("Bob Powell"),
				SentFromEmail:              harvest.String("bobpowell@example.com"),
				EventType:                  harvest.String("close"),
				IncludeLinkToClientInvoice: harvest.Bool(false),
				SendMeACopy:                harvest.Bool(false),
				ThankYou:                   harvest.Bool(false),
				Reminder:                   harvest.Bool(false),
				AttachPdf:                  harvest.Bool(false),
				Recipients:                 &[]harvest.InvoiceMessageRecipient{},
				CreatedAt:                  harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 25, 59, 0, time.UTC)),
				UpdatedAt:                  harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 25, 59, 0, time.UTC)),
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

			got, _, err := service.Invoice.MarkAsClosed(context.Background(), tt.invoiceID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestInvoice_MarkAsReopen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		invoiceID int64
		setupMock func(mux *http.ServeMux)
		want      *harvest.InvoiceMessage
		wantErr   bool
	}{
		{
			name:      "Valid Mark As Reopen",
			invoiceID: 13150403,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150403/messages", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testWriteResponse(t, w, "invoice_message/mark_as_reopen/response_1.json")
				})
			},
			want: &harvest.InvoiceMessage{
				ID:                         harvest.Int64(27835327),
				SentBy:                     harvest.String("Bob Powell"),
				SentByEmail:                harvest.String("bobpowell@example.com"),
				SentFrom:                   harvest.String("Bob Powell"),
				SentFromEmail:              harvest.String("bobpowell@example.com"),
				EventType:                  harvest.String("re-open"),
				IncludeLinkToClientInvoice: harvest.Bool(false),
				SendMeACopy:                harvest.Bool(false),
				ThankYou:                   harvest.Bool(false),
				Reminder:                   harvest.Bool(false),
				AttachPdf:                  harvest.Bool(false),
				Recipients:                 &[]harvest.InvoiceMessageRecipient{},
				CreatedAt:                  harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 25, 59, 0, time.UTC)),
				UpdatedAt:                  harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 25, 59, 0, time.UTC)),
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

			got, _, err := service.Invoice.MarkAsReopen(context.Background(), tt.invoiceID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
