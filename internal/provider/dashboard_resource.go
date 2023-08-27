// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
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
var _ resource.Resource = &dashboardResource{}
var _ resource.ResourceWithImportState = &dashboardResource{}

func NewDashboardResource() resource.Resource {
	return &dashboardResource{}
}

// dashboardResource defines the resource implementation.
type dashboardResource struct {
	client *timeplus.Client
}

// dashboardResourceModel describes the dashboard resource data model.
type dashboardResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Panels      types.String `tfsdk:"panels"`
}

func (r *dashboardResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dashboard"
}

func (r *dashboardResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "A dashboard is a set of one or more panels organized and arranged in one web page. A variety of panels are supported to make it easy to construct the visualization components so that you can create the dashboards for specific monitoring and analytics needs.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The dashboard immutable ID, generated by Timeplus",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The human-friendly name for the dashboard",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A detailed text describes the dashboard",
				Optional:            true,
			},
			"panels": schema.StringAttribute{
				MarkdownDescription: "A list of panels defined in a JSON array. The best way to generate such array is to copy it directly from the Timeplus console UI.",
				Required:            true,
			},
		},
	}
}

func (r *dashboardResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*timeplus.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected resource Configure Type",
			fmt.Sprintf("Expected *timeplus.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *dashboardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *dashboardResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var panels []timeplus.Panel
	if data.Panels.ValueString() != "" {
		err := json.Unmarshal([]byte(data.Panels.ValueString()), &panels)
		if err != nil {
			resp.Diagnostics.AddAttributeError(path.Root("panels"), "Invalid panels JSON", err.Error())
		}
	}

	s := timeplus.Dashboard{
		ID:          data.ID.ValueString(),
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Panels:      panels,
	}

	if err := r.client.CreateDashboard(&s); err != nil {
		resp.Diagnostics.AddError("Error Creating Dashboard", fmt.Sprintf("Unable to create dashboard %q, got error: %s", s.Name, err))
		return
	}

	// set Computed fields
	data.ID = types.StringValue(s.ID)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a timeplus_dashboard resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dashboardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *dashboardResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	s, err := r.client.GetDashboard(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Dashboard",
			fmt.Sprintf("Unable to read dashboard %q, got error: %s",
				data.Name.ValueString(), err))
		return
	}

	// required fields
	data.Name = types.StringValue(s.Name)

	bytes, err := json.Marshal(s.Panels)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("panels"), "Panel Parsing Failed", err.Error())
	}

	data.Panels = types.StringValue(string(bytes))

	// optional fields
	if !(data.Description.IsNull() && s.Description == "") {
		data.Description = types.StringValue(s.Description)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dashboardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *dashboardResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var panels []timeplus.Panel
	if data.Panels.ValueString() != "" {
		err := json.Unmarshal([]byte(data.Panels.ValueString()), &panels)
		if err != nil {
			resp.Diagnostics.AddAttributeError(path.Root("panels"), "Invalid panels JSON", err.Error())
		}
	}

	s := timeplus.Dashboard{
		ID:          data.ID.ValueString(),
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Panels:      panels,
	}

	if err := r.client.UpdateDashboard(&s); err != nil {
		resp.Diagnostics.AddError("Error Updating Dashboard", fmt.Sprintf("Unable to update dashboard %q, got error: %s", s.Name, err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dashboardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *dashboardResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDashboard(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Dashboard", fmt.Sprintf("Unable to delete dashboard %q, got error: %s", data.Name.ValueString(), err))
	}
}

func (r *dashboardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
