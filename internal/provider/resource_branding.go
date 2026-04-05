package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/queflyhq/terraform-provider-authfi/internal/client"
)

var _ resource.Resource = &BrandingResource{}

type BrandingResource struct {
	client *client.Client
}

type BrandingResourceModel struct {
	Name            types.String `tfsdk:"name"`
	LogoURL         types.String `tfsdk:"logo_url"`
	FaviconURL      types.String `tfsdk:"favicon_url"`
	PrimaryColor    types.String `tfsdk:"primary_color"`
	BackgroundColor types.String `tfsdk:"background_color"`
	TextColor       types.String `tfsdk:"text_color"`
	ButtonRadius    types.String `tfsdk:"button_radius"`
	FontFamily      types.String `tfsdk:"font_family"`
	ThemeMode       types.String `tfsdk:"theme_mode"`
	LoginLayout     types.String `tfsdk:"login_layout"`
	SocialPosition  types.String `tfsdk:"social_position"`
	WelcomeText     types.String `tfsdk:"welcome_text"`
	SubtitleText    types.String `tfsdk:"subtitle_text"`
	SupportEmail    types.String `tfsdk:"support_email"`
	PrivacyURL      types.String `tfsdk:"privacy_url"`
	TermsURL        types.String `tfsdk:"terms_url"`
	CustomCSS       types.String `tfsdk:"custom_css"`
	MFAPolicy       types.String `tfsdk:"mfa_policy"`
	PasswordMin     types.Int64  `tfsdk:"password_min_length"`
	ShowAttribution types.Bool   `tfsdk:"show_attribution"`
}

func NewBrandingResource() resource.Resource {
	return &BrandingResource{}
}

func (r *BrandingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branding"
}

func (r *BrandingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages branding and login UI configuration for an AuthFI tenant.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Tenant display name.",
			},
			"logo_url": schema.StringAttribute{
				Optional:    true,
				Description: "Logo URL.",
			},
			"favicon_url": schema.StringAttribute{
				Optional:    true,
				Description: "Favicon URL.",
			},
			"primary_color": schema.StringAttribute{
				Optional:    true,
				Description: "Primary brand color (hex).",
			},
			"background_color": schema.StringAttribute{
				Optional:    true,
				Description: "Background color (hex).",
			},
			"text_color": schema.StringAttribute{
				Optional:    true,
				Description: "Text color (hex).",
			},
			"button_radius": schema.StringAttribute{
				Optional:    true,
				Description: "Button border radius (e.g. '8px').",
			},
			"font_family": schema.StringAttribute{
				Optional:    true,
				Description: "Font family name.",
			},
			"theme_mode": schema.StringAttribute{
				Optional:    true,
				Description: "Theme mode: light, dark, auto.",
			},
			"login_layout": schema.StringAttribute{
				Optional:    true,
				Description: "Login page layout: centered, split.",
			},
			"social_position": schema.StringAttribute{
				Optional:    true,
				Description: "Social login button position: top, bottom.",
			},
			"welcome_text": schema.StringAttribute{
				Optional:    true,
				Description: "Welcome heading text on login page.",
			},
			"subtitle_text": schema.StringAttribute{
				Optional:    true,
				Description: "Subtitle text on login page.",
			},
			"support_email": schema.StringAttribute{
				Optional:    true,
				Description: "Support email displayed on login page.",
			},
			"privacy_url": schema.StringAttribute{
				Optional:    true,
				Description: "Privacy policy URL.",
			},
			"terms_url": schema.StringAttribute{
				Optional:    true,
				Description: "Terms of service URL.",
			},
			"custom_css": schema.StringAttribute{
				Optional:    true,
				Description: "Custom CSS injected into login page.",
			},
			"mfa_policy": schema.StringAttribute{
				Optional:    true,
				Description: "MFA policy: none, optional, required.",
			},
			"password_min_length": schema.Int64Attribute{
				Optional:    true,
				Description: "Minimum password length.",
			},
			"show_attribution": schema.BoolAttribute{
				Optional:    true,
				Description: "Show 'Powered by AuthFI' attribution.",
			},
		},
	}
}

