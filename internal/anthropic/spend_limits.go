package anthropic

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// SpendLimit returned from Set/Get. Scope is polymorphic but Set only
// accepts scope.type="user"; Get/List can return other scope types
// (seat_tier, rbac_group, organization_service, organization) which are
// configured in claude.ai and not via this API.
type SpendLimit struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Amount    string          `json:"amount"`
	Currency  string          `json:"currency"`
	Period    string          `json:"period"`
	Scope     SpendLimitScope `json:"scope"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

// SpendLimitScope flattens the polymorphic scope union. The active field
// depends on Type: "user" → UserID; "seat_tier" → SeatTier;
// "rbac_group" → RbacGroupID; "organization_service" → Service;
// "organization" → (no other field).
type SpendLimitScope struct {
	Type        string `json:"type"`
	UserID      string `json:"user_id,omitempty"`
	SeatTier    string `json:"seat_tier,omitempty"`
	RbacGroupID string `json:"rbac_group_id,omitempty"`
	Service     string `json:"service,omitempty"`
}

type SetSpendLimitRequest struct {
	Amount string          `json:"amount"`
	Scope  SpendLimitScope `json:"scope"`
	Period string          `json:"period,omitempty"`
}

type SpendLimitDeleteResponse struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// SpendSummary is one row of the effective-limits list — a member's
// resolved spend limit + period-to-date consumption.
type SpendSummary struct {
	Actor             SpendSummaryActor `json:"actor"`
	Amount            string            `json:"amount"`
	Currency          string            `json:"currency"`
	Period            string            `json:"period"`
	PeriodToDateSpend string            `json:"period_to_date_spend"`
	Scope             SpendLimitScope   `json:"scope"`
	Source            SpendLimitScope   `json:"source"`
}

type SpendSummaryActor struct {
	Type         string `json:"type"`
	UserID       string `json:"user_id"`
	EmailAddress string `json:"email_address"`
	Name         string `json:"name"`
	Deleted      bool   `json:"deleted"`
}

type ListEffectiveSpendLimitsParams struct {
	Limit   int
	Period  []string
	UserIDs []string
}

type ListEffectiveSpendLimitsResponse struct {
	Data     []SpendSummary `json:"data"`
	NextPage *string        `json:"next_page"`
}

func (c *Client) SetSpendLimit(ctx context.Context, in SetSpendLimitRequest) (*SpendLimit, error) {
	var out SpendLimit
	if err := c.do(ctx, http.MethodPost, "/v1/organizations/spend_limits", in, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetSpendLimit(ctx context.Context, id string) (*SpendLimit, error) {
	var out SpendLimit
	if err := c.do(ctx, http.MethodGet, "/v1/organizations/spend_limits/"+url.PathEscape(id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteSpendLimit(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, "/v1/organizations/spend_limits/"+url.PathEscape(id), nil, nil)
}

func (c *Client) ListEffectiveSpendLimits(ctx context.Context, p ListEffectiveSpendLimitsParams) ([]SpendSummary, error) {
	var all []SpendSummary
	var cursor string
	for {
		q := url.Values{}
		if cursor != "" {
			q.Set("page", cursor)
		}
		if p.Limit > 0 {
			q.Set("limit", strconv.Itoa(p.Limit))
		}
		for _, period := range p.Period {
			q.Add("period", period)
		}
		for _, uid := range p.UserIDs {
			q.Add("user_ids", uid)
		}
		path := "/v1/organizations/spend_limits/effective"
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
		var page ListEffectiveSpendLimitsResponse
		if err := c.do(ctx, http.MethodGet, path, nil, &page); err != nil {
			return nil, err
		}
		all = append(all, page.Data...)
		if page.NextPage == nil || *page.NextPage == "" {
			return all, nil
		}
		cursor = *page.NextPage
	}
}

// ---------- Increase Requests ----------

type SpendLimitIncreaseRequest struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	Status      string            `json:"status"`
	Actor       SpendSummaryActor `json:"actor"`
	Amount      string            `json:"amount,omitempty"`
	Period      string            `json:"period"`
	Currency    string            `json:"currency,omitempty"`
	CreatedAt   string            `json:"created_at"`
	ResolvedAt  string            `json:"resolved_at,omitempty"`
	ResolvedBy  *IncreaseResolver `json:"resolved_by,omitempty"`
	Description string            `json:"description,omitempty"`
}

// IncreaseResolver flattens the polymorphic resolved_by — either a user
// actor or a scoped API key actor.
type IncreaseResolver struct {
	Type           string `json:"type"`
	UserID         string `json:"user_id,omitempty"`
	EmailAddress   string `json:"email_address,omitempty"`
	Name           string `json:"name,omitempty"`
	Deleted        bool   `json:"deleted,omitempty"`
	ScopedAPIKeyID string `json:"scoped_api_key_id,omitempty"`
}

type ListIncreaseRequestsParams struct {
	ActorIDs []string
	Status   []string
	Limit    int
}

type ListIncreaseRequestsResponse struct {
	Data     []SpendLimitIncreaseRequest `json:"data"`
	NextPage *string                     `json:"next_page"`
}

type ApproveIncreaseRequest struct {
	Amount               string `json:"amount"`
	Period               string `json:"period,omitempty"`
	SuppressNotification bool   `json:"suppress_notification,omitempty"`
}

type DenyIncreaseRequest struct {
	SuppressNotification bool `json:"suppress_notification,omitempty"`
}

func (c *Client) ListSpendLimitIncreaseRequests(ctx context.Context, p ListIncreaseRequestsParams) ([]SpendLimitIncreaseRequest, error) {
	var all []SpendLimitIncreaseRequest
	var cursor string
	for {
		q := url.Values{}
		if cursor != "" {
			q.Set("page", cursor)
		}
		if p.Limit > 0 {
			q.Set("limit", strconv.Itoa(p.Limit))
		}
		for _, a := range p.ActorIDs {
			q.Add("actor_ids", a)
		}
		for _, s := range p.Status {
			q.Add("status", s)
		}
		path := "/v1/organizations/spend_limit_increase_requests"
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
		var page ListIncreaseRequestsResponse
		if err := c.do(ctx, http.MethodGet, path, nil, &page); err != nil {
			return nil, err
		}
		all = append(all, page.Data...)
		if page.NextPage == nil || *page.NextPage == "" {
			return all, nil
		}
		cursor = *page.NextPage
	}
}

func (c *Client) GetSpendLimitIncreaseRequest(ctx context.Context, id string) (*SpendLimitIncreaseRequest, error) {
	var out SpendLimitIncreaseRequest
	if err := c.do(ctx, http.MethodGet, "/v1/organizations/spend_limit_increase_requests/"+url.PathEscape(id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ApproveSpendLimitIncreaseRequest(ctx context.Context, id string, in ApproveIncreaseRequest) (*SpendLimitIncreaseRequest, error) {
	var out SpendLimitIncreaseRequest
	if err := c.do(ctx, http.MethodPost, "/v1/organizations/spend_limit_increase_requests/"+url.PathEscape(id)+"/approve", in, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DenySpendLimitIncreaseRequest(ctx context.Context, id string, in DenyIncreaseRequest) (*SpendLimitIncreaseRequest, error) {
	var out SpendLimitIncreaseRequest
	if err := c.do(ctx, http.MethodPost, "/v1/organizations/spend_limit_increase_requests/"+url.PathEscape(id)+"/deny", in, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
