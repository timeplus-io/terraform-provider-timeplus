// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/timeplus-io/terraform-provider-timeplus/internal/timeplus"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &streamResource{}
var _ resource.ResourceWithImportState = &streamResource{}

func NewStreamResource() resource.Resource {
	return &streamResource{}
}

// streamResource defines the resource implementation.
type streamResource struct {
	client *timeplus.Client
}

type columnModel struct {
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	Default         types.String `tfsdk:"default"`
	Codec           types.String `tfsdk:"codec"`
	CodecExpression types.String `tfsdk:"codec_expression"`
}

// streamResourceModel describes the stream resource data model.
type streamResourceModel struct {
	Name              types.String  `tfsdk:"name"`
	Description       types.String  `tfsdk:"description"`
	Columns           []columnModel `tfsdk:"column"`
	RetentionSize     types.Int64   `tfsdk:"retention_size"`
	RetentionPeriod   types.Int64   `tfsdk:"retention_period"`
	HistoricalDataTTL types.String  `tfsdk:"historical_data_ttl"`
}

func (r *streamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stream"
}

func (r *streamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Optional:            true,
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
		Blocks: map[string]schema.Block{
			"column": schema.ListNestedBlock{
				MarkdownDescription: "Define a column of the stream",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The column name",
							Required:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type name of the column",
							Required:            true,
						},
						"default": schema.StringAttribute{
							MarkdownDescription: "The default value for the column",
							Optional:            true,
							Computed:            true,
							Default:             stringdefault.StaticString(""),
						},
						"codec": schema.StringAttribute{
							MarkdownDescription: "The codec for value encoding",
							Optional:            true,
							Computed:            true,
							Default:             stringdefault.StaticString(""),
						},
						"codec_expression": schema.StringAttribute{
							MarkdownDescription: "The computed codec",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (r *streamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *streamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *streamResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	columns := make([]timeplus.Column, 0, len(data.Columns))
	for i := range data.Columns {
		columns = append(columns, timeplus.Column{
			Name:             data.Columns[i].Name.ValueString(),
			Type:             data.Columns[i].Type.ValueString(),
			Default:          data.Columns[i].Default.ValueString(),
			CompressionCodec: data.Columns[i].Codec.ValueString(),
		})
	}

	s := timeplus.Stream{
		Name:                    data.Name.ValueString(),
		Description:             data.Description.ValueString(),
		Columns:                 columns,
		RetentionBytes:          int(data.RetentionSize.ValueInt64()),
		RetentionMS:             int(data.RetentionPeriod.ValueInt64()),
		HistoricalTTLExpression: data.HistoricalDataTTL.ValueString(),
	}
	if err := r.client.CreateStream(&s); err != nil {
		resp.Diagnostics.AddError("Error Creating Stream", fmt.Sprintf("Unable to create stream %q, got error: %s", s.Name, err))
		return
	}

	// Computed fields
	data.RetentionSize = types.Int64Value(int64(s.RetentionBytes))
	data.RetentionPeriod = types.Int64Value(int64(s.RetentionMS))

	for i := range data.Columns {
		name := data.Columns[i].Name.ValueString()
		for j := range s.Columns {
			if s.Columns[j].Name == name {
				data.Columns[i].CodecExpression = types.StringValue(s.Columns[j].Codec)
				break
			}
		}
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a timeplus_stream resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

var spaces = regexp.MustCompile(`\s+`)

func (r *streamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *streamResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	s, err := r.client.GetStream(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Stream",
			fmt.Sprintf("Unable to read stream %q, got error: %s",
				data.Name.ValueString(), err))
		return
	}

	hasTpTimeColumn := false
	for i := range data.Columns {
		if data.Columns[i].Name.ValueString() == "_tp_time" {
			hasTpTimeColumn = true
			break
		}
	}

	// required fields
	data.Name = types.StringValue(s.Name)

	columnStates := data.Columns
	data.Columns = make([]columnModel, 0, len(s.Columns))
	for i := range s.Columns {
		name := s.Columns[i].Name
		// if `_tp_time` is not explicitely specified in the resource definition, don't show it
		if name == "_tp_time" && !hasTpTimeColumn {
			continue
		}

		codec := types.StringValue("")
		for j := range columnStates {
			if columnStates[j].Name.ValueString() == name {
				codec = columnStates[j].Codec
			}
		}
		data.Columns = append(data.Columns, columnModel{
			Name:            types.StringValue(name),
			Type:            types.StringValue(s.Columns[i].Type),
			Default:         types.StringValue(s.Columns[i].Default),
			Codec:           codec,
			CodecExpression: types.StringValue(s.Columns[i].Codec),
		})
	}

	// optional fields
	if !(data.Description.IsNull() && s.Description == "") {
		data.Description = types.StringValue(s.Description)
	}

	// the create stream API will set retention_bytes to 0 if it's not provided
	if !(data.RetentionSize.IsNull() && s.RetentionBytes == 0) {
		data.RetentionSize = types.Int64Value(int64(s.RetentionBytes))
	}

	// the create stream API will set retention_ms to 0 if it's not provided
	if !(data.RetentionPeriod.IsNull() && s.RetentionMS == 0) {
		data.RetentionPeriod = types.Int64Value(int64(s.RetentionMS))
	}

	if !(data.HistoricalDataTTL.IsNull() && s.HistoricalTTLExpression == "") {
		// assume TTL expressions are space insignificant
		if spaces.ReplaceAllString(data.HistoricalDataTTL.ValueString(), "") != spaces.ReplaceAllString(s.HistoricalTTLExpression, "") {
			data.HistoricalDataTTL = types.StringValue(s.HistoricalTTLExpression)
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *streamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *streamResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	columns := make([]timeplus.Column, 0, len(data.Columns))
	for i := range data.Columns {
		columns = append(columns, timeplus.Column{
			Name:    data.Columns[i].Name.ValueString(),
			Type:    data.Columns[i].Type.ValueString(),
			Default: data.Columns[i].Default.ValueString(),
		})
	}

	s := timeplus.Stream{
		Name:                    data.Name.ValueString(),
		Description:             data.Description.ValueString(),
		Columns:                 columns,
		RetentionBytes:          int(data.RetentionSize.ValueInt64()),
		RetentionMS:             int(data.RetentionPeriod.ValueInt64()),
		HistoricalTTLExpression: data.HistoricalDataTTL.ValueString(),
	}
	if err := r.client.UpdateStream(&s); err != nil {
		resp.Diagnostics.AddError("Error Updating Stream", fmt.Sprintf("Unable to update stream %q, got error: %s", s.Name, err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *streamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *streamResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteStream(&timeplus.Stream{Name: data.Name.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Stream", fmt.Sprintf("Unable to delete stream %q, got error: %s", data.Name.ValueString(), err))
	}
}

func (r *streamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
