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
var _ datasource.DataSource = &javascriptFunctionDataSource{}

func NewJavascriptFunctionDataSource() datasource.DataSource {
	return &javascriptFunctionDataSource{}
}

// javascriptFunctionDataSource defines the data source implementation.
type javascriptFunctionDataSource struct {
	client *timeplus.Client
}

// javascriptFunctionDataSourceModel describes the data source data model.
type javascriptFunctionDataSourceModel javascriptFunctionResourceModel

func (d *javascriptFunctionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_javascript_function"
}

func (d *javascriptFunctionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Timeplus javascript functions are one of the supported user defined function types. Javascript functions allow users to implement functions with the javascript programming language, and be called in queries.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The javascript function name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A detailed text describes the javascript function",
				Computed:            true,
			},
			"is_aggregate_function": schema.BoolAttribute{
				MarkdownDescription: "Indecates if the javascript function an aggregate function",
				Computed:            true,
			},
			"source": schema.StringAttribute{
				MarkdownDescription: "The javascript function source code",
				Computed:            true,
			},
			"return_type": schema.StringAttribute{
				MarkdownDescription: "The type of the function's return value",
				Computed:            true,
			},
			"arg": schema.ListNestedAttribute{
				MarkdownDescription: "The argument names and types the javascript function takes",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The argument name",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The argument type",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *javascriptFunctionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *javascriptFunctionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *javascriptFunctionDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	s, err := d.client.GetUDF(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading JavascriptFunction",
			fmt.Sprintf("Unable to read javascript function %q, got error: %s",
				data.Name.ValueString(), err))
		return
	}

	if s.Type != timeplus.UDFTypeJavascript {
		resp.Diagnostics.AddError(
			"Error Reading JavascriptFunction",
			fmt.Sprintf("Function with name %s is not a javascript function",
				data.Name.ValueString()))
		return
	}

	// required fields
	data.Name = types.StringValue(s.Name)
	data.IsAggrFunction = types.BoolValue(s.IsAggrFunction)
	data.Source = types.StringValue(s.Source)
	data.ReturnType = types.StringValue(s.ReturnType)

	// optional fields
	data.Arguments = make([]functionArgumentModel, 0, len(s.Arguments))
	for i := range s.Arguments {
		data.Arguments = append(data.Arguments, functionArgumentModel{
			Name: types.StringValue(s.Arguments[i].Name),
			Type: types.StringValue(s.Arguments[i].Type),
		})
	}

	if !(data.Description.IsNull() && s.Description == "") {
		data.Description = types.StringValue(s.Description)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
