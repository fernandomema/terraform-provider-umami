package provider

import (
	"context"
	"os"

	"github.com/fernandomema/terraform-provider-umami/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type umamiProvider struct{ version string }

type umamiProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	APIKey   types.String `tfsdk:"api_key"`
	Insecure types.Bool   `tfsdk:"insecure"`
}

// Data compartida con recursos: cliente y endpoint (para construir el snippet).
type providerData struct {
	Client   *client.Client
	Endpoint string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider { return &umamiProvider{version: version} }
}

func (p *umamiProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "umami"
	resp.Version = p.version
}

func (p *umamiProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an Umami instance (self-hosted or Cloud).",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{Optional: true, Description: "Base URL of the instance, e.g. https://umami.example.com. Environment variable: UMAMI_ENDPOINT."},
			"username": schema.StringAttribute{Optional: true, Description: "Username (self-hosted). Environment variable: UMAMI_USERNAME."},
			"password": schema.StringAttribute{Optional: true, Sensitive: true, Description: "Password (self-hosted). Environment variable: UMAMI_PASSWORD."},
			"api_key":  schema.StringAttribute{Optional: true, Sensitive: true, Description: "Umami Cloud API key, an alternative to username/password. Environment variable: UMAMI_API_KEY."},
			"insecure": schema.BoolAttribute{Optional: true, Description: "Skip TLS certificate verification."},
		},
	}
}

func orEnv(v types.String, env string) string {
	if !v.IsNull() && !v.IsUnknown() && v.ValueString() != "" {
		return v.ValueString()
	}
	return os.Getenv(env)
}

func (p *umamiProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg umamiProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	endpoint := orEnv(cfg.Endpoint, "UMAMI_ENDPOINT")
	user, pass, key := orEnv(cfg.Username, "UMAMI_USERNAME"), orEnv(cfg.Password, "UMAMI_PASSWORD"), orEnv(cfg.APIKey, "UMAMI_API_KEY")
	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(path.Root("endpoint"), "Missing endpoint", "Set endpoint or the UMAMI_ENDPOINT environment variable.")
	}
	if key == "" && (user == "" || pass == "") {
		resp.Diagnostics.AddError("Missing credentials", "Set api_key, or both username and password (or the UMAMI_* environment variables).")
	}
	if resp.Diagnostics.HasError() {
		return
	}
	c, err := client.New(ctx, client.Config{Endpoint: endpoint, Username: user, Password: pass, APIKey: key, Insecure: cfg.Insecure.ValueBool()})
	if err != nil {
		resp.Diagnostics.AddError("Unable to connect to Umami", err.Error())
		return
	}
	d := &providerData{Client: c, Endpoint: endpoint}
	resp.ResourceData, resp.DataSourceData = d, d
}

func (p *umamiProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{NewWebsiteResource}
}

func (p *umamiProvider) DataSources(context.Context) []func() datasource.DataSource { return nil }
