// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"
	"strings"
	"time"

	"github.com/HPE/terraform-provider-opsramp/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces.
var _ resource.Resource = &DeviceGroupResource{}

// DeviceGroupModel maps Terraform schema attributes to the provider model.
type DeviceGroupModel struct {
	Client      types.String `tfsdk:"client"`
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	EntityType  types.String `tfsdk:"entity_type"`
	Parent      types.String `tfsdk:"parent"`
	SearchQuery types.String `tfsdk:"search_query"`
	Resources   types.Set    `tfsdk:"resources"`
}

// DeviceGroupResource defines the resource implementation.
type DeviceGroupResource struct {
	apiClient *client.OpsRampClient
}

// NewDeviceGroup creates a new instance of the resource.
func NewDeviceGroup() resource.Resource {
	return &DeviceGroupResource{}
}

// Metadata returns the resource type name.
func (r *DeviceGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_group"
}

// Recursive schema definition example
func (r *DeviceGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client (tenant) UUID where this device group should be created. If not specified, uses the provider's tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Unique identifier for the group.",
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"parent": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"entity_type": schema.StringAttribute{
				Computed: true,
			},
			"search_query": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The search query for selecting resources.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resources": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Set of resource IDs to assign to this device group.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}

}

// Configure prepares the resource with client.
func (r *DeviceGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.OpsRampClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *client.OpsRampClient",
		)
		return
	}

	r.apiClient = client
}

func translatePlanToDeviceGroupModel(plan DeviceGroupModel) client.DeviceGroupAPI {
	var parent *client.Parent

	if plan.Parent.ValueString() != "" {
		parent = &client.Parent{
			Id: plan.Parent.ValueString(),
		}
	}

	var filterCriteria *client.FilterCriteria
	if !plan.SearchQuery.IsNull() && !plan.SearchQuery.IsUnknown() {
		filterCriteria = &client.FilterCriteria{
			SearchQuery: plan.SearchQuery.ValueString(),
		}
	}

	deviceGroup := client.DeviceGroupAPI{
		Id:             plan.Id.ValueString(),
		Name:           plan.Name.ValueString(),
		Parent:         parent,
		EntityType:     "DEVICE_GROUP",
		FilterCriteria: filterCriteria,
	}

	return deviceGroup
}

// setToStringSlice converts a types.Set of strings to a []string.
func setToStringSlice(s types.Set) []string {
	elements := s.Elements()
	result := make([]string, 0, len(elements))
	for _, e := range elements {
		if sv, ok := e.(types.String); ok {
			result = append(result, sv.ValueString())
		}
	}
	return result
}

// stringSetDiff returns elements in a that are not in b.
func stringSetDiff(a, b []string) []string {
	bSet := make(map[string]struct{}, len(b))
	for _, v := range b {
		bSet[v] = struct{}{}
	}
	var diff []string
	for _, v := range a {
		if _, found := bSet[v]; !found {
			diff = append(diff, v)
		}
	}
	return diff
}

func deviceGroupResourcesSet(ctx context.Context, ids []string) (types.Set, diag.Diagnostics) {
	if ids == nil {
		ids = []string{}
	}

	return types.SetValueFrom(ctx, types.StringType, ids)
}

func (r *DeviceGroupResource) readDeviceGroupResources(ctx context.Context, tenantId string, deviceGroupId string) (types.Set, diag.Diagnostics) {
	resources, err := r.apiClient.GetDeviceGroupChilds(tenantId, deviceGroupId)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("Read Error", err.Error())
		return types.SetUnknown(types.StringType), diags
	}

	ids := make([]string, 0, len(resources))
	for _, resource := range resources {
		// Empty type for child groups
		if resource.ResourceType != "" {
			ids = append(ids, resource.Uuid)
		}
	}

	return deviceGroupResourcesSet(ctx, ids)
}

func createDeviceGroupWithRetry(apiClient *client.OpsRampClient, tenantId string, req client.DeviceGroupAPI, hasParent bool) (*client.DeviceGroupAPI, error) {
	const maxAttempts = 4
	const retryDelay = 2 * time.Second

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		created, err := apiClient.CreateDeviceGroup(tenantId, req)
		if err == nil {
			return created, nil
		}

		lastErr = err
		// Retry only for transient server errors while attaching to a parent.
		if !hasParent || !strings.Contains(err.Error(), "status: 500") || attempt == maxAttempts {
			break
		}

		time.Sleep(retryDelay)
	}

	return nil, lastErr
}

