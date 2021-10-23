package harvest_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/becoded/go-harvest/harvest"
	"github.com/stretchr/testify/assert"
)

func TestTaskService_Create(t *testing.T) {
	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	tests := []struct {
		name       string
		args       *harvest.TaskCreateRequest
		method     string
		path       string
		formValues values
		body       string
		response   string
		want       *harvest.Task
		wantErr    error
	}{
		{
			name: "create task",
			args: &harvest.TaskCreateRequest{
				Name:              harvest.String("Task new"),
				BillableByDefault: harvest.Bool(true),
				DefaultHourlyRate: harvest.Float64(123),
				IsDefault:         harvest.Bool(true),
				IsActive:          harvest.Bool(true),
			},
			method:     "POST",
			path:       "/tasks",
			formValues: values{},
			body:       `{"name":"Task new","billable_by_default":true,"default_hourly_rate":123,"is_default":true,"is_active":true}` + "\n",
			response:   `{"id":1,"name":"Task new","billable_by_default":true,"default_hourly_rate":123,"is_default":true,"is_active":true, "created_at":"2018-01-31T20:34:30Z","updated_at":"2018-05-31T21:34:30Z"}`,
			want: &harvest.Task{
				ID:                harvest.Int64(1),
				Name:              harvest.String("Task new"),
				BillableByDefault: harvest.Bool(true),
				DefaultHourlyRate: harvest.Float64(123),
				IsDefault:         harvest.Bool(true),
				IsActive:          harvest.Bool(true),
				CreatedAt:         &createdOne,
				UpdatedAt:         &updatedOne,
			},
			wantErr: nil,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.method)
				testFormValues(t, r, tt.formValues)
				testBody(t, r, tt.body)
				_, err := fmt.Fprint(w, tt.response)
				assert.NoError(t, err)
			})

			task, _, err := service.Task.Create(context.Background(), tt.args)

			if tt.wantErr != nil {
				assert.Nil(t, task)
				assert.EqualError(t, err, tt.wantErr.Error())

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, task)
		})
	}
}

func TestTaskService_Delete(t *testing.T) {
	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	type args struct {
		taskID int64
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
				taskID: 1,
			},
			method:     "DELETE",
			path:       "/tasks/1",
			formValues: values{},
			body:       "",
			response:   "",
			wantErr:    nil,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.method)
				testFormValues(t, r, tt.formValues)
				testBody(t, r, tt.body)
				_, err := fmt.Fprint(w, tt.response)
				assert.NoError(t, err)
			})

			_, err := service.Task.Delete(context.Background(), tt.args.taskID)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())

				return
			}

			assert.NoError(t, err)
			assert.Nil(t, err)
		})
	}
}

func TestTaskService_Get(t *testing.T) {
	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	type args struct {
		taskID int64
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
		want       *harvest.Task
		wantErr    error
	}{
		{
			name: "get 1",
			args: args{
				taskID: 1,
			},
			method:     "GET",
			path:       "/tasks/1",
			formValues: values{},
			body:       "",
			response:   `{"id":1,"name":"Task new","billable_by_default":true,"default_hourly_rate":123,"is_default":true,"is_active":true, "created_at":"2018-01-31T20:34:30Z","updated_at":"2018-05-31T21:34:30Z"}`,
			want: &harvest.Task{
				ID:                harvest.Int64(1),
				Name:              harvest.String("Task new"),
				BillableByDefault: harvest.Bool(true),
				DefaultHourlyRate: harvest.Float64(123),
				IsDefault:         harvest.Bool(true),
				IsActive:          harvest.Bool(true),
				CreatedAt:         &createdOne,
				UpdatedAt:         &updatedOne,
			},
			wantErr: nil,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.method)
				testFormValues(t, r, tt.formValues)
				testBody(t, r, tt.body)
				_, err := fmt.Fprint(w, tt.response)
				assert.NoError(t, err)
			})

			task, _, err := service.Task.Get(context.Background(), tt.args.taskID)

			if tt.wantErr != nil {
				assert.Nil(t, task)
				assert.EqualError(t, err, tt.wantErr.Error())

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, task)
		})
	}
}

