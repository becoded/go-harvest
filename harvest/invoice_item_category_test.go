package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestInvoiceService_ListItemCategories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(mux *http.ServeMux)
		want      *harvest.InvoiceItemCategoryList
		wantErr   bool
	}{
		{
			name: "Valid List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoice_item_categories", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testWriteResponse(t, w, "invoice_item_category/list/response_1.json")
				})
			},
			want: &harvest.InvoiceItemCategoryList{
				InvoiceItemCategories: []*harvest.InvoiceItemCategory{
					{
						ID:           harvest.Int64(1466293),
						Name:         harvest.String("Product"),
						UseAsService: harvest.Bool(false),
						UseAsExpense: harvest.Bool(true),
						CreatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
						UpdatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
					},
					{
						ID:           harvest.Int64(1466292),
						Name:         harvest.String("Service"),
						UseAsService: harvest.Bool(true),
						UseAsExpense: harvest.Bool(false),
						CreatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
						UpdatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(2000),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
					Links: &harvest.PageLinks{
						First: harvest.String("https://api.harvestapp.com/v2/invoice_item_categories?page=1&per_page=2000"),
						Last:  harvest.String("https://api.harvestapp.com/v2/invoice_item_categories?page=1&per_page=2000"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Error Fetching List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoice_item_categories", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Invoice.ListItemCategories(context.Background(), &harvest.InvoiceItemCategoryListOptions{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestInvoiceService_GetItemCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        int64
		setupMock func(mux *http.ServeMux)
		want      *harvest.InvoiceItemCategory
		wantErr   bool
	}{
		{
			name: "Valid Get",
			id:   1466293,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoice_item_categories/1466293", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testWriteResponse(t, w, "invoice_item_category/get/response_1.json")
				})
			},
			want: &harvest.InvoiceItemCategory{
				ID:           harvest.Int64(1466293),
				Name:         harvest.String("Product"),
				UseAsService: harvest.Bool(false),
				UseAsExpense: harvest.Bool(true),
				CreatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
				UpdatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name: "Not Found",
			id:   999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoice_item_categories/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					http.Error(w, `{"message":"Not Found"}`, http.StatusNotFound)
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

			got, _, err := service.Invoice.GetItemCategory(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestInvoiceService_CreateItemCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		request   *harvest.InvoiceItemCategoryRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.InvoiceItemCategory
		wantErr   bool
	}{
		{
			name: "Valid Create",
			request: &harvest.InvoiceItemCategoryRequest{
				Name: harvest.String("Hosting"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoice_item_categories", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testBody(t, r, "invoice_item_category/create/body_1.json")
					testWriteResponse(t, w, "invoice_item_category/create/response_1.json")
				})
			},
			want: &harvest.InvoiceItemCategory{
				ID:           harvest.Int64(1467098),
				Name:         harvest.String("Hosting"),
				UseAsService: harvest.Bool(false),
				UseAsExpense: harvest.Bool(false),
				CreatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 20, 59, 0, time.UTC)),
				UpdatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 20, 59, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name: "Error Creating Item Category",
			request: &harvest.InvoiceItemCategoryRequest{
				Name: harvest.String("Hosting"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoice_item_categories", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Invoice.CreateItemCategory(context.Background(), tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestInvoiceService_UpdateItemCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        int64
		request   *harvest.InvoiceItemCategoryRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.InvoiceItemCategory
		wantErr   bool
	}{
		{
			name: "Valid Update",
			id:   1467098,
			request: &harvest.InvoiceItemCategoryRequest{
				Name: harvest.String("Expense"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoice_item_categories/1467098", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					testBody(t, r, "invoice_item_category/update/body_1.json")
					testWriteResponse(t, w, "invoice_item_category/update/response_1.json")
				})
			},
			want: &harvest.InvoiceItemCategory{
				ID:           harvest.Int64(1467098),
				Name:         harvest.String("Expense"),
				UseAsService: harvest.Bool(false),
				UseAsExpense: harvest.Bool(false),
				CreatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 20, 59, 0, time.UTC)),
				UpdatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 21, 26, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name: "Error Updating Item Category",
			id:   1467098,
			request: &harvest.InvoiceItemCategoryRequest{
				Name: harvest.String("Expense"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoice_item_categories/1467098", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Invoice.UpdateItemCategory(context.Background(), tt.id, tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestInvoiceService_DeleteItemCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        int64
		setupMock func(mux *http.ServeMux)
		wantErr   bool
	}{
		{
			name: "Valid Delete",
			id:   1467098,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoice_item_categories/1467098", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					w.WriteHeader(http.StatusNoContent)
				})
			},
			wantErr: false,
		},
		{
			name: "Not Found",
			id:   999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/invoice_item_categories/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					http.Error(w, `{"message":"Not Found"}`, http.StatusNotFound)
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

			_, err := service.Invoice.DeleteItemCategory(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInvoiceItemCategory_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.InvoiceItemCategory
		want string
	}{
		{
			name: "InvoiceItemCategory with all fields",
			in: harvest.InvoiceItemCategory{
				ID:           harvest.Int64(1466293),
				Name:         harvest.String("Product"),
				UseAsService: harvest.Bool(false),
				UseAsExpense: harvest.Bool(true),
				CreatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
				UpdatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 30, 0, time.UTC)),
			},
			want: `harvest.InvoiceItemCategory{ID:1466293, Name:"Product", UseAsService:false, UseAsExpense:true, CreatedAt:time.Time{2017-06-26 20:41:00 +0000 UTC}, UpdatedAt:time.Time{2017-06-26 20:41:30 +0000 UTC}}`, //nolint: lll
		},
		{
			name: "InvoiceItemCategory with minimal fields",
			in: harvest.InvoiceItemCategory{
				ID:   harvest.Int64(1466292),
				Name: harvest.String("Service"),
			},
			want: `harvest.InvoiceItemCategory{ID:1466292, Name:"Service"}`,
		},
		{
			name: "InvoiceItemCategory with boolean fields",
			in: harvest.InvoiceItemCategory{
				ID:           harvest.Int64(1466294),
				Name:         harvest.String("Both"),
				UseAsService: harvest.Bool(true),
				UseAsExpense: harvest.Bool(true),
			},
			want: `harvest.InvoiceItemCategory{ID:1466294, Name:"Both", UseAsService:true, UseAsExpense:true}`,
		},
		{
			name: "InvoiceItemCategory without timestamps",
			in: harvest.InvoiceItemCategory{
				ID:           harvest.Int64(1466295),
				Name:         harvest.String("NoTimestamps"),
				UseAsService: harvest.Bool(false),
				UseAsExpense: harvest.Bool(false),
			},
			want: `harvest.InvoiceItemCategory{ID:1466295, Name:"NoTimestamps", UseAsService:false, UseAsExpense:false}`,
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

func TestInvoiceItemCategoryList_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.InvoiceItemCategoryList
		want string
	}{
		{
			name: "InvoiceItemCategoryList with multiple categories",
			in: harvest.InvoiceItemCategoryList{
				InvoiceItemCategories: []*harvest.InvoiceItemCategory{
					{
						ID:           harvest.Int64(1466293),
						Name:         harvest.String("Product"),
						UseAsService: harvest.Bool(false),
						UseAsExpense: harvest.Bool(true),
						CreatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
						UpdatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
					},
					{
						ID:           harvest.Int64(1466292),
						Name:         harvest.String("Service"),
						UseAsService: harvest.Bool(true),
						UseAsExpense: harvest.Bool(false),
						CreatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
						UpdatedAt:    harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(2000),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.InvoiceItemCategoryList{InvoiceItemCategories:[harvest.InvoiceItemCategory{ID:1466293, Name:"Product", UseAsService:false, UseAsExpense:true, CreatedAt:time.Time{2017-06-26 20:41:00 +0000 UTC}, UpdatedAt:time.Time{2017-06-26 20:41:00 +0000 UTC}} harvest.InvoiceItemCategory{ID:1466292, Name:"Service", UseAsService:true, UseAsExpense:false, CreatedAt:time.Time{2017-06-26 20:41:00 +0000 UTC}, UpdatedAt:time.Time{2017-06-26 20:41:00 +0000 UTC}}], Pagination:harvest.Pagination{PerPage:2000, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
		},
		{
			name: "InvoiceItemCategoryList with single category",
			in: harvest.InvoiceItemCategoryList{
				InvoiceItemCategories: []*harvest.InvoiceItemCategory{
					{
						ID:   harvest.Int64(1466293),
						Name: harvest.String("Product"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(2000),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.InvoiceItemCategoryList{InvoiceItemCategories:[harvest.InvoiceItemCategory{ID:1466293, Name:"Product"}], Pagination:harvest.Pagination{PerPage:2000, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
		},
		{
			name: "Empty InvoiceItemCategoryList",
			in: harvest.InvoiceItemCategoryList{
				InvoiceItemCategories: []*harvest.InvoiceItemCategory{},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(2000),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(0),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.InvoiceItemCategoryList{InvoiceItemCategories:[], Pagination:harvest.Pagination{PerPage:2000, TotalPages:1, TotalEntries:0, Page:1}}`, //nolint: lll
		},
		{
			name: "InvoiceItemCategoryList with Links",
			in: harvest.InvoiceItemCategoryList{
				InvoiceItemCategories: []*harvest.InvoiceItemCategory{
					{
						ID:   harvest.Int64(1466293),
						Name: harvest.String("Product"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(3),
					TotalEntries: harvest.Int(250),
					Page:         harvest.Int(2),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/invoice_item_categories?page=1&per_page=100"),
						Next:     harvest.String("https://api.harvestapp.com/v2/invoice_item_categories?page=3&per_page=100"),
						Previous: harvest.String("https://api.harvestapp.com/v2/invoice_item_categories?page=1&per_page=100"),
						Last:     harvest.String("https://api.harvestapp.com/v2/invoice_item_categories?page=3&per_page=100"),
					},
				},
			},
			want: `harvest.InvoiceItemCategoryList{InvoiceItemCategories:[harvest.InvoiceItemCategory{ID:1466293, Name:"Product"}], Pagination:harvest.Pagination{PerPage:100, TotalPages:3, TotalEntries:250, Page:2, Links:harvest.PageLinks{First:"https://api.harvestapp.com/v2/invoice_item_categories?page=1&per_page=100", Next:"https://api.harvestapp.com/v2/invoice_item_categories?page=3&per_page=100", Previous:"https://api.harvestapp.com/v2/invoice_item_categories?page=1&per_page=100", Last:"https://api.harvestapp.com/v2/invoice_item_categories?page=3&per_page=100"}}}`, //nolint: lll
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
