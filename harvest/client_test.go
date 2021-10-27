package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/becoded/go-harvest/harvest"
	"github.com/stretchr/testify/assert"
)

func TestClientService_List(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		testBody(t, r, "client/list/body_1.json")
		testWriteResponse(t, w, "client/list/response_1.json")
	})

	clientList, _, err := service.Client.List(context.Background(), &harvest.ClientListOptions{})
	assert.NoError(t, err)

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)
	createdTwo := time.Date(
		2018, 3, 2, 10, 12, 13, 0, time.UTC)
	updatedTwo := time.Date(
		2018, 4, 30, 12, 13, 14, 0, time.UTC)

	want := &harvest.ClientList{
		Clients: []*harvest.Client{
			{
				ID:        harvest.Int64(1),
				Name:      harvest.String("Client 1"),
				IsActive:  harvest.Bool(true),
				Address:   harvest.String("Address line 1"),
				Currency:  harvest.String("EUR"),
				CreatedAt: &createdOne,
				UpdatedAt: &updatedOne,
			}, {
				ID:        harvest.Int64(2),
				Name:      harvest.String("Client 2"),
				IsActive:  harvest.Bool(false),
				Address:   harvest.String("Address line 2"),
				Currency:  harvest.String("EUR"),
				CreatedAt: &createdTwo,
				UpdatedAt: &updatedTwo,
			},
		},
		Pagination: harvest.Pagination{
			PerPage:      harvest.Int(100),
			TotalPages:   harvest.Int(1),
			TotalEntries: harvest.Int(2),
			NextPage:     nil,
			PreviousPage: nil,
			Page:         harvest.Int(1),
			Links: &harvest.PageLinks{
				First:    harvest.String("https://api.harvestapp.com/v2/clients?page=1&per_page=100"),
				Next:     nil,
				Previous: nil,
				Last:     harvest.String("https://api.harvestapp.com/v2/clients?page=1&per_page=100"),
			},
		},
	}

	assert.Equal(t, want, clientList)
}

func TestClientService_Get(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/clients/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		testBody(t, r, "client/get/body_1.json")
		testWriteResponse(t, w, "client/get/response_1.json")
	})

	client, _, err := service.Client.Get(context.Background(), 1)
	assert.NoError(t, err)

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &harvest.Client{
		ID:        harvest.Int64(1),
		Name:      harvest.String("Client 1"),
		IsActive:  harvest.Bool(true),
		Address:   harvest.String("Address line 1"),
		Currency:  harvest.String("EUR"),
		CreatedAt: &createdOne,
		UpdatedAt: &updatedOne,
	}

	assert.Equal(t, want, client)
}

func TestClientService_CreateClient(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{})
		testBody(t, r, "client/create/body_1.json")
		testWriteResponse(t, w, "client/create/response_1.json")
	})

	client, _, err := service.Client.Create(context.Background(), &harvest.ClientCreateRequest{
		Name:     harvest.String("Client new"),
		IsActive: harvest.Bool(true),
		Address:  harvest.String("Address line 1"),
		Currency: harvest.String("EUR"),
	})
	assert.NoError(t, err)

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &harvest.Client{
		ID:        harvest.Int64(1),
		Name:      harvest.String("Client 1"),
		IsActive:  harvest.Bool(true),
		Address:   harvest.String("Address line 1"),
		Currency:  harvest.String("EUR"),
		CreatedAt: &createdOne,
		UpdatedAt: &updatedOne,
	}

	assert.Equal(t, want, client)
}

func TestClientService_UpdateClient(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/clients/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		testBody(t, r, "client/update/body_1.json")
		testWriteResponse(t, w, "client/update/response_1.json")
	})

	client, _, err := service.Client.Update(context.Background(), 1, &harvest.ClientUpdateRequest{
		Name:     harvest.String("Client new"),
		IsActive: harvest.Bool(true),
		Address:  harvest.String("Address line 1"),
		Currency: harvest.String("EUR"),
	})
	if err != nil {
		t.Errorf("UpdateClient returned error: %v", err)
	}

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &harvest.Client{
		ID:        harvest.Int64(1),
		Name:      harvest.String("Client 1"),
		IsActive:  harvest.Bool(true),
		Address:   harvest.String("Address line 1"),
		Currency:  harvest.String("EUR"),
		CreatedAt: &createdOne,
		UpdatedAt: &updatedOne,
	}

	assert.Equal(t, want, client)
}

func TestClientService_DeleteClient(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/clients/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testFormValues(t, r, values{})
		testBody(t, r, "client/delete/body_1.json")
		testWriteResponse(t, w, "client/delete/response_1.json")
	})

	_, err := service.Client.Delete(context.Background(), 1)
	assert.NoError(t, err)
}
