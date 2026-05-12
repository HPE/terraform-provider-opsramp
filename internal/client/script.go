// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
)

// CreateScript creates a new RBA script inside a category
func (c *OpsRampClient) CreateScript(tenantId string, categoryId string, scriptData Script) (*Script, error) {
	script, err := json.Marshal(scriptData)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/rba/categories/%s/scripts", c.BaseUrl, tenantId, categoryId)

	body, err := c.NewJsonRequest("POST", apiUrl, script)
	if err != nil {
		return nil, err
	}

	var responseBody Script
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// GetScript retrieves a specific RBA script by ID
func (c *OpsRampClient) GetScript(tenantId string, categoryId string, scriptId string) (*Script, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/rba/categories/%s/scripts/%s", c.BaseUrl, tenantId, categoryId, scriptId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var responseBody Script
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// UpdateScript updates an existing RBA script
func (c *OpsRampClient) UpdateScript(tenantId string, categoryId string, scriptId string, scriptData Script) (*Script, error) {
	rb, err := json.Marshal(scriptData)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/rba/categories/%s/scripts/%s", c.BaseUrl, tenantId, categoryId, scriptId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody Script
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// DeleteScript deletes an RBA script
func (c *OpsRampClient) DeleteScript(tenantId string, categoryId string, scriptId string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/rba/categories/%s/scripts/%s", c.BaseUrl, tenantId, categoryId, scriptId)

	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)
	return err
}
