package anthropic

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// Analytics endpoints share a common time-bucketed envelope:
//   { data: [...], next_page: string }
// All require Claude Enterprise plan + an API key with read:analytics scope.

// ---------- Activity Summaries ----------

type ActivitySummary struct {
	StartingAt                   string   `json:"starting_at"`
	EndingAt                     string   `json:"ending_at"`
	AssignedSeatCount            *int64   `json:"assigned_seat_count"`
	DailyActiveUserCount         int64    `json:"daily_active_user_count"`
	WeeklyActiveUserCount        int64    `json:"weekly_active_user_count"`
	MonthlyActiveUserCount       int64    `json:"monthly_active_user_count"`
	DailyAdoptionRate            *float64 `json:"daily_adoption_rate"`
	WeeklyAdoptionRate           *float64 `json:"weekly_adoption_rate"`
	MonthlyAdoptionRate          *float64 `json:"monthly_adoption_rate"`
	CoworkDailyActiveUserCount   int64    `json:"cowork_daily_active_user_count"`
	CoworkWeeklyActiveUserCount  int64    `json:"cowork_weekly_active_user_count"`
	CoworkMonthlyActiveUserCount int64    `json:"cowork_monthly_active_user_count"`
	PendingInviteCount           *int64   `json:"pending_invite_count"`
}

type ActivitySummariesResponse struct {
	Summaries []ActivitySummary `json:"summaries"`
}

type GetActivitySummariesParams struct {
	StartingDate string
	EndingDate   string
}

