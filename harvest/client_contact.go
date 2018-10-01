package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Harvest API docs: https://help.getharvest.com/api-v2/clients-api/clients/contacts/

type ClientContact struct {
	Id          *int64     `json:"id,omitempty"`           // Unique ID for the contact.
	Client      *Client    `json:"client,omitempty"`       // An object containing the contact’s client id and name.
	Title       *string    `json:"title,omitempty"`        // The title of the contact.
	FirstName   *string    `json:"first_name,omitempty"`   // The first name of the contact.
	LastName    *string    `json:"last_name,omitempty"`    // The last name of the contact.
	Email       *string    `json:"email,omitempty"`        // The contact’s email address.
	PhoneOffice *string    `json:"phone_office,omitempty"` // The contact’s office phone number.
	PhoneMobile *string    `json:"phone_mobile,omitempty"` // The contact’s mobile phone number.
	Fax         *string    `json:"fax,omitempty"`          // The contact’s fax number.
	CreatedAt   *time.Time `json:"created_at,omitempty"`   // Date and time the contact was created.
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`   // Date and time the contact was last updated.
}

type ClientContactList struct {
	ClientContacts []*ClientContact `json:"contacts"`

	Pagination
}

func (p ClientContact) String() string {
	return Stringify(p)
}

func (p ClientContactList) String() string {
	return Stringify(p)
}

type ClientContactListOptions struct {
	// Only return contacts belonging to the client with the given ID.
	ClientId int64 `url:"client_id,omitempty"`
	// Only return contacts that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

type ClientContactCreateRequest struct {
	ClientId    *int64  `json:"client_id"`              // required	The ID of the client associated with this contact.
	Title       *string `json:"title,omitempty"`        // optional	The title of the contact.
	FirstName   *string `json:"first_name"`             // required	The first name of the contact.
	LastName    *string `json:"last_name,omitempty"`    // optional	The last name of the contact.
	Email       *string `json:"email,omitempty"`        // optional	The contact’s email address.
	PhoneOffice *string `json:"phone_office,omitempty"` // optional	The contact’s office phone number.
	PhoneMobile *string `json:"phone_mobile,omitempty"` // optional	The contact’s mobile phone number.
	Fax         *string `json:"fax,omitempty"`          // optional	The contact’s fax number.
}

type ClientContactUpdateRequest struct {
	ClientId    *int64  `json:"client_id,omitempty"`    // required	The ID of the client associated with this contact.
	Title       *string `json:"title,omitempty"`        // optional	The title of the contact.
	FirstName   *string `json:"first_name,omitempty"`   // required	The first name of the contact.
	LastName    *string `json:"last_name,omitempty"`    // optional	The last name of the contact.
	Email       *string `json:"email,omitempty"`        // optional	The contact’s email address.
	PhoneOffice *string `json:"phone_office,omitempty"` // optional	The contact’s office phone number.
	PhoneMobile *string `json:"phone_mobile,omitempty"` // optional	The contact’s mobile phone number.
	Fax         *string `json:"fax,omitempty"`          // optional	The contact’s fax number.
}

func (s *ClientService) ListContacts(ctx context.Context, opt *ClientContactListOptions) (*ClientContactList, *http.Response, error) {
	u := "contacts"
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	clientContactList := new(ClientContactList)
	resp, err := s.client.Do(ctx, req, &clientContactList)
	if err != nil {
		return nil, resp, err
	}

	return clientContactList, resp, nil
}

func (s *ClientService) GetContact(ctx context.Context, clientContactId int64) (*ClientContact, *http.Response, error) {
	u := fmt.Sprintf("contacts/%d", clientContactId)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	clientContact := new(ClientContact)
	resp, err := s.client.Do(ctx, req, clientContact)
	if err != nil {
		return nil, resp, err
	}

	return clientContact, resp, nil
}

func (s *ClientService) CreateClientContact(ctx context.Context, data *ClientContactCreateRequest) (*ClientContact, *http.Response, error) {
	u := "contacts"

	req, err := s.client.NewRequest("POST", u, data)
	if err != nil {
		return nil, nil, err
	}

	clientContact := new(ClientContact)
	resp, err := s.client.Do(ctx, req, clientContact)
	if err != nil {
		return nil, resp, err
	}

	return clientContact, resp, nil
}

func (s *ClientService) UpdateClientContact(ctx context.Context, contactId int64, data *ClientContactUpdateRequest) (*ClientContact, *http.Response, error) {
	u := fmt.Sprintf("contacts/%d", contactId)

	req, err := s.client.NewRequest("PATCH", u, data)
	if err != nil {
		return nil, nil, err
	}

	clientContact := new(ClientContact)
	resp, err := s.client.Do(ctx, req, clientContact)
	if err != nil {
		return nil, resp, err
	}

	return clientContact, resp, nil
}

func (s *ClientService) DeleteClientContact(ctx context.Context, contactId int64) (*http.Response, error) {
	u := fmt.Sprintf("contacts/%d", contactId)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
