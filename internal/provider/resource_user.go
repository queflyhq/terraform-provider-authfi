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

var _ resource.Resource = &UserResource{}

type UserResource struct {
	client *client.Client
}

type UserResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Email    types.String `tfsdk:"email"`
	Name     types.String `tfsdk:"name"`
	Password types.String `tfsdk:"password"`
}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

func (r *UserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an AuthFI user within a tenant.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Required:    true,
				Description: "User email address.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "User display name.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "User password. Only used on creation.",
			},
		},
	}
}

func (r *UserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	u := &client.User{
		Email:    plan.Email.ValueString(),
		Name:     plan.Name.ValueString(),
		Password: plan.Password.ValueString(),
	}

	created, err := r.client.CreateUser(u)
	if err != nil {
		resp.Diagnostics.AddError("Create User Failed", err.Error())
		return
	}

	plan.ID = types.StringValue(created.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	u, err := r.client.GetUser(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read User Failed", err.Error())
		return
	}
	if u == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Email = types.StringValue(u.Email)
	if u.Name != "" {
		state.Name = types.StringValue(u.Name)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	u := &client.User{
		Email: plan.Email.ValueString(),
		Name:  plan.Name.ValueString(),
	}

	updated, err := r.client.UpdateUser(plan.ID.ValueString(), u)
	if err != nil {
		resp.Diagnostics.AddError("Update User Failed", err.Error())
		return
	}

	plan.Email = types.StringValue(updated.Email)
	if updated.Name != "" {
		plan.Name = types.StringValue(updated.Name)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteUser(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete User Failed", err.Error())
	}
}
