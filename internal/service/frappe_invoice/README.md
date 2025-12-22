# Service Layer - NFe.io Integration

## Overview
This service layer has been adapted from the older API structure found in `docs/invoice/` to work with the current project architecture. The code maintains the same business logic while fitting into the new repository pattern.

## Files and Their Purpose

### 1. `tax.go` - Tax Calculation Service
**Adapted from:** `docs/invoice/build.go` (tax-related functions)

**Key Features:**
- `TaxService` struct with dependency-free tax calculations
- `CalculateTax()` - Main method to calculate all taxes for an item
- Individual tax calculators:
  - `calculateICMS()` - ICMS (state tax)
  - `calculatePIS()` - PIS (federal social contribution)
  - `calculateCOFINS()` - COFINS (federal social contribution)
  - `calculateIPI()` - IPI (federal product tax)
- `calculateDIFAL()` - DIFAL calculation (dual base method)
- `calculateDIFALSimple()` - DIFAL simple calculation (single base)
- `calculateBaseTax()` - Base value calculation for taxes

**What Changed:**
- Removed dependency on `dto.Items` struct
- Created `TaxInput` struct for clean input data
- Made all methods part of `TaxService` for better testability
- Properly formatted tax calculations with rounding

**Usage Example:**
```go
taxService := NewTaxService()
input := TaxInput{
    ItemValue:       100.00,
    Quantity:        2.0,
    AliqICMS:        18.0,
    AliquotaPis:     1.65,
    AliquotaCofins:  7.6,
}
tax := taxService.CalculateTax(input)
```

### 2. `issuer.go` - Main Issuer Service
**Adapted from:** `docs/invoice/index.go`

**Key Features:**
- `IssuerService` interface and implementation
- `IssueNoteForFrappeInvoice()` - Main orchestration method
- `mapFrappeToNFe()` - Converts Frappe data to NFe.io format
- Buyer information building with CPF/CNPJ detection
- Phone and tax number formatting utilities

**What Changed:**
- Integrated with repository pattern (FrappeRepository, NFeRepository)
- Added TaxService integration
- Simplified buyer building logic
- Removed direct HTTP calls (now handled by repositories)
- Added proper error handling throughout

**Flow:**
1. Fetch invoice from Frappe via `frappeRepo`
2. Map Frappe data to NFe.io format
3. Calculate taxes for each item
4. Send to NFe.io via `nfeRepo`
5. Return response with PDF/XML URLs

### 3. `builder.go` - Invoice Builder Service (NEW)
**Adapted from:** `docs/invoice/build.go` and `docs/invoice/transport.go`

**Key Features:**
- CFOP determination based on operation nature
- Transport information building
- Address building
- Volume and weight calculations
- Buyer type determination (CPF vs CNPJ)
- State tax number indicator determination

**Main Methods:**
- `DetermineCFOP()` - Determines CFOP code based on operation and states
- `resolveCFOPCode()` - Resolves CFOP template (x102 → 5102 or 6102)
- `DetermineOperationType()` - Sets operation type and serie
- `DetermineDestination()` - Internal vs interstate operation
- `BuildTransport()` - Creates transport information
- `BuildAddress()` - Creates address information
- `DetermineBuyerType()` - Identifies CPF or CNPJ
- `DetermineStateTaxNumberIndicator()` - Sets ICMS contributor status

**CFOP Logic:**
```
Outgoing + Same State → 5xxx
Outgoing + Different State → 6xxx
Incoming + Same State → 1xxx
Incoming + Different State → 2xxx
```

## Key Concepts from Old API (docs/)

### From `docs/invoice/build.go`:
1. **Operation Nature Constants:** naturalPerson, legalEntity, taxPayer, etc.
2. **Interface Pattern:** OperationNatureIE, ItemsIE, BuyerIE
3. **CFOP Template System:** Using "x" as placeholder (x102 → 5102/6102/1102/2102)
4. **Buyer Types:** PF (Pessoa Física) and PJ (Pessoa Jurídica)
5. **Tax Calculations:** ICMS, PIS, COFINS, IPI with proper rounding

### From `docs/invoice/transport.go`:
1. **Transport Modalities:** ByIssuer, ByBuyer, ThirdParty, etc.
2. **Volume Information:** Species, Brand, Numeration, Weights
3. **Carrier Information:** Full legal entity data required

### From `docs/invoice/consult.go`:
1. **Access Key Consultation:** Retrieve invoice by access key (not yet implemented)

## Integration Points

### With Frappe/ERPNext:
- `FrappeRepository` fetches invoice data
- Expected Frappe fields: customer_name, cnpj_cpf, items with NCM, CFOP, tax info

### With NFe.io:
- `NFeRepository` sends formatted invoice
- Returns ProductInvoiceResponse with ID, status, PDF URL, XML URL

## Configuration Required

1. **Environment Variables:**
   - Frappe API credentials
   - NFe.io API key
   - Company ID for NFe.io

2. **Tax Rates:**
   - Should be configured per item or per company
   - Currently using zeros as defaults (need configuration)

3. **State Mapping:**
   - Issuer state (your company's state)
   - Used for CFOP and destination determination

## Future Enhancements

1. **Complete Transport Integration:**
   - Currently stubbed in `mapFrappeToNFe()`
   - Use `BuilderService.BuildTransport()`

2. **Address Mapping:**
   - Implement full address building from Frappe
   - Use `BuilderService.BuildAddress()`

3. **Tax Configuration:**
   - Load tax rates from Frappe or configuration
   - Per-item or per-category tax profiles

4. **DIFAL Support:**
   - Implement DIFAL calculation for interstate operations to final consumers
   - Already have `calculateDIFAL()` methods available

5. **Additional Information:**
   - AdditionalInformation field support
   - Referenced documents support

6. **Webhook Integration:**
   - NFe.io status updates
   - Update Frappe invoice with authorization data

## Testing

To test the service layer:

```go
// Mock repositories
frappeRepo := &mockFrappeRepo{}
nfeRepo := &mockNFeRepo{}

// Create service
issuerService := NewIssuerService(frappeRepo, nfeRepo, "companyID")

// Issue invoice
response, err := issuerService.IssueNoteForFrappeInvoice("INV-001")
```

## Migration Notes from Old API

### What Was Removed:
- Direct HTTP calls (moved to repositories)
- Hardcoded modules (erpnext-go/pkg/modules)
- File writing (productInvoice.json) - now handled by caller if needed
- Inline error handling with log.Println - now returns errors properly

### What Was Added:
- Clean separation of concerns (tax, builder, issuer)
- Testable pure functions
- Proper error propagation
- Repository pattern integration
- Input/Output DTOs for each service

### What Stayed the Same:
- Business logic for tax calculations
- CFOP determination algorithm
- Buyer type logic (CPF/CNPJ)
- Tax formulas and rounding
- Constants and enumerations
