package models

// ProductInvoiceRequest represents the JSON payload for NFe.io V2
type ProductInvoiceRequest struct {
    CompanyID string   `json:"company_id"` // Your ID in NFe.io
    Reference string   `json:"reference"`  // External ID (e.g., ERPNext Invoice Name)
    Items     []Item   `json:"items"`
    Buyer     Buyer    `json:"buyer"`
}

type Item struct {
    Code        string  `json:"code"`
    Description string  `json:"description"`
    NCM         string  `json:"ncm"`
    CFOP        string  `json:"cfop"`
    Quantity    float64 `json:"quantity"`
    UnitValue   float64 `json:"unit_value"` // Value per item
    TotalValue  float64 `json:"total_value"`
}

type Buyer struct {
    Name             string  `json:"name"`
    FederalTaxNumber string  `json:"federalTaxNumber"` // CPF or CNPJ
    Email            string  `json:"email"`
    Address          Address `json:"address"`
}

type Address struct {
    Street       string `json:"street"`
    Number       string `json:"number"`
    Neighborhood string `json:"neighborhood"`
    CityCode     string `json:"cityCode"` // IBGE Code
    State        string `json:"state"`    // UF (e.g., SP)
    ZipCode      string `json:"postalCode"`
}

// Response from NFe.io
type ProductInvoiceResponse struct {
    ID              string `json:"id"`
    Status          string `json:"status"`
    Environment     string `json:"environment"`
    FlowStatus      string `json:"flowStatus"`
    PdfUrl          string `json:"pdf"`
    XmlUrl          string `json:"xml"`
}