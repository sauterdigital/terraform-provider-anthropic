package anthropic

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// All Service Account endpoints REQUIRE OAuth Bearer auth.
// Admin API keys are rejected — see the client's HasOAuth() guard.

type ServiceAccount struct {
	ID                string `json:"id"`
	Type              string `json:"type"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	OrganizationRole  string `json:"organization_role"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	ArchivedAt        string `json:"archived_at,omitempty"`
	CreatedByActorID  string `json:"created_by_actor_id"`
	UpdatedByActorID  string `json:"updated_by_actor_id"`
	ArchivedByActorID string `json:"archived_by_actor_id,omitempty"`
}

type CreateServiceAccountRequest struct {
	Name             string `json:"name"`
	Description      string `json:"description,omitempty"`
	OrganizationRole string `json:"organization_role,omitempty"`
}

type UpdateServiceAccountRequest struct {
	Description      *string `json:"description,omitempty"`
	OrganizationRole *string `json:"organization_role,omitempty"`
}

type ListServiceAccountsParams struct {
	Limit            int
	Page             string
	IncludeArchived  bool
	OrganizationRole string
}

type ListServiceAccountsResponse struct {
	Data     []ServiceAccount `json:"data"`
	NextPage *string          `json:"next_page"`
}

func (c *Client) CreateServiceAccount(ctx context.Context, in CreateServiceAccountRequest) (*ServiceAccount, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var out ServiceAccount
	if err := c.do(ctx, http.MethodPost, "/v1/organizations/service_accounts", in, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetServiceAccount(ctx context.Context, id string) (*ServiceAccount, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var out ServiceAccount
	if err := c.do(ctx, http.MethodGet, "/v1/organizations/service_accounts/"+url.PathEscape(id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateServiceAccount(ctx context.Context, id string, in UpdateServiceAccountRequest) (*ServiceAccount, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var out ServiceAccount
	if err := c.do(ctx, http.MethodPost, "/v1/organizations/service_accounts/"+url.PathEscape(id), in, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ArchiveServiceAccount(ctx context.Context, id string) (*ServiceAccount, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var out ServiceAccount
	if err := c.do(ctx, http.MethodPost, "/v1/organizations/service_accounts/"+url.PathEscape(id)+"/archive", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListServiceAccounts(ctx context.Context, p ListServiceAccountsParams) ([]ServiceAccount, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var all []ServiceAccount
	cursor := p.Page
	for {
		q := url.Values{}
		if cursor != "" {
			q.Set("page", cursor)
		}
		if p.Limit > 0 {
			q.Set("limit", strconv.Itoa(p.Limit))
		}
		if p.IncludeArchived {
			q.Set("include_archived", "true")
		}
		if p.OrganizationRole != "" {
			q.Set("organization_role", p.OrganizationRole)
		}
		path := "/v1/organizations/service_accounts"
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
		var page ListServiceAccountsResponse
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

// ---------- SA × Workspace bindings ----------

type ServiceAccountWorkspaceMember struct {
	Type             string `json:"type"`
	ServiceAccountID string `json:"service_account_id"`
	WorkspaceID      string `json:"workspace_id"`
	WorkspaceRole    string `json:"workspace_role"`
	Implicit         bool   `json:"implicit"`
	CreatedByActorID string `json:"created_by_actor_id,omitempty"`
}

type AddWorkspaceToServiceAccountRequest struct {
	WorkspaceID   string `json:"workspace_id"`
	WorkspaceRole string `json:"workspace_role"`
}

type UpdateWorkspaceServiceAccountRoleRequest struct {
	WorkspaceRole string `json:"workspace_role"`
}

type ListSAWorkspacesResponse struct {
	Data     []ServiceAccountWorkspaceMember `json:"data"`
	NextPage *string                         `json:"next_page"`
}

func (c *Client) AddWorkspaceToServiceAccount(ctx context.Context, saID string, in AddWorkspaceToServiceAccountRequest) (*ServiceAccountWorkspaceMember, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var out ServiceAccountWorkspaceMember
	if err := c.do(ctx, http.MethodPost, "/v1/organizations/service_accounts/"+url.PathEscape(saID)+"/workspaces", in, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateServiceAccountWorkspaceRole(ctx context.Context, workspaceID, saID, role string) (*ServiceAccountWorkspaceMember, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var out ServiceAccountWorkspaceMember
	path := "/v1/organizations/workspaces/" + url.PathEscape(workspaceID) + "/service_accounts/" + url.PathEscape(saID)
	if err := c.do(ctx, http.MethodPost, path, UpdateWorkspaceServiceAccountRoleRequest{WorkspaceRole: role}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetServiceAccountWorkspaceMember(ctx context.Context, workspaceID, saID string) (*ServiceAccountWorkspaceMember, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var out ServiceAccountWorkspaceMember
	path := "/v1/organizations/workspaces/" + url.PathEscape(workspaceID) + "/service_accounts/" + url.PathEscape(saID)
	if err := c.do(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) RemoveWorkspaceFromServiceAccount(ctx context.Context, saID, workspaceID string) error {
	if !c.HasOAuth() {
		return ErrOAuthRequired
	}
	return c.do(ctx, http.MethodDelete, "/v1/organizations/service_accounts/"+url.PathEscape(saID)+"/workspaces/"+url.PathEscape(workspaceID), nil, nil)
}

func (c *Client) ListServiceAccountWorkspaces(ctx context.Context, saID string) ([]ServiceAccountWorkspaceMember, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var all []ServiceAccountWorkspaceMember
	var cursor string
	for {
		q := url.Values{}
		if cursor != "" {
			q.Set("page", cursor)
		}
		path := "/v1/organizations/service_accounts/" + url.PathEscape(saID) + "/workspaces"
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
		var page ListSAWorkspacesResponse
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

func (c *Client) ListWorkspaceServiceAccountMembers(ctx context.Context, workspaceID string) ([]ServiceAccountWorkspaceMember, error) {
	if !c.HasOAuth() {
		return nil, ErrOAuthRequired
	}
	var all []ServiceAccountWorkspaceMember
	var cursor string
	for {
		q := url.Values{}
		if cursor != "" {
			q.Set("page", cursor)
		}
		path := "/v1/organizations/workspaces/" + url.PathEscape(workspaceID) + "/service_accounts"
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
		var page ListSAWorkspacesResponse
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
