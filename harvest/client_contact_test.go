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

	createdOne := time.Date(
		2017, 6, 26, 21, 20, 7, 0, time.UTC)
	updatedOne := time.Date(
		2017, 6, 26, 21, 27, 7, 0, time.UTC)
	createdTwo := time.Date(
		2017, 6, 26, 21, 0o6, 55, 0, time.UTC)
	updatedTwo := time.Date(
		2017, 6, 26, 21, 27, 20, 0, time.UTC)

	tests := []struct {
		name      string
		setupMock func(mux *http.ServeMux)
		want      *harvest.ClientContactList
		wantErr   bool
	}{
		{
			name: "Valid Contact List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "contact/list/body_1.json")
					testWriteResponse(t, w, "contact/list/response_1.json")
				})
			},
			want: &harvest.ClientContactList{
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
			},
			wantErr: false,
		},
		{
			name: "Error Fetching Contact List",
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
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

			got, _, err := service.Client.ListContacts(context.Background(), &harvest.ClientContactListOptions{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestClientService_GetContact(t *testing.T) {
	t.Parallel()

	createdOne := time.Date(
		2017, 6, 26, 21, 20, 0o7, 0, time.UTC)
	updatedOne := time.Date(
		2017, 6, 26, 21, 27, 0o7, 0, time.UTC)

	tests := []struct {
		name      string
		contactID int64
		setupMock func(mux *http.ServeMux)
		want      *harvest.ClientContact
		wantErr   bool
	}{
		{
			name:      "Valid Contact Retrieval",
			contactID: 4706479,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/contacts/4706479", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{})
					testBody(t, r, "contact/get/body_1.json")
					testWriteResponse(t, w, "contact/get/response_1.json")
				})
			},
			want: &harvest.ClientContact{
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
			wantErr: false,
		},
		{
			name:      "Contact Not Found",
			contactID: 999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/contacts/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					http.Error(w, `{"message":"Contact not found"}`, http.StatusNotFound)
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

			got, _, err := service.Client.GetContact(context.Background(), tt.contactID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// Similar table-driven test structures can be applied to `CreateClientContact`, `UpdateClientContact`, and `DeleteClientContact`.
