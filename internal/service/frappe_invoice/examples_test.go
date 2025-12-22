package service_test

import (
	"testing"

	"github.com/AnyGridTech/frappe-nfe-bridge/internal/models"
	"github.com/AnyGridTech/frappe-nfe-bridge/internal/service"
)

// Example: Using TaxService
func ExampleTaxService_CalculateTax() {
	taxService := service.NewTaxService()

	// Prepare tax input from your Frappe data
	input := service.TaxInput{
		ItemValue:       100.00, // Unit price
		Quantity:        2.0,    // Quantity
		FreightAmount:   10.00,  // Freight
		InsuranceAmount: 5.00,   // Insurance
		OthersAmount:    0.00,   // Other expenses
		DiscountAmount:  0.00,   // Discount

		// ICMS configuration
		AliqICMS:    18.0, // 18% ICMS
		OrigemICMS:  "0",  // National origin
		CSTICMS:     "00", // Normal taxation
		ModDetermBC: "3",  // Value of operation

		// PIS configuration
		AliquotaPis: 1.65, // 1.65% PIS
		CSTPis:      "01", // Normal taxation

		// COFINS configuration
		AliquotaCofins: 7.6,  // 7.6% COFINS
		CSTCofins:      "01", // Normal taxation

		// IPI configuration (if applicable)
		AliquotaIPI: 0.0,  // No IPI
		CSTIPI:      "50", // Outgoing - exempt
	}

	tax := taxService.CalculateTax(input)

	// Use the calculated tax in your NFe item
	_ = tax // Tax is ready to be used
}

// Example: Using BuilderService for CFOP
func ExampleBuilderService_DetermineCFOP() {
	builderService := service.NewBuilderService()

	// Example 1: Interstate sale
	cfop, err := builderService.DetermineCFOP(
		"venda",    // operation nature
		"outgoing", // operation type
		"RJ",       // buyer state
		"SP",       // issuer state
	)
	if err != nil {
		panic(err)
	}
	_ = cfop // Result: 6102 (interstate outgoing sale)

	// Example 2: Same state sale
	cfop, err = builderService.DetermineCFOP(
		"venda",
		"outgoing",
		"SP", // buyer in same state
		"SP", // issuer state
	)
	if err != nil {
		panic(err)
	}
	_ = cfop // Result: 5102 (same state outgoing sale)

	// Example 3: Return from warranty
	cfop, err = builderService.DetermineCFOP(
		"retorno de troca em garantia",
		"outgoing",
		"RJ",
		"SP",
	)
	if err != nil {
		panic(err)
	}
	_ = cfop // Result: 6949
}

// Example: Building Transport Information
func ExampleBuilderService_BuildTransport() {
	builderService := service.NewBuilderService()

	transportInput := service.TransportInput{
		CarrierName:           "Transportadora XYZ Ltda",
		CarrierCNPJ:           "12345678000190",
		CarrierStateTaxNumber: "123456789",
		CarrierEmail:          "contato@transportadora.com.br",
		CarrierPhone:          "11987654321",
		CarrierAddress:        "Rua das Flores",
		CarrierNumber:         "123",
		CarrierCity:           "São Paulo",
		CarrierCityCode:       "3550308", // IBGE code
		CarrierState:          "SP",
		CarrierPostalCode:     "01234567",
		CarrierNeighborhood:   "Centro",
		FreightModality:       "ByIssuer",
		SealNumber:            "SEAL123",
		VolumeQuantity:        1,
		Species:               "Caixa",
		Brand:                 "Marca XYZ",
		VolumeNumeration:      "VOL-001",
		NetWeight:             10.5,
		GrossWeight:           12.0,
	}

	transport, err := builderService.BuildTransport(transportInput)
	if err != nil {
		panic(err)
	}
	_ = transport // Use in NFe request
}

// Example: Building Address
func ExampleBuilderService_BuildAddress() {
	builderService := service.NewBuilderService()

	addressInput := service.AddressInput{
		Street:       "Rua das Palmeiras",
		Number:       "456",
		Neighborhood: "Jardim Paulista",
		City:         "São Paulo",
		CityCode:     "3550308", // IBGE code
		State:        "SP",
		PostalCode:   "01234567",
		Phone:        "11987654321",
		Country:      "BRA",
	}

	address, err := builderService.BuildAddress(addressInput)
	if err != nil {
		panic(err)
	}
	_ = address // Use for buyer or delivery address
}

// Example: Complete IssuerService Integration
func ExampleIssuerService_CompleteFlow() {
	// This would be in your actual application code

	// 1. Setup repositories (mocked here for example)
	// frappeRepo := repository.NewFrappeRepo(...)
	// nfeRepo := repository.NewNFeRepo(...)

	// 2. Create issuer service
	// issuerService := service.NewIssuerService(frappeRepo, nfeRepo, "companyID")

	// 3. Issue invoice
	// response, err := issuerService.IssueNoteForFrappeInvoice("INV-001")
	// if err != nil {
	//     log.Fatal(err)
	// }

	// 4. Use response
	// fmt.Printf("NFe ID: %s\n", response.ID)
	// fmt.Printf("Status: %s\n", response.Status)
	// fmt.Printf("PDF: %s\n", response.PdfUrl)
	// fmt.Printf("XML: %s\n", response.XmlUrl)
}

// Example: Determining Buyer Type
func ExampleBuilderService_DetermineBuyerType() {
	builderService := service.NewBuilderService()

	// CPF example
	buyerType, taxRegime, err := builderService.DetermineBuyerType("123.456.789-00")
	if err != nil {
		panic(err)
	}
	_ = buyerType // "naturalPerson"
	_ = taxRegime // "none"

	// CNPJ example
	buyerType, taxRegime, err = builderService.DetermineBuyerType("12.345.678/0001-90")
	if err != nil {
		panic(err)
	}
	_ = buyerType // "legalEntity"
	_ = taxRegime // "" (to be determined)
}

