package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api_v0 "github.com/hashicorp/terraform-provider-scaffolding-framework/api/v0"
)

// Ensure FelderaProvider satisfies various provider interfaces.
var _ provider.Provider = &FelderaProvider{}

// FelderaProvider defines the provider implementation.
type FelderaProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ScaffoldingProviderModel describes the provider data model.
type ScaffoldingProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

type FelderaProviderData struct {
	Client   api_v0.ClientWithResponsesInterface
	Endpoint string
}

func (p *FelderaProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "feldera"
	resp.Version = p.version
}

func (p *FelderaProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The endpoint to use for the Feldera API",
				Required:            true,
			},
		},
	}
}

func (p *FelderaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ScaffoldingProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Endpoint.IsNull() {
		resp.Diagnostics.AddError(
			"endpoint is required",
			"The endpoint must be set")
		return
	}

	client, err := api_v0.NewClientWithResponses(data.Endpoint.ValueString() + "/v0")
	if err != nil {
		resp.Diagnostics.AddError(
			"error creating client",
			fmt.Sprintf("There was an error while creating openAPI client for Feldera API, error: %s", err))
		return
	}

	respData := &FelderaProviderData{
		Client:   client,
		Endpoint: data.Endpoint.ValueString(),
	}
	resp.DataSourceData = respData
	resp.ResourceData = respData
}

func (p *FelderaProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProgramResource,
	}
}

func (p *FelderaProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &FelderaProvider{
			version: version,
		}
	}
}
