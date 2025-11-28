package models

import "time"

type ProductInvoiceRequest struct {
	Serie                    int                       `json:"serie"`
	Number                   int                       `json:"number,omitempty"`
	OperationOn              time.Time                 `json:"-"`
	OperationNature          string                    `json:"operationNature,omitempty"`
	OperationType            string                    `json:"operationType"`
	PrintType                string                    `json:"printType,omitempty"`
	ConsumerType             string                    `json:"consumerType"`
	PresenceType             string                    `json:"presenceType,omitempty"`
	ContingencyOn            time.Time                 `json:"-"`
	ContingencyJustification string                    `json:"contingencyJustification,omitempty"`
	Buyer                    Buyer                     `json:"buyer"`
	Payment                  []Payment                 `json:"payment"`
	Totals                   *Totals                   `json:"totals,omitempty"`
	Items                    []Items                   `json:"items"`
	Transport                Transport                 `json:"transport"`
	Billing                  []Billing                 `json:"billing,omitempty"`
	Issuer                   *Issuer                   `json:"issuer,omitempty"`
	TransactionIntermediate  []TransactionIntermediate `json:"transactionIntermediate,omitempty"`
	Delivery                 []Delivery                `json:"delivery,omitempty"`
	Withdrawal               []Withdrawal              `json:"withdrawal,omitempty"`
	AdditionalInformation    *AdditionalInformation    `json:"additionalInformation,omitempty"`
	Destination              string                    `json:"destination"`
	PurposeType              string                    `json:"purposeType"`
	Id                       string                    `json:"id,omitempty"`
}

