// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/timeplus-io/terraform-provider-timeplus/internal/timeplus"
	myValidator "github.com/timeplus-io/terraform-provider-timeplus/internal/validator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &javascriptFunctionResource{}
var _ resource.ResourceWithImportState = &javascriptFunctionResource{}

func NewJavascriptFunctionResource() resource.Resource {
	return &javascriptFunctionResource{}
}

// javascriptFunctionResource defines the resource implementation.
type javascriptFunctionResource struct {
	client *timeplus.Client
}

// javascriptFunctionResourceModel describes the stream resource data model.
type javascriptFunctionResourceModel struct {
	Name           types.String            `tfsdk:"name"`
	Description    types.String            `tfsdk:"description"`
	Arguments      []functionArgumentModel `tfsdk:"arg"`
	ReturnType     types.String            `tfsdk:"return_type"`
	Source         types.String            `tfsdk:"source"`
	IsAggrFunction types.Bool              `tfsdk:"is_aggregate_function"`
}

func (r *javascriptFunctionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_javascript_function"
}

func (r *javascriptFunctionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Timeplus javascript functions are one of the supported user defined function types. They allow users to register a HTTP webhook as a function which can be called in any queries.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The javascript function name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A detailed text describes the javascript function",
				Optional:            true,
			},
			"is_aggregate_function": schema.BoolAttribute{
				MarkdownDescription: "Indecates if the javascript function an aggregate function",
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"source": schema.StringAttribute{
				MarkdownDescription: "The javascript function source code",
				Required:            true,
				Validators: []validator.String{
					myValidator.URL(),
				},
			},
			"return_type": schema.StringAttribute{
				MarkdownDescription: "The type of the function's return value",
				Required:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"arg": schema.ListNestedBlock{
				MarkdownDescription: "Describe an argument of the javascript function, argument order matters",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The argument name",
							Required:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The argument type",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (r *javascriptFunctionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *javascriptFunctionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *javascriptFunctionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	args := make([]timeplus.UDFArgument, 0, len(data.Arguments))
	for i := range data.Arguments {
		args = append(args, timeplus.UDFArgument{
			Name: data.Arguments[i].Name.ValueString(),
			Type: data.Arguments[i].Type.ValueString(),
		})
	}

	f := timeplus.UDF{
		Name:           data.Name.ValueString(),
		Description:    data.Description.ValueString(),
		Type:           timeplus.UDFTypeJavascript,
		Arguments:      args,
		ReturnType:     data.ReturnType.ValueString(),
		IsAggrFunction: data.IsAggrFunction.ValueBool(),
		Source:         data.Source.ValueString(),
	}
	if err := r.client.CreateUDF(&f); err != nil {
		resp.Diagnostics.AddError("Error Creating JavascriptFunction", fmt.Sprintf("Unable to create javascript function %q, got error: %s", f.Name, err))
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a timeplus_remote_function resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *javascriptFunctionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *javascriptFunctionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	s, err := r.client.GetUDF(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading JavascriptFunction",
			fmt.Sprintf("Unable to read javascript function %q, got error: %s",
				data.Name.ValueString(), err))
		return
	}

	if s.Type != timeplus.UDFTypeJavascript {
		resp.Diagnostics.AddError(
			"Error Reading JavascriptFunction",
			fmt.Sprintf("Function with name %s is not a javascript function",
				data.Name.ValueString()))
		return
	}

	// required fields
	data.Name = types.StringValue(s.Name)
	data.Source = types.StringValue(s.Source)
	data.ReturnType = types.StringValue(s.ReturnType)

	// optional fields
	data.IsAggrFunction = types.BoolValue(s.IsAggrFunction)

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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *javascriptFunctionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *javascriptFunctionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	args := make([]timeplus.UDFArgument, 0, len(data.Arguments))
	for i := range data.Arguments {
		args = append(args, timeplus.UDFArgument{
			Name: data.Arguments[i].Name.ValueString(),
			Type: data.Arguments[i].Type.ValueString(),
		})
	}

	f := timeplus.UDF{
		Name:           data.Name.ValueString(),
		Description:    data.Description.ValueString(),
		Type:           timeplus.UDFTypeJavascript,
		Arguments:      args,
		ReturnType:     data.ReturnType.ValueString(),
		IsAggrFunction: data.IsAggrFunction.ValueBool(),
		Source:         data.Source.ValueString(),
	}
	if err := r.client.UpdateUDF(&f); err != nil {
		resp.Diagnostics.AddError("Error Updating JavascriptFunction", fmt.Sprintf("Unable to update javascript function %q, got error: %s", f.Name, err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *javascriptFunctionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *javascriptFunctionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUDF(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting JavascriptFunction", fmt.Sprintf("Unable to delete javascript function %q, got error: %s", data.Name.ValueString(), err))
	}
}

func (r *javascriptFunctionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
