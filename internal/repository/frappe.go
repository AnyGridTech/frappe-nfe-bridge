package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/AnyGridTech/frappe-nfe-bridge/internal/models"
)

type FrappeRepository interface {
	GetCustomInvoice(id string) (*models.CustomFrappeInvoice, error)
	GetInvoice(id string) (*models.Invoices, error)
	GetTax(id string) (*models.FrappeTax, error)
	GetCarrier(id string) (*models.Carrier, error)
	UpdateInvoice(id string, data map[string]interface{}) error
}

type frappeRepo struct {
	client      *http.Client
	baseURL     string
	apiKey      string
	apiSecret   string
	docTypeName string // Default DocType name (e.g., "Invoices")
}

func NewFrappeRepo(baseURL, key, secret, docType string) FrappeRepository {
	return &frappeRepo{
		client:      &http.Client{},
		baseURL:     baseURL,
		apiKey:      key,
		apiSecret:   secret,
		docTypeName: docType,
	}
}

// GetInvoice retrieves the full Invoices doctype
func (r *frappeRepo) GetInvoice(id string) (*models.Invoices, error) {
	escapedDocType := url.PathEscape("Invoices")
	endpoint := fmt.Sprintf("%s/api/resource/%s/%s", r.baseURL, escapedDocType, id)

	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Authorization", fmt.Sprintf("token %s:%s", r.apiKey, r.apiSecret))

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("frappe returned status: %d", resp.StatusCode)
	}

	var result struct {
		Data models.Invoices `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// GetCustomInvoice for backward compatibility
func (r *frappeRepo) GetCustomInvoice(id string) (*models.CustomFrappeInvoice, error) {
	// Handle spaces in DocType names (e.g., "Brazil Invoice" -> "Brazil%20Invoice")
	escapedDocType := url.PathEscape(r.docTypeName)
	endpoint := fmt.Sprintf("%s/api/resource/%s/%s", r.baseURL, escapedDocType, id)

	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Authorization", fmt.Sprintf("token %s:%s", r.apiKey, r.apiSecret))

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("frappe returned status: %d", resp.StatusCode)
	}

	var result struct {
		Data models.CustomFrappeInvoice `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// GetTax retrieves a Tax template by name
func (r *frappeRepo) GetTax(id string) (*models.FrappeTax, error) {
	escapedDocType := url.PathEscape("Tax")
	endpoint := fmt.Sprintf("%s/api/resource/%s/%s", r.baseURL, escapedDocType, id)

	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Authorization", fmt.Sprintf("token %s:%s", r.apiKey, r.apiSecret))

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("frappe returned status: %d", resp.StatusCode)
	}

	var result struct {
		Data models.FrappeTax `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// GetCarrier retrieves carrier information
func (r *frappeRepo) GetCarrier(id string) (*models.Carrier, error) {
	escapedDocType := url.PathEscape("Carrier")
	endpoint := fmt.Sprintf("%s/api/resource/%s/%s", r.baseURL, escapedDocType, id)

	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Authorization", fmt.Sprintf("token %s:%s", r.apiKey, r.apiSecret))

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("frappe returned status: %d", resp.StatusCode)
	}

	var result struct {
		Data models.Carrier `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// UpdateInvoice updates an invoice document with NFe.io response data
func (r *frappeRepo) UpdateInvoice(id string, data map[string]interface{}) error {
	escapedDocType := url.PathEscape("Invoices")
	endpoint := fmt.Sprintf("%s/api/resource/%s/%s", r.baseURL, escapedDocType, id)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("PUT", endpoint, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", fmt.Sprintf("token %s:%s", r.apiKey, r.apiSecret))
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("frappe returned status: %d", resp.StatusCode)
	}

	return nil
}