type NfeioWebhook struct {
	Action                  string                `json:"action,omitempty"`
	ID                      string                `json:"id,omitempty"`
	Serie                   int                   `json:"serie,omitempty"`
	Number                  int                   `json:"number,omitempty"`
	Status                  string                `json:"status,omitempty"`
	Authorization           Authorization         `json:"authorization,omitempty"`
	OperationNature         string                `json:"operationNature,omitempty"`
	CreatedOn               time.Time             `json:"createdOn,omitempty"`
	ModifiedOn              time.Time             `json:"modifiedOn,omitempty"`
	OperationOn             any                   `json:"operationOn,omitempty"`
	OperationType           string                `json:"operationType,omitempty"`
	EnvironmentType         string                `json:"environmentType,omitempty"`
	PurposeType             string                `json:"purposeType,omitempty"`
	Issuer                  Issuer                `json:"issuer,omitempty"`
	Buyer                   Buyer                 `json:"buyer,omitempty"`
	Totals                  Totals                `json:"totals,omitempty"`
	Transport               Transport             `json:"transport,omitempty"`
	AdditionalInformation   AdditionalInformation `json:"additionalInformation,omitempty"`
	Export                  any                   `json:"export,omitempty"`
	Billing                 any                   `json:"billing,omitempty"`
	Payment                 []Payment             `json:"payment,omitempty"`
	TransactionIntermediate any                   `json:"transactionIntermediate,omitempty"`
	Delivery                any                   `json:"delivery,omitempty"`
	Withdrawal              any                   `json:"withdrawal,omitempty"`
	LastEvents              LastEvents            `json:"lastEvents,omitempty"`
}
type Authorization struct {
	ReceiptOn time.Time `json:"receiptOn,omitempty"`
	AccessKey string    `json:"accessKey,omitempty"`
	Message   string    `json:"message,omitempty"`
}
type LastEvents struct {
	Events  []Events `json:"events,omitempty"`
	HasMore bool     `json:"hasMore,omitempty"`
}
type Events struct {
	Data     Data    `json:"data,omitempty"`
	Type     string  `json:"type,omitempty"`
	Sequence float32 `json:"sequence,omitempty"`
}
type Data struct {
	URI                 string    `json:"uri,omitempty"`
	ContentType         string    `json:"contentType,omitempty"`
	CreatedOn           time.Time `json:"createdOn,omitempty"`
	AccessKey           string    `json:"accessKey,omitempty"`
	ApplicationVersion  string    `json:"applicationVersion,omitempty"`
	Description         string    `json:"description,omitempty"`
	EnvironmentType     string    `json:"environmentType,omitempty"`
	ProtocolNumber      string    `json:"protocolNumber,omitempty"`
	ValidatorDigit      string    `json:"validatorDigit,omitempty"`
	StatusCode          float32   `json:"statusCode,omitempty"`
	Status              interface{}   `json:"status,omitempty"`
	Message             string    `json:"message,omitempty"`
	ReceiptNumber       string    `json:"receiptNumber,omitempty"`
	AccessKeyCheckDigit string    `json:"accessKeyCheckDigit,omitempty"`
	CheckCode           float32   `json:"checkCode,omitempty"`
	Serie               float32   `json:"serie,omitempty"`
	Number              float32   `json:"number,omitempty"`
	BatchID             float32   `json:"batchId,omitempty"`
}
type Buyer struct {
	Address                 Address     `json:"address"`
	StateTaxNumberIndicator interface{} `json:"stateTaxNumberIndicator,omitempty"`
	TradeName               string      `json:"tradeName,omitempty"`
	TaxRegime               string      `json:"taxRegime,omitempty"`
	StateTaxNumber          string      `json:"stateTaxNumber,omitempty"`
	ID                      string      `json:"id,omitempty"`
	Name                    string      `json:"name"`
	FederalTaxNumber        int         `json:"federalTaxNumber"`
	Email                   string      `json:"email,omitempty"`
	Type                    string      `json:"type"`
}
type Totals struct {
	Icms  Icms `json:"icms"`
	Issqn struct {
		TotalServiceNotTaxedICMS int       `json:"totalServiceNotTaxedICMS,omitempty"`
		BaseRateISS              int       `json:"baseRateISS,omitempty"`
		TotalISS                 int       `json:"totalISS,omitempty"`
		ValueServicePIS          int       `json:"valueServicePIS,omitempty"`
		ValueServiceCOFINS       int       `json:"valueServiceCOFINS,omitempty"`
		ProvisionService         time.Time `json:"provisionService,omitempty"`
		DeductionReductionBC     int       `json:"deductionReductionBC,omitempty"`
		ValueOtherRetention      int       `json:"valueOtherRetention,omitempty"`
		DiscountUnconditional    int       `json:"discountUnconditional,omitempty"`
		DiscountConditioning     int       `json:"discountConditioning,omitempty"`
		TotalRetentionISS        int       `json:"totalRetentionISS,omitempty"`
		CodeTaxRegime            int       `json:"codeTaxRegime,omitempty"`
	} `json:"issqn,omitempty"`
}
type Transport struct {
	TransportGroup TransportGroup `json:"transportGroup"`
	Reboque        struct {
		Plate string `json:"plate"`
		Uf    string `json:"uf"`
		Rntc  string `json:"rntc"`
		Wagon string `json:"wagon"`
		Ferry string `json:"ferry"`
	} `json:"-"`
	Volume           Volume `json:"volume"`
	TransportVehicle struct {
		Plate string `json:"plate"`
		State string `json:"state"`
		Rntc  string `json:"rntc"`
	} `json:"-"`
	SealNumber      string      `json:"sealNumber"`
	FreightModality interface{} `json:"freightModality,omitempty"`
	TranspRate      struct {
		ServiceAmount         int `json:"serviceAmount"`
		BcRetentionAmount     int `json:"bcRetentionAmount"`
		IcmsRetentionRate     int `json:"icmsRetentionRate"`
		IcmsRetentionAmount   int `json:"icmsRetentionAmount"`
		Cfop                  int `json:"cfop"`
		CityGeneratorFactCode int `json:"cityGeneratorFactCode"`
	} `json:"-"`
}

type TransportGroup struct {
	StateTaxNumber     string      `json:"stateTaxNumber"`
	TransportRetention string      `json:"transportRetention,omitempty"`
	ID                 string      `json:"id,omitempty"`
	Name               string      `json:"name"`
	FederalTaxNumber   interface{} `json:"federalTaxNumber"`
	Email              string      `json:"email"`
	FullAddress        string      `json:"fullAddress,omitempty"`
	Address            Address     `json:"address"`
	Type               string      `json:"type"`
	CityName           string      `json:"cityName,omitempty"`
}

type Volume struct {
	VolumeQuantity   int     `json:"volumeQuantity"`
	Species          string  `json:"species"`
	Brand            string  `json:"brand"`
	VolumeNumeration string  `json:"volumeNumeration"`
	NetWeight        float64 `json:"netWeight"`
	GrossWeight      float64 `json:"grossWeight"`
}
type Address struct {
	Phone                 string `json:"phone"`
	City                  City   `json:"city"`
	State                 string `json:"state"`
	Country               string `json:"country"`
	District              string `json:"district"`
	AdditionalInformation string `json:"additionalInformation,omitempty"`
	Street                string `json:"street"`
	Number                string `json:"number"`
	PostalCode            string `json:"postalCode"`
}
type City struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type TaxCouponInformation struct {
	ModelDocumentFiscal string `json:"modelDocumentFiscal,omitempty"`
	OrderECF            string `json:"orderECF,omitempty"`
	OrderCountOperation int    `json:"orderCountOperation,omitempty"`
}

