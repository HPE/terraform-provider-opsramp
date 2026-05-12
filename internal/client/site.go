// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
)

// CreateSite creates a new site
func (c *OpsRampClient) CreateSite(tenantId string, site Site) (*SiteResponse, error) {
	rb, err := json.Marshal(site)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/sites", c.BaseUrl, tenantId)
	method := "POST"

	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody SiteResponse
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// GetSite retrieves a specific site by ID
func (c *OpsRampClient) GetSite(tenantId string, siteId string) (*SiteResponse, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/sites/%s", c.BaseUrl, tenantId, siteId)
	method := "GET"

	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var responseBody SiteResponse
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// UpdateSite updates an existing site
func (c *OpsRampClient) UpdateSite(tenantId string, siteId string, site Site) (*SiteResponse, error) {
	rb, err := json.Marshal(site)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/sites/%s", c.BaseUrl, tenantId, siteId)
	method := "POST"

	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody SiteResponse
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// DeleteSite deletes a site
func (c *OpsRampClient) DeleteSite(tenantId string, siteId string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/sites/%s", c.BaseUrl, tenantId, siteId)
	method := "DELETE"

	_, err := c.NewJsonRequest(method, apiUrl, nil)
	return err
}

// GetSites retrieves all sites (minimal info)
func (c *OpsRampClient) GetSites(tenantId string) ([]SiteMinimal, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/sites/minimal", c.BaseUrl, tenantId)
	method := "GET"

	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var responseBody []SiteMinimal
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

// SearchSites searches for sites using a query
func (c *OpsRampClient) SearchSites(tenantId string, query string) ([]SiteResponse, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/sites/search?q=%s", c.BaseUrl, tenantId, query)
	method := "GET"

	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var responseBody []SiteResponse
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
