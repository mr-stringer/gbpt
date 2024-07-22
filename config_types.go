package main

import (
	"fmt"
	"slices"
)

const blank string = ""

// Config is a struct that represents the configuration supplied by the user
type Config struct {
	Currency        string
	Applications    []Application
	StorageProfiles []StorageProfile
}

// Application is a struct that represents an application configuration as
// supplied by the user
type Application struct {
	Name         string
	Environments []Environment
}

// Environment is a struct that represent an application environment as supplied
// by the user
type Environment struct {
	Name     string
	Location string
	Phase    uint
	VMs      []Vm
}

// Vm is a struct that represents a Virtual Machine within an application
// environment
type Vm struct {
	Name  string
	Qty   uint
	VmSku string
	// The following bit has some hazards.  Consumption MUST be set correctly and the use of RiTermYears and PaygHours are mutually exclusive
	// If the config is mismatched the program will quit
	Consumption    string //must be ri or payg
	RiTermYears    int    // must be 1 or 3, if not populated it will return fine 0 which is fine if PAYG is to be used
	PaygHours      int    // must be a positive integer up to 730, will be populated is 0 if not
	StorageProfile string // this string must match a named storage config in the StorageConfigurations array in the config

}

// StorageProfile is a struct that represents a one or more disks that will be
// that be applied to one or more VMs
type StorageProfile struct {
	Name  string
	Disks []Disk
}

// Disk is a struct the represents a disk type that will be applied by a
// StorageProfile
type Disk struct {
	Name string // a friendly label
	Type string // must match pssd or pssd_v2
	Qty  uint   // The number of this disk type
	Size uint   // The size of the disk in GiB - required for both types
	Iops uint   // required for pssd_v2 disks only, ignored if set for pssd but should raise a warning
	MBs  uint   // required for pssd_v2 disks only, ignored if set for pssd but should raise a warning
}

// PriceLine represents all the information to price a single item in the config
type PriceLine struct {
	Application string
	Environment string
	Location    string
	Item        string
	Qty         uint
	UnitPrice   float32
	LinePrice   float32
}

// String outputs a human readable representation of the type PriceLine
func (pl PriceLine) String() string {
	return fmt.Sprintf("Application: %s, Environment:%s, Location:%s, Item:%s, Qty:%d, UnitPrice:%0.2f, TotalPrice:%0.2f",
		pl.Application, pl.Environment, pl.Location, pl.Item, pl.Qty, pl.UnitPrice, pl.LinePrice)
}

// CvsString outputs a string used to populate the CSV config file
func (pl PriceLine) CsvString() string {
	return fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%d\",\"%0.2f\",\"%0.2f\"",
		pl.Application, pl.Environment, pl.Location, pl.Item, pl.Qty, pl.UnitPrice, pl.LinePrice)
}

// csvHeader outputs the String header for the CSV file
func csvHeader(currency string) string {
	return fmt.Sprintf("\"Application\",\"Environment\",\"Location\",\"Description\",\"Qty\",\"UnitPrice\",\"LinePrice\"\"Currency:%s\"", currency)
}

