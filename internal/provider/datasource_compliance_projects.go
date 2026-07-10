package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/sauterdigital/terraform-provider-claudeadmin/internal/anthropic"
)

var (
	_ datasource.DataSource              = &ComplianceProjectsDataSource{}
	_ datasource.DataSourceWithConfigure = &ComplianceProjectsDataSource{}
)

func NewComplianceProjectsDataSource() datasource.DataSource { return &ComplianceProjectsDataSource{} }

type ComplianceProjectsDataSource struct{ client *anthropic.Client }

type complianceProjectModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	UserID      types.String `tfsdk:"user_id"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	ChatCount   types.Int64  `tfsdk:"chat_count"`
}

type ComplianceProjectsModel struct {
	UserID   types.String             `tfsdk:"user_id"`
	Limit    types.Int64              `tfsdk:"limit"`
	Projects []complianceProjectModel `tfsdk:"projects"`
}

func (d *ComplianceProjectsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_compliance_projects"
}

func (d *ComplianceProjectsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists Claude Projects (`/v1/compliance/apps/projects`). Enterprise + Compliance Access Key with scope `read:compliance_user_data`.",
		Attributes: map[string]schema.Attribute{
			"user_id": schema.StringAttribute{Optional: true, Description: "Filter by owning user_id."},
			"limit":   schema.Int64Attribute{Optional: true, Description: "Per-page size (default 100)."},
			"projects": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.StringAttribute{Computed: true},
						"name":        schema.StringAttribute{Computed: true},
						"user_id":     schema.StringAttribute{Computed: true},
						"description": schema.StringAttribute{Computed: true},
						"created_at":  schema.StringAttribute{Computed: true},
						"updated_at":  schema.StringAttribute{Computed: true},
						"chat_count":  schema.Int64Attribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *ComplianceProjectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := clientFromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	d.client = c
}

func (d *ComplianceProjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var cfg ComplianceProjectsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	projs, err := d.client.ListComplianceProjects(ctx, anthropic.ListComplianceProjectsParams{
		UserID: cfg.UserID.ValueString(),
		Limit:  int(cfg.Limit.ValueInt64()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to list compliance projects", err.Error())
		return
	}
	out := make([]complianceProjectModel, 0, len(projs))
	for _, p := range projs {
		cc := types.Int64Null()
		if p.ChatCount != nil {
			cc = types.Int64Value(*p.ChatCount)
		}
		out = append(out, complianceProjectModel{
			ID:          types.StringValue(p.ID),
			Name:        types.StringValue(p.Name),
			UserID:      types.StringValue(p.UserID),
			Description: optionalStringValue(p.Description),
			CreatedAt:   types.StringValue(p.CreatedAt),
			UpdatedAt:   optionalStringValue(p.UpdatedAt),
			ChatCount:   cc,
		})
	}
	cfg.Projects = out
	resp.Diagnostics.Append(resp.State.Set(ctx, cfg)...)
}
