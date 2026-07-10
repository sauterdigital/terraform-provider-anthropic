package anthropic

import (
	"context"
	"net/http"
	"net/url"
)

// Compliance Content API surface (/v1/compliance/apps/*).
//
// Handles chat messages, attached files, generated files, artifacts, projects
// and project documents — the actual user-generated content, distinct from
// the directory / activity metadata under /v1/compliance/*.
//
// Auth: same Compliance Access Key (sk-ant-api01-...) via x-api-key, but
// requires the content-data scopes on the key:
//   read:compliance_user_data    — all GETs and lists
//   delete:compliance_user_data  — DELETE endpoints (eDiscovery / DLP)
//
// The Admin API key is rejected with 403 for these paths.

// ---- Types ----

type ComplianceChat struct {
	ID           string  `json:"id"`
	Type         string  `json:"type"`
	UserID       string  `json:"user_id"`
	Title        *string `json:"title"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    *string `json:"updated_at"`
	MessageCount *int64  `json:"message_count"`
	ProjectID    *string `json:"project_id"`
}

type ComplianceChatMessage struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	ChatID    string  `json:"chat_id"`
	Role      string  `json:"role"` // "user" | "assistant"
	Text      *string `json:"text"`
	Model     *string `json:"model"`
	CreatedAt string  `json:"created_at"`
}

type ComplianceProject struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	UserID      string  `json:"user_id"`
	Description *string `json:"description"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
	ChatCount   *int64  `json:"chat_count"`
}

type ComplianceProjectAttachment struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	ProjectID string  `json:"project_id"`
	Filename  *string `json:"filename"`
	MimeType  *string `json:"mime_type"`
	SizeBytes *int64  `json:"size_bytes"`
	CreatedAt string  `json:"created_at"`
}

// ---- Params ----

type ListComplianceChatsParams struct {
	UserID     string
	StartingAt string
	EndingAt   string
	Limit      int
}

type ListComplianceProjectsParams struct {
	UserID string
	Limit  int
}

type paginatedComplianceChats struct {
	Data    []ComplianceChat `json:"data"`
	HasMore bool             `json:"has_more"`
	LastID  string           `json:"last_id"`
}

type paginatedComplianceChatMessages struct {
	Data    []ComplianceChatMessage `json:"data"`
	HasMore bool                    `json:"has_more"`
	LastID  string                  `json:"last_id"`
}

type paginatedComplianceProjects struct {
	Data    []ComplianceProject `json:"data"`
	HasMore bool                `json:"has_more"`
	LastID  string              `json:"last_id"`
}

type paginatedComplianceProjectAttachments struct {
	Data    []ComplianceProjectAttachment `json:"data"`
	HasMore bool                          `json:"has_more"`
	LastID  string                        `json:"last_id"`
}

// ---- Reads ----

func (c *Client) ListComplianceChats(ctx context.Context, p ListComplianceChatsParams) ([]ComplianceChat, error) {
	if err := c.requireCompliance(); err != nil {
		return nil, err
	}
	var all []ComplianceChat
	var afterID string
	limit := p.Limit
	if limit <= 0 {
		limit = 100
	}
	for {
		q := url.Values{}
		q.Set("limit", intStr(limit))
		if afterID != "" {
			q.Set("after_id", afterID)
		}
		if p.UserID != "" {
			q.Set("user_id", p.UserID)
		}
		if p.StartingAt != "" {
			q.Set("starting_at", p.StartingAt)
		}
		if p.EndingAt != "" {
			q.Set("ending_at", p.EndingAt)
		}
		var page paginatedComplianceChats
		if err := c.do(complianceCtx(ctx), http.MethodGet, "/v1/compliance/apps/chats?"+q.Encode(), nil, &page); err != nil {
			return nil, err
		}
		all = append(all, page.Data...)
		if !page.HasMore || page.LastID == "" {
			return all, nil
		}
		afterID = page.LastID
	}
}

func (c *Client) ListComplianceChatMessages(ctx context.Context, chatID string) ([]ComplianceChatMessage, error) {
	if err := c.requireCompliance(); err != nil {
		return nil, err
	}
	var all []ComplianceChatMessage
	var afterID string
	for {
		q := url.Values{}
		q.Set("limit", "100")
		if afterID != "" {
			q.Set("after_id", afterID)
		}
		var page paginatedComplianceChatMessages
		path := "/v1/compliance/apps/chats/" + url.PathEscape(chatID) + "/messages?" + q.Encode()
		if err := c.do(complianceCtx(ctx), http.MethodGet, path, nil, &page); err != nil {
			return nil, err
		}
		all = append(all, page.Data...)
		if !page.HasMore || page.LastID == "" {
			return all, nil
		}
		afterID = page.LastID
	}
}

func (c *Client) ListComplianceProjects(ctx context.Context, p ListComplianceProjectsParams) ([]ComplianceProject, error) {
	if err := c.requireCompliance(); err != nil {
		return nil, err
	}
	var all []ComplianceProject
	var afterID string
	limit := p.Limit
	if limit <= 0 {
		limit = 100
	}
	for {
		q := url.Values{}
		q.Set("limit", intStr(limit))
		if afterID != "" {
			q.Set("after_id", afterID)
		}
		if p.UserID != "" {
			q.Set("user_id", p.UserID)
		}
		var page paginatedComplianceProjects
		if err := c.do(complianceCtx(ctx), http.MethodGet, "/v1/compliance/apps/projects?"+q.Encode(), nil, &page); err != nil {
			return nil, err
		}
		all = append(all, page.Data...)
		if !page.HasMore || page.LastID == "" {
			return all, nil
		}
		afterID = page.LastID
	}
}

func (c *Client) ListComplianceProjectAttachments(ctx context.Context, projectID string) ([]ComplianceProjectAttachment, error) {
	if err := c.requireCompliance(); err != nil {
		return nil, err
	}
	var all []ComplianceProjectAttachment
	var afterID string
	for {
		q := url.Values{}
		q.Set("limit", "100")
		if afterID != "" {
			q.Set("after_id", afterID)
		}
		var page paginatedComplianceProjectAttachments
		path := "/v1/compliance/apps/projects/" + url.PathEscape(projectID) + "/attachments?" + q.Encode()
		if err := c.do(complianceCtx(ctx), http.MethodGet, path, nil, &page); err != nil {
			return nil, err
		}
		all = append(all, page.Data...)
		if !page.HasMore || page.LastID == "" {
			return all, nil
		}
		afterID = page.LastID
	}
}

// ---- Deletes (eDiscovery / DLP; require delete:compliance_user_data scope) ----

// DeleteComplianceContent maps target_type -> API path. See resource
// docs for the accepted values.
func (c *Client) DeleteComplianceContent(ctx context.Context, targetType, targetID string) error {
	if err := c.requireCompliance(); err != nil {
		return err
	}
	var path string
	switch targetType {
	case "chat":
		path = "/v1/compliance/apps/chats/" + url.PathEscape(targetID)
	case "chat_file":
		path = "/v1/compliance/apps/chats/files/" + url.PathEscape(targetID)
	case "chat_generated_file":
		path = "/v1/compliance/apps/chats/generated_files/" + url.PathEscape(targetID)
	case "project":
		path = "/v1/compliance/apps/projects/" + url.PathEscape(targetID)
	case "project_document":
		path = "/v1/compliance/apps/projects/documents/" + url.PathEscape(targetID)
	default:
		return &APIError{StatusCode: 400, Type: "invalid_request_error", Message: "unknown compliance content target_type: " + targetType}
	}
	return c.do(complianceCtx(ctx), http.MethodDelete, path, nil, nil)
}
