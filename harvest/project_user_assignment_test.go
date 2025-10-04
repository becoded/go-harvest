package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestProjectUserAssignmentService_ListUserAssignments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		projectID int64
		setupMock func(mux *http.ServeMux)
		want      *harvest.ProjectUserAssignmentList
		wantErr   bool
	}{
		{
			name:      "Valid User Assignment List",
			projectID: 14308069,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/projects/14308069/user_assignments", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testWriteResponse(t, w, "project_user_assignment/list/response_1.json")
				})
			},
			want: &harvest.ProjectUserAssignmentList{
				UserAssignments: []*harvest.ProjectUserAssignment{
					{
						ID:               harvest.Int64(125068554),
						IsProjectManager: harvest.Bool(true),
						IsActive:         harvest.Bool(true),
						HourlyRate:       harvest.Float64(100.0),
						CreatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
						UpdatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
						Project: &harvest.Project{
							ID:   harvest.Int64(14308069),
							Name: harvest.String("Online Store - Phase 1"),
							Code: harvest.String("OS1"),
						},
						User: &harvest.User{
							ID:   harvest.Int64(1782959),
							Name: harvest.String("Kim Allen"),
						},
					},
					{
						ID:               harvest.Int64(125066109),
						IsProjectManager: harvest.Bool(true),
						IsActive:         harvest.Bool(true),
						HourlyRate:       harvest.Float64(100.0),
						CreatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
						UpdatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
						Project: &harvest.Project{
							ID:   harvest.Int64(14308069),
							Name: harvest.String("Online Store - Phase 1"),
							Code: harvest.String("OS1"),
						},
						User: &harvest.User{
							ID:   harvest.Int64(1782884),
							Name: harvest.String("Jeremy Israelsen"),
						},
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(2000),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
					Links: &harvest.PageLinks{
						First: harvest.String("https://api.harvestapp.com/v2/projects/14308069/user_assignments?page=1&per_page=2000"),
						Last:  harvest.String("https://api.harvestapp.com/v2/projects/14308069/user_assignments?page=1&per_page=2000"),
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "Error Fetching User Assignment List",
			projectID: 14308069,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/projects/14308069/user_assignments", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Project.ListUserAssignments(
				context.Background(),
				tt.projectID,
				&harvest.ProjectUserAssignmentListOptions{},
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

func TestProjectUserAssignmentService_GetUserAssignment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		projectID        int64
		userAssignmentID int64
		setupMock        func(mux *http.ServeMux)
		want             *harvest.ProjectUserAssignment
		wantErr          bool
	}{
		{
			name:             "Valid User Assignment Retrieval",
			projectID:        14308069,
			userAssignmentID: 125068554,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/projects/14308069/user_assignments/125068554", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testWriteResponse(t, w, "project_user_assignment/get/response_1.json")
				})
			},
			want: &harvest.ProjectUserAssignment{
				ID:               harvest.Int64(125068554),
				IsProjectManager: harvest.Bool(true),
				IsActive:         harvest.Bool(true),
				HourlyRate:       harvest.Float64(100.0),
				CreatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
				UpdatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
				Project: &harvest.Project{
					ID:   harvest.Int64(14308069),
					Name: harvest.String("Online Store - Phase 1"),
					Code: harvest.String("OS1"),
				},
				User: &harvest.User{
					ID:   harvest.Int64(1782959),
					Name: harvest.String("Kim Allen"),
				},
			},
			wantErr: false,
		},
		{
			name:             "User Assignment Not Found",
			projectID:        14308069,
			userAssignmentID: 999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/projects/14308069/user_assignments/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					http.Error(w, `{"message":"User assignment not found"}`, http.StatusNotFound)
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

			got, _, err := service.Project.GetUserAssignment(context.Background(), tt.projectID, tt.userAssignmentID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestProjectUserAssignment_String(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)
	updatedAt := time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)

	tests := []struct {
		name string
		in   harvest.ProjectUserAssignment
		want string
	}{
		{
			name: "Assignment with all fields",
			in: harvest.ProjectUserAssignment{
				ID:               harvest.Int64(125068554),
				IsProjectManager: harvest.Bool(true),
				IsActive:         harvest.Bool(true),
				HourlyRate:       harvest.Float64(100.0),
				Budget:           harvest.Float64(5000.0),
				CreatedAt:        &createdAt,
				UpdatedAt:        &updatedAt,
				Project: &harvest.Project{
					ID:   harvest.Int64(14308069),
					Name: harvest.String("Online Store - Phase 1"),
					Code: harvest.String("OS1"),
				},
				User: &harvest.User{
					ID:   harvest.Int64(1782959),
					Name: harvest.String("Kim Allen"),
				},
			},
			want: `harvest.ProjectUserAssignment{ID:125068554, Project:harvest.Project{ID:14308069, Name:"Online Store - Phase 1", Code:"OS1"}, User:harvest.User{ID:1782959, Name:"Kim Allen"}, IsActive:true, IsProjectManager:true, HourlyRate:100, Budget:5000, CreatedAt:time.Time{2017-06-26 22:32:52 +0000 UTC}, UpdatedAt:time.Time{2017-06-26 22:32:52 +0000 UTC}}`, //nolint: lll
		},
		{
			name: "Assignment with minimal fields",
			in: harvest.ProjectUserAssignment{
				ID:         harvest.Int64(125068554),
				IsActive:   harvest.Bool(true),
				HourlyRate: harvest.Float64(100.0),
			},
			want: `harvest.ProjectUserAssignment{ID:125068554, IsActive:true, HourlyRate:100}`,
		},
		{
			name: "Assignment with project and user only",
			in: harvest.ProjectUserAssignment{
				ID: harvest.Int64(125068554),
				Project: &harvest.Project{
					ID:   harvest.Int64(14308069),
					Name: harvest.String("Online Store - Phase 1"),
				},
				User: &harvest.User{
					ID:   harvest.Int64(1782959),
					Name: harvest.String("Kim Allen"),
				},
			},
			want: `harvest.ProjectUserAssignment{ID:125068554, Project:harvest.Project{ID:14308069, Name:"Online Store - Phase 1"}, User:harvest.User{ID:1782959, Name:"Kim Allen"}}`, //nolint: lll
		},
		{
			name: "Assignment with project manager",
			in: harvest.ProjectUserAssignment{
				ID:               harvest.Int64(125068554),
				IsActive:         harvest.Bool(true),
				IsProjectManager: harvest.Bool(true),
				HourlyRate:       harvest.Float64(150.0),
				Budget:           harvest.Float64(10000.0),
			},
			want: `harvest.ProjectUserAssignment{ID:125068554, IsActive:true, IsProjectManager:true, HourlyRate:150, Budget:10000}`, //nolint: lll
		},
		{
			name: "Empty Assignment",
			in:   harvest.ProjectUserAssignment{},
			want: `harvest.ProjectUserAssignment{}`,
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

func TestProjectUserAssignmentList_String(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)
	updatedAt := time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)

	tests := []struct {
		name string
		in   harvest.ProjectUserAssignmentList
		want string
	}{
		{
			name: "AssignmentList with multiple assignments",
			in: harvest.ProjectUserAssignmentList{
				UserAssignments: []*harvest.ProjectUserAssignment{
					{
						ID:         harvest.Int64(125068554),
						IsActive:   harvest.Bool(true),
						HourlyRate: harvest.Float64(100.0),
						Project: &harvest.Project{
							ID:   harvest.Int64(14308069),
							Name: harvest.String("Online Store - Phase 1"),
						},
						User: &harvest.User{
							ID:   harvest.Int64(1782959),
							Name: harvest.String("Kim Allen"),
						},
					},
					{
						ID:         harvest.Int64(125066109),
						IsActive:   harvest.Bool(true),
						HourlyRate: harvest.Float64(100.0),
						Project: &harvest.Project{
							ID:   harvest.Int64(14308069),
							Name: harvest.String("Online Store - Phase 1"),
						},
						User: &harvest.User{
							ID:   harvest.Int64(1782884),
							Name: harvest.String("Jeremy Israelsen"),
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
			want: `harvest.ProjectUserAssignmentList{UserAssignments:[harvest.ProjectUserAssignment{ID:125068554, Project:harvest.Project{ID:14308069, Name:"Online Store - Phase 1"}, User:harvest.User{ID:1782959, Name:"Kim Allen"}, IsActive:true, HourlyRate:100} harvest.ProjectUserAssignment{ID:125066109, Project:harvest.Project{ID:14308069, Name:"Online Store - Phase 1"}, User:harvest.User{ID:1782884, Name:"Jeremy Israelsen"}, IsActive:true, HourlyRate:100}], Pagination:harvest.Pagination{PerPage:100, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
		},
		{
			name: "AssignmentList with single assignment",
			in: harvest.ProjectUserAssignmentList{
				UserAssignments: []*harvest.ProjectUserAssignment{
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
			want: `harvest.ProjectUserAssignmentList{UserAssignments:[harvest.ProjectUserAssignment{ID:999, IsActive:true, IsProjectManager:true}], Pagination:harvest.Pagination{PerPage:50, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
		},
		{
			name: "Empty AssignmentList",
			in: harvest.ProjectUserAssignmentList{
				UserAssignments: []*harvest.ProjectUserAssignment{},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(0),
					TotalEntries: harvest.Int(0),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.ProjectUserAssignmentList{UserAssignments:[], Pagination:harvest.Pagination{PerPage:100, TotalPages:0, TotalEntries:0, Page:1}}`, //nolint: lll
		},
		{
			name: "AssignmentList with Links",
			in: harvest.ProjectUserAssignmentList{
				UserAssignments: []*harvest.ProjectUserAssignment{
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
						First:    harvest.String("https://api.harvestapp.com/v2/projects/14308069/user_assignments?page=1&per_page=25"),
						Next:     harvest.String("https://api.harvestapp.com/v2/projects/14308069/user_assignments?page=3&per_page=25"),
						Previous: harvest.String("https://api.harvestapp.com/v2/projects/14308069/user_assignments?page=1&per_page=25"),
						Last:     harvest.String("https://api.harvestapp.com/v2/projects/14308069/user_assignments?page=3&per_page=25"),
					},
				},
			},
			want: `harvest.ProjectUserAssignmentList{UserAssignments:[harvest.ProjectUserAssignment{ID:100, IsActive:true}], Pagination:harvest.Pagination{PerPage:25, TotalPages:3, TotalEntries:75, Page:2, Links:harvest.PageLinks{First:"https://api.harvestapp.com/v2/projects/14308069/user_assignments?page=1&per_page=25", Next:"https://api.harvestapp.com/v2/projects/14308069/user_assignments?page=3&per_page=25", Previous:"https://api.harvestapp.com/v2/projects/14308069/user_assignments?page=1&per_page=25", Last:"https://api.harvestapp.com/v2/projects/14308069/user_assignments?page=3&per_page=25"}}}`, //nolint: lll
		},
		{
			name: "AssignmentList with full details",
			in: harvest.ProjectUserAssignmentList{
				UserAssignments: []*harvest.ProjectUserAssignment{
					{
						ID:               harvest.Int64(1),
						IsProjectManager: harvest.Bool(true),
						IsActive:         harvest.Bool(true),
						HourlyRate:       harvest.Float64(150.0),
						Budget:           harvest.Float64(5000.0),
						CreatedAt:        &createdAt,
						UpdatedAt:        &updatedAt,
						Project: &harvest.Project{
							ID:   harvest.Int64(10),
							Name: harvest.String("Project Alpha"),
							Code: harvest.String("PA"),
						},
						User: &harvest.User{
							ID:   harvest.Int64(20),
							Name: harvest.String("John Doe"),
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
			want: `harvest.ProjectUserAssignmentList{UserAssignments:[harvest.ProjectUserAssignment{ID:1, Project:harvest.Project{ID:10, Name:"Project Alpha", Code:"PA"}, User:harvest.User{ID:20, Name:"John Doe"}, IsActive:true, IsProjectManager:true, HourlyRate:150, Budget:5000, CreatedAt:time.Time{2017-06-26 22:32:52 +0000 UTC}, UpdatedAt:time.Time{2017-06-26 22:32:52 +0000 UTC}}], Pagination:harvest.Pagination{PerPage:10, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
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
