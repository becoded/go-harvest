package harvest

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimesheetService_CreateTimeEntryViaDuration(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/time_entries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{})
		testBody(t, r, `{"user_id":1,"project_id":2,"task_id":3,"spent_date":"2018-03-30T22:24:10Z","hours":1.2,"notes":"Writing tests"}`+"\n")
		fmt.Fprint(w, `{"id":1,"name":"TimeEntry new","billable_by_default":true,"default_hourly_rate":123,"is_default":true,"is_active":true, "created_at":"2018-01-31T20:34:30Z","updated_at":"2018-05-31T21:34:30Z"}`)
	})

	spentDate := time.Date(2018, 3, 30, 22, 24, 10, 0, time.UTC)

	taskList, _, err := service.Timesheet.CreateTimeEntryViaDuration(context.Background(), &TimeEntryCreateViaDuration{
		UserId:    Int64(1),
		ProjectId: Int64(2),
		TaskId:    Int64(3),
		SpentDate: DateP(Date{spentDate}),
		Hours:     Float64(1.2),
		Notes:     String("Writing tests"),
	})
	if err != nil {
		t.Errorf("CreateTimeEntryViaDuration task returned error: %v", err)
	}

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &TimeEntry{
		Id:                Int64(1),
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
	service, mux, _, teardown := setup()
	defer teardown()

	// https://help.getharvest.com/api-v2/timesheets-api/timesheets/time-entries/#create-a-time-entry-via-start-and-end-time

	mux.HandleFunc("/time_entries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{})
		testBody(t, r, `{"user_id":1782959,"project_id":14307913,"task_id":8083365,"spent_date":"2017-03-21T22:24:10Z","started_time":"8:00am","ended_time":"9:00am","notes":"Writing tests"}`+"\n")
		fmt.Fprint(w, `{"id":1,"name":"TimeEntry new","billable_by_default":true,"default_hourly_rate":123,"is_default":true,"is_active":true, "created_at":"2018-01-31T20:34:30Z","updated_at":"2018-05-31T21:34:30Z"}`)
	})

	spentDate := time.Date(2017, 3, 21, 22, 24, 10, 0, time.UTC)

	taskList, _, err := service.Timesheet.CreateTimeEntryViaStartEndTime(context.Background(), &TimeEntryCreateViaStartEndTime{
		UserId:      Int64(1782959),
		ProjectId:   Int64(14307913),
		TaskId:      Int64(8083365),
		SpentDate:   DateP(Date{spentDate}),
		StartedTime: TimeP(Time{time.Date(0, 1, 1, 8, 00, 00, 0, time.UTC)}),
		EndedTime:   TimeP(Time{time.Date(0, 1, 1, 9, 00, 00, 0, time.UTC)}),
		Notes:       String("Writing tests"),
	})
	if err != nil {
		t.Errorf("CreateTimeEntry task returned error: %v", err)
	}

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &TimeEntry{
		Id:        Int64(1),
		CreatedAt: &createdOne,
		UpdatedAt: &updatedOne,
	}

	assert.Equal(t, taskList, want)
}

func TestTimesheetService_DeleteTimeEntry(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/time_entries/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testFormValues(t, r, values{})
		testBody(t, r, ``)
		fmt.Fprint(w, ``)
	})

	_, err := service.Timesheet.DeleteTimeEntry(context.Background(), 1)
	if err != nil {
		t.Errorf("DeleteTimeEntry returned error: %v", err)
	}
}

func TestTimesheetService_GetTimeEntry(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/time_entries/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `{"id":1,"name":"TimeEntry new","billable_by_default":true,"default_hourly_rate":123,"is_default":true,"is_active":true, "created_at":"2018-01-31T20:34:30Z","updated_at":"2018-05-31T21:34:30Z"}`)
	})

	timeEntry, _, err := service.Timesheet.Get(context.Background(), 1)
	if err != nil {
		t.Errorf("TimeEntry.GetTimeEntry returned error: %v", err)
	}

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &TimeEntry{
		Id:        Int64(1),
		CreatedAt: &createdOne,
		UpdatedAt: &updatedOne,
	}

	assert.Equal(t, timeEntry, want)
}

