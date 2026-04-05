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

var _ resource.Resource = &DomainResource{}

type DomainResource struct {
	client *client.Client
}

type DomainResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Domain   types.String `tfsdk:"domain"`
	Type     types.String `tfsdk:"type"`
	Verified types.Bool   `tfsdk:"verified"`
}

func NewDomainResource() resource.Resource {
	return &DomainResource{}
}

func (r *DomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (r *DomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a custom domain for an AuthFI tenant.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain": schema.StringAttribute{
				Required:    true,
				Description: "Custom domain name (e.g. 'auth.example.com').",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "Domain type: custom, subdomain.",
			},
			"verified": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the domain has been verified.",
			},
		},
	}
}

func (r *DomainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *DomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DomainResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	d := &client.Domain{
		Domain: plan.Domain.ValueString(),
		Type:   plan.Type.ValueString(),
	}

	created, err := r.client.CreateDomain(d)
	if err != nil {
		resp.Diagnostics.AddError("Create Domain Failed", err.Error())
		return
	}

	plan.ID = types.StringValue(created.ID)
	plan.Verified = types.BoolValue(created.Verified)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DomainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	d, err := r.client.GetDomain()
	if err != nil {
		resp.Diagnostics.AddError("Read Domain Failed", err.Error())
		return
	}
	if d == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Domain = types.StringValue(d.Domain)
	if d.Type != "" {
		state.Type = types.StringValue(d.Type)
	}
	state.Verified = types.BoolValue(d.Verified)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update Not Supported", "Domains cannot be updated. Delete and recreate.")
}

func (r *DomainResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Domain is a singleton on the tenant; removal from state only.
	// A dedicated DELETE endpoint would be needed for full removal.
}
