package harvest_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestTimesheetService_CreateTimeEntryViaDuration(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/time_entries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{})
		testBody(t, r, "time_entry/create/body_1.json")
		testWriteResponse(t, w, "time_entry/create/response_1.json")
	})

	spentDate := time.Date(2018, 3, 30, 22, 24, 10, 0, time.UTC)

	taskList, _, err := service.Timesheet.CreateTimeEntryViaDuration(
		context.Background(),
		&harvest.TimeEntryCreateViaDuration{
			UserID:    harvest.Int64(1),
			ProjectID: harvest.Int64(2),
			TaskID:    harvest.Int64(3),
			SpentDate: harvest.DateP(harvest.Date{spentDate}),
			Hours:     harvest.Float64(1.2),
			Notes:     harvest.String("Writing tests"),
		},
	)
	if err != nil {
		t.Errorf("CreateTimeEntryViaDuration task returned error: %v", err)
	}

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &harvest.TimeEntry{
		ID:                harvest.Int64(1),
		SpentDate:         nil,
		User:              nil,
		UserAssignment:    nil,
		Client:            nil,
		Project:           nil,
		Task:              nil,
		TaskAssignment:    nil,
		ExternalReference: nil,
		Invoice:           nil,
		Hours:             nil,
		Notes:             nil,
		IsLocked:          nil,
		LockedReason:      nil,
		IsClosed:          nil,
		IsBilled:          nil,
		TimerStartedAt:    nil,
		StartedTime:       nil,
		EndedTime:         nil,
		IsRunning:         nil,
		Billable:          nil,
		Budgeted:          nil,
		BillableRate:      nil,
		CostRate:          nil,
		CreatedAt:         &createdOne,
		UpdatedAt:         &updatedOne,
	}

	if !reflect.DeepEqual(taskList, want) {
		t.Errorf("TimesheetService.CreateTimeEntryViaDuration returned %+v, want %+v", taskList, want)
	}
}

func TestTimesheetService_CreateTimeEntryViaStartEndTime(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	// https://help.getharvest.com/api-v2/timesheets-api/timesheets/time-entries/
	// #create-a-time-entry-via-start-and-end-time

	mux.HandleFunc("/time_entries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{})
		testBody(t, r, "time_entry/create/body_2.json")
		testWriteResponse(t, w, "time_entry/create/response_2.json")
	})

	spentDate := time.Date(2017, 3, 21, 22, 24, 10, 0, time.UTC)

	taskList, _, err := service.Timesheet.CreateTimeEntryViaStartEndTime(
		context.Background(),
		&harvest.TimeEntryCreateViaStartEndTime{
			UserID:      harvest.Int64(1782959),
			ProjectID:   harvest.Int64(14307913),
			TaskID:      harvest.Int64(8083365),
			SpentDate:   harvest.DateP(harvest.Date{spentDate}),
			StartedTime: harvest.TimeP(harvest.Time{time.Date(0, 1, 1, 8, 0o0, 0o0, 0, time.UTC)}),
			EndedTime:   harvest.TimeP(harvest.Time{time.Date(0, 1, 1, 9, 0o0, 0o0, 0, time.UTC)}),
			Notes:       harvest.String("Writing tests"),
		},
	)
	if err != nil {
		t.Errorf("CreateTimeEntry task returned error: %v", err)
	}

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &harvest.TimeEntry{
		ID:        harvest.Int64(1),
		CreatedAt: &createdOne,
		UpdatedAt: &updatedOne,
	}

	assert.Equal(t, taskList, want)
}

func TestTimesheetService_DeleteTimeEntry(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/time_entries/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testFormValues(t, r, values{})
		testBody(t, r, "time_entry/delete/body_1.json")
		testWriteResponse(t, w, "time_entry/delete/response_1.json")
	})

	_, err := service.Timesheet.DeleteTimeEntry(context.Background(), 1)
	if err != nil {
		t.Errorf("DeleteTimeEntry returned error: %v", err)
	}
}

func TestTimesheetService_GetTimeEntry(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/time_entries/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		testBody(t, r, "time_entry/get/body_1.json")
		testWriteResponse(t, w, "time_entry/get/response_1.json")
	})

	timeEntry, _, err := service.Timesheet.Get(context.Background(), 1)
	if err != nil {
		t.Errorf("TimeEntry.GetTimeEntry returned error: %v", err)
	}

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &harvest.TimeEntry{
		ID:        harvest.Int64(1),
		CreatedAt: &createdOne,
		UpdatedAt: &updatedOne,
	}

	assert.Equal(t, timeEntry, want)
}

