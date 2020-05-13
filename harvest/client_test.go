package harvest

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestClientService_ListClients(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `{"clients":[{"id":1,"name":"Client 1","is_active":true,"address":"Address line 1","statement_key":"1234567890","created_at":"2018-01-31T20:34:30Z","updated_at":"2018-05-31T21:34:30Z","currency":"EUR"},{"id":2,"name":"Client 2","is_active":false,"address":"Address line 2","statement_key":"0987654321","created_at":"2018-03-02T10:12:13Z","updated_at":"2018-04-30T12:13:14Z","currency":"EUR"}],"per_page":100,"total_pages":1,"total_entries":2,"next_page":null,"previous_page":null,"page":1,"links":{"first":"https://api.harvestapp.com/v2/clients?page=1&per_page=100","next":null,"previous":null,"last":"https://api.harvestapp.com/v2/clients?page=1&per_page=100"}}`)
	})

	clientList, _, err := service.Client.ListClients(context.Background(), &ClientListOptions{})
	if err != nil {
		t.Errorf("Client.ListClient returned error: %v", err)
	}

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)
	createdTwo := time.Date(
		2018, 3, 2, 10, 12, 13, 0, time.UTC)
	updatedTwo := time.Date(
		2018, 4, 30, 12, 13, 14, 0, time.UTC)

	want := &ClientList{
		Clients: []*Client{
			{
				Id:        Int64(1),
				Name:      String("Client 1"),
				IsActive:  Bool(true),
				Address:   String("Address line 1"),
				Currency:  String("EUR"),
				CreatedAt: &createdOne,
				UpdatedAt: &updatedOne,
			}, {
				Id:        Int64(2),
				Name:      String("Client 2"),
				IsActive:  Bool(false),
				Address:   String("Address line 2"),
				Currency:  String("EUR"),
				CreatedAt: &createdTwo,
				UpdatedAt: &updatedTwo,
			}},
		Pagination: Pagination{
			PerPage:      Int(100),
			TotalPages:   Int(1),
			TotalEntries: Int(2),
			NextPage:     nil,
			PreviousPage: nil,
			Page:         Int(1),
			Links: &PageLinks{
				First:    String("https://api.harvestapp.com/v2/clients?page=1&per_page=100"),
				Next:     nil,
				Previous: nil,
				Last:     String("https://api.harvestapp.com/v2/clients?page=1&per_page=100"),
			},
		},
	}

	if !reflect.DeepEqual(clientList, want) {
		t.Errorf("Client.ListClients returned %+v, want %+v", clientList, want)
	}
}

func TestClientService_GetClient(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/clients/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `{"id":1,"name":"Client 1","is_active":true,"address":"Address line 1","statement_key":"1234567890","created_at":"2018-01-31T20:34:30Z","updated_at":"2018-05-31T21:34:30Z","currency":"EUR"}`)
	})

	clientList, _, err := service.Client.GetClient(context.Background(), 1)
	if err != nil {
		t.Errorf("Client.GetClient returned error: %v", err)
	}

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &Client{
		Id:        Int64(1),
		Name:      String("Client 1"),
		IsActive:  Bool(true),
		Address:   String("Address line 1"),
		Currency:  String("EUR"),
		CreatedAt: &createdOne,
		UpdatedAt: &updatedOne,
	}

	if !reflect.DeepEqual(clientList, want) {
		t.Errorf("Client.GetClient returned %+v, want %+v", clientList, want)
	}
}

func TestClientService_CreateClient(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{})
		testBody(t, r, `{"name":"Client new","is_active":true,"address":"Address line 1","currency":"EUR"}`+"\n")
		fmt.Fprint(w, `{"id":1,"name":"Client 1","is_active":true,"address":"Address line 1","statement_key":"1234567890","created_at":"2018-01-31T20:34:30Z","updated_at":"2018-05-31T21:34:30Z","currency":"EUR"}`)
	})

	clientList, _, err := service.Client.CreateClient(context.Background(), &ClientCreateRequest{
		Name:     String("Client new"),
		IsActive: Bool(true),
		Address:  String("Address line 1"),
		Currency: String("EUR"),
	})
	if err != nil {
		t.Errorf("CreateClient returned error: %v", err)
	}

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &Client{
		Id:        Int64(1),
		Name:      String("Client 1"),
		IsActive:  Bool(true),
		Address:   String("Address line 1"),
		Currency:  String("EUR"),
		CreatedAt: &createdOne,
		UpdatedAt: &updatedOne,
	}

	if !reflect.DeepEqual(clientList, want) {
		t.Errorf("Client.CreateClient returned %+v, want %+v", clientList, want)
	}
}

func TestClientService_UpdateClient(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/clients/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		testBody(t, r, `{"name":"Client new","is_active":true,"address":"Address line 1","currency":"EUR"}`+"\n")
		fmt.Fprint(w, `{"id":1,"name":"Client 1","is_active":true,"address":"Address line 1","statement_key":"1234567890","created_at":"2018-01-31T20:34:30Z","updated_at":"2018-05-31T21:34:30Z","currency":"EUR"}`)
	})

	clientList, _, err := service.Client.UpdateClient(context.Background(), 1, &ClientUpdateRequest{
		Name:     String("Client new"),
		IsActive: Bool(true),
		Address:  String("Address line 1"),
		Currency: String("EUR"),
	})
	if err != nil {
		t.Errorf("UpdateClient returned error: %v", err)
	}

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &Client{
		Id:        Int64(1),
		Name:      String("Client 1"),
		IsActive:  Bool(true),
		Address:   String("Address line 1"),
		Currency:  String("EUR"),
		CreatedAt: &createdOne,
		UpdatedAt: &updatedOne,
	}

	if !reflect.DeepEqual(clientList, want) {
		t.Errorf("Client.UpdateClient returned %+v, want %+v", clientList, want)
	}
}

func TestClientService_DeleteClient(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/clients/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testFormValues(t, r, values{})
		testBody(t, r, ``)
		fmt.Fprint(w, ``)
	})

	_, err := service.Client.DeleteClient(context.Background(), 1)
	if err != nil {
		t.Errorf("DeleteClient returned error: %v", err)
	}
}
