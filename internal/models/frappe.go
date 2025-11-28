package models

// Invoices represents the main "Invoices" DocType from Frappe Brazil Invoice app
type Invoices struct {
	Name            string `json:"name"`              // Auto-generated: INV-YYYY-####
	OperationType   string `json:"operation_type"`    // Natureza de Operação
	ClientType      string `json:"client_type"`       // PF or PJ
	FreightModality string `json:"freight_modality"`  // Modalidade de Frete
	NfRefSerie      string `json:"nf_ref_serie"`      // NF Ref. Série
	NfRefNum        string `json:"nf_ref_num"`        // NF Ref. Número
	NfRefAccessKey  string `json:"nf_ref_access_key"` // NF Ref. Chave de Acesso
	NfDeRetorno     int    `json:"nf_de_retorno"`     // NF de Retorno? (checkbox)

	// Client information
	ClientName        string `json:"client_name"`        // Nome / Razão Social
	ClientEmail       string `json:"client_email"`       // Email
	ContribuinteIcms  string `json:"contribuinte_icms"`  // Contribuinte ICMS status
	ClientPhone       string `json:"client_phone"`       // Telefone de Contato
	ClientIDNumber    string `json:"client_id_number"`   // CPF / CNPJ
	InscricaoEstadual string `json:"inscricao_estadual"` // Inscrição Estadual (IE)

	// Product/Transport information
	TaxTemplate        string        `json:"tax_template"`         // Tax Template link
	InvoicesTable      []ItemInvoice `json:"invoices_table"`       // Child table items
	ProductBrand       string        `json:"product_brand"`        // Marca
	ProductQuantity    string        `json:"product_quantity"`     // Quantidade
	ProductType        string        `json:"product_type"`         // Espécie
	Carrier            string        `json:"carrier"`              // Transportadora (Link to Carrier)
	ProductGrossWeight string        `json:"product_gross_weight"` // Peso Bruto
	ProductNetWeight   string        `json:"product_net_weight"`   // Peso Líquido

	// Additional information
	AdditionalInformation string `json:"additional_information"` // Informações Complementares

	// Totals
	TotalFreight   string  `json:"total_freight"`   // Frete
	TotalDiscount  string  `json:"total_discount"`  // Desconto
	TotalInsurance string  `json:"total_insurance"` // Seguro
	OtherExpenses  string  `json:"other_expenses"`  // Outras Despesas
	Total          float64 `json:"total"`           // Total
	TotalTax       float64 `json:"total_tax"`       // Total + Impostos

	// Delivery address
	DeliverySupervisor    string `json:"delivery_supervisor"`     // Responsável
	DeliveryCEP           string `json:"delivery_cep"`            // CEP
	DeliveryAddress       string `json:"delivery_address"`        // Endereço
	DeliveryNeighborhood  string `json:"delivery_neighborhood"`   // Bairro
	DeliveryIBGE          string `json:"delivery_ibge"`           // IBGE code
	DeliveryPhone         string `json:"delivery_phone"`          // Telefone
	DeliveryState         string `json:"delivery_state"`          // Estado
	City                  string `json:"city"`                    // Cidade
	DeliveryNumberAddress string `json:"delivery_number_address"` // Nº do Endereço
	DeliveryComplement    string `json:"delivery_complement"`     // Complemento

	// NFe.io response fields
	InvoiceID     string `json:"invoice_id"`     // Invoice ID from NFe.io
	InvoiceSerie  string `json:"invoice_serie"`  // Invoice Serie
	InvoiceLink   string `json:"invoice_link"`   // Invoice Link (PDF URL)
	InvoiceNumber string `json:"invoice_number"` // Invoice Number

	// Errors/Logs
	ErrorsField string `json:"errors_field"` // Logs

	// Standard Frappe fields
	AmendedFrom string `json:"amended_from"`
	DocStatus   int    `json:"docstatus"` // 0=Draft, 1=Submitted, 2=Cancelled
	Creation    string `json:"creation"`
	Modified    string `json:"modified"`
	ModifiedBy  string `json:"modified_by"`
	Owner       string `json:"owner"`
}

// ItemInvoice represents the "Item Invoice" child table
type ItemInvoice struct {
	Name        string `json:"name"`
	Parent      string `json:"parent"`
	ParentField string `json:"parentfield"`
	ParentType  string `json:"parenttype"`
	Idx         int    `json:"idx"`

	SerialNumber string  `json:"serial_number"` // Serial No (Link)
	ItemCode     string  `json:"item_code"`     // Item (Link)
	ItemName     string  `json:"item_name"`     // Item Name
	Rate         float64 `json:"rate"`          // Rate
	Quantity     int     `json:"quantity"`      // Quantity
	NCM          string  `json:"ncm"`           // NCM code

	// Tax information
	InvoiceTaxes string  `json:"invoice_taxes"` // Link to Tax template
	ICMSRate     string  `json:"icms_rate"`     // ICMS %
	IPIRate      string  `json:"ipi_rate"`      // IPI %
	RateTaxes    float64 `json:"rate_taxes"`    // Rate with Taxes
	PISRate      string  `json:"pis_rate"`      // PIS %
	COFINSRate   string  `json:"cofins_rate"`   // COFINS %
}