func TestTimesheetService_ListTimeEntries(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/time_entries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `{
  "time_entries":[
    {
      "id":636709355,
      "spent_date":"2017-03-02",
      "user":{
        "id":1782959,
        "name":"Kim Allen"
      },
      "client":{
        "id":5735774,
        "name":"ABC Corp"
      },
      "project":{
        "id":14307913,
        "name":"Marketing Website"
      },
      "task":{
        "id":8083365,
        "name":"Graphic Design"
      },
      "user_assignment":{
        "id":125068553,
        "is_project_manager":true,
        "is_active":true,
        "budget":null,
        "created_at":"2017-06-26T22:32:52Z",
        "updated_at":"2017-06-26T22:32:52Z",
        "hourly_rate":100.0
      },
      "task_assignment":{
        "id":155502709,
        "billable":true,
        "is_active":true,
        "created_at":"2017-06-26T21:36:23Z",
        "updated_at":"2017-06-26T21:36:23Z",
        "hourly_rate":100.0,
        "budget":null
      },
      "hours":2.11,
      "rounded_hours": 2.25,
      "notes":"Adding CSS styling",
      "created_at":"2017-06-27T15:50:15Z",
      "updated_at":"2017-06-27T16:47:14Z",
      "is_locked":true,
      "locked_reason":"Item Approved and Locked for this Time Period",
      "is_closed":true,
      "is_billed":false,
      "timer_started_at":null,
      "started_time":"3:00pm",
      "ended_time":"5:00pm",
      "is_running":false,
      "invoice":null,
      "external_reference":null,
      "billable":true,
      "budgeted":true,
      "billable_rate":100.0,
      "cost_rate":50.0
    },
    {
      "id":636708723,
      "spent_date":"2017-03-01",
      "user":{
        "id":1782959,
        "name":"Kim Allen"
      },
      "client":{
        "id":5735776,
        "name":"123 Industries"
      },
      "project":{
        "id":14308069,
        "name":"Online Store - Phase 1"
      },
      "task":{
        "id":8083366,
        "name":"Programming"
      },
      "user_assignment":{
        "id":125068554,
        "is_project_manager":true,
        "is_active":true,
        "budget":null,
        "created_at":"2017-06-26T22:32:52Z",
        "updated_at":"2017-06-26T22:32:52Z",
        "hourly_rate":100.0
      },
      "task_assignment":{
        "id":155505014,
        "billable":true,
        "is_active":true,
        "created_at":"2017-06-26T21:52:18Z",
        "updated_at":"2017-06-26T21:52:18Z",
        "hourly_rate":100.0,
        "budget":null
      },
      "hours":1.35,
      "rounded_hours":1.5,
      "notes":"Importing products",
      "created_at":"2017-06-27T15:49:28Z",
      "updated_at":"2017-06-27T16:47:14Z",
      "is_locked":true,
      "locked_reason":"Item Invoiced and Approved and Locked for this Time Period",
      "is_closed":true,
      "is_billed":true,
      "timer_started_at":null,
      "started_time":"1:00pm",
      "ended_time":"2:00pm",
      "is_running":false,
      "invoice":{
        "id":13150403,
        "number":"1001"
      },
      "external_reference":null,
      "billable":true,
      "budgeted":true,
      "billable_rate":100.0,
      "cost_rate":50.0
    },
    {
      "id":636708574,
      "spent_date":"2017-03-01",
      "user":{
        "id":1782959,
        "name":"Kim Allen"
      },
      "client":{
        "id":5735776,
        "name":"123 Industries"
      },
      "project":{
        "id":14308069,
        "name":"Online Store - Phase 1"
      },
      "task":{
        "id":8083369,
        "name":"Research"
      },
      "user_assignment":{
        "id":125068554,
        "is_project_manager":true,
        "is_active":true,
        "budget":null,
        "created_at":"2017-06-26T22:32:52Z",
        "updated_at":"2017-06-26T22:32:52Z",
        "hourly_rate":100.0
      },
      "task_assignment":{
        "id":155505016,
        "billable":false,
        "is_active":true,
        "created_at":"2017-06-26T21:52:18Z",
        "updated_at":"2017-06-26T21:54:06Z",
        "hourly_rate":100.0,
        "budget":null
      },
      "hours":1.0,
      "rounded_hours":1.0,
      "notes":"Evaluating 3rd party libraries",
      "created_at":"2017-06-27T15:49:17Z",
      "updated_at":"2017-06-27T16:47:14Z",
      "is_locked":true,
      "locked_reason":"Item Approved and Locked for this Time Period",
      "is_closed":true,
      "is_billed":false,
      "timer_started_at":null,
      "started_time":"11:00am",
      "ended_time":"12:00pm",
      "is_running":false,
      "invoice":null,
      "external_reference":null,
      "billable":false,
      "budgeted":true,
      "billable_rate":null,
      "cost_rate":50.0
    },
    {
      "id":636707831,
      "spent_date":"2017-03-01",
      "user":{
        "id":1782959,
        "name":"Kim Allen"
      },
      "client":{
        "id":5735776,
        "name":"123 Industries"
      },
      "project":{
        "id":14308069,
        "name":"Online Store - Phase 1"
      },
      "task":{
        "id":8083368,
        "name":"Project Management"
      },
      "user_assignment":{
        "id":125068554,
        "is_project_manager":true,
        "is_active":true,
        "budget":null,
        "created_at":"2017-06-26T22:32:52Z",
        "updated_at":"2017-06-26T22:32:52Z",
        "hourly_rate":100.0
      },
      "task_assignment":{
        "id":155505015,
        "billable":true,
        "is_active":true,
        "created_at":"2017-06-26T21:52:18Z",
        "updated_at":"2017-06-26T21:52:18Z",
        "hourly_rate":100.0,
        "budget":null
      },
      "hours":2.0,
      "rounded_hours":2.0,
      "notes":"Planning meetings",
      "created_at":"2017-06-27T15:48:24Z",
      "updated_at":"2017-06-27T16:47:14Z",
      "is_locked":true,
      "locked_reason":"Item Invoiced and Approved and Locked for this Time Period",
      "is_closed":true,
      "is_billed":true,
      "timer_started_at":null,
      "started_time":"9:00am",
      "ended_time":"11:00am",
      "is_running":false,
      "invoice":{
        "id":13150403,
        "number":"1001"
      },
      "external_reference":null,
      "billable":true,
      "budgeted":true,
      "billable_rate":100.0,
      "cost_rate":50.0
    }
  ],
  "per_page":100,
  "total_pages":1,
  "total_entries":4,
  "next_page":null,
  "previous_page":null,
  "page":1,
  "links":{
    "first":"https://api.harvestapp.com/v2/time_entries?page=1&per_page=100",
    "next":null,
    "previous":null,
    "last":"https://api.harvestapp.com/v2/time_entries?page=1&per_page=100"
  }
}`)
	})

	taskList, _, err := service.Timesheet.List(context.Background(), &TimeEntryListOptions{})
	assert.NoError(t, err)

	want := &TimeEntryList{
		TimeEntries: []*TimeEntry{
			{
				Id: Int64(636709355),
				SpentDate: DateP(Date{time.Date(
					2017, 3, 2, 0, 0, 0, 0, time.Local)}),
				User: &User{
					Id: Int64(1782959),
				},
				Client: &Client{
					Id:   Int64(5735774),
					Name: String("ABC Corp"),
				},
				Project: &Project{
					Id:   Int64(14307913),
					Name: String("Marketing Website"),
				},
				Task: &Task{
					Id:   Int64(8083365),
					Name: String("Graphic Design"),
				},
				UserAssignment: &ProjectUserAssignment{
					Id:               Int64(125068553),
					User:             nil,
					IsActive:         Bool(true),
					IsProjectManager: Bool(true),
					HourlyRate:       Float64(100),
					Budget:           nil,
					CreatedAt:        TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
					UpdatedAt:        TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
				},
				TaskAssignment: &ProjectTaskAssignment{
					Id:         Int64(155502709),
					Task:       nil,
					IsActive:   Bool(true),
					Billable:   Bool(true),
					HourlyRate: Float64(100),
					Budget:     nil,
					CreatedAt:  TimeTimeP(time.Date(2017, 6, 26, 21, 36, 23, 0, time.UTC)),
					UpdatedAt:  TimeTimeP(time.Date(2017, 6, 26, 21, 36, 23, 0, time.UTC)),
				},
				ExternalReference: nil,
				Invoice:           nil,
				Hours:             Float64(2.11),
				RoundedHours:      Float64(2.25),
				Notes:             String("Adding CSS styling"),
				IsLocked:          Bool(true),
				LockedReason:      String("Item Approved and Locked for this Time Period"),
				IsClosed:          Bool(true),
				IsBilled:          Bool(false),
				TimerStartedAt:    nil,
				StartedTime:       TimeP(Time{time.Date(0, 1, 1, 15, 00, 00, 0, time.Local)}),
				EndedTime:         TimeP(Time{time.Date(0, 1, 1, 17, 00, 00, 0, time.Local)}),
				IsRunning:         Bool(false),
				Billable:          Bool(true),
				Budgeted:          Bool(true),
				BillableRate:      Float64(100),
				CostRate:          Float64(50),
				CreatedAt:         TimeTimeP(time.Date(2017, 6, 27, 15, 50, 15, 0, time.UTC)),
				UpdatedAt:         TimeTimeP(time.Date(2017, 6, 27, 16, 47, 14, 0, time.UTC)),
			}, {
				Id: Int64(636708723),
				SpentDate: DateP(Date{time.Date(
					2017, 3, 1, 0, 0, 0, 0, time.Local)}),
				User: &User{
					Id: Int64(1782959),
				},
				Client: &Client{
					Id:   Int64(5735776),
					Name: String("123 Industries"),
				},
				Project: &Project{
					Id:   Int64(14308069),
					Name: String("Online Store - Phase 1"),
				},
				Task: &Task{
					Id:   Int64(8083366),
					Name: String("Programming"),
				},
				UserAssignment: &ProjectUserAssignment{
					Id:               Int64(125068554),
					IsActive:         Bool(true),
					IsProjectManager: Bool(true),
					HourlyRate:       Float64(100),
					CreatedAt:        TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
					UpdatedAt:        TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
				},
				TaskAssignment: &ProjectTaskAssignment{
					Id:         Int64(155505014),
					IsActive:   Bool(true),
					Billable:   Bool(true),
					HourlyRate: Float64(100),
					CreatedAt:  TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
					UpdatedAt:  TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
				},

				Invoice: &Invoice{
					Id:     Int64(13150403),
					Number: String("1001"),
				},
				Hours:        Float64(1.35),
				RoundedHours: Float64(1.5),
				Notes:        String("Importing products"),
				IsLocked:     Bool(true),
				LockedReason: String("Item Invoiced and Approved and Locked for this Time Period"),
				IsClosed:     Bool(true),
				IsBilled:     Bool(true),
				StartedTime:  TimeP(Time{time.Date(0, 1, 1, 13, 00, 00, 0, time.Local)}),
				EndedTime:    TimeP(Time{time.Date(0, 1, 1, 14, 00, 00, 0, time.Local)}),
				IsRunning:    Bool(false),

				Billable:     Bool(true),
				Budgeted:     Bool(true),
				BillableRate: Float64(100),
				CostRate:     Float64(50),
				CreatedAt:    TimeTimeP(time.Date(2017, 6, 27, 15, 49, 28, 0, time.UTC)),
				UpdatedAt:    TimeTimeP(time.Date(2017, 6, 27, 16, 47, 14, 0, time.UTC)),
			}, {
				Id: Int64(636708574),
				SpentDate: DateP(Date{time.Date(
					2017, 3, 1, 0, 0, 0, 0, time.Local)}),
				User: &User{
					Id: Int64(1782959),
				},
				Client: &Client{
					Id:   Int64(5735776),
					Name: String("123 Industries"),
				},
				Project: &Project{
					Id:   Int64(14308069),
					Name: String("Online Store - Phase 1"),
				},
				Task: &Task{
					Id:   Int64(8083369),
					Name: String("Research"),
				},
				UserAssignment: &ProjectUserAssignment{
					Id:               Int64(125068554),
					IsActive:         Bool(true),
					IsProjectManager: Bool(true),
					HourlyRate:       Float64(100),
					CreatedAt:        TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
					UpdatedAt:        TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
				},
				TaskAssignment: &ProjectTaskAssignment{
					Id:         Int64(155505016),
					IsActive:   Bool(true),
					Billable:   Bool(false),
					HourlyRate: Float64(100),
					CreatedAt:  TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
					UpdatedAt:  TimeTimeP(time.Date(2017, 6, 26, 21, 54, 06, 0, time.UTC)),
				},

				Invoice:      nil,
				Hours:        Float64(1),
				RoundedHours: Float64(1),
				Notes:        String("Evaluating 3rd party libraries"),
				IsLocked:     Bool(true),
				LockedReason: String("Item Approved and Locked for this Time Period"),
				IsClosed:     Bool(true),
				IsBilled:     Bool(false),
				StartedTime:  TimeP(Time{time.Date(0, 1, 1, 11, 00, 00, 0, time.Local)}),
				EndedTime:    TimeP(Time{time.Date(0, 1, 1, 12, 00, 00, 0, time.Local)}),
				IsRunning:    Bool(false),

				Billable:     Bool(false),
				Budgeted:     Bool(true),
				BillableRate: nil,
				CostRate:     Float64(50),
				CreatedAt:    TimeTimeP(time.Date(2017, 6, 27, 15, 49, 17, 0, time.UTC)),
				UpdatedAt:    TimeTimeP(time.Date(2017, 6, 27, 16, 47, 14, 0, time.UTC)),
			}, {
				Id: Int64(636707831),
				SpentDate: DateP(Date{time.Date(
					2017, 3, 1, 0, 0, 0, 0, time.Local)}),
				User: &User{
					Id: Int64(1782959),
				},
				Client: &Client{
					Id:   Int64(5735776),
					Name: String("123 Industries"),
				},
				Project: &Project{
					Id:   Int64(14308069),
					Name: String("Online Store - Phase 1"),
				},
				Task: &Task{
					Id:   Int64(8083368),
					Name: String("Project Management"),
				},
				UserAssignment: &ProjectUserAssignment{
					Id:               Int64(125068554),
					IsActive:         Bool(true),
					IsProjectManager: Bool(true),
					HourlyRate:       Float64(100),
					CreatedAt:        TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
					UpdatedAt:        TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
				},
				TaskAssignment: &ProjectTaskAssignment{
					Id:         Int64(155505015),
					IsActive:   Bool(true),
					Billable:   Bool(true),
					HourlyRate: Float64(100),
					CreatedAt:  TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
					UpdatedAt:  TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
				},

				Invoice: &Invoice{
					Id:     Int64(13150403),
					Number: String("1001"),
				},
				Hours:        Float64(2),
				RoundedHours: Float64(2),
				Notes:        String("Planning meetings"),
				IsLocked:     Bool(true),
				LockedReason: String("Item Invoiced and Approved and Locked for this Time Period"),
				IsClosed:     Bool(true),
				IsBilled:     Bool(true),
				StartedTime:  TimeP(Time{time.Date(0, 1, 1, 9, 00, 00, 0, time.Local)}),
				EndedTime:    TimeP(Time{time.Date(0, 1, 1, 11, 00, 00, 0, time.Local)}),
				IsRunning:    Bool(false),

				Billable:     Bool(true),
				Budgeted:     Bool(true),
				BillableRate: Float64(100),
				CostRate:     Float64(50),
				CreatedAt:    TimeTimeP(time.Date(2017, 6, 27, 15, 48, 24, 0, time.UTC)),
				UpdatedAt:    TimeTimeP(time.Date(2017, 6, 27, 16, 47, 14, 0, time.UTC)),
			},
		},

		Pagination: Pagination{
			PerPage:      Int(100),
			TotalPages:   Int(1),
			TotalEntries: Int(4),
			NextPage:     nil,
			PreviousPage: nil,
			Page:         Int(1),
			Links: &PageLinks{
				First:    String("https://api.harvestapp.com/v2/time_entries?page=1&per_page=100"),
				Next:     nil,
				Previous: nil,
				Last:     String("https://api.harvestapp.com/v2/time_entries?page=1&per_page=100"),
			},
		},
	}

	assert.Equal(t, want, taskList)
}

