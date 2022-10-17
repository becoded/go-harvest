package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/becoded/go-harvest/harvest"
	"github.com/stretchr/testify/assert"
)

func TestClientService_ListContacts(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		testBody(t, r, "contact/list/body_1.json")
		testWriteResponse(t, w, "contact/list/response_1.json")
	})

	clientContactList, _, err := service.Client.ListContacts(context.Background(), &harvest.ClientContactListOptions{})
	assert.NoError(t, err)

	createdOne := time.Date(
		2017, 6, 26, 21, 20, 7, 0, time.UTC)
	updatedOne := time.Date(
		2017, 6, 26, 21, 27, 7, 0, time.UTC)
	createdTwo := time.Date(
		2017, 6, 26, 21, 0o6, 55, 0, time.UTC)
	updatedTwo := time.Date(
		2017, 6, 26, 21, 27, 20, 0, time.UTC)

	want := &harvest.ClientContactList{
		ClientContacts: []*harvest.ClientContact{
			{
				ID: harvest.Int64(4706479),
				Client: &harvest.Client{
					ID:   harvest.Int64(5735774),
					Name: harvest.String("ABC Corp"),
				},
				Title:       harvest.String("Owner"),
				FirstName:   harvest.String("Jane"),
				LastName:    harvest.String("Doe"),
				Email:       harvest.String("janedoe@example.com"),
				PhoneOffice: harvest.String("(203) 697-8885"),
				PhoneMobile: harvest.String("(203) 697-8886"),
				Fax:         harvest.String("(203) 697-8887"),
				CreatedAt:   &createdOne,
				UpdatedAt:   &updatedOne,
			},
			{
				ID: harvest.Int64(4706453),
				Client: &harvest.Client{
					ID:   harvest.Int64(5735776),
					Name: harvest.String("123 Industries"),
				},
				Title:       harvest.String("Manager"),
				FirstName:   harvest.String("Richard"),
				LastName:    harvest.String("Roe"),
				Email:       harvest.String("richardroe@example.com"),
				PhoneOffice: harvest.String("(318) 515-5905"),
				PhoneMobile: harvest.String("(318) 515-5906"),
				Fax:         harvest.String("(318) 515-5907"),
				CreatedAt:   &createdTwo,
				UpdatedAt:   &updatedTwo,
			},
		},
		Pagination: harvest.Pagination{
			PerPage:      harvest.Int(100),
			TotalPages:   harvest.Int(1),
			TotalEntries: harvest.Int(2),
			Page:         harvest.Int(1),
			Links: &harvest.PageLinks{
				First:    harvest.String("https://api.harvestapp.com/v2/contacts?page=1&per_page=100"),
				Next:     nil,
				Previous: nil,
				Last:     harvest.String("https://api.harvestapp.com/v2/contacts?page=1&per_page=100"),
			},
		},
	}

	assert.Equal(t, want, clientContactList)
}

func TestClientService_GetContact(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/contacts/4706479", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		testBody(t, r, "contact/get/body_1.json")
		testWriteResponse(t, w, "contact/get/response_1.json")
	})

	client, _, err := service.Client.GetContact(context.Background(), 4706479)
	assert.NoError(t, err)

	createdOne := time.Date(
		2017, 6, 26, 21, 20, 0o7, 0, time.UTC)
	updatedOne := time.Date(
		2017, 6, 26, 21, 27, 0o7, 0, time.UTC)

	want := &harvest.ClientContact{
		ID: harvest.Int64(4706479),
		Client: &harvest.Client{
			ID:   harvest.Int64(5735774),
			Name: harvest.String("ABC Corp"),
		},
		Title:       harvest.String("Owner"),
		FirstName:   harvest.String("Jane"),
		LastName:    harvest.String("Doe"),
		Email:       harvest.String("janedoe@example.com"),
		PhoneOffice: harvest.String("(203) 697-8885"),
		PhoneMobile: harvest.String("(203) 697-8886"),
		Fax:         harvest.String("(203) 697-8887"),
		CreatedAt:   &createdOne,
		UpdatedAt:   &updatedOne,
	}

	assert.Equal(t, want, client)
}

func TestClientService_CreateClientContact(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{})
		testBody(t, r, "contact/create/body_1.json")
		testWriteResponse(t, w, "contact/create/response_1.json")
	})

	client, _, err := service.Client.CreateClientContact(
		context.Background(),
		&harvest.ClientContactCreateRequest{
			ClientID:  harvest.Int64(5735776),
			FirstName: harvest.String("George"),
			LastName:  harvest.String("Frank"),
			Email:     harvest.String("georgefrank@example.com"),
		},
	)
	assert.NoError(t, err)

	createdOne := time.Date(
		2019, 6, 26, 21, 44, 57, 0, time.UTC)
	updatedOne := time.Date(
		2019, 6, 26, 21, 44, 57, 0, time.UTC)

	want := &harvest.ClientContact{
		ID: harvest.Int64(4706510),
		Client: &harvest.Client{
			ID:   harvest.Int64(5735776),
			Name: harvest.String("123 Industries"),
		},
		FirstName:   harvest.String("George"),
		LastName:    harvest.String("Frank"),
		Email:       harvest.String("georgefrank@example.com"),
		PhoneMobile: harvest.String(""),
		PhoneOffice: harvest.String(""),
		Fax:         harvest.String(""),
		CreatedAt:   &createdOne,
		UpdatedAt:   &updatedOne,
	}

	assert.Equal(t, want, client)
}

func TestClientService_UpdateClientContact(t *testing.T) {
	// contacts/%d", contactID
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/contacts/4706510", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		testBody(t, r, "contact/update/body_1.json")
		testWriteResponse(t, w, "contact/update/response_1.json")
	})

	client, _, err := service.Client.UpdateClientContact(
		context.Background(),
		4706510,
		&harvest.ClientContactUpdateRequest{
			Title: harvest.String("Owner"),
		},
	)
	if err != nil {
		t.Errorf("UpdateClient returned error: %v", err)
	}

	createdOne := time.Date(
		2019, 6, 26, 21, 44, 57, 0, time.UTC)
	updatedOne := time.Date(
		2019, 6, 26, 21, 44, 57, 0, time.UTC)

	want := &harvest.ClientContact{
		ID: harvest.Int64(4706510),
		Client: &harvest.Client{
			ID:   harvest.Int64(5735776),
			Name: harvest.String("123 Industries"),
		},
		Title:       harvest.String("Owner"),
		FirstName:   harvest.String("George"),
		LastName:    harvest.String("Frank"),
		Email:       harvest.String("georgefrank@example.com"),
		PhoneMobile: harvest.String(""),
		PhoneOffice: harvest.String(""),
		Fax:         harvest.String(""),
		CreatedAt:   &createdOne,
		UpdatedAt:   &updatedOne,
	}

	assert.Equal(t, want, client)
}

func TestClientService_DeleteClientContact(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/contacts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testFormValues(t, r, values{})
		testBody(t, r, "contact/delete/body_1.json")
		testWriteResponse(t, w, "contact/delete/response_1.json")
	})

	_, err := service.Client.DeleteClientContact(context.Background(), 1)
	assert.NoError(t, err)
}
