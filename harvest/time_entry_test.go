package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestTimesheetService_CreateTimeEntryViaDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		request   *harvest.TimeEntryCreateViaDuration
		setupMock func(mux *http.ServeMux)
		want      *harvest.TimeEntry
		wantErr   bool
	}{
		{
			name: "Valid Time Entry Creation",
			request: &harvest.TimeEntryCreateViaDuration{
				UserID:    harvest.Int64(1),
				ProjectID: harvest.Int64(2),
				TaskID:    harvest.Int64(3),
				SpentDate: harvest.DateP(harvest.Date{Time: time.Date(2018, 3, 30, 0, 0, 0, 0, time.UTC)}),
				Hours:     harvest.Float64(1.2),
				Notes:     harvest.String("Writing tests"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testBody(t, r, "time_entry/create/body_1.json")
					testWriteResponse(t, w, "time_entry/create/response_1.json")
				})
			},
			want: &harvest.TimeEntry{
				ID:        harvest.Int64(1),
				CreatedAt: harvest.TimeTimeP(time.Date(2018, 1, 31, 20, 34, 30, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2018, 5, 31, 21, 34, 30, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name: "Invalid Time Entry Creation - Missing ProjectID",
			request: &harvest.TimeEntryCreateViaDuration{
				UserID: harvest.Int64(1),
				TaskID: harvest.Int64(3),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries", func(w http.ResponseWriter, _ *http.Request) {
					http.Error(w, `{"message":"ProjectID is required"}`, http.StatusBadRequest)
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

			got, _, err := service.Timesheet.CreateTimeEntryViaDuration(context.Background(), tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestTimesheetService_GetTimeEntry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		timeEntryID int64
		setupMock   func(mux *http.ServeMux)
		want        *harvest.TimeEntry
		wantErr     bool
	}{
		{
			name:        "Valid Time Entry Retrieval",
			timeEntryID: 1,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries/1", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testWriteResponse(t, w, "time_entry/get/response_1.json")
				})
			},
			want: &harvest.TimeEntry{
				ID:        harvest.Int64(1),
				CreatedAt: harvest.TimeTimeP(time.Date(2018, 1, 31, 20, 34, 30, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2018, 5, 31, 21, 34, 30, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name:        "Time Entry Not Found",
			timeEntryID: 999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries/999", func(w http.ResponseWriter, _ *http.Request) {
					http.Error(w, `{"message":"Time entry not found"}`, http.StatusNotFound)
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

			got, _, err := service.Timesheet.Get(context.Background(), tt.timeEntryID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
