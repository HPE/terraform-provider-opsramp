// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"context"
	"regexp"
	"strings"

	"github.com/HPE/terraform-provider-opsramp/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure implementation satisfies the expected interfaces
var _ resource.Resource = &AlertEscalationPolicyResource{}

// AlertEscalationPolicyResource defines the resource implementation.
type AlertEscalationPolicyResource struct {
	apiClient *client.OpsRampClient
}

// AlertEscalationPolicyModel maps Terraform schema attributes to the provider model.
type AlertEscalationPolicyModel struct {
	Client         types.String                   `tfsdk:"client"`
	Id             types.String                   `tfsdk:"id"`
	Name           types.String                   `tfsdk:"name"`
	Description    types.String                   `tfsdk:"description"`
	TenantScope    types.String                   `tfsdk:"tenant_scope"`
	Precedence     types.Int64                    `tfsdk:"precedence"`
	EscalationType types.String                   `tfsdk:"escalation_type"`
	PolicyType     types.String                   `tfsdk:"policy_type"`
	EnabledMode    types.String                   `tfsdk:"enabled_mode"`
	AllClients     types.Bool                     `tfsdk:"all_clients"`
	Scope          types.String                   `tfsdk:"scope"`
	Resources      []EscalationResourceModel      `tfsdk:"resources"`
	Escalations    []EscalationLevelModel         `tfsdk:"escalations"`
	FilterCriteria *EscalationFilterCriteriaModel `tfsdk:"filter_criteria"`
}

