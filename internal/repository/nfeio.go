package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AnyGridTech/frappe-nfe-bridge/internal/models"
)

type NFeRepository interface {
	CreateProductInvoice(req *models.ProductInvoiceRequest) (*models.ProductInvoiceResponse, error)
}

type nfeRepo struct {
	apiKey string
}

func NewNFeRepo(apiKey string) NFeRepository {
	return &nfeRepo{apiKey: apiKey}
}

func (r *nfeRepo) CreateProductInvoice(payload *models.ProductInvoiceRequest) (*models.ProductInvoiceResponse, error) {
	url := "https://api.nfe.io/v2/productinvoices" // Confirm exact endpoint in docs

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		// Read body to see error message from NFe.io
		return nil, fmt.Errorf("nfe.io API error: status %d", resp.StatusCode)
	}

	var result models.ProductInvoiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
