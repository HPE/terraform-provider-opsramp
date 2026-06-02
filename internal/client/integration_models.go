// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

// Integration models for the generic integration resource

// InstallIntegrationRequest represents the request body to install an integration.
// Fields are used selectively depending on the application type.
type InstallIntegrationRequest struct {
	DisplayName   string `json:"displayName,omitempty"`
	Description   string `json:"description,omitempty"`
	Category      string `json:"category,omitempty"`
	IPAddress     string `json:"ipAddress,omitempty"`
	CredentialSet string `json:"credentialSet,omitempty"`

	// For event-based integrations (CUSTOM-EVENT, NEWRELIC, etc.)
	AlertSource *AlertSource `json:"alertSource,omitempty"`

	// For configuration-based integrations (VMWARE, etc.)
	ProviderProps             map[string]any     `json:"providerProps,omitempty"`
	ConfigDetails             map[string]any     `json:"configDetails,omitempty"`
	MultiAppsDiscoveryEnabled bool               `json:"multiAppsDiscoveryEnabled,omitempty"`
	DiscoveryProfiles         []DiscoveryProfile `json:"discoveryProfiles,omitempty"`
}

// AlertSource identifies the third-party alert source for event integrations
type AlertSource struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	TechUID     string `json:"techUid,omitempty"`
}

// DiscoveryProfile represents a discovery scan profile for configuration integrations
type DiscoveryProfile struct {
	ID               int                `json:"id,omitempty"`
	MgmtProfileUUID  string             `json:"mgmtProfileUuid,omitempty"`
	CustomAttributes map[string]string  `json:"customAttributes,omitempty"`
	ScanNow          bool               `json:"scanNow,omitempty"`
	Schedule         *DiscoverySchedule `json:"schedule,omitempty"`
	Policy           *DiscoveryPolicy   `json:"policy,omitempty"`
}

// DiscoverySchedule defines when discovery scans run
type DiscoverySchedule struct {
	PatternType string `json:"patternType,omitempty"`
	Pattern     string `json:"pattern,omitempty"`
	StartTime   string `json:"startTime,omitempty"`
}

// DiscoveryPolicy defines which resources to manage from discovery
type DiscoveryPolicy struct {
	Name       string            `json:"name,omitempty"`
	EntityType string            `json:"entityType,omitempty"`
	MatchType  string            `json:"matchType,omitempty"`
	Rules      []DiscoveryRule   `json:"rules,omitempty"`
	Actions    []DiscoveryAction `json:"actions,omitempty"`
}

// DiscoveryRule is a filter rule within a discovery policy
type DiscoveryRule struct {
	FilterType   string   `json:"filterType,omitempty"`
	ResourceType []string `json:"resourceType,omitempty"`
}

// DiscoveryAction is an action taken on discovered resources
type DiscoveryAction struct {
	Action string `json:"action,omitempty"`
}

// IntegrationResponse represents the API response for an installed integration
type IntegrationResponse struct {
	ID            string                    `json:"id"`
	DisplayName   string                    `json:"displayName"`
	Description   string                    `json:"description,omitempty"`
	Integration   *IntegrationDetails       `json:"integration,omitempty"`
	InboundConfig *IntegrationInboundConfig `json:"inboundConfig,omitempty"`
	Category      string                    `json:"category,omitempty"`
	InstalledBy   string                    `json:"installedBy,omitempty"`
	InstalledTime string                    `json:"installedTime,omitempty"`
	Status        string                    `json:"status,omitempty"`
	IPAddress     string                    `json:"ipAddress,omitempty"`
	CredentialSet string                    `json:"credentialSet,omitempty"`
	ConfigFiles   []ConfigFile              `json:"configFiles,omitempty"`
	AlertSource   *AlertSource              `json:"alertSource,omitempty"`

	// Configuration integrations
	ProviderProps             map[string]any     `json:"providerProps,omitempty"`
	ConfigDetails             map[string]any     `json:"configDetails,omitempty"`
	MultiAppsDiscoveryEnabled bool               `json:"multiAppsDiscoveryEnabled,omitempty"`
	DiscoveryProfiles         []DiscoveryProfile `json:"discoveryProfiles,omitempty"`
	AdditionalProperties      map[string]string  `json:"additionalProperties,omitempty"`
}

// IntegrationDetails describes the integration application metadata
type IntegrationDetails struct {
	ID                    string `json:"id,omitempty"`
	Name                  string `json:"name,omitempty"`
	MultipleInstallations bool   `json:"multipleInstallations,omitempty"`
	IntegrationId         int    `json:"integrationId,omitempty"`
	SubCategory           string `json:"subCategory,omitempty"`
}

