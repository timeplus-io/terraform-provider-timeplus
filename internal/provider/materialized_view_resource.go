// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/timeplus-io/terraform-provider-timeplus/internal/timeplus"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &materializedViewResource{}
var _ resource.ResourceWithImportState = &materializedViewResource{}

func NewMaterializedViewResource() resource.Resource {
	return &materializedViewResource{}
}

// materializedViewResource defines the resource implementation.
type materializedViewResource struct {
	client *timeplus.Client
}

// materializedViewResourceModel describes the materialized view resource data model.
type materializedViewResourceModel struct {
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Query             types.String `tfsdk:"query"`
	TargetStream      types.String `tfsdk:"target_stream"`
	RetentionSize     types.Int64  `tfsdk:"retention_size"`
	RetentionPeriod   types.Int64  `tfsdk:"retention_period"`
	HistoricalDataTTL types.String `tfsdk:"historical_data_ttl"`
}

func (r *materializedViewResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_materialized_view"
}

func (r *materializedViewResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Timeplus materialized views are special views that persist its data. Once created, a materialized view will keep running in the background and continuously writes the query results to the underlying storage system.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The view name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A detailed text describes the view",
				Optional:            true,
			},
			"query": schema.StringAttribute{
				MarkdownDescription: "The query SQL of the view",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_stream": schema.StringAttribute{
				MarkdownDescription: "The optional stream name that the materialized view writes data to",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"retention_size": schema.Int64Attribute{
				MarkdownDescription: "The retention size threadhold in bytes indicates how many data could be kept in the streaming store",
				Optional:            true,
				Computed:            true,
			},
			"retention_period": schema.Int64Attribute{
				MarkdownDescription: "The retention period threadhold in millisecond indicates how long data could be kept in the streaming store",
				Optional:            true,
				Computed:            true,
			},
			"historical_data_ttl": schema.StringAttribute{
				MarkdownDescription: "A SQL expression defines the maximum age of data that are persisted in the historical store",
				Optional:            true,
			},
		},
	}
}

func (r *materializedViewResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*timeplus.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *timeplus.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *materializedViewResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *materializedViewResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	v := timeplus.MaterializedView{
		View: timeplus.View{
			Name:        data.Name.ValueString(),
			Description: data.Description.ValueString(),
			Query:       data.Query.ValueString(),
		},
		RetentionBytes:          int(data.RetentionSize.ValueInt64()),
		RetentionMS:             int(data.RetentionPeriod.ValueInt64()),
		HistoricalTTLExpression: data.HistoricalDataTTL.ValueString(),
	}
	if err := r.client.CreateMaterializedView(&v); err != nil {
		resp.Diagnostics.AddError("Error Creating Materialized View", fmt.Sprintf("Unable to create materialized view %q, got error: %s", v.Name, err))
		return
	}

	// Computed fields
	data.RetentionSize = types.Int64Value(int64(v.RetentionBytes))
	data.RetentionPeriod = types.Int64Value(int64(v.RetentionMS))

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a timeplus_materialized_view resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *materializedViewResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *materializedViewResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	v, err := r.client.GetMaterializedView(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Materialized View", fmt.Sprintf("Unable to read materialized view %q, got error: %s", data.Name.ValueString(), err))
		return
	}

	// required fields
	data.Name = types.StringValue(v.Name)
	data.Query = types.StringValue(v.Query)

	// optional fields
	if !(data.Description.IsNull() && v.Description == "") {
		data.Description = types.StringValue(v.Description)
	}

	if !(data.TargetStream.IsNull() && v.TargetStream == "") {
		data.TargetStream = types.StringValue(v.TargetStream)
	}

	// the create view API will set retention_bytes to -1 if it's not provided
	if !(data.RetentionSize.IsNull() && v.RetentionBytes == -1) {
		data.RetentionSize = types.Int64Value(int64(v.RetentionBytes))
	}

	// the create view API will set retention_ms to -1 if it's not provided
	if !(data.RetentionPeriod.IsNull() && v.RetentionMS == -1) {
		data.RetentionPeriod = types.Int64Value(int64(v.RetentionMS))
	}

	if !(data.HistoricalDataTTL.IsNull() && v.HistoricalTTLExpression == "") {
		data.HistoricalDataTTL = types.StringValue(v.HistoricalTTLExpression)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *materializedViewResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *materializedViewResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	v := timeplus.MaterializedView{
		View: timeplus.View{
			Name:        data.Name.ValueString(),
			Description: data.Description.ValueString(),
			Query:       data.Query.ValueString(),
		},
		RetentionBytes:          int(data.RetentionSize.ValueInt64()),
		RetentionMS:             int(data.RetentionPeriod.ValueInt64()),
		HistoricalTTLExpression: data.HistoricalDataTTL.ValueString(),
	}
	if err := r.client.UpdateMaterializedView(&v); err != nil {
		resp.Diagnostics.AddError("Error Updating Materialized View", fmt.Sprintf("Unable to update materialized view %q, got error: %s", v.Name, err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *materializedViewResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *materializedViewResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteMaterializedView(&timeplus.MaterializedView{
		View: timeplus.View{Name: data.Name.ValueString()},
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Materialized View", fmt.Sprintf("Unable to delete materialized view %q, got error: %s", data.Name.ValueString(), err))
	}
}

func (r *materializedViewResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