// Print outputs a human readable representation of the Config type
func (c Config) Print() {
	fmt.Printf("Currency: %s\n", c.Currency)
	fmt.Printf("Number of Applications: %d\n", len(c.Applications))
	for i, app := range c.Applications {
		fmt.Printf("Application_%d Name: %s\n", i, app.Name)
		fmt.Printf("  Application_%d Number of Environments: %d\n", i, len(app.Environments))
		for ei, env := range app.Environments {
			fmt.Printf("%2sEnvironment_%d Name: %s\n", blank, ei, env.Name)
			fmt.Printf("%2sEnvironment_%d Location: %s\n", blank, ei, env.Location)
			fmt.Printf("%2sEnvironment_%d Phase: %d\n", blank, ei, env.Phase)
			fmt.Printf("%2sEnvironment_%d Number of VMs: %d\n", blank, ei, len(env.VMs))
			for vi, vm := range env.VMs {
				fmt.Printf("%4sVM_%d Name: %s\n", blank, vi, vm.Name)
				fmt.Printf("%4sVM_%d Qty: %d\n", blank, vi, vm.Qty)
				fmt.Printf("%4sVM_%d VmSku: %s\n", blank, vi, vm.VmSku)
				fmt.Printf("%4sVM_%d Consumption Type: %s\n", blank, vi, vm.Consumption)
				if vm.Consumption == "ri" {
					fmt.Printf("%4sVM_%d RiTermYears: %d\n", blank, vi, vm.RiTermYears)
				} else {
					fmt.Printf("%4sVM_%d PaygHours: %d\n", blank, vi, vm.PaygHours)
				}
				fmt.Printf("%4sVM_%d StorageProfile: %s\n", blank, vi, vm.StorageProfile)
			}
		}
	}
	fmt.Printf("Number of StorageProfiles: %d\n", len(c.StorageProfiles))
	for i, sp := range c.StorageProfiles {
		fmt.Printf("%2sStorageProfile_%d Name: %s\n", blank, i, sp.Name)
		fmt.Printf("%2sStorageProfile_%d Number of Disks: %d\n", blank, i, len(sp.Disks))
		for di, disk := range sp.Disks {
			fmt.Printf("%4sDisk_%d Name: %s\n", blank, di, disk.Name)
			fmt.Printf("%4sDisk_%d Type: %s\n", blank, di, disk.Type)
			fmt.Printf("%4sDisk_%d Size(GiB): %d\n", blank, di, disk.Size)
			if disk.Type == "pssd_v2" {
				fmt.Printf("%4sDisk_%d Iops: %d\n", blank, di, disk.Iops)
				fmt.Printf("%4sDisk_%d MBs: %d\n", blank, di, disk.MBs)
			}

		}
	}
}

