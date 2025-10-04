package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestTaskService_Create(t *testing.T) {
	t.Parallel()

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	tests := []struct {
		name      string
		args      *harvest.TaskCreateRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.Task
		wantErr   bool
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
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testFormValues(t, r, values{})
					testBody(t, r, "task/create/body_1.json")
					testWriteResponse(t, w, "task/create/response_1.json")
				})
			},
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
			wantErr: false,
		},
		{
			name: "error creating task",
			args: &harvest.TaskCreateRequest{
				Name:              harvest.String("Task new"),
				BillableByDefault: harvest.Bool(true),
				DefaultHourlyRate: harvest.Float64(123),
				IsDefault:         harvest.Bool(true),
				IsActive:          harvest.Bool(true),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
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

			task, _, err := service.Task.Create(context.Background(), tt.args)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, task)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, task)
			}
		})
	}
}

func TestTaskService_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		taskID    int64
		setupMock func(mux *http.ServeMux)
		wantErr   bool
	}{
		{
			name:   "delete 1",
			taskID: 1,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/tasks/1", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					testFormValues(t, r, values{})
					testBody(t, r, "task/delete/body_1.json")
					testWriteResponse(t, w, "task/delete/response_1.json")
				})
			},
			wantErr: false,
		},
		{
			name:   "error deleting task",
			taskID: 1,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/tasks/1", func(w http.ResponseWriter, r *http.Request) {
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

			_, err := service.Task.Delete(context.Background(), tt.taskID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskService_Get(t *testing.T) {
	t.Parallel()

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	tests := []struct {
		name      string
		taskID    int64
		setupMock func(mux *http.ServeMux)
		want      *harvest.Task
		wantErr   bool
	}{
		{
			name:   "get 1",
			taskID: 1,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/tasks/1", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "task/get/body_1.json")
					testWriteResponse(t, w, "task/get/response_1.json")
				})
			},
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
			wantErr: false,
		},
		{
			name:   "error getting task",
			taskID: 1,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/tasks/1", func(w http.ResponseWriter, r *http.Request) {
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

			task, _, err := service.Task.Get(context.Background(), tt.taskID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, task)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, task)
			}
		})
	}
}

func TestTaskService_List(t *testing.T) {
	t.Parallel()

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)
	createdTwo := time.Date(
		2018, 3, 2, 10, 12, 13, 0, time.UTC)
	updatedTwo := time.Date(
		2018, 4, 30, 12, 13, 14, 0, time.UTC)

	tests := []struct {
		name      string
		setupMock func(mux *http.ServeMux)
		want      *harvest.TaskList
		wantErr   bool
	}{
		{
			name: "list tasks",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "task/list/body_1.json")
					testWriteResponse(t, w, "task/list/response_1.json")
				})
			},
			want: &harvest.TaskList{
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
			},
			wantErr: false,
		},
		{
			name: "error listing tasks",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
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

			taskList, _, err := service.Task.List(context.Background(), &harvest.TaskListOptions{})

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, taskList)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, taskList)
			}
		})
	}
}

func TestTaskService_Update(t *testing.T) {
	t.Parallel()

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	tests := []struct {
		name      string
		taskID    int64
		args      *harvest.TaskUpdateRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.Task
		wantErr   bool
	}{
		{
			name:   "update task",
			taskID: 1,
			args: &harvest.TaskUpdateRequest{
				Name:              harvest.String("Task update"),
				BillableByDefault: harvest.Bool(false),
				DefaultHourlyRate: harvest.Float64(213),
				IsDefault:         harvest.Bool(false),
				IsActive:          harvest.Bool(false),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/tasks/1", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					testFormValues(t, r, values{})
					testBody(t, r, "task/update/body_1.json")
					testWriteResponse(t, w, "task/update/response_1.json")
				})
			},
			want: &harvest.Task{
				ID:                harvest.Int64(1),
				Name:              harvest.String("Task update"),
				BillableByDefault: harvest.Bool(false),
				DefaultHourlyRate: harvest.Float64(213),
				IsDefault:         harvest.Bool(false),
				IsActive:          harvest.Bool(false),
				CreatedAt:         &createdOne,
				UpdatedAt:         &updatedOne,
			},
			wantErr: false,
		},
		{
			name:   "error updating task",
			taskID: 1,
			args: &harvest.TaskUpdateRequest{
				Name:              harvest.String("Task update"),
				BillableByDefault: harvest.Bool(false),
				DefaultHourlyRate: harvest.Float64(213),
				IsDefault:         harvest.Bool(false),
				IsActive:          harvest.Bool(false),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/tasks/1", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
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

			task, _, err := service.Task.Update(context.Background(), tt.taskID, tt.args)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, task)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, task)
			}
		})
	}
}

