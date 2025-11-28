package service

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/AnyGridTech/frappe-nfe-bridge/internal/models"
	"github.com/AnyGridTech/frappe-nfe-bridge/internal/repository"
)

const (
	naturalPerson = "naturalPerson"
	legalEntity   = "legalEntity"
	none          = "none"
	taxPayer      = "taxPayer"
	nonTaxPayer   = "nonTaxPayer"
	exempt        = "exempt"
	BRA           = "BRA"
)

type IssuerService interface {
	IssueNoteForFrappeInvoice(invoiceID string) (*models.ProductInvoiceResponse, error)
}

type issuerService struct {
	frappeRepo repository.FrappeRepository
	nfeRepo    repository.NFeRepository
	taxService *TaxService
	companyID  string // Your Company ID in NFe.io
}

func NewIssuerService(f repository.FrappeRepository, n repository.NFeRepository, companyID string) IssuerService {
	return &issuerService{
		frappeRepo: f,
		nfeRepo:    n,
		taxService: NewTaxService(),
		companyID:  companyID,
	}
}

func (s *issuerService) IssueNoteForFrappeInvoice(invoiceID string) (*models.ProductInvoiceResponse, error) {
	// 1. Get Data from ERPNext/Frappe
	frappeInv, err := s.frappeRepo.GetCustomInvoice(invoiceID)
	if err != nil {
		return nil, err
	}

	// 2. Map ERPNext -> NFe.io
	nfePayload, err := s.mapFrappeToNFe(frappeInv)
	if err != nil {
		return nil, err
	}

	// 3. Send to NFe.io
	response, err := s.nfeRepo.CreateProductInvoice(nfePayload)
	if err != nil {
		return nil, err
	}

	// 4. (Optional) Update ERPNext with the PDF URL
	// s.frappeRepo.UpdateInvoiceStatus(invoiceID, response.PdfUrl)

	return response, nil
}

// Mapper Function (Pure Logic) - Adapted from docs/invoice/index.go
func (s *issuerService) mapFrappeToNFe(inv *models.CustomFrappeInvoice) (*models.ProductInvoiceRequest, error) {
	var items []models.Items

	// Map items with tax calculations
	for i, item := range inv.Items {
		// Calculate taxes for this item
		taxInput := TaxInput{
			ItemValue:      item.Rate,
			Quantity:       item.Qty,
			AliqICMS:       0, // Should come from Frappe item or be configured
			CSTICMS:        "00",
			OrigemICMS:     "0",
			ModDetermBC:    "3",
			AliquotaPis:    0,
			CSTPis:         "01",
			AliquotaCofins: 0,
			CSTCofins:      "01",
			AliquotaIPI:    0,
			CSTIPI:         "50",
		}

		calculatedTax := s.taxService.CalculateTax(taxInput)

		nfeItem := models.Items{
			Code:        strconv.Itoa(i + 1),
			CodeGTIN:    "SEM GTIN",
			CodeTaxGTIN: "SEM GTIN",
			Description: item.ItemName,
			Ncm:         item.NCM,
			Cfop:        s.determineCFOP(item.CFOP),
			Unit:        "UN",
			Quantity:    int(item.Qty),
			UnitAmount:  item.Rate,
			TotalAmount: item.Amount,
			Cest:        item.CEST,
			Tax:         calculatedTax,
		}

		items = append(items, nfeItem)
	}

	// Build buyer information
	buyer, err := s.buildBuyer(inv)
	if err != nil {
		return nil, err
	}

	// Build payment
	payment := s.buildPayment()

	// Determine destination (internal or interstate)
	destination := s.determineDestination(inv)

	payload := &models.ProductInvoiceRequest{
		Serie:         1, // Configure based on operation type
		OperationType: "outgoing",
		ConsumerType:  s.determineConsumerType(inv),
		PurposeType:   "normal",
		Destination:   destination,
		Buyer:         *buyer,
		Items:         items,
		Payment:       payment,
		Transport:     models.Transport{},
	}

	return payload, nil
}

// buildBuyer creates buyer information from Frappe invoice
// Adapted from docs/invoice/build.go buyer methods
func (s *issuerService) buildBuyer(inv *models.CustomFrappeInvoice) (*models.Buyer, error) {
	buyer := &models.Buyer{
		Name: inv.CustomerName,
	}

	// Determine if it's CPF (natural person) or CNPJ (legal entity)
	cleanTaxNumber := s.cleanTaxNumber(inv.CNPJ)
	taxNumber, err := strconv.Atoi(cleanTaxNumber)
	if err != nil {
		return nil, fmt.Errorf("invalid tax number: %v", err)
	}

	buyer.FederalTaxNumber = taxNumber

	// CPF has 11 digits, CNPJ has 14 digits
	if len(cleanTaxNumber) == 11 {
		buyer.Type = naturalPerson
		buyer.TaxRegime = none
	} else if len(cleanTaxNumber) == 14 {
		buyer.Type = legalEntity
		// State tax number indicator should come from Frappe
		// buyer.StateTaxNumberIndicator = taxPayer // or nonTaxPayer or exempt
	} else {
		return nil, fmt.Errorf("invalid tax number length: %d", len(cleanTaxNumber))
	}

	// Address should be built from Frappe data
	// buyer.Address = s.buildAddress(inv)

	return buyer, nil
}

// buildPayment creates payment information
// Adapted from docs/invoice/build.go payment method
func (s *issuerService) buildPayment() []models.Payment {
	return []models.Payment{
		{
			PaymentDetail: []models.PaymentDetail{
				{
					Method: "withoutPayment",
					Amount: 0,
				},
			},
		},
	}
}

// determineDestination determines if it's internal or interstate operation
// Adapted from docs/invoice/build.go destination method
func (s *issuerService) determineDestination(inv *models.CustomFrappeInvoice) string {
	// This should compare buyer's state with issuer's state
	// For now, returning a default value
	return "interstate_Operation"
}

// determineConsumerType determines if it's final consumer or normal
// Adapted from docs/invoice/build.go consumerType methods
func (s *issuerService) determineConsumerType(inv *models.CustomFrappeInvoice) string {
	cleanTaxNumber := s.cleanTaxNumber(inv.CNPJ)

	// CPF (11 digits) = final consumer
	if len(cleanTaxNumber) == 11 {
		return "finalConsumer"
	}
	// CNPJ (14 digits) = normal
	return "normal"
}

// determineCFOP converts CFOP from string to int
func (s *issuerService) determineCFOP(cfop string) int {
	// Remove any non-numeric characters
	cleanCFOP := regexp.MustCompile(`\D`).ReplaceAllString(cfop, "")
	cfopInt, err := strconv.Atoi(cleanCFOP)
	if err != nil {
		return 5102 // Default CFOP for outgoing sales
	}
	return cfopInt
}

// cleanTaxNumber removes formatting from CPF/CNPJ
// Adapted from docs/invoice/index.go formatPhoneNumber concept
func (s *issuerService) cleanTaxNumber(taxNumber string) string {
	// Remove all non-numeric characters
	re := regexp.MustCompile(`\D`)
	cleaned := re.ReplaceAllString(taxNumber, "")
	return strings.TrimPrefix(cleaned, "+")
}

// formatPhoneNumber removes formatting from phone numbers
func (s *issuerService) formatPhoneNumber(phone string) string {
	re := regexp.MustCompile(`\D`)
	formatted := re.ReplaceAllString(phone, "")
	return strings.TrimPrefix(formatted, "+")
}
