package models

// CustomFrappeInvoice represents your custom DocType from the AnyGridTech app
// Check your .json file in the "fields" array for the exact "fieldname" values
type CustomFrappeInvoice struct {
    Name          string  `json:"name"` // The ID (e.g., BR-INV-2024-001)
    CustomerName  string  `json:"customer_name"`
    
    // Brazilian specific fields usually found in custom apps
    CNPJ          string  `json:"cnpj_cpf"`      // Verify if this is 'cpf', 'cnpj', or a combined field
    NaturezaOp    string  `json:"natureza_operacao"` 
    Status        string  `json:"status"`
    GrandTotal    float64 `json:"grand_total"`
    
    // Child Table (Items)
    Items         []CustomFrappeItem `json:"items"`
}

type CustomFrappeItem struct {
    ItemCode    string  `json:"item_code"`
    ItemName    string  `json:"item_name"`
    Qty         float64 `json:"qty"`
    Rate        float64 `json:"rate"`
    Amount      float64 `json:"amount"`
    
    // Fiscal fields are critical for NFe.io
    NCM         string  `json:"ncm"`         // The NCM code is mandatory
    CFOP        string  `json:"cfop"`        // e.g., "5102"
    CEST        string  `json:"cest"`        // Optional but common
    TotalTaxes  float64 `json:"total_taxes"` // If your app pre-calculates this
}