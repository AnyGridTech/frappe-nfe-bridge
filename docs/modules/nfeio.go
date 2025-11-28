//go:build ignore
// +build ignore

package modules

import (
	"encoding/json"
	"erpnext-go/pkg/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type NfeioModule struct {
	l *log.Logger
}

type Messages struct {
	Uri string `json:"uri"`
}

var (
	endpoint         string
	apiKey           string
	endpoint_consult string
)

func Initialize() {
	endpoint = os.Getenv("NFEIO_ENDPOINT")
	endpoint_consult = os.Getenv("NFEIO_ENDPOINT_CONSULT")
	apiKey = os.Getenv("NFEIO_APIKEY_ANYGRID_SOLAR")
}

func NewNfeioModule(l *log.Logger) *NfeioModule {
	Initialize()
	return &NfeioModule{l}
}

func (n *NfeioModule) CreateInvoice(invoice io.Reader, companyKey string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/productinvoices?apiKey=%s", endpoint, companyKey, apiKey)
	return http.Post(url, "application/json", invoice)
}

func (n *NfeioModule) GetInvoice(id, companyKey string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/productinvoices/%s?apiKey=%s", endpoint, companyKey, id, apiKey)
	return http.Get(url)
}

func (n *NfeioModule) CreateCorrectionLetter(reason io.Reader, id string, companyKey string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/productinvoices/%s/correctionletter?apiKey=%s", endpoint, companyKey, id, apiKey)
	req, err := http.NewRequest(http.MethodPut, url, reason)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return doRequest(req)
}

func (n *NfeioModule) DeleteInvoice(id, companyKey string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/productinvoices/%s?apikey=%s", endpoint, companyKey, id, apiKey)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return doRequest(req)
}

func (n *NfeioModule) GetInvoicePdf(id, companyKey string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/productinvoices/%s/pdf?apikey=%s", endpoint, companyKey, id, apiKey)
	n.l.Println("getting PDF uri...")
	return http.Get(url)

}

func (n *NfeioModule) GetInvoiceAccesKeyModule(accessKey string) *http.Response {
	url := fmt.Sprintf("%s/%s?apikey=%s", endpoint, accessKey, apiKey)
	n.l.Println("Getting Invoice")
	res, err := http.Get(url)
	if err != nil {
		n.l.Println(err)
	}
	return res
}

func (n *NfeioModule) GetCorrectionLetterPdf(id, companyKey string) (*http.Response, error) {
	urlPdf := fmt.Sprintf("%s/%s/productinvoices/%s/correctionletter/pdf?apikey=%s", endpoint, companyKey, id, apiKey)
	n.l.Println("getting PDF uri...")
	return http.Get(urlPdf)
}

func (n *NfeioModule) GetCorrectionLetterXml(id, companyKey string) (*http.Response, error) {
	urlPdf := fmt.Sprintf("%s/%s/productinvoices/%s/correctionletter/xml?apikey=%s", endpoint, companyKey, id, apiKey)
	n.l.Println("getting PDF uri...")
	return http.Get(urlPdf)
}

func (n *NfeioModule) GetInvoiceAccessKey(accessKey string) (*http.Response, error) {
	url := fmt.Sprintf("%s/productinvoices/%s?apikey=%s", endpoint_consult, accessKey, apiKey)
	n.l.Println("Getting Invoice...")
	return http.Get(url)
}

/*
DEPRECATED: This function
This function, GetViaCep, is a method of the NfeioModule struct.
It takes a string argument 'cep' and returns a Cep struct and an error.
The function fetches an environment variable for the endpoint, appends the cep to the URL and makes a GET request.
If there is an error in the request, it logs the error.
It reads the response body and closes the body.
If the status code of the response is greater than 299, it logs the status code and the body.
If there is an error in reading the body, it logs the error.
It unmarshals the body into a Cep struct and returns the struct and any error.
*/
func (n *NfeioModule) GetViaCep(cep string) (models.Cep, error) {
	endpoint := os.Getenv("OPERACAO_ENDPOINT")
	urlCep := endpoint + "/api/cep/viaCep/" + cep
	var response models.Cep
	res, err := http.Get(urlCep)
	if err != nil {
		n.l.Println(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		n.l.Printf("ViaCep: Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		n.l.Println(err)
	}
	json.Unmarshal([]byte(body), &response)
	return response, err
}
