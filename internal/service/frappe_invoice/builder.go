package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AnyGridTech/frappe-nfe-bridge/internal/models"
)

// BuilderService provides methods to build different parts of the NFe invoice
// Adapted from docs/invoice/build.go and docs/invoice/transport.go
type BuilderService struct{}

func NewBuilderService() *BuilderService {
	return &BuilderService{}
}

// BuilderInput contains all data needed to build an invoice
type BuilderInput struct {
	OperationNature string
	BuyerState      string
	IssuerState     string
	OperationType   string // "incoming" or "outgoing"
	Items           []models.Items
}

// DetermineCFOP determines the CFOP code based on operation nature and type
// Adapted from docs/invoice/build.go cfop methods
func (b *BuilderService) DetermineCFOP(operationNature, operationType, buyerState, issuerState string) (int, error) {
	operationNature = strings.ToLower(operationNature)
	isOutgoing := operationType == "outgoing"
	isSameState := buyerState == issuerState

	// CFOP mapping for common operations
	// Outgoing operations (sales)
	if isOutgoing {
		cfopMap := map[string]string{
			"retorno de remessa para conserto":       "x916",
			"retorno de troca em garantia":           "x949",
			"devolução de mercadoria de bonificação": "x949",
			"venda":                                "x102",
			"venda de produção do estabelecimento": "x101",
		}

		if cfopCode, exists := cfopMap[operationNature]; exists {
			return b.resolveCFOPCode(cfopCode, isSameState, isOutgoing), nil
		}

		// Default for sales
		return b.resolveCFOPCode("x102", isSameState, isOutgoing), nil
	}

	// Incoming operations (purchases/returns)
	cfopMap := map[string]string{
		"remessa para conserto": "x915",
		"troca em garantia":     "x949",
		"bonificação":           "x910",
		"compra":                "x102",
	}

	if cfopCode, exists := cfopMap[operationNature]; exists {
		return b.resolveCFOPCode(cfopCode, isSameState, isOutgoing), nil
	}

	return 0, fmt.Errorf("cfop not found for operation: %s", operationNature)
}

// resolveCFOPCode resolves the CFOP code based on state and operation type
// Adapted from docs/invoice/build.go verifyState methods
func (b *BuilderService) resolveCFOPCode(template string, sameState, isOutgoing bool) int {
	var prefix string

	if isOutgoing {
		// Outgoing operations: 5xxx (same state) or 6xxx (different state)
		if sameState {
			prefix = "5"
		} else {
			prefix = "6"
		}
	} else {
		// Incoming operations: 1xxx (same state) or 2xxx (different state)
		if sameState {
			prefix = "1"
		} else {
			prefix = "2"
		}
	}

	resolved := strings.Replace(template, "x", prefix, 1)
	cfop, _ := strconv.Atoi(resolved)
	return cfop
}

// DetermineOperationType sets the operation type and serie
// Adapted from docs/invoice/build.go OperationType methods
func (b *BuilderService) DetermineOperationType(operationType string) (string, int, error) {
	switch operationType {
	case "outgoing":
		return "outgoing", 11, nil
	case "incoming":
		return "incoming", 10, nil
	default:
		return "", 0, fmt.Errorf("invalid operation type: %s", operationType)
	}
}

// DetermineDestination determines if operation is internal or interstate
// Adapted from docs/invoice/build.go destination method
func (b *BuilderService) DetermineDestination(buyerState, issuerState string) (string, error) {
	if buyerState == issuerState {
		return "internal_Operation", nil
	}
	return "interstate_Operation", nil
}

// BuildTransport builds transport information
// Adapted from docs/invoice/transport.go
type TransportInput struct {
	CarrierName           string
	CarrierCNPJ           string
	CarrierStateTaxNumber string
	CarrierEmail          string
	CarrierPhone          string
	CarrierAddress        string
	CarrierCity           string
	CarrierCityCode       string
	CarrierState          string
	CarrierPostalCode     string
	CarrierNeighborhood   string
	CarrierNumber         string
	FreightModality       string // "ByIssuer", "ByBuyer", "ThirdParty", "OwnAccount", "NoFreight"
	SealNumber            string
	VolumeQuantity        int
	Species               string
	Brand                 string
	VolumeNumeration      string
	NetWeight             float64
	GrossWeight           float64
}

