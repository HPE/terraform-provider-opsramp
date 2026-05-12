// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

// Script models

// Script represents an RBA script for creation, update, and retrieval
type Script struct {
	Id             int                `json:"id,omitempty"`
	Name           string             `json:"name"`
	Description    string             `json:"description,omitempty"`
	Category       *ScriptCategoryRef `json:"category,omitempty"`
	Platforms      []string           `json:"platforms,omitempty"`
	Parameters     []ScriptParameter  `json:"parameters,omitempty"`
	ExecutionType  string             `json:"executionType,omitempty"`
	InstallTimeout int                `json:"installTimeout,omitempty"`
	Attachment     *ScriptAttachment  `json:"attachment,omitempty"`
	ScriptVersion  string             `json:"scriptVersion,omitempty"`
	RegistryPath   string             `json:"registryPath,omitempty"`
	RegistryValue  string             `json:"registryValue,omitempty"`
	ProcessName    string             `json:"processName,omitempty"`
	ServiceName    string             `json:"serviceName,omitempty"`
}

// ScriptCategoryRef represents a category reference inside a script
type ScriptCategoryRef struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// ScriptParameter represents a single parameter of a script
type ScriptParameter struct {
	Id           int    `json:"id,omitempty"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	DefaultValue string `json:"defaultValue,omitempty"`
	Type         string `json:"type,omitempty"`
	DataType     string `json:"dataType,omitempty"`
}

// ScriptAttachment represents the file attachment of a script
type ScriptAttachment struct {
	Id         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	ContentURL string `json:"contentURL,omitempty"`
	File       string `json:"file,omitempty"`
}
