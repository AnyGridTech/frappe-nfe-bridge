# NFe-Go: Brazilian Electronic Invoice Bridge

A Go-based microservice that bridges Frappe/ERPNext with the NFe.io API for generating Brazilian electronic invoices (Nota Fiscal Eletrônica).

[![Go Version](https://img.shields.io/badge/Go-1.x-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## 🎯 Overview

This service provides a clean, maintainable bridge between Frappe/ERPNext and the NFe.io API for Brazilian electronic invoice generation. It handles:

- ✅ Complete NFe.io API integration (create, retrieve, cancel, corrections)
- ✅ Brazilian tax calculations (ICMS, PIS, COFINS, IPI, DIFAL)
- ✅ Frappe custom app integration
- ✅ CFOP determination and fiscal rules
- ✅ Address and transport data handling
- ✅ PDF and XML document generation

## 📁 Project Structure

```
nfe-go/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── handler/
│   │   └── invoice_hdl.go      # HTTP handlers
│   ├── middleware/
│   │   └── auth.go             # Authentication middleware
│   ├── models/
│   │   ├── frappe.go           # Frappe DocType models
│   │   └── invoice.go          # NFe.io models
│   ├── repository/
│   │   ├── frappe.go           # Frappe API client
│   │   └── nfeio.go            # NFe.io API client
│   ├── router/
│   │   └── nfeio.go            # Route definitions
│   └── service/
│       ├── builder.go          # CFOP & address building
│       ├── issuer.go           # Main invoice orchestration
│       └── tax.go              # Tax calculation engine
├── docs/
│   └── frappe_brazil_invoice/  # Frappe custom app (reference)
├── FRAPPE_INTEGRATION.md       # Detailed integration guide
└── README.md                   # This file
```

## 🚀 Quick Start

### Prerequisites

- Go 1.x or higher
- Access to NFe.io API (API key and Company ID)
- Frappe/ERPNext instance with Brazil Invoice app
- PostgreSQL (optional, for local development)

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/AnyGridTech/frappe-nfe-bridge.git
   cd frappe-nfe-bridge
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Configure environment**:
   Create a `.env` file:
   ```bash
   # Frappe Configuration
   FRAPPE_URL=https://your-frappe-site.com
   FRAPPE_API_KEY=your_api_key
   FRAPPE_API_SECRET=your_api_secret

   # NFe.io Configuration
   NFEIO_API_KEY=your_nfeio_api_key
   COMPANY_ID=your_company_id
   NFE_ENDPOINT=https://api.nfe.io/v2/companies/{company_id}/productinvoices
   NFE_ENDPOINT_CONSULT=https://api.nfe.io/v2/companies/{company_id}/productinvoices

   # Server Configuration
   PORT=3000
   ENVIRONMENT=development  # or production
   ```

4. **Run the service**:
   ```bash
   go run cmd/main.go
   ```

   The service will start on `http://localhost:3000`

### Docker Deployment

```bash
# Build image
docker build -t nfe-go .

# Run container
docker run -p 3000:3000 --env-file .env nfe-go
```

## 📡 API Endpoints

### Invoice Operations

#### Create Invoice
```http
POST /issue
Content-Type: application/json

{
  "invoice_id": "INV-2025-0001"
}
```

**Response**:
```json
{
  "id": "nfe_io_invoice_id",
  "status": "waiting",
  "environment": "production",
  "flowStatus": "pending",
  "pdf": "https://api.nfe.io/.../pdf",
  "xml": "https://api.nfe.io/.../xml"
}
```

#### Get Invoice Status
```http
GET /invoice/:id
```

#### Get Invoice PDF
```http
GET /invoice/:id/pdf
```

#### Get Invoice XML
```http
GET /invoice/:id/xml
```

#### Cancel Invoice
```http
DELETE /invoice/:id
```

#### Create Correction Letter
```http
POST /invoice/:id/correction
Content-Type: application/json

{
  "correction": "Correction text here"
}
```

## 🏗️ Architecture

### Service Layer Pattern

```
┌─────────────────────────────────────────────┐
│            HTTP Handlers                     │
│         (invoice_hdl.go)                     │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│         Issuer Service                       │
│  • Orchestrates invoice creation            │
│  • Coordinates between services              │
│  • Maps data between systems                 │
└──────┬───────────┬──────────────┬───────────┘
       │           │              │
   ┌───▼────┐  ┌──▼──────┐  ┌───▼─────┐
   │  Tax   │  │ Builder │  │ Frappe  │
   │Service │  │ Service │  │  Repo   │
   └────────┘  └─────────┘  └─────────┘
                                 │
                            ┌────▼─────┐
                            │  NFe.io  │
                            │   Repo   │
                            └──────────┘
```

### Key Services

#### 1. **IssuerService** (`service/issuer.go`)
Main orchestrator that:
- Fetches invoice data from Frappe
- Retrieves tax templates and carrier info
- Maps Frappe data to NFe.io format
- Calculates taxes
- Sends to NFe.io
- Updates Frappe with response

#### 2. **TaxService** (`service/tax.go`)
Handles all Brazilian tax calculations:
- **ICMS**: State tax with various situations (00, 10, 20, 30, 40, 41, 50, 51, 60, 70, 90)
- **PIS**: Social Integration Program tax
- **COFINS**: Social Security Financing tax
- **IPI**: Industrialized Products Tax
- **DIFAL**: Interstate tax difference for final consumers

#### 3. **BuilderService** (`service/builder.go`)
Utility functions for:
- CFOP determination based on operation type and location
- Address building and formatting
- Transport information structuring

#### 4. **FrappeRepository** (`repository/frappe.go`)
Frappe API client with methods:
- `GetInvoice(id)` - Fetch invoice data
- `GetTax(name)` - Fetch tax template
- `GetCarrier(name)` - Fetch carrier information
- `UpdateInvoice(id, data)` - Update invoice with NFe.io response

#### 5. **NFeRepository** (`repository/nfeio.go`)
NFe.io API client with 10 methods:
- `CreateProductInvoice()` - Create new invoice
- `GetInvoice()` - Retrieve invoice details
- `GetInvoiceByAccessKey()` - Query by access key
- `DeleteInvoice()` - Cancel invoice
- `GetInvoicePDF()` - Download PDF
- `GetInvoiceXML()` - Download XML
- `CreateCorrectionLetter()` - Create correction
- `GetCorrectionLetterPDF()` - Download correction PDF
- `GetCorrectionLetterXML()` - Download correction XML

## 💾 Data Models

### Frappe Invoices DocType

```go
type Invoices struct {
    Name              string        // INV-YYYY-####
    OperationType     string        // Natureza de Operação
    ClientType        string        // PF or PJ
    ClientName        string        // Nome / Razão Social
    ClientIDNumber    string        // CPF / CNPJ
    ClientEmail       string
    TaxTemplate       string        // Link to Tax DocType
    Carrier           string        // Link to Carrier DocType
    InvoicesTable     []ItemInvoice // Child table with items
    
    // Delivery Address
    DeliveryAddress   string
    DeliveryCEP       string
    DeliveryState     string
    City              string
    
    // NFe.io Response (auto-filled)
    InvoiceID         string        // NFe.io invoice ID
    InvoiceLink       string        // PDF URL
    InvoiceSerie      string        // Serie number
    
    // Totals
    Total             float64       // Subtotal
    TotalTax          float64       // Total with taxes
}
```

### Tax Calculations

The service automatically calculates:

```
ICMS: Base Calculation × (ICMS Rate / 100)
PIS:  Base Calculation × (PIS Rate / 100)
COFINS: Base Calculation × (COFINS Rate / 100)
IPI:  Base Calculation × (IPI Rate / 100)

Where Base Calculation = Item Value × Quantity
```

## 🔌 Frappe Integration

### Installation

1. **Install the Frappe app**:
   ```bash
   cd ~/frappe-bench
   bench get-app https://github.com/AnyGridTech/frappe_brazil_invoice.git
   bench --site your-site install-app frappe_brazil_invoice
   bench migrate
   ```

2. **Configure site**:
   Add to `site_config.json`:
   ```json
   {
     "nfe_go_api_url": "http://localhost:3000"
   }
   ```

3. **Restart bench**:
   ```bash
   bench restart
   ```

### Usage in Frappe

1. Create a new **Invoice** document
2. Fill in client information (CPF/CNPJ, name, email)
3. Add delivery address (use CEP lookup for auto-fill)
4. Select **Tax Template**
5. Add items to the **invoices_table**
6. Select **Carrier** (optional)
7. **Submit** the document
8. Click **"Create NFe Invoice"** button
9. Wait for success message
10. Click **"View NFe PDF"** to see the generated invoice

### Form Buttons

- **Create NFe Invoice**: Sends invoice to NFe.io for processing
- **Check NFe Status**: Queries current status from NFe.io
- **View NFe PDF**: Opens generated PDF in new tab

## 📊 Tax Template Configuration

Create tax templates in Frappe with the following structure:

| Field | Description | Example |
|-------|-------------|---------|
| **ICMS** | | |
| CST ICMS | Tax situation code | 00, 10, 20, etc. |
| Aliq ICMS | ICMS rate (%) | 18.00 |
| Origin ICMS | Product origin | 0 (Nacional) |
| **IPI** | | |
| CST IPI | IPI situation | 50, 99 |
| Aliquota IPI | IPI rate (%) | 10.00 |
| **PIS** | | |
| CST PIS | PIS situation | 01, 02 |
| Aliquota PIS | PIS rate (%) | 1.65 |
| **COFINS** | | |
| CST COFINS | COFINS situation | 01, 02 |
| Aliquota COFINS | COFINS rate (%) | 7.60 |

## 🧪 Testing

### Unit Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/service
```

### Integration Testing

Use the provided Postman collection or test manually:

```bash
# Health check
curl http://localhost:3000/health

# Create invoice (replace with actual invoice ID from Frappe)
curl -X POST http://localhost:3000/issue \
  -H "Content-Type: application/json" \
  -d '{"invoice_id": "INV-2025-0001"}'

# Get invoice status
curl http://localhost:3000/invoice/{nfe_io_id}
```

## 🔐 Security

- API authentication using Frappe API keys
- NFe.io API key secured via environment variables
- HTTPS recommended for production
- Token-based auth middleware available

## 📝 Configuration Reference

### Environment Variables

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `FRAPPE_URL` | Yes | Frappe site URL | `https://mysite.erpnext.com` |
| `FRAPPE_API_KEY` | Yes | Frappe API key | `abc123...` |
| `FRAPPE_API_SECRET` | Yes | Frappe API secret | `def456...` |
| `NFEIO_API_KEY` | Yes | NFe.io API key | `xyz789...` |
| `COMPANY_ID` | Yes | NFe.io company ID | `123456` |
| `NFE_ENDPOINT` | Yes | NFe.io API endpoint | `https://api.nfe.io/v2/...` |
| `PORT` | No | Server port | `3000` (default) |
| `ENVIRONMENT` | No | Environment | `development` or `production` |

### CFOP Codes

The service automatically determines CFOP codes based on:
- Operation type (sale, return, transfer, etc.)
- Location (internal state vs interstate)
- Client type (PF vs PJ)

Common CFOPs:
- `5101/6101`: Sale to consumers
- `5102/6102`: Sale for resale
- `5405/6405`: Sale of fixed assets
- `5949/6949`: Other outbound operations

## 🐛 Troubleshooting

### Common Issues

**1. "Could not connect to Go API"**
- Check if service is running: `curl http://localhost:3000/health`
- Verify `nfe_go_api_url` in Frappe config
- Check firewall rules

**2. "Invalid tax number"**
- Ensure CPF has 11 digits (remove formatting)
- Ensure CNPJ has 14 digits (remove formatting)
- Use only numeric characters

**3. "Failed to get tax template"**
- Verify Tax template exists in Frappe
- Check template name spelling
- Ensure proper permissions

**4. NFe.io API errors**
- Check API key validity
- Verify company ID is correct
- Check NFe.io service status
- Review NFe.io documentation for error codes

### Debug Mode

Enable verbose logging:
```go
// In cmd/main.go
app.Use(logger.New(logger.Config{
    Format: "[${time}] ${status} - ${method} ${path}\n",
}))
```

## 📚 Documentation

- **[FRAPPE_INTEGRATION.md](FRAPPE_INTEGRATION.md)**: Detailed integration guide with architecture, data flow, and field mappings
- **[QUICK_START.md](QUICK_START.md)**: Step-by-step setup guide
- **[ADAPTATION_SUMMARY.md](ADAPTATION_SUMMARY.md)**: How legacy code was adapted
- **[NFEIO_ENHANCEMENT.md](NFEIO_ENHANCEMENT.md)**: NFe.io API enhancements

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👥 Authors

- **AnyGridTech Team** - Initial work

## 🙏 Acknowledgments

- NFe.io for the Brazilian electronic invoice API
- Frappe/ERPNext for the excellent ERP framework
- Go community for the amazing tools and libraries

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/AnyGridTech/frappe-nfe-bridge/issues)
- **Email**: support@anygridtech.com
- **Documentation**: [Wiki](https://github.com/AnyGridTech/frappe-nfe-bridge/wiki)

## 🗺️ Roadmap

- [ ] Add support for service invoices (NFSe)
- [ ] Implement batch invoice processing
- [ ] Add webhook support for NFe.io events
- [ ] Create admin dashboard
- [ ] Add invoice preview before sending
- [ ] Support for multiple companies
- [ ] Advanced tax scenario handling
- [ ] Integration with other ERP systems

---

**Made with ❤️ for the Brazilian market by AnyGridTech**
