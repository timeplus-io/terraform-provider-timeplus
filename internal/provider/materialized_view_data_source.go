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
var _ datasource.DataSource = &materializedViewDataSource{}

func NewMaterializedViewDataSource() datasource.DataSource {
	return &materializedViewDataSource{}
}

// materializedViewDataSource defines the data source implementation.
type materializedViewDataSource struct {
	client *timeplus.Client
}

// materializedViewDataSourceModel describes the data source data model.
type materializedViewDataSourceModel struct {
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Query          types.String `tfsdk:"query"`
	TargetStream   types.String `tfsdk:"target_stream"`
	RetentionBytes types.Int64  `tfsdk:"retention_bytes"`
	RetentionMS    types.Int64  `tfsdk:"retention_ms"`
	HistoryTTL     types.String `tfsdk:"history_ttl"`
}

func (d *materializedViewDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_materialized_view"
}

func (d *materializedViewDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Timeplus materialized views are special views that persist its data. Once created, a materialized view will keep running in the background and continuously writes the query results to the underlying storage system.",

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
			"target_stream": schema.StringAttribute{
				MarkdownDescription: "The optional stream name that the materialized view writes data to",
				Computed:            true,
			},
			"retention_bytes": schema.Int64Attribute{
				MarkdownDescription: "The retention size threadhold in bytes indicates how many data could be kept in the streaming store",
				Computed:            true,
			},
			"retention_ms": schema.Int64Attribute{
				MarkdownDescription: "The retention period threadhold in millisecond indicates how long data could be kept in the streaming store",
				Computed:            true,
			},
			"history_ttl": schema.StringAttribute{
				MarkdownDescription: "A SQL expression defines the maximum age of historical data",
				Computed:            true,
			},
		},
	}
}

func (d *materializedViewDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *materializedViewDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *materializedViewDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	v, err := d.client.GetMaterializedView(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Materialized View", fmt.Sprintf("Unable to read materialized view %q, got error: %s", data.Name.ValueString(), err))
		return
	}

	data.Name = types.StringValue(v.Name)
	data.Description = types.StringValue(v.Description)
	data.Query = types.StringValue(v.Query)
	data.TargetStream = types.StringValue(v.TargetStream)
	data.RetentionBytes = types.Int64Value(int64(v.RetentionBytes))
	data.RetentionMS = types.Int64Value(int64(v.RetentionMS))
	data.HistoryTTL = types.StringValue(v.TTLExpression)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
