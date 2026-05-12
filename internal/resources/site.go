// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/HPE/terraform-provider-opsramp/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &SiteResource{}
var _ resource.ResourceWithImportState = &SiteResource{}

// SiteResource defines the resource implementation.
type SiteResource struct {
	apiClient *client.OpsRampClient
}

// SiteModel maps Terraform schema attributes to the provider model.
type SiteModel struct {
	Client         types.String `tfsdk:"client"`
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Parent         types.String `tfsdk:"parent"`
	Description    types.String `tfsdk:"description"`
	Address        types.String `tfsdk:"address"`
	State          types.String `tfsdk:"state"`
	City           types.String `tfsdk:"city"`
	Country        types.String `tfsdk:"country"`
	Zip            types.String `tfsdk:"zip"`
	PrimaryContact types.Object `tfsdk:"primary_contact"`
	PhoneNumber    types.String `tfsdk:"phone_number"`
	PhoneExtension types.String `tfsdk:"phone_extension"`
	SearchQuery    types.String `tfsdk:"search_query"`
	Resources      types.Set    `tfsdk:"resources"`
}

// ContactModel represents the primary contact
type ContactModel struct {
	Id          types.String `tfsdk:"id"`
	LoginName   types.String `tfsdk:"login_name"`
	LastName    types.String `tfsdk:"last_name"`
	FirstName   types.String `tfsdk:"first_name"`
	Email       types.String `tfsdk:"email"`
	PhoneNumber types.String `tfsdk:"phone_number"`
}

var contactAttrTypes = map[string]attr.Type{
	"id":         types.StringType,
	"login_name": types.StringType,
	"first_name": types.StringType,
	"last_name":  types.StringType,
	"email":      types.StringType,
}

// NewSite creates a new instance of the resource.
func NewSite() resource.Resource {
	return &SiteResource{}
}

// Metadata returns the resource type name.
func (r *SiteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

// Schema defines the schema for the resource.
func (r *SiteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Site resource for grouping and organizing monitored resources by physical or logical location.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client (tenant) UUID where this site should be created. If not specified, uses the provider's tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the site.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the site.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "A description of the site.",
			},
			"address": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The physical address of the site.",
			},
			"city": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The city where the site is located.",
			},
			"state": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The state or region where the site is located.",
			},
			"country": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The country where the site is located.",
			},
			"zip": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The postal code for the site.",
			},
			"phone_number": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The primary phone number for the site.",
			},
			"phone_extension": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The phone extension for the site.",
			},
			"parent": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The parent site, if this site is nested under another.",
			},
			"primary_contact": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The primary contact person for the site.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The user ID of the primary contact.",
					},
					"login_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The login name of the primary contact.",
					},
					"first_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The first name of the primary contact.",
					},
					"last_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The last name of the primary contact.",
					},
					"email": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The email of the primary contact.",
					},
				},
			},
			"search_query": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The search query to filter resources for this site.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resources": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Set of resource IDs to assign to this site.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure prepares the resource with the API client.
