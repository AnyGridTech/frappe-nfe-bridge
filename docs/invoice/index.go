package invoice

import (
	"bytes"
	"encoding/json"
	"erpnext-go/internal/api/dto"
	DTO "erpnext-go/internal/api/dto"
	"erpnext-go/internal/lib"
	"erpnext-go/internal/namespace"
	"erpnext-go/pkg/models"
	"erpnext-go/pkg/modules"
	"erpnext-go/pkg/services/erpnext"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var keys = namespace.GeneralKey

type Invoice struct {
	Params         Params
	ProductInvoice models.ProductInvoice
	Data           models.Req
}

type Nfeio struct {
	ctx *fiber.Ctx
}

func NewNfeio(f *fiber.Ctx) *Nfeio {
	return &Nfeio{ctx: f}
}
func NewInvoice(params Params, req models.Req) *Invoice {
	return &Invoice{Params: params, Data: req}
}

func (n *Nfeio) CreateInvoice() error {
	dto, err := lib.GetFromCtx[interface{}](n.ctx, keys.Dto)
	if err != nil {
		return err
	}
	for index, item := range dto.(*DTO.WarrantyInvoice).Data.Items {
		tx := getInvoiceTaxes(item.InvoiceTaxesId)
		if tx == (DTO.InvoiceTaxes{}) {
			continue
		}
		dto.(*DTO.WarrantyInvoice).Data.Items[index].InvoiceTaxes = tx
	}
	body, err := buildInvoiceBody(dto.(*DTO.WarrantyInvoice))
	if err != nil {
		return err
	}
	nf := NewInvoice(Params{Op: Ongoing{msg: body.OperationNature}, Buyer: PF{}}, body)
	err = nf.Create()
	if err != nil {
		log.Println(err)
	}
	return nil
}
func buildInvoiceBody(dto *DTO.WarrantyInvoice) (models.Req, error) {
	var req models.Req
	//buyer
	req.OperationNature = dto.Data.OperationNature
	req.BuyerType = dto.Data.ClientType
	req.BuyerCEP = dto.Data.Cep
	req.BuyerCPF = dto.Data.Cpf
	req.BuyerCity = dto.Data.City
	req.BuyerAddress = dto.Data.Address
	req.BuyerNeighborhood = dto.Data.Neighborhood
	req.BuyerNumber = dto.Data.NumberAddress
	req.BuyerState = dto.Data.State
	req.BuyerName = dto.Data.Nome
	req.BuyerICMS = dto.Data.ContribuinteIcms
	req.BuyerStateTaxNumber = dto.Data.StateTaxNumber
	//delivery
	req.DeliveryAddress = dto.Data.Address
	req.DeliveryCity = dto.Data.City
	req.DeliveryNeighborhood = dto.Data.Neighborhood
	req.DeliveryNumber = dto.Data.NumberAddress
	req.DeliveryState = dto.Data.State
	req.DeliveryCEP = dto.Data.Cep
	req.DeliveryName = dto.Data.CollectGuy
	req.BuyerIBGE = dto.Data.Ibge
	req.DeliveryPhone = formatPhoneNumber(dto.Data.DeliveryPhone)
	//items
	for _, item := range dto.Data.Items {
		req.Product = append(req.Product, models.Product{
			Parent:          item.Parent,
			Name:            item.ItemName,
			Price:           strconv.FormatFloat(float64(item.Rate), 'f', 2, 64),
			NCM:             item.Ncm,
			NetWeight:       item.Netweight,
			GrossWeight:     item.Grossweight,
			PackagingHeight: item.PackagingHeight,
			PackagingWidth:  item.PackagingWidth,
			PackagingLength: item.PackagingLength,
			Description:     item.ItemName,
			Taxes: models.Tax{
				TotalTax: totalTax(item),
				Icms:     icmsTax(item),
				Pis:      pisTax(item),
				Cofins:   cofinsTax(item),
				Ipi: 	ipiTax(item),
			},
		})
	}
	//carrier
	req.Carrier = models.Carrier{
		Name:              dto.CarrierData.Data.FantasyName,
		Cnpj:              dto.CarrierData.Data.Cnpj,
		CompanyName:       dto.CarrierData.Data.RazaoSocial,
		StateRegistration: dto.CarrierData.Data.InscricaoEstadual,
		Address:           dto.CarrierData.Data.Address,
		City:              dto.CarrierData.Data.City,
		State:             dto.CarrierData.Data.State,
		Cep:               dto.CarrierData.Data.Cep,
		Phone:             formatPhoneNumber(dto.CarrierData.Data.Telefone),
		Ibge:              dto.CarrierData.Data.Ibge,
		Email:             dto.CarrierData.Data.Email,
		Neighborhood:      dto.CarrierData.Data.Neighborhood,
	}
	return req, nil
}

func getInvoiceTaxes(name string) dto.InvoiceTaxes {
	erp := erpnext.NewFrappe()
	var m map[string]json.RawMessage
	if err := erp.GetOneBy("Invoice Taxes", name, &m); err != nil {
		log.Println(err)
		return dto.InvoiceTaxes{}
	}
	var dt dto.InvoiceTaxes
	if err := json.Unmarshal(m["data"], &dt); err != nil {
		log.Println(err)
		return dto.InvoiceTaxes{}
	}
	return dt
}

func (i *Invoice) Create() error {

	log.Println(i.Data.BuyerType)
	i.ProductInvoice.PurposeType = "normal"
	i.ProductInvoice.Transport.SealNumber = i.Data.Product[0].Parent
	err := i.Params.Buyer.buyer(&i.ProductInvoice.Buyer, &i.Data)
	if err != nil {
		return err
	}
	err = i.Params.Op.items(&i.ProductInvoice.Items, &i.Data.Product)
	if err != nil {
		return err
	}
	err = address(&i.ProductInvoice.Buyer, &i.Data)
	if err != nil {
		return err
	}
	err = i.Params.Op.cfop(&i.ProductInvoice.Items)
	if err != nil {
		return err
	}
	err = i.tax()
	if err != nil {
		return err
	}
	err = i.Params.Op.OperationType(&i.ProductInvoice)
	if err != nil {
		return err
	}
	i.Params.Buyer.consumerType(&i.ProductInvoice)
	err = i.transport()
	if err != nil {
		return err
	}
	err = i.destination()
	if err != nil {
		return err
	}
	i.payment()

	file, err := json.Marshal(i.ProductInvoice)
	if err != nil {
		return err
	}

	err = os.WriteFile("productInvoice.json", file, 0644)
	if err != nil {
		return err
	}
	module := modules.NewNfeioModule(log.New(os.Stdout, "erpnext-go ", 0))
	res, err := module.CreateInvoice(bytes.NewBuffer(file), os.Getenv("NFEIO_COMPANYKEY"))
	if err != nil {
		return err
	}
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return fmt.Errorf("error reading response body: %v", readErr)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("error creating invoice, status: %v, body:\n %v ", res.StatusCode, string(body))
	} else {
		log.Println(string(body))

	}
	return nil
}

func formatPhoneNumber(s string) string {
	// Remove all non-numeric characters
	re := regexp.MustCompile(`\D`)
	formatted := re.ReplaceAllString(s, "")

	// Remove leading "+" if present
	formatted = strings.TrimPrefix(formatted, "+")

	return formatted
}
