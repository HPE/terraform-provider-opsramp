// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
)

// ListScriptCategories returns all RBA categories (tree structure) for a tenant.
func (c *OpsRampClient) ListScriptCategories(tenantId string) ([]ScriptCategory, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/rba/categories", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var responseBody []ScriptCategory
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

// findScriptCategoryInTree recursively searches a category tree for a given ID.
// Returns the matched category and the ID of its parent (0 if root-level).
func findScriptCategoryInTree(categories []ScriptCategory, id int, parentId int) (*ScriptCategory, int) {
	for i := range categories {
		if categories[i].Id == id {
			return &categories[i], parentId
		}
		if found, pid := findScriptCategoryInTree(categories[i].Childs, id, categories[i].Id); found != nil {
			return found, pid
		}
	}
	return nil, 0
}

// ScriptCategoryResult wraps a ScriptCategory with its resolved parent ID.
type ScriptCategoryResult struct {
	Category *ScriptCategory
	ParentId int // 0 if the category is at the root level
}

// GetScriptCategory retrieves a specific RBA category by ID by searching the full list.
func (c *OpsRampClient) GetScriptCategory(tenantId string, categoryId int) (*ScriptCategoryResult, error) {
	categories, err := c.ListScriptCategories(tenantId)
	if err != nil {
		return nil, err
	}

	cat, parentId := findScriptCategoryInTree(categories, categoryId, 0)
	if cat == nil {
		return nil, fmt.Errorf("script category with id %d not found", categoryId)
	}

	return &ScriptCategoryResult{Category: cat, ParentId: parentId}, nil
}

// CreateScriptCategory creates a new RBA category.
func (c *OpsRampClient) CreateScriptCategory(tenantId string, categoryData ScriptCategory) (*ScriptCategory, error) {
	rb, err := json.Marshal(categoryData)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/rba/categories", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody ScriptCategory
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// UpdateScriptCategory updates an existing RBA category.
// The category ID must be set in the categoryData struct; the PUT URL has no ID path segment.
func (c *OpsRampClient) UpdateScriptCategory(tenantId string, categoryData ScriptCategory) (*ScriptCategory, error) {
	rb, err := json.Marshal(categoryData)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/rba/categories", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("PUT", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody ScriptCategory
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// DeleteScriptCategory deletes an RBA category by ID.
func (c *OpsRampClient) DeleteScriptCategory(tenantId string, categoryId int) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/rba/categories/%d", c.BaseUrl, tenantId, categoryId)

	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)
	return err
}
