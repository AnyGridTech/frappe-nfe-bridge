package repository

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "nfe-go/internal/models"
)

type FrappeRepository interface {
    GetCustomInvoice(id string) (*models.CustomFrappeInvoice, error)
}

type frappeRepo struct {
    client      *http.Client
    baseURL     string
    apiKey      string
    apiSecret   string
    docTypeName string // e.g., "Brazil Invoice"
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