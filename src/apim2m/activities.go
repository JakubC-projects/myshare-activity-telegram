package apim2m

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
)

type Filter[T any] struct {
	Eq  *T `json:"_eq,omitempty"`
	Neq *T `json:"_ne,omitempty"`
	Gt  *T `json:"_gt,omitempty"`
	Lt  *T `json:"_lt,omitempty"`
	Gte *T `json:"_gte,omitempty"`
	Lte *T `json:"_lte,omitempty"`
}

type ActivitiesFilter struct {
	LastChanged *Filter[time.Time] `json:"lastChanged,omitempty"`
	Created     *Filter[time.Time] `json:"created,omitempty"`
	Start       *Filter[time.Time] `json:"start,omitempty"`
	TeamId      *Filter[int]       `json:"teamID,omitempty"`
	And         []ActivitiesFilter `json:"_and,omitempty"`
	Or          []ActivitiesFilter `json:"_or,omitempty"`
}

type ActivitiesQueryParams struct {
	Filter    *ActivitiesFilter
	OrderBy   OrderByField
	SortOrder SortDirection
}

type OrderByField string

var (
	LastChanged OrderByField = "lastChanged"
	Start       OrderByField = "start"
)

type SortDirection string

var (
	Asc  SortDirection = "asc"
	Desc SortDirection = "desc"
)

func (qp ActivitiesQueryParams) ToQuery() string {
	res := url.Values{}
	if qp.OrderBy != "" {
		res.Add("orderBy", string(qp.OrderBy))
	}
	if qp.SortOrder != "" {
		res.Add("sortOrder", string(qp.SortOrder))
	}
	if qp.Filter != nil {
		filterData, err := json.Marshal(qp.Filter)
		if err != nil {
			panic(fmt.Errorf("invalid filter %+v: %w", qp.Filter, err))
		}
		res.Add("filter", string(filterData))
	}

	return res.Encode()
}

func GetActivities(ctx context.Context, qp ActivitiesQueryParams) ([]models.ContributionsActivity, error) {
	url := fmt.Sprintf("%s/Projects?%s", config.Get().ContributionsAPI.BaseUrl, qp.ToQuery())

	fmt.Println(url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: invalid status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read body: %w", err)
	}
	var res struct {
		Data []models.ContributionsActivity
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal result: %w", err)
	}
	return res.Data, nil
}
