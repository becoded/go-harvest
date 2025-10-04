package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestUserService_ListProjectAssignments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		userID    int64
		setupMock func(mux *http.ServeMux)
		want      *harvest.UserProjectAssignmentList
		wantErr   bool
	}{
		{
			name:   "Valid Project Assignment List",
			userID: 1782959,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/users/1782959/project_assignments", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "user/list_project_assignments/body_1.json")
					testWriteResponse(t, w, "user/list_project_assignments/response_1.json")
				})
			},
			want: &harvest.UserProjectAssignmentList{
				ProjectAssignments: []*harvest.UserProjectAssignment{
					{
						ID:               harvest.Int64(125068554),
						IsProjectManager: harvest.Bool(true),
						IsActive:         harvest.Bool(true),
						UseDefaultRates:  harvest.Bool(true),
						Budget:           nil,
						CreatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
						UpdatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
						HourlyRate:       harvest.Float64(100.0),
						Project: &harvest.Project{
							ID:   harvest.Int64(14308069),
							Name: harvest.String("Online Store - Phase 1"),
							Code: harvest.String("OS1"),
						},
						Client: &harvest.Client{
							ID:   harvest.Int64(5735776),
							Name: harvest.String("123 Industries"),
						},
						TaskAssignments: &[]harvest.ProjectTaskAssignment{
							{
								ID:         harvest.Int64(155505013),
								Billable:   harvest.Bool(true),
								IsActive:   harvest.Bool(true),
								CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
								UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
								HourlyRate: harvest.Float64(100.0),
								Budget:     nil,
								Task: &harvest.Task{
									ID:   harvest.Int64(8083365),
									Name: harvest.String("Graphic Design"),
								},
							},
							{
								ID:         harvest.Int64(155505014),
								Billable:   harvest.Bool(true),
								IsActive:   harvest.Bool(true),
								CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
								UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
								HourlyRate: harvest.Float64(100.0),
								Budget:     nil,
								Task: &harvest.Task{
									ID:   harvest.Int64(8083366),
									Name: harvest.String("Programming"),
								},
							},
							{
								ID:         harvest.Int64(155505015),
								Billable:   harvest.Bool(true),
								IsActive:   harvest.Bool(true),
								CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
								UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
								HourlyRate: harvest.Float64(100.0),
								Budget:     nil,
								Task: &harvest.Task{
									ID:   harvest.Int64(8083368),
									Name: harvest.String("Project Management"),
								},
							},
							{
								ID:         harvest.Int64(155505016),
								Billable:   harvest.Bool(false),
								IsActive:   harvest.Bool(true),
								CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
								UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 54, 6, 0, time.UTC)),
								HourlyRate: harvest.Float64(100.0),
								Budget:     nil,
								Task: &harvest.Task{
									ID:   harvest.Int64(8083369),
									Name: harvest.String("Research"),
								},
							},
						},
					},
					{
						ID:               harvest.Int64(125068553),
						IsProjectManager: harvest.Bool(true),
						IsActive:         harvest.Bool(true),
						UseDefaultRates:  harvest.Bool(false),
						Budget:           nil,
						CreatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
						UpdatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
						HourlyRate:       harvest.Float64(100.0),
						Project: &harvest.Project{
							ID:   harvest.Int64(14307913),
							Name: harvest.String("Marketing Website"),
							Code: harvest.String("MW"),
						},
						Client: &harvest.Client{
							ID:   harvest.Int64(5735774),
							Name: harvest.String("ABC Corp"),
						},
						TaskAssignments: &[]harvest.ProjectTaskAssignment{
							{
								ID:         harvest.Int64(155502709),
								Billable:   harvest.Bool(true),
								IsActive:   harvest.Bool(true),
								CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 36, 23, 0, time.UTC)),
								UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 36, 23, 0, time.UTC)),
								HourlyRate: harvest.Float64(100.0),
								Budget:     nil,
								Task: &harvest.Task{
									ID:   harvest.Int64(8083365),
									Name: harvest.String("Graphic Design"),
								},
							},
							{
								ID:         harvest.Int64(155502710),
								Billable:   harvest.Bool(true),
								IsActive:   harvest.Bool(true),
								CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 36, 23, 0, time.UTC)),
								UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 36, 23, 0, time.UTC)),
								HourlyRate: harvest.Float64(100.0),
								Budget:     nil,
								Task: &harvest.Task{
									ID:   harvest.Int64(8083366),
									Name: harvest.String("Programming"),
								},
							},
							{
								ID:         harvest.Int64(155502711),
								Billable:   harvest.Bool(true),
								IsActive:   harvest.Bool(true),
								CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 36, 23, 0, time.UTC)),
								UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 36, 23, 0, time.UTC)),
								HourlyRate: harvest.Float64(100.0),
								Budget:     nil,
								Task: &harvest.Task{
									ID:   harvest.Int64(8083368),
									Name: harvest.String("Project Management"),
								},
							},
							{
								ID:         harvest.Int64(155505153),
								Billable:   harvest.Bool(false),
								IsActive:   harvest.Bool(true),
								CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 53, 20, 0, time.UTC)),
								UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 54, 31, 0, time.UTC)),
								HourlyRate: harvest.Float64(100.0),
								Budget:     nil,
								Task: &harvest.Task{
									ID:   harvest.Int64(8083369),
									Name: harvest.String("Research"),
								},
							},
						},
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(2000),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					NextPage:     nil,
					PreviousPage: nil,
					Page:         harvest.Int(1),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/users/1782959/project_assignments?page=1&per_page=2000"),
						Next:     nil,
						Previous: nil,
						Last:     harvest.String("https://api.harvestapp.com/v2/users/1782959/project_assignments?page=1&per_page=2000"),
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "Error Fetching Project Assignment List",
			userID: 1782959,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/users/1782959/project_assignments", func(w http.ResponseWriter, r *http.Request) {
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

			userProjectAssignments, _, err := service.User.ListProjectAssignments(context.Background(), tt.userID, nil)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, userProjectAssignments)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, userProjectAssignments)
			}
		})
	}
}