// ConfigFile represents a downloadable config file from an integration
type ConfigFile struct {
	ID                     int    `json:"id,omitempty"`
	Name                   string `json:"name,omitempty"`
	ContentURL             string `json:"contentURL,omitempty"`
	AllowSpecialCharacters bool   `json:"allowSpecialCharacters,omitempty"`
}

// IntegrationInboundConfig holds the inbound configuration state
type IntegrationInboundConfig struct {
	Authentication *IntegrationAuthentication `json:"authentication,omitempty"`
}

// IntegrationAuthentication represents the inbound authentication configuration
type IntegrationAuthentication struct {
	AuthType   string         `json:"authType,omitempty"`
	Token      string         `json:"token,omitempty"`
	WebhookURL string         `json:"webhookUrl,omitempty"`
	Role       *RoleClientRef `json:"role,omitempty"`
}

// SetInboundAuthRequest is the request to configure inbound authentication
type SetInboundAuthRequest struct {
	AuthType string         `json:"authType"`
	Role     *RoleClientRef `json:"role,omitempty"`
}

// InboundAuthResponse is the response from setting/getting inbound auth
type InboundAuthResponse struct {
	AuthType   string         `json:"authType"`
	Token      string         `json:"token,omitempty"`
	WebhookURL string         `json:"webhookUrl,omitempty"`
	Role       *RoleClientRef `json:"role,omitempty"`
}

// MappingAttributesRequest represents the request to configure attribute mapping
type MappingAttributesRequest struct {
	InboundConfig *MappingInboundConfig `json:"inboundConfig,omitempty"`
}

// MappingInboundConfig wraps the list of attribute mappings
type MappingInboundConfig struct {
	MapAttributes []MapAttribute `json:"mapAttributes,omitempty"`
}

// MapAttribute represents a single attribute mapping rule
type MapAttribute struct {
	EntityType           string             `json:"entityType,omitempty"`
	ThirdPartyEntityType string             `json:"thirdPartyEntityType,omitempty"`
	Name                 string             `json:"name"`
	ThirdPartyAttrName   string             `json:"thirdPartyAttrName,omitempty"`
	AttrName             string             `json:"attrName,omitempty"`
	AttrValues           []AttrValueMapping `json:"attrValues"`
	ParsingProperty      *ParsingProperty   `json:"parsingProperty"`
}

// AttrValueMapping maps third-party values to OpsRamp values
type AttrValueMapping struct {
	AttrValue           string `json:"attrValue,omitempty"`
	ThirdPartyAttrValue string `json:"thirdPartyAttrValue,omitempty"`
}

// ParsingProperty defines how to extract values from inbound data
type ParsingProperty struct {
	DefaultValue  string            `json:"defaultValue"`
	OprSet        []ParsingOperator `json:"oprset"`
	ValueMappings []any             `json:"valueMappings"`
}

// ParsingOperator defines a string parsing operation
type ParsingOperator struct {
	Operator  string `json:"operator,omitempty"`
	StartWord string `json:"startWord,omitempty"`
	EndWord   string `json:"endWord,omitempty"`
	RegexStr  string `json:"regexStr,omitempty"`
}

// NotifierRequest represents the outbound notifier configuration
type NotifierRequest struct {
	Type                string         `json:"type"`
	BaseURI             string         `json:"baseURI"`
	AuthType            string         `json:"authType"`
	GrantType           string         `json:"grantType,omitempty"`
	UserName            string         `json:"userName,omitempty"`
	Password            string         `json:"password,omitempty"`
	APIKey              string         `json:"apiKey,omitempty"`
	APISecret           string         `json:"apiSecret,omitempty"`
	AccessTokenURL      string         `json:"accessTokenURL,omitempty"`
	Scope               string         `json:"scope,omitempty"`
	TokenURL            string         `json:"tokenURL,omitempty"`
	TokenPayload        map[string]any `json:"tokenPayload,omitempty"`
	TokenHeaders        []KeyValuePair `json:"tokenHeaders,omitempty"`
	TokensPath          []KeyValuePair `json:"tokensPath,omitempty"`
	ResourceAuthHeaders []KeyValuePair `json:"resourceAuthHeaders,omitempty"`
}

// KeyValuePair represents a generic key-value configuration entry
type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// DeleteIntegrationRequest represents the request body to uninstall an integration
type DeleteIntegrationRequest struct {
	UninstallReason string `json:"uninstallReason,omitempty"`
}

// InboundEntityProperty represents an available property for inbound attribute mapping
type InboundEntityProperty struct {
	Entity       string `json:"entity"`
	Property     string `json:"property"`
	PropertyType string `json:"propertyType"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	DisplayGroup string `json:"displayGroup"`
	Mapable      bool   `json:"mapable"`
	ValueMapable bool   `json:"valueMapable"`
	Mapped       bool   `json:"mapped"`
	Mandatory    bool   `json:"mandatory"`
	Parsable     bool   `json:"parsable"`
}