type EscalationResourceModel struct {
	Id   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

type EscalationFilterCriteriaModel struct {
	MatchingType        types.String `tfsdk:"matching_type"`
	SearchQuery         types.String `tfsdk:"search_query"`
	ResourceSearchQuery types.String `tfsdk:"resource_search_query"`
}

type EscalationLevelModel struct {
	WaitMins               types.Int64                    `tfsdk:"wait_mins"`
	Action                 types.String                   `tfsdk:"action"`
	Priority               types.String                   `tfsdk:"priority"`
	RepeatFrequency        types.Int64                    `tfsdk:"repeat_frequency"`
	NotifyLimitCount       types.Int64                    `tfsdk:"notify_limit_count"`
	NotificationType       types.String                   `tfsdk:"notification_type"`
	NotificationTemplateId types.String                   `tfsdk:"notification_template_id"`
	Recipients             []EscalationRecipientModel     `tfsdk:"recipients"`
	Incident               *EscalationIncidentModel       `tfsdk:"incident"`
	UpdateIncident         *EscalationUpdateIncidentModel `tfsdk:"update_incident"`
}

type EscalationRecipientModel struct {
	Id   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

type EscalationIncidentModel struct {
	Priority            types.String           `tfsdk:"priority"`
	Subject             types.String           `tfsdk:"subject"`
	Description         types.String           `tfsdk:"description"`
	AssigneeGroupId     types.String           `tfsdk:"assignee_group_id"`
	AssignedUserId      types.String           `tfsdk:"assigned_user_id"`
	CategoryId          types.String           `tfsdk:"category_id"`
	SubCategoryId       types.String           `tfsdk:"sub_category_id"`
	BusinessImpactId    types.String           `tfsdk:"business_impact_id"`
	UrgencyId           types.String           `tfsdk:"urgency_id"`
	KnowledgeArticleIds []types.String         `tfsdk:"knowledge_article_ids"`
	Cc                  types.String           `tfsdk:"cc"`
	ToMail              *EscalationToMailModel `tfsdk:"to_mail"`
}

type EscalationToMailModel struct {
	UsersIds      []types.String `tfsdk:"users_ids"`
	UserGroupsIds []types.String `tfsdk:"user_groups_ids"`
	RostersIds    []types.String `tfsdk:"rosters_ids"`
}

type EscalationUpdateIncidentModel struct {
	UpdateWhenAlertStateChange         types.Bool                    `tfsdk:"update_when_alert_state_change"`
	UpdateForEveryRepeatAlert          types.Bool                    `tfsdk:"update_for_every_repeat_alert"`
	UpdateWithRuleWhenAlertStateChange types.Bool                    `tfsdk:"update_with_rule_when_alert_state_change"`
	UpdateWithRuleForEveryRepeatAlert  types.Bool                    `tfsdk:"update_with_rule_for_every_repeat_alert"`
	UpdateIncidentSubject              types.Bool                    `tfsdk:"update_incident_subject"`
	UpdateIncidentSubjectWithRule      types.Bool                    `tfsdk:"update_incident_subject_with_rule"`
	AutoResolveIncident                types.Bool                    `tfsdk:"auto_resolve_incident"`
	AutoResolveUnassignedIncident      types.Bool                    `tfsdk:"auto_resolve_unassigned_incident"`
	AutoHealWaitTime                   types.Int64                   `tfsdk:"auto_heal_wait_time"`
	UpdatePriorityByMLConfiguration    types.Bool                    `tfsdk:"update_priority_by_ml_configuration"`
	PriorityRules                      []EscalationPriorityRuleModel `tfsdk:"priority_rules"`
}

type EscalationPriorityRuleModel struct {
	Key              types.String `tfsdk:"key"`
	Operator         types.String `tfsdk:"operator"`
	Value            types.String `tfsdk:"value"`
	BusinessImpactId types.String `tfsdk:"business_impact_id"`
	UrgencyId        types.String `tfsdk:"urgency_id"`
	Priority         types.String `tfsdk:"priority"`
}

// NewAlertEscalationPolicy creates a new instance of the resource.
func NewAlertEscalationPolicy() resource.Resource {
	return &AlertEscalationPolicyResource{}
}

// Metadata returns the resource type name.
func (r *AlertEscalationPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_escalation_policy"
}

// Schema defines the schema for the resource.
func (r *AlertEscalationPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpsRamp Alert Escalation Policy. Defines automated notifications and incident creation for escalated alerts.",
		Attributes: map[string]schema.Attribute{
			"client": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The client (tenant) UUID where this policy should be created. If not specified, uses the provider's tenant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the alert escalation policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the alert escalation policy.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "A description of the alert escalation policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tenant_scope": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The tenant scope. Valid values: `CLIENT`, `MSP`.",
				Validators: []validator.String{
					stringvalidator.OneOf("CLIENT", "MSP"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"precedence": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The execution order of the policy.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"escalation_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Automated notification type. Valid values: `AUTOMATIC_UNTIL_ACKNOWLEDGED_CLOSED_SUPPRESSED_TICKETED`, `AUTOMATIC_UNTIL_ACKNOWLEDGED_CLOSED_SUPPRESSED`, `MANUAL`, `AUTOMATIC`.",
				Validators: []validator.String{
					stringvalidator.OneOf("AUTOMATIC_UNTIL_ACKNOWLEDGED_CLOSED_SUPPRESSED_TICKETED", "AUTOMATIC_UNTIL_ACKNOWLEDGED_CLOSED_SUPPRESSED", "MANUAL", "AUTOMATIC"),
				},
			},
			"policy_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The policy type. Value: `ESCALATION_POLICY`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("ESCALATION_POLICY"),
				},
			},
			"enabled_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The enabled mode. Valid values: `ON`, `OFF`, `RECOMMEND`, `OBSERVED`.",
				Validators: []validator.String{
					stringvalidator.OneOf("ON", "OFF", "RECOMMEND", "OBSERVED"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"all_clients": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the policy applies to all clients.",
			},
			"scope": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The scope (client or tenant UUID) of the escalation policy.",
			},
			"resources": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Resources the policy applies to.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The resource ID (e.g. client UUID).",
						},
						"type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The resource type (e.g. `CLIENT`, `PARTNER`).",
							Validators: []validator.String{
								stringvalidator.OneOf("CLIENT", "PARTNER"),
							},
						},
					},
				},
			},
			"escalations": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "Escalation steps. Each step can be a NOTIFICATION or INCIDENT action.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"wait_mins": schema.Int64Attribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "Minutes to wait before executing this escalation step.",
						},
						"action": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The action type. Valid values: `NOTIFICATION`, `INCIDENT`.",
							Validators: []validator.String{
								stringvalidator.OneOf("NOTIFICATION", "INCIDENT"),
							},
						},
						"priority": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The priority for notifications (e.g. `Normal`, `High`, `Low`, `Urgent`).",
							Validators: []validator.String{
								stringvalidator.OneOf("Normal", "Low", "High", "Urgent"),
							},
						},
						"repeat_frequency": schema.Int64Attribute{
							Optional:            true,
							MarkdownDescription: "Number of minutes between repeat notifications.",
						},
						"notify_limit_count": schema.Int64Attribute{
							Optional:            true,
							MarkdownDescription: "Maximum number of notifications to send.",
						},
						"notification_type": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "Notification template type (e.g. `basic`, `advanced`).",
							Validators: []validator.String{
								stringvalidator.OneOf("basic", "advanced"),
							},
						},
						"notification_template_id": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The ID of the notification template.",
						},
						"recipients": schema.SetNestedAttribute{
							Optional:            true,
							MarkdownDescription: "Recipients for notification actions.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "The recipient ID.",
									},
									"type": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "The recipient type (e.g. `USERGROUP`, `USER`, `USERGROUP_DL`).",
										Validators: []validator.String{
											stringvalidator.OneOf("USER", "USERGROUP", "USERGROUP_DL"),
										},
									},
								},
							},
						},
						"incident": schema.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "Incident creation settings. Used when action is `INCIDENT`.",
							Attributes: map[string]schema.Attribute{
								"priority": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Incident priority (e.g. `Normal`, `Low`, `High`, `Urgent`).",
									Validators: []validator.String{
										stringvalidator.OneOf("Normal", "Low", "High", "Urgent"),
									},
								},
								"subject": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Incident subject. Supports placeholders like `$alert.subject`.",
								},
								"description": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Incident description. Supports placeholders like `$alert.description`.",
								},
								"assignee_group_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The uniqueId of the assignee user group.",
								},
								"assigned_user_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The ID of the assigned user.",
								},
								"category_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The uniqueId of the service desk category.",
								},
								"sub_category_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The uniqueId of the service desk sub-category.",
								},
								"business_impact_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The uniqueId of the business impact.",
								},
								"urgency_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The uniqueId of the urgency.",
								},
								"knowledge_article_ids": schema.SetAttribute{
									Optional:            true,
									ElementType:         types.StringType,
									MarkdownDescription: "List of knowledge base article IDs to attach.",
								},
								"cc": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "CC email addresses for the incident.",
									Validators: []validator.String{
										stringvalidator.RegexMatches(
											regexp.MustCompile(`^([a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,})(\s*,\s*[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,})*$`),
											"must be a valid email address or a comma-separated list of valid email addresses",
										),
									},
								},
								"to_mail": schema.SingleNestedAttribute{
									Optional:            true,
									MarkdownDescription: "Email notification recipients for the incident.",
									Attributes: map[string]schema.Attribute{
										"users_ids": schema.SetAttribute{
											Optional:            true,
											ElementType:         types.StringType,
											MarkdownDescription: "List of user IDs for email notifications.",
										},
										"user_groups_ids": schema.SetAttribute{
											Optional:            true,
											ElementType:         types.StringType,
											MarkdownDescription: "List of user group IDs for email notifications.",
										},
										"rosters_ids": schema.SetAttribute{
											Optional:            true,
											ElementType:         types.StringType,
											MarkdownDescription: "List of roster IDs for email notifications.",
										},
									},
								},
							},
						},
						"update_incident": schema.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "Incident update settings. Used when action is `INCIDENT`.",
							Attributes: map[string]schema.Attribute{
								"update_when_alert_state_change": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: "Update incident when alert state changes.",
								},
								"update_for_every_repeat_alert": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: "Update incident for every repeat alert.",
								},
								"update_with_rule_when_alert_state_change": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: "Update with rule when alert state changes.",
								},
								"update_with_rule_for_every_repeat_alert": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: "Update with rule for every repeat alert.",
								},
								"update_incident_subject": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "Whether to update incident subject on alert changes.",
								},
								"update_incident_subject_with_rule": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "Whether to update incident subject using rules.",
								},
								"auto_resolve_incident": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: "Whether to auto-resolve the incident when the alert clears.",
								},
								"auto_resolve_unassigned_incident": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: "Whether to auto-resolve unassigned incidents.",
								},
								"auto_heal_wait_time": schema.Int64Attribute{
									Optional:            true,
									MarkdownDescription: "Wait time in minutes before auto-healing.",
								},
								"update_priority_by_ml_configuration": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: "Update priority by ML configuration.",
								},
								"priority_rules": schema.ListNestedAttribute{
									Optional:            true,
									MarkdownDescription: "Priority rules for incident updates.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"key": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "The alert property key (e.g. `currentState.code`).",
											},
											"operator": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "The operator (e.g. `Is`).",
											},
											"value": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "The value to match (e.g. `WARNING`, `CRITICAL`).",
											},
											"business_impact_id": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "The uniqueId of the business impact for this rule.",
											},
											"urgency_id": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "The uniqueId of the urgency for this rule.",
											},
											"priority": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "The resulting priority (e.g. `Low`, `Normal`, `High`).",
												Validators: []validator.String{
													stringvalidator.OneOf("Low", "Normal", "High"),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"filter_criteria": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Filter criteria to scope which alerts trigger this policy.",
				Attributes: map[string]schema.Attribute{
					"matching_type": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The matching type. Valid values: `ALL`, `ANY`.",
						Validators: []validator.String{
							stringvalidator.OneOf("ALL", "ANY"),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"search_query": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "OpsQL search query for alert filtering.",
					},
					"resource_search_query": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "OpsQL search query for resource filtering.",
					},
				},
			},
		},
	}
}