type DocumentInvoiceReference struct {
	State            int    `json:"state,omitempty"`
	YearMonth        string `json:"yearMonth,omitempty"`
	FederalTaxNumber string `json:"federalTaxNumber,omitempty"`
	Model            string `json:"model,omitempty"`
	Series           string `json:"series,omitempty"`
	Number           string `json:"number,omitempty"`
}
type DocumentElectronicInvoice struct {
	AccessKey string `json:"accessKey,omitempty"`
}

type TaxDocumentsReference struct {
	TaxCouponInformation      []TaxCouponInformation    `json:"taxCouponInformation,omitempty"`
	DocumentInvoiceReference  DocumentInvoiceReference  `json:"documentInvoiceReference,omitempty"`
	DocumentElectronicInvoice DocumentElectronicInvoice `json:"documentElectronicInvoice,omitempty"`
}
type ReferencedProcess []struct {
	IdentifierConcessory string `json:"identifierConcessory,omitempty"`
	IdentifierOrigin     int    `json:"identifierOrigin,omitempty"`
}

type TaxpayerComments []struct {
	Field string `json:"field,omitempty"`
	Text  string `json:"text,omitempty"`
}
type AdditionalInformation struct {
	Fisco                 string                  `json:"fisco,omitempty"`
	Taxpayer              string                  `json:"taxpayer,omitempty"`
	XMLAuthorized         []int                   `json:"xmlAuthorized,omitempty"`
	Effort                string                  `json:"effort,omitempty"`
	Order                 string                  `json:"order,omitempty"`
	Contract              string                  `json:"contract,omitempty"`
	TaxDocumentsReference []TaxDocumentsReference `json:"taxDocumentsReference,omitempty"`
	TaxpayerComments      TaxpayerComments        `json:"taxpayerComments,omitempty"`
	ReferencedProcess     ReferencedProcess       `json:"referencedProcess,omitempty"`
}
type Items struct {
	Code                  string               `json:"code"`
	CodeGTIN              string               `json:"codeGTIN"`
	CodeTaxGTIN           string               `json:"codeTaxGTIN"`
	Description           string               `json:"description"`
	Ncm                   string               `json:"ncm"`
	Cfop                  int                  `json:"cfop"`
	Unit                  string               `json:"unit"`
	Quantity              int                  `json:"quantity"`
	UnitAmount            float64              `json:"unitAmount"`
	TotalAmount           float64              `json:"totalAmount"`
	UnitTax               string               `json:"unitTax,omitempty"`
	QuantityTax           int                  `json:"quantityTax,omitempty"`
	TaxUnitAmount         int                  `json:"taxUnitAmount"`
	FreightAmount         int                  `json:"freightAmount,omitempty"`
	InsuranceAmount       int                  `json:"insuranceAmount,omitempty"`
	DiscountAmount        int                  `json:"discountAmount,omitempty"`
	OthersAmount          int                  `json:"othersAmount,omitempty"`
	TotalIndicator        bool                 `json:"-"`
	Cest                  string               `json:"cest,omitempty"`
	Tax                   Tax                  `json:"tax"`
	AdditionalInformation string               `json:"additionalInformation,omitempty"`
	NumberOrderBuy        string               `json:"numberOrderBuy,omitempty"`
	ItemNumberOrderBuy    int                  `json:"itemNumberOrderBuy,omitempty"`
	FuelDetail            []FuelDetail         `json:"fuelDetail,omitempty"`
	Benefit               string               `json:"benefit,omitempty"`
	ImportDeclarations    []ImportDeclarations `json:"importDeclarations,omitempty"`
}
type Billing struct {
	Bill struct {
		Number         string `json:"number"`
		OriginalAmount int    `json:"originalAmount"`
		DiscountAmount int    `json:"discountAmount"`
		NetAmount      int    `json:"netAmount"`
	} `json:"bill"`
	Duplicates []struct {
		Number       string    `json:"number"`
		ExpirationOn time.Time `json:"expirationOn"`
		Amount       int       `json:"amount"`
	} `json:"duplicates"`
}
type Issuer struct {
	StStateTaxNumber string `json:"stStateTaxNumber,omitempty"`
}
type TransactionIntermediate struct {
	FederalTaxNumber int    `json:"federalTaxNumber,omitempty"`
	Identifier       string `json:"identifier,omitempty"`
}
type Delivery struct {
	StateTaxNumber   string `json:"stateTaxNumber"`
	ID               string `json:"id"`
	Name             string `json:"name"`
	FederalTaxNumber int    `json:"federalTaxNumber"`
	Email            string `json:"email"`
	Address          struct {
		Phone string `json:"phone"`
		State string `json:"state"`
		City  struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"city"`
		District              string `json:"district"`
		AdditionalInformation string `json:"additionalInformation,omitempty"`
		Street                string `json:"street"`
		Number                string `json:"number"`
		PostalCode            string `json:"postalCode"`
		Country               string `json:"country"`
	} `json:"address"`
	Type string `json:"type"`
}
type Withdrawal struct {
	StateTaxNumber   string `json:"stateTaxNumber"`
	ID               string `json:"id"`
	Name             string `json:"name"`
	FederalTaxNumber int    `json:"federalTaxNumber"`
	Email            string `json:"email"`
	Address          struct {
		Phone string `json:"phone"`
		State string `json:"state"`
		City  struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"city"`
		District              string `json:"district"`
		AdditionalInformation string `json:"additionalInformation,omitempty"`
		Street                string `json:"street"`
		Number                string `json:"number"`
		PostalCode            string `json:"postalCode"`
		Country               string `json:"country"`
	} `json:"address"`
	Type string `json:"type"`
}
type Payment struct {
	PaymentDetail []PaymentDetail `json:"paymentDetail"`
	PayBack       int             `json:"payBack,omitempty"`
}
type PaymentDetail struct {
	Method string `json:"method"`
	Amount int    `json:"amount"`
	Card   []Card `json:"card,omitempty"`
}
type Card struct {
	FederalTaxNumber       string `json:"federalTaxNumber"`
	Flag                   string `json:"flag"`
	Authorization          string `json:"authorization"`
	IntegrationPaymentType string `json:"integrationPaymentType"`
}
type Tax struct {
	TotalTax float64`json:"totalTax"`
	Icms     Icms        `json:"icms"`
	Ipi      Ipi         `json:"-"`
	Ii       Ii          `json:"-"`
	Pis      Pis         `json:"pis"`
	Cofins   Cofins      `json:"cofins"`
}

