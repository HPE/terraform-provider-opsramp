// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/HPE/terraform-provider-opsramp/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &ClientResource{}
var _ resource.ResourceWithImportState = &ClientResource{}

// ClientResource defines the resource implementation.
type ClientResource struct {
	client *client.OpsRampClient
}

// ClientModel maps Terraform schema attributes to the provider model.
type ClientModel struct {
	Id          types.Int64  `tfsdk:"id"`
	UniqueId    types.String `tfsdk:"unique_id"`
	Name        types.String `tfsdk:"name"`
	Address     types.String `tfsdk:"address"`
	City        types.String `tfsdk:"city"`
	State       types.String `tfsdk:"state"`
	Country     types.String `tfsdk:"country"`
	Zip         types.String `tfsdk:"zip"`
	TimeZone    types.String `tfsdk:"time_zone"`
	PhoneNumber types.String `tfsdk:"phone_number"`
	ExtOrgId    types.String `tfsdk:"ext_org_id"`
	Addons      types.Set    `tfsdk:"addons"`
	Packages    types.Set    `tfsdk:"packages"`
}

func stringSliceFromSet(ctx context.Context, set types.Set) ([]string, error) {
	if set.IsNull() || set.IsUnknown() {
		return nil, nil
	}

	values := make([]string, 0)
	diags := set.ElementsAs(ctx, &values, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to parse set values")
	}

	return values, nil
}

func stringSetFromSlice(values []string) (types.Set, error) {
	elements := make([]attr.Value, len(values))
	for i, value := range values {
		elements[i] = types.StringValue(value)
	}

	set, diags := types.SetValue(types.StringType, elements)
	if diags.HasError() {
		return types.SetNull(types.StringType), fmt.Errorf("failed to build set value")
	}

	return set, nil
}

// NewClient creates a new instance of the resource.
func NewClient() resource.Resource {
	return &ClientResource{}
}

// Metadata returns the resource type name.
func (r *ClientResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client"
}

// Schema defines the schema for the resource.
func (r *ClientResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Client (sub-tenant) resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the client.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"unique_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "A unique identifier for the client.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the client.",
			},
			"address": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The address of the client.",
			},
			"city": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The city of the client.",
			},
			"state": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The state of the client.",
			},
			"country": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The country of the client.",
			},
			"zip": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The zip code of the client.",
			},
			"time_zone": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The time zone of the client.",
			},
			"phone_number": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The phone number of the client.",
			},
			"ext_org_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "External organization ID for the client.",
			},
			"addons": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of addons enabled in the Client.",
			},
			"packages": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of packages enabled in the Client.",
			},
		},
	}
}

// Configure prepares the resource with client.
func (r *ClientResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = c
}

