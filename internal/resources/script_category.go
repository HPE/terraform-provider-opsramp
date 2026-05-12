// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/HPE/terraform-provider-opsramp/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &ScriptCategoryResource{}
var _ resource.ResourceWithImportState = &ScriptCategoryResource{}

// ScriptCategoryResource defines the resource implementation.
type ScriptCategoryResource struct {
	apiClient *client.OpsRampClient
}

// ScriptCategoryModel maps Terraform schema attributes to the provider model.
type ScriptCategoryModel struct {
	Client   types.String `tfsdk:"client"`
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	ParentId types.Int64  `tfsdk:"parent_id"`
}

// NewScriptCategory creates a new instance of the resource.
func NewScriptCategory() resource.Resource {
	return &ScriptCategoryResource{}
}

// Metadata returns the resource type name.
func (r *ScriptCategoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_script_category"
}

// Schema defines the schema for the resource.
func (r *ScriptCategoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp RBA Script Category. Categories organize scripts and can be nested using `parent_id`.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client (tenant) UUID where this category should be created. If not specified, uses the provider's tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The numeric ID of the category (as a string), assigned by the API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the script category.",
			},
			"parent_id": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The numeric ID of the parent category. Omit or set to `0` for a root-level category.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure prepares the resource with the API client.
func (r *ScriptCategoryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.OpsRampClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *client.OpsRampClient",
		)
		return
	}

	r.apiClient = c
}

func (r *ScriptCategoryResource) resolveTenantId(clientAttr types.String) string {
	if !clientAttr.IsNull() && clientAttr.ValueString() != "" {
		return clientAttr.ValueString()
	}
	return r.apiClient.TenantId
}

// buildCategoryRequest converts the Terraform model into an API request body.
func buildCategoryRequest(plan ScriptCategoryModel) client.ScriptCategory {
	cat := client.ScriptCategory{
		Name: plan.Name.ValueString(),
	}
	if !plan.ParentId.IsNull() && !plan.ParentId.IsUnknown() && plan.ParentId.ValueInt64() != 0 {
		cat.Parent = &client.ScriptCategoryParentRef{Id: int(plan.ParentId.ValueInt64())}
	}
	return cat
}

// Create handles the creation of the resource.
func (r *ScriptCategoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ScriptCategoryModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(plan.Client)
	catReq := buildCategoryRequest(plan)

	created, err := r.apiClient.CreateScriptCategory(tenantId, catReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Script Category", fmt.Sprintf("Could not create script category: %s", err))
		return
	}

	plan.Id = types.StringValue(strconv.Itoa(created.Id))
	plan.Name = types.StringValue(created.Name)

	// Preserve the parent_id from plan since the API response may not echo it back
	if plan.ParentId.IsNull() || plan.ParentId.IsUnknown() {
		plan.ParentId = types.Int64Value(0)
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *ScriptCategoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ScriptCategoryModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(state.Client)
	categoryId, err := strconv.Atoi(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse category ID '%s': %s", state.Id.ValueString(), err))
		return
	}

	result, err := r.apiClient.GetScriptCategory(tenantId, categoryId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Script Category", fmt.Sprintf("Could not read script category %d: %s", categoryId, err))
		return
	}

	state.Name = types.StringValue(result.Category.Name)
	state.ParentId = types.Int64Value(int64(result.ParentId))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *ScriptCategoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ScriptCategoryModel
	var state ScriptCategoryModel
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	tenantId := r.resolveTenantId(plan.Client)

	categoryId, err := strconv.Atoi(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse category ID: %s", err))
		return
	}

	catReq := buildCategoryRequest(plan)
	catReq.Id = categoryId // required for PUT

	updated, err := r.apiClient.UpdateScriptCategory(tenantId, catReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Script Category", fmt.Sprintf("Could not update script category %d: %s", categoryId, err))
		return
	}

	plan.Id = types.StringValue(strconv.Itoa(updated.Id))
	plan.Name = types.StringValue(updated.Name)
	if plan.ParentId.IsNull() || plan.ParentId.IsUnknown() {
		plan.ParentId = types.Int64Value(0)
	}

	diags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *ScriptCategoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ScriptCategoryModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.resolveTenantId(state.Client)
	categoryId, err := strconv.Atoi(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse category ID: %s", err))
		return
	}

	err = r.apiClient.DeleteScriptCategory(tenantId, categoryId)
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Script Category", fmt.Sprintf("Could not delete script category %d: %s", categoryId, err))
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles importing an existing resource.
// Import ID format: "{categoryId}" or "{tenantId}/{categoryId}"
func (r *ScriptCategoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")

	var tenantId string
	var rawCategoryId string

	switch len(parts) {
	case 1:
		tenantId = r.apiClient.TenantId
		rawCategoryId = parts[0]
	case 2:
		tenantId = parts[0]
		rawCategoryId = parts[1]
	default:
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Expected format: '{categoryId}' or '{tenantId}/{categoryId}'",
		)
		return
	}

	categoryId, err := strconv.Atoi(rawCategoryId)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("Category ID must be numeric: %s", err))
		return
	}

	result, err := r.apiClient.GetScriptCategory(tenantId, categoryId)
	if err != nil {
		resp.Diagnostics.AddError("Error Importing Script Category", fmt.Sprintf("Could not import category: %s", err))
		return
	}

	state := ScriptCategoryModel{
		Id:       types.StringValue(strconv.Itoa(result.Category.Id)),
		Name:     types.StringValue(result.Category.Name),
		ParentId: types.Int64Value(int64(result.ParentId)),
	}

	if len(parts) == 2 {
		state.Client = types.StringValue(tenantId)
	} else {
		state.Client = types.StringNull()
	}

	resp.State.Set(ctx, &state)
}
