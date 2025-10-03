package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
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
			t.Parallel()

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
			t.Parallel()

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

func TestClientService_CreateClientContact(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2019, 6, 26, 21, 44, 57, 0, time.UTC)
	updatedAt := time.Date(2019, 6, 26, 21, 44, 57, 0, time.UTC)

	tests := []struct {
		name      string
		input     *harvest.ClientContactCreateRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.ClientContact
		wantErr   bool
	}{
		{
			name: "Valid Contact Creation",
			input: &harvest.ClientContactCreateRequest{
				ClientID:  harvest.Int64(5735776),
				FirstName: harvest.String("George"),
				LastName:  harvest.String("Frank"),
				Email:     harvest.String("georgefrank@example.com"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					testBody(t, r, "contact/create/body_1.json")
					testWriteResponse(t, w, "contact/create/response_1.json")
				})
			},
			want: &harvest.ClientContact{
				ID: harvest.Int64(4706510),
				Client: &harvest.Client{
					ID:   harvest.Int64(5735776),
					Name: harvest.String("123 Industries"),
				},
				FirstName:   harvest.String("George"),
				LastName:    harvest.String("Frank"),
				Email:       harvest.String("georgefrank@example.com"),
				PhoneOffice: harvest.String(""),
				PhoneMobile: harvest.String(""),
				Fax:         harvest.String(""),
				CreatedAt:   &createdAt,
				UpdatedAt:   &updatedAt,
			},
			wantErr: false,
		},
		{
			name: "Error Creating Contact",
			input: &harvest.ClientContactCreateRequest{
				ClientID:  harvest.Int64(5735776),
				FirstName: harvest.String("George"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "POST")
					http.Error(w, `{"message":"Invalid data"}`, http.StatusBadRequest)
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

			got, _, err := service.Client.CreateClientContact(context.Background(), tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestClientService_UpdateClientContact(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2019, 6, 26, 21, 44, 57, 0, time.UTC)
	updatedAt := time.Date(2019, 6, 26, 21, 44, 57, 0, time.UTC)

	tests := []struct {
		name      string
		contactID int64
		input     *harvest.ClientContactUpdateRequest
		setupMock func(mux *http.ServeMux)
		want      *harvest.ClientContact
		wantErr   bool
	}{
		{
			name:      "Valid Contact Update",
			contactID: 4706510,
			input: &harvest.ClientContactUpdateRequest{
				Title: harvest.String("Owner"),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/contacts/4706510", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					testBody(t, r, "contact/update/body_1.json")
					testWriteResponse(t, w, "contact/update/response_1.json")
				})
			},
			want: &harvest.ClientContact{
				ID: harvest.Int64(4706510),
				Client: &harvest.Client{
					ID:   harvest.Int64(5735776),
					Name: harvest.String("123 Industries"),
				},
				Title:       harvest.String("Owner"),
				FirstName:   harvest.String("George"),
				LastName:    harvest.String("Frank"),
				Email:       harvest.String("georgefrank@example.com"),
				PhoneOffice: harvest.String(""),
				PhoneMobile: harvest.String(""),
				Fax:         harvest.String(""),
				CreatedAt:   &createdAt,
				UpdatedAt:   &updatedAt,
			},
			wantErr: false,
		},
		{
			name:      "Error Updating Contact",
			contactID: 4706510,
			input: &harvest.ClientContactUpdateRequest{
				Title: harvest.String(""),
			},
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/contacts/4706510", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PATCH")
					http.Error(w, `{"message":"Invalid data"}`, http.StatusBadRequest)
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

			got, _, err := service.Client.UpdateClientContact(context.Background(), tt.contactID, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestClientService_DeleteClientContact(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		contactID int64
		setupMock func(mux *http.ServeMux)
		wantErr   bool
	}{
		{
			name:      "Valid Client Deletion",
			contactID: 1,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/contacts/1", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					testFormValues(t, r, values{})
					w.WriteHeader(http.StatusNoContent)
				})
			},
			wantErr: false,
		},
		{
			name:      "Client Not Found",
			contactID: 999,
			setupMock: func(mux *http.ServeMux) {
				mux.HandleFunc("/contacts/999", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
					http.Error(w, `{"message":"Client not found"}`, http.StatusNotFound)
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

			_, err := service.Client.DeleteClientContact(context.Background(), tt.contactID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClientContact_String(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2017, 6, 26, 21, 20, 7, 0, time.UTC)
	updatedAt := time.Date(2017, 6, 26, 21, 27, 7, 0, time.UTC)

	tests := []struct {
		name string
		in   harvest.ClientContact
		want string
	}{
		{
			name: "Contact with all fields",
			in: harvest.ClientContact{
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
				CreatedAt:   &createdAt,
				UpdatedAt:   &updatedAt,
			},
			want: `harvest.ClientContact{ID:4706479, Client:harvest.Client{ID:5735774, Name:"ABC Corp"}, Title:"Owner", FirstName:"Jane", LastName:"Doe", Email:"janedoe@example.com", PhoneOffice:"(203) 697-8885", PhoneMobile:"(203) 697-8886", Fax:"(203) 697-8887", CreatedAt:time.Time{2017-06-26 21:20:07 +0000 UTC}, UpdatedAt:time.Time{2017-06-26 21:27:07 +0000 UTC}}`, //nolint: lll
		},
		{
			name: "Contact with minimal fields",
			in: harvest.ClientContact{
				ID:        harvest.Int64(123456),
				FirstName: harvest.String("John"),
				LastName:  harvest.String("Smith"),
				Email:     harvest.String("john.smith@example.com"),
			},
			want: `harvest.ClientContact{ID:123456, FirstName:"John", LastName:"Smith", Email:"john.smith@example.com"}`,
		},
		{
			name: "Contact with client info only",
			in: harvest.ClientContact{
				ID: harvest.Int64(999),
				Client: &harvest.Client{
					ID:   harvest.Int64(100),
					Name: harvest.String("Test Company"),
				},
				FirstName: harvest.String("Test"),
			},
			want: `harvest.ClientContact{ID:999, Client:harvest.Client{ID:100, Name:"Test Company"}, FirstName:"Test"}`,
		},
		{
			name: "Contact with phone numbers only",
			in: harvest.ClientContact{
				ID:          harvest.Int64(555),
				FirstName:   harvest.String("Contact"),
				PhoneOffice: harvest.String("555-1234"),
				PhoneMobile: harvest.String("555-5678"),
				Fax:         harvest.String("555-9999"),
			},
			want: `harvest.ClientContact{ID:555, FirstName:"Contact", PhoneOffice:"555-1234", PhoneMobile:"555-5678", Fax:"555-9999"}`, //nolint: lll
		},
		{
			name: "Empty Contact",
			in:   harvest.ClientContact{},
			want: `harvest.ClientContact{}`,
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

func TestClientContactList_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   harvest.ClientContactList
		want string
	}{
		{
			name: "ContactList with multiple contacts",
			in: harvest.ClientContactList{
				ClientContacts: []*harvest.ClientContact{
					{
						ID:        harvest.Int64(4706479),
						FirstName: harvest.String("Jane"),
						LastName:  harvest.String("Doe"),
						Email:     harvest.String("jane@example.com"),
					},
					{
						ID:        harvest.Int64(4706453),
						FirstName: harvest.String("Richard"),
						LastName:  harvest.String("Roe"),
						Email:     harvest.String("richard@example.com"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.ClientContactList{ClientContacts:[harvest.ClientContact{ID:4706479, FirstName:"Jane", LastName:"Doe", Email:"jane@example.com"} harvest.ClientContact{ID:4706453, FirstName:"Richard", LastName:"Roe", Email:"richard@example.com"}], Pagination:harvest.Pagination{PerPage:100, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
		},
		{
			name: "ContactList with single contact",
			in: harvest.ClientContactList{
				ClientContacts: []*harvest.ClientContact{
					{
						ID:        harvest.Int64(999),
						FirstName: harvest.String("Solo"),
						LastName:  harvest.String("Contact"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(50),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(1),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.ClientContactList{ClientContacts:[harvest.ClientContact{ID:999, FirstName:"Solo", LastName:"Contact"}], Pagination:harvest.Pagination{PerPage:50, TotalPages:1, TotalEntries:1, Page:1}}`, //nolint: lll
		},
		{
			name: "Empty ContactList",
			in: harvest.ClientContactList{
				ClientContacts: []*harvest.ClientContact{},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(100),
					TotalPages:   harvest.Int(0),
					TotalEntries: harvest.Int(0),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.ClientContactList{ClientContacts:[], Pagination:harvest.Pagination{PerPage:100, TotalPages:0, TotalEntries:0, Page:1}}`, //nolint: lll
		},
		{
			name: "ContactList with Links",
			in: harvest.ClientContactList{
				ClientContacts: []*harvest.ClientContact{
					{
						ID:        harvest.Int64(100),
						FirstName: harvest.String("Test"),
						LastName:  harvest.String("User"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(25),
					TotalPages:   harvest.Int(3),
					TotalEntries: harvest.Int(75),
					Page:         harvest.Int(2),
					Links: &harvest.PageLinks{
						First:    harvest.String("https://api.harvestapp.com/v2/contacts?page=1&per_page=25"),
						Next:     harvest.String("https://api.harvestapp.com/v2/contacts?page=3&per_page=25"),
						Previous: harvest.String("https://api.harvestapp.com/v2/contacts?page=1&per_page=25"),
						Last:     harvest.String("https://api.harvestapp.com/v2/contacts?page=3&per_page=25"),
					},
				},
			},
			want: `harvest.ClientContactList{ClientContacts:[harvest.ClientContact{ID:100, FirstName:"Test", LastName:"User"}], Pagination:harvest.Pagination{PerPage:25, TotalPages:3, TotalEntries:75, Page:2, Links:harvest.PageLinks{First:"https://api.harvestapp.com/v2/contacts?page=1&per_page=25", Next:"https://api.harvestapp.com/v2/contacts?page=3&per_page=25", Previous:"https://api.harvestapp.com/v2/contacts?page=1&per_page=25", Last:"https://api.harvestapp.com/v2/contacts?page=3&per_page=25"}}}`, //nolint: lll
		},
		{
			name: "ContactList with contacts with client info",
			in: harvest.ClientContactList{
				ClientContacts: []*harvest.ClientContact{
					{
						ID: harvest.Int64(1),
						Client: &harvest.Client{
							ID:   harvest.Int64(10),
							Name: harvest.String("Company A"),
						},
						FirstName: harvest.String("Alice"),
					},
					{
						ID: harvest.Int64(2),
						Client: &harvest.Client{
							ID:   harvest.Int64(20),
							Name: harvest.String("Company B"),
						},
						FirstName: harvest.String("Bob"),
					},
				},
				Pagination: harvest.Pagination{
					PerPage:      harvest.Int(10),
					TotalPages:   harvest.Int(1),
					TotalEntries: harvest.Int(2),
					Page:         harvest.Int(1),
				},
			},
			want: `harvest.ClientContactList{ClientContacts:[harvest.ClientContact{ID:1, Client:harvest.Client{ID:10, Name:"Company A"}, FirstName:"Alice"} harvest.ClientContact{ID:2, Client:harvest.Client{ID:20, Name:"Company B"}, FirstName:"Bob"}], Pagination:harvest.Pagination{PerPage:10, TotalPages:1, TotalEntries:2, Page:1}}`, //nolint: lll
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