// Configure prepares the resource with the API client.
func (r *AlertEscalationPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func buildAlertEscalationPolicyRequest(plan AlertEscalationPolicyModel) client.AlertEscalationPolicy {
	policy := client.AlertEscalationPolicy{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		EscalationType: plan.EscalationType.ValueString(),
		PolicyType:     plan.PolicyType.ValueString(),
		EnabledMode:    plan.EnabledMode.ValueString(),
		AllClients:     plan.AllClients.ValueBool(),
		Resources:      []client.EscalationResource{},
	}

	if !plan.TenantScope.IsNull() && !plan.TenantScope.IsUnknown() {
		policy.TenantScope = plan.TenantScope.ValueString()
	}

	if !plan.Precedence.IsNull() && !plan.Precedence.IsUnknown() {
		policy.Precedence = int(plan.Precedence.ValueInt64())
	}

	if !plan.Scope.IsNull() && !plan.Scope.IsUnknown() {
		policy.Scope = &client.EscalationScope{
			Uuid: plan.Scope.ValueString(),
		}
	}

	for _, r := range plan.Resources {
		policy.Resources = append(policy.Resources, client.EscalationResource{
			Id:   r.Id.ValueString(),
			Type: r.Type.ValueString(),
		})
	}

	for _, e := range plan.Escalations {
		level := client.EscalationLevel{
			WaitMins: int(e.WaitMins.ValueInt64()),
			Action:   e.Action.ValueString(),
		}
		if !e.Priority.IsNull() {
			level.Priority = e.Priority.ValueString()
		}
		if !e.RepeatFrequency.IsNull() {
			level.RepeatFrequency = int(e.RepeatFrequency.ValueInt64())
		}
		if !e.NotifyLimitCount.IsNull() {
			level.NotifyLimitCount = int(e.NotifyLimitCount.ValueInt64())
		}
		if !e.NotificationType.IsNull() {
			level.NotificationType = e.NotificationType.ValueString()
		}
		if !e.NotificationTemplateId.IsNull() {
			level.NotificationTemplateId = e.NotificationTemplateId.ValueString()
		}
		for _, rec := range e.Recipients {
			level.Recipients = append(level.Recipients, client.EscalationRecipient{
				Id:   rec.Id.ValueString(),
				Type: rec.Type.ValueString(),
			})
		}
		if e.Incident != nil {
			inc := &client.EscalationIncident{
				Priority:    e.Incident.Priority.ValueString(),
				Subject:     e.Incident.Subject.ValueString(),
				Description: e.Incident.Description.ValueString(),
			}
			if !e.Incident.AssigneeGroupId.IsNull() && e.Incident.AssigneeGroupId.ValueString() != "" {
				inc.AssigneeGroup = &client.EscalationUniqueRef{UniqueId: e.Incident.AssigneeGroupId.ValueString()}
			}
			if !e.Incident.AssignedUserId.IsNull() && e.Incident.AssignedUserId.ValueString() != "" {
				inc.AssignedUser = &client.EscalationUserRef{
					Id: e.Incident.AssignedUserId.ValueString(),
				}
			}
			if !e.Incident.CategoryId.IsNull() && e.Incident.CategoryId.ValueString() != "" {
				inc.Category = &client.EscalationUniqueRef{UniqueId: e.Incident.CategoryId.ValueString()}
			}
			if !e.Incident.SubCategoryId.IsNull() && e.Incident.SubCategoryId.ValueString() != "" {
				inc.SubCategory = &client.EscalationUniqueRef{UniqueId: e.Incident.SubCategoryId.ValueString()}
			}
			if !e.Incident.BusinessImpactId.IsNull() && e.Incident.BusinessImpactId.ValueString() != "" {
				inc.BusinessImpact = &client.EscalationUniqueRef{UniqueId: e.Incident.BusinessImpactId.ValueString()}
			}
			if !e.Incident.UrgencyId.IsNull() && e.Incident.UrgencyId.ValueString() != "" {
				inc.Urgency = &client.EscalationUniqueRef{UniqueId: e.Incident.UrgencyId.ValueString()}
			}
			for _, kid := range e.Incident.KnowledgeArticleIds {
				inc.KnowledgeArticleIds = append(inc.KnowledgeArticleIds, kid.ValueString())
				inc.AttachedArticles = append(inc.AttachedArticles, client.EscalationArticleRef{Id: kid.ValueString()})
			}
			if !e.Incident.Cc.IsNull() && e.Incident.Cc.ValueString() != "" {
				inc.Cc = e.Incident.Cc.ValueString()
			}
			if e.Incident.ToMail != nil {
				tm := &client.EscalationToMail{}
				for _, uid := range e.Incident.ToMail.UsersIds {
					tm.Users = append(tm.Users, client.EscalationUserRef{Id: uid.ValueString()})
				}
				// Build comma-separated IDs for backward compatibility
				var userIds []string
				for _, uid := range e.Incident.ToMail.UsersIds {
					userIds = append(userIds, uid.ValueString())
				}
				if len(userIds) > 0 {
					inc.ToMailUserIds = strings.Join(userIds, ",")
				}
				var groupIds []string
				for _, gid := range e.Incident.ToMail.UserGroupsIds {
					groupIds = append(groupIds, gid.ValueString())
				}
				if len(groupIds) > 0 {
					inc.ToMailUserGroupIds = strings.Join(groupIds, ",")
				}
				var rosterIds []string
				for _, rid := range e.Incident.ToMail.RostersIds {
					rosterIds = append(rosterIds, rid.ValueString())
				}
				if len(rosterIds) > 0 {
					inc.ToMailRosterIds = strings.Join(rosterIds, ",")
				}
				inc.ToMail = tm
			}
			level.Incident = inc
		}
		if e.UpdateIncident != nil {
			ui := &client.EscalationUpdateIncident{
				UpdateWhenAlertStateChange:         e.UpdateIncident.UpdateWhenAlertStateChange.ValueBool(),
				UpdateForEveryRepeatAlert:          e.UpdateIncident.UpdateForEveryRepeatAlert.ValueBool(),
				UpdateWithRuleWhenAlertStateChange: e.UpdateIncident.UpdateWithRuleWhenAlertStateChange.ValueBool(),
				UpdateWithRuleForEveryRepeatAlert:  e.UpdateIncident.UpdateWithRuleForEveryRepeatAlert.ValueBool(),
				UpdateIncidentSubject:              e.UpdateIncident.UpdateIncidentSubject.ValueBool(),
				UpdateIncidentSubjectWithRule:      e.UpdateIncident.UpdateIncidentSubjectWithRule.ValueBool(),
				AutoResolveIncident:                e.UpdateIncident.AutoResolveIncident.ValueBool(),
				AutoResolveUnassignedIncident:      e.UpdateIncident.AutoResolveUnassignedIncident.ValueBool(),
				AutoHealWaitTime:                   int(e.UpdateIncident.AutoHealWaitTime.ValueInt64()),
				UpdatePriorityByMLConfiguration:    e.UpdateIncident.UpdatePriorityByMLConfiguration.ValueBool(),
			}
			for _, pr := range e.UpdateIncident.PriorityRules {
				rule := client.EscalationPriorityRule{
					Key:      pr.Key.ValueString(),
					Operator: pr.Operator.ValueString(),
					Value:    pr.Value.ValueString(),
					Priority: pr.Priority.ValueString(),
				}
				if !pr.BusinessImpactId.IsNull() && pr.BusinessImpactId.ValueString() != "" {
					rule.BusinessImpact = &client.EscalationUniqueRef{UniqueId: pr.BusinessImpactId.ValueString()}
				}
				if !pr.UrgencyId.IsNull() && pr.UrgencyId.ValueString() != "" {
					rule.Urgency = &client.EscalationUniqueRef{UniqueId: pr.UrgencyId.ValueString()}
				}
				ui.PriorityRules = append(ui.PriorityRules, rule)
			}
			level.UpdateIncident = ui
		}
		policy.Escalations = append(policy.Escalations, level)
	}

	if plan.FilterCriteria != nil {
		matchingType := "ALL"
		if !plan.FilterCriteria.MatchingType.IsNull() && !plan.FilterCriteria.MatchingType.IsUnknown() {
			if mt := strings.TrimSpace(plan.FilterCriteria.MatchingType.ValueString()); mt != "" {
				matchingType = mt
			}
		}

		policy.FilterCriteria = &client.EscalationFilterCriteria{
			MatchingType:        matchingType,
			SearchQuery:         plan.FilterCriteria.SearchQuery.ValueString(),
			ResourceSearchQuery: plan.FilterCriteria.ResourceSearchQuery.ValueString(),
		}
	}

	return policy
}

func mapAlertEscalationPolicyToState(resp *client.AlertEscalationPolicy, state *AlertEscalationPolicyModel) {
	state.Id = types.StringValue(resp.Id)
	state.Name = types.StringValue(resp.Name)
	state.Description = types.StringValue(resp.Description)
	state.EscalationType = types.StringValue(resp.EscalationType)
	state.PolicyType = types.StringValue(resp.PolicyType)
	state.EnabledMode = types.StringValue(resp.EnabledMode)
	state.AllClients = types.BoolValue(resp.AllClients)

	if resp.Precedence != 0 {
		state.Precedence = types.Int64Value(int64(resp.Precedence))
	}

	if len(resp.Resources) > 0 {
		state.Resources = make([]EscalationResourceModel, len(resp.Resources))
		for i, r := range resp.Resources {
			state.Resources[i] = EscalationResourceModel{
				Id:   types.StringValue(r.Id),
				Type: types.StringValue(r.Type),
			}
		}
	}

	if len(resp.Escalations) > 0 {
		state.Escalations = make([]EscalationLevelModel, len(resp.Escalations))
		for i, e := range resp.Escalations {
			level := EscalationLevelModel{
				WaitMins: types.Int64Value(int64(e.WaitMins)),
				Action:   types.StringValue(e.Action),
			}
			if e.Priority != "" {
				level.Priority = types.StringValue(e.Priority)
			}
			if e.RepeatFrequency != 0 {
				level.RepeatFrequency = types.Int64Value(int64(e.RepeatFrequency))
			}
			if e.NotifyLimitCount != 0 {
				level.NotifyLimitCount = types.Int64Value(int64(e.NotifyLimitCount))
			}
			if e.NotificationType != "" {
				level.NotificationType = types.StringValue(e.NotificationType)
			}
			if e.NotificationTemplateId != "" {
				level.NotificationTemplateId = types.StringValue(e.NotificationTemplateId)
			}
			for _, rec := range e.Recipients {
				level.Recipients = append(level.Recipients, EscalationRecipientModel{
					Id:   types.StringValue(rec.Id),
					Type: types.StringValue(rec.Type),
				})
			}
			if e.Incident != nil {
				inc := &EscalationIncidentModel{
					Priority:    types.StringValue(e.Incident.Priority),
					Subject:     types.StringValue(e.Incident.Subject),
					Description: types.StringValue(e.Incident.Description),
					Cc:          types.StringValue(e.Incident.Cc),
				}
				if e.Incident.AssigneeGroup != nil {
					inc.AssigneeGroupId = types.StringValue(e.Incident.AssigneeGroup.UniqueId)
				}
				if e.Incident.AssignedUser != nil {
					inc.AssignedUserId = types.StringValue(e.Incident.AssignedUser.Id)
				}
				if e.Incident.Category != nil {
					inc.CategoryId = types.StringValue(e.Incident.Category.UniqueId)
				}
				if e.Incident.SubCategory != nil {
					inc.SubCategoryId = types.StringValue(e.Incident.SubCategory.UniqueId)
				}
				if e.Incident.BusinessImpact != nil {
					inc.BusinessImpactId = types.StringValue(e.Incident.BusinessImpact.UniqueId)
				}
				if e.Incident.Urgency != nil {
					inc.UrgencyId = types.StringValue(e.Incident.Urgency.UniqueId)
				}
				for _, kid := range e.Incident.KnowledgeArticleIds {
					inc.KnowledgeArticleIds = append(inc.KnowledgeArticleIds, types.StringValue(kid))
				}
				// Map toMail:
				// - Users: prefer the structured ToMail.Users object; fall back to ToMailUserIds only when absent.
				// - UserGroups/Rosters: always read from the comma-separated fields (the structured object
				//   carries these as []interface{} which is unreliable; the comma-separated strings are canonical).
				if e.Incident.ToMail != nil || e.Incident.ToMailUserIds != "" || e.Incident.ToMailUserGroupIds != "" || e.Incident.ToMailRosterIds != "" {
					tm := &EscalationToMailModel{
						UsersIds:      make([]types.String, 0),
						UserGroupsIds: make([]types.String, 0),
						RostersIds:    make([]types.String, 0),
					}
					// Users — structured object is authoritative; only fall back to CSV if absent.
					if e.Incident.ToMail != nil && len(e.Incident.ToMail.Users) > 0 {
						for _, u := range e.Incident.ToMail.Users {
							tm.UsersIds = append(tm.UsersIds, types.StringValue(u.Id))
						}
					} else if e.Incident.ToMailUserIds != "" {
						for _, uid := range strings.Split(e.Incident.ToMailUserIds, ",") {
							uid = strings.TrimSpace(uid)
							if uid != "" {
								tm.UsersIds = append(tm.UsersIds, types.StringValue(uid))
							}
						}
					}
					// Groups — always from CSV.
					if e.Incident.ToMailUserGroupIds != "" {
						for _, gid := range strings.Split(e.Incident.ToMailUserGroupIds, ",") {
							gid = strings.TrimSpace(gid)
							if gid != "" {
								tm.UserGroupsIds = append(tm.UserGroupsIds, types.StringValue(gid))
							}
						}
					}
					// Rosters — always from CSV.
					if e.Incident.ToMailRosterIds != "" {
						for _, rid := range strings.Split(e.Incident.ToMailRosterIds, ",") {
							rid = strings.TrimSpace(rid)
							if rid != "" {
								tm.RostersIds = append(tm.RostersIds, types.StringValue(rid))
							}
						}
					}
					inc.ToMail = tm
				}
				level.Incident = inc
			}
			if e.UpdateIncident != nil {
				ui := &EscalationUpdateIncidentModel{
					UpdateWhenAlertStateChange:         types.BoolValue(e.UpdateIncident.UpdateWhenAlertStateChange),
					UpdateForEveryRepeatAlert:          types.BoolValue(e.UpdateIncident.UpdateForEveryRepeatAlert),
					UpdateWithRuleWhenAlertStateChange: types.BoolValue(e.UpdateIncident.UpdateWithRuleWhenAlertStateChange),
					UpdateWithRuleForEveryRepeatAlert:  types.BoolValue(e.UpdateIncident.UpdateWithRuleForEveryRepeatAlert),
					UpdateIncidentSubject:              types.BoolValue(e.UpdateIncident.UpdateIncidentSubject),
					UpdateIncidentSubjectWithRule:      types.BoolValue(e.UpdateIncident.UpdateIncidentSubjectWithRule),
					AutoResolveIncident:                types.BoolValue(e.UpdateIncident.AutoResolveIncident),
					AutoResolveUnassignedIncident:      types.BoolValue(e.UpdateIncident.AutoResolveUnassignedIncident),
					AutoHealWaitTime:                   types.Int64Value(int64(e.UpdateIncident.AutoHealWaitTime)),
					UpdatePriorityByMLConfiguration:    types.BoolValue(e.UpdateIncident.UpdatePriorityByMLConfiguration),
				}
				for _, pr := range e.UpdateIncident.PriorityRules {
					rule := EscalationPriorityRuleModel{
						Key:      types.StringValue(pr.Key),
						Operator: types.StringValue(pr.Operator),
						Value:    types.StringValue(pr.Value),
						Priority: types.StringValue(pr.Priority),
					}
					if pr.BusinessImpact != nil {
						rule.BusinessImpactId = types.StringValue(pr.BusinessImpact.UniqueId)
					}
					if pr.Urgency != nil {
						rule.UrgencyId = types.StringValue(pr.Urgency.UniqueId)
					}
					ui.PriorityRules = append(ui.PriorityRules, rule)
				}
				level.UpdateIncident = ui
			}
			state.Escalations[i] = level
		}
	}

	if resp.FilterCriteria != nil {
		state.FilterCriteria = &EscalationFilterCriteriaModel{
			MatchingType:        types.StringValue(resp.FilterCriteria.MatchingType),
			SearchQuery:         types.StringValue(resp.FilterCriteria.SearchQuery),
			ResourceSearchQuery: types.StringValue(resp.FilterCriteria.ResourceSearchQuery),
		}
	}
}

func normalizeEscalationsForState(escalations []EscalationLevelModel) []EscalationLevelModel {
	for i := range escalations {
		if escalations[i].Incident != nil && escalations[i].Incident.ToMail != nil {
			// Filter out null/empty values, but preserve empty collections as empty slices
			// (not nil) to match Terraform's empty set semantics in the plan.
			filtered := make([]types.String, 0, len(escalations[i].Incident.ToMail.RostersIds))
			for _, rid := range escalations[i].Incident.ToMail.RostersIds {
				if !rid.IsNull() && !rid.IsUnknown() && strings.TrimSpace(rid.ValueString()) != "" {
					filtered = append(filtered, rid)
				}
			}
			// Keep as empty slice, never nil, so state matches the planned empty set
			escalations[i].Incident.ToMail.RostersIds = filtered
		}

		if escalations[i].UpdateIncident != nil {
			for j := range escalations[i].UpdateIncident.PriorityRules {
				rule := &escalations[i].UpdateIncident.PriorityRules[j]

				// Empty strings are semantically null for optional fields in rules.
				if !rule.BusinessImpactId.IsUnknown() && !rule.BusinessImpactId.IsNull() && strings.TrimSpace(rule.BusinessImpactId.ValueString()) == "" {
					rule.BusinessImpactId = types.StringNull()
				}
				if !rule.UrgencyId.IsUnknown() && !rule.UrgencyId.IsNull() && strings.TrimSpace(rule.UrgencyId.ValueString()) == "" {
					rule.UrgencyId = types.StringNull()
				}
				if !rule.Priority.IsUnknown() && !rule.Priority.IsNull() && strings.TrimSpace(rule.Priority.ValueString()) == "" {
					rule.Priority = types.StringNull()
				}
			}
		}
	}

	return escalations
}

// Create handles the creation of the resource.
func (r *AlertEscalationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AlertEscalationPolicyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plannedEscalations := normalizeEscalationsForState(plan.Escalations)

	tenantId := r.apiClient.TenantId
	if !plan.Client.IsNull() && plan.Client.ValueString() != "" {
		tenantId = plan.Client.ValueString()
	}

	policy := buildAlertEscalationPolicyRequest(plan)

	created, err := r.apiClient.CreateAlertEscalationPolicy(tenantId, policy)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", err.Error())
		return
	}

	mapAlertEscalationPolicyToState(created, &plan)
	plan.Escalations = plannedEscalations

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the resource state.
func (r *AlertEscalationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AlertEscalationPolicyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	existing, err := r.apiClient.GetAlertEscalationPolicy(tenantId, state.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	// Preserve escalations from the existing state rather than overwriting with API values.
	// The API normalizes/transforms escalation fields in ways that don't round-trip cleanly
	// (e.g. booleans defaulted to false, empty collections omitted, etc.), which would cause
	// a perpetual diff on every plan. Non-escalation fields (id, name, description, etc.)
	// are safe to refresh from the API response.
	savedEscalations := state.Escalations
	mapAlertEscalationPolicyToState(existing, &state)
	state.Escalations = normalizeEscalationsForState(savedEscalations)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update handles updating the resource.
func (r *AlertEscalationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state AlertEscalationPolicyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plannedEscalations := normalizeEscalationsForState(plan.Escalations)

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	policy := buildAlertEscalationPolicyRequest(plan)

	updated, err := r.apiClient.UpdateAlertEscalationPolicy(tenantId, state.Id.ValueString(), policy)
	if err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())
		return
	}

	mapAlertEscalationPolicyToState(updated, &plan)
	plan.Escalations = plannedEscalations

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete handles deleting the resource.
func (r *AlertEscalationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AlertEscalationPolicyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantId := r.apiClient.TenantId
	if !state.Client.IsNull() && state.Client.ValueString() != "" {
		tenantId = state.Client.ValueString()
	}

	err := r.apiClient.DeleteAlertEscalationPolicy(tenantId, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}
