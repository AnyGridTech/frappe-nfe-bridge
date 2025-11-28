package service

import (
	"math"

	"github.com/AnyGridTech/frappe-nfe-bridge/internal/models"
)

// TaxService provides tax calculation methods for NFe items
type TaxService struct{}

func NewTaxService() *TaxService {
	return &TaxService{}
}

// Base de cálculo = (Valor do produto + Frete + Seguro + Outras Despesas Acessórias)
// DIFAL = Valor da Operação * (Alíquota interna – Alíquota interestadual) (BASE ÚNICA)
/*DIFAL (BASE DUPLA)
ICMS Interestadual = Valor da Operação * Alíquota Interestadual
Base de Cálculo 1 = Valor da Operação – Valor ICMS Interestadual
Base de Cálculo 2 = Base de Cálculo 1 / (1 – Alíquota Interna)
ICMS Interno = Base de Cálculo 2 * Alíquota Interna
DIFAL = ICMS Interno – ICMS Interestadual
*/

// TaxInput represents the input data for tax calculations
type TaxInput struct {
	ItemValue       float64
	Quantity        float64
	FreightAmount   float64
	InsuranceAmount float64
	OthersAmount    float64
	DiscountAmount  float64

	// Tax rates and codes from Frappe
	AliqICMS        float64
	AliqICMSDestino float64
	OrigemICMS      string
	CSTICMS         string
	ModDetermBC     string

	AliquotaPis float64
	CSTPis      string

	AliquotaCofins float64
	CSTCofins      string

	AliquotaIPI float64
	CSTIPI      string
}

// CalculateTax calculates all taxes for an item and returns a Tax struct
func (s *TaxService) CalculateTax(input TaxInput) models.Tax {
	baseValue := input.ItemValue * input.Quantity
	baseTax := s.calculateBaseTax(baseValue, input.FreightAmount, input.InsuranceAmount, input.OthersAmount, input.DiscountAmount)

	icms := s.calculateICMS(input, baseTax)
	pis := s.calculatePIS(input, baseTax)
	cofins := s.calculateCOFINS(input, baseTax)
	ipi := s.calculateIPI(input, baseTax)

	totalTax := icms.Amount + pis.Amount + cofins.Amount + ipi.Amount

	return models.Tax{
		TotalTax: totalTax,
		Icms:     icms,
		Pis:      pis,
		Cofins:   cofins,
		Ipi:      ipi,
	}
}

// calculateBaseTax calculates the base value for tax calculation
// Base = (Product Value + Freight + Insurance + Others) - Discount
func (s *TaxService) calculateBaseTax(productValue, freight, insurance, others, discount float64) float64 {
	return productValue + freight + insurance + others - discount
}

// calculateICMS calculates ICMS tax
func (s *TaxService) calculateICMS(input TaxInput, baseTax float64) models.Icms {
	amount := s.calcSimpleTax(input.AliqICMS, baseTax)

	origin := "0"
	if len(input.OrigemICMS) > 0 {
		origin = input.OrigemICMS[:1]
	}

	cst := "00"
	if len(input.CSTICMS) >= 2 {
		cst = input.CSTICMS[:2]
	}

	modality := "3"
	if len(input.ModDetermBC) > 0 {
		modality = input.ModDetermBC[:1]
	}

	return models.Icms{
		Origin:             origin,         // Origem da mercadoria
		Cst:                cst,            // Código de Situação Tributária
		BaseTaxModality:    modality,       // Modalidade de determinação da BC
		BaseTax:            baseTax,        // Valor da BC do ICMS (vBC)
		Rate:               input.AliqICMS, // pICMS Alíquota do imposto (pICMS)
		Amount:             amount,         // Valor do ICMS (vICMS)
		BaseTaxSTReduction: "0",
	}
}

// calculatePIS calculates PIS tax
func (s *TaxService) calculatePIS(input TaxInput, baseTax float64) models.Pis {
	amount := s.calcSimpleTax(input.AliquotaPis, baseTax)

	cst := "01"
	if len(input.CSTPis) >= 2 {
		cst = input.CSTPis[:2]
	}

	return models.Pis{
		Amount:  amount,            // Valor do PIS (vPIS)
		Rate:    input.AliquotaPis, // Alíquota do PIS (em percentual) (pPIS)
		BaseTax: baseTax,           // Valor da Base de Cálculo do PIS (vBC)
		Cst:     cst,               // Código de Situação Tributária do PIS (CST)
	}
}

// calculateCOFINS calculates COFINS tax
func (s *TaxService) calculateCOFINS(input TaxInput, baseTax float64) models.Cofins {
	amount := s.calcSimpleTax(input.AliquotaCofins, baseTax)

	cst := "01"
	if len(input.CSTCofins) >= 2 {
		cst = input.CSTCofins[:2]
	}

	return models.Cofins{
		Amount:  amount,               // Valor do Cofins (vCofins)
		Rate:    input.AliquotaCofins, // Alíquota do Cofins (em percentual) (pCofins)
		BaseTax: baseTax,              // Valor da Base de Cálculo do Cofins (vBC)
		Cst:     cst,                  // Código de Situação Tributária do Cofins (CST)
	}
}

// calculateIPI calculates IPI tax
func (s *TaxService) calculateIPI(input TaxInput, baseTax float64) models.Ipi {
	amount := s.calcSimpleTax(input.AliquotaIPI, baseTax)

	cst := "50"
	if len(input.CSTIPI) >= 2 {
		cst = input.CSTIPI[:2]
	}

	return models.Ipi{
		Amount: amount,
		Rate:   input.AliquotaIPI, // Alíquota do IPI (em percentual) (pIPI)
		Cst:    cst,               // Código de Situação Tributária do IPI (CST)
	}
}

// CalculateDIFAL calculates DIFAL (Diferencial de Alíquota) - dual base method
// Used for interstate operations to non-taxpayers
func (s *TaxService) CalculateDIFAL(operationValue, aliqInterestadual, aliqInterna float64) float64 {
	icmsInter := operationValue * aliqInterestadual
	bc1 := operationValue - icmsInter
	bc2 := bc1 / (1 - aliqInterna)
	icmsIntra := bc2 * aliqInterna
	return icmsIntra - icmsInter
}

// CalculateDIFALSimple calculates DIFAL using simple method (single base)
func (s *TaxService) CalculateDIFALSimple(operationValue, aliqInterestadual, aliqInterna float64) float64 {
	return operationValue * (aliqInterna - aliqInterestadual)
}

// calcSimpleTax calculates tax amount based on rate and base value
// Returns rounded value to 2 decimal places
func (s *TaxService) calcSimpleTax(aliq, base float64) float64 {
	return math.Round((aliq/100)*base*100) / 100
}

// CalculateTotalTax calculates the sum of all taxes
func (s *TaxService) CalculateTotalTax(tax models.Tax) float64 {
	return tax.Icms.Amount + tax.Pis.Amount + tax.Cofins.Amount + tax.Ipi.Amount
}
