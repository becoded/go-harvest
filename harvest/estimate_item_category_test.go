package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestEstimateService_ListItemCategories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(mux *http.ServeMux)
		want      *harvest.EstimateItemCategoryList
		wantErr   bool
	}{
		{
			name: "Valid List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimate_item_categories", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testWriteResponse(t, w, "estimate_item_category/list/response_1.json")
				})
			},
			want: &harvest.EstimateItemCategoryList{
				EstimateItemCategories: []*harvest.EstimateItemCategory{
					{
						ID:        harvest.Int64(1378704),
						Name:      harvest.String("Product"),
						CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
						UpdatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
					},
					{
						ID:        harvest.Int64(1378703),
						Name:      harvest.String("Service"),
						CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
						UpdatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
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
						First:    harvest.String("https://api.harvestapp.com/v2/estimate_item_categories?page=1&per_page=2000"),
						Next:     nil,
						Previous: nil,
						Last:     harvest.String("https://api.harvestapp.com/v2/estimate_item_categories?page=1&per_page=2000"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Error Fetching List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimate_item_categories", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Estimate.ListItemCategories(context.Background(), &harvest.EstimateItemCategoryListOptions{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestEstimateService_GetItemCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        int64
		setupMock func(mux *http.ServeMux)
		want      *harvest.EstimateItemCategory
		wantErr   bool
	}{
		{
			name: "Valid Get",
			id:   1378704,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimate_item_categories/1378704", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testWriteResponse(t, w, "estimate_item_category/get/response_1.json")
				})
			},
			want: &harvest.EstimateItemCategory{
				ID:        harvest.Int64(1378704),
				Name:      harvest.String("Product"),
				CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name: "Not Found",
			id:   999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimate_item_categories/999", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Estimate.GetItemCategory(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestEstimateService_CreateItemCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		request   *harvest.EstimateItemCategoryRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.EstimateItemCategory
		wantErr   bool
	}{
		{
			name: "Valid Create",
			request: &harvest.EstimateItemCategoryRequest{
				Name: harvest.String("Hosting"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimate_item_categories", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testBody(t, r, "estimate_item_category/create/body_1.json")
					testWriteResponse(t, w, "estimate_item_category/create/response_1.json")
				})
			},
			want: &harvest.EstimateItemCategory{
				ID:        harvest.Int64(1379244),
				Name:      harvest.String("Hosting"),
				CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 6, 35, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 6, 35, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name: "Error Creating Category",
			request: &harvest.EstimateItemCategoryRequest{
				Name: harvest.String("Hosting"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimate_item_categories", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Estimate.CreateItemCategory(context.Background(), tt.request)
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

func TestEstimateService_UpdateItemCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        int64
		request   *harvest.EstimateItemCategoryRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.EstimateItemCategory
		wantErr   bool
	}{
		{
			name: "Valid Update",
			id:   1379244,
			request: &harvest.EstimateItemCategoryRequest{
				Name: harvest.String("Transportation"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimate_item_categories/1379244", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					testBody(t, r, "estimate_item_category/update/body_1.json")
					testWriteResponse(t, w, "estimate_item_category/update/response_1.json")
				})
			},
			want: &harvest.EstimateItemCategory{
				ID:        harvest.Int64(1379244),
				Name:      harvest.String("Transportation"),
				CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 6, 35, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2017, 6, 27, 16, 7, 5, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name: "Error Updating Category",
			id:   1379244,
			request: &harvest.EstimateItemCategoryRequest{
				Name: harvest.String("Transportation"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimate_item_categories/1379244", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Estimate.UpdateItemCategory(context.Background(), tt.id, tt.request)
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

func TestEstimateService_DeleteItemCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        int64
		setupMock func(mux *http.ServeMux)
		wantErr   bool
	}{
		{
			name: "Valid Delete",
			id:   1,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimate_item_categories/1", func(w http.ResponseWriter, r *http.Request) {
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
				mux.HandleFunc("/estimate_item_categories/999", func(w http.ResponseWriter, r *http.Request) {
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

			_, err := service.Estimate.DeleteItemCategory(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEstimateItemCategory_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.EstimateItemCategory
		want string
	}{
		{
			name: "EstimateItemCategory with all fields",
			in: harvest.EstimateItemCategory{
				ID:        harvest.Int64(1378704),
				Name:      harvest.String("Product"),
				CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
				UpdatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 5, 0, time.UTC)),
			},
			want: `harvest.EstimateItemCategory{ID:1378704, Name:"Product", CreatedAt:time.Time{2017-06-26 20:41:00 +0000 UTC}, UpdatedAt:time.Time{2017-06-26 20:41:05 +0000 UTC}}`, //nolint: lll
		},
		{
			name: "EstimateItemCategory with minimal fields",
			in: harvest.EstimateItemCategory{
				ID:   harvest.Int64(999),
				Name: harvest.String("Service"),
			},
			want: `harvest.EstimateItemCategory{ID:999, Name:"Service"}`,
		},
		{
			name: "EstimateItemCategory with ID only",
			in: harvest.EstimateItemCategory{
				ID: harvest.Int64(123),
			},
			want: `harvest.EstimateItemCategory{ID:123}`,
		},
		{
			name: "Empty EstimateItemCategory",
			in:   harvest.EstimateItemCategory{},
			want: `harvest.EstimateItemCategory{}`,
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

func TestEstimateItemCategoryList_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.EstimateItemCategoryList
		want string
	}{
		{
			name: "EstimateItemCategoryList with multiple categories",
			in: harvest.EstimateItemCategoryList{
				EstimateItemCategories: []*harvest.EstimateItemCategory{
					{
						ID:        harvest.Int64(1378704),
						Name:      harvest.String("Product"),
						CreatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
						UpdatedAt: harvest.TimeTimeP(time.Date(2017, 6, 26, 20, 41, 0, 0, time.UTC)),
					},
					{
						ID:        harvest.Int64(1378703),
						Name:      harvest.String("Service"),
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
			want: `harvest.EstimateItemCategoryList{EstimateItemCategories:[harvest.EstimateItemCategory{ID:1378704, Name:"Product", CreatedAt:time.Time{2017-06-26 20:41:00 +0000 UTC}, UpdatedAt:time.Time{2017-06-26 20:41:00 +0000 UTC}} harvest.EstimateItemCategory{ID:1378703, Name:"Service", CreatedAt:time.Time{2017-06-26 20:41:00 +0000 UTC}, UpdatedAt:time.Time{2017-06-26 20:41:00 +0000 UTC}}], Pagination:harvest.Pagination{PerPage:2000, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
		},
		{
			name: "EstimateItemCategoryList with single category",
			in: harvest.EstimateItemCategoryList{
				EstimateItemCategories: []*harvest.EstimateItemCategory{
					{
						ID:   harvest.Int64(999),
						Name: harvest.String("Hosting"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.EstimateItemCategoryList{EstimateItemCategories:[harvest.EstimateItemCategory{ID:999, Name:"Hosting"}], Pagination:harvest.Pagination{PerPage:100, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
		},
		{
			name: "Empty EstimateItemCategoryList",
			in: harvest.EstimateItemCategoryList{
				EstimateItemCategories: []*harvest.EstimateItemCategory{},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(0),
					TotalEntries: harvest.Int(0),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.EstimateItemCategoryList{EstimateItemCategories:[], Pagination:harvest.Pagination{PerPage:100, TotalPages:0, TotalEntries:0, Page:1}}`, //nolint: lll
		},
		{
			name: "EstimateItemCategoryList with Links",
			in: harvest.EstimateItemCategoryList{
				EstimateItemCategories: []*harvest.EstimateItemCategory{
					{
						ID:   harvest.Int64(100),
						Name: harvest.String("Product"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(50),
					TotalPages:   harvest.Int(3),
					TotalEntries: harvest.Int(150),
					Page:         harvest.Int(2),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/estimate_item_categories?page=1&per_page=50"),
						Next:     harvest.String("https://api.harvestapp.com/v2/estimate_item_categories?page=3&per_page=50"),
						Previous: harvest.String("https://api.harvestapp.com/v2/estimate_item_categories?page=1&per_page=50"),
						Last:     harvest.String("https://api.harvestapp.com/v2/estimate_item_categories?page=3&per_page=50"),
					},
				},
			},
			want: `harvest.EstimateItemCategoryList{EstimateItemCategories:[harvest.EstimateItemCategory{ID:100, Name:"Product"}], Pagination:harvest.Pagination{PerPage:50, TotalPages:3, TotalEntries:150, Page:2, Links:harvest.PageLinks{First:"https://api.harvestapp.com/v2/estimate_item_categories?page=1&per_page=50", Next:"https://api.harvestapp.com/v2/estimate_item_categories?page=3&per_page=50", Previous:"https://api.harvestapp.com/v2/estimate_item_categories?page=1&per_page=50", Last:"https://api.harvestapp.com/v2/estimate_item_categories?page=3&per_page=50"}}}`, //nolint: lll
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
