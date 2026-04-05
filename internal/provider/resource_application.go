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

var _ resource.Resource = &ApplicationResource{}

type ApplicationResource struct {
	client *client.Client
}

type ApplicationResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	AppType      types.String `tfsdk:"app_type"`
	CallbackURLs types.List   `tfsdk:"callback_urls"`
	LogoutURLs   types.List   `tfsdk:"logout_urls"`
	Origins      types.List   `tfsdk:"allowed_origins"`
}

func NewApplicationResource() resource.Resource {
	return &ApplicationResource{}
}

func (r *ApplicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *ApplicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an AuthFI application — an OAuth2/OIDC client with credentials and redirect URIs.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Application name.",
			},
			"client_id": schema.StringAttribute{
				Computed:    true,
				Description: "OAuth2 client ID (auto-generated).",
			},
			"client_secret": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "OAuth2 client secret (auto-generated, shown once).",
			},
			"app_type": schema.StringAttribute{
				Optional:    true,
				Description: "Application type. Values: regular_web, spa, native, m2m. Default: regular_web.",
			},
			"callback_urls": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Allowed callback/redirect URIs.",
			},
			"logout_urls": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Allowed post-logout redirect URIs.",
			},
			"allowed_origins": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Allowed CORS origins.",
			},
		},
	}
}

func (r *ApplicationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *ApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ApplicationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	a := &client.Application{
		Name:    plan.Name.ValueString(),
		AppType: plan.AppType.ValueString(),
	}
	if a.AppType == "" {
		a.AppType = "regular_web"
	}

	created, err := r.client.CreateApplication(a)
	if err != nil {
		resp.Diagnostics.AddError("Create Application Failed", err.Error())
		return
	}

	plan.ID = types.StringValue(created.ID)
	plan.ClientID = types.StringValue(created.ClientID)
	plan.ClientSecret = types.StringValue(created.ClientSecret)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ApplicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	a, err := r.client.GetApplication(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Application Failed", err.Error())
		return
	}
	if a == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(a.Name)
	state.ClientID = types.StringValue(a.ClientID)
	state.AppType = types.StringValue(a.AppType)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ApplicationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	a := &client.Application{
		Name:    plan.Name.ValueString(),
		AppType: plan.AppType.ValueString(),
	}

	updated, err := r.client.UpdateApplication(plan.ID.ValueString(), a)
	if err != nil {
		resp.Diagnostics.AddError("Update Application Failed", err.Error())
		return
	}

	plan.ClientID = types.StringValue(updated.ClientID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ApplicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteApplication(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete Application Failed", err.Error())
	}
}
