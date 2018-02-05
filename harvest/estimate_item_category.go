package harvest

import (
	"context"
	"fmt"
	"time"
	"net/http"
)

// Harvest API docs: https://help.getharvest.com/api-v2/estimates-api/estimates/estimate-item-categories/

type EstimateItemCategory struct {
	Id *int64 `json:"id,omitempty"` // Unique ID for the estimate item category.
	Name *string `json:"name,omitempty"` // The name of the estimate item category.
	CreatedAt *time.Time `json:"created_at,omitempty"` // Date and time the estimate item category was created.
	UpdatedAt *time.Time `json:"updated_at,omitempty"` // Date and time the estimate item category was last updated.
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
	UpdatedSince	time.Time `url:"updated_since,omitempty"`

	ListOptions
}

func (s *EstimateService) ListItemCategories(ctx context.Context, opt *EstimateItemCategoryListOptions) (*EstimateItemCategoryList, *http.Response, error) {
	u := "estimate_item_categories"
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
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

func (s *EstimateService) GetItemCategory(ctx context.Context, estimateItemCategoryId int64) (*EstimateItemCategory, *http.Response, error) {
	u := fmt.Sprintf("estimate_item_categories/%d", estimateItemCategoryId)
	req, err := s.client.NewRequest("GET", u, nil)
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

func (s *EstimateService) CreateItemCategory(ctx context.Context, data *EstimateItemCategoryRequest) (*EstimateItemCategory, *http.Response, error) {
	u := "estimate_item_categories"

	req, err := s.client.NewRequest("POST", u, data)
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

func (s *EstimateService) UpdateItemCategory(ctx context.Context, estimateItemCategoryId int64, data *EstimateItemCategoryRequest) (*EstimateItemCategory, *http.Response, error) {
	u := fmt.Sprintf("estimate_item_categories/%d", estimateItemCategoryId)

	req, err := s.client.NewRequest("PATCH", u, data)
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

func (s *EstimateService) DeleteItemCategory(ctx context.Context, estimateItemCategoryId int64) (*http.Response, error) {
	u := fmt.Sprintf("estimate_item_categories/%d", estimateItemCategoryId)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
