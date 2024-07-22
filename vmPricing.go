package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

/* Need to iterate over all VMs to retrieve SKUs and consumption and location */
/* Ensure that no duplicates exists */
/* Minimize API traffic */

type vmReduction struct {
	Sku      string
	Location string
}

/* Reduce VMs will search the config for unique VM sku and location           */
/* matches. The function is not concerned with the consumption                */
/* IT outputs a slice of the type vmReduction. The vmReduction slice          */
/* is then used to retrieve the cost of a single PAYG hour, 1 & 3yr RIs       */
func (c Config) ReduceVms() []vmReduction {
	ret := []vmReduction{}

	/* Cycle through the applications searching for VMs */
	for _, app := range c.Applications {
		for _, env := range app.Environments {
			for _, vm := range env.VMs {
				match := false
				for _, red := range ret {
					log.Printf("")
					if red.Location == env.Location && red.Sku == vm.VmSku {
						match = true
					}
				}
				if !match {
					ret = append(ret, vmReduction{vm.VmSku, env.Location})
				}
			}
		}
	}
	return ret
}

func (c Config) PriceVms() ([]VmPrice, error) {
	vms := c.ReduceVms()
	vmp := make([]VmPrice, len(vms))
	for i, vm := range vms {
		log.Printf("Making VM API call %d of %d", i+1, len(vms))
		vmp[i].Currency = c.Currency
		vmp[i].Location = vm.Location
		vmp[i].VmSku = vm.Sku
		/* Format String */
		s1 := apiVmPriceString(c.Currency, vmp[i].Location, vmp[i].VmSku)
		log.Println(s1)
		resp, err := http.Get(s1)
		if err != nil {
			log.Print("Problem getting response from API.")
			return vmp, err
		}
		defer resp.Body.Close()
		ar := ApiResponse{}
		jdec := json.NewDecoder(resp.Body)
		err = jdec.Decode(&ar)
		if err != nil {
			log.Print("Problem with JSON decoder")
			return vmp, err
		}

		/* Ensure that something came out of the call */
		if ar.Count < 1 {
			return vmp, fmt.Errorf("API response contained no item, perhaps the SKU is wrong")
		}

		/* iterate and switch over the ar to find hourly, 1yrRi and 3yrRi costs */
		for _, v := range ar.Items {

			/* Don't use the Windows version */
			if v.Type == "Consumption" && !strings.Contains(v.ProductName, "Windows") {
				vmp[i].PaygHrRate = v.UnitPrice
			} else if v.Type == "Reservation" && v.ReservationTerm == "1 Year" {
				vmp[i].OneYrRi = v.UnitPrice
			} else if v.Type == "Reservation" && v.ReservationTerm == "3 Years" {
				vmp[i].ThreeYrRi = v.UnitPrice
			}
		}
	}
	return vmp, nil
}

func apiVmPriceString(c, l, s string) string {
	s1 := fmt.Sprintf("%s?api-version=%s&CurrencyCode=%s&$filter=armRegionName eq '%s' and serviceFamily eq 'Compute' and serviceName eq 'Virtual Machines' and skuName eq '%s'", ApiUrl, ApiPreview, c, l, s)
	/* Need url encoding */
	/* Due to MS API weirdness, can't use net/url */
	/* but the following does what is needed */
	s1 = strings.Replace(s1, " ", "%20", -1)
	s1 = strings.Replace(s1, "%3D", "=", -1)
	return s1

}
