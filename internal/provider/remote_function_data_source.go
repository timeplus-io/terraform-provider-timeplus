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
var _ datasource.DataSource = &remoteFunctionDataSource{}

func NewRemoteFunctionDataSource() datasource.DataSource {
	return &remoteFunctionDataSource{}
}

// remoteFunctionDataSource defines the data source implementation.
type remoteFunctionDataSource struct {
	client *timeplus.Client
}

// remoteFunctionDataSourceModel describes the data source data model.
type remoteFunctionDataSourceModel remoteFunctionResourceModel

func (d *remoteFunctionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_function"
}

func (d *remoteFunctionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Timeplus remote functions are one of the supported user defined function types. They allow users to register a HTTP webhook as a function which can be called in any queries.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The remote function name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A detailed text describes the remote function",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The HTTP endpoint to be used to call the function",
				Computed:            true,
			},
			"auth_header": schema.SingleNestedAttribute{
				MarkdownDescription: "The HTTP header and its value to be used as an authentication means to call the function. The remote function can use this information to determine if it's a valid call",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "The HTTP header name",
						Computed:            true,
					},
					"value": schema.StringAttribute{
						MarkdownDescription: "The value for the header",
						Computed:            true,
					},
				},
			},
			"return_type": schema.StringAttribute{
				MarkdownDescription: "The type of the function's return value",
				Computed:            true,
			},
			"arguments": schema.ListNestedAttribute{
				MarkdownDescription: "The argument names and types the remote function takes",
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

func (d *remoteFunctionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *remoteFunctionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *remoteFunctionDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	s, err := d.client.GetUDF(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading RemoteFunction",
			fmt.Sprintf("Unable to read remote function %q, got error: %s",
				data.Name.ValueString(), err))
		return
	}

	if s.Type != timeplus.UDFTypeRemote {
		resp.Diagnostics.AddError(
			"Error Reading RemoteFunction",
			fmt.Sprintf("Function with name %s is not a remote function",
				data.Name.ValueString()))
		return
	}

	// required fields
	data.Name = types.StringValue(s.Name)
	data.URL = types.StringValue(s.URL)

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

	if !(data.AuthHeader == nil && s.AuthMethod == timeplus.UDFAuthNone) {
		ctx := s.AuthContext
		if ctx == nil {
			ctx = &timeplus.UDFAuthContext{}
		}
		data.AuthHeader = &functionAuthHeaderModel{
			Name:  types.StringValue(ctx.Name),
			Value: types.StringValue(ctx.Value),
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
