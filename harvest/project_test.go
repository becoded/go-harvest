package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestProjectService_List(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name      string
		setupMock func(mux *http.ServeMux)
		want      *harvest.ProjectList
		wantErr   bool
	}{
		{
			name: "Valid Project List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "project/list/body_1.json")
					testWriteResponse(t, w, "project/list/response_1.json")
				})
			},
			want: &harvest.ProjectList{
				Projects: []*harvest.Project{
					{
						ID: harvest.Int64(14308069),
						Client: &harvest.Client{
							ID:       harvest.Int64(5735776),
							Name:     harvest.String("123 Industries"),
							Currency: harvest.String("EUR"),
						},
						Name:                             harvest.String("Online Store - Phase 1"),
						Code:                             harvest.String("OS1"),
						IsActive:                         harvest.Bool(true),
						IsBillable:                       harvest.Bool(true),
						IsFixedFee:                       harvest.Bool(false),
						BillBy:                           harvest.String("Project"),
						HourlyRate:                       harvest.Float64(100),
						Budget:                           harvest.Float64(200),
						BudgetBy:                         harvest.String("project"),
						BudgetIsMonthly:                  harvest.Bool(false),
						NotifyWhenOverBudget:             harvest.Bool(true),
						OverBudgetNotificationPercentage: harvest.Float64(80),
						ShowBudgetToAll:                  harvest.Bool(false),
						CostBudgetIncludeExpenses:        harvest.Bool(false),
						Notes:                            harvest.String(""),
						StartsOn:                         &harvest.Date{Time: time.Date(2017, 6, 1, 0, 0, 0, 0, time.UTC)},
						CreatedAt:                        harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
						UpdatedAt:                        harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 54, 6, 0, time.UTC)),
					},
					{
						ID: harvest.Int64(14307913),
						Client: &harvest.Client{
							ID:       harvest.Int64(5735774),
							Name:     harvest.String("ABC Corp"),
							Currency: harvest.String("USD"),
						},
						Name:                             harvest.String("Marketing Website"),
						Code:                             harvest.String("MW"),
						IsActive:                         harvest.Bool(true),
						IsBillable:                       harvest.Bool(true),
						IsFixedFee:                       harvest.Bool(false),
						BillBy:                           harvest.String("Project"),
						HourlyRate:                       harvest.Float64(100),
						Budget:                           harvest.Float64(50),
						BudgetBy:                         harvest.String("project"),
						BudgetIsMonthly:                  harvest.Bool(false),
						NotifyWhenOverBudget:             harvest.Bool(true),
						OverBudgetNotificationPercentage: harvest.Float64(80),
						ShowBudgetToAll:                  harvest.Bool(false),
						CostBudgetIncludeExpenses:        harvest.Bool(false),
						Notes:                            harvest.String(""),
						StartsOn: &harvest.Date{Time: time.Date(
							2017, 1, 1, 0, 0, 0, 0, time.UTC)},
						EndsOn: &harvest.Date{Time: time.Date(
							2017, 3, 31, 0, 0, 0, 0, time.UTC)},
						CreatedAt: harvest.TimeTimeP(time.Date(
							2017, 6, 26, 21, 36, 23, 0, time.UTC)),
						UpdatedAt: harvest.TimeTimeP(time.Date(
							2017, 6, 26, 21, 54, 46, 0, time.UTC)),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
					Links: &harvest.PageLinks{
						First: harvest.String("https://api.harvestapp.com/v2/projects?page=1&per_page=100"),
						Last:  harvest.String("https://api.harvestapp.com/v2/projects?page=1&per_page=100"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Error Fetching Project List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Project.List(context.Background(), &harvest.ProjectListOptions{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestProjectService_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		projectID int64
		setupMock func(mux *http.ServeMux)
		want      *harvest.Project
		wantErr   bool
	}{
		{
			name:      "Valid Project Retrieval",
			projectID: 14308069,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/projects/14308069", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "project/get/body_1.json")
					testWriteResponse(t, w, "project/get/response_1.json")
				})
			},
			want: &harvest.Project{
				ID: harvest.Int64(14308069),
				Client: &harvest.Client{
					ID:       harvest.Int64(5735776),
					Name:     harvest.String("123 Industries"),
					Currency: harvest.String("EUR"),
				},
				Name:                             harvest.String("Online Store - Phase 1"),
				Code:                             harvest.String("OS1"),
				IsActive:                         harvest.Bool(true),
				IsBillable:                       harvest.Bool(true),
				IsFixedFee:                       harvest.Bool(false),
				BillBy:                           harvest.String("Project"),
				HourlyRate:                       harvest.Float64(100),
				Budget:                           harvest.Float64(200),
				BudgetBy:                         harvest.String("project"),
				BudgetIsMonthly:                  harvest.Bool(false),
				NotifyWhenOverBudget:             harvest.Bool(true),
				OverBudgetNotificationPercentage: harvest.Float64(80),
				ShowBudgetToAll:                  harvest.Bool(false),
				CostBudgetIncludeExpenses:        harvest.Bool(false),
				Notes:                            harvest.String(""),
				StartsOn:                         &harvest.Date{Time: time.Date(2017, 6, 1, 0, 0, 0, 0, time.UTC)},
				CreatedAt:                        harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 52, 18, 0, time.UTC)),
				UpdatedAt:                        harvest.TimeTimeP(time.Date(2017, 6, 26, 21, 54, 6, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name:      "Project Not Found",
			projectID: 999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/projects/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					http.Error(w, `{"message":"Project not found"}`, http.StatusNotFound)
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

			got, _, err := service.Project.Get(context.Background(), tt.projectID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
