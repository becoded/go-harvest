package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// RoleService handles communication with the issue related
// methods of the Harvest API.
//
// Harvest API docs: https://help.getharvest.com/api-v2/roles-api/roles/roles/
type RoleService service

type Role struct {
	Id        *int64     `json:"id,omitempty"`         // Unique ID for the role.
	Name      *string    `json:"name,omitempty"`       // The name of the role.
	UserIds   *[]int64   `json:"user_ids,omitempty"`   // of integers	The IDs of the users assigned to this role.
	CreatedAt *time.Time `json:"created_at,omitempty"` // Date and time the role was created.
	UpdatedAt *time.Time `json:"updated_at,omitempty"` // Date and time the role was last updated.
}

type RoleList struct {
	Roles []*Role `json:"roles"`

	Pagination
}

func (p Role) String() string {
	return Stringify(p)
}

func (p RoleList) String() string {
	return Stringify(p)
}

type RoleListOptions struct {
	ListOptions
}

type RoleCreateRequest struct {
	Name    *string  `json:"name"`               // required	The name of the role.
	UserIds *[]int64 `json:"user_ids,omitempty"` // The IDs of the users assigned to this role.
}

type RoleUpdateRequest struct {
	Name    *string  `json:"name,omitempty"`     // The name of the role.
	UserIds *[]int64 `json:"user_ids,omitempty"` // The IDs of the users assigned to this role.
}

func (s *RoleService) ListRoles(ctx context.Context, opt *RoleListOptions) (*RoleList, *http.Response, error) {
	u := "roles"
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	roleList := new(RoleList)
	resp, err := s.client.Do(ctx, req, &roleList)
	if err != nil {
		return nil, resp, err
	}

	return roleList, resp, nil
}

func (s *RoleService) GetRole(ctx context.Context, roleId int64) (*Role, *http.Response, error) {
	u := fmt.Sprintf("roles/%d", roleId)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	role := new(Role)
	resp, err := s.client.Do(ctx, req, role)
	if err != nil {
		return nil, resp, err
	}

	return role, resp, nil
}

func (s *RoleService) CreateRole(ctx context.Context, data *RoleCreateRequest) (*Role, *http.Response, error) {
	u := "roles"

	req, err := s.client.NewRequest("POST", u, data)
	if err != nil {
		return nil, nil, err
	}

	role := new(Role)
	resp, err := s.client.Do(ctx, req, role)
	if err != nil {
		return nil, resp, err
	}

	return role, resp, nil
}

func (s *RoleService) UpdateRole(ctx context.Context, roleId int64, data *RoleUpdateRequest) (*Role, *http.Response, error) {
	u := fmt.Sprintf("roles/%d", roleId)

	req, err := s.client.NewRequest("PATCH", u, data)
	if err != nil {
		return nil, nil, err
	}

	role := new(Role)
	resp, err := s.client.Do(ctx, req, role)
	if err != nil {
		return nil, resp, err
	}

	return role, resp, nil
}

func (s *RoleService) DeleteRole(ctx context.Context, roleId int64) (*http.Response, error) {
	u := fmt.Sprintf("roles/%d", roleId)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
