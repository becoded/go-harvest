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
		{
			name:      "Error Fetching List",
			invoiceID: 13150403,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150403/messages", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Invoice.ListInvoiceMessages(
				context.Background(),
				tt.invoiceID,
				&harvest.InvoiceMessageListOptions{},
			)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
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
		{
			name:      "Error Creating Message",
			invoiceID: 13150403,
			request: &harvest.InvoiceMessageCreateRequest{
				Subject: harvest.String("Invoice #1001"),
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

			got, _, err := service.Invoice.CreateInvoiceMessage(context.Background(), tt.invoiceID, tt.request)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
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
		{
			name:      "Error Deleting Message",
			invoiceID: 13150403,
			messageID: 27835324,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150403/messages/27835324", func(w http.ResponseWriter, r *http.Request) {
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
		{
			name:      "Error Marking As Sent",
			invoiceID: 13150403,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150403/messages", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Invoice.MarkAsSent(context.Background(), tt.invoiceID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
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
		{
			name:      "Error Marking As Draft",
			invoiceID: 13150403,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150403/messages", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Invoice.MarkAsDraft(context.Background(), tt.invoiceID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
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
		{
			name:      "Error Marking As Closed",
			invoiceID: 13150403,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150403/messages", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Invoice.MarkAsClosed(context.Background(), tt.invoiceID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
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
		{
			name:      "Error Marking As Reopen",
			invoiceID: 13150403,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoices/13150403/messages", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Invoice.MarkAsReopen(context.Background(), tt.invoiceID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestInvoiceMessage_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.InvoiceMessage
		want string
	}{
		{
			name: "InvoiceMessage with all fields",
			in: harvest.InvoiceMessage{
				ID:            harvest.Int64(27835209),
				SentBy:        harvest.String("Bob Powell"),
				SentByEmail:   harvest.String("bobpowell@example.com"),
				SentFrom:      harvest.String("Bob Powell"),
				SentFromEmail: harvest.String("bobpowell@example.com"),
				Recipients: &[]harvest.InvoiceMessageRecipient{
					{
						Name:  harvest.String("Richard Roe"),
						Email: harvest.String("richardroe@example.com"),
					},
				},
				Subject:                    harvest.String("Invoice #1001"),
				Body:                       harvest.String("The invoice is attached below."),
				IncludeLinkToClientInvoice: harvest.Bool(true),
				AttachPdf:                  harvest.Bool(true),
				SendMeACopy:                harvest.Bool(true),
				ThankYou:                   harvest.Bool(false),
				EventType:                  harvest.String("send"),
				Reminder:                   harvest.Bool(false),
				SendReminderOn: &harvest.Date{
					Time: time.Date(2017, 9, 1, 0, 0, 0, 0, time.UTC),
				},
				CreatedAt: harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 15, 6, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 15, 6, 0, time.UTC)),
			},
			want: `harvest.InvoiceMessage{ID:27835209, SentBy:"Bob Powell", SentByEmail:"bobpowell@example.com", SentFrom:"Bob Powell", SentFromEmail:"bobpowell@example.com", Recipients:[harvest.InvoiceMessageRecipient{Name:"Richard Roe", Email:"richardroe@example.com"}], Subject:"Invoice #1001", Body:"The invoice is attached below.", IncludeLinkToClientInvoice:true, AttachPdf:true, SendMeACopy:true, ThankYou:false, EventType:"send", Reminder:false, SendReminderOn:harvest.Date{{2017-09-01 00:00:00 +0000 UTC}}, CreatedAt:time.Time{2017-08-23 22:15:06 +0000 UTC}, UpdatedAt:time.Time{2017-08-23 22:15:06 +0000 UTC}}`, //nolint: lll
		},
		{
			name: "InvoiceMessage with minimal fields",
			in: harvest.InvoiceMessage{
				ID:      harvest.Int64(27835324),
				Subject: harvest.String("Invoice #1001"),
			},
			want: `harvest.InvoiceMessage{ID:27835324, Subject:"Invoice #1001"}`,
		},
		{
			name: "InvoiceMessage with multiple recipients",
			in: harvest.InvoiceMessage{
				ID: harvest.Int64(27835207),
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
			},
			want: `harvest.InvoiceMessage{ID:27835207, Recipients:[harvest.InvoiceMessageRecipient{Name:"Richard Roe", Email:"richardroe@example.com"} harvest.InvoiceMessageRecipient{Name:"Bob Powell", Email:"bobpowell@example.com"}], Subject:"Invoice #1001 from API Examples"}`, //nolint: lll
		},
		{
			name: "InvoiceMessage with event type",
			in: harvest.InvoiceMessage{
				ID:                         harvest.Int64(27835325),
				EventType:                  harvest.String("close"),
				IncludeLinkToClientInvoice: harvest.Bool(false),
				SendMeACopy:                harvest.Bool(false),
				ThankYou:                   harvest.Bool(false),
				Reminder:                   harvest.Bool(false),
				AttachPdf:                  harvest.Bool(false),
			},
			want: `harvest.InvoiceMessage{ID:27835325, IncludeLinkToClientInvoice:false, AttachPdf:false, SendMeACopy:false, ThankYou:false, EventType:"close", Reminder:false}`, //nolint: lll
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

func TestInvoiceMessageList_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.InvoiceMessageList
		want string
	}{
		{
			name: "InvoiceMessageList with multiple messages",
			in: harvest.InvoiceMessageList{
				InvoiceMessages: []*harvest.InvoiceMessage{
					{
						ID:            harvest.Int64(27835209),
						SentBy:        harvest.String("Bob Powell"),
						SentByEmail:   harvest.String("bobpowell@example.com"),
						Subject:       harvest.String("Invoice #1001"),
						AttachPdf:     harvest.Bool(true),
						SendMeACopy:   harvest.Bool(false),
						ThankYou:      harvest.Bool(false),
						Reminder:      harvest.Bool(false),
						CreatedAt:     harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 15, 6, 0, time.UTC)),
						UpdatedAt:     harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 15, 6, 0, time.UTC)),
					},
					{
						ID:          harvest.Int64(27835207),
						SentBy:      harvest.String("Bob Powell"),
						SentByEmail: harvest.String("bobpowell@example.com"),
						Subject:     harvest.String("Invoice #1002"),
						AttachPdf:   harvest.Bool(false),
						SendMeACopy: harvest.Bool(true),
						ThankYou:    harvest.Bool(false),
						Reminder:    harvest.Bool(false),
						CreatedAt:   harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 14, 49, 0, time.UTC)),
						UpdatedAt:   harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 14, 49, 0, time.UTC)),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(2000),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.InvoiceMessageList{InvoiceMessages:[harvest.InvoiceMessage{ID:27835209, SentBy:"Bob Powell", SentByEmail:"bobpowell@example.com", Subject:"Invoice #1001", AttachPdf:true, SendMeACopy:false, ThankYou:false, Reminder:false, CreatedAt:time.Time{2017-08-23 22:15:06 +0000 UTC}, UpdatedAt:time.Time{2017-08-23 22:15:06 +0000 UTC}} harvest.InvoiceMessage{ID:27835207, SentBy:"Bob Powell", SentByEmail:"bobpowell@example.com", Subject:"Invoice #1002", AttachPdf:false, SendMeACopy:true, ThankYou:false, Reminder:false, CreatedAt:time.Time{2017-08-23 22:14:49 +0000 UTC}, UpdatedAt:time.Time{2017-08-23 22:14:49 +0000 UTC}}], Pagination:harvest.Pagination{PerPage:2000, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
		},
		{
			name: "InvoiceMessageList with single message",
			in: harvest.InvoiceMessageList{
				InvoiceMessages: []*harvest.InvoiceMessage{
					{
						ID:      harvest.Int64(27835324),
						Subject: harvest.String("Invoice #1001"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(2000),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.InvoiceMessageList{InvoiceMessages:[harvest.InvoiceMessage{ID:27835324, Subject:"Invoice #1001"}], Pagination:harvest.Pagination{PerPage:2000, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
		},
		{
			name: "Empty InvoiceMessageList",
			in: harvest.InvoiceMessageList{
				InvoiceMessages: []*harvest.InvoiceMessage{},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(2000),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(0),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.InvoiceMessageList{InvoiceMessages:[], Pagination:harvest.Pagination{PerPage:2000, TotalPages:1, TotalEntries:0, Page:1}}`, //nolint: lll
		},
		{
			name: "InvoiceMessageList with Links",
			in: harvest.InvoiceMessageList{
				InvoiceMessages: []*harvest.InvoiceMessage{
					{
						ID:      harvest.Int64(27835324),
						Subject: harvest.String("Invoice #1001"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(3),
					TotalEntries: harvest.Int(250),
					Page:         harvest.Int(2),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/invoices/13150403/messages?page=1&per_page=100"),
						Next:     harvest.String("https://api.harvestapp.com/v2/invoices/13150403/messages?page=3&per_page=100"),
						Previous: harvest.String("https://api.harvestapp.com/v2/invoices/13150403/messages?page=1&per_page=100"),
						Last:     harvest.String("https://api.harvestapp.com/v2/invoices/13150403/messages?page=3&per_page=100"),
					},
				},
			},
			want: `harvest.InvoiceMessageList{InvoiceMessages:[harvest.InvoiceMessage{ID:27835324, Subject:"Invoice #1001"}], Pagination:harvest.Pagination{PerPage:100, TotalPages:3, TotalEntries:250, Page:2, Links:harvest.PageLinks{First:"https://api.harvestapp.com/v2/invoices/13150403/messages?page=1&per_page=100", Next:"https://api.harvestapp.com/v2/invoices/13150403/messages?page=3&per_page=100", Previous:"https://api.harvestapp.com/v2/invoices/13150403/messages?page=1&per_page=100", Last:"https://api.harvestapp.com/v2/invoices/13150403/messages?page=3&per_page=100"}}}`, //nolint: lll
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
