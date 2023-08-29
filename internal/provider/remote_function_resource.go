// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/timeplus-io/terraform-provider-timeplus/internal/timeplus"
	myValidator "github.com/timeplus-io/terraform-provider-timeplus/internal/validator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &remoteFunctionResource{}
var _ resource.ResourceWithImportState = &remoteFunctionResource{}

func NewRemoteFunctionResource() resource.Resource {
	return &remoteFunctionResource{}
}

// remoteFunctionResource defines the resource implementation.
type remoteFunctionResource struct {
	client *timeplus.Client
}

type functionArgumentModel struct {
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

type functionAuthHeaderModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

// remoteFunctionResourceModel describes the stream resource data model.
type remoteFunctionResourceModel struct {
	Name        types.String             `tfsdk:"name"`
	Description types.String             `tfsdk:"description"`
	Arguments   []functionArgumentModel  `tfsdk:"arg"`
	ReturnType  types.String             `tfsdk:"return_type"`
	URL         types.String             `tfsdk:"url"`
	AuthHeader  *functionAuthHeaderModel `tfsdk:"auth_header"`
}

func (r *remoteFunctionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_function"
}

func (r *remoteFunctionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Timeplus remote functions are one of the supported user defined function types. Remote functions allow users to register a HTTP webhook as a function which can be called in queries.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The remote function name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A detailed text describes the remote function",
				Optional:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The HTTP endpoint to be used to call the function",
				Required:            true,
				Validators: []validator.String{
					myValidator.URL(),
				},
			},
			"auth_header": schema.SingleNestedAttribute{
				MarkdownDescription: "The HTTP header and its value to be used as an authentication means to call the function. The remote function can use this information to determine if it's a valid call",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "The HTTP header name",
						Required:            true,
					},
					"value": schema.StringAttribute{
						MarkdownDescription: "The value for the header",
						Required:            true,
					},
				},
			},
			"return_type": schema.StringAttribute{
				MarkdownDescription: "The type of the function's return value",
				Required:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"arg": schema.ListNestedBlock{
				MarkdownDescription: "Describe an argument of the remote function, argument order matters",
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

func (r *remoteFunctionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *remoteFunctionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *remoteFunctionResourceModel

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

	var (
		authCtx    *timeplus.UDFAuthContext
		authMethod timeplus.UDFAuthMethod
	)
	if data.AuthHeader != nil {
		authMethod = timeplus.UDFAuthHeader
		authCtx = &timeplus.UDFAuthContext{
			Name:  data.AuthHeader.Name.ValueString(),
			Value: data.AuthHeader.Value.ValueString(),
		}
	}
	f := timeplus.UDF{
		Name:           data.Name.ValueString(),
		Description:    data.Description.ValueString(),
		Type:           timeplus.UDFTypeRemote,
		Arguments:      args,
		ReturnType:     data.ReturnType.ValueString(),
		URL:            data.URL.ValueString(),
		AuthMethod:     authMethod,
		AuthContext:    authCtx,
		IsAggrFunction: false,
		Source:         "",
	}
	if err := r.client.CreateUDF(&f); err != nil {
		resp.Diagnostics.AddError("Error Creating RemoteFunction", fmt.Sprintf("Unable to create remote function %q, got error: %s", f.Name, err))
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a timeplus_remote_function resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *remoteFunctionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *remoteFunctionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	s, err := r.client.GetUDF(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading RemoteFunction",
			fmt.Sprintf("Unable to read remote function %q, got error: %s",
				data.Name.ValueString(), err))
		return
	}

	if s.Type != timeplus.UDFTypeRemote {
		resp.Diagnostics.AddError(
			"Error Reading RemoteFunction",
			fmt.Sprintf("Function with name %s is not a remote function",
				data.Name.ValueString()))
		return
	}

	// required fields
	data.Name = types.StringValue(s.Name)
	data.URL = types.StringValue(s.URL)
	data.ReturnType = types.StringValue(s.ReturnType)

	// optional fields
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

	if !(data.AuthHeader == nil && (s.AuthMethod == timeplus.UDFAuthNone || s.AuthMethod == "")) {
		ctx := s.AuthContext
		if ctx == nil {
			ctx = &timeplus.UDFAuthContext{}
		}
		data.AuthHeader = &functionAuthHeaderModel{
			Name:  types.StringValue(ctx.Name),
			Value: types.StringValue(ctx.Value),
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *remoteFunctionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *remoteFunctionResourceModel

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

	var (
		authCtx    *timeplus.UDFAuthContext
		authMethod timeplus.UDFAuthMethod
	)
	if data.AuthHeader != nil {
		authMethod = timeplus.UDFAuthHeader
		authCtx = &timeplus.UDFAuthContext{
			Name:  data.AuthHeader.Name.ValueString(),
			Value: data.AuthHeader.Value.ValueString(),
		}
	}
	f := timeplus.UDF{
		Name:           data.Name.ValueString(),
		Description:    data.Description.ValueString(),
		Type:           timeplus.UDFTypeRemote,
		Arguments:      args,
		ReturnType:     data.ReturnType.ValueString(),
		URL:            data.URL.ValueString(),
		AuthMethod:     authMethod,
		AuthContext:    authCtx,
		IsAggrFunction: false,
		Source:         "",
	}
	if err := r.client.UpdateUDF(&f); err != nil {
		resp.Diagnostics.AddError("Error Updating RemoteFunction", fmt.Sprintf("Unable to update remote function %q, got error: %s", f.Name, err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *remoteFunctionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *remoteFunctionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUDF(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting RemoteFunction", fmt.Sprintf("Unable to delete remote function %q, got error: %s", data.Name.ValueString(), err))
	}
}

func (r *remoteFunctionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
