// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/timeplus-io/terraform-provider-timeplus/internal/timeplus"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &alertDataSource{}

func NewAlertDataSource() datasource.DataSource {
	return &alertDataSource{}
}

// alertDataSource defines the data source implementation.
type alertDataSource struct {
	client *timeplus.Client
}

// alertDataSourceModel describes the data source data model.
type alertDataSourceModel alertResourceModel

func (d *alertDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert"
}

func (d *alertDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Timeplus alerts run queries in background and send query results to the target system continuously.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The alert immutable ID, generated by Timeplus",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The human-friendly name for the alert",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A detailed text describes the alert",
				Computed:            true,
			},
			"severity": schema.Int64Attribute{
				MarkdownDescription: "A number indicates how serious this alert is",
				Computed:            true,
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "The type of action the alert should take, i.e. the name of the target system, like 'slack', 'email', etc. Please refer to the Timeplus document for supported alert action types",
				Computed:            true,
			},
			// since Terraform does not have built-in support for map[string]any with the framework library, we use JSON as a simple solution
			"properties": schema.StringAttribute{
				MarkdownDescription: "a JSON object defines the configurations for the specific alert action. The properites could contain sensitive information like password, secret, etc.",
				Computed:            true,
				Sensitive:           true,
			},
			"trigger_sql": schema.StringAttribute{
				MarkdownDescription: "The query the alert uses to generate events that trigger the alert",
				Computed:            true,
			},
			"resolve_sql": schema.StringAttribute{
				MarkdownDescription: "The query the alert uses to generate events that resolve the alert",
				Computed:            true,
			},
		},
	}
}

func (d *alertDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *alertDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *alertDataSourceModel

	// Read Terraform prior state data into the model
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	s, err := d.client.GetAlert(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Alert",
			fmt.Sprintf("Unable to read alert %q, got error: %s",
				data.Name.ValueString(), err))
		return
	}

	data.Name = types.StringValue(s.Name)
	data.Description = types.StringValue(s.Description)
	data.Severity = types.Int64Value(int64(s.Severity))
	data.Action = types.StringValue(s.Action)
	data.TriggerSQL = types.StringValue(s.TriggerSQL)
	data.ResolveSQL = types.StringValue(s.ResolveSQL)

	propsBytes, err := json.Marshal(s.Properties)
	if err != nil {
		resp.Diagnostics.AddError("Bad Alert Properties", fmt.Sprintf("Unable to encode alert properties into JSON, got error: %s", err))
		return
	}
	data.Properties = types.StringValue(string(propsBytes))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
