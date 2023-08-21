// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/timeplus-io/terraform-provider-timeplus/internal/timeplus"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &viewResource{}
var _ resource.ResourceWithImportState = &viewResource{}

func NewViewResource() resource.Resource {
	return &viewResource{}
}

// viewResource defines the resource implementation.
type viewResource struct {
	client *timeplus.Client
}

// viewResourceModel describes the stream resource data model.
type viewResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Query       types.String `tfsdk:"query"`
}

func (r *viewResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_view"
}

func (r *viewResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Timeplus views are named queries. When you create a view, you basically create a query and assign a name to the query. Therefore, a view is useful for wrapping a commonly used complex query.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The view name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A detailed text describes the view",
				Optional:            true,
			},
			"query": schema.StringAttribute{
				MarkdownDescription: "The query SQL of the view",
				Required:            true,
			},
		},
	}
}

func (r *viewResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *viewResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *viewResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	v := timeplus.View{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Query:       data.Query.ValueString(),
	}
	if err := r.client.CreateView(&v); err != nil {
		resp.Diagnostics.AddError("Error Creating View", fmt.Sprintf("Unable to create view %q, got error: %s", v.Name, err))
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a timeplus_view resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *viewResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *viewResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	v, err := r.client.GetView(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading View",
			fmt.Sprintf("Unable to read view %q, got error: %s",
				data.Name.ValueString(), err))
		return
	}

	// required fields
	data.Name = types.StringValue(v.Name)
	data.Query = types.StringValue(v.Query)

	// optional fields
	if !(data.Description.IsNull() && v.Description == "") {
		data.Description = types.StringValue(v.Description)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *viewResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *viewResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	v := timeplus.View{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Query:       data.Query.ValueString(),
	}
	if err := r.client.UpdateView(&v); err != nil {
		resp.Diagnostics.AddError("Error Updating View", fmt.Sprintf("Unable to update view %q, got error: %s", v.Name, err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *viewResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *viewResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteView(&timeplus.View{Name: data.Name.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting View", fmt.Sprintf("Unable to delete view %q, got error: %s", data.Name.ValueString(), err))
	}
}

func (r *viewResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