func (r *SiteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func buildSiteResourceRefs(resources types.Set) []client.SiteResource {
	if resources.IsNull() || resources.IsUnknown() {
		return nil
	}

	ids := setToStringSlice(resources)
	refs := make([]client.SiteResource, 0, len(ids))
	for _, id := range ids {
		refs = append(refs, client.SiteResource{Id: id})
	}

	return refs
}

func mapSiteResourceRefsToSet(ctx context.Context, refs []client.SiteResource) (types.Set, error) {
	ids := make([]string, 0, len(refs))
	for _, ref := range refs {
		ids = append(ids, ref.Id)
	}

	set, diags := types.SetValueFrom(ctx, types.StringType, ids)
	if diags.HasError() {
		return types.SetNull(types.StringType), fmt.Errorf("failed to map site resources to state")
	}

	return set, nil
}

// Create handles the creation of the resource.
func (r *SiteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SiteModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	// Build the site object
	site := client.Site{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		Address:        plan.Address.ValueString(),
		City:           plan.City.ValueString(),
		State:          plan.State.ValueString(),
		Country:        plan.Country.ValueString(),
		Zip:            plan.Zip.ValueString(),
		PhoneNumber:    plan.PhoneNumber.ValueString(),
		PhoneExtension: plan.PhoneExtension.ValueString(),
	}

	// Handle parent
	if !plan.Parent.IsNull() && !plan.Parent.IsUnknown() {
		parentId := plan.Parent.ValueString()

		if parentId != "" {
			parentIdInt, err := strconv.ParseInt(parentId, 10, 64)
			if err != nil {
				resp.Diagnostics.AddError("Invalid Parent ID", fmt.Sprintf("Parent ID must be a number: %s", err))
				return
			}
			site.Parent = &client.SiteParent{
				Id: parentIdInt,
			}
		}
	}

	// Handle primary contact
	if !plan.PrimaryContact.IsNull() && !plan.PrimaryContact.IsUnknown() {
		attrs := plan.PrimaryContact.Attributes()
		var contactId string
		var contactLogin string
		if id, ok := attrs["id"]; ok && !id.IsNull() {
			contactId = id.(types.String).ValueString()
		}
		if login, ok := attrs["login_name"]; ok && !login.IsNull() {
			contactLogin = login.(types.String).ValueString()
		}
		if contactId != "" {
			site.PrimaryContact = &client.SiteContact{
				Id:        contactId,
				LoginName: contactLogin,
			}
		}
	}

	// Handle search query
	if !plan.SearchQuery.IsNull() && !plan.SearchQuery.IsUnknown() && plan.SearchQuery.ValueString() != "" {
		site.FilterCriteria = &client.SiteFilter{
			SearchQuery: plan.SearchQuery.ValueString(),
		}
	}

	site.Resources = buildSiteResourceRefs(plan.Resources)

	created, err := r.apiClient.CreateSite(tenantId, site)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", err.Error())
		return
	}

	plan.Id = types.StringValue(fmt.Sprintf("%d", created.Id))
	plan.Name = types.StringValue(created.Name)
	plan.Description = types.StringValue(created.Description)
	plan.Address = types.StringValue(created.Address)
	plan.City = types.StringValue(created.City)
	plan.State = types.StringValue(created.State)
	plan.Country = types.StringValue(created.Country)
	plan.Zip = types.StringValue(created.Zip)
	plan.PhoneNumber = types.StringValue(created.PhoneNumber)
	plan.PhoneExtension = types.StringValue(created.PhoneExtension)

	// Set parent if returned
	if created.Parent != nil && created.Parent.Id != 0 {
		plan.Parent = types.StringValue(fmt.Sprintf("%d", created.Parent.Id))
	} else {
		plan.Parent = types.StringValue("")
	}

	// Set primary contact if returned
	if created.PrimaryContact != nil && created.PrimaryContact.Id != "" {
		contactObj, diags := types.ObjectValue(contactAttrTypes, map[string]attr.Value{
			"id":         types.StringValue(created.PrimaryContact.Id),
			"login_name": types.StringValue(created.PrimaryContact.LoginName),
			"first_name": types.StringValue(created.PrimaryContact.FirstName),
			"last_name":  types.StringValue(created.PrimaryContact.LastName),
			"email":      types.StringValue(created.PrimaryContact.Email),
		})
		resp.Diagnostics.Append(diags...)
		plan.PrimaryContact = contactObj
	} else {
		plan.PrimaryContact = types.ObjectNull(contactAttrTypes)
	}

	// Set search query if returned
	plan.SearchQuery = types.StringValue("")
	if created.FilterCriteria != nil && created.FilterCriteria.SearchQuery != "" {
		plan.SearchQuery = types.StringValue(created.FilterCriteria.SearchQuery)
	}

	plan.Resources, err = mapSiteResourceRefsToSet(ctx, created.Resources)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *SiteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SiteModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	existing, err := r.apiClient.GetSite(tenantId, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "No site found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	// Update state from API response
	state.Name = types.StringValue(existing.Name)
	state.Description = types.StringValue(existing.Description)
	state.Address = types.StringValue(existing.Address)
	state.City = types.StringValue(existing.City)
	state.State = types.StringValue(existing.State)
	state.Country = types.StringValue(existing.Country)
	state.Zip = types.StringValue(existing.Zip)
	state.PhoneNumber = types.StringValue(existing.PhoneNumber)
	state.PhoneExtension = types.StringValue(existing.PhoneExtension)

	// Set parent if returned
	if existing.Parent != nil && existing.Parent.Id != 0 {
		state.Parent = types.StringValue(fmt.Sprintf("%d", existing.Parent.Id))
	} else {
		state.Parent = types.StringValue("")
	}

	// Set primary contact if returned
	if existing.PrimaryContact != nil && existing.PrimaryContact.Id != "" {
		contactObj, diags := types.ObjectValue(contactAttrTypes, map[string]attr.Value{
			"id":         types.StringValue(existing.PrimaryContact.Id),
			"login_name": types.StringValue(existing.PrimaryContact.LoginName),
			"first_name": types.StringValue(existing.PrimaryContact.FirstName),
			"last_name":  types.StringValue(existing.PrimaryContact.LastName),
			"email":      types.StringValue(existing.PrimaryContact.Email),
		})
		resp.Diagnostics.Append(diags...)
		state.PrimaryContact = contactObj
	} else {
		state.PrimaryContact = types.ObjectNull(contactAttrTypes)
	}

	// Set search query if returned
	state.SearchQuery = types.StringValue("")
	if existing.FilterCriteria != nil && existing.FilterCriteria.SearchQuery != "" {
		state.SearchQuery = types.StringValue(existing.FilterCriteria.SearchQuery)
	}

	state.Resources, err = mapSiteResourceRefsToSet(ctx, existing.Resources)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *SiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state SiteModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	// Build the site object
	site := client.Site{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		Address:        plan.Address.ValueString(),
		City:           plan.City.ValueString(),
		State:          plan.State.ValueString(),
		Country:        plan.Country.ValueString(),
		Zip:            plan.Zip.ValueString(),
		PhoneNumber:    plan.PhoneNumber.ValueString(),
		PhoneExtension: plan.PhoneExtension.ValueString(),
	}

	// Handle parent
	if !plan.Parent.IsNull() && !plan.Parent.IsUnknown() {
		parentId := plan.Parent.ValueString()
		if parentId != "" {
			parentIdInt, err := strconv.ParseInt(parentId, 10, 64)
			if err != nil {
				resp.Diagnostics.AddError("Invalid Parent ID", fmt.Sprintf("Parent ID must be a number: %s", err))
				return
			}
			site.Parent = &client.SiteParent{
				Id: parentIdInt,
			}
		}
	}

	// Handle primary contact
	if !plan.PrimaryContact.IsNull() && !plan.PrimaryContact.IsUnknown() {
		attrs := plan.PrimaryContact.Attributes()
		var contactId string
		var contactLogin string
		if id, ok := attrs["id"]; ok && !id.IsNull() {
			contactId = id.(types.String).ValueString()
		}
		if login, ok := attrs["login_name"]; ok && !login.IsNull() {
			contactLogin = login.(types.String).ValueString()
		}
		if contactId != "" {
			site.PrimaryContact = &client.SiteContact{
				Id:        contactId,
				LoginName: contactLogin,
			}
		}
	}

	// Handle search query
	if !plan.SearchQuery.IsNull() && !plan.SearchQuery.IsUnknown() && plan.SearchQuery.ValueString() != "" {
		site.FilterCriteria = &client.SiteFilter{
			SearchQuery: plan.SearchQuery.ValueString(),
		}
	}

	site.Resources = buildSiteResourceRefs(plan.Resources)

	updated, err := r.apiClient.UpdateSite(tenantId, state.Id.ValueString(), site)
	if err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())
		return
	}

	plan.Id = types.StringValue(fmt.Sprintf("%d", updated.Id))
	plan.Client = state.Client
	plan.Name = types.StringValue(updated.Name)
	plan.Description = types.StringValue(updated.Description)
	plan.Address = types.StringValue(updated.Address)
	plan.City = types.StringValue(updated.City)
	plan.State = types.StringValue(updated.State)
	plan.Country = types.StringValue(updated.Country)
	plan.Zip = types.StringValue(updated.Zip)
	plan.PhoneNumber = types.StringValue(updated.PhoneNumber)
	plan.PhoneExtension = types.StringValue(updated.PhoneExtension)

	// Set parent if returned
	if updated.Parent != nil && updated.Parent.Id != 0 {
		plan.Parent = types.StringValue(fmt.Sprintf("%d", updated.Parent.Id))
	} else {
		plan.Parent = types.StringValue("")
	}

	// Set primary contact if returned
	if updated.PrimaryContact != nil && updated.PrimaryContact.Id != "" {
		contactObj, diags := types.ObjectValue(contactAttrTypes, map[string]attr.Value{
			"id":         types.StringValue(updated.PrimaryContact.Id),
			"login_name": types.StringValue(updated.PrimaryContact.LoginName),
			"first_name": types.StringValue(updated.PrimaryContact.FirstName),
			"last_name":  types.StringValue(updated.PrimaryContact.LastName),
			"email":      types.StringValue(updated.PrimaryContact.Email),
		})
		resp.Diagnostics.Append(diags...)
		plan.PrimaryContact = contactObj
	} else {
		plan.PrimaryContact = types.ObjectNull(contactAttrTypes)
	}

	// Set search query if returned
	plan.SearchQuery = types.StringValue("")
	if updated.FilterCriteria != nil && updated.FilterCriteria.SearchQuery != "" {
		plan.SearchQuery = types.StringValue(updated.FilterCriteria.SearchQuery)
	}

	plan.Resources, err = mapSiteResourceRefsToSet(ctx, updated.Resources)
	if err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *SiteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SiteModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	err := r.apiClient.DeleteSite(tenantId, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState handles importing an existing resource.
func (r *SiteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	siteId := req.ID

	tenantId := r.apiClient.TenantId

	existing, err := r.apiClient.GetSite(tenantId, siteId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Site",
			fmt.Sprintf("Could not import site with ID '%s': %s", siteId, err),
		)
		return
	}

	state := SiteModel{
		Client:         types.StringValue(tenantId),
		Id:             types.StringValue(fmt.Sprintf("%d", existing.Id)),
		Name:           types.StringValue(existing.Name),
		Description:    types.StringValue(existing.Description),
		Address:        types.StringValue(existing.Address),
		City:           types.StringValue(existing.City),
		State:          types.StringValue(existing.State),
		Country:        types.StringValue(existing.Country),
		Zip:            types.StringValue(existing.Zip),
		PhoneNumber:    types.StringValue(existing.PhoneNumber),
		PhoneExtension: types.StringValue(existing.PhoneExtension),
	}

	// Set parent if returned
	if existing.Parent != nil && existing.Parent.Id != 0 {
		state.Parent = types.StringValue(fmt.Sprintf("%d", existing.Parent.Id))
	}

	// Set primary contact if returned
	if existing.PrimaryContact != nil && existing.PrimaryContact.Id != "" {
		contactObj, diags := types.ObjectValue(contactAttrTypes, map[string]attr.Value{
			"id":         types.StringValue(existing.PrimaryContact.Id),
			"login_name": types.StringValue(existing.PrimaryContact.LoginName),
			"first_name": types.StringValue(existing.PrimaryContact.FirstName),
			"last_name":  types.StringValue(existing.PrimaryContact.LastName),
			"email":      types.StringValue(existing.PrimaryContact.Email),
		})
		resp.Diagnostics.Append(diags...)
		state.PrimaryContact = contactObj
	}

	// Set search query if returned
	if existing.FilterCriteria != nil && existing.FilterCriteria.SearchQuery != "" {
		state.SearchQuery = types.StringValue(existing.FilterCriteria.SearchQuery)
	} else {
		state.SearchQuery = types.StringValue("")
	}

	resources, err := mapSiteResourceRefsToSet(ctx, existing.Resources)
	if err != nil {
		resp.Diagnostics.AddError("Error Importing Site", err.Error())
		return
	}
	state.Resources = resources

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
