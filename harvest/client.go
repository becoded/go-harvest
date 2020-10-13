package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// ClientService handles communication with the client related
// methods of the Harvest API.
//
// Harvest API docs: https://help.getharvest.com/api-v2/clients-api/clients/clients/
type ClientService service

type Client struct {
	Id        *int64     `json:"id,omitempty"`         // Unique ID for the client.
	Name      *string    `json:"name,omitempty"`       // A textual description of the client.
	IsActive  *bool      `json:"is_active,omitempty"`  // Whether the client is active or archived.
	Address   *string    `json:"address,omitempty"`    // The physical address for the client.
	Currency  *string    `json:"currency,omitempty"`   // The currency code associated with this client.
	CreatedAt *time.Time `json:"created_at,omitempty"` // Date and time the client was created.
	UpdatedAt *time.Time `json:"updated_at,omitempty"` // Date and time the client was last updated.
}

type ClientList struct {
	Clients []*Client `json:"clients"`

	Pagination
}

func (c Client) String() string {
	return Stringify(c)
}

func (c ClientList) String() string {
	return Stringify(c)
}

type ClientListOptions struct {
	// Pass true to only return active projects and false to return inactive projects.
	IsActive bool `url:"is_active,omitempty"`
	// Only return projects that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

type ClientCreateRequest struct {
	Name     *string `json:"name"`                // required	A textual description of the client.
	IsActive *bool   `json:"is_active,omitempty"` // optional	Whether the client is active, or archived. Defaults to true.
	Address  *string `json:"address,omitempty"`   // optional	A textual representation of the client’s physical address. May include new line characters.
	Currency *string `json:"currency,omitempty"`  // optional	The currency used by the client. If not provided, the company’s currency will be used. See a list of supported currencies
}

type ClientUpdateRequest struct {
	Name     *string `json:"name,omitempty"`      // A textual description of the client.
	IsActive *bool   `json:"is_active,omitempty"` // Whether the client is active, or archived. Defaults to true.
	Address  *string `json:"address,omitempty"`   // A textual representation of the client’s physical address. May include new line characters.
	Currency *string `json:"currency,omitempty"`  // The currency used by the client. If not provided, the company’s currency will be used. See a list of supported currencies
}

// ListClients returns a list of your clients. The clients are returned sorted by creation date, with the most recently created clients appearing first.
func (s *ClientService) List(ctx context.Context, opt *ClientListOptions) (*ClientList, *http.Response, error) {
	u := "clients"
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	clientList := new(ClientList)
	resp, err := s.client.Do(ctx, req, &clientList)
	if err != nil {
		return nil, resp, err
	}

	return clientList, resp, nil
}

// GetClient retrieves the client with the given ID.
// Returns a client object and a 200 OK response code if a valid identifier was provided.
func (s *ClientService) Get(ctx context.Context, clientId int64) (*Client, *http.Response, error) {
	u := fmt.Sprintf("clients/%d", clientId)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	client := new(Client)
	resp, err := s.client.Do(ctx, req, &client)
	if err != nil {
		return nil, resp, err
	}

	return client, resp, nil
}

// CreateClient creates a new client object.
// Returns a client object and a 201 Created response code if the call succeeded.
func (s *ClientService) Create(ctx context.Context, data *ClientCreateRequest) (*Client, *http.Response, error) {
	u := "clients"

	req, err := s.client.NewRequest("POST", u, data)
	if err != nil {
		return nil, nil, err
	}

	client := new(Client)
	resp, err := s.client.Do(ctx, req, client)
	if err != nil {
		return nil, resp, err
	}

	return client, resp, nil
}

// UpdateClient updates the specific client by setting the values of the parameters passed.
// Any parameters not provided will be left unchanged. Returns a client object and a 200 OK response code if the call succeeded.
func (s *ClientService) Update(ctx context.Context, clientId int64, data *ClientUpdateRequest) (*Client, *http.Response, error) {
	u := fmt.Sprintf("clients/%d", clientId)

	req, err := s.client.NewRequest("PATCH", u, data)
	if err != nil {
		return nil, nil, err
	}

	client := new(Client)
	resp, err := s.client.Do(ctx, req, client)
	if err != nil {
		return nil, resp, err
	}

	return client, resp, nil
}

// DeleteClient deletes a specific client. Deleting a client is only possible if it has no projects, invoices, or estimates associated with it.
// Returns a 200 OK response code if the call succeeded.
func (s *ClientService) Delete(ctx context.Context, clientId int64) (*http.Response, error) {
	u := fmt.Sprintf("clients/%d", clientId)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
