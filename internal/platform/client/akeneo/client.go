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

// ReferenceEntity represents a Reference Entity definition
type ReferenceEntity map[string]interface{}

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

// Get ReferenceEntity retrieves a Reference Entity definition
func (c *Client) GetReferenceEntity(entityCode string) (ReferenceEntity, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/rest/v1/reference-entities/%s", c.config.Host, entityCode)

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

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("reference entity '%s' not found", entityCode)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error fetching reference entity: %d - %s", resp.StatusCode, string(body))
	}

	var entity ReferenceEntity
	if err := json.NewDecoder(resp.Body).Decode(&entity); err != nil {
		return nil, err
	}

	return entity, nil
}

// PatchReferenceEntity creates or updates a Reference Entity definition
func (c *Client) PatchReferenceEntity(entityCode string, entity ReferenceEntity) error {
	if err := c.ensureValidToken(); err != nil {
		return err
	}

	// Clean fields that should not be sent
	cleanEntity := c.cleanReferenceEntity(entity)

	jsonData, err := json.Marshal(cleanEntity)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/rest/v1/reference-entities/%s", c.config.Host, entityCode)

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
				return fmt.Errorf("validation error in reference entity %s: %s", entityCode, c.formatAkeneoErrors(errorResponse))
			}
		}

		return fmt.Errorf("error updating reference entity %s: %d - %s", entityCode, resp.StatusCode, string(body))
	}

	return nil
}

