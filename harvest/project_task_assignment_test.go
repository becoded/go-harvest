package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestProjectTaskAssignmentService_CreateTaskAssignment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		projectID int64
		request   *harvest.ProjectTaskAssignmentCreateRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.ProjectTaskAssignment
		wantErr   bool
	}{
		{
			name:      "Valid Task Assignment Creation",
			projectID: 14308069,
			request: &harvest.ProjectTaskAssignmentCreateRequest{
				TaskID:     harvest.Int64(8083800),
				IsActive:   harvest.Bool(true),
				Billable:   harvest.Bool(true),
				HourlyRate: harvest.Float64(75.5),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/projects/14308069/task_assignments", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testBody(t, r, "project_task_assignment/create/body_1.json")
					testWriteResponse(t, w, "project_task_assignment/create/response_1.json")
				})
			},
			want: &harvest.ProjectTaskAssignment{
				ID: harvest.Int64(155506339),
				Task: &harvest.Task{
					ID:   harvest.Int64(8083800),
					Name: harvest.String("Business Development"),
				},
				Project: &harvest.Project{
					ID:   harvest.Int64(14308069),
					Name: harvest.String("Online Store - Phase 1"),
					Code: harvest.String("OS1"),
				},
				IsActive:   harvest.Bool(true),
				Billable:   harvest.Bool(true),
				HourlyRate: harvest.Float64(75.5),
				CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 10, 43, 0, time.UTC)),
				UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 10, 43, 0, time.UTC)),
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

			got, _, err := service.Project.CreateTaskAssignment(context.Background(), tt.projectID, tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestProjectTaskAssignmentService_ListTaskAssignments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		projectID int64
		setupMock func(mux *http.ServeMux)
		want      *harvest.ProjectTaskAssignmentList
		wantErr   bool
	}{
		{
			name:      "Valid Task Assignment List",
			projectID: 14308069,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/projects/14308069/task_assignments", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testWriteResponse(t, w, "project_task_assignment/list/response_1.json")
				})
			},
			want: &harvest.ProjectTaskAssignmentList{
				TaskAssignments: []*harvest.ProjectTaskAssignment{
					{
						ID: harvest.Int64(155505016),
						Task: &harvest.Task{
							ID:   harvest.Int64(8083369),
							Name: harvest.String("Research"),
						},
						Project: &harvest.Project{
							ID:   harvest.Int64(14308069),
							Name: harvest.String("Online Store - Phase 1"),
							Code: harvest.String("OS1"),
						},
						IsActive:   harvest.Bool(true),
						Billable:   harvest.Bool(false),
						HourlyRate: harvest.Float64(100.0),
						CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
						UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 54, 6, 0, time.UTC)),
					},
					{
						ID: harvest.Int64(155505015),
						Task: &harvest.Task{
							ID:   harvest.Int64(8083368),
							Name: harvest.String("Project Management"),
						},
						Project: &harvest.Project{
							ID:   harvest.Int64(14308069),
							Name: harvest.String("Online Store - Phase 1"),
							Code: harvest.String("OS1"),
						},
						IsActive:   harvest.Bool(true),
						Billable:   harvest.Bool(true),
						HourlyRate: harvest.Float64(100.0),
						CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
						UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
					},
					{
						ID: harvest.Int64(155505014),
						Task: &harvest.Task{
							ID:   harvest.Int64(8083366),
							Name: harvest.String("Programming"),
						},
						Project: &harvest.Project{
							ID:   harvest.Int64(14308069),
							Name: harvest.String("Online Store - Phase 1"),
							Code: harvest.String("OS1"),
						},
						IsActive:   harvest.Bool(true),
						Billable:   harvest.Bool(true),
						HourlyRate: harvest.Float64(100.0),
						CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
						UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
					},
					{
						ID: harvest.Int64(155505013),
						Task: &harvest.Task{
							ID:   harvest.Int64(8083365),
							Name: harvest.String("Graphic Design"),
						},
						Project: &harvest.Project{
							ID:   harvest.Int64(14308069),
							Name: harvest.String("Online Store - Phase 1"),
							Code: harvest.String("OS1"),
						},
						IsActive:   harvest.Bool(true),
						Billable:   harvest.Bool(true),
						HourlyRate: harvest.Float64(100.0),
						CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
						UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(2000),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(4),
					NextPage:     nil,
					PreviousPage: nil,
					Page:         harvest.Int(1),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/projects/14308069/task_assignments?page=1&per_page=2000"),
						Next:     nil,
						Previous: nil,
						Last:     harvest.String("https://api.harvestapp.com/v2/projects/14308069/task_assignments?page=1&per_page=2000"),
					},
				},
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

			got, _, err := service.Project.ListTaskAssignments(
				context.Background(),
				tt.projectID,
				&harvest.ProjectTaskAssignmentListOptions{},
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

func TestProjectTaskAssignmentService_GetTaskAssignment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		projectID        int64
		taskAssignmentID int64
		setupMock        func(mux *http.ServeMux)
		want             *harvest.ProjectTaskAssignment
		wantErr          bool
	}{
		{
			name:             "Valid Task Assignment Retrieval",
			projectID:        14308069,
			taskAssignmentID: 155505016,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/projects/14308069/task_assignments/155505016", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testWriteResponse(t, w, "project_task_assignment/get/response_1.json")
				})
			},
			want: &harvest.ProjectTaskAssignment{
				ID: harvest.Int64(155505016),
				Task: &harvest.Task{
					ID:   harvest.Int64(8083369),
					Name: harvest.String("Research"),
				},
				Project: &harvest.Project{
					ID:   harvest.Int64(14308069),
					Name: harvest.String("Online Store - Phase 1"),
					Code: harvest.String("OS1"),
				},
				IsActive:   harvest.Bool(true),
				Billable:   harvest.Bool(false),
				HourlyRate: harvest.Float64(100.0),
				CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
				UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 54, 6, 0, time.UTC)),
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

			got, _, err := service.Project.GetTaskAssignment(
				context.Background(),
				tt.projectID,
				tt.taskAssignmentID,
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
