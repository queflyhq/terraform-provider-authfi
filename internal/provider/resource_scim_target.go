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

var _ resource.Resource = &SCIMTargetResource{}

type SCIMTargetResource struct {
	client *client.Client
}

type SCIMTargetResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
}

func NewSCIMTargetResource() resource.Resource {
	return &SCIMTargetResource{}
}

func (r *SCIMTargetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scim_target"
}

func (r *SCIMTargetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a SCIM provisioning target in AuthFI.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "SCIM target name.",
			},
			"endpoint": schema.StringAttribute{
				Required:    true,
				Description: "SCIM endpoint URL.",
			},
			"token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Bearer token for SCIM endpoint authentication.",
			},
		},
	}
}

func (r *SCIMTargetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *SCIMTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SCIMTargetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	s := &client.SCIMTarget{
		Name:     plan.Name.ValueString(),
		Endpoint: plan.Endpoint.ValueString(),
		Token:    plan.Token.ValueString(),
	}

	created, err := r.client.CreateSCIMTarget(s)
	if err != nil {
		resp.Diagnostics.AddError("Create SCIM Target Failed", err.Error())
		return
	}

	plan.ID = types.StringValue(created.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SCIMTargetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SCIMTargetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	s, err := r.client.GetSCIMTarget(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read SCIM Target Failed", err.Error())
		return
	}
	if s == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(s.Name)
	state.Endpoint = types.StringValue(s.Endpoint)
	// Token is write-only; don't overwrite from API response.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SCIMTargetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update Not Supported", "SCIM targets cannot be updated via Terraform yet. Delete and recreate.")
}

func (r *SCIMTargetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SCIMTargetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSCIMTarget(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete SCIM Target Failed", err.Error())
	}
}
