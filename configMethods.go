package main

import (
	"fmt"
	"log/slog"
	"math"
)

// PriceConfig is called after validation and builds the pricing and outputs a
// file.
func (c Config) PriceConfig() error {
	slog.Info("Pricing configuration")
	vmp, err := c.PriceVms()
	if err != nil {
		return err
	}

	dskp, err := c.PriceDisks()
	if err != nil {
		return err
	}

	for _, vm := range vmp {
		slog.Debug("Reduced VM information", "SKU", vm.VmSku, "Location", vm.Location)
	}

	for _, disk := range dskp {
		slog.Debug("Reduced Disk information", "DiskType", disk.DiskType, "Location", disk.Location)

	}

	pl := []PriceLine{}
	slog.Info("Populating price list")
	/* Iterate over the	config */
	for _, app := range c.Applications {
		slog.Debug("Pricing application:", "Application", app.Name)
		for _, env := range app.Environments {
			slog.Debug("Pricing environment", "Application", app.Name, "Environment", env.Name)
			for _, vm := range env.VMs {
				/* Search for VM price */
				for _, pr := range vmp {
					slog.Debug("Searching reduction for matching VM")
					if pr.VmSku == vm.VmSku && pr.Location == env.Location {
						/* found one */
						slog.Debug("Found match", "VmSku", pr.VmSku, "Location", pr.Location)
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
						slog.Debug("PriceLine added for", "Application", app.Name, "Environment", app.Environments, "VmSku", pr.VmSku)
						pl = append(pl, line)
						continue
					}
				}
				/* Find the storage profile that matches */
				for _, sp := range c.StorageProfiles {
					slog.Debug("Pricing storage profile for VM", "Application", app.Name, "Environment", env.Name, "Vm", vm.Name, "StorageProfile", sp.Name)
					if vm.StorageProfile == sp.Name {
						/* We got a match */
						/* Now iterate over the disks */
						for _, disk := range sp.Disks {
							slog.Debug("StorageProfile match found")
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
									slog.Debug("Disk match found for disk", "Location", dp.Location, "Disk Type", dp.DiskType)
									var itemLine string
									var unitPrice float32

									/* price pssd disk */
									if dp.DiskType == pssd {
										/* Check to see if pssd is correct one */
										pssdStr := getPssdFromSize(disk.Size)
										if pssdStr != dp.Pssd {
											/*If it isn't a match, move to the next position in the loop*/
											continue
										}
										slog.Debug("Found ")
										itemLine = fmt.Sprintf("Premium SSD, Size:%dGiB (%s)", disk.Size, pssdStr)
										unitPrice = dp.Price

									} else {
										/* gotta be pssdv2 */
										itemLine = fmt.Sprintf("Premium SSD v2, Size:%dGiB, IOPS:%d, MBps:%d", disk.Size, disk.Iops, disk.MBs)
										unitPrice = float32(dp.GBs) * float32(disk.Size) * 730
										/* Ensure to remove included IOPS and Throughput */
										unitPrice += (float32(disk.Iops) - 3000.0) * dp.Iops * 730
										unitPrice += (float32(disk.MBs) - 125.0) * dp.MBps * 730
										slog.Debug("pssdv2 disk capacity price", "Size(GiB)", disk.Size, "Price", float32(dp.GBs)*float32(disk.Size)*730)
										slog.Debug("pssdv2 disk iops price", "Iops", disk.Iops, "Price", (float32(disk.Iops)-3000.0)*dp.Iops*730)
										slog.Debug("pssdv2 disk throughput price", "MB/s", disk.MBs, "Price", (float32(disk.MBs)-125.0)*dp.MBps*730)

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
									slog.Debug("Adding price line")
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
		slog.Debug(pl[i].String())
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