func (b *BuilderService) BuildTransport(input TransportInput) (models.Transport, error) {
	if input.CarrierCNPJ == "" {
		return models.Transport{}, fmt.Errorf("carrier CNPJ is required")
	}

	// Default freight modality
	if input.FreightModality == "" {
		input.FreightModality = "ByIssuer"
	}

	transport := models.Transport{
		FreightModality: input.FreightModality,
		SealNumber:      input.SealNumber,
		TransportGroup: models.TransportGroup{
			Name:             input.CarrierName,
			FederalTaxNumber: input.CarrierCNPJ,
			StateTaxNumber:   input.CarrierStateTaxNumber,
			Email:            input.CarrierEmail,
			Type:             legalEntity,
			Address: models.Address{
				Phone: input.CarrierPhone,
				State: input.CarrierState,
				City: models.City{
					Name: input.CarrierCity,
					Code: input.CarrierCityCode,
				},
				District:   input.CarrierNeighborhood,
				Street:     input.CarrierAddress,
				Number:     input.CarrierNumber,
				PostalCode: input.CarrierPostalCode,
				Country:    BRA,
			},
		},
		Volume: models.Volume{
			VolumeQuantity:   input.VolumeQuantity,
			Species:          input.Species,
			Brand:            input.Brand,
			VolumeNumeration: input.VolumeNumeration,
			NetWeight:        input.NetWeight,
			GrossWeight:      input.GrossWeight,
		},
	}

	return transport, nil
}

// CalculateVolumeWeight calculates total weight from items
// Adapted from docs/invoice/transport.go volume method
func (b *BuilderService) CalculateVolumeWeight(items []models.CustomFrappeItem) (netWeight, grossWeight float64) {
	for _, item := range items {
		// Assuming the item has weight fields from Frappe
		// Adjust according to your actual Frappe item structure
		netWeight += item.Qty * 1.0   // Replace with actual net weight field
		grossWeight += item.Qty * 1.2 // Replace with actual gross weight field
	}
	return
}

// BuildAddress builds address information
// Adapted from docs/invoice/build.go address method
type AddressInput struct {
	Street       string
	Number       string
	Neighborhood string
	City         string
	CityCode     string // IBGE code
	State        string
	PostalCode   string
	Phone        string
	Country      string
}

func (b *BuilderService) BuildAddress(input AddressInput) (models.Address, error) {
	if input.City == "" || input.CityCode == "" {
		return models.Address{}, fmt.Errorf("city and city code (IBGE) are required")
	}

	country := input.Country
	if country == "" {
		country = BRA
	}

	address := models.Address{
		Street:   input.Street,
		Number:   input.Number,
		District: input.Neighborhood,
		City: models.City{
			Name: input.City,
			Code: input.CityCode,
		},
		State:      input.State,
		PostalCode: input.PostalCode,
		Phone:      input.Phone,
		Country:    country,
	}

	return address, nil
}

// DetermineBuyerType determines if buyer is natural person or legal entity
// and sets appropriate tax regime
func (b *BuilderService) DetermineBuyerType(taxNumber string) (buyerType, taxRegime string, err error) {
	// Remove formatting
	cleanNumber := strings.ReplaceAll(taxNumber, ".", "")
	cleanNumber = strings.ReplaceAll(cleanNumber, "-", "")
	cleanNumber = strings.ReplaceAll(cleanNumber, "/", "")

	length := len(cleanNumber)

	switch length {
	case 11: // CPF
		return naturalPerson, none, nil
	case 14: // CNPJ
		return legalEntity, "", nil // Tax regime should be determined separately
	default:
		return "", "", fmt.Errorf("invalid tax number length: %d", length)
	}
}

// DetermineStateTaxNumberIndicator determines state tax number indicator for legal entities
// Adapted from docs/invoice/build.go buyer method for PJ
func (b *BuilderService) DetermineStateTaxNumberIndicator(icmsStatus string) (string, error) {
	switch icmsStatus {
	case "Contribuinte":
		return taxPayer, nil
	case "Não Contribuinte":
		return nonTaxPayer, nil
	case "Contribuinte Isento":
		return exempt, nil
	default:
		return "", fmt.Errorf("invalid ICMS status: %s", icmsStatus)
	}
}