func TestTimesheetService_RestartTimeEntry(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/time_entries/662202797/restart", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `{"id":662202797,"spent_date":"2017-03-21","user":{"id":1795925,"name":"Jane Smith"},"client":{"id":5735776,"name":"123 Industries"},"project":{"id":14808188,"name":"Task Force"},"task":{"id":8083366,"name":"Programming"},"user_assignment":{"id":130403296,"is_project_manager":true,"is_active":true,"budget":null,"created_at":"2017-08-22T17:36:54Z","updated_at":"2017-08-22T17:36:54Z","hourly_rate":100},"task_assignment":{"id":160726645,"billable":true,"is_active":true,"created_at":"2017-08-22T17:36:54Z","updated_at":"2017-08-22T17:36:54Z","hourly_rate":100,"budget":null},"hours":0,"rounded_hours":0,"notes":null,"created_at":"2017-08-22T17:40:24Z","updated_at":"2017-08-22T17:40:24Z","is_locked":false,"locked_reason":null,"is_closed":false,"is_billed":false,"timer_started_at":"2017-08-22T17:40:24Z","started_time":"11:40am","ended_time":null,"is_running":true,"invoice":null,"external_reference":null,"billable":true,"budgeted":false,"billable_rate":100,"cost_rate":75}`)
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
	startedTime := Time{time.Date(0, 1, 1, 11, 40, 00, 0, time.Local)}

	want := &TimeEntry{
		Id:        Int64(662202797),
		SpentDate: DateP(Date{spentDate}),
		User: &User{
			Id: Int64(1795925),
		},
		Client: &Client{
			Id:   Int64(5735776),
			Name: String("123 Industries"),
		},
		Project: &Project{
			Id:   Int64(14808188),
			Name: String("Task Force"),
		},
		Task: &Task{
			Id:   Int64(8083366),
			Name: String("Programming"),
		},
		UserAssignment: &ProjectUserAssignment{
			Id:               Int64(130403296),
			IsProjectManager: Bool(true),
			IsActive:         Bool(true),
			Budget:           nil,
			CreatedAt:        TimeTimeP(assignmentTime),
			UpdatedAt:        TimeTimeP(assignmentTime),
			HourlyRate:       Float64(100),
		},
		TaskAssignment: &ProjectTaskAssignment{
			Id:         Int64(160726645),
			Billable:   Bool(true),
			IsActive:   Bool(true),
			CreatedAt:  TimeTimeP(assignmentTime),
			UpdatedAt:  TimeTimeP(assignmentTime),
			HourlyRate: Float64(100),
			Budget:     nil,
		},
		Hours:             Float64(0),
		RoundedHours:      Float64(0),
		Notes:             nil,
		CreatedAt:         &createdOne,
		UpdatedAt:         &updatedOne,
		IsLocked:          Bool(false),
		LockedReason:      nil,
		IsClosed:          Bool(false),
		IsBilled:          Bool(false),
		TimerStartedAt:    TimeTimeP(timerStartedAt),
		StartedTime:       TimeP(startedTime),
		EndedTime:         nil,
		IsRunning:         Bool(true),
		Invoice:           nil,
		ExternalReference: nil,
		Billable:          Bool(true),
		Budgeted:          Bool(false),
		BillableRate:      Float64(100),
		CostRate:          Float64(75),
	}

	assert.Equal(t, want, timeEntry)
}

