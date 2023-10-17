package provider

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api_v0 "github.com/hashicorp/terraform-provider-scaffolding-framework/api/v0"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ProgramResource{}
var _ resource.ResourceWithImportState = &ProgramResource{}

func NewProgramResource() resource.Resource {
	return &ProgramResource{}
}

// ProgramResource defines the resource implementation.
type ProgramResource struct {
	providerData *FelderaProviderData
}

// ProgramResourceModel describes the resource data model.
type ProgramResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name" json:"name"`
	Description   types.String `tfsdk:"description" json:"description"`
	Code          types.String `tfsdk:"code" json:"code"`
	ShouldCompile types.Bool   `tfsdk:"should_compile"`
	Version       types.Int64  `tfsdk:"version"`
}

func (r *ProgramResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_program"
}

func (r *ProgramResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Program resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique program identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:            false,
				MarkdownDescription: "Unique program name.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Computed:            false,
				MarkdownDescription: "Program description.",
				Required:            true,
			},
			"code": schema.StringAttribute{
				Computed:            false,
				MarkdownDescription: "SQL code of the program.",
				Required:            true,
			},
			"should_compile": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Should the program be marked for compilation.",
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"version": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Version of the program.",
			},
		},
	}
}

func (r *ProgramResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*FelderaProviderData)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *FelderaProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.providerData = providerData
}

func (r *ProgramResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProgramResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpRes, err := r.providerData.Client.NewProgramWithResponse(ctx, api_v0.NewProgramJSONRequestBody{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Code:        data.Code.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create program, got error: %s", err))
		return
	}

	if httpRes.JSON201 == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create program, got status: %s resp: %s", httpRes.Status(), httpRes.Body))
		return
	}

	program := httpRes.JSON201

	// save into the Terraform state.
	data.Id = types.StringValue(program.ProgramId.String())
	data.Version = types.Int64Value(program.Version)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// TODO compile

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProgramResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProgramResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	uid, err := uuid.Parse(data.Id.String())
	if err != nil {
		resp.Diagnostics.AddError("Parser Error", fmt.Sprintf("Unable to parse program id, got error: %s", err))
		return
	}

	clientResp, err := r.providerData.Client.GetProgramWithResponse(ctx, uid, &api_v0.GetProgramParams{})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read program, got error: %s", err))
		return
	} else if clientResp.JSON200 == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read program, got status: %s resp: %s", clientResp.Status(), clientResp.Body))
		return
	}

	data.Name = types.StringValue(clientResp.JSON200.Name)
	data.Version = types.Int64Value(clientResp.JSON200.Version)
	data.Description = types.StringValue(clientResp.JSON200.Description)
	if clientResp.JSON200.Code != nil {
		data.Code = types.StringValue(*clientResp.JSON200.Code)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProgramResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ProgramResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	uid, err := uuid.Parse(data.Id.String())
	if err != nil {
		resp.Diagnostics.AddError("Parser Error", fmt.Sprintf("Unable to parse program id, got error: %s", err))
		return
	}

	httpRes, err := r.providerData.Client.UpdateProgramWithResponse(ctx, uid, api_v0.UpdateProgramJSONRequestBody{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueStringPointer(),
		Code:        data.Code.ValueStringPointer(),
	})

	if err != nil || httpRes.JSON200 == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update program, got error: %s", err))
		return
	}

	data.Version = types.Int64Value(httpRes.JSON200.Version)

	// TODO compile

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProgramResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProgramResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	uid, err := uuid.FromBytes([]byte(data.Id.String()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to parse program id, got error: %s", err))
		return
	}

	_, err = r.providerData.Client.DeleteProgramWithResponse(ctx, uid, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete program, got error: %s", err))
		return
	}

}

func (r *ProgramResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
