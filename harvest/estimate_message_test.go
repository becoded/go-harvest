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
		{
			name:       "Error Marking As Sent",
			estimateID: 123,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimates/123/messages", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Estimate.MarkAsSent(context.Background(), tt.estimateID)
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
		{
			name:       "Error Marking As Accepted",
			estimateID: 123,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimates/123/messages", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Estimate.MarkAsAccepted(context.Background(), tt.estimateID)
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
		{
			name:       "Error Marking As Declined",
			estimateID: 123,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimates/123/messages", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Estimate.MarkAsDeclined(context.Background(), tt.estimateID)
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
		{
			name:       "Error Marking As Reopen",
			estimateID: 123,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimates/123/messages", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Estimate.MarkAsReopen(context.Background(), tt.estimateID)
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

func TestEstimateMessage_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.EstimateMessage
		want string
	}{
		{
			name: "EstimateMessage with all fields",
			in: harvest.EstimateMessage{
				ID:            harvest.Int64(2666236),
				SentBy:        harvest.String("Bob Powell"),
				SentByEmail:   harvest.String("bobpowell@example.com"),
				SentFrom:      harvest.String("Bob Powell"),
				SentFromEmail: harvest.String("bobpowell@example.com"),
				Recipients: &[]harvest.EstimateMessageRecipient{
					{
						Name:  harvest.String("Richard Roe"),
						Email: harvest.String("richardroe@example.com"),
					},
				},
				Subject:     harvest.String("Estimate #1001 from API Examples"),
				Body:        harvest.String("Please review the estimate."),
				SendMeACopy: harvest.Bool(true),
				EventType:   harvest.String("send"),
				CreatedAt:   harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 25, 59, 0, time.UTC)),
				UpdatedAt:   harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 25, 59, 0, time.UTC)),
			},
			want: `harvest.EstimateMessage{ID:2666236, SentBy:"Bob Powell", SentByEmail:"bobpowell@example.com", SentFrom:"Bob Powell", SentFromEmail:"bobpowell@example.com", Recipients:[harvest.EstimateMessageRecipient{Name:"Richard Roe", Email:"richardroe@example.com"}], Subject:"Estimate #1001 from API Examples", Body:"Please review the estimate.", SendMeACopy:true, EventType:"send", CreatedAt:time.Time{2017-08-23 22:25:59 +0000 UTC}, UpdatedAt:time.Time{2017-08-23 22:25:59 +0000 UTC}}`, //nolint: lll
		},
		{
			name: "EstimateMessage with minimal fields",
			in: harvest.EstimateMessage{
				ID:      harvest.Int64(999),
				Subject: harvest.String("Test Estimate"),
			},
			want: `harvest.EstimateMessage{ID:999, Subject:"Test Estimate"}`,
		},
		{
			name: "EstimateMessage with multiple recipients",
			in: harvest.EstimateMessage{
				ID: harvest.Int64(2666236),
				Recipients: &[]harvest.EstimateMessageRecipient{
					{
						Name:  harvest.String("Richard Roe"),
						Email: harvest.String("richardroe@example.com"),
					},
					{
						Name:  harvest.String("Jane Doe"),
						Email: harvest.String("janedoe@example.com"),
					},
				},
				Subject:   harvest.String("Estimate #1001"),
				EventType: harvest.String("send"),
			},
			want: `harvest.EstimateMessage{ID:2666236, Recipients:[harvest.EstimateMessageRecipient{Name:"Richard Roe", Email:"richardroe@example.com"} harvest.EstimateMessageRecipient{Name:"Jane Doe", Email:"janedoe@example.com"}], Subject:"Estimate #1001", EventType:"send"}`, //nolint: lll
		},
		{
			name: "EstimateMessage with event type",
			in: harvest.EstimateMessage{
				ID:          harvest.Int64(2666237),
				Subject:     harvest.String("Estimate Accepted"),
				EventType:   harvest.String("accept"),
				SendMeACopy: harvest.Bool(false),
			},
			want: `harvest.EstimateMessage{ID:2666237, Subject:"Estimate Accepted", SendMeACopy:false, EventType:"accept"}`,
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

func TestEstimateMessageList_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.EstimateMessageList
		want string
	}{
		{
			name: "EstimateMessageList with multiple messages",
			in: harvest.EstimateMessageList{
				EstimateMessages: []*harvest.EstimateMessage{
					{
						ID:          harvest.Int64(2666236),
						SentBy:      harvest.String("Bob Powell"),
						SentByEmail: harvest.String("bobpowell@example.com"),
						Subject:     harvest.String("Estimate #1001 from API Examples"),
						Body:        harvest.String("Please review the estimate."),
						EventType:   harvest.String("send"),
						CreatedAt:   harvest.TimeTimeP(time.Date(2017, 8, 23, 22, 25, 59, 0, time.UTC)),
					},
					{
						ID:        harvest.Int64(2666237),
						SentBy:    harvest.String("Bob Powell"),
						Subject:   harvest.String("Estimate Accepted"),
						EventType: harvest.String("accept"),
						CreatedAt: harvest.TimeTimeP(time.Date(2017, 8, 24, 10, 15, 30, 0, time.UTC)),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.EstimateMessageList{EstimateMessages:[harvest.EstimateMessage{ID:2666236, SentBy:"Bob Powell", SentByEmail:"bobpowell@example.com", Subject:"Estimate #1001 from API Examples", Body:"Please review the estimate.", EventType:"send", CreatedAt:time.Time{2017-08-23 22:25:59 +0000 UTC}} harvest.EstimateMessage{ID:2666237, SentBy:"Bob Powell", Subject:"Estimate Accepted", EventType:"accept", CreatedAt:time.Time{2017-08-24 10:15:30 +0000 UTC}}], Pagination:harvest.Pagination{PerPage:100, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
		},
		{
			name: "EstimateMessageList with single message",
			in: harvest.EstimateMessageList{
				EstimateMessages: []*harvest.EstimateMessage{
					{
						ID:      harvest.Int64(999),
						Subject: harvest.String("Test Message"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.EstimateMessageList{EstimateMessages:[harvest.EstimateMessage{ID:999, Subject:"Test Message"}], Pagination:harvest.Pagination{PerPage:100, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
		},
		{
			name: "Empty EstimateMessageList",
			in: harvest.EstimateMessageList{
				EstimateMessages: []*harvest.EstimateMessage{},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(0),
					TotalEntries: harvest.Int(0),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.EstimateMessageList{EstimateMessages:[], Pagination:harvest.Pagination{PerPage:100, TotalPages:0, TotalEntries:0, Page:1}}`, //nolint: lll
		},
		{
			name: "EstimateMessageList with Links",
			in: harvest.EstimateMessageList{
				EstimateMessages: []*harvest.EstimateMessage{
					{
						ID:      harvest.Int64(100),
						Subject: harvest.String("Message"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(50),
					TotalPages:   harvest.Int(3),
					TotalEntries: harvest.Int(150),
					Page:         harvest.Int(2),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/estimates/1/messages?page=1&per_page=50"),
						Next:     harvest.String("https://api.harvestapp.com/v2/estimates/1/messages?page=3&per_page=50"),
						Previous: harvest.String("https://api.harvestapp.com/v2/estimates/1/messages?page=1&per_page=50"),
						Last:     harvest.String("https://api.harvestapp.com/v2/estimates/1/messages?page=3&per_page=50"),
					},
				},
			},
			want: `harvest.EstimateMessageList{EstimateMessages:[harvest.EstimateMessage{ID:100, Subject:"Message"}], Pagination:harvest.Pagination{PerPage:50, TotalPages:3, TotalEntries:150, Page:2, Links:harvest.PageLinks{First:"https://api.harvestapp.com/v2/estimates/1/messages?page=1&per_page=50", Next:"https://api.harvestapp.com/v2/estimates/1/messages?page=3&per_page=50", Previous:"https://api.harvestapp.com/v2/estimates/1/messages?page=1&per_page=50", Last:"https://api.harvestapp.com/v2/estimates/1/messages?page=3&per_page=50"}}}`, //nolint: lll
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