func (c *Client) GetActivitySummaries(ctx context.Context, p GetActivitySummariesParams) (*ActivitySummariesResponse, error) {
	q := url.Values{}
	q.Set("starting_date", p.StartingDate)
	if p.EndingDate != "" {
		q.Set("ending_date", p.EndingDate)
	}
	var out ActivitySummariesResponse
	if err := c.do(ctx, http.MethodGet, "/v1/organizations/analytics/summaries?"+q.Encode(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Token / Cost over time + per-user — generic envelope ----------

// AnalyticsParams are filters shared across time-bucketed endpoints.
type AnalyticsParams struct {
	StartingAt     string
	EndingAt       string
	BucketWidth    string
	GroupBy        []string
	Products       []string
	Models         []string
	ContextWindows []string
	InferenceGeos  []string
	Speeds         []string
	UserIDs        []string
	Limit          int
	Page           string
}

// TokenUsageBucket — one row in /analytics/usage_report response.
type TokenUsageBucket struct {
	StartingAt string             `json:"starting_at"`
	EndingAt   string             `json:"ending_at"`
	Results    []TokenUsageResult `json:"results"`
}

type TokenUsageResult struct {
	UncachedInputTokens  int64         `json:"uncached_input_tokens"`
	CacheReadInputTokens int64         `json:"cache_read_input_tokens"`
	OutputTokens         int64         `json:"output_tokens"`
	CacheCreation        CacheCreation `json:"cache_creation"`
	ServerToolUse        ServerToolUse `json:"server_tool_use"`
	Requests             int64         `json:"requests"`
	Product              *string       `json:"product"`
	Model                *string       `json:"model"`
	ContextWindow        *string       `json:"context_window"`
	InferenceGeo         *string       `json:"inference_geo"`
	Speed                *string       `json:"speed"`
	UserID               *string       `json:"user_id"`
}

type TokenUsageResponse struct {
	Data     []TokenUsageBucket `json:"data"`
	NextPage *string            `json:"next_page"`
}

// CostBucketV2 — analytics cost_report bucket.
type CostBucketV2 struct {
	StartingAt string         `json:"starting_at"`
	EndingAt   string         `json:"ending_at"`
	Results    []CostResultV2 `json:"results"`
}

type CostResultV2 struct {
	Amount        string  `json:"amount"`
	Currency      string  `json:"currency"`
	Product       *string `json:"product"`
	Model         *string `json:"model"`
	ContextWindow *string `json:"context_window"`
	InferenceGeo  *string `json:"inference_geo"`
	Speed         *string `json:"speed"`
	UserID        *string `json:"user_id"`
	TokenType     *string `json:"token_type"`
}

type CostResponseV2 struct {
	Data     []CostBucketV2 `json:"data"`
	NextPage *string        `json:"next_page"`
}

func (c *Client) GetTokenUsageOverTime(ctx context.Context, p AnalyticsParams) (*TokenUsageResponse, error) {
	var out TokenUsageResponse
	q := buildAnalyticsQuery(p)
	if err := c.do(ctx, http.MethodGet, "/v1/organizations/analytics/usage_report?"+q.Encode(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetPerUserTokenUsage(ctx context.Context, p AnalyticsParams) (*TokenUsageResponse, error) {
	var out TokenUsageResponse
	q := buildAnalyticsQuery(p)
	if err := c.do(ctx, http.MethodGet, "/v1/organizations/analytics/user_usage_report?"+q.Encode(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetCostOverTime(ctx context.Context, p AnalyticsParams) (*CostResponseV2, error) {
	var out CostResponseV2
	q := buildAnalyticsQuery(p)
	if err := c.do(ctx, http.MethodGet, "/v1/organizations/analytics/cost_report?"+q.Encode(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetPerUserCost(ctx context.Context, p AnalyticsParams) (*CostResponseV2, error) {
	var out CostResponseV2
	q := buildAnalyticsQuery(p)
	if err := c.do(ctx, http.MethodGet, "/v1/organizations/analytics/user_cost_report?"+q.Encode(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func buildAnalyticsQuery(p AnalyticsParams) url.Values {
	q := url.Values{}
	if p.StartingAt != "" {
		q.Set("starting_at", p.StartingAt)
	}
	if p.EndingAt != "" {
		q.Set("ending_at", p.EndingAt)
	}
	if p.BucketWidth != "" {
		q.Set("bucket_width", p.BucketWidth)
	}
	if p.Limit > 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Page != "" {
		q.Set("page", p.Page)
	}
	addList(q, "group_by", p.GroupBy)
	addList(q, "products", p.Products)
	addList(q, "models", p.Models)
	addList(q, "context_windows", p.ContextWindows)
	addList(q, "inference_geos", p.InferenceGeos)
	addList(q, "speeds", p.Speeds)
	addList(q, "user_ids", p.UserIDs)
	return q
}

// ---------- User Activity ----------

// AnalyticsUser is one row of /analytics/users — per-user activity snapshot.
type AnalyticsUser struct {
	UserID            string  `json:"user_id"`
	EmailAddress      string  `json:"email_address"`
	Name              string  `json:"name"`
	SeatTier          *string `json:"seat_tier"`
	LastActiveAt      *string `json:"last_active_at"`
	TotalRequests     int64   `json:"total_requests"`
	TotalInputTokens  int64   `json:"total_input_tokens"`
	TotalOutputTokens int64   `json:"total_output_tokens"`
	TotalCostAmount   *string `json:"total_cost_amount"`
	Currency          *string `json:"currency"`
}

type ListUserActivityResponse struct {
	Data     []AnalyticsUser `json:"data"`
	NextPage *string         `json:"next_page"`
}

type ListUserActivityParams struct {
	StartingDate string
	EndingDate   string
	UserIDs      []string
	Limit        int
}

func (c *Client) ListUserActivity(ctx context.Context, p ListUserActivityParams) (*ListUserActivityResponse, error) {
	q := url.Values{}
	if p.StartingDate != "" {
		q.Set("starting_date", p.StartingDate)
	}
	if p.EndingDate != "" {
		q.Set("ending_date", p.EndingDate)
	}
	if p.Limit > 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	addList(q, "user_ids", p.UserIDs)
	var out ListUserActivityResponse
	if err := c.do(ctx, http.MethodGet, "/v1/organizations/analytics/users?"+q.Encode(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- Skills / Connectors / Chat Projects usage ----------

// SkillUsageEntry is one row of /analytics/skills.
type SkillUsageEntry struct {
	SkillName       string `json:"skill_name"`
	InvocationCount int64  `json:"invocation_count"`
	SuccessCount    int64  `json:"success_count"`
	FailureCount    int64  `json:"failure_count"`
}

type SkillUsageResponse struct {
	Data     []SkillUsageEntry `json:"data"`
	NextPage *string           `json:"next_page"`
}

type ConnectorUsageEntry struct {
	ConnectorName   string `json:"connector_name"`
	InvocationCount int64  `json:"invocation_count"`
	UniqueUsers     int64  `json:"unique_users"`
}

type ConnectorUsageResponse struct {
	Data     []ConnectorUsageEntry `json:"data"`
	NextPage *string               `json:"next_page"`
}

type ChatProjectUsageEntry struct {
	ProjectID    string `json:"project_id"`
	ProjectName  string `json:"project_name"`
	MessageCount int64  `json:"message_count"`
	UniqueUsers  int64  `json:"unique_users"`
}

type ChatProjectUsageResponse struct {
	Data     []ChatProjectUsageEntry `json:"data"`
	NextPage *string                 `json:"next_page"`
}

type SimpleAnalyticsParams struct {
	StartingDate string
	EndingDate   string
	Limit        int
}

func (c *Client) GetSkillsUsage(ctx context.Context, p SimpleAnalyticsParams) (*SkillUsageResponse, error) {
	q := simpleAnalyticsQuery(p)
	var out SkillUsageResponse
	if err := c.do(ctx, http.MethodGet, "/v1/organizations/analytics/skills?"+q.Encode(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetConnectorsUsage(ctx context.Context, p SimpleAnalyticsParams) (*ConnectorUsageResponse, error) {
	q := simpleAnalyticsQuery(p)
	var out ConnectorUsageResponse
	if err := c.do(ctx, http.MethodGet, "/v1/organizations/analytics/connectors?"+q.Encode(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetChatProjectsUsage(ctx context.Context, p SimpleAnalyticsParams) (*ChatProjectUsageResponse, error) {
	q := simpleAnalyticsQuery(p)
	var out ChatProjectUsageResponse
	if err := c.do(ctx, http.MethodGet, "/v1/organizations/analytics/apps/chat/projects?"+q.Encode(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func simpleAnalyticsQuery(p SimpleAnalyticsParams) url.Values {
	q := url.Values{}
	if p.StartingDate != "" {
		q.Set("starting_date", p.StartingDate)
	}
	if p.EndingDate != "" {
		q.Set("ending_date", p.EndingDate)
	}
	if p.Limit > 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	return q
}
