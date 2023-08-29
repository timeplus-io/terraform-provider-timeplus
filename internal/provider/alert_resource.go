// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/timeplus-io/terraform-provider-timeplus/internal/timeplus"
	myvalidator "github.com/timeplus-io/terraform-provider-timeplus/internal/validator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &alertResource{}
var _ resource.ResourceWithImportState = &alertResource{}

func NewAlertResource() resource.Resource {
	return &alertResource{}
}

// alertResource defines the resource implementation.
type alertResource struct {
	client *timeplus.Client
}

// alertResourceModel describes the alert resource data model.
type alertResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Severity    types.Int64  `tfsdk:"severity"`
	Action      types.String `tfsdk:"action"`
	Properties  types.String `tfsdk:"properties"`
	TriggerSQL  types.String `tfsdk:"trigger_sql"`
	ResolveSQL  types.String `tfsdk:"resolve_sql"`
}

func (r *alertResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert"
}

func (r *alertResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: `Alerts are like sinks, they are used to send data to external systems. How alerts are different is that alerts have two statuses: 'triggerred' and 'resolved'.

An alert runs two queries in background to detect if the status should be triggerred or resolved. Once an alert is in one status, it won't send the same kind of events to the target until it gets a different kind of event.`,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The alert immutable ID, generated by Timeplus",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The human-friendly name for the alert",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A detailed text describes the alert",
				Optional:            true,
			},
			"severity": schema.Int64Attribute{
				MarkdownDescription: "A number indicates how serious this alert is",
				Optional:            true,
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "The type of action the alert should take, i.e. the name of the target system, like 'slack', 'email', etc. Please refer to the Timeplus document for supported alert action types",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// since Terraform does not have built-in support for map[string]any with the framework library, we use JSON as a simple solution
			"properties": schema.StringAttribute{
				MarkdownDescription: "a JSON object defines the configurations for the specific alert action. The properites could contain sensitive information like password, secret, etc.",
				Required:            true,
				Sensitive:           true,
				Validators: []validator.String{
					myvalidator.JsonObject(),
				},
			},
			"trigger_sql": schema.StringAttribute{
				MarkdownDescription: "The query the alert uses to generate events that trigger the alert",
				Required:            true,
			},
			"resolve_sql": schema.StringAttribute{
				MarkdownDescription: "The query the alert uses to generate events that resolve the alert",
				Optional:            true,
			},
		},
	}
}

func (r *alertResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *alertResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *alertResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	props := make(map[string]any)
	// data.Properties uses the JsonObject validator, so it's guaranteed that it's a valid JSON object
	if data.Properties.ValueString() != "" {
		_ = json.Unmarshal([]byte(data.Properties.ValueString()), &props)
	}

	s := timeplus.Alert{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Severity:    int(data.Severity.ValueInt64()),
		Action:      data.Action.ValueString(),
		Properties:  props,
		TriggerSQL:  data.TriggerSQL.ValueString(),
		ResolveSQL:  data.ResolveSQL.ValueString(),
	}

	if err := r.client.CreateAlert(&s); err != nil {
		resp.Diagnostics.AddError("Error Creating Alert", fmt.Sprintf("Unable to create alert %q, got error: %s", s.Name, err))
		return
	}

	// set Computed fields
	data.ID = types.StringValue(s.ID)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a timeplus_alert resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *alertResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *alertResourceModel

	// Read Terraform prior state data into the model
	if resp.Diagnostics.Append(req.State.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	s, err := r.client.GetAlert(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Alert",
			fmt.Sprintf("Unable to read alert %q, got error: %s",
				data.Name.ValueString(), err))
		return
	}

	// required fields
	data.Name = types.StringValue(s.Name)
	data.Severity = types.Int64Value(int64(s.Severity))
	data.Action = types.StringValue(s.Action)
	data.TriggerSQL = types.StringValue(s.TriggerSQL)

	// We can't compare the JSON directly since order is not guaranteed, need a bit more work to detect if properties are acutally changed
	props := make(map[string]any)
	if data.Properties.ValueString() != "" {
		_ = json.Unmarshal([]byte(data.Properties.ValueString()), &props)
	}

	clone := maps.Clone(props)

	// API does not return sensitive fields, thus we can't simply use s.Properties to replace data.Properties
	maps.Copy(props, s.Properties)

	if !reflect.DeepEqual(clone, props) {
		propsBytes, err := json.Marshal(s.Properties)
		if err != nil {
			resp.Diagnostics.AddError("Bad Alert Properties", fmt.Sprintf("Unable to encode alert properties into JSON, got error: %s", err))
			return
		}
		data.Properties = types.StringValue(string(propsBytes))
	}

	// optional fields
	if !(data.Description.IsNull() && s.Description == "") {
		data.Description = types.StringValue(s.Description)
	}

	if !(data.ResolveSQL.IsNull() && s.ResolveSQL == "") {
		data.ResolveSQL = types.StringValue(s.ResolveSQL)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *alertResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *alertResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	props := make(map[string]any)
	// data.Properties uses the JsonObject validator, so it's guaranteed that it's a valid JSON object
	if data.Properties.ValueString() != "" {
		_ = json.Unmarshal([]byte(data.Properties.ValueString()), &props)
	}

	s := timeplus.Alert{
		ID:          data.ID.ValueString(),
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Severity:    int(data.Severity.ValueInt64()),
		Action:      data.Action.ValueString(),
		Properties:  props,
		TriggerSQL:  data.TriggerSQL.ValueString(),
		ResolveSQL:  data.ResolveSQL.ValueString(),
	}

	if err := r.client.UpdateAlert(&s); err != nil {
		resp.Diagnostics.AddError("Error Updating Alert", fmt.Sprintf("Unable to update alert %q, got error: %s", s.Name, err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *alertResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *alertResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAlert(&timeplus.Alert{ID: data.ID.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Alert", fmt.Sprintf("Unable to delete alert %q, got error: %s", data.Name.ValueString(), err))
	}
}

func (r *alertResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