// Validate ensures that the user input is valid. It will check all aspects of
// the config and report on all errors.
func (c Config) Validate() ([]string, error) {
	/* Ensure currency is correct */
	eStrings := []string{}
	if !slices.Contains(supCur, c.Currency) {
		eStrings = append(eStrings, fmt.Sprintf("the value %s is not a supported currency", c.Currency))
	}
	/* Ensure that the Application slice is populated */
	if len(c.Applications) < 1 {
		eStrings = append(eStrings, "the Applications slice contains no elements")
	}
	/* Ensure that the StorageProfile Array is populated */
	if len(c.StorageProfiles) < 1 {
		eStrings = append(eStrings, "the StorageProfile slice contains no elements")
	}
	/* Start the traversal of the StorageProfiles (do these first as the apps depend on them) */
	for si, sp := range c.StorageProfiles {
		/*Ensure storage profile is named */
		if sp.Name == "" {
			eStrings = append(eStrings, fmt.Sprintf("StorageProfile:%d is not named", si))
		}
		if len(sp.Disks) == 0 {
			eStrings = append(eStrings, fmt.Sprintf("StorageProfile:%d slice contains no elements", si))
		}
		/* Check disks */
		for di, disk := range sp.Disks {
			/* Ensure disk has a name */
			if disk.Name == "" {
				eStrings = append(eStrings, fmt.Sprintf("StorageProfile:%d Disk:%d is not named", si, di))
			}
			if disk.Type != "pssd" && disk.Type != "pssd_v2" {

				eStrings = append(eStrings, fmt.Sprintf("StorageProfile:%d Disk:%d type set to '%s' but only 'pssd' and 'pssd_v2' are allowed", si, di, disk.Type))
			}
			/* Ensure disk has Qty */
			if disk.Qty < 1 {
				eStrings = append(eStrings, fmt.Sprintf("StorageProfile:%d Disk:%d Qty is zero but must be 1 or more.", si, di))
			}
			/* Checks for pssd disks*/
			if disk.Type == "pssd" {
				/* Check size */
				if disk.Size < 1 || disk.Size > (32*1024)-1 {
					eStrings = append(eStrings, fmt.Sprintf("StorageProfile:%d Disk:%d size set to '%d' but pssd devices must be between 1 and %d GiB (inclusive)", si, di, disk.Size, (32*1024)-1))
				}
			}
			/* Check for pssd_v2 :| */
			if disk.Type == "pssd_v2" {
				if disk.Size < 1 || disk.Size > (64*1024)-1 {
					eStrings = append(eStrings, fmt.Sprintf("StorageProfile:%d Disk:%d size set to '%d' but pssd_v2 devices must be between 1 and %d GiB (inclusive)", si, di, disk.Size, (64*1024)-1))
				}
				if disk.Iops%500 != 0 {
					eStrings = append(eStrings, fmt.Sprintf("StorageProfile:%d Disk:%d IOPS set to '%d' but must be a multiple of '500'", si, di, disk.Iops))
				}
				/* Ensure IOPS are in range */
				if disk.Iops < 3000 || disk.Iops > 80000 {
					eStrings = append(eStrings, fmt.Sprintf("StorageProfile:%d Disk:%d IOPS set to '%d' but must be between '3000' and '80000'", si, di, disk.Iops))
				}
				/* Find max iops for disk size */
				var maxIops (uint) = 3000
				if disk.Size > 6 && disk.Size <= 160 {
					maxIops = (disk.Size) * 500
				} else if disk.Size > 160 {
					maxIops = 80000
				}
				/* Ensure the requested IOPS are possible for requested disk size*/
				if disk.Iops > maxIops {
					eStrings = append(eStrings, fmt.Sprintf("StorageProfile:%d Disk:%d Iops set to '%d' but maximum IOPS for a %dGiB is %d", si, di, disk.Iops, disk.Size, maxIops))
				}
				/* Ensure that MBs are in range */
				if disk.MBs < 125 || disk.MBs > 1200 {
					eStrings = append(eStrings, fmt.Sprintf("StorageProfile:%d Disk:%d MBs set to '%d' but must be between '125' and '1200' inclusive", si, di, disk.MBs))
				}
				/* Ensure MBs are possible for requested IOPS */
				if disk.MBs > disk.Iops/4 {
					eStrings = append(eStrings, fmt.Sprintf("StorageProfile:%d Disk:%d MBs set to '%d' but maximum MBs for %d IOPS is %d", si, di, disk.MBs, disk.Iops, disk.Iops/4))
				}
			}
		}
	}
	/* Start the traversal of the Application slice */
	for ai, app := range c.Applications {
		/* Ensure application is named */
		if app.Name == "" {
			eStrings = append(eStrings, fmt.Sprintf("Application:%d is not named", ai))
		}
		/* Ensure that Environments is populated */
		if len(app.Environments) < 1 {
			eStrings = append(eStrings, fmt.Sprintf("Application:%d Environment slice contains no elements", ai))
		}
		for ei, env := range app.Environments {
			/* Ensure that Environment is named */
			if env.Name == "" {
				eStrings = append(eStrings, fmt.Sprintf("Application:%d Environment:%d is not named", ai, ei))
			}
			if !slices.Contains(supReg, env.Location) {
				eStrings = append(eStrings, fmt.Sprintf("Application:%d Environment:%d location '%s' is not valid", ai, ei, env.Location))
			}
			/* Phase will default to zero, we don't need to check it */
			/* Check VM length */
			if len(env.VMs) < 1 {
				eStrings = append(eStrings, fmt.Sprintf("Application:%d Environment:%d No VMs found", ai, ei))
			}
			/* Now check VMs */
			for vi, vm := range env.VMs {
				/* Ensure VMs are named*/
				if vm.Name == "" {
					eStrings = append(eStrings, fmt.Sprintf("Application:%d Environment:%d VM:%d is not named", ai, ei, vi))
				}
				/* Ensure Qty is not 0 */
				if vm.Qty < 1 {
					eStrings = append(eStrings, fmt.Sprintf("Application:%d Environment:%d VM:%d Qty is zero", ai, ei, vi))
				}
				/* Too many VM SKUs to check, better for the API calls to fail if the SKU is incorrect*/
				/* if ri, ensure that term is correctly set */
				if vm.Consumption == "ri" {
					if vm.RiTermYears != 1 && vm.RiTermYears != 3 {
						eStrings = append(eStrings, fmt.Sprintf("Application:%d Environment:%d VM:%d RiTerms set to '%d' but must be either '1' or '3'", ai, ei, vi, vm.RiTermYears))
					}
				}
				if vm.Consumption == "payg" {
					if vm.PaygHours < 1 || vm.PaygHours > 720 {
						eStrings = append(eStrings, fmt.Sprintf("Application:%d Environment:%d VM:%d PaygHours set to '%d' but must be an integer between 1 and 720", ai, ei, vi, vm.PaygHours))
					}
				}
				if vm.Consumption != "ri" && vm.Consumption != "payg" {
					eStrings = append(eStrings, fmt.Sprintf("Application:%d Environment:%d VM:%d Consumption set to '%s' but must be either 'payg' or 'ri'", ai, ei, vi, vm.Consumption))
				}
				/* Get a slice of Storage Profile names */
				spNames := make([]string, len(c.StorageProfiles))
				/* populate slice */
				for k, v := range c.StorageProfiles {
					spNames[k] = v.Name
				}
				if !slices.Contains(spNames, vm.StorageProfile) {
					eStrings = append(eStrings, fmt.Sprintf("Application:%d Environment:%d VM:%d StorageProfile '%s' not found, check configuration", ai, ei, vi, vm.StorageProfile))
				}
			}
		}
	}

	if len(eStrings) > 0 {
		if len(eStrings) == 1 {
			return eStrings, fmt.Errorf("1 configuration error was found")
		} else {
			return eStrings, fmt.Errorf("%d configuration errors were found", len(eStrings))
		}
	}

	return eStrings, nil
}

