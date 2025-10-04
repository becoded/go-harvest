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
		{
			name:      "Error Creating Task Assignment",
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
		{
			name:      "Error Fetching Task Assignment List",
			projectID: 14308069,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/projects/14308069/task_assignments", func(w http.ResponseWriter, r *http.Request) {
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
		{
			name:             "Error Fetching Task Assignment",
			projectID:        14308069,
			taskAssignmentID: 155505016,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/projects/14308069/task_assignments/155505016", func(w http.ResponseWriter, r *http.Request) {
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

func TestProjectTaskAssignment_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.ProjectTaskAssignment
		want string
	}{
		{
			name: "ProjectTaskAssignment with all fields",
			in: harvest.ProjectTaskAssignment{
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
				Budget:     harvest.Float64(5000.0),
				CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
				UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 54, 6, 0, time.UTC)),
			},
			want: `harvest.ProjectTaskAssignment{ID:155505016, Project:harvest.Project{ID:14308069, Name:"Online Store - Phase 1", Code:"OS1"}, Task:harvest.Task{ID:8083369, Name:"Research"}, IsActive:true, Billable:false, HourlyRate:100, Budget:5000, CreatedAt:time.Time{2017-06-26 21:52:18 +0000 UTC}, UpdatedAt:time.Time{2017-06-26 21:54:06 +0000 UTC}}`, //nolint: lll
		},
		{
			name: "ProjectTaskAssignment with minimal fields",
			in: harvest.ProjectTaskAssignment{
				ID: harvest.Int64(999),
				Task: &harvest.Task{
					ID: harvest.Int64(123),
				},
			},
			want: `harvest.ProjectTaskAssignment{ID:999, Task:harvest.Task{ID:123}}`,
		},
		{
			name: "ProjectTaskAssignment without budget",
			in: harvest.ProjectTaskAssignment{
				ID: harvest.Int64(155505015),
				Task: &harvest.Task{
					ID:   harvest.Int64(8083368),
					Name: harvest.String("Project Management"),
				},
				Project: &harvest.Project{
					ID:   harvest.Int64(14308069),
					Name: harvest.String("Online Store - Phase 1"),
				},
				IsActive:   harvest.Bool(true),
				Billable:   harvest.Bool(true),
				HourlyRate: harvest.Float64(100.0),
			},
			want: `harvest.ProjectTaskAssignment{ID:155505015, Project:harvest.Project{ID:14308069, Name:"Online Store - Phase 1"}, Task:harvest.Task{ID:8083368, Name:"Project Management"}, IsActive:true, Billable:true, HourlyRate:100}`, //nolint: lll
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

func TestProjectTaskAssignmentList_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.ProjectTaskAssignmentList
		want string
	}{
		{
			name: "ProjectTaskAssignmentList with multiple assignments",
			in: harvest.ProjectTaskAssignmentList{
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
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(2000),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.ProjectTaskAssignmentList{TaskAssignments:[harvest.ProjectTaskAssignment{ID:155505016, Project:harvest.Project{ID:14308069, Name:"Online Store - Phase 1", Code:"OS1"}, Task:harvest.Task{ID:8083369, Name:"Research"}, IsActive:true, Billable:false, HourlyRate:100} harvest.ProjectTaskAssignment{ID:155505015, Project:harvest.Project{ID:14308069, Name:"Online Store - Phase 1", Code:"OS1"}, Task:harvest.Task{ID:8083368, Name:"Project Management"}, IsActive:true, Billable:true, HourlyRate:100}], Pagination:harvest.Pagination{PerPage:2000, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
		},
		{
			name: "ProjectTaskAssignmentList with single assignment",
			in: harvest.ProjectTaskAssignmentList{
				TaskAssignments: []*harvest.ProjectTaskAssignment{
					{
						ID: harvest.Int64(999),
						Task: &harvest.Task{
							ID:   harvest.Int64(123),
							Name: harvest.String("Development"),
						},
						IsActive: harvest.Bool(true),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.ProjectTaskAssignmentList{TaskAssignments:[harvest.ProjectTaskAssignment{ID:999, Task:harvest.Task{ID:123, Name:"Development"}, IsActive:true}], Pagination:harvest.Pagination{PerPage:100, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
		},
		{
			name: "Empty ProjectTaskAssignmentList",
			in: harvest.ProjectTaskAssignmentList{
				TaskAssignments: []*harvest.ProjectTaskAssignment{},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(2000),
					TotalPages:   harvest.Int(0),
					TotalEntries: harvest.Int(0),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.ProjectTaskAssignmentList{TaskAssignments:[], Pagination:harvest.Pagination{PerPage:2000, TotalPages:0, TotalEntries:0, Page:1}}`, //nolint: lll
		},
		{
			name: "ProjectTaskAssignmentList with Links",
			in: harvest.ProjectTaskAssignmentList{
				TaskAssignments: []*harvest.ProjectTaskAssignment{
					{
						ID: harvest.Int64(100),
						Task: &harvest.Task{
							ID: harvest.Int64(50),
						},
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(50),
					TotalPages:   harvest.Int(3),
					TotalEntries: harvest.Int(150),
					Page:         harvest.Int(2),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/projects/14308069/task_assignments?page=1&per_page=50"),
						Next:     harvest.String("https://api.harvestapp.com/v2/projects/14308069/task_assignments?page=3&per_page=50"),
						Previous: harvest.String("https://api.harvestapp.com/v2/projects/14308069/task_assignments?page=1&per_page=50"),
						Last:     harvest.String("https://api.harvestapp.com/v2/projects/14308069/task_assignments?page=3&per_page=50"),
					},
				},
			},
			want: `harvest.ProjectTaskAssignmentList{TaskAssignments:[harvest.ProjectTaskAssignment{ID:100, Task:harvest.Task{ID:50}}], Pagination:harvest.Pagination{PerPage:50, TotalPages:3, TotalEntries:150, Page:2, Links:harvest.PageLinks{First:"https://api.harvestapp.com/v2/projects/14308069/task_assignments?page=1&per_page=50", Next:"https://api.harvestapp.com/v2/projects/14308069/task_assignments?page=3&per_page=50", Previous:"https://api.harvestapp.com/v2/projects/14308069/task_assignments?page=1&per_page=50", Last:"https://api.harvestapp.com/v2/projects/14308069/task_assignments?page=3&per_page=50"}}}`, //nolint: lll
		},
		{
			name: "ProjectTaskAssignmentList with Budget",
			in: harvest.ProjectTaskAssignmentList{
				TaskAssignments: []*harvest.ProjectTaskAssignment{
					{
						ID: harvest.Int64(1),
						Task: &harvest.Task{
							ID:   harvest.Int64(10),
							Name: harvest.String("Design"),
						},
						Budget:     harvest.Float64(2500.0),
						HourlyRate: harvest.Float64(75.0),
					},
					{
						ID: harvest.Int64(2),
						Task: &harvest.Task{
							ID:   harvest.Int64(20),
							Name: harvest.String("Testing"),
						},
						Budget:     harvest.Float64(1500.0),
						HourlyRate: harvest.Float64(50.0),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(10),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.ProjectTaskAssignmentList{TaskAssignments:[harvest.ProjectTaskAssignment{ID:1, Task:harvest.Task{ID:10, Name:"Design"}, HourlyRate:75, Budget:2500} harvest.ProjectTaskAssignment{ID:2, Task:harvest.Task{ID:20, Name:"Testing"}, HourlyRate:50, Budget:1500}], Pagination:harvest.Pagination{PerPage:10, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
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