func TestUserService_GetMyProjectAssignments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(mux *http.ServeMux)
		want      *harvest.UserProjectAssignmentList
		wantErr   bool
	}{
		{
			name: "Valid My Project Assignment List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/users/me/project_assignments", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "user/current_list_project_assignments/body_1.json")
					testWriteResponse(t, w, "user/current_list_project_assignments/response_1.json")
				})
			},
			want: &harvest.UserProjectAssignmentList{
				ProjectAssignments: []*harvest.UserProjectAssignment{
					{
						ID:               harvest.Int64(125066109),
						IsProjectManager: harvest.Bool(true),
						IsActive:         harvest.Bool(true),
						UseDefaultRates:  harvest.Bool(true),
						Budget:           nil,
						CreatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
						UpdatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
						HourlyRate:       harvest.Float64(100.0),
						Project: &harvest.Project{
							ID:   harvest.Int64(14308069),
							Name: harvest.String("Online Store - Phase 1"),
							Code: harvest.String("OS1"),
						},
						Client: &harvest.Client{
							ID:   harvest.Int64(5735776),
							Name: harvest.String("123 Industries"),
						},
						TaskAssignments: &[]harvest.ProjectTaskAssignment{
							{
								ID:         harvest.Int64(155505013),
								Billable:   harvest.Bool(true),
								IsActive:   harvest.Bool(true),
								CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
								UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
								HourlyRate: harvest.Float64(100.0),
								Budget:     nil,
								Task: &harvest.Task{
									ID:   harvest.Int64(8083365),
									Name: harvest.String("Graphic Design"),
								},
							},
							{
								ID:         harvest.Int64(155505014),
								Billable:   harvest.Bool(true),
								IsActive:   harvest.Bool(true),
								CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
								UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
								HourlyRate: harvest.Float64(100.0),
								Budget:     nil,
								Task: &harvest.Task{
									ID:   harvest.Int64(8083366),
									Name: harvest.String("Programming"),
								},
							},
							{
								ID:         harvest.Int64(155505015),
								Billable:   harvest.Bool(true),
								IsActive:   harvest.Bool(true),
								CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
								UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
								HourlyRate: harvest.Float64(100.0),
								Budget:     nil,
								Task: &harvest.Task{
									ID:   harvest.Int64(8083368),
									Name: harvest.String("Project Management"),
								},
							},
							{
								ID:         harvest.Int64(155505016),
								Billable:   harvest.Bool(false),
								IsActive:   harvest.Bool(true),
								CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
								UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 54, 6, 0, time.UTC)),
								HourlyRate: harvest.Float64(100.0),
								Budget:     nil,
								Task: &harvest.Task{
									ID:   harvest.Int64(8083369),
									Name: harvest.String("Research"),
								},
							},
						},
					},
					{
						ID:               harvest.Int64(125063975),
						IsProjectManager: harvest.Bool(true),
						IsActive:         harvest.Bool(true),
						UseDefaultRates:  harvest.Bool(false),
						Budget:           nil,
						CreatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 36, 23, 0, time.UTC)),
						UpdatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 36, 23, 0, time.UTC)),
						HourlyRate:       harvest.Float64(100.0),
						Project: &harvest.Project{
							ID:   harvest.Int64(14307913),
							Name: harvest.String("Marketing Website"),
							Code: harvest.String("MW"),
						},
						Client: &harvest.Client{
							ID:   harvest.Int64(5735774),
							Name: harvest.String("ABC Corp"),
						},
						TaskAssignments: &[]harvest.ProjectTaskAssignment{
							{
								ID:         harvest.Int64(155502709),
								Billable:   harvest.Bool(true),
								IsActive:   harvest.Bool(true),
								CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 36, 23, 0, time.UTC)),
								UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 36, 23, 0, time.UTC)),
								HourlyRate: harvest.Float64(100.0),
								Budget:     nil,
								Task: &harvest.Task{
									ID:   harvest.Int64(8083365),
									Name: harvest.String("Graphic Design"),
								},
							},
							{
								ID:         harvest.Int64(155502710),
								Billable:   harvest.Bool(true),
								IsActive:   harvest.Bool(true),
								CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 36, 23, 0, time.UTC)),
								UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 36, 23, 0, time.UTC)),
								HourlyRate: harvest.Float64(100.0),
								Budget:     nil,
								Task: &harvest.Task{
									ID:   harvest.Int64(8083366),
									Name: harvest.String("Programming"),
								},
							},
							{
								ID:         harvest.Int64(155502711),
								Billable:   harvest.Bool(true),
								IsActive:   harvest.Bool(true),
								CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 36, 23, 0, time.UTC)),
								UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 36, 23, 0, time.UTC)),
								HourlyRate: harvest.Float64(100.0),
								Budget:     nil,
								Task: &harvest.Task{
									ID:   harvest.Int64(8083368),
									Name: harvest.String("Project Management"),
								},
							},
							{
								ID:         harvest.Int64(155505153),
								Billable:   harvest.Bool(false),
								IsActive:   harvest.Bool(true),
								CreatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 53, 20, 0, time.UTC)),
								UpdatedAt:  harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 54, 31, 0, time.UTC)),
								HourlyRate: harvest.Float64(100.0),
								Budget:     nil,
								Task: &harvest.Task{
									ID:   harvest.Int64(8083369),
									Name: harvest.String("Research"),
								},
							},
						},
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(2000),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					NextPage:     nil,
					PreviousPage: nil,
					Page:         harvest.Int(1),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/users/1782884/project_assignments?page=1&per_page=2000"),
						Next:     nil,
						Previous: nil,
						Last:     harvest.String("https://api.harvestapp.com/v2/users/1782884/project_assignments?page=1&per_page=2000"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Error Fetching My Project Assignment List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/users/me/project_assignments", func(w http.ResponseWriter, r *http.Request) {
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

			userProjectAssignments, _, err := service.User.GetMyProjectAssignments(context.Background(), nil)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, userProjectAssignments)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, userProjectAssignments)
			}
		})
	}
}

