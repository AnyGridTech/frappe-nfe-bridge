//go:build ignore
// +build ignore

package invoice

import (
	"encoding/json"
	"erpnext-go/pkg/modules"
	"log"
	"os"
)

func (n *Nfeio) getInvoiceViaAccessKey(accessKey string, test bool) (*models.NfeioInvoiceResponse, error) {

	nfeio := modules.NewNfeioModule(log.New(os.Stdout, "erpnext-go ", 0))
	resp, err := nfeio.GetInvoiceAccessKey(accessKey)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var invoiceResponse models.NfeioInvoiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&invoiceResponse); err != nil {
		return nil, err
	}
	return &invoiceResponse, nil
}