func TestTimesheetService_ListTimeEntries(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/time_entries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		testBody(t, r, "time_entry/list/body_1.json")
		testWriteResponse(t, w, "time_entry/list/response_1.json")
	})

	taskList, _, err := service.Timesheet.List(context.Background(), &harvest.TimeEntryListOptions{})
	assert.NoError(t, err)

	want := &harvest.TimeEntryList{
		TimeEntries: []*harvest.TimeEntry{
			{
				ID: harvest.Int64(636709355),
				SpentDate: harvest.DateP(harvest.Date{time.Date(
					2017, 3, 2, 0, 0, 0, 0, time.Local)}),
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
				UserAssignment: &harvest.ProjectUserAssignment{
					ID:               harvest.Int64(125068553),
					User:             nil,
					IsActive:         harvest.Bool(true),
					IsProjectManager: harvest.Bool(true),
					HourlyRate:       harvest.Float64(100),
					Budget:           nil,
					CreatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
					UpdatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
				},
				TaskAssignment: &harvest.ProjectTaskAssignment{
					ID:         harvest.Int64(155502709),
					Task:       nil,
					IsActive:   harvest.Bool(true),
					Billable:   harvest.Bool(true),
					HourlyRate: harvest.Float64(100),
					Budget:     nil,
					CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 36, 23, 0, time.UTC)),
					UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 36, 23, 0, time.UTC)),
				},
				ExternalReference: nil,
				Invoice:           nil,
				Hours:             harvest.Float64(2.11),
				RoundedHours:      harvest.Float64(2.25),
				Notes:             harvest.String("Adding CSS styling"),
				IsLocked:          harvest.Bool(true),
				LockedReason:      harvest.String("Item Approved and Locked for this Time Period"),
				IsClosed:          harvest.Bool(true),
				IsBilled:          harvest.Bool(false),
				TimerStartedAt:    nil,
				StartedTime:       harvest.TimeP(harvest.Time{time.Date(0, 1, 1, 15, 0o0, 0o0, 0, time.Local)}),
				EndedTime:         harvest.TimeP(harvest.Time{time.Date(0, 1, 1, 17, 0o0, 0o0, 0, time.Local)}),
				IsRunning:         harvest.Bool(false),
				Billable:          harvest.Bool(true),
				Budgeted:          harvest.Bool(true),
				BillableRate:      harvest.Float64(100),
				CostRate:          harvest.Float64(50),
				CreatedAt:         harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 50, 15, 0, time.UTC)),
				UpdatedAt:         harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 47, 14, 0, time.UTC)),
			}, {
				ID: harvest.Int64(636708723),
				SpentDate: harvest.DateP(harvest.Date{time.Date(
					2017, 3, 1, 0, 0, 0, 0, time.Local)}),
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
				UserAssignment: &harvest.ProjectUserAssignment{
					ID:               harvest.Int64(125068554),
					IsActive:         harvest.Bool(true),
					IsProjectManager: harvest.Bool(true),
					HourlyRate:       harvest.Float64(100),
					CreatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
					UpdatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
				},
				TaskAssignment: &harvest.ProjectTaskAssignment{
					ID:         harvest.Int64(155505014),
					IsActive:   harvest.Bool(true),
					Billable:   harvest.Bool(true),
					HourlyRate: harvest.Float64(100),
					CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
					UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
				},

				Invoice: &harvest.Invoice{
					ID:     harvest.Int64(13150403),
					Number: harvest.String("1001"),
				},
				Hours:        harvest.Float64(1.35),
				RoundedHours: harvest.Float64(1.5),
				Notes:        harvest.String("Importing products"),
				IsLocked:     harvest.Bool(true),
				LockedReason: harvest.String("Item Invoiced and Approved and Locked for this Time Period"),
				IsClosed:     harvest.Bool(true),
				IsBilled:     harvest.Bool(true),
				StartedTime:  harvest.TimeP(harvest.Time{time.Date(0, 1, 1, 13, 0o0, 0o0, 0, time.Local)}),
				EndedTime:    harvest.TimeP(harvest.Time{time.Date(0, 1, 1, 14, 0o0, 0o0, 0, time.Local)}),
				IsRunning:    harvest.Bool(false),
				Billable:     harvest.Bool(true),
				Budgeted:     harvest.Bool(true),
				BillableRate: harvest.Float64(100),
				CostRate:     harvest.Float64(50),
				CreatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 49, 28, 0, time.UTC)),
				UpdatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 47, 14, 0, time.UTC)),
			}, {
				ID: harvest.Int64(636708574),
				SpentDate: harvest.DateP(harvest.Date{time.Date(
					2017, 3, 1, 0, 0, 0, 0, time.Local)}),
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
					ID:   harvest.Int64(8083369),
					Name: harvest.String("Research"),
				},
				UserAssignment: &harvest.ProjectUserAssignment{
					ID:               harvest.Int64(125068554),
					IsActive:         harvest.Bool(true),
					IsProjectManager: harvest.Bool(true),
					HourlyRate:       harvest.Float64(100),
					CreatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
					UpdatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
				},
				TaskAssignment: &harvest.ProjectTaskAssignment{
					ID:         harvest.Int64(155505016),
					IsActive:   harvest.Bool(true),
					Billable:   harvest.Bool(false),
					HourlyRate: harvest.Float64(100),
					CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
					UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 54, 0o6, 0, time.UTC)),
				},

				Invoice:      nil,
				Hours:        harvest.Float64(1),
				RoundedHours: harvest.Float64(1),
				Notes:        harvest.String("Evaluating 3rd party libraries"),
				IsLocked:     harvest.Bool(true),
				LockedReason: harvest.String("Item Approved and Locked for this Time Period"),
				IsClosed:     harvest.Bool(true),
				IsBilled:     harvest.Bool(false),
				StartedTime:  harvest.TimeP(harvest.Time{time.Date(0, 1, 1, 11, 0o0, 0o0, 0, time.Local)}),
				EndedTime:    harvest.TimeP(harvest.Time{time.Date(0, 1, 1, 12, 0o0, 0o0, 0, time.Local)}),
				IsRunning:    harvest.Bool(false),

				Billable:     harvest.Bool(false),
				Budgeted:     harvest.Bool(true),
				BillableRate: nil,
				CostRate:     harvest.Float64(50),
				CreatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 49, 17, 0, time.UTC)),
				UpdatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 47, 14, 0, time.UTC)),
			}, {
				ID: harvest.Int64(636707831),
				SpentDate: harvest.DateP(harvest.Date{time.Date(
					2017, 3, 1, 0, 0, 0, 0, time.Local)}),
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
					ID:   harvest.Int64(8083368),
					Name: harvest.String("Project Management"),
				},
				UserAssignment: &harvest.ProjectUserAssignment{
					ID:               harvest.Int64(125068554),
					IsActive:         harvest.Bool(true),
					IsProjectManager: harvest.Bool(true),
					HourlyRate:       harvest.Float64(100),
					CreatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
					UpdatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
				},
				TaskAssignment: &harvest.ProjectTaskAssignment{
					ID:         harvest.Int64(155505015),
					IsActive:   harvest.Bool(true),
					Billable:   harvest.Bool(true),
					HourlyRate: harvest.Float64(100),
					CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
					UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
				},

				Invoice: &harvest.Invoice{
					ID:     harvest.Int64(13150403),
					Number: harvest.String("1001"),
				},
				Hours:        harvest.Float64(2),
				RoundedHours: harvest.Float64(2),
				Notes:        harvest.String("Planning meetings"),
				IsLocked:     harvest.Bool(true),
				LockedReason: harvest.String("Item Invoiced and Approved and Locked for this Time Period"),
				IsClosed:     harvest.Bool(true),
				IsBilled:     harvest.Bool(true),
				StartedTime:  harvest.TimeP(harvest.Time{time.Date(0, 1, 1, 9, 0o0, 0o0, 0, time.Local)}),
				EndedTime:    harvest.TimeP(harvest.Time{time.Date(0, 1, 1, 11, 0o0, 0o0, 0, time.Local)}),
				IsRunning:    harvest.Bool(false),

				Billable:     harvest.Bool(true),
				Budgeted:     harvest.Bool(true),
				BillableRate: harvest.Float64(100),
				CostRate:     harvest.Float64(50),
				CreatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 48, 24, 0, time.UTC)),
				UpdatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 47, 14, 0, time.UTC)),
			},
		},

		Pagination: harvest.Pagination{
			PerPage:      harvest.Int(100),
			TotalPages:   harvest.Int(1),
			TotalEntries: harvest.Int(4),
			NextPage:     nil,
			PreviousPage: nil,
			Page:         harvest.Int(1),
			Links: &harvest.PageLinks{
				First:    harvest.String("https://api.harvestapp.com/v2/time_entries?page=1&per_page=100"),
				Next:     nil,
				Previous: nil,
				Last:     harvest.String("https://api.harvestapp.com/v2/time_entries?page=1&per_page=100"),
			},
		},
	}

	assert.Equal(t, want, taskList)
}

