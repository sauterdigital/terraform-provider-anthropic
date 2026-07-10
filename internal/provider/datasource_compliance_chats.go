package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/sauterdigital/terraform-provider-claudeadmin/internal/anthropic"
)

var (
	_ datasource.DataSource              = &ComplianceChatsDataSource{}
	_ datasource.DataSourceWithConfigure = &ComplianceChatsDataSource{}
)

func NewComplianceChatsDataSource() datasource.DataSource { return &ComplianceChatsDataSource{} }

type ComplianceChatsDataSource struct{ client *anthropic.Client }

type complianceChatModel struct {
	ID           types.String `tfsdk:"id"`
	UserID       types.String `tfsdk:"user_id"`
	Title        types.String `tfsdk:"title"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
	MessageCount types.Int64  `tfsdk:"message_count"`
	ProjectID    types.String `tfsdk:"project_id"`
}

type ComplianceChatsModel struct {
	UserID     types.String          `tfsdk:"user_id"`
	StartingAt types.String          `tfsdk:"starting_at"`
	EndingAt   types.String          `tfsdk:"ending_at"`
	Limit      types.Int64           `tfsdk:"limit"`
	Chats      []complianceChatModel `tfsdk:"chats"`
}

func (d *ComplianceChatsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_compliance_chats"
}

func (d *ComplianceChatsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists Claude chats (`/v1/compliance/apps/chats`) for eDiscovery, DLP scans, and audit. Enterprise + Compliance Access Key with scope `read:compliance_user_data`. All filters optional — narrow by `user_id` + time window for busy orgs. Provider paginates until exhausted; the full result set materializes in state.",
		Attributes: map[string]schema.Attribute{
			"user_id":     schema.StringAttribute{Optional: true, Description: "Filter by owning user_id."},
			"starting_at": schema.StringAttribute{Optional: true, Description: "RFC3339 lower bound (inclusive) on `created_at`."},
			"ending_at":   schema.StringAttribute{Optional: true, Description: "RFC3339 upper bound (exclusive) on `created_at`."},
			"limit":       schema.Int64Attribute{Optional: true, Description: "Per-page size (default 100)."},
			"chats": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":            schema.StringAttribute{Computed: true},
						"user_id":       schema.StringAttribute{Computed: true},
						"title":         schema.StringAttribute{Computed: true},
						"created_at":    schema.StringAttribute{Computed: true},
						"updated_at":    schema.StringAttribute{Computed: true},
						"message_count": schema.Int64Attribute{Computed: true},
						"project_id":    schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *ComplianceChatsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := clientFromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	d.client = c
}

func (d *ComplianceChatsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var cfg ComplianceChatsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	chats, err := d.client.ListComplianceChats(ctx, anthropic.ListComplianceChatsParams{
		UserID:     cfg.UserID.ValueString(),
		StartingAt: cfg.StartingAt.ValueString(),
		EndingAt:   cfg.EndingAt.ValueString(),
		Limit:      int(cfg.Limit.ValueInt64()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to list compliance chats", err.Error())
		return
	}
	out := make([]complianceChatModel, 0, len(chats))
	for _, c := range chats {
		mc := types.Int64Null()
		if c.MessageCount != nil {
			mc = types.Int64Value(*c.MessageCount)
		}
		out = append(out, complianceChatModel{
			ID:           types.StringValue(c.ID),
			UserID:       types.StringValue(c.UserID),
			Title:        optionalStringValue(c.Title),
			CreatedAt:    types.StringValue(c.CreatedAt),
			UpdatedAt:    optionalStringValue(c.UpdatedAt),
			MessageCount: mc,
			ProjectID:    optionalStringValue(c.ProjectID),
		})
	}
	cfg.Chats = out
	resp.Diagnostics.Append(resp.State.Set(ctx, cfg)...)
}
