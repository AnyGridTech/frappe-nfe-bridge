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
	frappeInv, err := s.frappeRepo.GetInvoice(invoiceID)
	if err != nil {
		return nil, err
	}

	// 2. Get Tax Template if specified
	var taxTemplate *models.FrappeTax
	if frappeInv.TaxTemplate != "" {
		taxTemplate, err = s.frappeRepo.GetTax(frappeInv.TaxTemplate)
		if err != nil {
			return nil, fmt.Errorf("failed to get tax template: %v", err)
		}
	}

	// 3. Get Carrier if specified
	var carrier *models.Carrier
	if frappeInv.Carrier != "" {
		carrier, err = s.frappeRepo.GetCarrier(frappeInv.Carrier)
		if err != nil {
			return nil, fmt.Errorf("failed to get carrier: %v", err)
		}
	}

	// 4. Map ERPNext -> NFe.io
	nfePayload, err := s.mapFrappeToNFe(frappeInv, taxTemplate, carrier)
	if err != nil {
		return nil, err
	}

	// 5. Send to NFe.io
	response, err := s.nfeRepo.CreateProductInvoice(nfePayload)
	if err != nil {
		return nil, err
	}

	// 6. Update Frappe invoice with NFe.io response
	updateData := map[string]interface{}{
		"invoice_id":    response.ID,
		"invoice_serie": nfePayload.Serie,
		"invoice_link":  response.PdfUrl,
	}

	err = s.frappeRepo.UpdateInvoice(invoiceID, updateData)
	if err != nil {
		// Log error but don't fail - invoice was created successfully
		fmt.Printf("Warning: Failed to update Frappe invoice: %v\n", err)
	}

	return response, nil
} // Mapper Function (Pure Logic) - Adapted from docs/invoice/index.go
func (s *issuerService) mapFrappeToNFe(inv *models.Invoices, taxTemplate *models.FrappeTax, carrier *models.Carrier) (*models.ProductInvoiceRequest, error) {
	var items []models.Items

	// Map items with tax calculations
	for i, item := range inv.InvoicesTable {
		// Use tax template values if available, otherwise use item-specific values
		var taxInput TaxInput
		if taxTemplate != nil {
			taxInput = TaxInput{
				ItemValue:      item.Rate,
				Quantity:       float64(item.Quantity),
				AliqICMS:       taxTemplate.AliqICMS,
				CSTICMS:        taxTemplate.CSTICMS,
				OrigemICMS:     taxTemplate.OriginICMS,
				ModDetermBC:    taxTemplate.ModDetermBC,
				AliquotaPis:    taxTemplate.AliquotaPis,
				CSTPis:         taxTemplate.CSTPis,
				AliquotaCofins: taxTemplate.AliquotaCofins,
				CSTCofins:      taxTemplate.CSTCofins,
				AliquotaIPI:    taxTemplate.AliquotaIPI,
				CSTIPI:         taxTemplate.CSTIPI,
			}
		} else {
			// Parse rates from item if no template
			icmsRate, _ := strconv.ParseFloat(item.ICMSRate, 64)
			pisRate, _ := strconv.ParseFloat(item.PISRate, 64)
			cofinsRate, _ := strconv.ParseFloat(item.COFINSRate, 64)
			ipiRate, _ := strconv.ParseFloat(item.IPIRate, 64)

			taxInput = TaxInput{
				ItemValue:      item.Rate,
				Quantity:       float64(item.Quantity),
				AliqICMS:       icmsRate,
				CSTICMS:        "00",
				OrigemICMS:     "0",
				ModDetermBC:    "3",
				AliquotaPis:    pisRate,
				CSTPis:         "01",
				AliquotaCofins: cofinsRate,
				CSTCofins:      "01",
				AliquotaIPI:    ipiRate,
				CSTIPI:         "50",
			}
		}

		calculatedTax := s.taxService.CalculateTax(taxInput)

		nfeItem := models.Items{
			Code:        strconv.Itoa(i + 1),
			CodeGTIN:    "SEM GTIN",
			CodeTaxGTIN: "SEM GTIN",
			Description: item.ItemName,
			Ncm:         item.NCM,
			Cfop:        5102, // Default, should be determined based on operation
			Unit:        "UN",
			Quantity:    item.Quantity,
			UnitAmount:  item.Rate,
			TotalAmount: item.Rate * float64(item.Quantity),
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

	// Build transport if carrier is provided
	transport := s.buildTransport(inv, carrier)

	// Determine destination (internal or interstate)
	destination := s.determineDestination(inv)

	// Determine operation type based on Frappe data
	operationType := "outgoing"
	if inv.OperationType != "" {
		// Aqui tem que dar um erro caso o OperationType seja vazio ou invalido
		// Somente outgoing & ongoing que sao validos.
		// Map Frappe operation type to NFe.io format if needed
		operationType = "outgoing"
	}

	payload := &models.ProductInvoiceRequest{
		Serie:         1, // Configure based on operation type
		OperationType: operationType,
		ConsumerType:  s.determineConsumerType(inv),
		PurposeType:   "normal",
		Destination:   destination,
		Buyer:         *buyer,
		Items:         items,
		Payment:       payment,
		Transport:     transport,
	}

	// Add additional information if present
	if inv.AdditionalInformation != "" {
		payload.AdditionalInformation = &models.AdditionalInformation{
			Taxpayer: inv.AdditionalInformation,
		}
	}

	return payload, nil
}

// buildBuyer creates buyer information from Frappe invoice
// Adapted from docs/invoice/build.go buyer methods
func (s *issuerService) buildBuyer(inv *models.Invoices) (*models.Buyer, error) {
	buyer := &models.Buyer{
		Name: inv.ClientName,
	}

	// Determine if it's CPF (natural person) or CNPJ (legal entity)
	cleanTaxNumber := s.cleanTaxNumber(inv.ClientIDNumber)
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
		// Map Frappe contribuinte_icms to NFe.io format
		if inv.ContribuinteIcms == "Sim" || inv.ContribuinteIcms == "1" {
			buyer.StateTaxNumberIndicator = taxPayer
			if inv.InscricaoEstadual != "" {
				buyer.StateTaxNumber = s.cleanTaxNumber(inv.InscricaoEstadual)
			}
		} else if inv.ContribuinteIcms == "Isento" || inv.ContribuinteIcms == "2" {
			buyer.StateTaxNumberIndicator = exempt
		} else {
			buyer.StateTaxNumberIndicator = nonTaxPayer
		}
	} else {
		return nil, fmt.Errorf("invalid tax number length: %d", len(cleanTaxNumber))
	}

	// Build address from delivery information
	if inv.DeliveryAddress != "" {
		buyer.Address = s.buildAddress(inv)
	}

	// Add email if present
	if inv.ClientEmail != "" {
		buyer.Email = inv.ClientEmail
	}

	return buyer, nil
}

// buildAddress creates address from Frappe delivery information
func (s *issuerService) buildAddress(inv *models.Invoices) models.Address {
	return models.Address{
		Country:    BRA,
		PostalCode: s.cleanTaxNumber(inv.DeliveryCEP),
		Street:     inv.DeliveryAddress,
		Number:     inv.DeliveryNumberAddress,
		District:   inv.DeliveryNeighborhood,
		City: models.City{
			Name: inv.City,
			Code: inv.DeliveryIBGE,
		},
		State:                 inv.DeliveryState,
		AdditionalInformation: inv.DeliveryComplement,
		Phone:                 s.formatPhoneNumber(inv.DeliveryPhone),
	}
}

// buildTransport creates transport information from carrier data
func (s *issuerService) buildTransport(inv *models.Invoices, carrier *models.Carrier) models.Transport {
	transport := models.Transport{
		FreightModality: s.determineShippingModality(inv.FreightModality),
	}

	// Add carrier information if present
	if carrier != nil {
		cleanCNPJ := s.cleanTaxNumber(carrier.CNPJ)
		cnpjNum, _ := strconv.Atoi(cleanCNPJ)

		transport.TransportGroup = models.TransportGroup{
			Name:             carrier.CarrierName,
			FederalTaxNumber: cnpjNum,
			StateTaxNumber:   carrier.StateRegistration,
			Email:            carrier.Email,
			Type:             legalEntity,
			Address: models.Address{
				Country:    BRA,
				PostalCode: s.cleanTaxNumber(carrier.CEP),
				Street:     carrier.Address,
				District:   carrier.Neighborhood,
				City: models.City{
					Name: carrier.City,
					Code: carrier.IBGE,
				},
				State: carrier.State,
				Phone: s.formatPhoneNumber(carrier.Phone),
			},
		}
	}

	// Add volume information if present
	if inv.ProductQuantity != "" {
		qty, _ := strconv.Atoi(inv.ProductQuantity)
		grossWeight, _ := strconv.ParseFloat(inv.ProductGrossWeight, 64)
		netWeight, _ := strconv.ParseFloat(inv.ProductNetWeight, 64)

		transport.Volume = models.Volume{
			VolumeQuantity: qty,
			Species:        inv.ProductType,
			Brand:          inv.ProductBrand,
			GrossWeight:    grossWeight,
			NetWeight:      netWeight,
		}
	}

	return transport
}

// determineShippingModality maps Frappe freight modality to NFe.io format
func (s *issuerService) determineShippingModality(modality string) string {
	switch modality {
	case "0", "Por conta do emitente":
		return "byIssuer"
	case "1", "Por conta do destinatário":
		return "byReceiver"
	case "2", "Por conta de terceiros":
		return "byThirdParty"
	case "9", "Sem frete":
		return "noShipping"
	default:
		return "noShipping"
	}
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
func (s *issuerService) determineDestination(inv *models.Invoices) string {
	// This should compare buyer's state with issuer's state
	// For now, returning a default value
	return "interstate_Operation"
}

// determineConsumerType determines if it's final consumer or normal
// Adapted from docs/invoice/build.go consumerType methods
func (s *issuerService) determineConsumerType(inv *models.Invoices) string {
	cleanTaxNumber := s.cleanTaxNumber(inv.ClientIDNumber)

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
