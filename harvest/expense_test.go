package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
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
			t.Parallel()

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
			t.Parallel()

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
				mux.HandleFunc("/expenses", func(w http.ResponseWriter, _ *http.Request) {
					http.Error(w, `{"message":"ProjectID is required"}`, http.StatusBadRequest)
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
				ID:           harvest.Int64(15297032),
				Notes:        harvest.String("Dinner"),
				TotalCost:    harvest.Float64(13.59),
				Units:        harvest.Float64(1.0),
				IsClosed:     harvest.Bool(false),
				IsLocked:     harvest.Bool(false),
				IsBilled:     harvest.Bool(false),
				LockedReason: nil,
				SpentDate:    harvest.DateP(harvest.Date{Time: time.Date(2017, 3, 1, 0, 0, 0, 0, time.UTC)}),
				CreatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 42, 27, 0, time.UTC)),
				UpdatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 45, 51, 0, time.UTC)),
				Billable:     harvest.Bool(true),
				Receipt: &harvest.Receipt{
					URL:         harvest.String("https://{ACCOUNT_SUBDOMAIN}.harvestapp.com/expenses/15297032/receipt"),
					FileName:    harvest.String("dinner_receipt.gif"),
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
					ID:   harvest.Int64(14308069),
					Name: harvest.String("Online Store - Phase 1"),
					Code: harvest.String("OS1"),
				},
				ExpenseCategory: &harvest.ExpenseCategory{
					ID:        harvest.Int64(4195926),
					Name:      harvest.String("Meals"),
					UnitPrice: nil,
					UnitName:  nil,
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
			t.Parallel()

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
			t.Parallel()

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

func TestExpense_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.Expense
		want string
	}{
		{
			name: "Expense with all fields",
			in: harvest.Expense{
				ID: harvest.Int64(15296442),
				Client: &harvest.Client{
					ID:       harvest.Int64(5735776),
					Name:     harvest.String("ABC Corp"),
					Currency: harvest.String("USD"),
				},
				Project: &harvest.Project{
					ID:   harvest.Int64(14308069),
					Name: harvest.String("Marketing"),
					Code: harvest.String("MKT"),
				},
				ExpenseCategory: &harvest.ExpenseCategory{
					ID:        harvest.Int64(4195926),
					Name:      harvest.String("Meals"),
					UnitPrice: harvest.Float64(0),
					UnitName:  harvest.String(""),
				},
				User: &harvest.User{
					ID:   harvest.Int64(1782959),
					Name: harvest.String("Kim Allen"),
				},
				Receipt: &harvest.Receipt{
					URL:         harvest.String("https://cache.harvestapp.com/receipts/1.pdf"),
					FileName:    harvest.String("receipt.pdf"),
					FileSize:    harvest.Int64(1024),
					ContentType: harvest.String("application/pdf"),
				},
				Notes:     harvest.String("Lunch with client"),
				Units:     harvest.Float64(1.0),
				TotalCost: harvest.Float64(45.50),
				Billable:  harvest.Bool(true),
				IsClosed:  harvest.Bool(false),
				IsLocked:  harvest.Bool(false),
				IsBilled:  harvest.Bool(false),
				SpentDate: &harvest.Date{Time: time.Date(2017, 3, 21, 0, 0, 0, 0, time.UTC)},
				CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 49, 28, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 47, 14, 0, time.UTC)),
			},
			want: `harvest.Expense{ID:15296442, Client:harvest.Client{ID:5735776, Name:"ABC Corp", Currency:"USD"}, Project:harvest.Project{ID:14308069, Name:"Marketing", Code:"MKT"}, ExpenseCategory:harvest.ExpenseCategory{ID:4195926, Name:"Meals", UnitName:"", UnitPrice:0}, User:harvest.User{ID:1782959, Name:"Kim Allen"}, Receipt:harvest.Receipt{URL:"https://cache.harvestapp.com/receipts/1.pdf", FileName:"receipt.pdf", FileSize:1024, ContentType:"application/pdf"}, Notes:"Lunch with client", Units:1, TotalCost:45.5, Billable:true, IsClosed:false, IsLocked:false, IsBilled:false, SpentDate:harvest.Date{{2017-03-21 00:00:00 +0000 UTC}}, CreatedAt:time.Time{2017-06-27 15:49:28 +0000 UTC}, UpdatedAt:time.Time{2017-06-27 16:47:14 +0000 UTC}}`, //nolint: lll
		},
		{
			name: "Expense with minimal fields",
			in: harvest.Expense{
				ID:        harvest.Int64(999),
				TotalCost: harvest.Float64(100.00),
			},
			want: `harvest.Expense{ID:999, TotalCost:100}`,
		},
		{
			name: "Expense with boolean flags",
			in: harvest.Expense{
				ID:        harvest.Int64(15296442),
				TotalCost: harvest.Float64(45.50),
				Billable:  harvest.Bool(true),
				IsClosed:  harvest.Bool(false),
				IsLocked:  harvest.Bool(true),
				IsBilled:  harvest.Bool(false),
				SpentDate: &harvest.Date{Time: time.Date(2017, 3, 21, 0, 0, 0, 0, time.UTC)},
			},
			want: `harvest.Expense{ID:15296442, TotalCost:45.5, Billable:true, IsClosed:false, IsLocked:true, IsBilled:false, SpentDate:harvest.Date{{2017-03-21 00:00:00 +0000 UTC}}}`, //nolint: lll
		},
		{
			name: "Expense with invoice",
			in: harvest.Expense{
				ID:        harvest.Int64(15296442),
				TotalCost: harvest.Float64(45.50),
				Invoice: &harvest.Invoice{
					ID:     harvest.Int64(13150403),
					Number: harvest.String("1001"),
				},
				IsBilled: harvest.Bool(true),
			},
			want: `harvest.Expense{ID:15296442, Invoice:harvest.Invoice{ID:13150403, Number:"1001"}, TotalCost:45.5, IsBilled:true}`, //nolint: lll
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

func TestReceipt_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.Receipt
		want string
	}{
		{
			name: "Receipt with all fields",
			in: harvest.Receipt{
				URL:         harvest.String("https://cache.harvestapp.com/receipts/1.pdf"),
				FileName:    harvest.String("receipt.pdf"),
				FileSize:    harvest.Int64(1024),
				ContentType: harvest.String("application/pdf"),
			},
			want: `harvest.Receipt{URL:"https://cache.harvestapp.com/receipts/1.pdf", FileName:"receipt.pdf", FileSize:1024, ContentType:"application/pdf"}`, //nolint: lll
		},
		{
			name: "Receipt with minimal fields",
			in: harvest.Receipt{
				URL:      harvest.String("https://example.com/receipt.jpg"),
				FileName: harvest.String("receipt.jpg"),
			},
			want: `harvest.Receipt{URL:"https://example.com/receipt.jpg", FileName:"receipt.jpg"}`,
		},
		{
			name: "Empty Receipt",
			in:   harvest.Receipt{},
			want: `harvest.Receipt{}`,
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

func TestExpenseList_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.ExpenseList
		want string
	}{
		{
			name: "ExpenseList with multiple expenses",
			in: harvest.ExpenseList{
				Expenses: []*harvest.Expense{
					{
						ID:        harvest.Int64(15296442),
						TotalCost: harvest.Float64(45.50),
						Notes:     harvest.String("Lunch with client"),
						Billable:  harvest.Bool(true),
						SpentDate: &harvest.Date{Time: time.Date(2017, 3, 21, 0, 0, 0, 0, time.UTC)},
						CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 27, 15, 49, 28, 0, time.UTC)),
					},
					{
						ID:        harvest.Int64(15296443),
						TotalCost: harvest.Float64(120.00),
						Notes:     harvest.String("Hotel"),
						Billable:  harvest.Bool(true),
						SpentDate: &harvest.Date{Time: time.Date(2017, 3, 22, 0, 0, 0, 0, time.UTC)},
						CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 28, 10, 30, 0, 0, time.UTC)),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.ExpenseList{Expenses:[harvest.Expense{ID:15296442, Notes:"Lunch with client", TotalCost:45.5, Billable:true, SpentDate:harvest.Date{{2017-03-21 00:00:00 +0000 UTC}}, CreatedAt:time.Time{2017-06-27 15:49:28 +0000 UTC}} harvest.Expense{ID:15296443, Notes:"Hotel", TotalCost:120, Billable:true, SpentDate:harvest.Date{{2017-03-22 00:00:00 +0000 UTC}}, CreatedAt:time.Time{2017-06-28 10:30:00 +0000 UTC}}], Pagination:harvest.Pagination{PerPage:100, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
		},
		{
			name: "ExpenseList with single expense",
			in: harvest.ExpenseList{
				Expenses: []*harvest.Expense{
					{
						ID:        harvest.Int64(999),
						TotalCost: harvest.Float64(50.00),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.ExpenseList{Expenses:[harvest.Expense{ID:999, TotalCost:50}], Pagination:harvest.Pagination{PerPage:100, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
		},
		{
			name: "Empty ExpenseList",
			in: harvest.ExpenseList{
				Expenses: []*harvest.Expense{},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(0),
					TotalEntries: harvest.Int(0),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.ExpenseList{Expenses:[], Pagination:harvest.Pagination{PerPage:100, TotalPages:0, TotalEntries:0, Page:1}}`, //nolint: lll
		},
		{
			name: "ExpenseList with Links",
			in: harvest.ExpenseList{
				Expenses: []*harvest.Expense{
					{
						ID:        harvest.Int64(100),
						TotalCost: harvest.Float64(25.00),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(50),
					TotalPages:   harvest.Int(3),
					TotalEntries: harvest.Int(150),
					Page:         harvest.Int(2),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/expenses?page=1&per_page=50"),
						Next:     harvest.String("https://api.harvestapp.com/v2/expenses?page=3&per_page=50"),
						Previous: harvest.String("https://api.harvestapp.com/v2/expenses?page=1&per_page=50"),
						Last:     harvest.String("https://api.harvestapp.com/v2/expenses?page=3&per_page=50"),
					},
				},
			},
			want: `harvest.ExpenseList{Expenses:[harvest.Expense{ID:100, TotalCost:25}], Pagination:harvest.Pagination{PerPage:50, TotalPages:3, TotalEntries:150, Page:2, Links:harvest.PageLinks{First:"https://api.harvestapp.com/v2/expenses?page=1&per_page=50", Next:"https://api.harvestapp.com/v2/expenses?page=3&per_page=50", Previous:"https://api.harvestapp.com/v2/expenses?page=1&per_page=50", Last:"https://api.harvestapp.com/v2/expenses?page=3&per_page=50"}}}`, //nolint: lll
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
