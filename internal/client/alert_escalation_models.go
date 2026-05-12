// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

// AlertEscalationPolicy represents an alert escalation policy
type AlertEscalationPolicy struct {
	Id             string                    `json:"id,omitempty"`
	Name           string                    `json:"name"`
	Description    string                    `json:"description"`
	TenantScope    string                    `json:"tenantScope,omitempty"`
	Precedence     int                       `json:"precedence,omitempty"`
	EscalationType string                    `json:"escalationType"`
	PolicyType     string                    `json:"policyType,omitempty"`
	EnabledMode    string                    `json:"enabledMode,omitempty"`
	Active         bool                      `json:"active,omitempty"`
	AllClients     bool                      `json:"allClients"`
	Scope          *EscalationScope          `json:"scope,omitempty"`
	Resources      []EscalationResource      `json:"resources"`
	Escalations    []EscalationLevel         `json:"escalations,omitempty"`
	FilterCriteria *EscalationFilterCriteria `json:"filterCriteria,omitempty"`
}

// EscalationScope defines the scope of the policy
type EscalationScope struct {
	Uuid string `json:"uuid,omitempty"`
}

// EscalationResource defines a resource in scope
type EscalationResource struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

// EscalationFilterCriteria defines filter criteria for the policy
type EscalationFilterCriteria struct {
	MatchingType        string `json:"matchingType,omitempty"`
	SearchQuery         string `json:"searchQuery,omitempty"`
	ResourceSearchQuery string `json:"resourceSearchQuery,omitempty"`
}

// EscalationLevel defines an escalation level/step
type EscalationLevel struct {
	WaitMins               int                       `json:"waitMins"`
	Action                 string                    `json:"action"`
	Priority               string                    `json:"priority,omitempty"`
	RepeatFrequency        int                       `json:"repeatFrequency,omitempty"`
	NotifyLimitCount       int                       `json:"notifyLimitCount,omitempty"`
	NotificationType       string                    `json:"notificationType,omitempty"`
	NotificationTemplateId string                    `json:"notificationTemplateId,omitempty"`
	Recipients             []EscalationRecipient     `json:"recipients,omitempty"`
	Incident               *EscalationIncident       `json:"incident,omitempty"`
	UpdateIncident         *EscalationUpdateIncident `json:"updateIncident,omitempty"`
}

// EscalationRecipient defines a notification recipient
type EscalationRecipient struct {
	Name string `json:"name,omitempty"`
	Id   string `json:"id"`
	Type string `json:"type"`
}

// EscalationIncident defines an incident to create
type EscalationIncident struct {
	Priority            string                 `json:"priority,omitempty"`
	Subject             string                 `json:"subject,omitempty"`
	Description         string                 `json:"description,omitempty"`
	AssigneeGroup       *EscalationUniqueRef   `json:"assigneeGroup,omitempty"`
	AssignedUser        *EscalationUserRef     `json:"assignedUser,omitempty"`
	Category            *EscalationUniqueRef   `json:"category,omitempty"`
	SubCategory         *EscalationUniqueRef   `json:"subCategory,omitempty"`
	BusinessImpact      *EscalationUniqueRef   `json:"businessImpact,omitempty"`
	Urgency             *EscalationUniqueRef   `json:"urgency,omitempty"`
	AttachedArticles    []EscalationArticleRef `json:"attachedArticles,omitempty"`
	KnowledgeArticleIds []string               `json:"knowledgeArticleIds,omitempty"`
	Cc                  string                 `json:"cc,omitempty"`
	ToMail              *EscalationToMail      `json:"toMail,omitempty"`
	ToMailUserIds       string                 `json:"toMailUserIds,omitempty"`
	ToMailUserGroupIds  string                 `json:"toMailUserGroupIds,omitempty"`
	ToMailRosterIds     string                 `json:"toMailRosterIds,omitempty"`
}

// EscalationUniqueRef is a reference by uniqueId
type EscalationUniqueRef struct {
	UniqueId string `json:"uniqueId"`
}

// EscalationUserRef is a reference to a user
type EscalationUserRef struct {
	Id        string `json:"id"`
	LoginName string `json:"loginName,omitempty"`
}

// EscalationArticleRef is a reference to a KB article
type EscalationArticleRef struct {
	Id      string `json:"id"`
	Subject string `json:"subject,omitempty"`
}

// EscalationToMail defines mail recipients
type EscalationToMail struct {
	Users      []EscalationUserRef `json:"users,omitempty"`
	UserGroups []interface{}       `json:"userGroups,omitempty"`
	Rosters    []interface{}       `json:"rosters,omitempty"`
}

// EscalationUpdateIncident defines incident update settings
type EscalationUpdateIncident struct {
	UpdateWhenAlertStateChange         bool                     `json:"updateWhenAlertStateChange,omitempty"`
	UpdateForEveryRepeatAlert          bool                     `json:"updateForEveryRepeatAlert,omitempty"`
	UpdateWithRuleWhenAlertStateChange bool                     `json:"updateWithRuleWhenAlertStateChange,omitempty"`
	UpdateWithRuleForEveryRepeatAlert  bool                     `json:"updateWithRuleForEveryRepeatAlert,omitempty"`
	UpdateIncidentSubject              bool                     `json:"updateIncidentSubject,omitempty"`
	UpdateIncidentSubjectWithRule      bool                     `json:"updateIncidentSubjectWithRule,omitempty"`
	AutoResolveIncident                bool                     `json:"autoResolveIncident,omitempty"`
	AutoResolveUnassignedIncident      bool                     `json:"autoResolveUnassignedIncident,omitempty"`
	AutoHealWaitTime                   int                      `json:"autoHealWaitTime,omitempty"`
	UpdatePriorityByMLConfiguration    bool                     `json:"updatePriorityByMLConfiguration,omitempty"`
	PriorityRules                      []EscalationPriorityRule `json:"priorityRules,omitempty"`
}

// EscalationPriorityRule defines a priority rule for incident updates
type EscalationPriorityRule struct {
	Key            string               `json:"key"`
	Operator       string               `json:"operator"`
	Value          string               `json:"value"`
	BusinessImpact *EscalationUniqueRef `json:"businessImpact,omitempty"`
	Urgency        *EscalationUniqueRef `json:"urgency,omitempty"`
	Priority       string               `json:"priority,omitempty"`
}
