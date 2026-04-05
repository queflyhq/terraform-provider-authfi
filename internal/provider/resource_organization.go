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

var _ resource.Resource = &OrganizationResource{}

type OrganizationResource struct {
	client *client.Client
}

type OrganizationResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	EnvType      types.String `tfsdk:"env_type"`
	LogoURL      types.String `tfsdk:"logo_url"`
	PrimaryColor types.String `tfsdk:"primary_color"`
	Tags         types.Map    `tfsdk:"tags"`
}

func NewOrganizationResource() resource.Resource {
	return &OrganizationResource{}
}

func (r *OrganizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (r *OrganizationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an AuthFI organization — a B2B sub-group within a project with optional branding override.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Organization name.",
			},
			"slug": schema.StringAttribute{
				Required:    true,
				Description: "Organization slug (unique within project).",
			},
			"env_type": schema.StringAttribute{
				Optional:    true,
				Description: "Environment type. Values: development, production.",
			},
			"logo_url": schema.StringAttribute{
				Optional:    true,
				Description: "Logo URL override. If not set, inherits from tenant.",
			},
			"primary_color": schema.StringAttribute{
				Optional:    true,
				Description: "Primary color override (hex). If not set, inherits from tenant.",
			},
			"tags": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Key-value tags.",
			},
		},
	}
}

func (r *OrganizationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *OrganizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OrganizationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	o := &client.Organization{
		Name:    plan.Name.ValueString(),
		Slug:    plan.Slug.ValueString(),
		EnvType: plan.EnvType.ValueString(),
	}
	if !plan.LogoURL.IsNull() {
		v := plan.LogoURL.ValueString()
		o.LogoURL = &v
	}
	if !plan.PrimaryColor.IsNull() {
		v := plan.PrimaryColor.ValueString()
		o.Primary = &v
	}

	created, err := r.client.CreateOrganization(o)
	if err != nil {
		resp.Diagnostics.AddError("Create Organization Failed", err.Error())
		return
	}

	plan.ID = types.StringValue(created.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OrganizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OrganizationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	o, err := r.client.GetOrganization(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Organization Failed", err.Error())
		return
	}
	if o == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(o.Name)
	state.Slug = types.StringValue(o.Slug)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OrganizationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	o := &client.Organization{
		Name: plan.Name.ValueString(),
	}
	if !plan.LogoURL.IsNull() {
		v := plan.LogoURL.ValueString()
		o.LogoURL = &v
	}
	if !plan.PrimaryColor.IsNull() {
		v := plan.PrimaryColor.ValueString()
		o.Primary = &v
	}

	updated, err := r.client.UpdateOrganization(plan.ID.ValueString(), o)
	if err != nil {
		resp.Diagnostics.AddError("Update Organization Failed", err.Error())
		return
	}

	plan.Name = types.StringValue(updated.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OrganizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OrganizationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteOrganization(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete Organization Failed", err.Error())
	}
}