// FrappeTax represents the "Tax" DocType with detailed tax configuration
type FrappeTax struct {
	Name string `json:"name"`

	// ICMS
	OriginICMS      string  `json:"origin_icms"`
	CSTICMS         string  `json:"cst_icms"`
	ModDetermBC     string  `json:"mod_determ_bc"`
	BaseCalcICMS    float64 `json:"base_calc_icms"`
	AliqICMS        float64 `json:"aliq_icms"`
	AliqFCP         float64 `json:"aliq_fcp"`
	BaseCalcICMSFCP float64 `json:"base_calc_icms_fcp"`
	ICMSValueFCP    float64 `json:"icms_value_fcp"`

	// ICMS calculation flags
	CalcularAutomaticamenteICMS bool `json:"calcular_automaticamente_icms"`
	AdicionaOutrasDespesasICMS  bool `json:"adiciona_outras_despesas_icms"`
	AdicionaFreteICMS           bool `json:"adiciona_frete_icms"`
	AdicionaIPIICMS             bool `json:"adiciona_ipi_icms"`
	AdicionaSeguroICMS          bool `json:"adiciona_seguro_icms"`
	AplicarAliqAutoICMS         bool `json:"aplicar_aliq_auto_icms"`

	// Substituição Tributária
	ModBaseCalcICMSTrib string  `json:"mod_base_calc_icms_trib"`
	CestICMSTrib        string  `json:"cest_icms_trib"`
	BaseICMSTrib        float64 `json:"base_icms_trib"`
	MVATrib             float64 `json:"mva_trib"`
	CreditoTrib         float64 `json:"credito_trib"`
	ReducaoTrib         float64 `json:"reducao_trib"`

	// Partilha ICMS (DIFAL)
	ValorBaseCalculoICMSDestino float64 `json:"valor_da_base_de_calculo_icms_no_destino"`
	AliqICMSDestino             float64 `json:"aliq_do_icms_do_estado_de_destino"`
	AliqICMSInterestadual       float64 `json:"aliq_do_icms_interestadual"`
	ValorBaseCalculoFCPDestino  float64 `json:"valor_da_base_de_calculo_fcp_na_uf_destino"`
	ValorICMSInterDIFAL         float64 `json:"valor_icms_inter_puf_de_destino_difal"`
	AliqFundoPobre              float64 `json:"aliq_fundo_pobre"`

	// IPI
	CSTIPI                     string  `json:"cst_ipi"`
	BaseDeCalculoIPI           float64 `json:"base_de_calculo_ipi"`
	ValorDoIPI                 float64 `json:"valor_do_ipi"`
	CodEnquadramento           float64 `json:"cod_enquadramento"`
	AliquotaIPI                float64 `json:"aliquota_ipi"`
	CalcularAutomaticamenteIPI bool    `json:"calcular_automaticamente_ipi"`

	// COFINS
	CSTCofins                     string  `json:"cst_cofins"`
	BaseDeCalculoCofins           float64 `json:"base_de_calculo_cofins"`
	ValorCofins                   float64 `json:"valor_cofins"`
	AliquotaCofins                float64 `json:"aliquota_cofins"`
	CalcularAutomaticamenteCofins bool    `json:"calcular_automaticamente_cofins"`

	// PIS
	CSTPis                     string  `json:"cst_pis"`
	BaseDeCalculoPis           float64 `json:"base_de_calculo_pis"`
	ValorPis                   float64 `json:"valor_pis"`
	AliquotaPis                float64 `json:"aliquota_pis"`
	CalcularAutomaticamentePis bool    `json:"calcular_automaticamente_pis"`
}

// Carrier represents transporter information (if you have a Carrier doctype)
type Carrier struct {
	Name              string `json:"name"`
	CarrierName       string `json:"carrier_name"`
	CNPJ              string `json:"cnpj"`
	StateRegistration string `json:"state_registration"`
	Address           string `json:"address"`
	City              string `json:"city"`
	State             string `json:"state"`
	CEP               string `json:"cep"`
	Phone             string `json:"phone"`
	Email             string `json:"email"`
	IBGE              string `json:"ibge"`
	Neighborhood      string `json:"neighborhood"`
}

// CustomFrappeInvoice kept for backward compatibility - maps to Invoices
type CustomFrappeInvoice struct {
	Name         string             `json:"name"`
	CustomerName string             `json:"customer_name"`
	CNPJ         string             `json:"cnpj_cpf"`
	NaturezaOp   string             `json:"natureza_operacao"`
	Status       string             `json:"status"`
	GrandTotal   float64            `json:"grand_total"`
	Items        []CustomFrappeItem `json:"items"`
}

// CustomFrappeItem kept for backward compatibility
type CustomFrappeItem struct {
	ItemCode   string  `json:"item_code"`
	ItemName   string  `json:"item_name"`
	Qty        float64 `json:"qty"`
	Rate       float64 `json:"rate"`
	Amount     float64 `json:"amount"`
	NCM        string  `json:"ncm"`
	CFOP       string  `json:"cfop"`
	CEST       string  `json:"cest"`
	TotalTaxes float64 `json:"total_taxes"`
}
