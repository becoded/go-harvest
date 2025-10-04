package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestEstimateService_List(t *testing.T) {
	t.Parallel()

	issueOne := time.Date(
		2017, 6, 1, 0, 0, 0, 0, time.UTC)
	sentOne := time.Date(
		2017, 6, 27, 16, 11, 33, 0, time.UTC)
	createdOne := time.Date(
		2017, 6, 27, 16, 11, 24, 0, time.UTC)
	updatedOne := time.Date(
		2017, 6, 27, 16, 13, 56, 0, time.UTC)
	issueTwo := time.Date(
		2017, 1, 1, 0, 0, 0, 0, time.UTC)
	sentTwo := time.Date(
		2017, 6, 27, 16, 10, 30, 0, time.UTC)
	createdTwo := time.Date(
		2017, 6, 27, 16, 9, 33, 0, time.UTC)
	updatedTwo := time.Date(
		2017, 6, 27, 16, 12, 0o0, 0, time.UTC)
	acceptedTwo := time.Date(
		2017, 6, 27, 16, 10, 32, 0, time.UTC)

	lineItemsOne := []harvest.EstimateLineItem{
		{
			ID:          harvest.Int64(53334195),
			Kind:        harvest.String("Service"),
			Description: harvest.String("Phase 2 of the Online Store"),
			Quantity:    harvest.Int64(100),
			UnitPrice:   harvest.Float64(100.0),
			Amount:      harvest.Float64(10000),
			Taxed:       harvest.Bool(true),
			Taxed2:      harvest.Bool(true),
		},
	}

	lineItemsTwo := []harvest.EstimateLineItem{
		{
			ID:          harvest.Int64(57531966),
			Kind:        harvest.String("Service"),
			Description: harvest.String("Phase 1 of the Online Store"),
			Quantity:    harvest.Int64(1),
			UnitPrice:   harvest.Float64(20000),
			Amount:      harvest.Float64(20000),
			Taxed:       harvest.Bool(true),
			Taxed2:      harvest.Bool(false),
		},
	}

	tests := []struct {
		name      string
		setupMock func(mux *http.ServeMux)
		want      *harvest.EstimateList
		wantErr   bool
	}{
		{
			name: "Valid Estimate List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimates", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "estimate/list/body_1.json")
					testWriteResponse(t, w, "estimate/list/response_1.json")
				})
			},
			want: &harvest.EstimateList{
				Estimates: []*harvest.Estimate{
					{
						ID: harvest.Int64(1439818),
						Client: &harvest.Client{
							ID:   harvest.Int64(5735776),
							Name: harvest.String("123 Industries"),
						},
						LineItems: &lineItemsOne,
						Creator: &harvest.User{
							ID:   harvest.Int64(1782884),
							Name: harvest.String("Bob Powell"),
						},
						ClientKey:      harvest.String("13dc088aa7d51ec687f186b146730c3c75dc7423"),
						Number:         harvest.String("1001"),
						PurchaseOrder:  harvest.String("5678"),
						Amount:         harvest.Float64(9630.0),
						Tax:            harvest.Float64(5.0),
						TaxAmount:      harvest.Float64(450.0),
						Tax2:           harvest.Float64(2.0),
						Tax2Amount:     harvest.Float64(180.0),
						Discount:       harvest.Float64(10.0),
						DiscountAmount: harvest.Float64(1000.0),
						Subject:        harvest.String("Online Store - Phase 2"),
						Notes:          harvest.String("Some notes about the estimate"),
						Currency:       harvest.String("USD"),
						State:          harvest.String("sent"),
						IssueDate: &harvest.Date{
							Time: issueOne,
						},
						SentAt:    &sentOne,
						CreatedAt: &createdOne,
						UpdatedAt: &updatedOne,
					}, {
						ID: harvest.Int64(1439814),
						Client: &harvest.Client{
							ID:   harvest.Int64(5735776),
							Name: harvest.String("123 Industries"),
						},
						LineItems: &lineItemsTwo,
						Creator: &harvest.User{
							ID:   harvest.Int64(1782884),
							Name: harvest.String("Bob Powell"),
						},
						ClientKey:      harvest.String("a5ffaeb30c55776270fcd3992b70332d769f97e7"),
						Number:         harvest.String("1000"),
						PurchaseOrder:  harvest.String("1234"),
						Amount:         harvest.Float64(21000.0),
						Tax:            harvest.Float64(5.0),
						TaxAmount:      harvest.Float64(1000.0),
						Tax2Amount:     harvest.Float64(0.0),
						DiscountAmount: harvest.Float64(0.0),
						Subject:        harvest.String("Online Store - Phase 1"),
						Notes:          harvest.String("Some notes about the estimate"),
						Currency:       harvest.String("USD"),
						State:          harvest.String("accepted"),
						IssueDate: &harvest.Date{
							Time: issueTwo,
						},
						SentAt:     &sentTwo,
						CreatedAt:  &createdTwo,
						UpdatedAt:  &updatedTwo,
						AcceptedAt: &acceptedTwo,
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
						First:    harvest.String("https://api.harvestapp.com/v2/estimates?page=1&per_page=2000"),
						Next:     nil,
						Previous: nil,
						Last:     harvest.String("https://api.harvestapp.com/v2/estimates?page=1&per_page=2000"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Error Fetching Estimate List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimates", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Estimate.List(context.Background(), &harvest.EstimateListOptions{})

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

func TestEstimateService_Get(t *testing.T) {
	t.Parallel()

	issueOne := time.Date(
		2017, 6, 1, 0, 0, 0, 0, time.UTC)
	sentOne := time.Date(
		2017, 6, 27, 16, 11, 33, 0, time.UTC)
	createdOne := time.Date(
		2017, 6, 27, 16, 11, 24, 0, time.UTC)
	updatedOne := time.Date(
		2017, 6, 27, 16, 13, 56, 0, time.UTC)

	lineItemsOne := []harvest.EstimateLineItem{
		{
			ID:          harvest.Int64(53334195),
			Kind:        harvest.String("Service"),
			Description: harvest.String("Phase 2 of the Online Store"),
			Quantity:    harvest.Int64(100),
			UnitPrice:   harvest.Float64(100.0),
			Amount:      harvest.Float64(10000),
			Taxed:       harvest.Bool(true),
			Taxed2:      harvest.Bool(true),
		},
	}

	tests := []struct {
		name       string
		estimateID int64
		setupMock  func(mux *http.ServeMux)
		want       *harvest.Estimate
		wantErr    bool
	}{
		{
			name:       "Valid Estimate",
			estimateID: 1,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimates/1", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "estimate/get/body_1.json")
					testWriteResponse(t, w, "estimate/get/response_1.json")
				})
			},
			want: &harvest.Estimate{
				ID: harvest.Int64(1439818),
				Client: &harvest.Client{
					ID:   harvest.Int64(5735776),
					Name: harvest.String("123 Industries"),
				},
				LineItems: &lineItemsOne,
				Creator: &harvest.User{
					ID:   harvest.Int64(1782884),
					Name: harvest.String("Bob Powell"),
				},
				ClientKey:      harvest.String("13dc088aa7d51ec687f186b146730c3c75dc7423"),
				Number:         harvest.String("1001"),
				PurchaseOrder:  harvest.String("5678"),
				Amount:         harvest.Float64(9630.0),
				Tax:            harvest.Float64(5.0),
				TaxAmount:      harvest.Float64(450.0),
				Tax2:           harvest.Float64(2.0),
				Tax2Amount:     harvest.Float64(180.0),
				Discount:       harvest.Float64(10.0),
				DiscountAmount: harvest.Float64(1000.0),
				Subject:        harvest.String("Online Store - Phase 2"),
				Notes:          harvest.String("Some notes about the estimate"),
				Currency:       harvest.String("USD"),
				State:          harvest.String("sent"),
				IssueDate: &harvest.Date{
					Time: issueOne,
				},
				SentAt:    &sentOne,
				CreatedAt: &createdOne,
				UpdatedAt: &updatedOne,
			},
			wantErr: false,
		},
		{
			name:       "Estimate Not Found",
			estimateID: 999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/estimates/999", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Estimate.Get(context.Background(), tt.estimateID)

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

func TestEstimate_String(t *testing.T) {
	t.Parallel()

	issueDate := time.Date(2017, 6, 1, 0, 0, 0, 0, time.UTC)
	sentAt := time.Date(2017, 6, 27, 16, 11, 33, 0, time.UTC)
	createdAt := time.Date(2017, 6, 27, 16, 11, 24, 0, time.UTC)
	updatedAt := time.Date(2017, 6, 27, 16, 13, 56, 0, time.UTC)

	tests := []struct {
		name string
		in   harvest.Estimate
		want string
	}{
		{
			name: "Estimate with all fields",
			in: harvest.Estimate{
				ID: harvest.Int64(1439818),
				Client: &harvest.Client{
					ID:   harvest.Int64(5735776),
					Name: harvest.String("123 Industries"),
				},
				LineItems: &[]harvest.EstimateLineItem{
					{
						ID:          harvest.Int64(53334195),
						Kind:        harvest.String("Service"),
						Description: harvest.String("Phase 2 of the Online Store"),
						Quantity:    harvest.Int64(100),
						UnitPrice:   harvest.Float64(100.0),
						Amount:      harvest.Float64(10000),
						Taxed:       harvest.Bool(true),
						Taxed2:      harvest.Bool(true),
					},
				},
				Creator: &harvest.User{
					ID:   harvest.Int64(1782884),
					Name: harvest.String("Bob Powell"),
				},
				ClientKey:      harvest.String("13dc088aa7d51ec687f186b146730c3c75dc7423"),
				Number:         harvest.String("1001"),
				PurchaseOrder:  harvest.String("5678"),
				Amount:         harvest.Float64(9630.0),
				Tax:            harvest.Float64(5.0),
				TaxAmount:      harvest.Float64(450.0),
				Tax2:           harvest.Float64(2.0),
				Tax2Amount:     harvest.Float64(180.0),
				Discount:       harvest.Float64(10.0),
				DiscountAmount: harvest.Float64(1000.0),
				Subject:        harvest.String("Online Store - Phase 2"),
				Notes:          harvest.String("Some notes about the estimate"),
				Currency:       harvest.String("USD"),
				State:          harvest.String("sent"),
				IssueDate:      &harvest.Date{Time: issueDate},
				SentAt:         &sentAt,
				CreatedAt:      &createdAt,
				UpdatedAt:      &updatedAt,
			},
			want: `harvest.Estimate{ID:1439818, Client:harvest.Client{ID:5735776, Name:"123 Industries"}, LineItems:[harvest.EstimateLineItem{ID:53334195, Kind:"Service", Description:"Phase 2 of the Online Store", Quantity:100, UnitPrice:100, Amount:10000, Taxed:true, Taxed2:true}], Creator:harvest.User{ID:1782884, Name:"Bob Powell"}, ClientKey:"13dc088aa7d51ec687f186b146730c3c75dc7423", Number:"1001", PurchaseOrder:"5678", Amount:9630, Tax:5, TaxAmount:450, Tax2:2, Tax2Amount:180, Discount:10, DiscountAmount:1000, Subject:"Online Store - Phase 2", Notes:"Some notes about the estimate", Currency:"USD", State:"sent", IssueDate:harvest.Date{{2017-06-01 00:00:00 +0000 UTC}}, SentAt:time.Time{2017-06-27 16:11:33 +0000 UTC}, CreatedAt:time.Time{2017-06-27 16:11:24 +0000 UTC}, UpdatedAt:time.Time{2017-06-27 16:13:56 +0000 UTC}}`, //nolint: lll
		},
		{
			name: "Estimate with minimal fields",
			in: harvest.Estimate{
				ID:       harvest.Int64(1439818),
				Number:   harvest.String("1001"),
				Amount:   harvest.Float64(9630.0),
				Currency: harvest.String("USD"),
				State:    harvest.String("draft"),
			},
			want: `harvest.Estimate{ID:1439818, Number:"1001", Amount:9630, Currency:"USD", State:"draft"}`,
		},
		{
			name: "Estimate with client only",
			in: harvest.Estimate{
				ID: harvest.Int64(1439818),
				Client: &harvest.Client{
					ID:   harvest.Int64(5735776),
					Name: harvest.String("123 Industries"),
				},
				Number: harvest.String("1001"),
				State:  harvest.String("sent"),
			},
			want: `harvest.Estimate{ID:1439818, Client:harvest.Client{ID:5735776, Name:"123 Industries"}, Number:"1001", State:"sent"}`, //nolint: lll
		},
		{
			name: "Estimate with line items",
			in: harvest.Estimate{
				ID:     harvest.Int64(1439818),
				Number: harvest.String("1001"),
				LineItems: &[]harvest.EstimateLineItem{
					{
						ID:          harvest.Int64(53334195),
						Kind:        harvest.String("Service"),
						Description: harvest.String("Consulting"),
						Quantity:    harvest.Int64(10),
						UnitPrice:   harvest.Float64(150.0),
						Amount:      harvest.Float64(1500),
					},
					{
						ID:          harvest.Int64(53334196),
						Kind:        harvest.String("Product"),
						Description: harvest.String("Software License"),
						Quantity:    harvest.Int64(1),
						UnitPrice:   harvest.Float64(500.0),
						Amount:      harvest.Float64(500),
					},
				},
			},
			want: `harvest.Estimate{ID:1439818, LineItems:[harvest.EstimateLineItem{ID:53334195, Kind:"Service", Description:"Consulting", Quantity:10, UnitPrice:150, Amount:1500} harvest.EstimateLineItem{ID:53334196, Kind:"Product", Description:"Software License", Quantity:1, UnitPrice:500, Amount:500}], Number:"1001"}`, //nolint: lll
		},
		{
			name: "Empty Estimate",
			in:   harvest.Estimate{},
			want: `harvest.Estimate{}`,
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

func TestEstimateList_String(t *testing.T) {
	t.Parallel()

	issueDate := time.Date(2017, 6, 1, 0, 0, 0, 0, time.UTC)
	sentAt := time.Date(2017, 6, 27, 16, 11, 33, 0, time.UTC)
	createdAt := time.Date(2017, 6, 27, 16, 11, 24, 0, time.UTC)
	updatedAt := time.Date(2017, 6, 27, 16, 13, 56, 0, time.UTC)

	tests := []struct {
		name string
		in   harvest.EstimateList
		want string
	}{
		{
			name: "EstimateList with multiple estimates",
			in: harvest.EstimateList{
				Estimates: []*harvest.Estimate{
					{
						ID:     harvest.Int64(1439818),
						Number: harvest.String("1001"),
						Amount: harvest.Float64(9630.0),
						State:  harvest.String("sent"),
						Client: &harvest.Client{
							ID:   harvest.Int64(5735776),
							Name: harvest.String("123 Industries"),
						},
					},
					{
						ID:     harvest.Int64(1439814),
						Number: harvest.String("1000"),
						Amount: harvest.Float64(21000.0),
						State:  harvest.String("accepted"),
						Client: &harvest.Client{
							ID:   harvest.Int64(5735776),
							Name: harvest.String("123 Industries"),
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
			want: `harvest.EstimateList{Estimates:[harvest.Estimate{ID:1439818, Client:harvest.Client{ID:5735776, Name:"123 Industries"}, Number:"1001", Amount:9630, State:"sent"} harvest.Estimate{ID:1439814, Client:harvest.Client{ID:5735776, Name:"123 Industries"}, Number:"1000", Amount:21000, State:"accepted"}], Pagination:harvest.Pagination{PerPage:100, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
		},
		{
			name: "EstimateList with single estimate",
			in: harvest.EstimateList{
				Estimates: []*harvest.Estimate{
					{
						ID:     harvest.Int64(999),
						Number: harvest.String("TEST-001"),
						State:  harvest.String("draft"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(50),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.EstimateList{Estimates:[harvest.Estimate{ID:999, Number:"TEST-001", State:"draft"}], Pagination:harvest.Pagination{PerPage:50, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
		},
		{
			name: "Empty EstimateList",
			in: harvest.EstimateList{
				Estimates: []*harvest.Estimate{},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(0),
					TotalEntries: harvest.Int(0),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.EstimateList{Estimates:[], Pagination:harvest.Pagination{PerPage:100, TotalPages:0, TotalEntries:0, Page:1}}`, //nolint: lll
		},
		{
			name: "EstimateList with Links",
			in: harvest.EstimateList{
				Estimates: []*harvest.Estimate{
					{
						ID:    harvest.Int64(100),
						State: harvest.String("sent"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(25),
					TotalPages:   harvest.Int(3),
					TotalEntries: harvest.Int(75),
					Page:         harvest.Int(2),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/estimates?page=1&per_page=25"),
						Next:     harvest.String("https://api.harvestapp.com/v2/estimates?page=3&per_page=25"),
						Previous: harvest.String("https://api.harvestapp.com/v2/estimates?page=1&per_page=25"),
						Last:     harvest.String("https://api.harvestapp.com/v2/estimates?page=3&per_page=25"),
					},
				},
			},
			want: `harvest.EstimateList{Estimates:[harvest.Estimate{ID:100, State:"sent"}], Pagination:harvest.Pagination{PerPage:25, TotalPages:3, TotalEntries:75, Page:2, Links:harvest.PageLinks{First:"https://api.harvestapp.com/v2/estimates?page=1&per_page=25", Next:"https://api.harvestapp.com/v2/estimates?page=3&per_page=25", Previous:"https://api.harvestapp.com/v2/estimates?page=1&per_page=25", Last:"https://api.harvestapp.com/v2/estimates?page=3&per_page=25"}}}`, //nolint: lll
		},
		{
			name: "EstimateList with full details",
			in: harvest.EstimateList{
				Estimates: []*harvest.Estimate{
					{
						ID: harvest.Int64(1439818),
						Client: &harvest.Client{
							ID:   harvest.Int64(5735776),
							Name: harvest.String("123 Industries"),
						},
						Number:    harvest.String("1001"),
						Amount:    harvest.Float64(9630.0),
						Currency:  harvest.String("USD"),
						State:     harvest.String("sent"),
						IssueDate: &harvest.Date{Time: issueDate},
						SentAt:    &sentAt,
						CreatedAt: &createdAt,
						UpdatedAt: &updatedAt,
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(10),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.EstimateList{Estimates:[harvest.Estimate{ID:1439818, Client:harvest.Client{ID:5735776, Name:"123 Industries"}, Number:"1001", Amount:9630, Currency:"USD", State:"sent", IssueDate:harvest.Date{{2017-06-01 00:00:00 +0000 UTC}}, SentAt:time.Time{2017-06-27 16:11:33 +0000 UTC}, CreatedAt:time.Time{2017-06-27 16:11:24 +0000 UTC}, UpdatedAt:time.Time{2017-06-27 16:13:56 +0000 UTC}}], Pagination:harvest.Pagination{PerPage:10, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
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