func TestUserProjectAssignment_String(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)
	updatedAt := time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)

	tests := []struct {
		name string
		in   harvest.UserProjectAssignment
		want string
	}{
		{
			name: "Assignment with all fields",
			in: harvest.UserProjectAssignment{
				ID:               harvest.Int64(125068554),
				IsProjectManager: harvest.Bool(true),
				IsActive:         harvest.Bool(true),
				UseDefaultRates:  harvest.Bool(true),
				Budget:           harvest.Float64(1000.0),
				HourlyRate:       harvest.Float64(100.0),
				CreatedAt:        &createdAt,
				UpdatedAt:        &updatedAt,
				Project: &harvest.Project{
					ID:   harvest.Int64(14308069),
					Name: harvest.String("Online Store - Phase 1"),
					Code: harvest.String("OS1"),
				},
				Client: &harvest.Client{
					ID:   harvest.Int64(5735776),
					Name: harvest.String("123 Industries"),
				},
				TaskAssignments: &[]harvest.ProjectTaskAssignment{
					{
						ID:         harvest.Int64(155505013),
						Billable:   harvest.Bool(true),
						IsActive:   harvest.Bool(true),
						HourlyRate: harvest.Float64(100.0),
						Task: &harvest.Task{
							ID:   harvest.Int64(8083365),
							Name: harvest.String("Graphic Design"),
						},
					},
				},
			},
			want: `harvest.UserProjectAssignment{ID:125068554, IsActive:true, IsProjectManager:true, UseDefaultRates:true, HourlyRate:100, Budget:1000, CreatedAt:time.Time{2017-06-26 22:32:52 +0000 UTC}, UpdatedAt:time.Time{2017-06-26 22:32:52 +0000 UTC}, Project:harvest.Project{ID:14308069, Name:"Online Store - Phase 1", Code:"OS1"}, Client:harvest.Client{ID:5735776, Name:"123 Industries"}, TaskAssignments:[harvest.ProjectTaskAssignment{ID:155505013, Task:harvest.Task{ID:8083365, Name:"Graphic Design"}, IsActive:true, Billable:true, HourlyRate:100}]}`, //nolint: lll
		},
		{
			name: "Assignment with minimal fields",
			in: harvest.UserProjectAssignment{
				ID:         harvest.Int64(125068554),
				IsActive:   harvest.Bool(true),
				HourlyRate: harvest.Float64(100.0),
			},
			want: `harvest.UserProjectAssignment{ID:125068554, IsActive:true, HourlyRate:100}`,
		},
		{
			name: "Assignment with project and client only",
			in: harvest.UserProjectAssignment{
				ID: harvest.Int64(125068554),
				Project: &harvest.Project{
					ID:   harvest.Int64(14308069),
					Name: harvest.String("Online Store - Phase 1"),
				},
				Client: &harvest.Client{
					ID:   harvest.Int64(5735776),
					Name: harvest.String("123 Industries"),
				},
			},
			want: `harvest.UserProjectAssignment{ID:125068554, Project:harvest.Project{ID:14308069, Name:"Online Store - Phase 1"}, Client:harvest.Client{ID:5735776, Name:"123 Industries"}}`, //nolint: lll
		},
		{
			name: "Assignment with task assignments",
			in: harvest.UserProjectAssignment{
				ID:       harvest.Int64(125068554),
				IsActive: harvest.Bool(true),
				TaskAssignments: &[]harvest.ProjectTaskAssignment{
					{
						ID:       harvest.Int64(155505013),
						Billable: harvest.Bool(true),
						Task: &harvest.Task{
							ID:   harvest.Int64(8083365),
							Name: harvest.String("Graphic Design"),
						},
					},
					{
						ID:       harvest.Int64(155505014),
						Billable: harvest.Bool(true),
						Task: &harvest.Task{
							ID:   harvest.Int64(8083366),
							Name: harvest.String("Programming"),
						},
					},
				},
			},
			want: `harvest.UserProjectAssignment{ID:125068554, IsActive:true, TaskAssignments:[harvest.ProjectTaskAssignment{ID:155505013, Task:harvest.Task{ID:8083365, Name:"Graphic Design"}, Billable:true} harvest.ProjectTaskAssignment{ID:155505014, Task:harvest.Task{ID:8083366, Name:"Programming"}, Billable:true}]}`, //nolint: lll
		},
		{
			name: "Empty Assignment",
			in:   harvest.UserProjectAssignment{},
			want: `harvest.UserProjectAssignment{}`,
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

func TestUserProjectAssignmentList_String(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)
	updatedAt := time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)

	tests := []struct {
		name string
		in   harvest.UserProjectAssignmentList
		want string
	}{
		{
			name: "AssignmentList with multiple assignments",
			in: harvest.UserProjectAssignmentList{
				ProjectAssignments: []*harvest.UserProjectAssignment{
					{
						ID:         harvest.Int64(125068554),
						IsActive:   harvest.Bool(true),
						HourlyRate: harvest.Float64(100.0),
						Project: &harvest.Project{
							ID:   harvest.Int64(14308069),
							Name: harvest.String("Online Store - Phase 1"),
						},
					},
					{
						ID:         harvest.Int64(125068553),
						IsActive:   harvest.Bool(true),
						HourlyRate: harvest.Float64(100.0),
						Project: &harvest.Project{
							ID:   harvest.Int64(14307913),
							Name: harvest.String("Marketing Website"),
						},
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.UserProjectAssignmentList{ProjectAssignments:[harvest.UserProjectAssignment{ID:125068554, IsActive:true, HourlyRate:100, Project:harvest.Project{ID:14308069, Name:"Online Store - Phase 1"}} harvest.UserProjectAssignment{ID:125068553, IsActive:true, HourlyRate:100, Project:harvest.Project{ID:14307913, Name:"Marketing Website"}}], Pagination:harvest.Pagination{PerPage:100, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
		},
		{
			name: "AssignmentList with single assignment",
			in: harvest.UserProjectAssignmentList{
				ProjectAssignments: []*harvest.UserProjectAssignment{
					{
						ID:               harvest.Int64(999),
						IsProjectManager: harvest.Bool(true),
						IsActive:         harvest.Bool(true),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(50),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.UserProjectAssignmentList{ProjectAssignments:[harvest.UserProjectAssignment{ID:999, IsActive:true, IsProjectManager:true}], Pagination:harvest.Pagination{PerPage:50, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
		},
		{
			name: "Empty AssignmentList",
			in: harvest.UserProjectAssignmentList{
				ProjectAssignments: []*harvest.UserProjectAssignment{},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(0),
					TotalEntries: harvest.Int(0),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.UserProjectAssignmentList{ProjectAssignments:[], Pagination:harvest.Pagination{PerPage:100, TotalPages:0, TotalEntries:0, Page:1}}`, //nolint: lll
		},
		{
			name: "AssignmentList with Links",
			in: harvest.UserProjectAssignmentList{
				ProjectAssignments: []*harvest.UserProjectAssignment{
					{
						ID:       harvest.Int64(100),
						IsActive: harvest.Bool(true),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(25),
					TotalPages:   harvest.Int(3),
					TotalEntries: harvest.Int(75),
					Page:         harvest.Int(2),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/users/1782959/project_assignments?page=1&per_page=25"),
						Next:     harvest.String("https://api.harvestapp.com/v2/users/1782959/project_assignments?page=3&per_page=25"),
						Previous: harvest.String("https://api.harvestapp.com/v2/users/1782959/project_assignments?page=1&per_page=25"),
						Last:     harvest.String("https://api.harvestapp.com/v2/users/1782959/project_assignments?page=3&per_page=25"),
					},
				},
			},
			want: `harvest.UserProjectAssignmentList{ProjectAssignments:[harvest.UserProjectAssignment{ID:100, IsActive:true}], Pagination:harvest.Pagination{PerPage:25, TotalPages:3, TotalEntries:75, Page:2, Links:harvest.PageLinks{First:"https://api.harvestapp.com/v2/users/1782959/project_assignments?page=1&per_page=25", Next:"https://api.harvestapp.com/v2/users/1782959/project_assignments?page=3&per_page=25", Previous:"https://api.harvestapp.com/v2/users/1782959/project_assignments?page=1&per_page=25", Last:"https://api.harvestapp.com/v2/users/1782959/project_assignments?page=3&per_page=25"}}}`, //nolint: lll
		},
		{
			name: "AssignmentList with full details",
			in: harvest.UserProjectAssignmentList{
				ProjectAssignments: []*harvest.UserProjectAssignment{
					{
						ID:               harvest.Int64(1),
						IsProjectManager: harvest.Bool(true),
						IsActive:         harvest.Bool(true),
						UseDefaultRates:  harvest.Bool(false),
						HourlyRate:       harvest.Float64(150.0),
						Budget:           harvest.Float64(5000.0),
						CreatedAt:        &createdAt,
						UpdatedAt:        &updatedAt,
						Project: &harvest.Project{
							ID:   harvest.Int64(10),
							Name: harvest.String("Project Alpha"),
							Code: harvest.String("PA"),
						},
						Client: &harvest.Client{
							ID:   harvest.Int64(20),
							Name: harvest.String("Client A"),
						},
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(10),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.UserProjectAssignmentList{ProjectAssignments:[harvest.UserProjectAssignment{ID:1, IsActive:true, IsProjectManager:true, UseDefaultRates:false, HourlyRate:150, Budget:5000, CreatedAt:time.Time{2017-06-26 22:32:52 +0000 UTC}, UpdatedAt:time.Time{2017-06-26 22:32:52 +0000 UTC}, Project:harvest.Project{ID:10, Name:"Project Alpha", Code:"PA"}, Client:harvest.Client{ID:20, Name:"Client A"}}], Pagination:harvest.Pagination{PerPage:10, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
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