type Ipi struct {
	Classification     string  `json:"classification,omitempty"`
	ProducerCNPJ       string  `json:"producerCNPJ,omitempty"`
	StampCode          string  `json:"stampCode,omitempty"`
	StampQuantity      float64 `json:"stampQuantity,omitempty"`
	ClassificationCode string  `json:"classificationCode,omitempty"`
	Cst                string  `json:"cst,omitempty"`
	Base               float64 `json:"base,omitempty"`
	Rate               float64 `json:"rate,omitempty"`
	UnitQuantity       float64 `json:"unitQuantity,omitempty"`
	UnitAmount         float64 `json:"unitAmount,omitempty"`
	Amount             float64 `json:"amount,omitempty"`
}
type Ii struct {
	BaseTax                  string `json:"baseTax"`
	CustomsExpenditureAmount string `json:"customsExpenditureAmount"`
	Amount                   int    `json:"amount"`
	IofAmount                int    `json:"iofAmount"`
}

type FuelDetail struct {
	CodeANP        string `json:"codeANP"`
	PercentageNG   int    `json:"percentageNG"`
	DescriptionANP string `json:"descriptionANP"`
	PercentageGLP  int    `json:"percentageGLP"`
	PercentageNGn  int    `json:"percentageNGn"`
	PercentageGNi  int    `json:"percentageGNi"`
	StartingAmount int    `json:"startingAmount"`
	Codif          string `json:"codif"`
	AmountTemp     int    `json:"amountTemp"`
	StateBuyer     string `json:"stateBuyer"`
	Cide           []Cide `json:"cide"`
	Pump           []Pump `json:"pump"`
}
type Pump struct {
	SpoutNumber     int `json:"spoutNumber"`
	Number          int `json:"number"`
	TankNumber      int `json:"tankNumber"`
	BeginningAmount int `json:"beginningAmount"`
	EndAmount       int `json:"endAmount"`
}
type Cide struct {
	Bc         int `json:"bc"`
	Rate       int `json:"rate"`
	CideAmount int `json:"cideAmount"`
}
type ImportDeclarations []struct {
	Code                  string    `json:"code"`
	RegisteredOn          time.Time `json:"registeredOn"`
	CustomsClearanceName  string    `json:"customsClearanceName"`
	CustomsClearanceState string    `json:"customsClearanceState"`
	CustomsClearancedOn   time.Time `json:"customsClearancedOn"`
	Additions             []struct {
		Code         int    `json:"code"`
		Manufacturer string `json:"manufacturer"`
		Amount       int    `json:"amount"`
		Drawback     int    `json:"drawback"`
	} `json:"additions"`
	Exporter               string `json:"exporter"`
	InternationalTransport string `json:"internationalTransport"`
	Intermediation         string `json:"intermediation"`
}

