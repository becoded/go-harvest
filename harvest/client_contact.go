package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Harvest API docs: https://help.getharvest.com/api-v2/clients-api/clients/contacts/

type ClientContact struct {
	// Unique ID for the contact.
	ID *int64 `json:"id,omitempty"`
	// An object containing the contact’s client id and name.
	Client *Client `json:"client,omitempty"`
	// The title of the contact.
	Title *string `json:"title,omitempty"`
	// The first name of the contact.
	FirstName *string `json:"first_name,omitempty"`
	// The last name of the contact.
	LastName *string `json:"last_name,omitempty"`
	// The contact’s email address.
	Email *string `json:"email,omitempty"`
	// The contact’s office phone number.
	PhoneOffice *string `json:"phone_office,omitempty"`
	// The contact’s mobile phone number.
	PhoneMobile *string `json:"phone_mobile,omitempty"`
	// The contact’s fax number.
	Fax *string `json:"fax,omitempty"`
	// Date and time the contact was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Date and time the contact was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
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
	ClientID int64 `url:"client_id,omitempty"`
	// Only return contacts that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

type ClientContactCreateRequest struct {
	// required	The ID of the client associated with this contact.
	ClientID *int64 `json:"client_id"`
	// optional	The title of the contact.
	Title *string `json:"title,omitempty"`
	// required	The first name of the contact.
	FirstName *string `json:"first_name"`
	// optional	The last name of the contact.
	LastName *string `json:"last_name,omitempty"`
	// optional	The contact’s email address.
	Email *string `json:"email,omitempty"`
	// optional	The contact’s office phone number.
	PhoneOffice *string `json:"phone_office,omitempty"`
	// optional	The contact’s mobile phone number.
	PhoneMobile *string `json:"phone_mobile,omitempty"`
	// optional	The contact’s fax number.
	Fax *string `json:"fax,omitempty"`
}

type ClientContactUpdateRequest struct {
	// required	The ID of the client associated with this contact.
	ClientID *int64 `json:"client_id,omitempty"`
	// optional	The title of the contact.
	Title *string `json:"title,omitempty"`
	// required	The first name of the contact.
	FirstName *string `json:"first_name,omitempty"`
	// optional	The last name of the contact.
	LastName *string `json:"last_name,omitempty"`
	// optional	The contact’s email address.
	Email *string `json:"email,omitempty"`
	// optional	The contact’s office phone number.
	PhoneOffice *string `json:"phone_office,omitempty"`
	// optional	The contact’s mobile phone number.
	PhoneMobile *string `json:"phone_mobile,omitempty"`
	// optional	The contact’s fax number.
	Fax *string `json:"fax,omitempty"`
}

func (s *ClientService) ListContacts(
	ctx context.Context,
	opt *ClientContactListOptions,
) (*ClientContactList, *http.Response, error) {
	u := "contacts"

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	clientContactList := new(ClientContactList)

	resp, err := s.client.Do(ctx, req, clientContactList)
	if err != nil {
		return nil, resp, err
	}

	return clientContactList, resp, nil
}

func (s *ClientService) GetContact(ctx context.Context, clientContactID int64) (*ClientContact, *http.Response, error) {
	u := fmt.Sprintf("contacts/%d", clientContactID)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
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

func (s *ClientService) CreateClientContact(
	ctx context.Context,
	data *ClientContactCreateRequest,
) (*ClientContact, *http.Response, error) {
	u := "contacts"

	req, err := s.client.NewRequest(ctx, "POST", u, data)
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

func (s *ClientService) UpdateClientContact(
	ctx context.Context,
	contactID int64,
	data *ClientContactUpdateRequest,
) (*ClientContact, *http.Response, error) {
	u := fmt.Sprintf("contacts/%d", contactID)

	req, err := s.client.NewRequest(ctx, "PATCH", u, data)
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

func (s *ClientService) DeleteClientContact(ctx context.Context, contactID int64) (*http.Response, error) {
	u := fmt.Sprintf("contacts/%d", contactID)

	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
