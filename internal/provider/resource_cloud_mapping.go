package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/queflyhq/terraform-provider-authfi/internal/client"
)

var _ resource.Resource = &CloudMappingResource{}

type CloudMappingResource struct {
	client *client.Client
}

type CloudMappingResourceModel struct {
	ID             types.String `tfsdk:"id"`
	CloudAccountID types.String `tfsdk:"cloud_account_id"`
	RoleID         types.String `tfsdk:"role_id"`
	CloudRole      types.String `tfsdk:"cloud_role"`
	Description    types.String `tfsdk:"description"`
}

func NewCloudMappingResource() resource.Resource {
	return &CloudMappingResource{}
}

func (r *CloudMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_mapping"
}

func (r *CloudMappingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Maps an AuthFI role to a cloud provider role for federated access.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloud_account_id": schema.StringAttribute{
				Required:    true,
				Description: "Cloud account ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role_id": schema.StringAttribute{
				Required:    true,
				Description: "AuthFI role ID to map.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cloud_role": schema.StringAttribute{
				Required:    true,
				Description: "Cloud provider role ARN or identifier.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of this mapping.",
			},
		},
	}
}

func (r *CloudMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *CloudMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CloudMappingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	m := &client.CloudMapping{
		CloudAccountID: plan.CloudAccountID.ValueString(),
		RoleID:         plan.RoleID.ValueString(),
		CloudRole:      plan.CloudRole.ValueString(),
		Description:    plan.Description.ValueString(),
	}

	created, err := r.client.CreateCloudMapping(m)
	if err != nil {
		resp.Diagnostics.AddError("Create Cloud Mapping Failed", err.Error())
		return
	}

	plan.ID = types.StringValue(created.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CloudMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CloudMappingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	m, err := r.client.GetCloudMapping(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Cloud Mapping Failed", err.Error())
		return
	}
	if m == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.CloudAccountID = types.StringValue(m.CloudAccountID)
	state.RoleID = types.StringValue(m.RoleID)
	state.CloudRole = types.StringValue(m.CloudRole)
	if m.Description != "" {
		state.Description = types.StringValue(m.Description)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CloudMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update Not Supported", "Cloud mappings cannot be updated via Terraform yet. Delete and recreate.")
}

func (r *CloudMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CloudMappingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteCloudMapping(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete Cloud Mapping Failed", err.Error())
	}
}
