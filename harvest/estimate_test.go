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

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/estimates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		testBody(t, r, "estimate/list/body_1.json")
		testWriteResponse(t, w, "estimate/list/response_1.json")
	})

	clientList, _, err := service.Estimate.List(context.Background(), &harvest.EstimateListOptions{})
	assert.NoError(t, err)

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

	want := &harvest.EstimateList{
		Estimates: []*harvest.Estimate{
			{
				ID: harvest.Int64(1439818),
				Client: &harvest.Client{
					ID:   harvest.Int64(5735776),
					Name: harvest.String("123 Industries"),
				},
				LineItems: &lineItemsOne,
				Creator: &harvest.User{
					ID: harvest.Int64(1782884),
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
					ID: harvest.Int64(1782884),
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
	}

	assert.Equal(t, want, clientList)
}

func TestEstimateService_Get(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/estimates/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		testBody(t, r, "estimate/get/body_1.json")
		testWriteResponse(t, w, "estimate/get/response_1.json")
	})

	client, _, err := service.Estimate.Get(context.Background(), 1)
	assert.NoError(t, err)

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

	want := &harvest.Estimate{
		ID: harvest.Int64(1439818),
		Client: &harvest.Client{
			ID:   harvest.Int64(5735776),
			Name: harvest.String("123 Industries"),
		},
		LineItems: &lineItemsOne,
		Creator: &harvest.User{
			ID: harvest.Int64(1782884),
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
	}

	assert.Equal(t, want, client)
}