func TestTaskService_List(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `{"tasks":[{"id":1,"name":"Task 1","billable_by_default":true,"default_hourly_rate":123,"is_default":true,"is_active":true, "created_at":"2018-01-31T20:34:30Z","updated_at":"2018-05-31T21:34:30Z"},{"id":2,"name":"Task 2","billable_by_default":false,"default_hourly_rate":321,"is_default":false,"is_active":false, "created_at":"2018-03-02T10:12:13Z","updated_at":"2018-04-30T12:13:14Z"}],"per_page":100,"total_pages":1,"total_entries":2,"next_page":null,"previous_page":null,"page":1,"links":{"first":"https://api.harvestapp.com/v2/tasks?page=1&per_page=100","next":null,"previous":null,"last":"https://api.harvestapp.com/v2/tasks?page=1&per_page=100"}}`)
	})

	taskList, _, err := service.Task.List(context.Background(), &harvest.TaskListOptions{})
	assert.NoError(t, err)

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)
	createdTwo := time.Date(
		2018, 3, 2, 10, 12, 13, 0, time.UTC)
	updatedTwo := time.Date(
		2018, 4, 30, 12, 13, 14, 0, time.UTC)

	want := &harvest.TaskList{
		Tasks: []*harvest.Task{
			{
				ID:                harvest.Int64(1),
				Name:              harvest.String("Task 1"),
				BillableByDefault: harvest.Bool(true),
				DefaultHourlyRate: harvest.Float64(123),
				IsDefault:         harvest.Bool(true),
				IsActive:          harvest.Bool(true),
				CreatedAt:         &createdOne,
				UpdatedAt:         &updatedOne,
			}, {
				ID:                harvest.Int64(2),
				Name:              harvest.String("Task 2"),
				BillableByDefault: harvest.Bool(false),
				DefaultHourlyRate: harvest.Float64(321),
				IsDefault:         harvest.Bool(false),
				IsActive:          harvest.Bool(false),
				CreatedAt:         &createdTwo,
				UpdatedAt:         &updatedTwo,
			},
		},
		Pagination: harvest.Pagination{
			PerPage:      harvest.Int(100),
			TotalPages:   harvest.Int(1),
			TotalEntries: harvest.Int(2),
			NextPage:     nil,
			PreviousPage: nil,
			Page:         harvest.Int(1),
			Links: &harvest.PageLinks{
				First:    harvest.String("https://api.harvestapp.com/v2/tasks?page=1&per_page=100"),
				Next:     nil,
				Previous: nil,
				Last:     harvest.String("https://api.harvestapp.com/v2/tasks?page=1&per_page=100"),
			},
		},
	}

	assert.Equal(t, want, taskList)
}

func TestTaskService_Update(t *testing.T) {
	t.Parallel()
	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/tasks/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		testBody(t, r, `{"name":"Task update","billable_by_default":false,"default_hourly_rate":213,"is_default":false,"is_active":false}`+"\n")
		fmt.Fprint(w, `{"id":1,"name":"Task update","billable_by_default":false,"default_hourly_rate":213,"is_default":false,"is_active":false, "created_at":"2018-01-31T20:34:30Z","updated_at":"2018-05-31T21:34:30Z"},{"id":2,"name":"Task 2","billable_by_default":false,"default_hourly_rate":321,"is_default":false,"is_active":false, "created_at":"2018-03-02T10:12:13Z","updated_at":"2018-04-30T12:13:14Z"}`)
	})

	task, _, err := service.Task.Update(context.Background(), 1, &harvest.TaskUpdateRequest{
		Name:              harvest.String("Task update"),
		BillableByDefault: harvest.Bool(false),
		DefaultHourlyRate: harvest.Float64(213),
		IsDefault:         harvest.Bool(false),
		IsActive:          harvest.Bool(false),
	})
	assert.NoError(t, err)

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &harvest.Task{
		ID:                harvest.Int64(1),
		Name:              harvest.String("Task update"),
		BillableByDefault: harvest.Bool(false),
		DefaultHourlyRate: harvest.Float64(213),
		IsDefault:         harvest.Bool(false),
		IsActive:          harvest.Bool(false),
		CreatedAt:         &createdOne,
		UpdatedAt:         &updatedOne,
	}

	assert.Equal(t, want, task)
}
