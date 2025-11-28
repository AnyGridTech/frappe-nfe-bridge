# NFe.io Repository - Complete API Integration

## Overview

The `internal/repository/nfeio.go` file has been enhanced with comprehensive NFe.io API integration, adapted from `docs/modules/nfeio.go`. It now provides full support for all invoice operations including creation, consultation, deletion, PDF/XML retrieval, and correction letters.

## Features

### ✅ Product Invoice Operations
- Create new product invoices
- Retrieve invoice by ID
- Retrieve invoice by access key
- Delete (cancel) invoices

### ✅ Document Operations
- Get invoice PDF
- Get invoice XML

### ✅ Correction Letter Operations
- Create correction letters
- Get correction letter PDF
- Get correction letter XML

## Interface

```go
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
```

## Configuration

### Environment Variables

Add these to your `.env` file:

```env
# NFe.io Configuration
NFE_API_KEY=your-nfe-io-api-key-here
NFE_COMPANY_ID=your-company-id
NFE_ENDPOINT=https://api.nfe.io/v2
NFE_ENDPOINT_CONSULT=https://api.nfe.io/v2

# For testing/sandbox
# NFE_ENDPOINT=https://sandbox.nfe.io/v2
# NFE_ENDPOINT_CONSULT=https://sandbox.nfe.io/v2
```

### Initialization

The repository is now initialized with endpoints:

```go
nfeRepo := repository.NewNFeRepo(
    "https://api.nfe.io/v2",           // endpoint
    "https://api.nfe.io/v2",           // endpoint consult
    "your-api-key",                    // API key
)
```

This is already configured in `cmd/main.go` to read from environment variables.

## Usage Examples

### 1. Create Invoice (Already Implemented)

```go
response, err := nfeRepo.CreateProductInvoice(payload)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Invoice ID: %s\n", response.ID)
```

### 2. Get Invoice by ID

```go
invoice, err := nfeRepo.GetInvoice("invoice-id-123", "company-key")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Status: %s\n", invoice.Status)
```

### 3. Get Invoice by Access Key

```go
// Access key is the 44-digit NFe key
accessKey := "35210812345678000190550010000001001234567890"
invoice, err := nfeRepo.GetInvoiceByAccessKey(accessKey)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Found invoice: %s\n", invoice.ID)
```

### 4. Download Invoice PDF

```go
pdfBytes, err := nfeRepo.GetInvoicePDF("invoice-id-123", "company-key")
if err != nil {
    log.Fatal(err)
}

// Save to file
err = os.WriteFile("invoice.pdf", pdfBytes, 0644)
if err != nil {
    log.Fatal(err)
}
```

### 5. Download Invoice XML

```go
xmlBytes, err := nfeRepo.GetInvoiceXML("invoice-id-123", "company-key")
if err != nil {
    log.Fatal(err)
}

// Save to file or send to Frappe
err = os.WriteFile("invoice.xml", xmlBytes, 0644)
```

### 6. Cancel Invoice

```go
err := nfeRepo.DeleteInvoice("invoice-id-123", "company-key")
if err != nil {
    log.Fatal(err)
}
fmt.Println("Invoice cancelled successfully")
```

### 7. Create Correction Letter

```go
reason := "Correção de dados do cliente - nome incorreto"
response, err := nfeRepo.CreateCorrectionLetter(
    "invoice-id-123",
    "company-key",
    reason,
)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Correction letter created: %s\n", response.ID)
```

### 8. Download Correction Letter PDF

```go
pdfBytes, err := nfeRepo.GetCorrectionLetterPDF("invoice-id-123", "company-key")
if err != nil {
    log.Fatal(err)
}
os.WriteFile("correction_letter.pdf", pdfBytes, 0644)
```

## Integration with Service Layer

### Enhance IssuerService

You can now enhance the `IssuerService` to use these new methods:

```go
type IssuerService interface {
    IssueNoteForFrappeInvoice(invoiceID string) (*models.ProductInvoiceResponse, error)
    GetInvoiceStatus(nfeID string) (*models.ProductInvoiceResponse, error)
    DownloadInvoicePDF(nfeID string) ([]byte, error)
    CancelInvoice(nfeID string) error
    CreateCorrection(nfeID, reason string) (*models.ProductInvoiceResponse, error)
}
```

Example implementation:

```go
func (s *issuerService) GetInvoiceStatus(nfeID string) (*models.ProductInvoiceResponse, error) {
    return s.nfeRepo.GetInvoice(nfeID, s.companyID)
}

func (s *issuerService) DownloadInvoicePDF(nfeID string) ([]byte, error) {
    return s.nfeRepo.GetInvoicePDF(nfeID, s.companyID)
}

func (s *issuerService) CancelInvoice(nfeID string) error {
    return s.nfeRepo.DeleteInvoice(nfeID, s.companyID)
}

func (s *issuerService) CreateCorrection(nfeID, reason string) (*models.ProductInvoiceResponse, error) {
    return s.nfeRepo.CreateCorrectionLetter(nfeID, s.companyID, reason)
}
```

