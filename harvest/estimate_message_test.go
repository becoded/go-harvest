package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestEstimate_ListEstimateMessages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		estimateID int64
		setupMock  func(mux *http.ServeMux)
		want       *harvest.EstimateMessageList
		wantErr    bool
	}{
		{
			name:       "Valid List",
			estimateID: 123,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimates/123/messages", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testBody(t, r, "estimate_message/list/body_1.json")
					testWriteResponse(t, w, "estimate_message/list/response_1.json")
				})
			},
			want: &harvest.EstimateMessageList{
				EstimateMessages: []*harvest.EstimateMessage{
					{
						ID:            harvest.Int64(2666236),
						SentBy:        harvest.String("Bob Powell"),
						SentByEmail:   harvest.String("bobpowell@example.com"),
						SentFrom:      harvest.String("Bob Powell"),
						SentFromEmail: harvest.String("bobpowell@example.com"),
						SendMeACopy:   harvest.Bool(true),
						CreatedAt:     harvest.TimeTimeP(time.Date(2017, 8, 25, 21, 23, 40, 0, time.UTC)),
						UpdatedAt:     harvest.TimeTimeP(time.Date(2017, 8, 25, 21, 23, 40, 0, time.UTC)),
						Recipients: &[]harvest.EstimateMessageRecipient{
							{
								Name:  harvest.String("Richard Roe"),
								Email: harvest.String("richardroe@example.com"),
							},
							{
								Name:  harvest.String("Bob Powell"),
								Email: harvest.String("bobpowell@example.com"),
							},
						},
						Subject: harvest.String("Estimate #1001 from API Examples"),
						Body: harvest.String("---------------------------------------------\r\n" +
							"Estimate Summary\r\n" +
							"---------------------------------------------\r\n" +
							"Estimate ID: 1001\r\n" +
							"Estimate Date: 06/01/2017\r\n" +
							"Client: 123 Industries\r\n" +
							"P.O. Number: 5678\r\n" +
							"Amount: $9,630.00\r\n\r\n" +
							"You can view the estimate here:\r\n\r\n" +
							"%estimate_url%\r\n\r\n" +
							"Thank you!\r\n---------------------------------------------"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(2000),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
					Links: &harvest.PageLinks{
						First: harvest.String(
							"https://api.harvestapp.com/v2/estimates/1439818/messages?page=1&per_page=2000",
						),
						Next:     nil,
						Previous: nil,
						Last: harvest.String(
							"https://api.harvestapp.com/v2/estimates/1439818/messages?page=1&per_page=2000",
						),
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "Error Fetching List",
			estimateID: 123,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimates/123/messages", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Estimate.ListEstimateMessages(
				context.Background(),
				tt.estimateID,
				&harvest.EstimateMessageListOptions{},
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

func TestEstimate_CreateEstimateMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		estimateID int64
		request    *harvest.EstimateMessageCreateRequest
		setupMock  func(mux *http.ServeMux)
		want       *harvest.EstimateMessage
		wantErr    bool
	}{
		{
			name:       "Valid Creation",
			estimateID: 1439818,
			request: &harvest.EstimateMessageCreateRequest{
				Subject:     harvest.String("Estimate #1001"),
				Body:        harvest.String("Here is our estimate."),
				SendMeACopy: harvest.Bool(true),
				Recipients: &[]harvest.EstimateMessageRecipientCreateRequest{
					{
						Name:  harvest.String("Richard Roe"),
						Email: harvest.String("richardroe@example.com"),
					},
				},
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimates/1439818/messages", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testBody(t, r, "estimate_message/create/body_1.json")
					testWriteResponse(t, w, "estimate_message/create/response_1.json")
				})
			},
			want: &harvest.EstimateMessage{
				ID:            harvest.Int64(2666240),
				SentBy:        harvest.String("Bob Powell"),
				SentByEmail:   harvest.String("bobpowell@example.com"),
				SentFrom:      harvest.String("Bob Powell"),
				SentFromEmail: harvest.String("bobpowell@example.com"),
				SendMeACopy:   harvest.Bool(true),
				CreatedAt:     harvest.TimeTimeP(time.Date(2017, 8, 25, 21, 27, 52, 0, time.UTC)),
				UpdatedAt:     harvest.TimeTimeP(time.Date(2017, 8, 25, 21, 27, 52, 0, time.UTC)),
				Recipients: &[]harvest.EstimateMessageRecipient{
					{
						Name:  harvest.String("Richard Roe"),
						Email: harvest.String("richardroe@example.com"),
					},
				},
				Subject:   harvest.String("Estimate #1001"),
				Body:      harvest.String("Here is our estimate."),
				EventType: nil,
			},
			wantErr: false,
		},
		{
			name:       "Invalid Creation - Missing Subject",
			estimateID: 1439818,
			request: &harvest.EstimateMessageCreateRequest{
				Body: harvest.String("Here is our estimate."),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimates/1439818/messages", func(w http.ResponseWriter, _ *http.Request) {
					http.Error(w, `{"message":"Subject is required"}`, http.StatusBadRequest)
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

			got, _, err := service.Estimate.CreateEstimateMessage(context.Background(), tt.estimateID, tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestEstimate_DeleteEstimateMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		estimateID int64
		messageID  int64
		setupMock  func(mux *http.ServeMux)
		wantErr    bool
	}{
		{
			name:       "Valid Deletion",
			estimateID: 123,
			messageID:  456,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimates/123/messages/456", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					w.WriteHeader(http.StatusNoContent)
				})
			},
			wantErr: false,
		},
		{
			name:       "Message Not Found",
			estimateID: 123,
			messageID:  999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimates/123/messages/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					http.Error(w, `{"message":"Message not found"}`, http.StatusNotFound)
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

			_, err := service.Estimate.DeleteEstimateMessage(context.Background(), tt.estimateID, tt.messageID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEstimate_MarkAsSent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		estimateID int64
		setupMock  func(mux *http.ServeMux)
		want       *harvest.EstimateMessage
		wantErr    bool
	}{
		{
			name:       "Valid Mark As Sent",
			estimateID: 123,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimates/123/messages", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testBody(t, r, "estimate_message/mark_as_sent/body_1.json")
					testWriteResponse(t, w, "estimate_message/mark_as_sent/response_1.json")
				})
			},
			want: &harvest.EstimateMessage{
				ID:            harvest.Int64(2666241),
				SentBy:        harvest.String("Bob Powell"),
				SentByEmail:   harvest.String("bobpowell@example.com"),
				SentFrom:      harvest.String("Bob Powell"),
				SentFromEmail: harvest.String("bobpowell@example.com"),
				SendMeACopy:   harvest.Bool(false),
				CreatedAt:     harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 25, 59, 0, time.UTC)),
				UpdatedAt:     harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 25, 59, 0, time.UTC)),
				EventType:     harvest.String("send"),
				Recipients:    &[]harvest.EstimateMessageRecipient{},
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

			got, _, err := service.Estimate.MarkAsSent(context.Background(), tt.estimateID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestEstimate_MarkAsAccepted(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		estimateID int64
		setupMock  func(mux *http.ServeMux)
		want       *harvest.EstimateMessage
		wantErr    bool
	}{
		{
			name:       "Valid Mark As Accepted",
			estimateID: 123,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimates/123/messages", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testBody(t, r, "estimate_message/mark_as_accepted/body_1.json")
					testWriteResponse(t, w, "estimate_message/mark_as_accepted/response_1.json")
				})
			},
			want: &harvest.EstimateMessage{
				ID:            harvest.Int64(2666244),
				SentBy:        harvest.String("Bob Powell"),
				SentByEmail:   harvest.String("bobpowell@example.com"),
				SentFrom:      harvest.String("Bob Powell"),
				SentFromEmail: harvest.String("bobpowell@example.com"),
				SendMeACopy:   harvest.Bool(false),
				CreatedAt:     harvest.TimeTimeP(time.Date(2017, 8, 25, 21, 31, 55, 0, time.UTC)),
				UpdatedAt:     harvest.TimeTimeP(time.Date(2017, 8, 25, 21, 31, 55, 0, time.UTC)),
				EventType:     harvest.String("accept"),
				Recipients:    &[]harvest.EstimateMessageRecipient{},
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

			got, _, err := service.Estimate.MarkAsAccepted(context.Background(), tt.estimateID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestEstimate_MarkAsDeclined(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		estimateID int64
		setupMock  func(mux *http.ServeMux)
		want       *harvest.EstimateMessage
		wantErr    bool
	}{
		{
			name:       "Valid Mark As Declined",
			estimateID: 123,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimates/123/messages", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testBody(t, r, "estimate_message/mark_as_declined/body_1.json")
					testWriteResponse(t, w, "estimate_message/mark_as_declined/response_1.json")
				})
			},
			want: &harvest.EstimateMessage{
				ID:            harvest.Int64(2666245),
				SentBy:        harvest.String("Bob Powell"),
				SentByEmail:   harvest.String("bobpowell@example.com"),
				SentFrom:      harvest.String("Bob Powell"),
				SentFromEmail: harvest.String("bobpowell@example.com"),
				SendMeACopy:   harvest.Bool(false),
				CreatedAt:     harvest.TimeTimeP(time.Date(2017, 8, 25, 21, 31, 55, 0, time.UTC)),
				UpdatedAt:     harvest.TimeTimeP(time.Date(2017, 8, 25, 21, 31, 55, 0, time.UTC)),
				EventType:     harvest.String("decline"),
				Recipients:    &[]harvest.EstimateMessageRecipient{},
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

			got, _, err := service.Estimate.MarkAsDeclined(context.Background(), tt.estimateID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestEstimate_MarkAsReopen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		estimateID int64
		setupMock  func(mux *http.ServeMux)
		want       *harvest.EstimateMessage
		wantErr    bool
	}{
		{
			name:       "Valid Mark As Reopen",
			estimateID: 123,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimates/123/messages", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testBody(t, r, "estimate_message/mark_as_reopen/body_1.json")
					testWriteResponse(t, w, "estimate_message/mark_as_reopen/response_1.json")
				})
			},
			want: &harvest.EstimateMessage{
				ID:            harvest.Int64(2666246),
				SentBy:        harvest.String("Bob Powell"),
				SentByEmail:   harvest.String("bobpowell@example.com"),
				SentFrom:      harvest.String("Bob Powell"),
				SentFromEmail: harvest.String("bobpowell@example.com"),
				SendMeACopy:   harvest.Bool(false),
				CreatedAt:     harvest.TimeTimeP(time.Date(2017, 8, 25, 21, 31, 55, 0, time.UTC)),
				UpdatedAt:     harvest.TimeTimeP(time.Date(2017, 8, 25, 21, 31, 55, 0, time.UTC)),
				EventType:     harvest.String("re-open"),
				Recipients:    &[]harvest.EstimateMessageRecipient{},
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

			got, _, err := service.Estimate.MarkAsReopen(context.Background(), tt.estimateID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
