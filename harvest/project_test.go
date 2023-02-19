package harvest_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/becoded/go-harvest/harvest"
	"github.com/stretchr/testify/assert"
)

func TestProjectService_List(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		testBody(t, r, "project/list/body_1.json")
		testWriteResponse(t, w, "project/list/response_1.json")
	})

	projectList, _, err := service.Project.List(context.Background(), &harvest.ProjectListOptions{})
	assert.NoError(t, err)

	fmt.Println(projectList)

	createdOne := time.Date(
		2017, 6, 26, 21, 52, 18, 0, time.UTC)
	updatedOne := time.Date(
		2017, 6, 26, 21, 54, 6, 0, time.UTC)
	startOnOne := time.Date(
		2017, 6, 1, 0, 0, 0, 0, time.UTC)

	createdTwo := time.Date(
		2017, 6, 26, 21, 36, 23, 0, time.UTC)
	updatedTwo := time.Date(
		2017, 6, 26, 21, 54, 46, 0, time.UTC)
	startOnTwo := time.Date(
		2017, 1, 1, 0, 0, 0, 0, time.UTC)
	endOnTwo := time.Date(
		2017, 3, 31, 0, 0, 0, 0, time.UTC)

	want := &harvest.ProjectList{
		Projects: []*harvest.Project{
			{
				ID: harvest.Int64(14308069),
				Client: &harvest.Client{
					ID:       harvest.Int64(5735776),
					Name:     harvest.String("123 Industries"),
					Currency: harvest.String("EUR"),
				},
				Name:                             harvest.String("Online Store - Phase 1"),
				Code:                             harvest.String("OS1"),
				IsActive:                         harvest.Bool(true),
				IsBillable:                       harvest.Bool(true),
				IsFixedFee:                       harvest.Bool(false),
				BillBy:                           harvest.String("Project"),
				HourlyRate:                       harvest.Float64(100),
				Budget:                           harvest.Float64(200),
				BudgetBy:                         harvest.String("project"),
				BudgetIsMonthly:                  harvest.Bool(false),
				NotifyWhenOverBudget:             harvest.Bool(true),
				OverBudgetNotificationPercentage: harvest.Float64(80),
				ShowBudgetToAll:                  harvest.Bool(false),
				CostBudgetIncludeExpenses:        harvest.Bool(false),
				Notes:                            harvest.String(""),
				StartsOn:                         &harvest.Date{Time: startOnOne},
				CreatedAt:                        &createdOne,
				UpdatedAt:                        &updatedOne,
			},
			{
				ID: harvest.Int64(14307913),
				Client: &harvest.Client{
					ID:       harvest.Int64(5735774),
					Name:     harvest.String("ABC Corp"),
					Currency: harvest.String("USD"),
				},
				Name:                             harvest.String("Marketing Website"),
				Code:                             harvest.String("MW"),
				IsActive:                         harvest.Bool(true),
				IsBillable:                       harvest.Bool(true),
				IsFixedFee:                       harvest.Bool(false),
				BillBy:                           harvest.String("Project"),
				HourlyRate:                       harvest.Float64(100),
				Budget:                           harvest.Float64(50),
				BudgetBy:                         harvest.String("project"),
				BudgetIsMonthly:                  harvest.Bool(false),
				NotifyWhenOverBudget:             harvest.Bool(true),
				OverBudgetNotificationPercentage: harvest.Float64(80),
				ShowBudgetToAll:                  harvest.Bool(false),
				CostBudgetIncludeExpenses:        harvest.Bool(false),
				Notes:                            harvest.String(""),
				StartsOn:                         &harvest.Date{Time: startOnTwo},
				EndsOn:                           &harvest.Date{Time: endOnTwo},
				CreatedAt:                        &createdTwo,
				UpdatedAt:                        &updatedTwo,
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
				First:    harvest.String("https://api.harvestapp.com/v2/projects?page=1&per_page=100"),
				Next:     nil,
				Previous: nil,
				Last:     harvest.String("https://api.harvestapp.com/v2/projects?page=1&per_page=100"),
			},
		},
	}

	assert.Equal(t, want, projectList)
}

func TestProjectService_Get(t *testing.T) {
	t.Parallel()

	service, mux, teardown := setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/projects/14308069", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		testBody(t, r, "project/get/body_1.json")
		testWriteResponse(t, w, "project/get/response_1.json")
	})

	project, _, err := service.Project.Get(context.Background(), 14308069)
	assert.NoError(t, err)

	createdOne := time.Date(
		2017, 6, 26, 21, 52, 18, 0, time.UTC)
	updatedOne := time.Date(
		2017, 6, 26, 21, 54, 6, 0, time.UTC)
	startOn := time.Date(
		2017, 6, 1, 0, 0, 0, 0, time.UTC)

	want := &harvest.Project{
		ID: harvest.Int64(14308069),
		Client: &harvest.Client{
			ID:       harvest.Int64(5735776),
			Name:     harvest.String("123 Industries"),
			Currency: harvest.String("EUR"),
		},
		Name:                             harvest.String("Online Store - Phase 1"),
		Code:                             harvest.String("OS1"),
		IsActive:                         harvest.Bool(true),
		IsBillable:                       harvest.Bool(true),
		IsFixedFee:                       harvest.Bool(false),
		BillBy:                           harvest.String("Project"),
		HourlyRate:                       harvest.Float64(100),
		Budget:                           harvest.Float64(200),
		BudgetBy:                         harvest.String("project"),
		BudgetIsMonthly:                  harvest.Bool(false),
		NotifyWhenOverBudget:             harvest.Bool(true),
		OverBudgetNotificationPercentage: harvest.Float64(80),
		ShowBudgetToAll:                  harvest.Bool(false),
		CostBudgetIncludeExpenses:        harvest.Bool(false),
		Notes:                            harvest.String(""),
		StartsOn:                         &harvest.Date{Time: startOn},
		CreatedAt:                        &createdOne,
		UpdatedAt:                        &updatedOne,
	}

	assert.Equal(t, want, project)
}
