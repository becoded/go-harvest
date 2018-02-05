package harvest


import (
"context"
"fmt"
"time"
	"net/http"
)

// RetainerService handles communication with the issue related
// methods of the Harvest API.
//
// Harvest API docs: https://help.getharvest.com/api-v2/retainers-api/retainers/retainers/
type RetainerService service


type Retainer struct {
	Id *int64 `json:"id,omitempty"` //Unique ID for the retainer.
	Name *string `json:"name,omitempty"` //Unique name for the retainer.
}

type RetainerList struct {
	Retainers []*Retainer `json:"retainers"`

	Pagination
}

func (p Retainer) String() string {
	return Stringify(p)
}

func (p RetainerList) String() string {
	return Stringify(p)
}

type RetainerListOptions struct {
	// Pass true to only return active retainers and false to return inactive retainers.
	IsActive	bool `url:"is_active,omitempty"`
	// Only return retainers that have been updated since the given date and time.
	UpdatedSince	time.Time `url:"updated_since,omitempty"`

	ListOptions
}

func (s *RetainerService) List(ctx context.Context, opt *RetainerListOptions) (*RetainerList, *http.Response, error) {
	u := "retainers"
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	retainerList := new(RetainerList)
	resp, err := s.client.Do(ctx, req, &retainerList)
	if err != nil {
		return nil, resp, err
	}

	return retainerList, resp, nil
}

func (s *RetainerService) Get(ctx context.Context, retainerId int64) (*Retainer, *http.Response, error) {
	u := fmt.Sprintf("retainers/%d", retainerId)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	retainer := new(Retainer)
	resp, err := s.client.Do(ctx, req, retainer)
	if err != nil {
		return nil, resp, err
	}

	return retainer, resp, nil
}