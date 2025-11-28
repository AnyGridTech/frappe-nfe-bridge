//go:build ignore
// +build ignore

package invoice

import (
	"erpnext-go/pkg/models"
	"fmt"
)

type TransportIE interface {
	transport() error
	getCarrier() models.InverterDb
}

func (i *Invoice) transport() error {
	if i.Data.Carrier.Cnpj == "" {
		return fmt.Errorf("carrier data is missing")
	}
	i.ProductInvoice.Transport.FreightModality = "ByIssuer"
	i.ProductInvoice.Transport.TransportGroup = models.TransportGroup{
		StateTaxNumber:   i.Data.Carrier.StateRegistration,
		Name:             i.Data.Carrier.CompanyName,
		FederalTaxNumber: i.Data.Carrier.Cnpj,
		Email:            i.Data.Carrier.Email,
		Address: models.Address{
			Phone: formatPhoneNumber(i.Data.Carrier.Phone),
			State: i.Data.Carrier.State,
			City: models.City{
				Name: i.Data.Carrier.City,
				Code: i.Data.Carrier.Ibge,
			},

			District:              i.Data.Carrier.Neighborhood,
			AdditionalInformation: "",
			Street:                i.Data.Carrier.Address,
			Number:                i.Data.Carrier.AddNumber,
			PostalCode:            i.Data.Carrier.Cep,
			Country:               "BRA",
		},

		Type: "legalEntity",
	}
	i.volume()
	return nil
}

func ( i *Invoice) volume() {
	vol := &i.ProductInvoice.Transport.Volume
	vol.VolumeQuantity = 1 //needs verify
	vol.Species = "Caixa"
	vol.Brand = "Growatt"
	vol.VolumeNumeration = i.Data.Product[0].Name
	vol.GrossWeight = func() float64 {
		var total float64
		for _, item := range i.Data.Product {
			total += item.GrossWeight
		}
		return total
	}()
}