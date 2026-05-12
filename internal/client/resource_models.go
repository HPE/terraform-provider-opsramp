// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

type DeviceType struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

type CreateResource struct {
	ResourceName string `json:"resourceName"`
	ResourceType string `json:"resourceType"`
	HostName     string `json:"hostName"`
}

type UpdateResource struct {
	ResourceType string `json:"resourceType"`
}

type ResourceCreated struct {
	ResourceUuid string `json:"resourceUUID"`
	TenantId     string `json:"tenantID"`
}

type GeneralInfo struct {
	HostName     string `json:"hostName"`
	ResourceName string `json:"resourceName"`
	ResourceType string `json:"resourceType"`
}

type GetResource struct {
	Uuid           string      `json:"id"`
	AgentInstalled bool        `json:"agentInstalled"`
	GeneralInfo    GeneralInfo `json:"generalInfo"`
}

type SearchResources struct {
	Results []SearchResource `json:"results"`
}

type SearchResource struct {
	Uuid         string `json:"id"`
	ResourceName string `json:"resourceName"`
	ResourceType string `json:"resourceType"`
}
