package anthropic

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// All Federation endpoints REQUIRE OAuth Bearer auth.

// ---------- Federation Issuers ----------

type FederationIssuer struct {
	ID                    string     `json:"id"`
	Type                  string     `json:"type"`
	Name                  string     `json:"name"`
	IssuerURL             string     `json:"issuer_url"`
	CheckJTI              bool       `json:"check_jti"`
	MaxJWTLifetimeSeconds int64      `json:"max_jwt_lifetime_seconds"`
	JWKS                  IssuerJWKS `json:"jwks"`
	JWKSPollingDisabledAt string     `json:"jwks_polling_disabled_at,omitempty"`
	CreatedAt             string     `json:"created_at"`
	UpdatedAt             string     `json:"updated_at,omitempty"`
	CreatedByActorID      string     `json:"created_by_actor_id,omitempty"`
	UpdatedByActorID      string     `json:"updated_by_actor_id,omitempty"`
	ArchivedAt            string     `json:"archived_at,omitempty"`
	ArchivedByActorID     string     `json:"archived_by_actor_id,omitempty"`
}

// IssuerJWKS flattens the oneOf union (discovery / explicit_url / inline).
// Discriminated by Type. Only the fields matching that Type are populated.
type IssuerJWKS struct {
	Type          string            `json:"type"`
	DiscoveryBase string            `json:"discovery_base,omitempty"`
	URL           string            `json:"url,omitempty"`
	CACertPEM     string            `json:"ca_cert_pem,omitempty"`
	Keys          []json.RawMessage `json:"keys,omitempty"`
}

type CreateFederationIssuerRequest struct {
	IssuerURL             string      `json:"issuer_url"`
	Name                  string      `json:"name"`
	CheckJTI              *bool       `json:"check_jti,omitempty"`
	JWKS                  *IssuerJWKS `json:"jwks,omitempty"`
	MaxJWTLifetimeSeconds *int64      `json:"max_jwt_lifetime_seconds,omitempty"`
}

type UpdateFederationIssuerRequest struct {
	IssuerURL             *string     `json:"issuer_url,omitempty"`
	CheckJTI              *bool       `json:"check_jti,omitempty"`
	JWKS                  *IssuerJWKS `json:"jwks,omitempty"`
	MaxJWTLifetimeSeconds *int64      `json:"max_jwt_lifetime_seconds,omitempty"`
	JWKSPollingDisabled   *bool       `json:"jwks_polling_disabled,omitempty"`
}

type ListFederationIssuersResponse struct {
	Data     []FederationIssuer `json:"data"`
	NextPage *string            `json:"next_page"`
}