// Create handles the creation of the resource.
func (r *ClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ClientModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addons, err := stringSliceFromSet(ctx, plan.Addons)
	if err != nil {
		resp.Diagnostics.AddError("Invalid addons value", err.Error())
		return
	}

	packages, err := stringSliceFromSet(ctx, plan.Packages)
	if err != nil {
		resp.Diagnostics.AddError("Invalid packages value", err.Error())
		return
	}

	createClient := client.CreateClient{
		Name:        plan.Name.ValueString(),
		UniqueId:    plan.UniqueId.ValueString(),
		Address:     plan.Address.ValueString(),
		City:        plan.City.ValueString(),
		State:       plan.State.ValueString(),
		Country:     plan.Country.ValueString(),
		Zip:         plan.Zip.ValueString(),
		TimeZone:    plan.TimeZone.ValueString(),
		PhoneNumber: plan.PhoneNumber.ValueString(),
		ExtOrgId:    plan.ExtOrgId.ValueString(),
		Addons:      addons,
		Packages:    packages,
	}

	created, err := r.client.CreateClient(createClient)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", err.Error())
		return
	}

	plan.Id = types.Int64Value(int64(created.Id))
	plan.UniqueId = types.StringValue(created.UniqueId)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *ClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ClientModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	existing, err := r.client.GetClient(state.UniqueId.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	state.Id = types.Int64Value(int64(existing.Id))
	state.UniqueId = types.StringValue(existing.UniqueId)
	state.Name = types.StringValue(existing.Name)
	state.Address = types.StringValue(existing.Address)
	state.Country = types.StringValue(existing.Country)
	state.TimeZone = types.StringValue(existing.TimeZone)

	// Only update optional fields if API returns non-empty values
	if existing.City != "" {
		state.City = types.StringValue(existing.City)
	}
	if existing.State != "" {
		state.State = types.StringValue(existing.State)
	}
	if existing.Zip != "" {
		state.Zip = types.StringValue(existing.Zip)
	}
	if existing.PhoneNumber != "" {
		state.PhoneNumber = types.StringValue(existing.PhoneNumber)
	}
	if existing.ExtOrgId != "" {
		state.ExtOrgId = types.StringValue(existing.ExtOrgId)
	}

	state.Addons, err = stringSetFromSlice(existing.Addons)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	state.Packages, err = stringSetFromSlice(existing.Packages)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *ClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ClientModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addons, err := stringSliceFromSet(ctx, plan.Addons)
	if err != nil {
		resp.Diagnostics.AddError("Invalid addons value", err.Error())
		return
	}

	packages, err := stringSliceFromSet(ctx, plan.Packages)
	if err != nil {
		resp.Diagnostics.AddError("Invalid packages value", err.Error())
		return
	}

	updateClient := client.CreateClient{
		Name:        plan.Name.ValueString(),
		Address:     plan.Address.ValueString(),
		City:        plan.City.ValueString(),
		State:       plan.State.ValueString(),
		Country:     plan.Country.ValueString(),
		Zip:         plan.Zip.ValueString(),
		TimeZone:    plan.TimeZone.ValueString(),
		PhoneNumber: plan.PhoneNumber.ValueString(),
		ExtOrgId:    plan.ExtOrgId.ValueString(),
		Addons:      addons,
		Packages:    packages,
	}

	updated, err := r.client.UpdateClient(state.UniqueId.ValueString(), updateClient)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Client",
			fmt.Sprintf("Could not update client: %s", err),
		)
		return
	}

	plan.Id = state.Id
	plan.UniqueId = types.StringValue(updated.UniqueId)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *ClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ClientModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteClient(state.UniqueId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles importing an existing resource.
func (r *ClientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importId := req.ID

	existing, err := r.client.GetClient(importId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Client",
			fmt.Sprintf("Could not import client with ID '%s': %s", importId, err),
		)
		return
	}

	state := ClientModel{
		Id:       types.Int64Value(int64(existing.Id)),
		UniqueId: types.StringValue(existing.UniqueId),
		Name:     types.StringValue(existing.Name),
		Address:  types.StringValue(existing.Address),
		Country:  types.StringValue(existing.Country),
		TimeZone: types.StringValue(existing.TimeZone),
	}

	// Only set optional fields if API returns non-empty values
	if existing.City != "" {
		state.City = types.StringValue(existing.City)
	}
	if existing.State != "" {
		state.State = types.StringValue(existing.State)
	}
	if existing.Zip != "" {
		state.Zip = types.StringValue(existing.Zip)
	}
	if existing.PhoneNumber != "" {
		state.PhoneNumber = types.StringValue(existing.PhoneNumber)
	}
	if existing.ExtOrgId != "" {
		state.ExtOrgId = types.StringValue(existing.ExtOrgId)
	}

	state.Addons, err = stringSetFromSlice(existing.Addons)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Client",
			fmt.Sprintf("Could not map addons for client with ID '%s': %s", importId, err),
		)
		return
	}

	state.Packages, err = stringSetFromSlice(existing.Packages)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Client",
			fmt.Sprintf("Could not map packages for client with ID '%s': %s", importId, err),
		)
		return
	}

	resp.State.Set(ctx, &state)
}
