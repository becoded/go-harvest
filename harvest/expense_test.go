package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/becoded/go-harvest/harvest"
	"github.com/stretchr/testify/assert"
)

func TestExpenseService_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(mux *http.ServeMux)
		want      *harvest.ExpenseList
		wantErr   bool
	}{
		{
			name: "Valid Expense List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expenses", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testWriteResponse(t, w, "expenses/list/response_1.json")
				})
			},
			want: &harvest.ExpenseList{
				Expenses: []*harvest.Expense{
					{
						ID:           harvest.Int64(15296442),
						Notes:        harvest.String("Lunch with client"),
						TotalCost:    harvest.Float64(33.35),
						Units:        harvest.Float64(1.0),
						IsClosed:     harvest.Bool(false),
						IsLocked:     harvest.Bool(true),
						IsBilled:     harvest.Bool(true),
						LockedReason: harvest.String("Expense is invoiced."),
						SpentDate:    harvest.DateP(harvest.Date{Time: time.Date(2017, 3, 3, 0, 0, 0, 0, time.UTC)}),
						CreatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 9, 54, 0, time.UTC)),
						UpdatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 47, 14, 0, time.UTC)),
						Billable:     harvest.Bool(true),
						Receipt: &harvest.Receipt{
							URL:         harvest.String("https://{ACCOUNT_SUBDOMAIN}.harvestapp.com/expenses/15296442/receipt"),
							FileName:    harvest.String("lunch_receipt.gif"),
							FileSize:    harvest.Int64(39410),
							ContentType: harvest.String("image/gif"),
						},
						User: &harvest.User{
							ID:   harvest.Int64(1782959),
							Name: harvest.String("Kim Allen"),
						},
						UserAssignment: &harvest.ProjectUserAssignment{
							ID:               harvest.Int64(125068553),
							IsProjectManager: harvest.Bool(true),
							IsActive:         harvest.Bool(true),
							CreatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
							UpdatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
							HourlyRate:       harvest.Float64(100.0),
						},
						Project: &harvest.Project{
							ID:   harvest.Int64(14307913),
							Name: harvest.String("Marketing Website"),
							Code: harvest.String("MW"),
						},
						ExpenseCategory: &harvest.ExpenseCategory{
							ID:   harvest.Int64(4195926),
							Name: harvest.String("Meals"),
						},
						Client: &harvest.Client{
							ID:       harvest.Int64(5735774),
							Name:     harvest.String("ABC Corp"),
							Currency: harvest.String("USD"),
						},
						Invoice: &harvest.Invoice{
							ID:     harvest.Int64(13150403),
							Number: harvest.String("1001"),
						},
					},
					{
						ID:           harvest.Int64(15296423),
						Notes:        harvest.String("Hotel stay for meeting"),
						TotalCost:    harvest.Float64(100.0),
						Units:        harvest.Float64(1.0),
						IsClosed:     harvest.Bool(true),
						IsLocked:     harvest.Bool(true),
						IsBilled:     harvest.Bool(false),
						LockedReason: harvest.String("The project is locked for this time period."),
						SpentDate:    harvest.DateP(harvest.Date{Time: time.Date(2017, 3, 1, 0, 0, 0, 0, time.UTC)}),
						CreatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 9, 17, 0, time.UTC)),
						UpdatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 47, 14, 0, time.UTC)),
						Billable:     harvest.Bool(true),
						Receipt:      nil,
						User: &harvest.User{
							ID:   harvest.Int64(1782959),
							Name: harvest.String("Kim Allen"),
						},
						UserAssignment: &harvest.ProjectUserAssignment{
							ID:               harvest.Int64(125068554),
							IsProjectManager: harvest.Bool(true),
							IsActive:         harvest.Bool(true),
							CreatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
							UpdatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
							HourlyRate:       harvest.Float64(100.0),
						},
						Project: &harvest.Project{
							ID:   harvest.Int64(14308069),
							Name: harvest.String("Online Store - Phase 1"),
							Code: harvest.String("OS1"),
						},
						ExpenseCategory: &harvest.ExpenseCategory{
							ID:   harvest.Int64(4197501),
							Name: harvest.String("Lodging"),
						},
						Client: &harvest.Client{
							ID:       harvest.Int64(5735776),
							Name:     harvest.String("123 Industries"),
							Currency: harvest.String("EUR"),
						},
						Invoice: nil,
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
						First:    harvest.String("https://api.harvestapp.com/v2/expenses?page=1&per_page=2000"),
						Next:     nil,
						Previous: nil,
						Last:     harvest.String("https://api.harvestapp.com/v2/expenses?page=1&per_page=2000"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Error Fetching Expense List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expenses", func(w http.ResponseWriter, r *http.Request) {
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
			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			tt.setupMock(mux)

			got, _, err := service.Expense.List(context.Background(), &harvest.ExpenseListOptions{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestExpenseService_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		expenseID int64
		setupMock func(mux *http.ServeMux)
		want      *harvest.Expense
		wantErr   bool
	}{
		{
			name:      "Valid Expense Retrieval",
			expenseID: 15296442,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expenses/15296442", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testWriteResponse(t, w, "expenses/get/response_1.json")
				})
			},
			want: &harvest.Expense{
				ID:           harvest.Int64(15296442),
				Notes:        harvest.String("Lunch with client"),
				TotalCost:    harvest.Float64(33.35),
				Units:        harvest.Float64(1.0),
				IsClosed:     harvest.Bool(false),
				IsLocked:     harvest.Bool(true),
				IsBilled:     harvest.Bool(true),
				LockedReason: harvest.String("Expense is invoiced."),
				SpentDate:    harvest.DateP(harvest.Date{Time: time.Date(2017, 3, 3, 0, 0, 0, 0, time.UTC)}),
				CreatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 9, 54, 0, time.UTC)),
				UpdatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 47, 14, 0, time.UTC)),
				Billable:     harvest.Bool(true),
				Receipt: &harvest.Receipt{
					URL:         harvest.String("https://{ACCOUNT_SUBDOMAIN}.harvestapp.com/expenses/15296442/receipt"),
					FileName:    harvest.String("lunch_receipt.gif"),
					FileSize:    harvest.Int64(39410),
					ContentType: harvest.String("image/gif"),
				},
				User: &harvest.User{
					ID:   harvest.Int64(1782959),
					Name: harvest.String("Kim Allen"),
				},
				UserAssignment: &harvest.ProjectUserAssignment{
					ID:               harvest.Int64(125068553),
					IsProjectManager: harvest.Bool(true),
					IsActive:         harvest.Bool(true),
					CreatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
					UpdatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
					HourlyRate:       harvest.Float64(100.0),
				},
				Project: &harvest.Project{
					ID:   harvest.Int64(14307913),
					Name: harvest.String("Marketing Website"),
					Code: harvest.String("MW"),
				},
				ExpenseCategory: &harvest.ExpenseCategory{
					ID:   harvest.Int64(4195926),
					Name: harvest.String("Meals"),
				},
				Client: &harvest.Client{
					ID:       harvest.Int64(5735774),
					Name:     harvest.String("ABC Corp"),
					Currency: harvest.String("USD"),
				},
				Invoice: &harvest.Invoice{
					ID:     harvest.Int64(13150403),
					Number: harvest.String("1001"),
				},
			},
			wantErr: false,
		},
		{
			name:      "Expense Not Found",
			expenseID: 999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expenses/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					http.Error(w, `{"message":"Expense not found"}`, http.StatusNotFound)
				})
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			tt.setupMock(mux)

			got, _, err := service.Expense.Get(context.Background(), tt.expenseID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestExpenseService_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		request   *harvest.ExpenseCreateRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.Expense
		wantErr   bool
	}{
		{
			name: "Valid Expense Creation",
			request: &harvest.ExpenseCreateRequest{
				UserID:            harvest.Int64(1782959),
				ProjectID:         harvest.Int64(14308069),
				ExpenseCategoryID: harvest.Int64(4195926),
				SpentDate:         harvest.DateP(harvest.Date{Time: time.Date(2017, 3, 1, 0, 0, 0, 0, time.UTC)}),
				TotalCost:         harvest.Float64(13.59),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expenses", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testBody(t, r, "expenses/create/body_1.json")
					testWriteResponse(t, w, "expenses/create/response_1.json")
				})
			},
			want: &harvest.Expense{
				ID:           harvest.Int64(15297032),
				Notes:        nil,
				TotalCost:    harvest.Float64(13.59),
				Units:        harvest.Float64(1.0),
				IsClosed:     harvest.Bool(false),
				IsLocked:     harvest.Bool(false),
				IsBilled:     harvest.Bool(false),
				LockedReason: nil,
				SpentDate:    harvest.DateP(harvest.Date{Time: time.Date(2017, 3, 1, 0, 0, 0, 0, time.UTC)}),
				CreatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 42, 27, 0, time.UTC)),
				UpdatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 42, 27, 0, time.UTC)),
				Billable:     harvest.Bool(true),
				Receipt:      nil,
				User: &harvest.User{
					ID:   harvest.Int64(1782959),
					Name: harvest.String("Kim Allen"),
				},
				UserAssignment: &harvest.ProjectUserAssignment{
					ID:               harvest.Int64(125068553),
					IsProjectManager: harvest.Bool(true),
					IsActive:         harvest.Bool(true),
					CreatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
					UpdatedAt:        harvest.TimeTimeP(time.Date(2017, 6, 26, 22, 32, 52, 0, time.UTC)),
					HourlyRate:       harvest.Float64(100.0),
				},
				Project: &harvest.Project{
					ID:   harvest.Int64(14308069),
					Name: harvest.String("Online Store - Phase 1"),
					Code: harvest.String("OS1"),
				},
				ExpenseCategory: &harvest.ExpenseCategory{
					ID:   harvest.Int64(4195926),
					Name: harvest.String("Meals"),
				},
				Client: &harvest.Client{
					ID:       harvest.Int64(5735776),
					Name:     harvest.String("123 Industries"),
					Currency: harvest.String("EUR"),
				},
				Invoice: nil,
			},
			wantErr: false,
		},
		{
			name: "Invalid Expense Creation - Missing ProjectID",
			request: &harvest.ExpenseCreateRequest{
				ExpenseCategoryID: harvest.Int64(4195926),
				SpentDate:         harvest.DateP(harvest.Date{Time: time.Date(2017, 3, 1, 0, 0, 0, 0, time.UTC)}),
				TotalCost:         harvest.Float64(13.59),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expenses", func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, `{"message":"ProjectID is required"}`, http.StatusBadRequest)
				})
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			tt.setupMock(mux)

			got, _, err := service.Expense.Create(context.Background(), tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestExpenseService_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		expenseID int64
		request   *harvest.ExpenseUpdateRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.Expense
		wantErr   bool
	}{
		{
			name:      "Valid Expense Update",
			expenseID: 15297032,
			request: &harvest.ExpenseUpdateRequest{
				Notes: harvest.String("Dinner"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expenses/15297032", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					testBody(t, r, "expenses/update/body_1.json")
					testWriteResponse(t, w, "expenses/update/response_1.json")
				})
			},
			want: &harvest.Expense{
				ID:        harvest.Int64(15297032),
				Notes:     harvest.String("Dinner"),
				TotalCost: harvest.Float64(13.59),
			},
			wantErr: false,
		},
		{
			name:      "Expense Not Found",
			expenseID: 999,
			request: &harvest.ExpenseUpdateRequest{
				Notes: harvest.String("Dinner"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expenses/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					http.Error(w, `{"message":"Expense not found"}`, http.StatusNotFound)
				})
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			tt.setupMock(mux)

			got, _, err := service.Expense.Update(context.Background(), tt.expenseID, tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestExpenseService_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		expenseID int64
		setupMock func(mux *http.ServeMux)
		wantErr   bool
	}{
		{
			name:      "Valid Expense Deletion",
			expenseID: 15297032,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expenses/15297032", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					w.WriteHeader(http.StatusNoContent)
				})
			},
			wantErr: false,
		},
		{
			name:      "Expense Not Found",
			expenseID: 999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expenses/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					http.Error(w, `{"message":"Expense not found"}`, http.StatusNotFound)
				})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mux, teardown := setup(t)
			t.Cleanup(teardown)

			tt.setupMock(mux)

			_, err := service.Expense.Delete(context.Background(), tt.expenseID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
