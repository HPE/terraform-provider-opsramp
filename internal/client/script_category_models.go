// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

// ScriptCategory models

// ScriptCategory represents an RBA category for creation, update, and retrieval.
// The API returns a tree structure where each category may have child categories.
type ScriptCategory struct {
	Id     int                      `json:"id,omitempty"`
	Name   string                   `json:"name"`
	Parent *ScriptCategoryParentRef `json:"parent,omitempty"`
	Childs []ScriptCategory         `json:"childs,omitempty"`
}

// ScriptCategoryParentRef is the nested parent object used in create/update requests.
type ScriptCategoryParentRef struct {
	Id int `json:"id"`
}
