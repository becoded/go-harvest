package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestExpenseCategoryService_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(mux *http.ServeMux)
		want      *harvest.ExpenseCategoryList
		wantErr   bool
	}{
		{
			name: "Valid Expense Category List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expense_categories", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testWriteResponse(t, w, "expense_categories/list/response_1.json")
				})
			},
			want: &harvest.ExpenseCategoryList{
				ExpenseCategories: []*harvest.ExpenseCategory{
					{
						ID:        harvest.Int64(4197501),
						Name:      harvest.String("Lodging"),
						UnitName:  nil,
						UnitPrice: nil,
						IsActive:  harvest.Bool(true),
						CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 1, 32, 0, time.UTC)),
						UpdatedAt: harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 1, 32, 0, time.UTC)),
					},
					{
						ID:        harvest.Int64(4195930),
						Name:      harvest.String("Mileage"),
						UnitName:  harvest.String("mile"),
						UnitPrice: harvest.Float64(0.535),
						IsActive:  harvest.Bool(true),
						CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
						UpdatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
					},
					{
						ID:        harvest.Int64(4195928),
						Name:      harvest.String("Transportation"),
						UnitName:  nil,
						UnitPrice: nil,
						IsActive:  harvest.Bool(true),
						CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
						UpdatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
					},
					{
						ID:        harvest.Int64(4195926),
						Name:      harvest.String("Meals"),
						UnitName:  nil,
						UnitPrice: nil,
						IsActive:  harvest.Bool(true),
						CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
						UpdatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
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
						First:    harvest.String("https://api.harvestapp.com/v2/expense_categories?page=1&per_page=2000"),
						Next:     nil,
						Previous: nil,
						Last:     harvest.String("https://api.harvestapp.com/v2/expense_categories?page=1&per_page=2000"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Error Fetching Expense Category List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expense_categories", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Expense.ListExpenseCategories(context.Background(), &harvest.ExpenseCategoryListOptions{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestExpenseCategoryService_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		expenseCategoryID int64
		setupMock         func(mux *http.ServeMux)
		want              *harvest.ExpenseCategory
		wantErr           bool
	}{
		{
			name:              "Valid Expense Category Retrieval",
			expenseCategoryID: 4197501,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expense_categories/4197501", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testWriteResponse(t, w, "expense_categories/get/response_1.json")
				})
			},
			want: &harvest.ExpenseCategory{
				ID:        harvest.Int64(4197501),
				Name:      harvest.String("Lodging"),
				IsActive:  harvest.Bool(true),
				CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 1, 32, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 1, 32, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name:              "Expense Category Not Found",
			expenseCategoryID: 999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expense_categories/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					http.Error(w, `{"message":"Expense category not found"}`, http.StatusNotFound)
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

			got, _, err := service.Expense.GetExpenseCategory(context.Background(), tt.expenseCategoryID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestExpenseCategoryService_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		request   *harvest.ExpenseCategoryRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.ExpenseCategory
		wantErr   bool
	}{
		{
			name: "Valid Expense Category Creation",
			request: &harvest.ExpenseCategoryRequest{
				Name: harvest.String("Other"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expense_categories", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testBody(t, r, "expense_categories/create/body_1.json")
					testWriteResponse(t, w, "expense_categories/create/response_1.json")
				})
			},
			want: &harvest.ExpenseCategory{
				ID:        harvest.Int64(4197514),
				Name:      harvest.String("Other"),
				IsActive:  harvest.Bool(true),
				CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 4, 23, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 4, 23, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name: "Error Creating Expense Category",
			request: &harvest.ExpenseCategoryRequest{
				Name: harvest.String("Other"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expense_categories", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Expense.CreateExpenseCategory(context.Background(), tt.request)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestExpenseCategoryService_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		expenseCategoryID int64
		request           *harvest.ExpenseCategoryRequest
		setupMock         func(mux *http.ServeMux)
		want              *harvest.ExpenseCategory
		wantErr           bool
	}{
		{
			name:              "Valid Expense Category Update",
			expenseCategoryID: 4197514,
			request: &harvest.ExpenseCategoryRequest{
				IsActive: harvest.Bool(false),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expense_categories/4197514", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					testBody(t, r, "expense_categories/update/body_1.json")
					testWriteResponse(t, w, "expense_categories/update/response_1.json")
				})
			},
			want: &harvest.ExpenseCategory{
				ID:        harvest.Int64(4197514),
				Name:      harvest.String("Other"),
				IsActive:  harvest.Bool(false),
				CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 4, 23, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 4, 58, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name:              "Error Updating Expense Category",
			expenseCategoryID: 4197514,
			request: &harvest.ExpenseCategoryRequest{
				IsActive: harvest.Bool(false),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expense_categories/4197514", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Expense.UpdateExpenseCategory(context.Background(), tt.expenseCategoryID, tt.request)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestExpenseCategoryService_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		expenseCategoryID int64
		setupMock         func(mux *http.ServeMux)
		wantErr           bool
	}{
		{
			name:              "Valid Expense Category Deletion",
			expenseCategoryID: 4197514,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expense_categories/4197514", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					w.WriteHeader(http.StatusNoContent)
				})
			},
			wantErr: false,
		},
		{
			name:              "Expense Category Not Found",
			expenseCategoryID: 999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/expense_categories/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					http.Error(w, `{"message":"Expense category not found"}`, http.StatusNotFound)
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

			_, err := service.Expense.DeleteExpenseCategory(context.Background(), tt.expenseCategoryID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExpenseCategory_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.ExpenseCategory
		want string
	}{
		{
			name: "ExpenseCategory with all fields",
			in: harvest.ExpenseCategory{
				ID:        harvest.Int64(4195930),
				Name:      harvest.String("Mileage"),
				UnitName:  harvest.String("mile"),
				UnitPrice: harvest.Float64(0.535),
				IsActive:  harvest.Bool(true),
				CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 5, 0, time.UTC)),
			},
			want: `harvest.ExpenseCategory{ID:4195930, Name:"Mileage", UnitName:"mile", UnitPrice:0.535, IsActive:true, CreatedAt:time.Time{2017-06-26 20:41:00 +0000 UTC}, UpdatedAt:time.Time{2017-06-26 20:41:05 +0000 UTC}}`, //nolint: lll
		},
		{
			name: "ExpenseCategory with minimal fields",
			in: harvest.ExpenseCategory{
				ID:       harvest.Int64(999),
				Name:     harvest.String("Lodging"),
				IsActive: harvest.Bool(true),
			},
			want: `harvest.ExpenseCategory{ID:999, Name:"Lodging", IsActive:true}`,
		},
		{
			name: "ExpenseCategory without unit fields",
			in: harvest.ExpenseCategory{
				ID:        harvest.Int64(4197501),
				Name:      harvest.String("Meals"),
				IsActive:  harvest.Bool(true),
				CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 1, 32, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 1, 32, 0, time.UTC)),
			},
			want: `harvest.ExpenseCategory{ID:4197501, Name:"Meals", IsActive:true, CreatedAt:time.Time{2017-06-27 15:01:32 +0000 UTC}, UpdatedAt:time.Time{2017-06-27 15:01:32 +0000 UTC}}`, //nolint: lll
		},
		{
			name: "ExpenseCategory inactive",
			in: harvest.ExpenseCategory{
				ID:       harvest.Int64(4197514),
				Name:     harvest.String("Other"),
				IsActive: harvest.Bool(false),
			},
			want: `harvest.ExpenseCategory{ID:4197514, Name:"Other", IsActive:false}`,
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

func TestExpenseCategoryList_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.ExpenseCategoryList
		want string
	}{
		{
			name: "ExpenseCategoryList with multiple categories",
			in: harvest.ExpenseCategoryList{
				ExpenseCategories: []*harvest.ExpenseCategory{
					{
						ID:        harvest.Int64(4197501),
						Name:      harvest.String("Lodging"),
						IsActive:  harvest.Bool(true),
						CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 1, 32, 0, time.UTC)),
						UpdatedAt: harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 1, 32, 0, time.UTC)),
					},
					{
						ID:        harvest.Int64(4195930),
						Name:      harvest.String("Mileage"),
						UnitName:  harvest.String("mile"),
						UnitPrice: harvest.Float64(0.535),
						IsActive:  harvest.Bool(true),
						CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
						UpdatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(2000),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.ExpenseCategoryList{ExpenseCategories:[harvest.ExpenseCategory{ID:4197501, Name:"Lodging", IsActive:true, CreatedAt:time.Time{2017-06-27 15:01:32 +0000 UTC}, UpdatedAt:time.Time{2017-06-27 15:01:32 +0000 UTC}} harvest.ExpenseCategory{ID:4195930, Name:"Mileage", UnitName:"mile", UnitPrice:0.535, IsActive:true, CreatedAt:time.Time{2017-06-26 20:41:00 +0000 UTC}, UpdatedAt:time.Time{2017-06-26 20:41:00 +0000 UTC}}], Pagination:harvest.Pagination{PerPage:2000, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
		},
		{
			name: "ExpenseCategoryList with single category",
			in: harvest.ExpenseCategoryList{
				ExpenseCategories: []*harvest.ExpenseCategory{
					{
						ID:       harvest.Int64(999),
						Name:     harvest.String("Transportation"),
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
			want: `harvest.ExpenseCategoryList{ExpenseCategories:[harvest.ExpenseCategory{ID:999, Name:"Transportation", IsActive:true}], Pagination:harvest.Pagination{PerPage:100, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
		},
		{
			name: "Empty ExpenseCategoryList",
			in: harvest.ExpenseCategoryList{
				ExpenseCategories: []*harvest.ExpenseCategory{},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(0),
					TotalEntries: harvest.Int(0),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.ExpenseCategoryList{ExpenseCategories:[], Pagination:harvest.Pagination{PerPage:100, TotalPages:0, TotalEntries:0, Page:1}}`, //nolint: lll
		},
		{
			name: "ExpenseCategoryList with Links",
			in: harvest.ExpenseCategoryList{
				ExpenseCategories: []*harvest.ExpenseCategory{
					{
						ID:       harvest.Int64(100),
						Name:     harvest.String("Meals"),
						IsActive: harvest.Bool(true),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(50),
					TotalPages:   harvest.Int(3),
					TotalEntries: harvest.Int(150),
					Page:         harvest.Int(2),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/expense_categories?page=1&per_page=50"),
						Next:     harvest.String("https://api.harvestapp.com/v2/expense_categories?page=3&per_page=50"),
						Previous: harvest.String("https://api.harvestapp.com/v2/expense_categories?page=1&per_page=50"),
						Last:     harvest.String("https://api.harvestapp.com/v2/expense_categories?page=3&per_page=50"),
					},
				},
			},
			want: `harvest.ExpenseCategoryList{ExpenseCategories:[harvest.ExpenseCategory{ID:100, Name:"Meals", IsActive:true}], Pagination:harvest.Pagination{PerPage:50, TotalPages:3, TotalEntries:150, Page:2, Links:harvest.PageLinks{First:"https://api.harvestapp.com/v2/expense_categories?page=1&per_page=50", Next:"https://api.harvestapp.com/v2/expense_categories?page=3&per_page=50", Previous:"https://api.harvestapp.com/v2/expense_categories?page=1&per_page=50", Last:"https://api.harvestapp.com/v2/expense_categories?page=3&per_page=50"}}}`, //nolint: lll
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