func TestTimesheetService_RestartTimeEntry(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/time_entries/662202797/restart", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		testBody(t, r, "time_entry/restart/body_1.json")
		testWriteResponse(t, w, "time_entry/restart/response_1.json")
	})

	timeEntry, _, err := service.Timesheet.RestartTimeEntry(context.Background(), 662202797)
	assert.NoError(t, err)

	spentDate := time.Date(2017, 3, 21, 0, 0, 0, 0, time.Local)

	createdOne := time.Date(
		2017, 8, 22, 17, 40, 24, 0, time.UTC)
	updatedOne := time.Date(
		2017, 8, 22, 17, 40, 24, 0, time.UTC)

	assignmentTime := time.Date(
		2017, 8, 22, 17, 36, 54, 0, time.UTC)

	timerStartedAt := time.Date(
		2017, 8, 22, 17, 40, 24, 0, time.UTC)
	startedTime := harvest.Time{time.Date(0, 1, 1, 11, 40, 0o0, 0, time.Local)}

	want := &harvest.TimeEntry{
		ID:        harvest.Int64(662202797),
		SpentDate: harvest.DateP(harvest.Date{spentDate}),
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
		UserAssignment: &harvest.ProjectUserAssignment{
			ID:               harvest.Int64(130403296),
			IsProjectManager: harvest.Bool(true),
			IsActive:         harvest.Bool(true),
			Budget:           nil,
			CreatedAt:        harvest.TimeTimeP(assignmentTime),
			UpdatedAt:        harvest.TimeTimeP(assignmentTime),
			HourlyRate:       harvest.Float64(100),
		},
		TaskAssignment: &harvest.ProjectTaskAssignment{
			ID:         harvest.Int64(160726645),
			Billable:   harvest.Bool(true),
			IsActive:   harvest.Bool(true),
			CreatedAt:  harvest.TimeTimeP(assignmentTime),
			UpdatedAt:  harvest.TimeTimeP(assignmentTime),
			HourlyRate: harvest.Float64(100),
			Budget:     nil,
		},
		Hours:             harvest.Float64(0),
		RoundedHours:      harvest.Float64(0),
		Notes:             nil,
		CreatedAt:         &createdOne,
		UpdatedAt:         &updatedOne,
		IsLocked:          harvest.Bool(false),
		LockedReason:      nil,
		IsClosed:          harvest.Bool(false),
		IsBilled:          harvest.Bool(false),
		TimerStartedAt:    harvest.TimeTimeP(timerStartedAt),
		StartedTime:       harvest.TimeP(startedTime),
		EndedTime:         nil,
		IsRunning:         harvest.Bool(true),
		Invoice:           nil,
		ExternalReference: nil,
		Billable:          harvest.Bool(true),
		Budgeted:          harvest.Bool(false),
		BillableRate:      harvest.Float64(100),
		CostRate:          harvest.Float64(75),
	}

	assert.Equal(t, want, timeEntry)
}

