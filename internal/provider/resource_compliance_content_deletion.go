package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/sauterdigital/terraform-provider-claudeadmin/internal/anthropic"
)

var (
	_ resource.Resource                = &ComplianceContentDeletionResource{}
	_ resource.ResourceWithConfigure   = &ComplianceContentDeletionResource{}
	_ resource.ResourceWithImportState = &ComplianceContentDeletionResource{}
)

func NewComplianceContentDeletionResource() resource.Resource {
	return &ComplianceContentDeletionResource{}
}

type ComplianceContentDeletionResource struct{ client *anthropic.Client }

type ComplianceContentDeletionModel struct {
	ID         types.String `tfsdk:"id"`
	TargetType types.String `tfsdk:"target_type"`
	TargetID   types.String `tfsdk:"target_id"`
	Reason     types.String `tfsdk:"reason"`
}

func (r *ComplianceContentDeletionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_compliance_content_deletion"
}

func (r *ComplianceContentDeletionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Performs a hard-delete of user-generated content (eDiscovery / DLP / GDPR erasure) via the Compliance Content API. **One-way, irreversible operation.** Requires Compliance Access Key with scope `delete:compliance_user_data`. Changing `target_type` / `target_id` / `reason` triggers a replace, which would attempt to delete again — usually failing because the content is already gone. Removing from state is a no-op.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "Composite identifier `<target_type>:<target_id>`.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"target_type": schema.StringAttribute{
				Description:   "What to delete. One of: `chat`, `chat_file`, `chat_generated_file`, `project`, `project_document`.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					stringvalidator.OneOf("chat", "chat_file", "chat_generated_file", "project", "project_document"),
				},
			},
			"target_id": schema.StringAttribute{
				Description:   "ID of the resource to delete (chat id, file id, project id, etc).",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"reason": schema.StringAttribute{
				Description:   "Optional free-form justification for audit purposes (LGPD art. 18 erasure request, DLP incident ticket, etc). Not sent to Anthropic — stored in state only.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *ComplianceContentDeletionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := clientFromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	r.client = c
}

func (r *ComplianceContentDeletionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ComplianceContentDeletionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteComplianceContent(ctx, plan.TargetType.ValueString(), plan.TargetID.ValueString()); err != nil {
		if !anthropic.IsNotFound(err) {
			resp.Diagnostics.AddError("Failed to delete compliance content", err.Error())
			return
		}
	}
	plan.ID = types.StringValue(plan.TargetType.ValueString() + ":" + plan.TargetID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ComplianceContentDeletionResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	// One-way action. No corresponding GET to refresh state — success is
	// recorded in state at Create time and not re-verified.
}

func (r *ComplianceContentDeletionResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// All mutable attributes RequiresReplace; Update unreachable.
}

func (r *ComplianceContentDeletionResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// No-op: the content is already gone in the API. Terraform state removal
	// is the only side effect.
}

func (r *ComplianceContentDeletionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