func (r *BrandingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func strPtr(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	s := v.ValueString()
	return &s
}

func intPtr(v types.Int64) *int {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	i := int(v.ValueInt64())
	return &i
}

func boolPtr(v types.Bool) *bool {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	b := v.ValueBool()
	return &b
}

func (r *BrandingResource) modelToBranding(m *BrandingResourceModel) *client.Branding {
	return &client.Branding{
		Name:            m.Name.ValueString(),
		LogoURL:         strPtr(m.LogoURL),
		FaviconURL:      strPtr(m.FaviconURL),
		PrimaryColor:    strPtr(m.PrimaryColor),
		BackgroundColor: strPtr(m.BackgroundColor),
		TextColor:       strPtr(m.TextColor),
		ButtonRadius:    strPtr(m.ButtonRadius),
		FontFamily:      strPtr(m.FontFamily),
		ThemeMode:       strPtr(m.ThemeMode),
		LoginLayout:     strPtr(m.LoginLayout),
		SocialPosition:  strPtr(m.SocialPosition),
		WelcomeText:     strPtr(m.WelcomeText),
		SubtitleText:    strPtr(m.SubtitleText),
		SupportEmail:    strPtr(m.SupportEmail),
		PrivacyURL:      strPtr(m.PrivacyURL),
		TermsURL:        strPtr(m.TermsURL),
		CustomCSS:       strPtr(m.CustomCSS),
		MFAPolicy:       strPtr(m.MFAPolicy),
		PasswordMin:     intPtr(m.PasswordMin),
		ShowAttribution: boolPtr(m.ShowAttribution),
	}
}

func optionalString(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}

func optionalInt(i *int) types.Int64 {
	if i == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*i))
}

func optionalBool(b *bool) types.Bool {
	if b == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*b)
}

func (r *BrandingResource) brandingToModel(b *client.Branding) BrandingResourceModel {
	m := BrandingResourceModel{
		LogoURL:         optionalString(b.LogoURL),
		FaviconURL:      optionalString(b.FaviconURL),
		PrimaryColor:    optionalString(b.PrimaryColor),
		BackgroundColor: optionalString(b.BackgroundColor),
		TextColor:       optionalString(b.TextColor),
		ButtonRadius:    optionalString(b.ButtonRadius),
		FontFamily:      optionalString(b.FontFamily),
		ThemeMode:       optionalString(b.ThemeMode),
		LoginLayout:     optionalString(b.LoginLayout),
		SocialPosition:  optionalString(b.SocialPosition),
		WelcomeText:     optionalString(b.WelcomeText),
		SubtitleText:    optionalString(b.SubtitleText),
		SupportEmail:    optionalString(b.SupportEmail),
		PrivacyURL:      optionalString(b.PrivacyURL),
		TermsURL:        optionalString(b.TermsURL),
		CustomCSS:       optionalString(b.CustomCSS),
		MFAPolicy:       optionalString(b.MFAPolicy),
		PasswordMin:     optionalInt(b.PasswordMin),
		ShowAttribution: optionalBool(b.ShowAttribution),
	}
	if b.Name != "" {
		m.Name = types.StringValue(b.Name)
	}
	return m
}

func (r *BrandingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan BrandingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateBranding(r.modelToBranding(&plan)); err != nil {
		resp.Diagnostics.AddError("Create Branding Failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BrandingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state BrandingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	b, err := r.client.GetBranding()
	if err != nil {
		resp.Diagnostics.AddError("Read Branding Failed", err.Error())
		return
	}

	model := r.brandingToModel(b)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *BrandingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan BrandingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateBranding(r.modelToBranding(&plan)); err != nil {
		resp.Diagnostics.AddError("Update Branding Failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BrandingResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Branding is a singleton on the tenant; "delete" is a no-op.
	// The tenant still exists; we just remove it from Terraform state.
}
