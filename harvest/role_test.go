package harvest_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/becoded/go-harvest/harvest"
)

func TestRoleService_CreateRole(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/roles", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{})
		testBody(t, r, "role/create/body_1.json")
		testWriteResponse(t, w, "role/create/response_1.json")
	})

	role, _, err := service.Role.Create(context.Background(), &harvest.RoleCreateRequest{
		Name:    harvest.String("Role new"),
		UserIDs: harvest.Ints64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
	})
	assert.NoError(t, err)

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &harvest.Role{
		ID:        harvest.Int64(1),
		Name:      harvest.String("Role new"),
		UserIDs:   harvest.Ints64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		CreatedAt: &createdOne,
		UpdatedAt: &updatedOne,
	}

	assert.Equal(t, want, role)
}

func TestRoleService_DeleteRole(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/roles/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testFormValues(t, r, values{})
		testBody(t, r, "role/delete/body_1.json")
		testWriteResponse(t, w, "role/delete/response_1.json")
	})

	_, err := service.Role.Delete(context.Background(), 1)
	assert.NoError(t, err)
}

func TestRoleService_GetRole(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/roles/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		testBody(t, r, "role/get/body_1.json")
		testWriteResponse(t, w, "role/get/response_1.json")
	})

	role, _, err := service.Role.Get(context.Background(), 1)
	assert.NoError(t, err)

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &harvest.Role{
		ID:        harvest.Int64(1),
		Name:      harvest.String("Role 1"),
		UserIDs:   harvest.Ints64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		CreatedAt: &createdOne,
		UpdatedAt: &updatedOne,
	}

	assert.Equal(t, want, role)
}

func TestRoleService_ListRoles(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/roles", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		testBody(t, r, "role/list/body_1.json")
		testWriteResponse(t, w, "role/list/response_1.json")
	})

	role, _, err := service.Role.List(context.Background(), &harvest.RoleListOptions{})
	assert.NoError(t, err)

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)
	createdTwo := time.Date(
		2018, 3, 2, 10, 12, 13, 0, time.UTC)
	updatedTwo := time.Date(
		2018, 4, 30, 12, 13, 14, 0, time.UTC)

	want := &harvest.RoleList{
		Roles: []*harvest.Role{
			{
				ID:        harvest.Int64(1),
				Name:      harvest.String("Role 1"),
				UserIDs:   harvest.Ints64([]int64{1, 2, 3, 4, 5}),
				CreatedAt: &createdOne,
				UpdatedAt: &updatedOne,
			}, {
				ID:        harvest.Int64(2),
				Name:      harvest.String("Role 2"),
				UserIDs:   harvest.Ints64([]int64{6, 7, 8, 9, 10}),
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
				First:    harvest.String("https://api.harvestapp.com/v2/roles?page=1&per_page=100"),
				Next:     nil,
				Previous: nil,
				Last:     harvest.String("https://api.harvestapp.com/v2/roles?page=1&per_page=100"),
			},
		},
	}

	assert.Equal(t, want, role)
}

func TestRoleService_UpdateRole(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/roles/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		testBody(t, r, "role/update/body_1.json")
		testWriteResponse(t, w, "role/update/response_1.json")
	})

	role, _, err := service.Role.Update(context.Background(), 1, &harvest.RoleUpdateRequest{
		Name:    harvest.String("Role update"),
		UserIDs: harvest.Ints64([]int64{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}),
	})
	assert.NoError(t, err)

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &harvest.Role{
		ID:        harvest.Int64(1),
		Name:      harvest.String("Role update"),
		UserIDs:   harvest.Ints64([]int64{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}),
		CreatedAt: &createdOne,
		UpdatedAt: &updatedOne,
	}

	assert.Equal(t, want, role)
}