func TestTimesheetService_StopTimeEntry(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/time_entries/662202797/stop", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		testBody(t, r, "time_entry/stop/body_1.json")
		testWriteResponse(t, w, "time_entry/stop/response_1.json")
	})

	timeEntry, _, err := service.Timesheet.StopTimeEntry(context.Background(), 662202797)
	assert.NoError(t, err)

	spentDate := time.Date(2017, 3, 21, 0, 0, 0, 0, time.Local)

	createdOne := time.Date(2017, 8, 22, 17, 37, 13, 0, time.UTC)
	updatedOne := time.Date(2017, 8, 22, 17, 38, 31, 0, time.UTC)

	assignmentTime := time.Date(2017, 8, 22, 17, 36, 54, 0, time.UTC)

	startedTime := harvest.Time{time.Date(0, 1, 1, 11, 37, 0o0, 0, time.Local)}
	endedTime := harvest.Time{time.Date(0, 1, 1, 11, 38, 0o0, 0, time.Local)}

	want := &harvest.TimeEntry{
		ID:        harvest.Int64(662202797),
		SpentDate: harvest.DateP(harvest.Date{spentDate}),
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
		UserAssignment: &harvest.ProjectUserAssignment{
			ID:               harvest.Int64(130403296),
			IsProjectManager: harvest.Bool(true),
			IsActive:         harvest.Bool(true),
			Budget:           nil,
			CreatedAt:        harvest.TimeTimeP(assignmentTime),
			UpdatedAt:        harvest.TimeTimeP(assignmentTime),
			HourlyRate:       harvest.Float64(100),
		},
		TaskAssignment: &harvest.ProjectTaskAssignment{
			ID:         harvest.Int64(160726645),
			Billable:   harvest.Bool(true),
			IsActive:   harvest.Bool(true),
			CreatedAt:  harvest.TimeTimeP(assignmentTime),
			UpdatedAt:  harvest.TimeTimeP(assignmentTime),
			HourlyRate: harvest.Float64(100),
			Budget:     nil,
		},
		Hours:             harvest.Float64(0.02),
		RoundedHours:      harvest.Float64(0.25),
		Notes:             nil,
		CreatedAt:         &createdOne,
		UpdatedAt:         &updatedOne,
		IsLocked:          harvest.Bool(false),
		LockedReason:      nil,
		IsClosed:          harvest.Bool(false),
		IsBilled:          harvest.Bool(false),
		TimerStartedAt:    nil,
		StartedTime:       harvest.TimeP(startedTime),
		EndedTime:         harvest.TimeP(endedTime),
		IsRunning:         harvest.Bool(false),
		Invoice:           nil,
		ExternalReference: nil,
		Billable:          harvest.Bool(true),
		Budgeted:          harvest.Bool(false),
		BillableRate:      harvest.Float64(100),
		CostRate:          harvest.Float64(75),
	}

	assert.Equal(t, want, timeEntry)
}

