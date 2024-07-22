package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

// SpLoc (Storage Profile Location) is a struct used to help reduce calls to the
// Azure Retail Pricing API
type SpLoc struct {
	sp  string
	loc string
}

// ReduceDisks is a function that inspects the configuration in order to
// minimize API calls for disks. It returns an slice of the DiskPrice type which
// contains unique disk/location combinations.
func (c Config) ReduceDisks() []DiskPrice {
	/* The purpose here is to find all of the different disks in different */
	/* locations we need and reduce them */

	/* collect unique combinations of storage profiles and locations */
	sp1 := []SpLoc{}

	for _, sp := range c.StorageProfiles {
		for _, ap := range c.Applications {
			for _, env := range ap.Environments {
				for _, vm := range env.VMs {
					if sp.Name == vm.StorageProfile {
						var match = false
						for _, v := range sp1 {
							if v.sp == sp.Name && env.Location == v.loc {
								match = true
							}
						}
						if !match {
							sp1 = append(sp1, SpLoc{sp.Name, env.Location})
						}
					}
				}
			}
		}
	}
	/* Now we have the unique storage profile and location combos */
	/* creating unique disk/location combos should be easier */
	rd1 := []DiskPrice{}

	for _, sp := range c.StorageProfiles {
		for _, sp1 := range sp1 {
			if sp1.sp != sp.Name {
				/* if the storage profiles don't match do a skip */
				continue
			}
			for _, disk := range sp.Disks {
				dp := DiskPrice{}
				dp.Location = sp1.loc
				if disk.Type == "pssd" {
					dp.DiskType = pssd
					dp.Pssd = getPssdFromSize(disk.Size)
				} else {
					dp.DiskType = pssdv2
				}
				match := false
				/* Check for a match in the returning slice */
				for _, rd := range rd1 {
					if rd.Location == dp.Location && rd.DiskType == dp.DiskType && rd.Pssd == dp.Pssd {
						match = true
					}
				}
				if !match {
					rd1 = append(rd1, dp)
				}

			}
		}
	}
	slog.Debug("ReduceDisks", "Unique Disk Configuration Found", len(rd1))
	for _, rd := range rd1 {
		slog.Debug("Disk", "Disk Type", rd.DiskType, "PssdType", rd.Pssd, "Location", rd.Location)
	}

	return rd1
}

// PriceDisks retrieves disks prices from the Azure Retail Price API
func (c Config) PriceDisks() ([]DiskPrice, error) {
	dp := c.ReduceDisks()
	slog.Debug("PriceDisks", "Disks to price", len(dp))
	for i, disk := range dp {
		/*Price pssd */
		if disk.DiskType == pssd {
			s1 := ApiPssdPriceString(c.Currency, disk.Location, disk.Pssd)
			resp, err := http.Get(s1)
			if err != nil {
				return []DiskPrice{}, err
			}
			defer resp.Body.Close()
			ar := ApiResponse{}
			jdec := json.NewDecoder(resp.Body)
			err = jdec.Decode(&ar)
			if err != nil {
				return []DiskPrice{}, nil
			}

			if len(ar.Items) == 0 {
				return []DiskPrice{}, fmt.Errorf("the API response contained no items, maybe the API has changed")
			} else if len(ar.Items) > 1 {
				return []DiskPrice{}, fmt.Errorf("the API response contained more than one item, this is unexpected, maybe the API has changed")
			}
			dp[i].Price = ar.Items[0].RetailPrice
		}
		if disk.DiskType == pssdv2 {
			s1 := ApiPssdv2PriceString(c.Currency, disk.Location)
			resp, err := http.Get(s1)
			if err != nil {
				return []DiskPrice{}, nil
			}
			ar := ApiResponse{}
			defer resp.Body.Close()
			jdec := json.NewDecoder(resp.Body)
			err = jdec.Decode(&ar)
			if err != nil {
				return []DiskPrice{}, err
			}

			if len(ar.Items) == 0 {
				return []DiskPrice{}, fmt.Errorf("the API response contained no items, maybe the API has changed")

			}
			for _, item := range ar.Items {
				/* immediately discard if price is empty */
				switch {
				case item.RetailPrice == 0:
					continue
				case item.MeterName == "Premium LRS Provisioned IOPS":
					dp[i].Iops = item.RetailPrice
				case item.MeterName == "Premium LRS Provisioned Throughput (MBps)":
					dp[i].MBps = item.RetailPrice
				case item.MeterName == "Premium LRS Provisioned Capacity":
					dp[i].GBs = item.RetailPrice
				}
			}
			slog.Debug("PriceDisks price found", "Type", "pssdv2", "Location", disk.Location, "MB_Price", disk.GBs, "IOPS_Price", disk.Iops, "MB/s_Price", disk.MBps)

			if dp[i].Iops == 0 || dp[i].MBps == 0 || dp[i].GBs == 0 {
				return []DiskPrice{}, fmt.Errorf("was not able to to price all parts of pssdv2 disk")
			}
		}
	}
	return dp, nil
}

// ApiPssdPriceSting crates a sting that can be used to call the Azure Retail
// Pricing API to look up the price of a Premium SSD disk. It takes three
// string arguments, c (currency), l (location) and p (P disk type) and
// returns a string.
func ApiPssdPriceString(c, l, p string) string {
	s1 := fmt.Sprintf("%s?api-version=%s&CurrencyCode=%s&$filter=armRegionName eq '%s' and serviceFamily eq 'Storage' and skuName eq '%s LRS' and productName eq 'Premium SSD Managed Disks' and meterName eq '%s LRS Disk' and priceType eq 'Consumption'", ApiUrl, ApiPreview, c, l, p, p)
	/* Need url encoding */
	/* Due to MS API weirdness, can't use net/url */
	/* but the following does what is needed */
	s1 = strings.Replace(s1, " ", "%20", -1)
	s1 = strings.Replace(s1, "%3D", "=", -1)
	slog.Debug("ApiPssdPriceString", "url", s1)
	return s1
}

// ApiPssdv2PriceString creates a string that can be used to call the Azure
// Retail API to look up price of a Premium SSD v2 disk. It takes two arguments,
// c (currency), l (location) and returns a string.
func ApiPssdv2PriceString(c, l string) string {
	s1 := fmt.Sprintf("%s?api-version=%s&CurrencyCode=%s&$filter=armRegionName eq '%s' and serviceFamily eq 'Storage' and priceType eq 'Consumption' and productName eq 'Azure Premium SSD v2'", ApiUrl, ApiPreview, c, l)
	/* Need url encoding */
	/* Due to MS API weirdness, can't use net/url */
	/* but the following does what is needed */
	s1 = strings.Replace(s1, " ", "%20", -1)
	slog.Debug("ApiPssdv2PriceString", "url", s1)
	return s1
}
