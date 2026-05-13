// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// OpsRampClient handles API communication.
type OpsRampClient struct {
	BaseUrl       string
	TenantId      string
	AccessToken   string
	muAccessToken sync.RWMutex
	Client        *http.Client
}

// OAuthTokenResponse stores API response
type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// NewOpsRampClient creates a new client and retrieves the OAuth token.
func NewOpsRampClient(clientID string, clientSecret string, endpoint string, tenant string) (*OpsRampClient, error) {
	rateLimitedTransport := NewRateLimitedTransport(&http.Transport{})

	client := &http.Client{
		Transport: rateLimitedTransport,
	}

	baseUrl := fmt.Sprintf("https://%s", endpoint)
	resource := "/tenancy/auth/oauth/token"

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "client_credentials")

	u, _ := url.ParseRequestURI(baseUrl)

	u.Path = resource
	urlStr := u.String()

	resp, err := client.Post(urlStr, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to retrieve OAuth token")
	}

	var tokenResponse OAuthTokenResponse
	json.NewDecoder(resp.Body).Decode(&tokenResponse)

	return &OpsRampClient{BaseUrl: baseUrl, TenantId: tenant, AccessToken: tokenResponse.AccessToken, Client: client}, nil
}

// Thread-safe getter for AccessToken
func (c *OpsRampClient) GetAccessToken() string {
	c.muAccessToken.RLock()
	defer c.muAccessToken.RUnlock()
	return c.AccessToken
}

// Thread-safe setter for AccessToken
func (c *OpsRampClient) SetAccessToken(token string) {
	c.muAccessToken.Lock()
	defer c.muAccessToken.Unlock()
	c.AccessToken = token
}

func (c *OpsRampClient) NewJsonRequest(method string, apiUrl string, payload []byte) ([]byte, error) {
	validStatuses := map[int]struct{}{
		http.StatusOK:       {},
		http.StatusAccepted: {},
		http.StatusCreated:  {},
	}

	const maxRetries = 3

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Build a fresh request for each attempt.
		payloadReader := strings.NewReader(string(payload))
		req, err := http.NewRequest(method, apiUrl, payloadReader)
		if err != nil {
			return nil, err
		}

		// Add the proper Headers
		req.Header.Set("Authorization", "Bearer "+c.GetAccessToken())
		req.Header.Add("Content-Type", "application/json")

		// Make the Request
		res, err := c.Client.Do(req)
		if err != nil {
			return nil, err
		}

		body, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return nil, readErr
		}

		if _, ok := validStatuses[res.StatusCode]; ok {
			return body, nil
		}

		// Retries are needed in some cases where previous resources
		// ex: parents, groups, etc.. are not fully ready and API returns 500 with code 0005.
		// This is a workaround to mitigate that issue.
		if attempt < maxRetries && shouldRetryOpsRampServerError(res.StatusCode, body) {
			delay := time.Duration(attempt+1) * time.Second
			time.Sleep(delay)
			continue
		}

		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return nil, fmt.Errorf("request failed after retries")
}

func shouldRetryOpsRampServerError(statusCode int, body []byte) bool {
	if statusCode != http.StatusInternalServerError {
		return false
	}

	var apiError struct {
		Code string `json:"code"`
	}

	if err := json.Unmarshal(body, &apiError); err == nil {
		return apiError.Code == "0005"
	}

	// Fallback for non-standard error payloads.
	return strings.Contains(string(body), `"code":"0005"`)
}
