package harvest

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestRoleService_CreateRole(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/roles", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{})
		testBody(t, r, `{"name":"Role new","user_ids":[1,2,3,4,5,6,7,8,9,10]}`+"\n")
		fmt.Fprint(w, `{"id":1,"name":"Role new","user_ids":[1,2,3,4,5,6,7,8,9,10],"created_at":"2018-01-31T20:34:30Z","updated_at":"2018-05-31T21:34:30Z"}`)
	})

	roleList, _, err := service.Role.Create(context.Background(), &RoleCreateRequest{
		Name:    String("Role new"),
		UserIds: Ints64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
	})
	if err != nil {
		t.Errorf("CreateRole role returned error: %v", err)
	}

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &Role{
		Id:        Int64(1),
		Name:      String("Role new"),
		UserIds:   Ints64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		CreatedAt: &createdOne,
		UpdatedAt: &updatedOne,
	}

	if !reflect.DeepEqual(roleList, want) {
		t.Errorf("Role.GetRole returned %+v, want %+v", roleList, want)
	}
}

func TestRoleService_DeleteRole(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/roles/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testFormValues(t, r, values{})
		testBody(t, r, ``)
		fmt.Fprint(w, ``)
	})

	_, err := service.Role.Delete(context.Background(), 1)
	if err != nil {
		t.Errorf("DeleteRole role returned error: %v", err)
	}
}

func TestRoleService_GetRole(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/roles/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `{"id":1,"name":"Role 1","user_ids":[1,2,3,4,5,6,7,8,9,10],"created_at":"2018-01-31T20:34:30Z","updated_at":"2018-05-31T21:34:30Z"}`)
	})

	roleList, _, err := service.Role.Get(context.Background(), 1)
	if err != nil {
		t.Errorf("Role.GetRole returned error: %v", err)
	}

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &Role{
		Id:        Int64(1),
		Name:      String("Role 1"),
		UserIds:   Ints64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		CreatedAt: &createdOne,
		UpdatedAt: &updatedOne,
	}

	if !reflect.DeepEqual(roleList, want) {
		t.Errorf("Role.GetRole returned %+v, want %+v", roleList, want)
	}
}

func TestRoleService_ListRoles(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/roles", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `{"roles":[{"id":1,"name":"Role 1","user_ids":[1,2,3,4,5],"created_at":"2018-01-31T20:34:30Z","updated_at":"2018-05-31T21:34:30Z"},{"id":2,"name":"Role 2","user_ids":[6,7,8,9,10],"created_at":"2018-03-02T10:12:13Z","updated_at":"2018-04-30T12:13:14Z"}],"per_page":100,"total_pages":1,"total_entries":2,"next_page":null,"previous_page":null,"page":1,"links":{"first":"https://api.harvestapp.com/v2/roles?page=1&per_page=100","next":null,"previous":null,"last":"https://api.harvestapp.com/v2/roles?page=1&per_page=100"}}`)
	})

	roleList, _, err := service.Role.List(context.Background(), &RoleListOptions{})
	if err != nil {
		t.Errorf("Role.ListRoles returned error: %v", err)
	}

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)
	createdTwo := time.Date(
		2018, 3, 2, 10, 12, 13, 0, time.UTC)
	updatedTwo := time.Date(
		2018, 4, 30, 12, 13, 14, 0, time.UTC)

	want := &RoleList{
		Roles: []*Role{
			{
				Id:        Int64(1),
				Name:      String("Role 1"),
				UserIds:   Ints64([]int64{1, 2, 3, 4, 5}),
				CreatedAt: &createdOne,
				UpdatedAt: &updatedOne,
			}, {
				Id:        Int64(2),
				Name:      String("Role 2"),
				UserIds:   Ints64([]int64{6, 7, 8, 9, 10}),
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
				First:    String("https://api.harvestapp.com/v2/roles?page=1&per_page=100"),
				Next:     nil,
				Previous: nil,
				Last:     String("https://api.harvestapp.com/v2/roles?page=1&per_page=100"),
			},
		},
	}

	if !reflect.DeepEqual(roleList, want) {
		t.Errorf("Role.ListRoles returned %+v, want %+v", roleList, want)
	}
}

func TestRoleService_UpdateRole(t *testing.T) {
	service, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/roles/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, values{})
		testBody(t, r, `{"name":"Role update","user_ids":[11,12,13,14,15,16,17,18,19,20]}`+"\n")
		fmt.Fprint(w, `{"id":1,"name":"Role update","is_active":true,"user_ids":[11,12,13,14,15,16,17,18,19,20],"created_at":"2018-01-31T20:34:30Z","updated_at":"2018-05-31T21:34:30Z"}`)
	})

	roleList, _, err := service.Role.Update(context.Background(), 1, &RoleUpdateRequest{
		Name:    String("Role update"),
		UserIds: Ints64([]int64{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}),
	})
	if err != nil {
		t.Errorf("CreateRole role returned error: %v", err)
	}

	createdOne := time.Date(
		2018, 1, 31, 20, 34, 30, 0, time.UTC)
	updatedOne := time.Date(
		2018, 5, 31, 21, 34, 30, 0, time.UTC)

	want := &Role{
		Id:        Int64(1),
		Name:      String("Role update"),
		UserIds:   Ints64([]int64{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}),
		CreatedAt: &createdOne,
		UpdatedAt: &updatedOne,
	}

	if !reflect.DeepEqual(roleList, want) {
		t.Errorf("Role.UpdateRole returned %+v, want %+v", roleList, want)
	}
}
