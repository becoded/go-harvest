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

func TestTaskService_CreateTask(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()
	
	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	tests := []struct {
		name       string
		args       *TaskCreateRequest
		method     string
		path       string
		formValues values
		body       string
		response   string
		want       *Task
		wantErr    error
	}{
		{
			name: "create task",
			args: &TaskCreateRequest{
				Name:              String("Task new"),
				BillableByDefault: Bool(true),
				DefaultHourlyRate: Float64(123),
				IsDefault:         Bool(true),
				IsActive:          Bool(true),
			},
			method:     "POST",
			path:       "/tasks",
			formValues: values{},
			body:       `{"name":"Task new","billable_by_default":true,"default_hourly_rate":123,"is_default":true,"is_active":true}`+"\n",
			response:   `{"id":1,"name":"Task new","billable_by_default":true,"default_hourly_rate":123,"is_default":true,"is_active":true, "created_at":"2018-01-31T20:34:30Z","updated_at":"2018-05-31T21:34:30Z"}`,
			want: &Task{
				Id:                Int64(1),
				Name:              String("Task new"),
				BillableByDefault: Bool(true),
				DefaultHourlyRate: Float64(123),
				IsDefault:         Bool(true),
				IsActive:          Bool(true),
				CreatedAt:         &createdOne,
				UpdatedAt:         &updatedOne,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.method)
				testFormValues(t, r, tt.formValues)
				testBody(t, r, tt.body)
				_, err := fmt.Fprint(w, tt.response)
				assert.NoError(t, err)
			})

			task, _, err := service.Task.CreateTask(context.Background(), tt.args)

			if tt.wantErr != nil {
				assert.Nil(t, task)
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}

			assert.Equal(t, tt.want, task)
		})
	}
}

func TestTaskService_DeleteTask(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	type args struct {
		taskId int64
	}

	tests := []struct {
		name       string
		args       args
		method     string
		path       string
		formValues values
		body       string
		response   string
		wantErr    error
	}{
		{
			name: "delete 1",
			args: args{
				taskId: 1,
			},
			method:     "DELETE",
			path:       "/tasks/1",
			formValues: values{},
			body:       "",
			response:   "",
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.method)
				testFormValues(t, r, tt.formValues)
				testBody(t, r, tt.body)
				_, err := fmt.Fprint(w, tt.response)
				assert.NoError(t, err)
			})

			_, err := service.Task.DeleteTask(context.Background(), tt.args.taskId)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestTaskService_GetTask(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	type args struct {
		taskId int64
	}

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	tests := []struct {
		name       string
		args       args
		method     string
		path       string
		formValues values
		body       string
		response   string
		want       *Task
		wantErr    error
	}{
		{
			name: "get 1",
			args: args{
				taskId: 1,
			},
			method:     "GET",
			path:       "/tasks/1",
			formValues: values{},
			body:       "",
			response:   `{"id":1,"name":"Task new","billable_by_default":true,"default_hourly_rate":123,"is_default":true,"is_active":true, "created_at":"2018-01-31T20:34:30Z","updated_at":"2018-05-31T21:34:30Z"}`,
			want: &Task{
				Id:                Int64(1),
				Name:              String("Task new"),
				BillableByDefault: Bool(true),
				DefaultHourlyRate: Float64(123),
				IsDefault:         Bool(true),
				IsActive:          Bool(true),
				CreatedAt:         &createdOne,
				UpdatedAt:         &updatedOne,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.method)
				testFormValues(t, r, tt.formValues)
				testBody(t, r, tt.body)
				_, err := fmt.Fprint(w, tt.response)
				assert.NoError(t, err)
			})

			task, _, err := service.Task.Get(context.Background(), tt.args.taskId)

			if tt.wantErr != nil {
				assert.Nil(t, task)
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}

			assert.Equal(t, tt.want, task)
		})
	}
}

