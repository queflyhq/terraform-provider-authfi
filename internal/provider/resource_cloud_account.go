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

var _ resource.Resource = &CloudAccountResource{}

type CloudAccountResource struct {
	client *client.Client
}

type CloudAccountResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	ProviderType types.String `tfsdk:"provider_type"`
	AccessType   types.String `tfsdk:"access_type"`
	Config       types.Map    `tfsdk:"config"`
}

func NewCloudAccountResource() resource.Resource {
	return &CloudAccountResource{}
}

func (r *CloudAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_account"
}

func (r *CloudAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a cloud account integration in AuthFI (GCP, AWS, Azure).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Cloud account name.",
			},
			"provider_type": schema.StringAttribute{
				Required:    true,
				Description: "Cloud provider: gcp, aws, azure.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"access_type": schema.StringAttribute{
				Optional:    true,
				Description: "Access type: wif, service_account, role_arn.",
			},
			"config": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Provider-specific configuration key-value pairs.",
			},
		},
	}
}

func (r *CloudAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *CloudAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CloudAccountResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	a := &client.CloudAccount{
		Name:         plan.Name.ValueString(),
		ProviderType: plan.ProviderType.ValueString(),
		AccessType:   plan.AccessType.ValueString(),
	}

	if !plan.Config.IsNull() {
		cfg := make(map[string]string)
		resp.Diagnostics.Append(plan.Config.ElementsAs(ctx, &cfg, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		a.Config = cfg
	}

	created, err := r.client.CreateCloudAccount(a)
	if err != nil {
		resp.Diagnostics.AddError("Create Cloud Account Failed", err.Error())
		return
	}

	plan.ID = types.StringValue(created.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CloudAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CloudAccountResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	a, err := r.client.GetCloudAccount(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Cloud Account Failed", err.Error())
		return
	}
	if a == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(a.Name)
	state.ProviderType = types.StringValue(a.ProviderType)
	if a.AccessType != "" {
		state.AccessType = types.StringValue(a.AccessType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CloudAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update Not Supported", "Cloud accounts cannot be updated via Terraform yet. Delete and recreate.")
}

func (r *CloudAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CloudAccountResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteCloudAccount(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete Cloud Account Failed", err.Error())
	}
}