func TestTask_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.Task
		want string
	}{
		{
			name: "Task with all fields",
			in: harvest.Task{
				ID:                harvest.Int64(1),
				Name:              harvest.String("Programming"),
				BillableByDefault: harvest.Bool(true),
				DefaultHourlyRate: harvest.Float64(100.0),
				IsDefault:         harvest.Bool(true),
				IsActive:          harvest.Bool(true),
				CreatedAt:         harvest.TimeTimeP(time.Date(2018, 1, 31, 20, 34, 30, 0, time.UTC)),
				UpdatedAt:         harvest.TimeTimeP(time.Date(2018, 5, 31, 21, 34, 30, 0, time.UTC)),
			},
			want: `harvest.Task{ID:1, Name:"Programming", BillableByDefault:true, DefaultHourlyRate:100, IsDefault:true, IsActive:true, CreatedAt:time.Time{2018-01-31 20:34:30 +0000 UTC}, UpdatedAt:time.Time{2018-05-31 21:34:30 +0000 UTC}}`, //nolint: lll
		},
		{
			name: "Task with minimal fields",
			in: harvest.Task{
				ID:   harvest.Int64(999),
				Name: harvest.String("Design"),
			},
			want: `harvest.Task{ID:999, Name:"Design"}`,
		},
		{
			name: "Task with false boolean values",
			in: harvest.Task{
				ID:                harvest.Int64(2),
				Name:              harvest.String("Research"),
				BillableByDefault: harvest.Bool(false),
				DefaultHourlyRate: harvest.Float64(75.5),
				IsDefault:         harvest.Bool(false),
				IsActive:          harvest.Bool(false),
			},
			want: `harvest.Task{ID:2, Name:"Research", BillableByDefault:false, DefaultHourlyRate:75.5, IsDefault:false, IsActive:false}`, //nolint: lll
		},
		{
			name: "Task without timestamps",
			in: harvest.Task{
				ID:                harvest.Int64(3),
				Name:              harvest.String("Testing"),
				BillableByDefault: harvest.Bool(true),
				IsActive:          harvest.Bool(true),
			},
			want: `harvest.Task{ID:3, Name:"Testing", BillableByDefault:true, IsActive:true}`,
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

func TestTaskList_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.TaskList
		want string
	}{
		{
			name: "TaskList with multiple tasks",
			in: harvest.TaskList{
				Tasks: []*harvest.Task{
					{
						ID:                harvest.Int64(1),
						Name:              harvest.String("Programming"),
						BillableByDefault: harvest.Bool(true),
						DefaultHourlyRate: harvest.Float64(100.0),
						IsDefault:         harvest.Bool(true),
						IsActive:          harvest.Bool(true),
					},
					{
						ID:                harvest.Int64(2),
						Name:              harvest.String("Design"),
						BillableByDefault: harvest.Bool(false),
						DefaultHourlyRate: harvest.Float64(75.0),
						IsDefault:         harvest.Bool(false),
						IsActive:          harvest.Bool(true),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.TaskList{Tasks:[harvest.Task{ID:1, Name:"Programming", BillableByDefault:true, DefaultHourlyRate:100, IsDefault:true, IsActive:true} harvest.Task{ID:2, Name:"Design", BillableByDefault:false, DefaultHourlyRate:75, IsDefault:false, IsActive:true}], Pagination:harvest.Pagination{PerPage:100, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
		},
		{
			name: "TaskList with single task",
			in: harvest.TaskList{
				Tasks: []*harvest.Task{
					{
						ID:       harvest.Int64(999),
						Name:     harvest.String("Testing"),
						IsActive: harvest.Bool(true),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(50),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.TaskList{Tasks:[harvest.Task{ID:999, Name:"Testing", IsActive:true}], Pagination:harvest.Pagination{PerPage:50, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
		},
		{
			name: "Empty TaskList",
			in: harvest.TaskList{
				Tasks: []*harvest.Task{},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(0),
					TotalEntries: harvest.Int(0),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.TaskList{Tasks:[], Pagination:harvest.Pagination{PerPage:100, TotalPages:0, TotalEntries:0, Page:1}}`, //nolint: lll
		},
		{
			name: "TaskList with Links",
			in: harvest.TaskList{
				Tasks: []*harvest.Task{
					{
						ID:   harvest.Int64(100),
						Name: harvest.String("Management"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(25),
					TotalPages:   harvest.Int(3),
					TotalEntries: harvest.Int(75),
					Page:         harvest.Int(2),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/tasks?page=1&per_page=25"),
						Next:     harvest.String("https://api.harvestapp.com/v2/tasks?page=3&per_page=25"),
						Previous: harvest.String("https://api.harvestapp.com/v2/tasks?page=1&per_page=25"),
						Last:     harvest.String("https://api.harvestapp.com/v2/tasks?page=3&per_page=25"),
					},
				},
			},
			want: `harvest.TaskList{Tasks:[harvest.Task{ID:100, Name:"Management"}], Pagination:harvest.Pagination{PerPage:25, TotalPages:3, TotalEntries:75, Page:2, Links:harvest.PageLinks{First:"https://api.harvestapp.com/v2/tasks?page=1&per_page=25", Next:"https://api.harvestapp.com/v2/tasks?page=3&per_page=25", Previous:"https://api.harvestapp.com/v2/tasks?page=1&per_page=25", Last:"https://api.harvestapp.com/v2/tasks?page=3&per_page=25"}}}`, //nolint: lll
		},
		{
			name: "TaskList with timestamps",
			in: harvest.TaskList{
				Tasks: []*harvest.Task{
					{
						ID:        harvest.Int64(1),
						Name:      harvest.String("Development"),
						IsActive:  harvest.Bool(true),
						CreatedAt: harvest.TimeTimeP(time.Date(2018, 1, 31, 20, 34, 30, 0, time.UTC)),
						UpdatedAt: harvest.TimeTimeP(time.Date(2018, 5, 31, 21, 34, 30, 0, time.UTC)),
					},
					{
						ID:        harvest.Int64(2),
						Name:      harvest.String("Support"),
						IsActive:  harvest.Bool(false),
						CreatedAt: harvest.TimeTimeP(time.Date(2018, 3, 2, 10, 12, 13, 0, time.UTC)),
						UpdatedAt: harvest.TimeTimeP(time.Date(2018, 4, 30, 12, 13, 14, 0, time.UTC)),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(10),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.TaskList{Tasks:[harvest.Task{ID:1, Name:"Development", IsActive:true, CreatedAt:time.Time{2018-01-31 20:34:30 +0000 UTC}, UpdatedAt:time.Time{2018-05-31 21:34:30 +0000 UTC}} harvest.Task{ID:2, Name:"Support", IsActive:false, CreatedAt:time.Time{2018-03-02 10:12:13 +0000 UTC}, UpdatedAt:time.Time{2018-04-30 12:13:14 +0000 UTC}}], Pagination:harvest.Pagination{PerPage:10, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
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