// cleanReferenceEntity removes fields that should not be sent in write operations
func (c *Client) cleanReferenceEntity(entity ReferenceEntity) ReferenceEntity {
	cleaned := make(ReferenceEntity)

	// List of fields to exclude
	excludedFields := map[string]bool{
		"_links": true,
	}

	for key, value := range entity {
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

// ReferenceEntityAttribute represents a Reference Entity attribute definition
type ReferenceEntityAttribute map[string]interface{}

// GetReferenceEntityAttributes retrieves all attributes from a Reference Entity
func (c *Client) GetReferenceEntityAttributes(entityCode string) ([]ReferenceEntityAttribute, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/rest/v1/reference-entities/%s/attributes", c.config.Host, entityCode)

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
		return nil, fmt.Errorf("error fetching reference entity attributes: %d - %s", resp.StatusCode, string(body))
	}

	// Read the body first to handle different response formats
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Debug: print raw response
	fmt.Printf("🔍 DEBUG - Raw attributes response:\n%s\n", string(body))

	// Try to unmarshal as array first (most common format)
	var attributes []ReferenceEntityAttribute
	if err := json.Unmarshal(body, &attributes); err == nil {
		// Normalize each attribute after unmarshalling
		for i := range attributes {
			attributes[i] = c.normalizeAttributeAfterFetch(attributes[i])
		}
		return attributes, nil
	}

	// If that fails, try the _embedded format
	var response struct {
		Embedded struct {
			Items []ReferenceEntityAttribute `json:"items"`
		} `json:"_embedded"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		// If both fail, return the original error with the body for debugging
		return nil, fmt.Errorf("error decoding attributes response: %w. Body: %s", err, string(body))
	}

	// Normalize each attribute after unmarshalling
	for i := range response.Embedded.Items {
		response.Embedded.Items[i] = c.normalizeAttributeAfterFetch(response.Embedded.Items[i])
	}

	return response.Embedded.Items, nil
}

// PatchReferenceEntityAttribute creates or updates a Reference Entity attribute
func (c *Client) PatchReferenceEntityAttribute(entityCode, attributeCode string, attribute ReferenceEntityAttribute) error {
	if err := c.ensureValidToken(); err != nil {
		return err
	}

	// Debug: print original attribute
	originalJSON, _ := json.MarshalIndent(attribute, "", "  ")
	fmt.Printf("🔍 DEBUG - Original attribute %s:\n%s\n", attributeCode, string(originalJSON))

	// Clean fields that should not be sent
	cleanAttribute := c.cleanReferenceEntityAttribute(attribute)

	// Use a custom encoder to ensure proper JSON formatting
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(cleanAttribute); err != nil {
		return err
	}

	jsonData := buf.Bytes()

	// Debug: print what we're sending
	fmt.Printf("🔍 DEBUG - Sending attribute %s:\n%s\n", attributeCode, string(jsonData))

	// Additional debug: verify by unmarshalling back
	var debugCheck map[string]interface{}
	json.Unmarshal(jsonData, &debugCheck)
	fmt.Printf("🔍 DEBUG - Labels type in JSON: %T, value: %v\n", debugCheck["labels"], debugCheck["labels"])

	// Extra debug: check raw bytes of labels field
	labelsJSON, _ := json.Marshal(debugCheck["labels"])
	fmt.Printf("🔍 DEBUG - Labels as JSON bytes: %s\n", string(labelsJSON))

	url := fmt.Sprintf("%s/api/rest/v1/reference-entities/%s/attributes/%s",
		c.config.Host, entityCode, attributeCode)

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

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
				return fmt.Errorf("validation error in attribute %s: %s", attributeCode, c.formatAkeneoErrors(errorResponse))
			}
		}

		return fmt.Errorf("error updating attribute %s: %d - %s", attributeCode, resp.StatusCode, string(body))
	}

	return nil
}

// cleanReferenceEntityAttribute removes fields that should not be sent in write operations
func (c *Client) cleanReferenceEntityAttribute(attribute ReferenceEntityAttribute) ReferenceEntityAttribute {
	cleaned := make(ReferenceEntityAttribute)

	// Get attribute code for default labels
	attributeCode, _ := attribute["code"].(string)

	// List of fields to exclude (metadata and null-only fields)
	excludedFields := map[string]bool{
		"_links":             true,
		"max_characters":     true, // Only include if not null
		"validation_regexp":  true, // Only include if not null
		"min_value":          true, // Only include if not null
		"max_value":          true, // Only include if not null
		"decimals_allowed":   true, // Only include if not null
		"max_file_size":      true, // Only include if not null
		"allowed_extensions": true, // Only include if not null/empty
	}

	// Required fields that should always be included
	requiredFields := map[string]bool{
		"code":                         true,
		"type":                         true,
		"labels":                       true,
		"value_per_locale":             true,
		"value_per_channel":            true,
		"is_required_for_completeness": true,
	}

	// Type-specific fields
	textFields := map[string]bool{
		"is_textarea":         true,
		"is_rich_text_editor": true,
		"validation_rule":     true,
		"max_characters":      true,
		"validation_regexp":   true,
	}

	numberFields := map[string]bool{
		"decimals_allowed": true,
		"min_value":        true,
		"max_value":        true,
	}

	imageFields := map[string]bool{
		"max_file_size":      true,
		"allowed_extensions": true,
	}

	// Get attribute type to know which fields to include
	attributeType, _ := attribute["type"].(string)

	for key, value := range attribute {
		// Skip excluded fields
		if excludedFields[key] {
			continue
		}

		// Skip null values unless it's a required field
		if value == nil && !requiredFields[key] {
			continue
		}

		// Include required fields
		if requiredFields[key] {
			if key == "labels" {
				cleaned[key] = c.normalizeLabels(value, attributeCode)
			} else {
				cleaned[key] = value
			}
			continue
		}

		// Include type-specific fields based on attribute type
		includeField := false
		switch attributeType {
		case "text":
			includeField = textFields[key]
		case "number":
			includeField = numberFields[key]
		case "image":
			includeField = imageFields[key]
		}

		if includeField && value != nil {
			if key == "allowed_extensions" {
				cleaned[key] = c.normalizeArray(value)
			} else {
				cleaned[key] = value
			}
		}
	}

	return cleaned
}

// normalizeAttributeAfterFetch normalizes an attribute right after fetching from API
func (c *Client) normalizeAttributeAfterFetch(attribute ReferenceEntityAttribute) ReferenceEntityAttribute {
	// Normalize labels if present
	if labels, exists := attribute["labels"]; exists {
		// Check if labels is an empty array and convert to empty map
		if labelsArray, ok := labels.([]interface{}); ok && len(labelsArray) == 0 {
			attribute["labels"] = map[string]string{}
		}
	}
	return attribute
}

// normalizeLabels converts labels from array format to object format if needed
func (c *Client) normalizeLabels(labels interface{}, attributeCode string) map[string]string {
	// Always return a proper map[string]string to ensure JSON marshals as object
	result := make(map[string]string)

	// If labels is nil, use attribute code as default label
	if labels == nil {
		result["en_US"] = attributeCode
		return result
	}

	// If labels is already a map[string]interface{}, convert it
	if labelsMap, ok := labels.(map[string]interface{}); ok {
		for k, v := range labelsMap {
			if strVal, ok := v.(string); ok {
				result[k] = strVal
			}
		}
		// If map is empty, add default label
		if len(result) == 0 {
			result["en_US"] = attributeCode
		}
		return result
	}

	// If labels is already a map[string]string, return it
	if labelsMap, ok := labels.(map[string]string); ok {
		// If map is empty, add default label
		if len(labelsMap) == 0 {
			result["en_US"] = attributeCode
			return result
		}
		return labelsMap
	}

	// If labels is an array, convert to map
	if labelsArray, ok := labels.([]interface{}); ok {
		for _, item := range labelsArray {
			if labelItem, ok := item.(map[string]interface{}); ok {
				// Try to get locale and label
				var locale, label string
				var hasLocale, hasLabel bool

				if loc, ok := labelItem["locale"].(string); ok {
					locale = loc
					hasLocale = true
				}

				if lbl, ok := labelItem["label"].(string); ok {
					label = lbl
					hasLabel = true
				}

				// Only add if we have both locale and label
				if hasLocale && hasLabel {
					result[locale] = label
				}
			}
		}
		// If no labels were extracted, add default
		if len(result) == 0 {
			result["en_US"] = attributeCode
		}
		return result
	}

	// If we couldn't convert, log warning and return default label
	fmt.Printf("⚠️  Warning: Could not normalize labels, using default. Original type: %T, value: %v\n", labels, labels)
	result["en_US"] = attributeCode
	return result
}

// normalizeArray ensures the value is an array
func (c *Client) normalizeArray(value interface{}) interface{} {
	// If it's already an array, return as is
	if _, ok := value.([]interface{}); ok {
		return value
	}

	// If it's nil, return empty array
	if value == nil {
		return []interface{}{}
	}

	// Return as is
	return value
}

// Product represents a product
type Product map[string]interface{}

// GetProduct retrieves a product by its identifier
func (c *Client) GetProduct(identifier string) (Product, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/rest/v1/products/%s", c.config.Host, identifier)

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

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("product '%s' not found", identifier)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error fetching product: %d - %s", resp.StatusCode, string(body))
	}

	var productData Product
	if err := json.NewDecoder(resp.Body).Decode(&productData); err != nil {
		return nil, err
	}

	return productData, nil
}

// PatchProduct creates or updates a product
func (c *Client) PatchProduct(identifier string, productData Product) error {
	if err := c.ensureValidToken(); err != nil {
		return err
	}

	// Clean fields that should not be sent
	cleanProduct := c.cleanProduct(productData)

	jsonData, err := json.Marshal(cleanProduct)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/rest/v1/products/%s", c.config.Host, identifier)

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(jsonData))
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
				return fmt.Errorf("validation error in product %s: %s", identifier, c.formatAkeneoErrors(errorResponse))
			}
		}

		return fmt.Errorf("error updating product %s: %d - %s", identifier, resp.StatusCode, string(body))
	}

	return nil
}

// cleanProduct removes fields that should not be sent in write operations
func (c *Client) cleanProduct(productData Product) Product {
	cleaned := make(Product)

	// List of fields to exclude
	excludedFields := map[string]bool{
		"_links":  true,
		"created": true,
		"updated": true,
	}

	for key, value := range productData {
		// Exclude metadata fields
		if !excludedFields[key] && value != nil {
			cleaned[key] = value
		}
	}

	return cleaned
}

// ProductModel represents a product model
type ProductModel map[string]interface{}

// GetProductModel retrieves a product model by its code
func (c *Client) GetProductModel(code string) (ProductModel, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/rest/v1/product-models/%s", c.config.Host, code)

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

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("product model '%s' not found", code)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error fetching product model: %d - %s", resp.StatusCode, string(body))
	}

	var model ProductModel
	if err := json.NewDecoder(resp.Body).Decode(&model); err != nil {
		return nil, err
	}

	return model, nil
}

// PatchProductModel creates or updates a product model
func (c *Client) PatchProductModel(code string, model ProductModel) error {
	if err := c.ensureValidToken(); err != nil {
		return err
	}

	// Clean fields that should not be sent
	cleanModel := c.cleanProductModel(model)

	jsonData, err := json.Marshal(cleanModel)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/rest/v1/product-models/%s", c.config.Host, code)

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(jsonData))
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

		if resp.StatusCode == http.StatusUnprocessableEntity {
			var errorResponse AkeneoErrorResponse
			if parseErr := json.Unmarshal(body, &errorResponse); parseErr == nil {
				return fmt.Errorf("validation error in product model %s: %s", code, c.formatAkeneoErrors(errorResponse))
			}
		}

		return fmt.Errorf("error updating product model %s: %d - %s", code, resp.StatusCode, string(body))
	}

	return nil
}

// GetProductsByParent retrieves all products with a specific parent
func (c *Client) GetProductsByParent(parentCode string) ([]Product, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}

	var allProducts []Product
	page := 1
	limit := 100

	for {
		url := fmt.Sprintf("%s/api/rest/v1/products?search={\"parent\":[{\"operator\":\"=\",\"value\":\"%s\"}]}&page=%d&limit=%d",
			c.config.Host, parentCode, page, limit)

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
			return nil, fmt.Errorf("error fetching products by parent: %d - %s", resp.StatusCode, string(body))
		}

		var response struct {
			Embedded struct {
				Items []Product `json:"items"`
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

		allProducts = append(allProducts, response.Embedded.Items...)

		if response.Links.Next == nil {
			break
		}

		page++
	}

	return allProducts, nil
}

// GetProductModelsByParent retrieves all product models with a specific parent
func (c *Client) GetProductModelsByParent(parentCode string) ([]ProductModel, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}

	var allModels []ProductModel
	page := 1
	limit := 100

	for {
		url := fmt.Sprintf("%s/api/rest/v1/product-models?search={\"parent\":[{\"operator\":\"=\",\"value\":\"%s\"}]}&page=%d&limit=%d",
			c.config.Host, parentCode, page, limit)

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
			return nil, fmt.Errorf("error fetching product models by parent: %d - %s", resp.StatusCode, string(body))
		}

		var response struct {
			Embedded struct {
				Items []ProductModel `json:"items"`
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

		allModels = append(allModels, response.Embedded.Items...)

		if response.Links.Next == nil {
			break
		}

		page++
	}

	return allModels, nil
}

// cleanProductModel removes fields that should not be sent in write operations
func (c *Client) cleanProductModel(model ProductModel) ProductModel {
	cleaned := make(ProductModel)

	// List of fields to exclude
	excludedFields := map[string]bool{
		"_links":  true,
		"created": true,
		"updated": true,
	}

	for key, value := range model {
		if !excludedFields[key] && value != nil {
			cleaned[key] = value
		}
	}

	return cleaned
}