func TestTimesheetService_StopTimeEntry(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/time_entries/662202797/stop", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `{"id":662202797,"spent_date":"2017-03-21","user":{"id":1795925,"name":"Jane Smith"},"client":{"id":5735776,"name":"123 Industries"},"project":{"id":14808188,"name":"Task Force"},"task":{"id":8083366,"name":"Programming"},"user_assignment":{"id":130403296,"is_project_manager":true,"is_active":true,"budget":null,"created_at":"2017-08-22T17:36:54Z","updated_at":"2017-08-22T17:36:54Z","hourly_rate":100},"task_assignment":{"id":160726645,"billable":true,"is_active":true,"created_at":"2017-08-22T17:36:54Z","updated_at":"2017-08-22T17:36:54Z","hourly_rate":100,"budget":null},"hours":0.02,"rounded_hours":0.25,"notes":null,"created_at":"2017-08-22T17:37:13Z","updated_at":"2017-08-22T17:38:31Z","is_locked":false,"locked_reason":null,"is_closed":false,"is_billed":false,"timer_started_at":null,"started_time":"11:37am","ended_time":"11:38am","is_running":false,"invoice":null,"external_reference":null,"billable":true,"budgeted":false,"billable_rate":100,"cost_rate":75}`)
	})

	timeEntry, _, err := service.Timesheet.StopTimeEntry(context.Background(), 662202797)
	assert.NoError(t, err)

	spentDate := time.Date(2017, 3, 21, 0, 0, 0, 0, time.Local)

	createdOne := time.Date(2017, 8, 22, 17, 37, 13, 0, time.UTC)
	updatedOne := time.Date(2017, 8, 22, 17, 38, 31, 0, time.UTC)

	assignmentTime := time.Date(2017, 8, 22, 17, 36, 54, 0, time.UTC)

	startedTime := Time{time.Date(0, 1, 1, 11, 37, 00, 0, time.Local)}
	endedTime := Time{time.Date(0, 1, 1, 11, 38, 00, 0, time.Local)}

	want := &TimeEntry{
		Id:        Int64(662202797),
		SpentDate: DateP(Date{spentDate}),
		User: &User{
			Id: Int64(1795925),
		},
		Client: &Client{
			Id:   Int64(5735776),
			Name: String("123 Industries"),
		},
		Project: &Project{
			Id:   Int64(14808188),
			Name: String("Task Force"),
		},
		Task: &Task{
			Id:   Int64(8083366),
			Name: String("Programming"),
		},
		UserAssignment: &ProjectUserAssignment{
			Id:               Int64(130403296),
			IsProjectManager: Bool(true),
			IsActive:         Bool(true),
			Budget:           nil,
			CreatedAt:        TimeTimeP(assignmentTime),
			UpdatedAt:        TimeTimeP(assignmentTime),
			HourlyRate:       Float64(100),
		},
		TaskAssignment: &ProjectTaskAssignment{
			Id:         Int64(160726645),
			Billable:   Bool(true),
			IsActive:   Bool(true),
			CreatedAt:  TimeTimeP(assignmentTime),
			UpdatedAt:  TimeTimeP(assignmentTime),
			HourlyRate: Float64(100),
			Budget:     nil,
		},
		Hours:             Float64(0.02),
		RoundedHours:      Float64(0.25),
		Notes:             nil,
		CreatedAt:         &createdOne,
		UpdatedAt:         &updatedOne,
		IsLocked:          Bool(false),
		LockedReason:      nil,
		IsClosed:          Bool(false),
		IsBilled:          Bool(false),
		TimerStartedAt:    nil,
		StartedTime:       TimeP(startedTime),
		EndedTime:         TimeP(endedTime),
		IsRunning:         Bool(false),
		Invoice:           nil,
		ExternalReference: nil,
		Billable:          Bool(true),
		Budgeted:          Bool(false),
		BillableRate:      Float64(100),
		CostRate:          Float64(75),
	}

	assert.Equal(t, want, timeEntry)
}

