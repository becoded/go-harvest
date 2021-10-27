package harvest

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskService_Create(t *testing.T) {
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
			body:       "task/create/body_1.json",
			response:   "task/create/response_1.json",
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
				testWriteResponse(t, w, tt.response)
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
			body:       "task/delete/body_1.json",
			response:   "task/delete/response_1.json",
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
				testWriteResponse(t, w, tt.response)
			})

			_, err := service.Task.Delete(context.Background(), tt.args.taskId)

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
			body:       "task/get/body_1.json",
			response:   "task/get/response_1.json",
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
				testWriteResponse(t, w, tt.response)
			})

			task, _, err := service.Task.Get(context.Background(), tt.args.taskId)

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
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		testBody(t, r, "task/list/body_1.json")
		testWriteResponse(t, w, "task/list/response_1.json")
	})

	taskList, _, err := service.Task.List(context.Background(), &TaskListOptions{})
	assert.NoError(t, err)

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

	assert.Equal(t, want, taskList)
}

func TestTaskService_Update(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/tasks/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		testBody(t, r, "task/update/body_1.json")
		testWriteResponse(t, w, "task/update/response_1.json")
	})

	task, _, err := service.Task.Update(context.Background(), 1, &TaskUpdateRequest{
		Name:              String("Task update"),
		BillableByDefault: Bool(false),
		DefaultHourlyRate: Float64(213),
		IsDefault:         Bool(false),
		IsActive:          Bool(false),
	})
	assert.NoError(t, err)

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

	assert.Equal(t, want, task)
}
