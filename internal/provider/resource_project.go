package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/queflyhq/terraform-provider-authfi/internal/client"
)

var _ resource.Resource = &ProjectResource{}

type ProjectResource struct {
	client *client.Client
}

type ProjectResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Slug    types.String `tfsdk:"slug"`
	Region  types.String `tfsdk:"region"`
	EnvType types.String `tfsdk:"env_type"`
	Tags    types.Map    `tfsdk:"tags"`
}

func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

func (r *ProjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an AuthFI project — a workspace with data residency where users, apps, and connections live.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Project ID (UUID).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Project name.",
			},
			"slug": schema.StringAttribute{
				Computed:    true,
				Description: "Project slug (auto-generated from name). Used in URLs: {slug}.authfi.app",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region": schema.StringAttribute{
				Required:    true,
				Description: "Data residency region. Cannot be changed after creation. Values: in, us, eu, au.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"env_type": schema.StringAttribute{
				Optional:    true,
				Description: "Environment type marker. Affects rate limits and security defaults. Values: development, production. Default: production.",
			},
			"tags": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Key-value tags for labeling and filtering.",
			},
		},
	}
}

func (r *ProjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	p := &client.Project{
		Name:    plan.Name.ValueString(),
		Region:  plan.Region.ValueString(),
		EnvType: plan.EnvType.ValueString(),
	}

	created, err := r.client.CreateProject(p)
	if err != nil {
		resp.Diagnostics.AddError("Create Project Failed", err.Error())
		return
	}

	plan.ID = types.StringValue(created.ID)
	plan.Slug = types.StringValue(created.Slug)
	if created.EnvType != "" {
		plan.EnvType = types.StringValue(created.EnvType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	p, err := r.client.GetProject(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Project Failed", err.Error())
		return
	}
	if p == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(p.Name)
	state.Slug = types.StringValue(p.Slug)
	if p.Region != "" {
		state.Region = types.StringValue(p.Region)
	}
	if p.EnvType != "" {
		state.EnvType = types.StringValue(p.EnvType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update Not Supported", "Projects cannot be updated. Delete and recreate instead.")
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteProject(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete Project Failed", err.Error())
	}
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	p, err := r.client.GetProject(req.ID)
	if err != nil || p == nil {
		resp.Diagnostics.AddError("Import Failed", fmt.Sprintf("Project %s not found", req.ID))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, types.ObjectType{}.AttributePath(), &ProjectResourceModel{
		ID:   types.StringValue(p.ID),
		Name: types.StringValue(p.Name),
		Slug: types.StringValue(p.Slug),
	})...)
}
