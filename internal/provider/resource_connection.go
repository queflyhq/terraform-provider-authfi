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

var _ resource.Resource = &ConnectionResource{}

type ConnectionResource struct {
	client *client.Client
}

type ConnectionResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Provider types.String `tfsdk:"provider_type"`
	Type     types.String `tfsdk:"type"`
	ClientID types.String `tfsdk:"client_id"`
}

func NewConnectionResource() resource.Resource {
	return &ConnectionResource{}
}

func (r *ConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connection"
}

func (r *ConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an AuthFI SSO connection — Google, GitHub, SAML, OIDC, etc.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Connection name.",
			},
			"provider_type": schema.StringAttribute{
				Required:    true,
				Description: "SSO provider. Values: google, github, microsoft, saml, oidc, apple, gitlab.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "Connection type. Values: social, enterprise. Default: social.",
			},
			"client_id": schema.StringAttribute{
				Optional:    true,
				Description: "OAuth client ID (for social providers).",
			},
		},
	}
}

func (r *ConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *ConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ConnectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := &client.Connection{
		Name:     plan.Name.ValueString(),
		Provider: plan.Provider.ValueString(),
		Type:     plan.Type.ValueString(),
		ClientID: plan.ClientID.ValueString(),
	}

	created, err := r.client.CreateConnection(c)
	if err != nil {
		resp.Diagnostics.AddError("Create Connection Failed", err.Error())
		return
	}

	plan.ID = types.StringValue(created.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ConnectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, err := r.client.GetConnection(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Connection Failed", err.Error())
		return
	}
	if c == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(c.Name)
	state.Provider = types.StringValue(c.Provider)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update Not Supported", "Connections cannot be updated. Delete and recreate instead.")
}

func (r *ConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ConnectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteConnection(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete Connection Failed", err.Error())
	}
}
