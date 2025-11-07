package akeneo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ClientConfig contains the configuration for the Akeneo client
type ClientConfig struct {
	Host     string
	ClientID string
	Secret   string
	Username string
	Password string
}

// Client represents a client for the Akeneo API
type Client struct {
	config      ClientConfig
	httpClient  *http.Client
	accessToken string
	tokenExpiry time.Time
}

// TokenResponse represents the authentication endpoint response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// ReferenceEntityRecord represents a Reference Entity record
type ReferenceEntityRecord map[string]interface{}

// AkeneoErrorResponse represents an Akeneo error response
type AkeneoErrorResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Errors  []AkeneoFieldError `json:"errors,omitempty"`
}

// AkeneoFieldError represents a field-specific error
type AkeneoFieldError struct {
	Property string `json:"property"`
	Message  string `json:"message"`
}

// NewClient creates a new Akeneo client
func NewClient(config ClientConfig) (*Client, error) {
	client := &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Get access token
	if err := client.authenticate(); err != nil {
		return nil, fmt.Errorf("authentication error: %w", err)
	}

	return client, nil
}

// authenticate obtains an OAuth2 access token
func (c *Client) authenticate() error {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", c.config.Username)
	data.Set("password", c.config.Password)

	req, err := http.NewRequest("POST", c.config.Host+"/api/oauth/v1/token", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.config.ClientID, c.config.Secret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("authentication error: %d - %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return err
	}

	c.accessToken = tokenResp.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}

// ensureValidToken verifies the token is valid and renews it if necessary
func (c *Client) ensureValidToken() error {
	if time.Now().After(c.tokenExpiry.Add(-5 * time.Minute)) {
		return c.authenticate()
	}
	return nil
}

// GetReferenceEntityRecords retrieves all records from a Reference Entity
func (c *Client) GetReferenceEntityRecords(entityName string) ([]ReferenceEntityRecord, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}

	var allRecords []ReferenceEntityRecord
	page := 1
	limit := 100

	for {
		url := fmt.Sprintf("%s/api/rest/v1/reference-entities/%s/records?page=%d&limit=%d", 
			c.config.Host, entityName, page, limit)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+c.accessToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("error fetching records: %d - %s", resp.StatusCode, string(body))
		}

		var response struct {
			Embedded struct {
				Items []ReferenceEntityRecord `json:"items"`
			} `json:"_embedded"`
			Links struct {
				Next *struct {
					Href string `json:"href"`
				} `json:"next"`
			} `json:"_links"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, err
		}

		allRecords = append(allRecords, response.Embedded.Items...)

		// If no next page, finish
		if response.Links.Next == nil {
			break
		}

		page++
	}

	return allRecords, nil
}

// PatchReferenceEntityRecord creates or updates a record in a Reference Entity
func (c *Client) PatchReferenceEntityRecord(entityName, code string, record ReferenceEntityRecord) error {
	if err := c.ensureValidToken(); err != nil {
		return err
	}

	// Clean fields that should not be sent
	cleanRecord := c.cleanRecord(record)

	jsonData, err := json.Marshal(cleanRecord)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/rest/v1/reference-entities/%s/records/%s", 
		c.config.Host, entityName, code)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		
		// For 422 errors, try to parse Akeneo error response
		if resp.StatusCode == http.StatusUnprocessableEntity {
			var errorResponse AkeneoErrorResponse
			if parseErr := json.Unmarshal(body, &errorResponse); parseErr == nil {
				return fmt.Errorf("validation error in record %s: %s", code, c.formatAkeneoErrors(errorResponse))
			}
		}
		
		return fmt.Errorf("error updating record %s: %d - %s", code, resp.StatusCode, string(body))
	}

	return nil
}

// cleanRecord removes fields that should not be sent in write operations
func (c *Client) cleanRecord(record ReferenceEntityRecord) ReferenceEntityRecord {
	cleaned := make(ReferenceEntityRecord)
	
	// List of fields to exclude
	excludedFields := map[string]bool{
		"_links":                true,
		"created":               true,
		"updated":               true,
		"reference_entity_code": true, // Metadata field that causes 422 error
	}
	
	for key, value := range record {
		// Exclude metadata fields
		if !excludedFields[key] {
			// Clean null or empty values that may cause issues
			if value != nil {
				cleaned[key] = value
			}
		}
	}
	
	return cleaned
}

// formatAkeneoErrors formats Akeneo errors to display useful information
func (c *Client) formatAkeneoErrors(errorResponse AkeneoErrorResponse) string {
	if len(errorResponse.Errors) == 0 {
		return errorResponse.Message
	}
	
	var errorMessages []string
	for _, fieldError := range errorResponse.Errors {
		errorMessages = append(errorMessages, fmt.Sprintf("Field '%s': %s", fieldError.Property, fieldError.Message))
	}
	
	return fmt.Sprintf("%s. Details: %s", errorResponse.Message, strings.Join(errorMessages, "; "))
}

// DebugRecord prints the content of a record for debugging purposes
func (c *Client) DebugRecord(entityName, code string, record ReferenceEntityRecord) {
	cleanRecord := c.cleanRecord(record)
	jsonData, _ := json.MarshalIndent(cleanRecord, "", "  ")
	fmt.Printf("🔍 DEBUG - Record %s/%s:\n%s\n", entityName, code, string(jsonData))
}