func TestTimesheetService_UpdateTimeEntry(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/time_entries/636718192", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		testBody(t, r, `{"project_id":1234,"task_id":2345,"spent_date":"2017-03-21T00:00:00Z","started_time":"11:40am","ended_time":"12:45pm","hours":1,"notes":"new notes"}`+"\n")
		fmt.Fprint(w, `{"id":636718192,"spent_date":"2017-03-21","user":{"id":1782959,"name":"Kim Allen"},"client":{"id":5735774,"name":"ABC Corp"},"project":{"id":14307913,"name":"Marketing Website"},"task":{"id":8083365,"name":"Graphic Design"},"user_assignment":{"id":125068553,"is_project_manager":true,"is_active":true,"budget":null,"created_at":"2017-06-26T22:32:52Z","updated_at":"2017-06-26T22:32:52Z","hourly_rate":100},"task_assignment":{"id":155502709,"billable":true,"is_active":true,"created_at":"2017-06-26T21:36:23Z","updated_at":"2017-06-26T21:36:23Z","hourly_rate":100,"budget":null},"hours":1,"rounded_hours":1,"notes":"Updated notes","created_at":"2017-06-27T16:01:23Z","updated_at":"2017-06-27T16:02:40Z","is_locked":false,"locked_reason":null,"is_closed":false,"is_billed":false,"timer_started_at":null,"started_time":null,"ended_time":null,"is_running":false,"invoice":null,"external_reference":null,"billable":true,"budgeted":true,"billable_rate":100,"cost_rate":50}`)
	})

	spentDate := time.Date(2017, 3, 21, 0, 0, 0, 0, time.Local)
	startedTime := Time{time.Date(0, 1, 1, 11, 40, 00, 0, time.UTC)}
	endedTime := Time{time.Date(0, 1, 1, 12, 45, 10, 0, time.UTC)}

	timeEntry, _, err := service.Timesheet.UpdateTimeEntry(context.Background(), 636718192, &TimeEntryUpdate{
		ProjectId:   Int64(1234),
		TaskId:      Int64(2345),
		SpentDate:   DateP(Date{spentDate}),
		StartedTime: &startedTime,
		EndedTime:   &endedTime,
		Hours:       Float64(1),
		Notes:       String("new notes"),
	})
	assert.NoError(t, err)

	createdTimeEntry := time.Date(
		2017, 6, 27, 16, 01, 23, 0, time.UTC)
	updatedTimeEntry := time.Date(
		2017, 6, 27, 16, 02, 40, 0, time.UTC)
	createdUserAssignment := time.Date(
		2017, 6, 26, 22, 32, 52, 0, time.UTC)
	updatedUserAssignment := time.Date(
		2017, 6, 26, 22, 32, 52, 0, time.UTC)
	createdTaskAssignment := time.Date(
		2017, 6, 26, 21, 36, 23, 0, time.UTC)
	updatedTaskAssignment := time.Date(
		2017, 6, 26, 21, 36, 23, 0, time.UTC)

	want := &TimeEntry{
		Id:        Int64(636718192),
		CreatedAt: &createdTimeEntry,
		UpdatedAt: &updatedTimeEntry,
		SpentDate: DateP(Date{spentDate}),
		User: &User{
			Id: Int64(1782959),
		},
		Client: &Client{
			Id:   Int64(5735774),
			Name: String("ABC Corp"),
		},
		Project: &Project{
			Id:   Int64(14307913),
			Name: String("Marketing Website"),
		},
		Task: &Task{
			Id:   Int64(8083365),
			Name: String("Graphic Design"),
		},
		UserAssignment: &ProjectUserAssignment{
			Id:               Int64(125068553),
			IsProjectManager: Bool(true),
			IsActive:         Bool(true),
			CreatedAt:        &createdUserAssignment,
			UpdatedAt:        &updatedUserAssignment,
			HourlyRate:       Float64(100),
		},
		TaskAssignment: &ProjectTaskAssignment{
			Id:         Int64(155502709),
			Billable:   Bool(true),
			IsActive:   Bool(true),
			CreatedAt:  &createdTaskAssignment,
			UpdatedAt:  &updatedTaskAssignment,
			HourlyRate: Float64(100),
		},
		Hours:        Float64(1),
		RoundedHours: Float64(1),
		Notes:        String("Updated notes"),
		IsLocked:     Bool(false),
		IsClosed:     Bool(false),
		IsBilled:     Bool(false),
		IsRunning:    Bool(false),
		Billable:     Bool(true),
		Budgeted:     Bool(true),
		BillableRate: Float64(100),
		CostRate:     Float64(50),
	}

	assert.Equal(t, want, timeEntry)
}
