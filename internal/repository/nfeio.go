package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AnyGridTech/frappe-nfe-bridge/internal/models"
)

// NFeRepository handles all NFe.io API operations
// Adapted from docs/modules/nfeio.go
type NFeRepository interface {
	// Product Invoice operations
	CreateProductInvoice(req *models.ProductInvoiceRequest) (*models.ProductInvoiceResponse, error)
	GetInvoice(id, companyKey string) (*models.ProductInvoiceResponse, error)
	GetInvoiceByAccessKey(accessKey string) (*models.ProductInvoiceResponse, error)
	DeleteInvoice(id, companyKey string) error

	// PDF and XML operations
	GetInvoicePDF(id, companyKey string) ([]byte, error)
	GetInvoiceXML(id, companyKey string) ([]byte, error)

	// Correction letter operations
	CreateCorrectionLetter(id, companyKey, reason string) (*models.ProductInvoiceResponse, error)
	GetCorrectionLetterPDF(id, companyKey string) ([]byte, error)
	GetCorrectionLetterXML(id, companyKey string) ([]byte, error)
}

type nfeRepo struct {
	apiKey          string
	endpoint        string
	endpointConsult string
	client          *http.Client
}

// NewNFeRepo creates a new NFe.io repository instance
// endpoint: base URL for NFe.io API (e.g., "https://api.nfe.io/v2")
// endpointConsult: base URL for consultation API
// apiKey: your NFe.io API key
func NewNFeRepo(endpoint, endpointConsult, apiKey string) NFeRepository {
	return &nfeRepo{
		apiKey:          apiKey,
		endpoint:        endpoint,
		endpointConsult: endpointConsult,
		client:          &http.Client{},
	}
}

// CreateProductInvoice creates a new product invoice
func (r *nfeRepo) CreateProductInvoice(payload *models.ProductInvoiceRequest) (*models.ProductInvoiceResponse, error) {
	url := fmt.Sprintf("%s/productinvoices?apiKey=%s", r.endpoint, r.apiKey)

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("nfe.io API error: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result models.ProductInvoiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetInvoice retrieves an invoice by ID
func (r *nfeRepo) GetInvoice(id, companyKey string) (*models.ProductInvoiceResponse, error) {
	url := fmt.Sprintf("%s/%s/productinvoices/%s?apiKey=%s", r.endpoint, companyKey, id, r.apiKey)

	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("nfe.io API error: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result models.ProductInvoiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetInvoiceByAccessKey retrieves an invoice by access key
func (r *nfeRepo) GetInvoiceByAccessKey(accessKey string) (*models.ProductInvoiceResponse, error) {
	url := fmt.Sprintf("%s/productinvoices/%s?apikey=%s", r.endpointConsult, accessKey, r.apiKey)

	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice by access key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("nfe.io API error: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result models.ProductInvoiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// DeleteInvoice deletes (cancels) an invoice
func (r *nfeRepo) DeleteInvoice(id, companyKey string) error {
	url := fmt.Sprintf("%s/%s/productinvoices/%s?apikey=%s", r.endpoint, companyKey, id, r.apiKey)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete invoice: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("nfe.io API error: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// GetInvoicePDF retrieves the PDF of an invoice
func (r *nfeRepo) GetInvoicePDF(id, companyKey string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/productinvoices/%s/pdf?apikey=%s", r.endpoint, companyKey, id, r.apiKey)

	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice PDF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("nfe.io API error: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	pdfBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF content: %w", err)
	}

	return pdfBytes, nil
}

// GetInvoiceXML retrieves the XML of an invoice
func (r *nfeRepo) GetInvoiceXML(id, companyKey string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/productinvoices/%s/xml?apikey=%s", r.endpoint, companyKey, id, r.apiKey)

	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice XML: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("nfe.io API error: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	xmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML content: %w", err)
	}

	return xmlBytes, nil
}

// CreateCorrectionLetter creates a correction letter for an invoice
func (r *nfeRepo) CreateCorrectionLetter(id, companyKey, reason string) (*models.ProductInvoiceResponse, error) {
	url := fmt.Sprintf("%s/%s/productinvoices/%s/correctionletter?apiKey=%s", r.endpoint, companyKey, id, r.apiKey)

	payload := map[string]string{"reason": reason}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal reason: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create correction letter: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("nfe.io API error: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result models.ProductInvoiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetCorrectionLetterPDF retrieves the PDF of a correction letter
func (r *nfeRepo) GetCorrectionLetterPDF(id, companyKey string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/productinvoices/%s/correctionletter/pdf?apikey=%s", r.endpoint, companyKey, id, r.apiKey)

	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get correction letter PDF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("nfe.io API error: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	pdfBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF content: %w", err)
	}

	return pdfBytes, nil
}

// GetCorrectionLetterXML retrieves the XML of a correction letter
func (r *nfeRepo) GetCorrectionLetterXML(id, companyKey string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/productinvoices/%s/correctionletter/xml?apikey=%s", r.endpoint, companyKey, id, r.apiKey)

	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get correction letter XML: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("nfe.io API error: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	xmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML content: %w", err)
	}

	return xmlBytes, nil
}
