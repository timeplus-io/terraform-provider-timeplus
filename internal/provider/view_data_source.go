// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/timeplus-io/terraform-provider-timeplus/internal/timeplus"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &viewDataSource{}

func NewViewDataSource() datasource.DataSource {
	return &viewDataSource{}
}

// viewDataSource defines the data source implementation.
type viewDataSource struct {
	client *timeplus.Client
}

// viewDataSourceModel describes the data source data model.
type viewDataSourceModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Query       types.String `tfsdk:"query"`
}

func (d *viewDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_view"
}

func (d *viewDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Timeplus views are named queries. When you create a view, you basically create a query and assign a name to the query. Therefore, a view is useful for wrapping a commonly used complex query.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The view name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A detailed text describes the view",
				Computed:            true,
			},
			"query": schema.StringAttribute{
				MarkdownDescription: "The query SQL of the view",
				Computed:            true,
			},
		},
	}
}

func (d *viewDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*timeplus.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *timeplus.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *viewDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *viewDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	v, err := d.client.GetView(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading View", fmt.Sprintf("Unable to read view %q, got error: %s", data.Name.ValueString(), err))
		return
	}

	data.Name = types.StringValue(v.Name)
	data.Description = types.StringValue(v.Description)
	data.Query = types.StringValue(v.Query)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
