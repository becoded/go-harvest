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
	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	type args struct {
		userID int64
	}

	tests := []struct {
		name       string
		args       args
		method     string
		path       string
		formValues values
		body       string
		response   string
		want       *harvest.UserProjectAssignmentList
		wantErr    error
	}{
		{
			name: "get 1",
			args: args{
				userID: 1782959,
			},
			method:     "GET",
			path:       "/users/1782959/project_assignments",
			formValues: values{},
			body:       "user/list_project_assignments/body_1.json",
			response:   "user/list_project_assignments/response_1.json",
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
			wantErr: nil,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.method)
				testFormValues(t, r, tt.formValues)
				testBody(t, r, tt.body)
				testWriteResponse(t, w, tt.response)
			})

			userProjectAssignments, _, err := service.User.ListProjectAssignments(context.Background(), tt.args.userID, nil)

			if tt.wantErr != nil {
				assert.Nil(t, userProjectAssignments)
				assert.EqualError(t, err, tt.wantErr.Error())

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.String(), userProjectAssignments.String())
			t.Log(userProjectAssignments.String())
		})
	}
}

func TestUserService_GetMyProjectAssignments(t *testing.T) {
	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	tests := []struct {
		name       string
		method     string
		path       string
		formValues values
		body       string
		response   string
		want       *harvest.UserProjectAssignmentList
		wantErr    error
	}{
		{
			name:       "get 1",
			method:     "GET",
			path:       "/users/me/project_assignments",
			formValues: values{},
			body:       "user/current_list_project_assignments/body_1.json",
			response:   "user/current_list_project_assignments/response_1.json",
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
			wantErr: nil,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.method)
				testFormValues(t, r, tt.formValues)
				testBody(t, r, tt.body)
				testWriteResponse(t, w, tt.response)
			})

			userProjectAssignments, _, err := service.User.GetMyProjectAssignments(context.Background(), nil)

			if tt.wantErr != nil {
				assert.Nil(t, userProjectAssignments)
				assert.EqualError(t, err, tt.wantErr.Error())

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.String(), userProjectAssignments.String())
			t.Log(userProjectAssignments.String())
		})
	}
}