func TestTimesheetService_UpdateTimeEntry(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/time_entries/636718192", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		testBody(t, r, "time_entry/update/body_1.json")
		testWriteResponse(t, w, "time_entry/update/response_1.json")
	})

	spentDate := time.Date(2017, 3, 21, 0, 0, 0, 0, time.Local)
	startedTime := harvest.Time{time.Date(0, 1, 1, 11, 40, 0o0, 0, time.UTC)}
	endedTime := harvest.Time{time.Date(0, 1, 1, 12, 45, 10, 0, time.UTC)}

	timeEntry, _, err := service.Timesheet.UpdateTimeEntry(context.Background(), 636718192, &harvest.TimeEntryUpdate{
		ProjectID:   harvest.Int64(1234),
		TaskID:      harvest.Int64(2345),
		SpentDate:   harvest.DateP(harvest.Date{spentDate}),
		StartedTime: &startedTime,
		EndedTime:   &endedTime,
		Hours:       harvest.Float64(1),
		Notes:       harvest.String("new notes"),
	})
	assert.NoError(t, err)

	createdTimeEntry := time.Date(
		2017, 6, 27, 16, 0o1, 23, 0, time.UTC)
	updatedTimeEntry := time.Date(
		2017, 6, 27, 16, 0o2, 40, 0, time.UTC)
	createdUserAssignment := time.Date(
		2017, 6, 26, 22, 32, 52, 0, time.UTC)
	updatedUserAssignment := time.Date(
		2017, 6, 26, 22, 32, 52, 0, time.UTC)
	createdTaskAssignment := time.Date(
		2017, 6, 26, 21, 36, 23, 0, time.UTC)
	updatedTaskAssignment := time.Date(
		2017, 6, 26, 21, 36, 23, 0, time.UTC)

	want := &harvest.TimeEntry{
		ID:        harvest.Int64(636718192),
		CreatedAt: &createdTimeEntry,
		UpdatedAt: &updatedTimeEntry,
		SpentDate: harvest.DateP(harvest.Date{spentDate}),
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
		UserAssignment: &harvest.ProjectUserAssignment{
			ID:               harvest.Int64(125068553),
			IsProjectManager: harvest.Bool(true),
			IsActive:         harvest.Bool(true),
			CreatedAt:        &createdUserAssignment,
			UpdatedAt:        &updatedUserAssignment,
			HourlyRate:       harvest.Float64(100),
		},
		TaskAssignment: &harvest.ProjectTaskAssignment{
			ID:         harvest.Int64(155502709),
			Billable:   harvest.Bool(true),
			IsActive:   harvest.Bool(true),
			CreatedAt:  &createdTaskAssignment,
			UpdatedAt:  &updatedTaskAssignment,
			HourlyRate: harvest.Float64(100),
		},
		Hours:        harvest.Float64(1),
		RoundedHours: harvest.Float64(1),
		Notes:        harvest.String("Updated notes"),
		IsLocked:     harvest.Bool(false),
		IsClosed:     harvest.Bool(false),
		IsBilled:     harvest.Bool(false),
		IsRunning:    harvest.Bool(false),
		Billable:     harvest.Bool(true),
		Budgeted:     harvest.Bool(true),
		BillableRate: harvest.Float64(100),
		CostRate:     harvest.Float64(50),
	}

	assert.Equal(t, want, timeEntry)
}
