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