type Icms struct {
	Origin                     string `json:"origin,omitempty"`
	Cst                        string `json:"cst,omitempty"`
	BaseTaxModality            string `json:"baseTaxModality"`
	BaseTax                    float64    `json:"baseTax"`
	BaseTaxSTModality          string `json:"baseTaxSTModality,omitempty"`
	BaseTaxSTReduction         string `json:"baseTaxSTReduction"`
	BaseTaxST                  int    `json:"baseTaxST,omitempty"`
	BaseTaxReduction           int    `json:"baseTaxReduction,omitempty"`
	StRate                     int    `json:"stRate,omitempty"`
	StAmount                   int    `json:"stAmount,omitempty"`
	StMarginAmount             int    `json:"stMarginAmount,omitempty"`
	Csosn                      string `json:"csosn,omitempty"`
	Rate                       float64    `json:"rate"`
	Amount                     float64    `json:"amount,omitempty"`
	Percentual                 int    `json:"percentual,omitempty"`
	SnCreditRate               int    `json:"snCreditRate,omitempty"`
	SnCreditAmount             int    `json:"snCreditAmount,omitempty"`
	StMarginAddedAmount        string `json:"stMarginAddedAmount,omitempty"`
	StRetentionAmount          string `json:"stRetentionAmount,omitempty"`
	BaseSTRetentionAmount      string `json:"baseSTRetentionAmount,omitempty"`
	BaseTaxOperationPercentual string `json:"baseTaxOperationPercentual,omitempty"`
	Ufst                       string `json:"ufst,omitempty"`
	AmountSTReason             string `json:"amountSTReason,omitempty"`
	BaseSNRetentionAmount      string `json:"baseSNRetentionAmount,omitempty"`
	SnRetentionAmount          string `json:"snRetentionAmount,omitempty"`
	AmountOperation            string `json:"amountOperation,omitempty"`
	PercentualDeferment        string `json:"percentualDeferment,omitempty"`
	BaseDeferred               string `json:"baseDeferred,omitempty"`
	ExemptAmount               int    `json:"exemptAmount,omitempty"`
	ExemptReason               string `json:"exemptReason,omitempty"`
	FcpRate                    int    `json:"fcpRate,omitempty"`
	FcpAmount                  int    `json:"fcpAmount,omitempty"`
	FcpstRate                  int    `json:"fcpstRate,omitempty"`
	FcpstAmount                int    `json:"fcpstAmount,omitempty"`
	FcpstRetRate               int    `json:"fcpstRetRate,omitempty"`
	FcpstRetAmount             int    `json:"fcpstRetAmount,omitempty"`
	BaseTaxFCPSTAmount         int    `json:"baseTaxFCPSTAmount,omitempty"`
	SubstituteAmount           int    `json:"substituteAmount,omitempty"`
}
type Pis struct {
	Cst                    string  `json:"cst"`
	BaseTax                float64 `json:"baseTax"`
	Rate                   float64 `json:"rate"`
	Amount                 float64 `json:"amount"`
	BaseTaxProductQuantity float64 `json:"baseTaxProductQuantity,omitempty"`
	ProductRate            float64 `json:"productRate,omitempty"`
}
type Cofins struct {
	Cst                    string  `json:"cst"`
	BaseTax                float64 `json:"baseTax"`
	Rate                   float64 `json:"rate"`
	Amount                 float64 `json:"amount"`
	BaseTaxProductQuantity float64 `json:"baseTaxProductQuantity,omitempty"`
	ProductRate            float64 `json:"productRate,omitempty"`
}


// Response from NFe.io
type ProductInvoiceResponse struct {
    ID              string `json:"id"`
    Status          string `json:"status"`
    Environment     string `json:"environment"`
    FlowStatus      string `json:"flowStatus"`
    PdfUrl          string `json:"pdf"`
    XmlUrl          string `json:"xml"`
}