## Handler Examples

### Add New Endpoints

```go
// Get invoice status
v1.Get("/invoices/:id", invoiceHandler.HandleGetInvoice)

// Download PDF
v1.Get("/invoices/:id/pdf", invoiceHandler.HandleDownloadPDF)

// Cancel invoice
v1.Delete("/invoices/:id", invoiceHandler.HandleCancelInvoice)

// Create correction letter
v1.Post("/invoices/:id/correction", invoiceHandler.HandleCreateCorrection)
```

Handler implementations:

```go
func (h *InvoiceHandler) HandleGetInvoice(c *fiber.Ctx) error {
    id := c.Params("id")
    
    invoice, err := h.svc.GetInvoiceStatus(id)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.JSON(invoice)
}

func (h *InvoiceHandler) HandleDownloadPDF(c *fiber.Ctx) error {
    id := c.Params("id")
    
    pdfBytes, err := h.svc.DownloadInvoicePDF(id)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    c.Set("Content-Type", "application/pdf")
    c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=invoice_%s.pdf", id))
    return c.Send(pdfBytes)
}

func (h *InvoiceHandler) HandleCancelInvoice(c *fiber.Ctx) error {
    id := c.Params("id")
    
    err := h.svc.CancelInvoice(id)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.JSON(fiber.Map{
        "message": "Invoice cancelled successfully",
        "id":      id,
    })
}

func (h *InvoiceHandler) HandleCreateCorrection(c *fiber.Ctx) error {
    id := c.Params("id")
    
    var req struct {
        Reason string `json:"reason"`
    }
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    if req.Reason == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Reason is required",
        })
    }
    
    response, err := h.svc.CreateCorrection(id, req.Reason)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.JSON(fiber.Map{
        "message": "Correction letter created",
        "data":    response,
    })
}
```

## Error Handling

All methods now include comprehensive error handling with:
- Proper error wrapping with context
- HTTP status code checking
- Response body reading on errors
- Clear error messages

Example error response:
```
nfe.io API error: status 400, body: {"message": "Invalid CNPJ format"}
```

## Improvements Over Original Code

### ✅ Better Error Handling
- All errors are wrapped with context
- HTTP response bodies are read on errors
- Clear error messages for debugging

### ✅ Type Safety
- Proper method signatures
- Returns typed responses instead of raw http.Response
- No more manual JSON decoding in calling code

### ✅ Consistent API
- All methods follow the same pattern
- Consistent error handling
- Predictable return types

### ✅ Cleaner Code
- Removed deprecated ViaCep method (not NFe.io related)
- Removed logger dependency (errors are returned)
- Single responsibility principle

### ✅ Testability
- Easy to mock interface
- No global state
- Dependency injection ready

## Migration Notes

### From Old Code (`docs/modules/nfeio.go`)

**Before:**
```go
nfeio := modules.NewNfeioModule(logger)
resp, err := nfeio.CreateInvoice(invoice, companyKey)
// Manual response handling
body, _ := io.ReadAll(resp.Body)
json.Unmarshal(body, &result)
```

**After:**
```go
nfeRepo := repository.NewNFeRepo(endpoint, endpointConsult, apiKey)
result, err := nfeRepo.CreateProductInvoice(invoice)
// result is already parsed and ready to use
```

## Testing

### Unit Tests

Create mocks for testing:

```go
type MockNFeRepo struct {
    CreateProductInvoiceFunc func(*models.ProductInvoiceRequest) (*models.ProductInvoiceResponse, error)
    GetInvoiceFunc           func(string, string) (*models.ProductInvoiceResponse, error)
    // ... other methods
}

func (m *MockNFeRepo) CreateProductInvoice(req *models.ProductInvoiceRequest) (*models.ProductInvoiceResponse, error) {
    return m.CreateProductInvoiceFunc(req)
}
```

### Integration Tests

Test against NFe.io sandbox:

```go
func TestNFeRepoIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    repo := repository.NewNFeRepo(
        "https://sandbox.nfe.io/v2",
        "https://sandbox.nfe.io/v2",
        os.Getenv("NFE_SANDBOX_KEY"),
    )
    
    // Test create invoice
    response, err := repo.CreateProductInvoice(testPayload)
    assert.NoError(t, err)
    assert.NotEmpty(t, response.ID)
}
```

## Next Steps

1. **Add Webhook Handler**: Implement NFe.io webhook handler to receive status updates
2. **Frappe Integration**: Update Frappe with PDF/XML URLs after issuance
3. **Status Polling**: Implement background job to check invoice status
4. **Retry Logic**: Add retry mechanism for failed API calls
5. **Caching**: Cache frequently accessed invoices
6. **Metrics**: Add logging/metrics for API calls

## API Documentation

For complete NFe.io API documentation, visit:
- Production: https://api.nfe.io/docs
- Sandbox: https://sandbox.nfe.io/docs

## Support

For issues related to NFe.io API, contact their support or check their documentation.