/******************************************************************************/
/* Notes on SSD v2 disks                                                      */
/* --- Size ---                                                               */
/* SSD v2 disks are available from 1GiB to 64TiB                              */
/* --- IOPS ---                                                               */
/* All disks are deployed with a baseline of 3,000 IOPS                       */
/* Maximum IOPS per disk is 80,000                                            */
/* Above capacity reaches 6GiB, an additional 500 IOPS per GB is supported    */
/* per GiB. A disk requiring 80,000 IOPS must be be at least 160GiB.          */
/* --- Throughput                                                             */
/* All disks are deployed with 125MiB/s of throughput                         */
/* There is a ratio of 4K:1 between IOPS and throughput, or in other words,   */
/* For every 1000 IOPSs provisioned, a maximum of 250MiB/s can be provisioned */
/* --- PseudoCode ---                                                         */
/* //Check disk size                                                          */
/* if size <1GiB or >64TiB                                                    */
/*   fail - disk too big or too small                                         */
/* //ensure IOPS is a multiple of 500                                         */
/* if iops % 500 != 0                                                         */
/*   fail - iops must be a multiple of 500                                    */
/* //ensure IOPS in range                                                     */
/* if iops < 3000 or > 80000                                                  */
/*    fail - iops too big or too small                                        */
/* //Checks IOPS fit in disk size                                             */
/* achievable_iops = 3000                                                     */
/* if size > 6 and < 160                                                      */
/*   achievable_iops = (size -6 ) * 500                                       */
/* elseif size => 160                                                         */
/*   achievable_iops = 80000                                                  */
/* if iops > achievable_iops                                                  */
/*   fail - IOPS not achievable with this disk size                           */
/* // ensure throughput is in bounds                                          */
/* if throughput < 125 or > 1200                                              */
/*   fail - throughput too big or too small                                   */
/* // ensure throughput is achievable with IOPS                               /*
/* if throughput > (iops/4)                                                   /*
/*   fail - not enough iops provisioned                                       /*
/******************************************************************************/
