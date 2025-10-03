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
		{
			name:        "Unauthorized Access",
			timeEntryID: 1,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries/1", func(w http.ResponseWriter, _ *http.Request) {
					http.Error(w, `{"message":"Unauthorized"}`, http.StatusUnauthorized)
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

func TestTimesheetService_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(mux *http.ServeMux)
		want      int // Just check the count of entries
		wantErr   bool
	}{
		{
			name: "Valid Time Entry List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "time_entry/list/body_1.json")
					testWriteResponse(t, w, "time_entry/list/response_1.json")
				})
			},
			want:    4, // The response file has 4 entries
			wantErr: false,
		},
		{
			name: "Error Fetching Time Entry List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					http.Error(w, `{"message":"Internal Server Error"}`, http.StatusInternalServerError)
				})
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Empty Time Entry List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{
						"time_entries": [],
						"per_page": 100,
						"total_pages": 0,
						"total_entries": 0,
						"page": 1,
						"links": {
							"first": "https://api.harvestapp.com/v2/time_entries?page=1&per_page=100",
							"last": "https://api.harvestapp.com/v2/time_entries?page=1&per_page=100"
						}
					}`))
				})
			},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			tt.setupMock(mux)

			got, _, err := service.Timesheet.List(context.Background(), &harvest.TimeEntryListOptions{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, len(got.TimeEntries))
			}
		})
	}
}

func TestTimesheetService_CreateTimeEntryViaStartEndTime(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2017, 6, 27, 15, 49, 28, 0, time.UTC)
	updatedAt := time.Date(2017, 6, 27, 16, 47, 14, 0, time.UTC)
	spentDate := harvest.Date{Time: time.Date(2017, 3, 21, 0, 0, 0, 0, time.UTC)}

	tests := []struct {
		name      string
		input     *harvest.TimeEntryCreateViaStartEndTime
		setupMock func(mux *http.ServeMux)
		want      *harvest.TimeEntry
		wantErr   bool
	}{
		{
			name: "Valid Time Entry Creation via Start/End Time",
			input: &harvest.TimeEntryCreateViaStartEndTime{
				ProjectID:   harvest.Int64(14308069),
				TaskID:      harvest.Int64(8083366),
				SpentDate:   &spentDate,
				StartedTime: &harvest.Time{Time: time.Date(0, time.January, 1, 3, 0, 0, 0, time.Local)},
				EndedTime:   &harvest.Time{Time: time.Date(0, time.January, 1, 5, 0, 0, 0, time.Local)},
				Notes:       harvest.String("Importing products"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testWriteResponse(t, w, "time_entry/create/response_via_start_end.json")
				})
			},
			want: &harvest.TimeEntry{
				ID:           harvest.Int64(636708723),
				SpentDate:    &spentDate,
				Hours:        harvest.Float64(2.0),
				RoundedHours: harvest.Float64(2.0),
				Notes:        harvest.String("Importing products"),
				IsLocked:     harvest.Bool(false),
				LockedReason: nil,
				IsClosed:     harvest.Bool(false),
				IsBilled:     harvest.Bool(false),
				StartedTime:  &harvest.Time{Time: time.Date(0, time.January, 1, 3, 0, 0, 0, time.Local)},
				EndedTime:    &harvest.Time{Time: time.Date(0, time.January, 1, 5, 0, 0, 0, time.Local)},
				IsRunning:    harvest.Bool(false),
				Billable:     harvest.Bool(true),
				Budgeted:     harvest.Bool(true),
				BillableRate: harvest.Float64(100.0),
				CostRate:     harvest.Float64(50.0),
				CreatedAt:    &createdAt,
				UpdatedAt:    &updatedAt,
				User: &harvest.User{
					ID:   harvest.Int64(1782959),
					Name: harvest.String("Kim Allen"),
				},
				Client: &harvest.Client{
					ID:   harvest.Int64(5735776),
					Name: harvest.String("123 Industries"),
				},
				Project: &harvest.Project{
					ID:   harvest.Int64(14308069),
					Name: harvest.String("Online Store - Phase 1"),
				},
				Task: &harvest.Task{
					ID:   harvest.Int64(8083366),
					Name: harvest.String("Programming"),
				},
			},
			wantErr: false,
		},
		{
			name: "Error Creating Time Entry",
			input: &harvest.TimeEntryCreateViaStartEndTime{
				ProjectID: harvest.Int64(14307913),
				TaskID:    harvest.Int64(8083365),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					http.Error(w, `{"message":"Invalid data"}`, http.StatusBadRequest)
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

			got, _, err := service.Timesheet.CreateTimeEntryViaStartEndTime(context.Background(), tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.Notes, got.Notes)
			}
		})
	}
}

func TestTimesheetService_UpdateTimeEntry(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2017, 6, 27, 16, 1, 23, 0, time.UTC)
	updatedAt := time.Date(2017, 6, 27, 16, 2, 40, 0, time.UTC)
	spentDate := harvest.Date{Time: time.Date(2017, 3, 21, 0, 0, 0, 0, time.UTC)}

	tests := []struct {
		name        string
		timeEntryID int64
		input       *harvest.TimeEntryUpdate
		setupMock   func(mux *http.ServeMux)
		want        *harvest.TimeEntry
		wantErr     bool
	}{
		{
			name:        "Valid Time Entry Update",
			timeEntryID: 636718192,
			input: &harvest.TimeEntryUpdate{
				ProjectID: harvest.Int64(14307913),
				TaskID:    harvest.Int64(8083365),
				SpentDate: &spentDate,
				Hours:     harvest.Float64(1.0),
				Notes:     harvest.String("Updated notes"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries/636718192", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					testWriteResponse(t, w, "time_entry/update/response_1.json")
				})
			},
			want: &harvest.TimeEntry{
				ID:           harvest.Int64(636718192),
				SpentDate:    &spentDate,
				Hours:        harvest.Float64(1.0),
				RoundedHours: harvest.Float64(1.0),
				Notes:        harvest.String("Updated notes"),
				IsLocked:     harvest.Bool(false),
				LockedReason: nil,
				IsClosed:     harvest.Bool(false),
				IsBilled:     harvest.Bool(false),
				IsRunning:    harvest.Bool(false),
				Billable:     harvest.Bool(true),
				Budgeted:     harvest.Bool(true),
				BillableRate: harvest.Float64(100.0),
				CostRate:     harvest.Float64(50.0),
				CreatedAt:    &createdAt,
				UpdatedAt:    &updatedAt,
				User: &harvest.User{
					ID:   harvest.Int64(1782959),
					Name: harvest.String("Kim Allen"),
				},
				Client: &harvest.Client{
					ID:   harvest.Int64(5735774),
					Name: harvest.String("ABC Corp"),
				},
				Project: &harvest.Project{
					ID:   harvest.Int64(14307913),
					Name: harvest.String("Marketing Website"),
				},
				Task: &harvest.Task{
					ID:   harvest.Int64(8083365),
					Name: harvest.String("Graphic Design"),
				},
			},
			wantErr: false,
		},
		{
			name:        "Error Updating Time Entry",
			timeEntryID: 999,
			input: &harvest.TimeEntryUpdate{
				ProjectID: harvest.Int64(14307913),
				TaskID:    harvest.Int64(8083365),
				SpentDate: &spentDate,
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					http.Error(w, `{"message":"Time entry not found"}`, http.StatusNotFound)
				})
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:        "Forbidden Update",
			timeEntryID: 636709355,
			input: &harvest.TimeEntryUpdate{
				ProjectID: harvest.Int64(14307913),
				TaskID:    harvest.Int64(8083365),
				SpentDate: &spentDate,
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries/636709355", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					http.Error(w, `{"message":"Forbidden"}`, http.StatusForbidden)
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

			got, _, err := service.Timesheet.UpdateTimeEntry(context.Background(), tt.timeEntryID, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.Notes, got.Notes)
			}
		})
	}
}

func TestTimesheetService_DeleteTimeEntry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		timeEntryID int64
		setupMock   func(mux *http.ServeMux)
		wantErr     bool
	}{
		{
			name:        "Valid Time Entry Deletion",
			timeEntryID: 636709355,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries/636709355", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					w.WriteHeader(http.StatusNoContent)
				})
			},
			wantErr: false,
		},
		{
			name:        "Time Entry Not Found",
			timeEntryID: 999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					http.Error(w, `{"message":"Time entry not found"}`, http.StatusNotFound)
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

			_, err := service.Timesheet.DeleteTimeEntry(context.Background(), tt.timeEntryID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTimesheetService_RestartTimeEntry(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2017, 8, 22, 17, 40, 24, 0, time.UTC)
	updatedAt := time.Date(2017, 8, 22, 17, 40, 24, 0, time.UTC)
	timerStartedAt := time.Date(2017, 8, 22, 17, 40, 24, 0, time.UTC)
	spentDate := harvest.Date{Time: time.Date(2017, 3, 21, 0, 0, 0, 0, time.UTC)}

	tests := []struct {
		name        string
		timeEntryID int64
		setupMock   func(mux *http.ServeMux)
		want        *harvest.TimeEntry
		wantErr     bool
	}{
		{
			name:        "Valid Time Entry Restart",
			timeEntryID: 662202797,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries/662202797/restart", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					testBody(t, r, "time_entry/restart/body_1.json")
					testWriteResponse(t, w, "time_entry/restart/response_1.json")
				})
			},
			want: &harvest.TimeEntry{
				ID:             harvest.Int64(662202797),
				SpentDate:      &spentDate,
				Hours:          harvest.Float64(0.0),
				RoundedHours:   harvest.Float64(0.0),
				Notes:          nil,
				IsLocked:       harvest.Bool(false),
				LockedReason:   nil,
				IsClosed:       harvest.Bool(false),
				IsBilled:       harvest.Bool(false),
				TimerStartedAt: &timerStartedAt,
				IsRunning:      harvest.Bool(true),
				Billable:       harvest.Bool(true),
				Budgeted:       harvest.Bool(false),
				BillableRate:   harvest.Float64(100.0),
				CostRate:       harvest.Float64(75.0),
				CreatedAt:      &createdAt,
				UpdatedAt:      &updatedAt,
				User: &harvest.User{
					ID:   harvest.Int64(1795925),
					Name: harvest.String("Jane Smith"),
				},
				Client: &harvest.Client{
					ID:   harvest.Int64(5735776),
					Name: harvest.String("123 Industries"),
				},
				Project: &harvest.Project{
					ID:   harvest.Int64(14808188),
					Name: harvest.String("Task Force"),
				},
				Task: &harvest.Task{
					ID:   harvest.Int64(8083366),
					Name: harvest.String("Programming"),
				},
			},
			wantErr: false,
		},
		{
			name:        "Time Entry Not Found",
			timeEntryID: 999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries/999/restart", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					http.Error(w, `{"message":"Time entry not found"}`, http.StatusNotFound)
				})
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:        "Already Running",
			timeEntryID: 636709355,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries/636709355/restart", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					http.Error(w, `{"message":"Time entry is already running"}`, http.StatusUnprocessableEntity)
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

			got, _, err := service.Timesheet.RestartTimeEntry(context.Background(), tt.timeEntryID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.IsRunning, got.IsRunning)
				assert.NotNil(t, got.TimerStartedAt)
			}
		})
	}
}

func TestTimesheetService_StopTimeEntry(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2017, 8, 22, 17, 37, 13, 0, time.UTC)
	updatedAt := time.Date(2017, 8, 22, 17, 38, 31, 0, time.UTC)
	spentDate := harvest.Date{Time: time.Date(2017, 3, 21, 0, 0, 0, 0, time.UTC)}

	tests := []struct {
		name        string
		timeEntryID int64
		setupMock   func(mux *http.ServeMux)
		want        *harvest.TimeEntry
		wantErr     bool
	}{
		{
			name:        "Valid Time Entry Stop",
			timeEntryID: 662202797,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries/662202797/stop", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					testBody(t, r, "time_entry/stop/body_1.json")
					testWriteResponse(t, w, "time_entry/stop/response_1.json")
				})
			},
			want: &harvest.TimeEntry{
				ID:             harvest.Int64(662202797),
				SpentDate:      &spentDate,
				Hours:          harvest.Float64(0.02),
				RoundedHours:   harvest.Float64(0.25),
				Notes:          nil,
				IsLocked:       harvest.Bool(false),
				LockedReason:   nil,
				IsClosed:       harvest.Bool(false),
				IsBilled:       harvest.Bool(false),
				TimerStartedAt: nil,
				IsRunning:      harvest.Bool(false),
				Billable:       harvest.Bool(true),
				Budgeted:       harvest.Bool(false),
				BillableRate:   harvest.Float64(100.0),
				CostRate:       harvest.Float64(75.0),
				CreatedAt:      &createdAt,
				UpdatedAt:      &updatedAt,
				User: &harvest.User{
					ID:   harvest.Int64(1795925),
					Name: harvest.String("Jane Smith"),
				},
				Client: &harvest.Client{
					ID:   harvest.Int64(5735776),
					Name: harvest.String("123 Industries"),
				},
				Project: &harvest.Project{
					ID:   harvest.Int64(14808188),
					Name: harvest.String("Task Force"),
				},
				Task: &harvest.Task{
					ID:   harvest.Int64(8083366),
					Name: harvest.String("Programming"),
				},
			},
			wantErr: false,
		},
		{
			name:        "Time Entry Not Found",
			timeEntryID: 999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/time_entries/999/stop", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
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

			got, _, err := service.Timesheet.StopTimeEntry(context.Background(), tt.timeEntryID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.IsRunning, got.IsRunning)
				assert.Nil(t, got.TimerStartedAt)
			}
		})
	}
}

func TestTimeEntry_String(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2017, 6, 27, 15, 49, 28, 0, time.UTC)
	updatedAt := time.Date(2017, 6, 27, 16, 47, 14, 0, time.UTC)
	spentDate := harvest.Date{Time: time.Date(2017, 3, 21, 0, 0, 0, 0, time.UTC)}

	tests := []struct {
		name string
		in   harvest.TimeEntry
		want string
	}{
		{
			name: "Time Entry with all fields",
			in: harvest.TimeEntry{
				ID:           harvest.Int64(636709355),
				SpentDate:    &spentDate,
				Hours:        harvest.Float64(2.11),
				RoundedHours: harvest.Float64(2.25),
				Notes:        harvest.String("Adding CSS styling"),
				IsLocked:     harvest.Bool(false),
				IsClosed:     harvest.Bool(false),
				IsBilled:     harvest.Bool(false),
				IsRunning:    harvest.Bool(false),
				Billable:     harvest.Bool(true),
				Budgeted:     harvest.Bool(true),
				BillableRate: harvest.Float64(100.0),
				CostRate:     harvest.Float64(50.0),
				CreatedAt:    &createdAt,
				UpdatedAt:    &updatedAt,
				User: &harvest.User{
					ID:   harvest.Int64(1782959),
					Name: harvest.String("Kim Allen"),
				},
				Client: &harvest.Client{
					ID:   harvest.Int64(5735774),
					Name: harvest.String("ABC Corp"),
				},
				Project: &harvest.Project{
					ID:   harvest.Int64(14307913),
					Name: harvest.String("Marketing Website"),
				},
				Task: &harvest.Task{
					ID:   harvest.Int64(8083365),
					Name: harvest.String("Graphic Design"),
				},
			},
			want: `harvest.TimeEntry{ID:636709355, SpentDate:harvest.Date{{2017-03-21 00:00:00 +0000 UTC}}, User:harvest.User{ID:1782959, Name:"Kim Allen"}, Client:harvest.Client{ID:5735774, Name:"ABC Corp"}, Project:harvest.Project{ID:14307913, Name:"Marketing Website"}, Task:harvest.Task{ID:8083365, Name:"Graphic Design"}, Hours:2.11, RoundedHours:2.25, Notes:"Adding CSS styling", IsLocked:false, IsClosed:false, IsBilled:false, IsRunning:false, Billable:true, Budgeted:true, BillableRate:100, CostRate:50, CreatedAt:time.Time{2017-06-27 15:49:28 +0000 UTC}, UpdatedAt:time.Time{2017-06-27 16:47:14 +0000 UTC}}`, //nolint: lll
		},
		{
			name: "Time Entry with minimal fields",
			in: harvest.TimeEntry{
				ID:        harvest.Int64(636709355),
				SpentDate: &spentDate,
				Hours:     harvest.Float64(2.11),
			},
			want: `harvest.TimeEntry{ID:636709355, SpentDate:harvest.Date{{2017-03-21 00:00:00 +0000 UTC}}, Hours:2.11}`,
		},
		{
			name: "Time Entry with nested objects",
			in: harvest.TimeEntry{
				ID:        harvest.Int64(636709355),
				SpentDate: &spentDate,
				User: &harvest.User{
					ID:   harvest.Int64(1782959),
					Name: harvest.String("Kim Allen"),
				},
				Client: &harvest.Client{
					ID:   harvest.Int64(5735774),
					Name: harvest.String("ABC Corp"),
				},
			},
			want: `harvest.TimeEntry{ID:636709355, SpentDate:harvest.Date{{2017-03-21 00:00:00 +0000 UTC}}, User:harvest.User{ID:1782959, Name:"Kim Allen"}, Client:harvest.Client{ID:5735774, Name:"ABC Corp"}}`, //nolint: lll
		},
		{
			name: "Time Entry with timer running",
			in: harvest.TimeEntry{
				ID:             harvest.Int64(636709355),
				SpentDate:      &spentDate,
				Hours:          harvest.Float64(0.0),
				IsRunning:      harvest.Bool(true),
				TimerStartedAt: &createdAt,
			},
			want: `harvest.TimeEntry{ID:636709355, SpentDate:harvest.Date{{2017-03-21 00:00:00 +0000 UTC}}, Hours:0, TimerStartedAt:time.Time{2017-06-27 15:49:28 +0000 UTC}, IsRunning:true}`, //nolint: lll
		},
		{
			name: "Empty Time Entry",
			in:   harvest.TimeEntry{},
			want: `harvest.TimeEntry{}`,
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

func TestTimeEntryList_String(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2017, 6, 27, 15, 49, 28, 0, time.UTC)
	updatedAt := time.Date(2017, 6, 27, 16, 47, 14, 0, time.UTC)
	spentDate := harvest.Date{Time: time.Date(2017, 3, 21, 0, 0, 0, 0, time.UTC)}

	tests := []struct {
		name string
		in   harvest.TimeEntryList
		want string
	}{
		{
			name: "Time Entry List with all fields",
			in: harvest.TimeEntryList{
				TimeEntries: []*harvest.TimeEntry{
					{
						ID:           harvest.Int64(636709355),
						SpentDate:    &spentDate,
						Hours:        harvest.Float64(2.11),
						RoundedHours: harvest.Float64(2.25),
						Notes:        harvest.String("Adding CSS styling"),
						IsLocked:     harvest.Bool(false),
						IsClosed:     harvest.Bool(false),
						IsBilled:     harvest.Bool(false),
						IsRunning:    harvest.Bool(false),
						Billable:     harvest.Bool(true),
						Budgeted:     harvest.Bool(true),
						BillableRate: harvest.Float64(100.0),
						CostRate:     harvest.Float64(50.0),
						CreatedAt:    &createdAt,
						UpdatedAt:    &updatedAt,
						User: &harvest.User{
							ID:   harvest.Int64(1782959),
							Name: harvest.String("Kim Allen"),
						},
						Client: &harvest.Client{
							ID:   harvest.Int64(5735774),
							Name: harvest.String("ABC Corp"),
						},
						Project: &harvest.Project{
							ID:   harvest.Int64(14307913),
							Name: harvest.String("Marketing Website"),
						},
						Task: &harvest.Task{
							ID:   harvest.Int64(8083365),
							Name: harvest.String("Graphic Design"),
						},
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.TimeEntryList{TimeEntries:[harvest.TimeEntry{ID:636709355, SpentDate:harvest.Date{{2017-03-21 00:00:00 +0000 UTC}}, User:harvest.User{ID:1782959, Name:"Kim Allen"}, Client:harvest.Client{ID:5735774, Name:"ABC Corp"}, Project:harvest.Project{ID:14307913, Name:"Marketing Website"}, Task:harvest.Task{ID:8083365, Name:"Graphic Design"}, Hours:2.11, RoundedHours:2.25, Notes:"Adding CSS styling", IsLocked:false, IsClosed:false, IsBilled:false, IsRunning:false, Billable:true, Budgeted:true, BillableRate:100, CostRate:50, CreatedAt:time.Time{2017-06-27 15:49:28 +0000 UTC}, UpdatedAt:time.Time{2017-06-27 16:47:14 +0000 UTC}}], Pagination:harvest.Pagination{PerPage:100, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
		},
		{
			name: "Time Entry List with minimal fields",
			in: harvest.TimeEntryList{
				TimeEntries: []*harvest.TimeEntry{
					{
						ID:        harvest.Int64(636709355),
						SpentDate: &spentDate,
						Hours:     harvest.Float64(2.11),
					},
				},
				Pagination: harvest.Pagination{
					Page: harvest.Int(1),
				},
			},
			want: `harvest.TimeEntryList{TimeEntries:[harvest.TimeEntry{ID:636709355, SpentDate:harvest.Date{{2017-03-21 00:00:00 +0000 UTC}}, Hours:2.11}], Pagination:harvest.Pagination{Page:1}}`, //nolint: lll
		},
		{
			name: "Time Entry List with multiple entries",
			in: harvest.TimeEntryList{
				TimeEntries: []*harvest.TimeEntry{
					{
						ID:        harvest.Int64(636709355),
						SpentDate: &spentDate,
						Hours:     harvest.Float64(2.11),
					},
					{
						ID:        harvest.Int64(636709356),
						SpentDate: &spentDate,
						Hours:     harvest.Float64(3.5),
					},
				},
				Pagination: harvest.Pagination{
					Page:         harvest.Int(1),
					TotalEntries: harvest.Int(2),
				},
			},
			want: `harvest.TimeEntryList{TimeEntries:[harvest.TimeEntry{ID:636709355, SpentDate:harvest.Date{{2017-03-21 00:00:00 +0000 UTC}}, Hours:2.11} harvest.TimeEntry{ID:636709356, SpentDate:harvest.Date{{2017-03-21 00:00:00 +0000 UTC}}, Hours:3.5}], Pagination:harvest.Pagination{TotalEntries:2, Page:1}}`, //nolint: lll
		},
		{
			name: "Empty Time Entry List",
			in:   harvest.TimeEntryList{},
			want: `harvest.TimeEntryList{Pagination:harvest.Pagination{}}`,
		},
		{
			name: "Time Entry List with empty entries slice",
			in: harvest.TimeEntryList{
				TimeEntries: []*harvest.TimeEntry{},
				Pagination: harvest.Pagination{
					Page:         harvest.Int(1),
					TotalEntries: harvest.Int(0),
				},
			},
			want: `harvest.TimeEntryList{TimeEntries:[], Pagination:harvest.Pagination{TotalEntries:0, Page:1}}`,
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
