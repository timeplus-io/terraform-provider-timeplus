// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/timeplus-io/terraform-provider-timeplus/internal/timeplus"
	myValidator "github.com/timeplus-io/terraform-provider-timeplus/internal/validator"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &TimeplusProvider{}

// TimeplusProvider defines the provider implementation.
type TimeplusProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// TimeplusProviderModel describes the provider data model.
type TimeplusProviderModel struct {
	Endpoint  types.String `tfsdk:"endpoint"`
	Workspace types.String `tfsdk:"workspace"`
	ApiKey    types.String `tfsdk:"api_key"`
	Username  types.String `tfsdk:"username"`
	Password  types.String `tfsdk:"password"`
}

func (p *TimeplusProvider) Metadata(ctx context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "timeplus"
	resp.Version = p.version
}

func (p *TimeplusProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The Timeplus provider is used to interact with the resources supported by [Timeplus](https://www.timeplus.com/) in a workspace. The provider needs to be configured with an API key before it can be used.

Use the navigation to the left to read about the available resources.`,
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The base URL endpoint for connecting to the Timeplus workspace. When it's not set, `https://us-west-2.timeplus.cloud` will be used.",
				Optional:            true,
				Validators:          []validator.String{myValidator.URL()},
			},
			"workspace": schema.StringAttribute{
				MarkdownDescription: "The ID of the workspace in which the provider manages resources.",
				Required:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "[Cloud] The API key to be used to call Timeplus Enterprise Cloud.",
				Optional:            true,
				Sensitive:           true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "[Onprem] The username.",
				Optional:            true,
				Sensitive:           false,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "[Onprem] The password.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *TimeplusProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data TimeplusProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	client, err := timeplus.NewClient(data.Workspace.ValueString(), data.ApiKey.ValueString(), data.Username.ValueString(), data.Password.ValueString(), timeplus.ClientOptions{
		BaseURL: data.Endpoint.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("failed to create Timeplus client", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *TimeplusProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewStreamResource,
		NewViewResource,
		NewMaterializedViewResource,
		NewSinkResource,
		NewAlertResource,
		NewSourceResource,
		NewRemoteFunctionResource,
		NewJavascriptFunctionResource,
		NewDashboardResource,
	}
}

func (p *TimeplusProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewStreamDataSource,
		NewViewDataSource,
		NewMaterializedViewDataSource,
		NewSinkDataSource,
		NewAlertDataSource,
		NewSourceDataSource,
		NewRemoteFunctionDataSource,
		NewJavascriptFunctionDataSource,
		NewDashboardDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TimeplusProvider{
			version: version,
		}
	}
}
