//go:build ignore
// +build ignore

package invoice

import (
	"erpnext-go/pkg/models"
	"fmt"
	"log"
	"strconv"
	"strings"
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

type Params struct {
	Op    OperationNatureIE
	Items ItemsIE
	Buyer BuyerIE
}
type OperationNatureIE interface {
	cfop(items *[]models.Items) error
	OperationType(invoice *models.ProductInvoice) error
	ItemsIE
}
type ItemsIE interface {
	items(items *[]models.Items, product *[]models.Product) error
}
type BuyerIE interface {
	buyer(buyer *models.Buyer, r *models.Req) error
	consumerType(invoice *models.ProductInvoice)
}

type Ongoing struct {
	msg string
}
type Outgoing struct {
	msg string
}
type PF struct{}
type PJ struct{}
type Sale struct{}

func (i PF) buyer(buyer *models.Buyer, r *models.Req) error {
	return setBuyerDetails(buyer, r.BuyerCPF, r.BuyerName, naturalPerson, none)
}

func (i PJ) buyer(buyer *models.Buyer, r *models.Req) error {
	if err := setBuyerDetails(buyer, r.BuyerCNPJ, r.BuyerName, legalEntity, ""); err != nil {
		return err
	}
	if err := address(buyer, r); err != nil {
		return err
	}
	buyer.StateTaxNumber = r.BuyerStateTaxNumber
	switch r.BuyerICMS {
	case "Contribuinte":
		buyer.StateTaxNumberIndicator = taxPayer
	case "Não Contribuinte":
		buyer.StateTaxNumberIndicator = nonTaxPayer
	case "Contribuinte Isento":
		buyer.StateTaxNumberIndicator = exempt
	default:
		return fmt.Errorf("buyer ICMS is not under acceptable values! buyerType: %s", r.BuyerType)
	}
	return nil
}

func (i PF) consumerType(invoice *models.ProductInvoice) {
	invoice.ConsumerType = "finalConsumer"
}
func (i PJ) consumerType(invoice *models.ProductInvoice) {
	invoice.ConsumerType = "normal"
}
func setBuyerDetails(buyer *models.Buyer, taxNumberStr, name, buyerType, taxRegime string) error {
	taxNumber, err := strconv.Atoi(taxNumberStr)
	if err != nil {
		return err
	}
	buyer.FederalTaxNumber = taxNumber
	buyer.Name = name
	buyer.Type = buyerType
	buyer.TaxRegime = taxRegime
	return nil
}

func address(buyer *models.Buyer, r *models.Req) error {
	var address models.Address
	address.City = models.City{Code: r.BuyerIBGE, Name: r.BuyerCity}
	if address.City.Code == "" || address.City.Name == "" {
		return fmt.Errorf("IBGE ou cidade não encontrada, Cidade: %s, %s, IBGE: %s", r.BuyerCity, r.BuyerState, r.BuyerIBGE)
	}
	address.Phone = formatPhoneNumber(r.DeliveryPhone)
	address.Country = "BRA"
	address.State = r.BuyerState
	address.District = r.BuyerNeighborhood
	address.Street = r.BuyerAddress
	address.PostalCode = r.BuyerCEP
	address.Number = r.BuyerNumber
	buyer.Address = address
	return nil
}

// Entrada: 0 | saida: 1
func (i Ongoing) cfop(items *[]models.Items) error {
	return setCfop(items, i.msg, false, map[string]string{
		"remessa para conserto": "x915",
		"troca em garantia":     "x949",
		"bonificação":           "x910",
	})
}

func (i Outgoing) cfop(items *[]models.Items) error {
	return setCfop(items, i.msg, true, map[string]string{
		"retorno de remessa para conserto":       "x916",
		"retorno de troca em garantia":           "x949",
		"devolução de mercadoria de bonificação": "x949",
	})
}

func (i Outgoing) OperationType(invoice *models.ProductInvoice) error {
	if invoice == nil {
		return fmt.Errorf("invoice is nil")
	}
	invoice.OperationType = "outgoing"
	invoice.Serie = 11
	return nil
}

func (i Ongoing) OperationType(invoice *models.ProductInvoice) error {
	if invoice == nil {
		return fmt.Errorf("invoice is nil")
	}
	invoice.OperationType = "incoming"
	invoice.Serie = 10
	return nil
}

func setCfop(items *[]models.Items, msg string, state bool, cfopMap map[string]string) error {
	msg = strings.ToLower(msg)
	log.Println(msg)
	vr := cfopMap[msg]
	if vr == "" {
		return fmt.Errorf("cfop not found")
	}
	vr = Ongoing{}.verifyState(vr, state)
	cfop, err := strconv.Atoi(vr)
	if err != nil {
		return err
	}
	(*items)[0].Cfop = cfop
	return nil
}

func (i Ongoing) items(items *[]models.Items, products *[]models.Product) error {
	var itemList []models.Items
	code := 1
	for _, product := range *products {
		item := models.Items{
			Code:        strconv.Itoa(code),
			CodeGTIN:    "SEM GTIN",
			CodeTaxGTIN: "SEM GTIN",
			Quantity:    1,
			Unit:        "UN",
			Cest:        "",
			Description: product.Description,
			Ncm:         product.NCM[:2],
		}
		var err error
		item.UnitAmount, err = strconv.ParseFloat(product.Price, 64)
		if err != nil {
			return err
		}
		item.TotalAmount, err = strconv.ParseFloat(product.Price, 64)
		if err != nil {
			return err
		}
		itemList = append(itemList, item)
		code++
	}
	*items = itemList
	return nil
}

func (i *Invoice) tax() error {
	tax := models.Tax{
		TotalTax: 0,
		Icms: models.Icms{
			Origin:             "2",
			Cst:                "41",
			BaseTaxModality:    "3",
			BaseTax:            0,
			BaseTaxSTReduction: "0",
		},
		Pis: models.Pis{
			Amount:  0,
			Rate:    0,
			BaseTax: 0,
			Cst:     "08",
		},
		Cofins: models.Cofins{
			Amount:  0,
			Rate:    0,
			BaseTax: 0,
			Cst:     "08",
		},
	}

	for d := range i.ProductInvoice.Items {
		i.ProductInvoice.Items[d].Tax = tax
	}
	return nil
}

func (i *Invoice) payment() {
	i.ProductInvoice.Payment = []models.Payment{
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

func (i *Invoice) destination() error {
	if i.Data.BuyerState == "SP" {
		i.ProductInvoice.Destination = "internal_Operation"
		return nil
	} else if i.Data.BuyerState != "SP" || len(i.Data.BuyerState) > 1 {
		i.ProductInvoice.Destination = "interstate_Operation"
		return nil
	} else {
		return fmt.Errorf("could not calculate the destination")
	}
}

func verifyState(vr string, state bool) string {
	if state {
		vr = strings.Replace(vr, "x", "6", 1)
	} else {
		vr = strings.Replace(vr, "x", "5", 1)
	}
	return vr
}

func (Ongoing) verifyState(vr string, state bool) string {
	if state {
		vr = strings.Replace(vr, "x", "2", 1)
	} else {
		vr = strings.Replace(vr, "x", "1", 1)
	}
	return vr
}