func (r *DeviceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DeviceGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceGroup := translatePlanToDeviceGroupModel(plan)

	// Use client parameter if set, otherwise provider's tenant ID
	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	// Create the device group. Parent-scoped creation may intermittently return 500,
	// so retry a few times for that specific transient error.
	newDeviceGroup, err := createDeviceGroupWithRetry(r.apiClient, tenantId, deviceGroup, !plan.Parent.IsNull() && plan.Parent.ValueString() != "")
	if err != nil {
		resp.Diagnostics.AddError("Create Error", err.Error())
		return
	}

	// Assign the backend response directly into the state
	plan.Id = types.StringValue(newDeviceGroup.Id)
	plan.EntityType = types.StringValue(newDeviceGroup.EntityType)

	// Set search_query from API response
	plan.SearchQuery = types.StringValue("")
	if newDeviceGroup.FilterCriteria != nil && newDeviceGroup.FilterCriteria.SearchQuery != "" {
		plan.SearchQuery = types.StringValue(newDeviceGroup.FilterCriteria.SearchQuery)
	}

	if newDeviceGroup.Parent != nil {
		plan.Parent = types.StringValue(newDeviceGroup.Parent.Id)
	}

	// Add child resources if specified
	if !plan.Resources.IsNull() && !plan.Resources.IsUnknown() && len(plan.Resources.Elements()) > 0 {
		ids := setToStringSlice(plan.Resources)
		if err := r.apiClient.AddDeviceGroupChilds(tenantId, newDeviceGroup.Id, ids); err != nil {
			resp.Diagnostics.AddError("Error adding resources to device group", err.Error())
			return
		}
	}

	if !plan.Resources.IsNull() && !plan.Resources.IsUnknown() {
		// Preserve configured resources in state after successful assignment so
		// Terraform sees the same values it planned during apply.
	} else {
		plan.Resources, diags = r.readDeviceGroupResources(ctx, tenantId, newDeviceGroup.Id)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Save new state back to Terraform
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// --- Read ---
func (r *DeviceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DeviceGroupModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use client parameter if set, otherwise provider's tenant ID
	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	backendDeviceGroup, err := r.apiClient.GetDeviceGroup(tenantId, state.Id.ValueString())
	if err != nil && !strings.Contains(err.Error(), "No Device group exists") {
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	if backendDeviceGroup == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// has parent?
	var parentId types.String
	if backendDeviceGroup.Parent != nil {
		parentId = types.StringValue(backendDeviceGroup.Parent.Id)
	}

	// Set search_query from API response
	var searchQuery types.String
	if backendDeviceGroup.FilterCriteria != nil && backendDeviceGroup.FilterCriteria.SearchQuery != "" {
		searchQuery = types.StringValue(backendDeviceGroup.FilterCriteria.SearchQuery)
	}

	resources, diags := r.readDeviceGroupResources(ctx, tenantId, backendDeviceGroup.Id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState := DeviceGroupModel{
		Client:      state.Client,
		Id:          types.StringValue(backendDeviceGroup.Id),
		Name:        types.StringValue(backendDeviceGroup.Name),
		EntityType:  types.StringValue(backendDeviceGroup.EntityType),
		Parent:      parentId,
		SearchQuery: searchQuery,
		Resources:   resources,
	}

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
}

// Update modifies an existing resource.
func (r *DeviceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state DeviceGroupModel
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	deviceGroup := translatePlanToDeviceGroupModel(plan)

	// Use client parameter if set, otherwise provider's tenant ID
	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	// Update the device group
	newDeviceGroup, err := r.apiClient.CreateDeviceGroup(tenantId, deviceGroup)
	if err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())
		return
	}

	// Assign the backend response directly into the state
	state.Id = types.StringValue(newDeviceGroup.Id)
	state.EntityType = types.StringValue(newDeviceGroup.EntityType)

	// Set search_query from API response
	if newDeviceGroup.FilterCriteria != nil && newDeviceGroup.FilterCriteria.SearchQuery != "" {
		state.SearchQuery = types.StringValue(newDeviceGroup.FilterCriteria.SearchQuery)
	}

	if newDeviceGroup.Parent != nil {
		state.Parent = types.StringValue(newDeviceGroup.Parent.Id)
	}

	// Reconcile child resources only when the desired plan is explicitly known.
	if !plan.Resources.IsNull() && !plan.Resources.IsUnknown() {
		oldIds := setToStringSlice(state.Resources)
		newIds := setToStringSlice(plan.Resources)
		toAdd := stringSetDiff(newIds, oldIds)
		toRemove := stringSetDiff(oldIds, newIds)

		if len(toAdd) > 0 {
			if err := r.apiClient.AddDeviceGroupChilds(tenantId, newDeviceGroup.Id, toAdd); err != nil {
				resp.Diagnostics.AddError("Error adding resources to device group", err.Error())
				return
			}
		}
		if len(toRemove) > 0 {
			if err := r.apiClient.RemoveDeviceGroupChilds(tenantId, newDeviceGroup.Id, toRemove); err != nil {
				resp.Diagnostics.AddError("Error removing resources from device group", err.Error())
				return
			}
		}
	}

	var diags diag.Diagnostics
	if !plan.Resources.IsNull() && !plan.Resources.IsUnknown() {
		state.Resources = plan.Resources
	} else {
		state.Resources, diags = r.readDeviceGroupResources(ctx, tenantId, newDeviceGroup.Id)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Save
	resp.State.Set(ctx, &state)
}

// Delete removes the resource from the API.
func (r *DeviceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DeviceGroupModel
	req.State.Get(ctx, &state)

	// Use client parameter if set, otherwise provider's tenant ID
	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	err := r.apiClient.DeleteDeviceGroup(tenantId, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}
