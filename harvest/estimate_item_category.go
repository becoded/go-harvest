package harvest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Harvest API docs: https://help.getharvest.com/api-v2/estimates-api/estimates/estimate-item-categories/

type EstimateItemCategory struct {
	// Unique ID for the estimate item category.
	ID *int64 `json:"id,omitempty"`
	// The name of the estimate item category.
	Name *string `json:"name,omitempty"`
	// Date and time the estimate item category was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Date and time the estimate item category was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type EstimateItemCategoryList struct {
	EstimateItemCategories []*EstimateItemCategory `json:"estimate_item_categories"`

	Pagination
}

type EstimateItemCategoryRequest struct {
	Name *string `json:"name,omitempty"` // required	The name of the estimate item category.
}

func (p EstimateItemCategory) String() string {
	return Stringify(p)
}

func (p EstimateItemCategoryList) String() string {
	return Stringify(p)
}

type EstimateItemCategoryListOptions struct {
	// Only return estimate item categories that have been updated since the given date and time.
	UpdatedSince time.Time `url:"updated_since,omitempty"`

	ListOptions
}

func (s *EstimateService) ListItemCategories(
	ctx context.Context,
	opt *EstimateItemCategoryListOptions,
) (*EstimateItemCategoryList, *http.Response, error) {
	u := "estimate_item_categories"

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	estimateItemCategoryList := new(EstimateItemCategoryList)

	resp, err := s.client.Do(ctx, req, &estimateItemCategoryList)
	if err != nil {
		return nil, resp, err
	}

	return estimateItemCategoryList, resp, nil
}

func (s *EstimateService) GetItemCategory(
	ctx context.Context,
	estimateItemCategoryID int64,
) (*EstimateItemCategory, *http.Response, error) {
	u := fmt.Sprintf("estimate_item_categories/%d", estimateItemCategoryID)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	estimateItemCategory := new(EstimateItemCategory)

	resp, err := s.client.Do(ctx, req, estimateItemCategory)
	if err != nil {
		return nil, resp, err
	}

	return estimateItemCategory, resp, nil
}

func (s *EstimateService) CreateItemCategory(
	ctx context.Context,
	data *EstimateItemCategoryRequest,
) (*EstimateItemCategory, *http.Response, error) {
	u := "estimate_item_categories"

	req, err := s.client.NewRequest(ctx, "POST", u, data)
	if err != nil {
		return nil, nil, err
	}

	estimateItemCategory := new(EstimateItemCategory)

	resp, err := s.client.Do(ctx, req, estimateItemCategory)
	if err != nil {
		return nil, resp, err
	}

	return estimateItemCategory, resp, nil
}

func (s *EstimateService) UpdateItemCategory(
	ctx context.Context,
	estimateItemCategoryID int64,
	data *EstimateItemCategoryRequest,
) (*EstimateItemCategory, *http.Response, error) {
	u := fmt.Sprintf("estimate_item_categories/%d", estimateItemCategoryID)

	req, err := s.client.NewRequest(ctx, "PATCH", u, data)
	if err != nil {
		return nil, nil, err
	}

	estimateItemCategory := new(EstimateItemCategory)

	resp, err := s.client.Do(ctx, req, estimateItemCategory)
	if err != nil {
		return nil, resp, err
	}

	return estimateItemCategory, resp, nil
}

func (s *EstimateService) DeleteItemCategory(
	ctx context.Context,
	estimateItemCategoryID int64,
) (*http.Response, error) {
	u := fmt.Sprintf("estimate_item_categories/%d", estimateItemCategoryID)

	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
