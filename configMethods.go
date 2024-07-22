package main

import (
	"fmt"
	"log"
	"math"
)

// PriceConfig is called after validation and builds the pricing and outputs a
// file.
func (c Config) PriceConfig() error {
	vmp, err := c.PriceVms()
	if err != nil {
		return err
	}

	dskp, err := c.PriceDisks()
	if err != nil {
		return err
	}

	for _, vm := range vmp {
		log.Print(vm.VmSku, ":", vm.Location)
	}

	for _, disk := range dskp {
		log.Print(disk.DiskType, ":", disk.Location)

	}

	pl := []PriceLine{}
	/* Iterate over the	config */
	for _, app := range c.Applications {
		log.Printf("Pricing application: %s", app.Name)
		for _, env := range app.Environments {
			log.Printf("Pricing %s:%s", app.Name, env.Name)
			for _, vm := range env.VMs {
				/* Search for VM price */
				for _, pr := range vmp {
					if pr.VmSku == vm.VmSku && pr.Location == env.Location {
						/* found one */
						log.Printf("Found match pr.VmSKU=%s", pr.VmSku)
						var itemLine string
						var unitPrice float32
						if vm.Consumption == "ri" {
							itemLine = fmt.Sprintf("Virtual Machine: %s Reserved Instance Term %d year/s", pr.VmSku, vm.RiTermYears)
							if vm.RiTermYears == 1 {
								unitPrice = pr.OneYrRi / 12
							} else {
								unitPrice = pr.ThreeYrRi / 36
							}
						} else {
							itemLine = fmt.Sprintf("Virtual Machine: %s Pay as You Go %d hours per month", vm.VmSku, vm.PaygHours)
							unitPrice = pr.PaygHrRate * float32(vm.PaygHours)
						}
						line := PriceLine{
							Application: app.Name,
							Environment: env.Name,
							Location:    env.Location,
							Item:        itemLine,
							Qty:         vm.Qty,
							UnitPrice:   float32(roundFloat(float64(unitPrice), 2)),
							LinePrice:   float32(roundFloat(float64(unitPrice), 2)) * float32(vm.Qty),
						}
						log.Print("Adding line")
						pl = append(pl, line)
						continue
					}
				}
				/* Find the storage profile that matches */
				for _, sp := range c.StorageProfiles {
					if vm.StorageProfile == sp.Name {
						/* We got a match */
						/* Now iterate over the disks */
						for _, disk := range sp.Disks {
							/* Iterate over the disk prices */
							for _, dp := range dskp {
								/* Now match the disks with the location */
								var dType uint
								if disk.Type == "pssd" {
									dType = pssd
								} else {
									dType = pssdv2
								}
								if env.Location == dp.Location && dType == dp.DiskType {
									/* we got out match */
									var itemLine string
									var unitPrice float32

									/* price pssd disk */
									if dp.DiskType == pssd {
										/* Check to see if pssd is correct one */
										pssdStr := getPssdFromSize(disk.Size)
										if pssdStr != dp.Pssd {
											log.Printf("%s!=%s", pssdStr, dp.Pssd)
											continue
										}
										itemLine = fmt.Sprintf("Premium SSD, Size:%dGiB (%s)", disk.Size, pssdStr)
										unitPrice = dp.Price

									} else {
										/* gotta be pssdv2 */
										itemLine = fmt.Sprintf("Premium SSD v2, Size:%dGiB, IOPS:%d, MBps:%d", disk.Size, disk.Iops, disk.MBs)
										unitPrice = float32(dp.GBs) * float32(disk.Size) * 730
										/* Ensure to remove included IOPS and Throughput */
										unitPrice += (float32(disk.Iops) - 3000.0) * dp.Iops * 730
										unitPrice += (float32(disk.MBs) - 125.0) * dp.MBps * 730
										log.Printf("pssdv2 disk capacity price %dGiB:%0.2f", disk.Size, float32(dp.GBs)*float32(disk.Size)*730)
										log.Printf("pssdv2 disk iops price %dIOPS:%0.2f", disk.Iops, (float32(disk.Iops)-3000.0)*dp.Iops*730)
										log.Printf("pssdv2 disk throughput price %dMB/s:%0.2f", disk.MBs, (float32(disk.MBs)-125.0)*dp.MBps*730)

									}
									line := PriceLine{
										Application: app.Name,
										Environment: env.Name,
										Location:    env.Location,
										Item:        itemLine,
										Qty:         disk.Qty,
										UnitPrice:   float32(roundFloat(float64(unitPrice), 2)),
										LinePrice:   float32(roundFloat(float64(unitPrice), 2)) * float32(disk.Qty),
									}
									log.Print("Adding line")
									pl = append(pl, line)
									continue
								}
							}
						}
					}
				}
			}
		}
	}

	csvData := []string{csvHeader(c.Currency)}
	var TotalCost float32
	for i := 0; i < len(pl); i++ {
		log.Print(pl[i].String())
		csvData = append(csvData, pl[i].CsvString())
		TotalCost += pl[i].LinePrice
	}

	csvData = append(csvData, "\"\",\"\",\"\",\"\",\"\",\"\",\"\"")
	csvData = append(csvData, fmt.Sprintf("\"\",\"\",\"\",\"\",\"\",\"Total Cost\",\"%0.2f\"", TotalCost))

	WriteToFile(&csvData, "/tmp/testOutput")

	return nil
}

// roundFloat sets the precision of a float to a specific level */
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
