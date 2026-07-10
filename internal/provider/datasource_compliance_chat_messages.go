package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/sauterdigital/terraform-provider-claudeadmin/internal/anthropic"
)

var (
	_ datasource.DataSource              = &ComplianceChatMessagesDataSource{}
	_ datasource.DataSourceWithConfigure = &ComplianceChatMessagesDataSource{}
)

func NewComplianceChatMessagesDataSource() datasource.DataSource {
	return &ComplianceChatMessagesDataSource{}
}

type ComplianceChatMessagesDataSource struct{ client *anthropic.Client }

type complianceChatMessageModel struct {
	ID        types.String `tfsdk:"id"`
	Role      types.String `tfsdk:"role"`
	Text      types.String `tfsdk:"text"`
	Model     types.String `tfsdk:"model"`
	CreatedAt types.String `tfsdk:"created_at"`
}

type ComplianceChatMessagesModel struct {
	ChatID   types.String                 `tfsdk:"chat_id"`
	Messages []complianceChatMessageModel `tfsdk:"messages"`
}

func (d *ComplianceChatMessagesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_compliance_chat_messages"
}

func (d *ComplianceChatMessagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads all messages in a specific chat (`/v1/compliance/apps/chats/{id}/messages`). Text content is materialized in state — DO NOT store this state anywhere the messages themselves shouldn't live (LGPD / PII / retention concerns). Enterprise + Compliance Access Key with scope `read:compliance_user_data`.",
		Attributes: map[string]schema.Attribute{
			"chat_id": schema.StringAttribute{Required: true},
			"messages": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":         schema.StringAttribute{Computed: true},
						"role":       schema.StringAttribute{Computed: true},
						"text":       schema.StringAttribute{Computed: true, Sensitive: true},
						"model":      schema.StringAttribute{Computed: true},
						"created_at": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *ComplianceChatMessagesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := clientFromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	d.client = c
}

func (d *ComplianceChatMessagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var cfg ComplianceChatMessagesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	msgs, err := d.client.ListComplianceChatMessages(ctx, cfg.ChatID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to list compliance chat messages", err.Error())
		return
	}
	out := make([]complianceChatMessageModel, 0, len(msgs))
	for _, m := range msgs {
		out = append(out, complianceChatMessageModel{
			ID:        types.StringValue(m.ID),
			Role:      types.StringValue(m.Role),
			Text:      optionalStringValue(m.Text),
			Model:     optionalStringValue(m.Model),
			CreatedAt: types.StringValue(m.CreatedAt),
		})
	}
	cfg.Messages = out
	resp.Diagnostics.Append(resp.State.Set(ctx, cfg)...)
}