// Example: Tax Calculation with Item Data
func TestTaxCalculation(t *testing.T) {
	taxService := service.NewTaxService()

	// Simulating data from a Frappe item
	itemValue := 1000.00
	quantity := 2.0

	input := service.TaxInput{
		ItemValue:       itemValue,
		Quantity:        quantity,
		FreightAmount:   50.00,
		InsuranceAmount: 25.00,
		OthersAmount:    0.00,
		DiscountAmount:  100.00,
		AliqICMS:        18.0,
		OrigemICMS:      "0",
		CSTICMS:         "00",
		ModDetermBC:     "3",
		AliquotaPis:     1.65,
		CSTPis:          "01",
		AliquotaCofins:  7.6,
		CSTCofins:       "01",
		AliquotaIPI:     0.0,
		CSTIPI:          "50",
	}

	tax := taxService.CalculateTax(input)

	// Base Tax = (1000 * 2 + 50 + 25) - 100 = 1975
	expectedBaseTax := 1975.00

	if tax.Icms.BaseTax != expectedBaseTax {
		t.Errorf("Expected base tax %.2f, got %.2f", expectedBaseTax, tax.Icms.BaseTax)
	}

	// ICMS = 1975 * 18% = 355.50
	expectedICMS := 355.50
	if tax.Icms.Amount != expectedICMS {
		t.Errorf("Expected ICMS %.2f, got %.2f", expectedICMS, tax.Icms.Amount)
	}

	// PIS = 1975 * 1.65% = 32.59
	expectedPIS := 32.59
	if tax.Pis.Amount != expectedPIS {
		t.Errorf("Expected PIS %.2f, got %.2f", expectedPIS, tax.Pis.Amount)
	}

	// COFINS = 1975 * 7.6% = 150.10
	expectedCOFINS := 150.10
	if tax.Cofins.Amount != expectedCOFINS {
		t.Errorf("Expected COFINS %.2f, got %.2f", expectedCOFINS, tax.Cofins.Amount)
	}

	// Total Tax = 355.50 + 32.59 + 150.10 + 0 = 538.19
	expectedTotal := 538.19
	if tax.TotalTax != expectedTotal {
		t.Errorf("Expected total tax %.2f, got %.2f", expectedTotal, tax.TotalTax)
	}
}

// Example: DIFAL Calculation
func ExampleTaxService_CalculateDIFAL() {
	taxService := service.NewTaxService()

	operationValue := 1000.00
	aliqInterestadual := 0.12 // 12% interstate rate
	aliqInterna := 0.18       // 18% internal rate

	// Dual base method (more precise)
	difal := taxService.CalculateDIFAL(operationValue, aliqInterestadual, aliqInterna)
	_ = difal

	// Simple method
	difalSimple := taxService.CalculateDIFALSimple(operationValue, aliqInterestadual, aliqInterna)
	_ = difalSimple
}

// Example: State Tax Number Indicator
func ExampleBuilderService_DetermineStateTaxNumberIndicator() {
	builderService := service.NewBuilderService()

	// Contributor
	indicator, err := builderService.DetermineStateTaxNumberIndicator("Contribuinte")
	if err != nil {
		panic(err)
	}
	_ = indicator // "taxPayer"

	// Non-contributor
	indicator, err = builderService.DetermineStateTaxNumberIndicator("Não Contribuinte")
	if err != nil {
		panic(err)
	}
	_ = indicator // "nonTaxPayer"

	// Exempt
	indicator, err = builderService.DetermineStateTaxNumberIndicator("Contribuinte Isento")
	if err != nil {
		panic(err)
	}
	_ = indicator // "exempt"
}

// Example: Multiple Items with Different Tax Profiles
func ExampleMultipleItemsWithTaxes() {
	taxService := service.NewTaxService()

	// Item 1: Normal taxation
	item1Tax := taxService.CalculateTax(service.TaxInput{
		ItemValue:      500.00,
		Quantity:       1.0,
		AliqICMS:       18.0,
		AliquotaPis:    1.65,
		AliquotaCofins: 7.6,
		AliquotaIPI:    0.0,
		OrigemICMS:     "0",
		CSTICMS:        "00",
		CSTPis:         "01",
		CSTCofins:      "01",
		CSTIPI:         "50",
		ModDetermBC:    "3",
	})

	// Item 2: With IPI
	item2Tax := taxService.CalculateTax(service.TaxInput{
		ItemValue:      300.00,
		Quantity:       2.0,
		AliqICMS:       18.0,
		AliquotaPis:    1.65,
		AliquotaCofins: 7.6,
		AliquotaIPI:    10.0, // 10% IPI
		OrigemICMS:     "0",
		CSTICMS:        "00",
		CSTPis:         "01",
		CSTCofins:      "01",
		CSTIPI:         "00", // Normal taxation
		ModDetermBC:    "3",
	})

	items := []models.Items{
		{
			Code:        "1",
			Description: "Product 1",
			Quantity:    1,
			UnitAmount:  500.00,
			TotalAmount: 500.00,
			Tax:         item1Tax,
		},
		{
			Code:        "2",
			Description: "Product 2",
			Quantity:    2,
			UnitAmount:  300.00,
			TotalAmount: 600.00,
			Tax:         item2Tax,
		},
	}

	_ = items // Use in ProductInvoiceRequest
}
