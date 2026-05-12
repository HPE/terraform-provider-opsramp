// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
	"time"
)

// CreateServicemap - Create new Servicemap (unofficial)
// api: POST v3 /api/v3/tenants/{tenantId}/serviceGroup
func (c *OpsRampClient) CreateServicemap(tenantId string, servicemap CreateServicemap, retries_optional ...int) (*CreateServicemap, error) {

	retries := 3
	if len(retries_optional) > 0 {
		retries = retries_optional[0]
	}

	if retries < 0 {
		return nil, fmt.Errorf("could not create servicemap, retries reched")
	}

	// Save and remove child services
	childs := servicemap.Services
	resources := servicemap.Resources
	servicemap.Services = make([]CreateServicemap, 0)
	servicemap.Resources = nil

	// Assert servicemap parameters
	if servicemap.Type == "service" && len(servicemap.Resources) > 0 {
		return nil, fmt.Errorf("servicemap node can't have 'type=Service' and 'resources' at the same time")
	}

	// Add required parameters in case type = "resource"

	if servicemap.ServiceAvailabilityMonitor == nil {
		servicemap.ServiceAvailabilityMonitor = &ServiceAvailabilityMonitor{
			AlertType: "availability",
			Metrics:   make([]string, 0),
		}
	}

	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(servicemap)
	if err != nil {
		return nil, err
	}

	// Prepare the URL, Method and Payload fo the Client
	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/serviceGroup", c.BaseUrl, tenantId)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var createdServicemap CreateServicemap
	err = json.Unmarshal([]byte(body), &createdServicemap)
	if err != nil {
		return nil, err
	}

	// Create child resources
	if len(resources) > 0 {

		err = c.CreateServicemapResources(tenantId, createdServicemap.Id, resources)
		if err != nil {
			return nil, err
		}
	}

	// Recursive call to create children
	for _, child := range childs {

		// Adjust parent
		child.Parent = &Parent{
			Id: createdServicemap.Id,
		}

		createdChild, err := c.CreateServicemap(tenantId, child)
		if err != nil {
			return nil, err
		}

		createdServicemap.Services = append(createdServicemap.Services, *createdChild)
	}

	// Restore resources
	servicemap.Resources = resources

	// Return ID of the record created
	return &createdServicemap, nil
}

func (c *OpsRampClient) CreateServicemapResources(tenantId string, parentId string, resources []string) error {
	var data []struct {
		Id string `json:"id"`
	}

	// Convert Request Data/Body to JSON
	for _, resource := range resources {
		res := struct {
			Id string `json:"id"`
		}{
			Id: resource,
		}
		data = append(data, res)
	}

	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Prepare the URL, Method and Payload fo the Client
	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/serviceGroups/%s/childs", c.BaseUrl, tenantId, parentId)
	method := "POST"

	// Create a new Request
	_, err = c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	return nil
}

// GetServicemap - Get existing Servicemap
func (c *OpsRampClient) GetServicemap(tenantId string, servicemapId string) (*CreateServicemap, error) {

	// Prepare the URL, Method and Payload fo the Client
	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/service-maps/%s", c.BaseUrl, tenantId, servicemapId)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var responseBody ReadServicemap
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	translatedResponse := TranslateGetToCreateServicemapModels(responseBody)

	// Return ID of the record created
	return &translatedResponse, nil

}

// DeleteServicemap - DeleteServicemap
func (c *OpsRampClient) DeleteServicemap(tenantId string, servicemapId string) error {

	// Prepare the URL, Method and Payload fo the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/serviceGroups/%s", c.BaseUrl, tenantId, servicemapId)
	method := "DELETE"

	// Create a new Request
	_, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return err
	}

	return nil
}
