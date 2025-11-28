# 🎯 Quick Integration Guide

## What You Have Now

```
internal/service/
├── 📄 tax.go              - Tax calculation engine
├── 📄 builder.go          - Invoice building utilities  
├── 📄 issuer.go           - Main orchestration service
├── 📄 examples_test.go    - Usage examples & tests
└── 📝 README.md           - Detailed documentation
```

## 🚀 Quick Start

### 1. Calculate Taxes for an Item

```go
import "github.com/AnyGridTech/frappe-nfe-bridge/internal/service"

taxService := service.NewTaxService()

tax := taxService.CalculateTax(service.TaxInput{
    ItemValue:       100.00,
    Quantity:        2.0,
    AliqICMS:        18.0,
    AliquotaPis:     1.65,
    AliquotaCofins:  7.6,
    OrigemICMS:      "0",
    CSTICMS:         "00",
    CSTPis:          "01",
    CSTCofins:       "01",
    ModDetermBC:     "3",
})

// tax.TotalTax contains the sum of all taxes
// tax.Icms, tax.Pis, tax.Cofins, tax.Ipi are ready to use
```

### 2. Determine CFOP Automatically

```go
builderService := service.NewBuilderService()

cfop, err := builderService.DetermineCFOP(
    "venda",      // operation nature
    "outgoing",   // operation type
    "RJ",         // buyer state
    "SP",         // issuer state
)
// Returns: 6102 (interstate sale)
```

### 3. Issue a Complete Invoice

```go
// Setup (do this once in main.go or dependency injection)
frappeRepo := repository.NewFrappeRepo(
    "https://your-erp.com",
    "apiKey",
    "apiSecret",
    "Brazil Invoice", // DocType name
)

nfeRepo := repository.NewNFeRepo("your-nfe-io-api-key")

issuerService := service.NewIssuerService(
    frappeRepo,
    nfeRepo,
    "your-company-id",
)

// Issue invoice
response, err := issuerService.IssueNoteForFrappeInvoice("INV-001")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("✅ Invoice issued: %s\n", response.ID)
fmt.Printf("📄 PDF: %s\n", response.PdfUrl)
fmt.Printf("📋 XML: %s\n", response.XmlUrl)
```

## 📊 Tax Calculation Formula

```
Base Tax = (Unit Value × Quantity + Freight + Insurance + Others) - Discount

ICMS Amount    = Base Tax × ICMS Rate
PIS Amount     = Base Tax × PIS Rate
COFINS Amount  = Base Tax × COFINS Rate
IPI Amount     = Base Tax × IPI Rate

Total Tax = ICMS + PIS + COFINS + IPI
```

## 🗺️ CFOP Logic Map

```
┌─────────────────────────────────────────────────┐
│              Operation Type                      │
├─────────────────────────────────────────────────┤
│                                                  │
│  OUTGOING (Sales/Returns)                       │
│  ├─ Same State    → 5xxx (e.g., 5102)          │
│  └─ Other State   → 6xxx (e.g., 6102)          │
│                                                  │
│  INCOMING (Purchases/Receipts)                  │
│  ├─ Same State    → 1xxx (e.g., 1102)          │
│  └─ Other State   → 2xxx (e.g., 2102)          │
│                                                  │
└─────────────────────────────────────────────────┘
```

## 🔧 Common Operations

### Get Buyer Type (CPF vs CNPJ)

```go
builderService := service.NewBuilderService()

// CPF (11 digits)
buyerType, taxRegime, _ := builderService.DetermineBuyerType("12345678900")
// Returns: "naturalPerson", "none"

// CNPJ (14 digits)  
buyerType, taxRegime, _ := builderService.DetermineBuyerType("12345678000190")
// Returns: "legalEntity", ""
```

### Build Address

```go
address, err := builderService.BuildAddress(service.AddressInput{
    Street:       "Rua das Flores",
    Number:       "123",
    Neighborhood: "Centro",
    City:         "São Paulo",
    CityCode:     "3550308", // IBGE code
    State:        "SP",
    PostalCode:   "01234567",
    Phone:        "11987654321",
})
```

### Build Transport

```go
transport, err := builderService.BuildTransport(service.TransportInput{
    CarrierName:     "Transportadora XYZ",
    CarrierCNPJ:     "12345678000190",
    CarrierEmail:    "contato@transportadora.com",
    FreightModality: "ByIssuer",
    VolumeQuantity:  1,
    Species:         "Caixa",
    NetWeight:       10.5,
    GrossWeight:     12.0,
    // ... other fields
})
```

## ⚙️ Configuration Checklist

- [ ] Set Frappe API credentials
- [ ] Set NFe.io API key
- [ ] Set Company ID for NFe.io
- [ ] Configure tax rates per item/category
- [ ] Map Frappe fields to tax configuration
- [ ] Configure issuer state for CFOP logic
- [ ] Test with NFe.io sandbox first

## 🧪 Testing

Run the example tests:

```powershell
cd c:\Users\Luigi\Desktop\nfe-go
go test ./internal/service/... -v
```

## 📝 Next Steps

1. **Customize Tax Configuration**
   - Load from Frappe item master
   - Or use company-wide defaults
   - Store in config file

2. **Complete Address Mapping**
   - Map Frappe customer fields
   - Add delivery address support
   - Validate IBGE codes

3. **Add Transport Integration**
   - Fetch carrier from Frappe
   - Calculate volumes automatically
   - Support multiple carriers

4. **Implement Webhooks**
   - Listen to NFe.io status updates
   - Update Frappe with authorization data
   - Handle cancellations

## 🆘 Troubleshooting

### Tax Calculation Returns 0

```go
// ❌ Wrong - rate as decimal
AliqICMS: 0.18

// ✅ Correct - rate as percentage
AliqICMS: 18.0
```

### CFOP Error

Make sure states are 2-letter codes: "SP", "RJ", "MG", etc.

### Buyer Type Error

Tax number must be clean digits only:
- CPF: 11 digits
- CNPJ: 14 digits

## 📚 More Information

- See `README.md` for complete documentation
- See `examples_test.go` for more usage examples
- Check `ADAPTATION_SUMMARY.md` for migration details

## 🎉 You're Ready!

All business logic from the old API is now adapted and ready to use in your clean architecture!
