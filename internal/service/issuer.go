package service

import (
    "nfe-go/internal/models"
    "nfe-go/internal/repository"
)

type IssuerService interface {
    IssueNoteForFrappeInvoice(invoiceID string) (*models.ProductInvoiceResponse, error)
}

type issuerService struct {
    frappeRepo repository.FrappeRepository
    nfeRepo    repository.NFeRepository
    companyID  string // Your Company ID in NFe.io
}

func NewIssuerService(f repository.FrappeRepository, n repository.NFeRepository, companyID string) IssuerService {
    return &issuerService{
        frappeRepo: f,
        nfeRepo:    n,
        companyID:  companyID,
    }
}

func (s *issuerService) IssueNoteForFrappeInvoice(invoiceID string) (*models.ProductInvoiceResponse, error) {
    // 1. Get Data from ERPNext
    frappeInv, err := s.frappeRepo.GetCustomInvoice(invoiceID)
    if err != nil {
        return nil, err
    }

    // 2. Map ERPNext -> NFe.io
    nfePayload := s.mapFrappeToNFe(frappeInv)

    // 3. Send to NFe.io
    response, err := s.nfeRepo.CreateProductInvoice(nfePayload)
    if err != nil {
        return nil, err
    }

    // 4. (Optional) Update ERPNext with the PDF URL
    // s.frappeRepo.UpdateInvoiceStatus(invoiceID, response.PdfUrl)

    return response, nil
}

// Mapper Function (Pure Logic)
func (s *issuerService) mapFrappeToNFe(inv *models.CustomFrappeInvoice) *models.ProductInvoiceRequest {
    var items []models.Item
    for _, item := range inv.Items {
        items = append(items, models.Item{
            Code:        item.ItemCode,
            Description: item.ItemName,
            Quantity:    item.Qty,
            UnitValue:   item.Rate,
            TotalValue:  item.Amount,
            NCM:         item.NCM, // Critical for Product Invoices
            CFOP:        "5102",   // Example: Sale of merchandise
        })
    }

    return &models.ProductInvoiceRequest{
        CompanyID: s.companyID,
        Reference: inv.Name,
        Items:     items,
        Buyer: models.Buyer{
            Name:             inv.CustomerName,
            FederalTaxNumber: inv.CNPJ,
            // Address mapping would go here
        },
    }
}