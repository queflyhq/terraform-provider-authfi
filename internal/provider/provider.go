package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/queflyhq/terraform-provider-authfi/internal/client"
)

var _ provider.Provider = &AuthFIProvider{}

type AuthFIProvider struct {
	version string
}

type AuthFIProviderModel struct {
	APIKey types.String `tfsdk:"api_key"`
	Tenant types.String `tfsdk:"tenant"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AuthFIProvider{version: version}
	}
}

func (p *AuthFIProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "authfi"
	resp.Version = p.version
}

func (p *AuthFIProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The AuthFI provider manages identity infrastructure — projects, organizations, applications, connections, roles, and users.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "AuthFI Management API key (sk_xxx). Can also be set via AUTHFI_API_KEY env var.",
				Optional:    true,
				Sensitive:   true,
			},
			"tenant": schema.StringAttribute{
				Description: "Tenant slug (e.g. 'ayush'). Can also be set via AUTHFI_TENANT env var.",
				Optional:    true,
			},
		},
	}
}

func (p *AuthFIProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config AuthFIProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := os.Getenv("AUTHFI_API_KEY")
	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}
	if apiKey == "" {
		resp.Diagnostics.AddError("Missing API Key", "Set api_key in provider config or AUTHFI_API_KEY env var")
		return
	}

	tenant := os.Getenv("AUTHFI_TENANT")
	if !config.Tenant.IsNull() {
		tenant = config.Tenant.ValueString()
	}
	if tenant == "" {
		resp.Diagnostics.AddError("Missing Tenant", "Set tenant in provider config or AUTHFI_TENANT env var")
		return
	}

	c := client.New(apiKey, tenant)
	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *AuthFIProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
		NewOrganizationResource,
		NewApplicationResource,
		NewConnectionResource,
		NewRoleResource,
	}
}

func (p *AuthFIProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