func (c *Client) CreateFederationIssuer(ctx context.Context, in CreateFederationIssuerRequest) (*FederationIssuer, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var out FederationIssuer
	if err := c.do(ctx, http.MethodPost, "/v1/organizations/federation_issuers", in, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetFederationIssuer(ctx context.Context, id string) (*FederationIssuer, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var out FederationIssuer
	if err := c.do(ctx, http.MethodGet, "/v1/organizations/federation_issuers/"+url.PathEscape(id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateFederationIssuer(ctx context.Context, id string, in UpdateFederationIssuerRequest) (*FederationIssuer, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var out FederationIssuer
	if err := c.do(ctx, http.MethodPost, "/v1/organizations/federation_issuers/"+url.PathEscape(id), in, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ArchiveFederationIssuer(ctx context.Context, id string) (*FederationIssuer, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var out FederationIssuer
	if err := c.do(ctx, http.MethodPost, "/v1/organizations/federation_issuers/"+url.PathEscape(id)+"/archive", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListFederationIssuers(ctx context.Context) ([]FederationIssuer, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var all []FederationIssuer
	var cursor string
	for {
		q := url.Values{}
		if cursor != "" {
			q.Set("page", cursor)
		}
		path := "/v1/organizations/federation_issuers"
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
		var page ListFederationIssuersResponse
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

// ---------- Federation Rules ----------

type FederationRule struct {
	ID                     string           `json:"id"`
	Type                   string           `json:"type"`
	Name                   string           `json:"name"`
	Description            string           `json:"description,omitempty"`
	IssuerID               string           `json:"issuer_id"`
	OAuthScope             string           `json:"oauth_scope"`
	Match                  FederationMatch  `json:"match"`
	Target                 FederationTarget `json:"target"`
	WorkspaceID            string           `json:"workspace_id,omitempty"`
	AppliesToAllWorkspaces bool             `json:"applies_to_all_workspaces"`
	TokenLifetimeSeconds   int64            `json:"token_lifetime_seconds"`
	CreatedAt              string           `json:"created_at"`
	UpdatedAt              string           `json:"updated_at,omitempty"`
	ArchivedAt             string           `json:"archived_at,omitempty"`
	CreatedByActorID       string           `json:"created_by_actor_id,omitempty"`
}

type FederationMatch struct {
	Audience      string            `json:"audience,omitempty"`
	Claims        map[string]string `json:"claims,omitempty"`
	Condition     string            `json:"condition,omitempty"`
	SubjectPrefix string            `json:"subject_prefix,omitempty"`
}

type FederationTarget struct {
	Type               string `json:"type"`
	ServiceAccountID   string `json:"service_account_id,omitempty"`
	ServiceAccountName string `json:"service_account_name,omitempty"`
}

type CreateFederationRuleRequest struct {
	IssuerID               string           `json:"issuer_id"`
	Name                   string           `json:"name"`
	OAuthScope             string           `json:"oauth_scope"`
	Match                  FederationMatch  `json:"match"`
	Target                 FederationTarget `json:"target"`
	Description            string           `json:"description,omitempty"`
	AppliesToAllWorkspaces *bool            `json:"applies_to_all_workspaces,omitempty"`
	WorkspaceID            string           `json:"workspace_id,omitempty"`
	TokenLifetimeSeconds   *int64           `json:"token_lifetime_seconds,omitempty"`
}

type UpdateFederationRuleRequest struct {
	Description            *string           `json:"description,omitempty"`
	OAuthScope             *string           `json:"oauth_scope,omitempty"`
	Match                  *FederationMatch  `json:"match,omitempty"`
	Target                 *FederationTarget `json:"target,omitempty"`
	AppliesToAllWorkspaces *bool             `json:"applies_to_all_workspaces,omitempty"`
	TokenLifetimeSeconds   *int64            `json:"token_lifetime_seconds,omitempty"`
}

type ListFederationRulesResponse struct {
	Data     []FederationRule `json:"data"`
	NextPage *string          `json:"next_page"`
}

func (c *Client) CreateFederationRule(ctx context.Context, in CreateFederationRuleRequest) (*FederationRule, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var out FederationRule
	if err := c.do(ctx, http.MethodPost, "/v1/organizations/federation_rules", in, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetFederationRule(ctx context.Context, id string) (*FederationRule, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var out FederationRule
	if err := c.do(ctx, http.MethodGet, "/v1/organizations/federation_rules/"+url.PathEscape(id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateFederationRule(ctx context.Context, id string, in UpdateFederationRuleRequest) (*FederationRule, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var out FederationRule
	if err := c.do(ctx, http.MethodPost, "/v1/organizations/federation_rules/"+url.PathEscape(id), in, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ArchiveFederationRule(ctx context.Context, id string) (*FederationRule, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var out FederationRule
	if err := c.do(ctx, http.MethodPost, "/v1/organizations/federation_rules/"+url.PathEscape(id)+"/archive", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListFederationRules(ctx context.Context) ([]FederationRule, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var all []FederationRule
	var cursor string
	for {
		q := url.Values{}
		if cursor != "" {
			q.Set("page", cursor)
		}
		path := "/v1/organizations/federation_rules"
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
		var page ListFederationRulesResponse
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

// ---------- Federation Rule × Workspace bindings ----------

type FederationRuleWorkspace struct {
	Type             string `json:"type"`
	FederationRuleID string `json:"federation_rule_id"`
	WorkspaceID      string `json:"workspace_id"`
}

type AddFederationRuleWorkspaceRequest struct {
	WorkspaceID string `json:"workspace_id"`
}

type ListFederationRuleWorkspacesResponse struct {
	Data     []FederationRuleWorkspace `json:"data"`
	NextPage *string                   `json:"next_page"`
}

func (c *Client) AddFederationRuleWorkspace(ctx context.Context, ruleID, workspaceID string) (*FederationRuleWorkspace, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var out FederationRuleWorkspace
	path := "/v1/organizations/federation_rules/" + url.PathEscape(ruleID) + "/workspaces"
	if err := c.do(ctx, http.MethodPost, path, AddFederationRuleWorkspaceRequest{WorkspaceID: workspaceID}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) RemoveFederationRuleWorkspace(ctx context.Context, ruleID, workspaceID string) error {
	if !c.HasOAuth() {
		return ErrOAuthRequired
	}
	path := "/v1/organizations/federation_rules/" + url.PathEscape(ruleID) + "/workspaces/" + url.PathEscape(workspaceID)
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) ListFederationRuleWorkspaces(ctx context.Context, ruleID string) ([]FederationRuleWorkspace, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var all []FederationRuleWorkspace
	var cursor string
	for {
		q := url.Values{}
		if cursor != "" {
			q.Set("page", cursor)
		}
		if x := strconv.Itoa(0); x == "" {
			_ = x
		}
		path := "/v1/organizations/federation_rules/" + url.PathEscape(ruleID) + "/workspaces"
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
		var page ListFederationRuleWorkspacesResponse
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
