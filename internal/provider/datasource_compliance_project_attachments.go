package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/sauterdigital/terraform-provider-claudeadmin/internal/anthropic"
)

var (
	_ datasource.DataSource              = &ComplianceProjectAttachmentsDataSource{}
	_ datasource.DataSourceWithConfigure = &ComplianceProjectAttachmentsDataSource{}
)

func NewComplianceProjectAttachmentsDataSource() datasource.DataSource {
	return &ComplianceProjectAttachmentsDataSource{}
}

type ComplianceProjectAttachmentsDataSource struct{ client *anthropic.Client }

type complianceProjectAttachmentModel struct {
	ID        types.String `tfsdk:"id"`
	Filename  types.String `tfsdk:"filename"`
	MimeType  types.String `tfsdk:"mime_type"`
	SizeBytes types.Int64  `tfsdk:"size_bytes"`
	CreatedAt types.String `tfsdk:"created_at"`
}

type ComplianceProjectAttachmentsModel struct {
	ProjectID   types.String                       `tfsdk:"project_id"`
	Attachments []complianceProjectAttachmentModel `tfsdk:"attachments"`
}

func (d *ComplianceProjectAttachmentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_compliance_project_attachments"
}

func (d *ComplianceProjectAttachmentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists metadata for files attached to a project (`/v1/compliance/apps/projects/{id}/attachments`). Content download is out of scope for Terraform state — use the Compliance API directly for binaries. Enterprise + Compliance Access Key with scope `read:compliance_user_data`.",
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{Required: true},
			"attachments": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":         schema.StringAttribute{Computed: true},
						"filename":   schema.StringAttribute{Computed: true},
						"mime_type":  schema.StringAttribute{Computed: true},
						"size_bytes": schema.Int64Attribute{Computed: true},
						"created_at": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *ComplianceProjectAttachmentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := clientFromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	d.client = c
}

func (d *ComplianceProjectAttachmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var cfg ComplianceProjectAttachmentsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	atts, err := d.client.ListComplianceProjectAttachments(ctx, cfg.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to list compliance project attachments", err.Error())
		return
	}
	out := make([]complianceProjectAttachmentModel, 0, len(atts))
	for _, a := range atts {
		sz := types.Int64Null()
		if a.SizeBytes != nil {
			sz = types.Int64Value(*a.SizeBytes)
		}
		out = append(out, complianceProjectAttachmentModel{
			ID:        types.StringValue(a.ID),
			Filename:  optionalStringValue(a.Filename),
			MimeType:  optionalStringValue(a.MimeType),
			SizeBytes: sz,
			CreatedAt: types.StringValue(a.CreatedAt),
		})
	}
	cfg.Attachments = out
	resp.Diagnostics.Append(resp.State.Set(ctx, cfg)...)
}
