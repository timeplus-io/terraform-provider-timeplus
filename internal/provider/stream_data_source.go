// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/timeplus-io/terraform-provider-timeplus/internal/timeplus"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &streamDataSource{}

func NewStreamDataSource() datasource.DataSource {
	return &streamDataSource{}
}

// streamDataSource defines the data source implementation.
type streamDataSource struct {
	client *timeplus.Client
}

// streamDataSourceModel describes the data source data model.
type streamDataSourceModel struct {
	Name           types.String  `tfsdk:"name"`
	Description    types.String  `tfsdk:"description"`
	Columns        []columnModel `tfsdk:"columns"`
	RetentionBytes types.Int64   `tfsdk:"retention_bytes"`
	RetentionMS    types.Int64   `tfsdk:"retention_ms"`
	HistoryTTL     types.String  `tfsdk:"history_ttl"`
}

func (d *streamDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stream"
}

func (d *streamDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Timeplus streams are similar to tables in the traditional SQL databases. Both of them are essentially datasets. The key difference is that Timeplus stream is an append-only (by default), unbounded, constantly changing events group.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The stream name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A detailed text describes the stream",
				Computed:            true,
			},
			"columns": schema.ListNestedAttribute{
				MarkdownDescription: "The columns of the stream",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The column name",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type name of the column",
							Computed:            true,
						},
						"default": schema.StringAttribute{
							MarkdownDescription: "The default value for the column",
							Computed:            true,
						},
						"codec": schema.StringAttribute{
							MarkdownDescription: "The codec for value encoding",
							Computed:            true,
						},
					},
				},
			},
			"retention_bytes": schema.Int64Attribute{
				MarkdownDescription: "The retention size threadhold in bytes indicates how many data could be kept in the streaming store",
				Optional:            true,
				Computed:            true,
			},
			"retention_ms": schema.Int64Attribute{
				MarkdownDescription: "The retention period threadhold in millisecond indicates how long data could be kept in the streaming store",
				Optional:            true,
				Computed:            true,
			},
			"history_ttl": schema.StringAttribute{
				MarkdownDescription: "A SQL expression defines the maximum age of data that are persisted in the historical store",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (d *streamDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *streamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *streamDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	s, err := d.client.GetStream(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Stream", fmt.Sprintf("Unable to read stream %q, got error: %s", data.Name.ValueString(), err))
		return
	}

	data.Name = types.StringValue(s.Name)
	data.Description = types.StringValue(s.Description)
	data.RetentionBytes = types.Int64Value(int64(s.RetentionBytes))
	data.RetentionMS = types.Int64Value(int64(s.RetentionMS))
	data.HistoryTTL = types.StringValue(s.HistoricalTTLExpression)

	data.Columns = make([]columnModel, 0, len(s.Columns))
	for i := range s.Columns {
		codec := types.StringValue(strings.TrimSuffix(strings.TrimPrefix(s.Columns[i].Codec, "CODEC("), ")"))
		data.Columns = append(data.Columns, columnModel{
			Name:    types.StringValue(s.Columns[i].Name),
			Type:    types.StringValue(s.Columns[i].Type),
			Default: types.StringValue(s.Columns[i].Default),
			Codec:   codec,
		})
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
