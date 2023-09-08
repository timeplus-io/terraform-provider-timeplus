// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

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
	Name           types.String `tfsdk:"name"`
	Type           types.String `tfsdk:"type"`
	Default        types.String `tfsdk:"default"`
	Codec          types.String `tfsdk:"codec"`
	UseAsEventTime types.Bool   `tfsdk:"use_as_event_time"`
	PrimaryKey     types.Bool   `tfsdk:"primary_key"`
}

// streamResourceModel describes the stream resource data model.
type streamResourceModel struct {
	Name           types.String  `tfsdk:"name"`
	Description    types.String  `tfsdk:"description"`
	Columns        []columnModel `tfsdk:"column"`
	RetentionBytes types.Int64   `tfsdk:"retention_bytes"`
	RetentionMS    types.Int64   `tfsdk:"retention_ms"`
	HistoryTTL     types.String  `tfsdk:"history_ttl"`
	Mode           types.String  `tfsdk:"mode"`
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
			"mode": schema.StringAttribute{
				MarkdownDescription: "The stream mode. Options: append, changelog, changelog_kv, versioned_kv. Default: \"append\"",
				Optional:            true,
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
			},
		},
		Blocks: map[string]schema.Block{
			"column": schema.ListNestedBlock{
				MarkdownDescription: "Define the columns of the stream",
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
						"use_as_event_time": schema.BoolAttribute{
							MarkdownDescription: "If set to `true`, this column will be used as the event time column (by default ingest time will be used as event time). Only one column can be marked as the event time column in a stream.",
							Optional:            true,
						},
						"primary_key": schema.BoolAttribute{
							MarkdownDescription: "If set to `true`, this column will be used as the primary key, or part of the combined primary key if multiple columns are marked as primary keys.",
							Optional:            true,
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

	mode, err := timeplus.StreamModeFrom(data.Mode.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("mode"), "Invalid Mode", err.Error())
		return
	}

	if len(data.Columns) == 0 {
		resp.Diagnostics.AddAttributeError(path.Root("column"), "No Columns", "At least one column must be defined for a stream.")
		return
	}

	eventTimeColumn := ""
	for i := range data.Columns {
		if data.Columns[i].UseAsEventTime.ValueBool() {
			if eventTimeColumn != "" {
				resp.Diagnostics.AddAttributeError(path.Root(fmt.Sprintf("column[%d]", i)), "Too Many EventTime Columns", "Only one column can be marked as event time column.")
				return
			}
			eventTimeColumn = data.Columns[i].Name.ValueString()
		}
	}

	primaryKeys := []string{}
	columns := make([]timeplus.Column, 0, len(data.Columns))
	for i := range data.Columns {
		if data.Columns[i].PrimaryKey.ValueBool() {
			primaryKeys = append(primaryKeys, "`"+data.Columns[i].Name.ValueString()+"`")
		}
		columns = append(columns, timeplus.Column{
			Name:    data.Columns[i].Name.ValueString(),
			Type:    data.Columns[i].Type.ValueString(),
			Default: data.Columns[i].Default.ValueString(),
			Codec:   data.Columns[i].Codec.ValueString(),
		})
	}

	s := timeplus.Stream{
		Name:                    data.Name.ValueString(),
		Description:             data.Description.ValueString(),
		Columns:                 columns,
		RetentionBytes:          int(data.RetentionBytes.ValueInt64()),
		RetentionMS:             int(data.RetentionMS.ValueInt64()),
		HistoricalTTLExpression: data.HistoryTTL.ValueString(),
		Mode:                    string(mode),
	}

	if len(primaryKeys) > 0 {
		s.PrimaryKey = fmt.Sprintf("(%s)", strings.Join(primaryKeys, ","))
	}

	if eventTimeColumn != "" {
		s.EventTimeColumn = eventTimeColumn
	}

	if err := r.client.CreateStream(&s); err != nil {
		resp.Diagnostics.AddError("Error Creating Stream", fmt.Sprintf("Unable to create stream %q, got error: %s", s.Name, err))
		return
	}

	// Computed fields
	data.RetentionBytes = types.Int64Value(int64(s.RetentionBytes))
	data.RetentionMS = types.Int64Value(int64(s.RetentionMS))

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

	// required fields
	data.Name = types.StringValue(s.Name)

	// we need to handle the case that the `_tp_time` column is explicitely defined by users
	hasTpTimeColumn := false
	// the column lookup map
	dataColumns := make(map[string]columnModel, len(data.Columns))
	for _, col := range data.Columns {
		name := col.Name.ValueString()

		if name == "_tp_time" {
			hasTpTimeColumn = true
		}

		dataColumns[name] = col
	}

	pKeys := map[string]struct{}{}
	for _, k := range strings.Split(s.PrimaryKey, ",") {
		// remove the the quotes "`" if they exists
		k = strings.TrimSuffix(
			strings.TrimPrefix(
				strings.TrimSpace(k), "`"), "`")
		pKeys[k] = struct{}{}
	}

	data.Columns = make([]columnModel, 0, len(s.Columns))
	for i := range s.Columns {
		name := s.Columns[i].Name
		// if `_tp_time` is not explicitely specified in the resource definition, don't show it
		if name == "_tp_time" && !hasTpTimeColumn {
			continue
		}

		// `codec` returned by the API contains the `CODEC()` function call, like `CODEC(LZ4)`.
		// Removing the surrounding `CODEC()` to match the input.
		codec := types.StringValue(strings.TrimSuffix(strings.TrimPrefix(s.Columns[i].Codec, "CODEC("), ")"))

		col := columnModel{
			Name:    types.StringValue(name),
			Type:    types.StringValue(s.Columns[i].Type),
			Default: types.StringValue(s.Columns[i].Default),
			Codec:   codec,
		}

		if _, ok := pKeys[name]; ok {
			col.PrimaryKey = types.BoolValue(true)
		}

		if dataColumn, ok := dataColumns[name]; ok {
			// FIXME parse the default value of `_tp_time` to figure out which column is used as event time column
			col.UseAsEventTime = dataColumn.UseAsEventTime
		}

		data.Columns = append(data.Columns, col)
	}

	// optional fields
	if !(data.Description.IsNull() && s.Description == "") {
		data.Description = types.StringValue(s.Description)
	}

	// the create stream API will set retention_bytes to 0 if it's not provided
	if !(data.RetentionBytes.IsNull() && s.RetentionBytes == 0) {
		data.RetentionBytes = types.Int64Value(int64(s.RetentionBytes))
	}

	// the create stream API will set retention_ms to 0 if it's not provided
	if !(data.RetentionMS.IsNull() && s.RetentionMS == 0) {
		data.RetentionMS = types.Int64Value(int64(s.RetentionMS))
	}

	if !(data.HistoryTTL.IsNull() && s.HistoricalTTLExpression == "") {
		// assume TTL expressions are space insignificant
		if spaces.ReplaceAllString(data.HistoryTTL.ValueString(), "") != spaces.ReplaceAllString(s.HistoricalTTLExpression, "") {
			data.HistoryTTL = types.StringValue(s.HistoricalTTLExpression)
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

	mode, err := timeplus.StreamModeFrom(data.Mode.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("mode"), "Invalid Mode", err.Error())
		return
	}

	if len(data.Columns) == 0 {
		resp.Diagnostics.AddAttributeError(path.Root("column"), "No Columns", "At least one column must be defined for a stream.")
		return
	}

	eventTimeColumn := ""
	for i := range data.Columns {
		if data.Columns[i].UseAsEventTime.ValueBool() {
			if eventTimeColumn != "" {
				resp.Diagnostics.AddAttributeError(path.Root(fmt.Sprintf("column[%d]", i)), "Too Many EventTime Columns", "Only one column can be marked as event time column.")
				return
			}
			eventTimeColumn = data.Columns[i].Name.ValueString()
		}
	}

	primaryKeys := []string{}
	columns := make([]timeplus.Column, 0, len(data.Columns))
	for i := range data.Columns {
		if data.Columns[i].PrimaryKey.ValueBool() {
			primaryKeys = append(primaryKeys, data.Columns[i].Name.ValueString())
		}
		columns = append(columns, timeplus.Column{
			Name:    data.Columns[i].Name.ValueString(),
			Type:    data.Columns[i].Type.ValueString(),
			Default: data.Columns[i].Default.ValueString(),
			Codec:   data.Columns[i].Codec.ValueString(),
		})
	}

	s := timeplus.Stream{
		Name:                    data.Name.ValueString(),
		Description:             data.Description.ValueString(),
		Columns:                 columns,
		RetentionBytes:          int(data.RetentionBytes.ValueInt64()),
		RetentionMS:             int(data.RetentionMS.ValueInt64()),
		HistoricalTTLExpression: data.HistoryTTL.ValueString(),
		Mode:                    string(mode),
	}

	if len(primaryKeys) > 0 {
		s.PrimaryKey = fmt.Sprintf("(%s)", strings.Join(primaryKeys, ","))
	}

	if eventTimeColumn != "" {
		s.EventTimeColumn = eventTimeColumn
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