func TestTaskService_ListTasks(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `{"tasks":[{"id":1,"name":"Task 1","billable_by_default":true,"default_hourly_rate":123,"is_default":true,"is_active":true, "created_at":"2018-01-31T20:34:30Z","updated_at":"2018-05-31T21:34:30Z"},{"id":2,"name":"Task 2","billable_by_default":false,"default_hourly_rate":321,"is_default":false,"is_active":false, "created_at":"2018-03-02T10:12:13Z","updated_at":"2018-04-30T12:13:14Z"}],"per_page":100,"total_pages":1,"total_entries":2,"next_page":null,"previous_page":null,"page":1,"links":{"first":"https://api.harvestapp.com/v2/tasks?page=1&per_page=100","next":null,"previous":null,"last":"https://api.harvestapp.com/v2/tasks?page=1&per_page=100"}}`)
	})

	taskList, _, err := service.Task.List(context.Background(), &TaskListOptions{})
	if err != nil {
		t.Errorf("Task.ListTakes returned error: %v", err)
	}

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)
	createdTwo := time.Date(
		2018, 3, 2, 10, 12, 13, 0, time.UTC)
	updatedTwo := time.Date(
		2018, 4, 30, 12, 13, 14, 0, time.UTC)

	want := &TaskList{
		Tasks: []*Task{
			{
				Id:                Int64(1),
				Name:              String("Task 1"),
				BillableByDefault: Bool(true),
				DefaultHourlyRate: Float64(123),
				IsDefault:         Bool(true),
				IsActive:          Bool(true),
				CreatedAt:         &createdOne,
				UpdatedAt:         &updatedOne,
			}, {
				Id:                Int64(2),
				Name:              String("Task 2"),
				BillableByDefault: Bool(false),
				DefaultHourlyRate: Float64(321),
				IsDefault:         Bool(false),
				IsActive:          Bool(false),
				CreatedAt:         &createdTwo,
				UpdatedAt:         &updatedTwo,
			}},
		Pagination: Pagination{
			PerPage:      Int(100),
			TotalPages:   Int(1),
			TotalEntries: Int(2),
			NextPage:     nil,
			PreviousPage: nil,
			Page:         Int(1),
			Links: &PageLinks{
				First:    String("https://api.harvestapp.com/v2/tasks?page=1&per_page=100"),
				Next:     nil,
				Previous: nil,
				Last:     String("https://api.harvestapp.com/v2/tasks?page=1&per_page=100"),
			},
		},
	}

	if !reflect.DeepEqual(taskList, want) {
		t.Errorf("Task.ListTakes returned %+v, response %+v", taskList, want)
	}
}

func TestTaskService_UpdateTask(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/tasks/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		testBody(t, r, `{"name":"Task update","billable_by_default":false,"default_hourly_rate":213,"is_default":false,"is_active":false}`+"\n")
		fmt.Fprint(w, `{"id":1,"name":"Task update","billable_by_default":false,"default_hourly_rate":213,"is_default":false,"is_active":false, "created_at":"2018-01-31T20:34:30Z","updated_at":"2018-05-31T21:34:30Z"},{"id":2,"name":"Task 2","billable_by_default":false,"default_hourly_rate":321,"is_default":false,"is_active":false, "created_at":"2018-03-02T10:12:13Z","updated_at":"2018-04-30T12:13:14Z"}`)
	})

	taskList, _, err := service.Task.UpdateTask(context.Background(), 1, &TaskUpdateRequest{
		Name:              String("Task update"),
		BillableByDefault: Bool(false),
		DefaultHourlyRate: Float64(213),
		IsDefault:         Bool(false),
		IsActive:          Bool(false),
	})
	if err != nil {
		t.Errorf("UpdateTask returned error: %v", err)
	}

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &Task{
		Id:                Int64(1),
		Name:              String("Task update"),
		BillableByDefault: Bool(false),
		DefaultHourlyRate: Float64(213),
		IsDefault:         Bool(false),
		IsActive:          Bool(false),
		CreatedAt:         &createdOne,
		UpdatedAt:         &updatedOne,
	}

	if !reflect.DeepEqual(taskList, want) {
		t.Errorf("Task.UpdateTask returned %+v, response %+v", taskList, want)
	}
}
