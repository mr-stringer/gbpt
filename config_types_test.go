package main

import (
	"fmt"
	"testing"
)

var cfgGd01 = Config{
	Currency: "GBP",
	Applications: []Application{
		{
			Name: "NFS",
			Environments: []Environment{
				{
					Name:     "Production",
					Location: "uksouth",
					Phase:    0,
					VMs: []Vm{
						{
							Name:           "NFS Cluster",
							Qty:            2,
							VmSku:          "Standard_ds2dsv5",
							Consumption:    "ri",
							RiTermYears:    3,
							StorageProfile: "NfsProd",
						},
						{
							Name:           "NFS Quorum",
							Qty:            1,
							VmSku:          "Standard_ds2dsv5",
							Consumption:    "payg",
							PaygHours:      500,
							StorageProfile: "NfsProd",
						},
					},
				},
			},
		},
	},
	StorageProfiles: []StorageProfile{
		{
			Name: "NfsProd",
			Disks: []Disk{
				{
					Name: "OS Disk",
					Type: "pssd",
					Qty:  1,
					Size: 128,
				},
				{
					Name: "Other Disk",
					Type: "pssd_v2",
					Qty:  1,
					Size: 1024,
					Iops: 5000,
					MBs:  500,
				},
			},
		},
	},
}

var cfgGd02 = Config{
	Currency: "GBP",
	Applications: []Application{
		{
			Name: "NFS",
			Environments: []Environment{
				{
					Name:     "Production",
					Location: "uksouth",
					Phase:    0,
					VMs: []Vm{
						{
							Name:           "NFS Cluster",
							Qty:            2,
							VmSku:          "Standard_ds2dsv5",
							Consumption:    "ri",
							RiTermYears:    3,
							StorageProfile: "NfsProd",
						},
						{
							Name:           "NFS Quorum",
							Qty:            1,
							VmSku:          "Standard_ds2dsv5",
							Consumption:    "payg",
							PaygHours:      500,
							StorageProfile: "NfsProd",
						},
					},
				}, {
					Name:     "Development",
					Location: "ukwest",
					Phase:    0,
					VMs: []Vm{
						{
							Name:           "NFS Server",
							Qty:            2,
							VmSku:          "Standard_ds2dsv5",
							Consumption:    "payg",
							PaygHours:      500,
							StorageProfile: "NfsNonProd",
						},
					},
				},
			},
		},
	},
	StorageProfiles: []StorageProfile{
		{
			Name: "NfsProd",
			Disks: []Disk{
				{
					Name: "OS Disk",
					Type: "pssd",
					Qty:  1,
					Size: 128,
				},
				{
					Name: "Other Disk",
					Type: "pssd_v2",
					Qty:  1,
					Size: 1024,
					Iops: 5000,
					MBs:  500,
				},
			},
		}, {
			Name: "NfsNonProd",
			Disks: []Disk{
				{
					Name: "OS Disk",
					Type: "pssd",
					Qty:  1,
					Size: 64,
				},
				{
					Name: "Same Size App Disk",
					Type: "pssd",
					Qty:  2,
					Size: 64,
				},

				{
					Name: "Other Disk",
					Type: "pssd",
					Qty:  1,
					Size: 512,
				},
			},
		},
	},
}

func TestConfig_Validate(t *testing.T) {
	/* Configure test vars */
	wrongCurrency := cfgGd01
	wrongCurrency.Currency = "Shillings"
	/* Empty App struct */
	noApplications := cfgGd01
	noApplications.Applications = nil
	noApplications.Applications = []Application{}
	/* Empty Storage Struct*/
	noStorage := cfgGd01
	noStorage.StorageProfiles = nil
	noStorage.StorageProfiles = []StorageProfile{}
	/* Unnamed Storage Profile */
	unnamedStorageProfile := cfgGd01
	unnamedStorageProfile.StorageProfiles = nil
	unnamedStorageProfile.StorageProfiles = append(unnamedStorageProfile.StorageProfiles, cfgGd01.StorageProfiles...)
	unnamedStorageProfile.StorageProfiles[0].Name = ""
	/* No Disks in profile */
	noDisksInProfile := cfgGd01
	noDisksInProfile.StorageProfiles = nil
	noDisksInProfile.StorageProfiles = append(noDisksInProfile.StorageProfiles, cfgGd01.StorageProfiles...)
	noDisksInProfile.StorageProfiles[0].Disks = []Disk{}
	/* Unnamed Disk */
	unnamedDiskInProfile := cfgGd01
	unnamedDiskInProfile.StorageProfiles = nil
	unnamedDiskInProfile.StorageProfiles = []StorageProfile{{Name: "NfsProd", Disks: []Disk{{Name: "", Type: "pssd", Qty: 1, Size: 128}}}}
	/* Unsupported Disk YType */
	unsupportedDiskType := cfgGd01
	unsupportedDiskType.StorageProfiles = nil
	unsupportedDiskType.StorageProfiles = []StorageProfile{{Name: "NfsProd", Disks: []Disk{{Name: "OS", Type: "hdd", Qty: 1, Size: 128}}}}
	/* No disk qty */
	NoDiskQtyDiskInProfile := cfgGd01
	NoDiskQtyDiskInProfile.StorageProfiles = nil
	NoDiskQtyDiskInProfile.StorageProfiles = []StorageProfile{{Name: "NfsProd", Disks: []Disk{{Name: "OS", Type: "pssd", Size: 128}}}}
	/* pssd zero qty */
	pssdZeroSize := cfgGd01
	pssdZeroSize.StorageProfiles = nil
	pssdZeroSize.StorageProfiles = []StorageProfile{{Name: "NfsProd", Disks: []Disk{{Name: "OS", Type: "pssd", Qty: 1}}}}
	/* pssd too big */
	pssdTooBig := cfgGd01
	pssdTooBig.StorageProfiles = nil
	pssdTooBig.StorageProfiles = []StorageProfile{{Name: "NfsProd", Disks: []Disk{{Name: "OS", Type: "pssd", Qty: 1, Size: 40 * 1024}}}}
	/* pssdv2 too zero size */
	pssv2dZeroSize := cfgGd01
	pssv2dZeroSize.StorageProfiles = nil
	pssv2dZeroSize.StorageProfiles = []StorageProfile{{Name: "NfsProd", Disks: []Disk{{Name: "OS", Type: "pssd_v2", Qty: 1}}}}
	/* pssdv2 too incorrect IOPS */
	pssv2dIncorrectIops := cfgGd01
	pssv2dIncorrectIops.StorageProfiles = nil
	pssv2dIncorrectIops.StorageProfiles = []StorageProfile{{Name: "NfsProd", Disks: []Disk{{Name: "OS", Type: "pssd_v2", Qty: 1, Size: 1024, Iops: 9999, MBs: 1000}}}}
	/* pssdv2 too low IOPS */
	pssv2dLowIops := cfgGd01
	pssv2dLowIops.StorageProfiles = nil
	pssv2dLowIops.StorageProfiles = []StorageProfile{{Name: "NfsProd", Disks: []Disk{{Name: "OS", Type: "pssd_v2", Qty: 1, Size: 1000, Iops: 1000, MBs: 1000}}}}
	/* pssdv2 too high IOPS */
	pssv2dHighIops := cfgGd01
	pssv2dHighIops.StorageProfiles = nil
	pssv2dHighIops.StorageProfiles = []StorageProfile{{Name: "NfsProd", Disks: []Disk{{Name: "OS", Type: "pssd_v2", Qty: 1, Size: 1000, Iops: 100000, MBs: 1000}}}}
	/* pssdv2 Not enough IOPs */
	pssv2dInsufficientIops := cfgGd01
	pssv2dInsufficientIops.StorageProfiles = nil
	pssv2dInsufficientIops.StorageProfiles = []StorageProfile{{Name: "NfsProd", Disks: []Disk{{Name: "OS", Type: "pssd_v2", Qty: 1, Size: 10, Iops: 50000, MBs: 125}}}}
	/* pssdv2 Not enough IOPs */
	pssv2dMBsTooLow := cfgGd01
	pssv2dMBsTooLow.StorageProfiles = nil
	pssv2dMBsTooLow.StorageProfiles = []StorageProfile{{Name: "NfsProd", Disks: []Disk{{Name: "OS", Type: "pssd_v2", Qty: 1, Size: 10, Iops: 3000, MBs: 30}}}}
	/* pssdv2 Not enough IOPs */
	pssv2dMBsTooHigh := cfgGd01
	pssv2dMBsTooHigh.StorageProfiles = nil
	pssv2dMBsTooHigh.StorageProfiles = []StorageProfile{{Name: "NfsProd", Disks: []Disk{{Name: "OS", Type: "pssd_v2", Qty: 1, Size: 1024, Iops: 50000, MBs: 3000}}}}
	/* pssdv2 Not enough IOPs */
	pssv2dMBsNotInRange := cfgGd01
	pssv2dMBsNotInRange.StorageProfiles = nil
	pssv2dMBsNotInRange.StorageProfiles = []StorageProfile{{Name: "NfsProd", Disks: []Disk{{Name: "OS", Type: "pssd_v2", Qty: 1, Size: 1024, Iops: 3000, MBs: 1200}}}}
	/* Ensure that applications are named */
	applicationNotNamed := cfgGd01
	applicationNotNamed.Applications = nil
	applicationNotNamed.Applications = append(applicationNotNamed.Applications, cfgGd01.Applications...)
	applicationNotNamed.Applications[0].Name = ""
	/* Ensure environments are populated */
	environmentNotPopulated := cfgGd01
	environmentNotPopulated.Applications = nil
	environmentNotPopulated.Applications = append(applicationNotNamed.Applications, cfgGd01.Applications...)
	environmentNotPopulated.Applications[0].Environments = nil
	/* Ensure that applications environments exist */
	environmentNotNamed := cfgGd01
	environmentNotNamed.Applications = nil
	environmentNotNamed.Applications = append(applicationNotNamed.Applications, cfgGd01.Applications...)
	environmentNotNamed.Applications[0].Environments = nil
	environmentNotNamed.Applications[0].Environments = append(applicationNotNamed.Applications[0].Environments, cfgGd01.Applications[0].Environments...)
	environmentNotNamed.Applications[0].Environments[0].Name = ""
	/* Ensure that location exist */
	environmentWrongLoc := cfgGd01
	environmentWrongLoc.Applications = nil
	environmentWrongLoc.Applications = append(applicationNotNamed.Applications, cfgGd01.Applications...)
	environmentWrongLoc.Applications[0].Environments = nil
	environmentWrongLoc.Applications[0].Environments = append(applicationNotNamed.Applications[0].Environments, cfgGd01.Applications[0].Environments...)
	environmentWrongLoc.Applications[0].Environments[0].Location = "bolivia"
	/* Ensure VMs are populated */
	vmsNotPopulated := cfgGd01
	vmsNotPopulated.Applications = nil
	vmsNotPopulated.Applications = append(applicationNotNamed.Applications, cfgGd01.Applications...)
	vmsNotPopulated.Applications[0].Environments = nil
	vmsNotPopulated.Applications[0].Environments = append(applicationNotNamed.Applications[0].Environments, cfgGd01.Applications[0].Environments...)
	vmsNotPopulated.Applications[0].Environments[0].VMs = nil
	/* Ensure VMs are named */
	vmNotNamed := cfgGd01
	vmNotNamed.Applications = nil
	vmNotNamed.Applications = append(applicationNotNamed.Applications, cfgGd01.Applications...)
	vmNotNamed.Applications[0].Environments = nil
	vmNotNamed.Applications[0].Environments = append(applicationNotNamed.Applications[0].Environments, cfgGd01.Applications[0].Environments...)
	vmNotNamed.Applications[0].Environments[0].VMs = nil
	vmNotNamed.Applications[0].Environments[0].VMs = append(vmsNotPopulated.Applications[0].Environments[0].VMs, cfgGd01.Applications[0].Environments[0].VMs...)
	vmNotNamed.Applications[0].Environments[0].VMs[0].Name = ""
	/* Ensure VMs Qty is set */
	vmNotQtySet := cfgGd01
	vmNotQtySet.Applications = nil
	vmNotQtySet.Applications = append(applicationNotNamed.Applications, cfgGd01.Applications...)
	vmNotQtySet.Applications[0].Environments = nil
	vmNotQtySet.Applications[0].Environments = append(applicationNotNamed.Applications[0].Environments, cfgGd01.Applications[0].Environments...)
	vmNotQtySet.Applications[0].Environments[0].VMs = nil
	vmNotQtySet.Applications[0].Environments[0].VMs = append(vmsNotPopulated.Applications[0].Environments[0].VMs, cfgGd01.Applications[0].Environments[0].VMs...)
	vmNotQtySet.Applications[0].Environments[0].VMs[0].Qty = 0
	/* Ensure Ri years is correct */
	vmRiTermIncorrect := cfgGd01
	vmRiTermIncorrect.Applications = nil
	vmRiTermIncorrect.Applications = append(applicationNotNamed.Applications, cfgGd01.Applications...)
	vmRiTermIncorrect.Applications[0].Environments = nil
	vmRiTermIncorrect.Applications[0].Environments = append(applicationNotNamed.Applications[0].Environments, cfgGd01.Applications[0].Environments...)
	vmRiTermIncorrect.Applications[0].Environments[0].VMs = nil
	vmRiTermIncorrect.Applications[0].Environments[0].VMs = append(vmsNotPopulated.Applications[0].Environments[0].VMs, cfgGd01.Applications[0].Environments[0].VMs...)
	vmRiTermIncorrect.Applications[0].Environments[0].VMs[0].RiTermYears = 2
	/* Ensure PAYG hours is correct */
	vmPaygHoursIncorrect := cfgGd01
	vmPaygHoursIncorrect.Applications = nil
	vmPaygHoursIncorrect.Applications = append(applicationNotNamed.Applications, cfgGd01.Applications...)
	vmPaygHoursIncorrect.Applications[0].Environments = nil
	vmPaygHoursIncorrect.Applications[0].Environments = append(applicationNotNamed.Applications[0].Environments, cfgGd01.Applications[0].Environments...)
	vmPaygHoursIncorrect.Applications[0].Environments[0].VMs = nil
	vmPaygHoursIncorrect.Applications[0].Environments[0].VMs = append(vmsNotPopulated.Applications[0].Environments[0].VMs, cfgGd01.Applications[0].Environments[0].VMs...)
	vmPaygHoursIncorrect.Applications[0].Environments[0].VMs[0].Consumption = "payg"
	vmPaygHoursIncorrect.Applications[0].Environments[0].VMs[0].PaygHours = 1000
	/* Ensure Consumption is correctly set  */
	vmConsumptionIncorrect := cfgGd01
	vmConsumptionIncorrect.Applications = nil
	vmConsumptionIncorrect.Applications = append(applicationNotNamed.Applications, cfgGd01.Applications...)
	vmConsumptionIncorrect.Applications[0].Environments = nil
	vmConsumptionIncorrect.Applications[0].Environments = append(applicationNotNamed.Applications[0].Environments, cfgGd01.Applications[0].Environments...)
	vmConsumptionIncorrect.Applications[0].Environments[0].VMs = nil
	vmConsumptionIncorrect.Applications[0].Environments[0].VMs = append(vmsNotPopulated.Applications[0].Environments[0].VMs, cfgGd01.Applications[0].Environments[0].VMs...)
	vmConsumptionIncorrect.Applications[0].Environments[0].VMs[0].Consumption = "monthly"

	type fields struct {
		Currency        string
		Applications    []Application
		StorageProfiles []StorageProfile
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"Good01", fields(cfgGd01), false},
		{"Incorrect Currency", fields(wrongCurrency), true},
		{"No Applications Defined", fields(noApplications), true},
		{"No Storage Profiles Defined", fields(noStorage), true},
		{"Unnamed Storage Profile", fields(unnamedStorageProfile), true},
		{"Unsupported Disk Type", fields(unsupportedDiskType), true},
		{"No Disks in Storage Profile", fields(noDisksInProfile), true},
		{"Unnamed Disk in Storage Profile", fields(unnamedDiskInProfile), true},
		{"Disk Quantity 0 in Storage Profile", fields(NoDiskQtyDiskInProfile), true},
		{"PSSD Zero Size", fields(pssdZeroSize), true},
		{"PSSD Too Big", fields(pssdTooBig), true},
		{"PSSD_V2 Zero Size", fields(pssv2dZeroSize), true},
		{"PSSD_V2 Incorrect IOPS Multiple", fields(pssv2dIncorrectIops), true},
		{"PSSD_V2 IOPS Too Low", fields(pssv2dLowIops), true},
		{"PSSD_V2 IOPS Too High", fields(pssv2dHighIops), true},
		{"PSSD_V2 Insufficient IOPS", fields(pssv2dInsufficientIops), true},
		{"PSSD_V2 Throughput Too Low", fields(pssv2dMBsTooLow), true},
		{"PSSD_V2 Throughput Too High", fields(pssv2dMBsTooHigh), true},
		{"PSSD_V2 Throughput Not in Range for IOPS", fields(pssv2dMBsNotInRange), true},
		{"Application Not Named", fields(applicationNotNamed), true},
		{"Environment Not Populated", fields(environmentNotPopulated), true},
		{"Environment Not Named", fields(environmentNotNamed), true},
		{"Incorrect Location", fields(environmentWrongLoc), true},
		{"VMs Not Populated", fields(vmsNotPopulated), true},
		{"VM Not Named", fields(vmNotNamed), true},
		{"VM Qty Not Set", fields(vmNotQtySet), true},
		{"VM Reserved Instance Years Incorrectly Set", fields(vmRiTermIncorrect), true},
		{"VM PAYG Hours Incorrectly Set", fields(vmPaygHoursIncorrect), true},
		{"VM Consumption Incorrectly Set", fields(vmConsumptionIncorrect), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Currency:        tt.fields.Currency,
				Applications:    tt.fields.Applications,
				StorageProfiles: tt.fields.StorageProfiles,
			}
			es, err := c.Validate()
			fmt.Println(es)
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			/* We don't care about precise the error messages */
			/*
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Config.Validate() = %v, want %v", got, tt.want)
				}
			*/
		})
	}
}

func TestConfig_Print(t *testing.T) {
	type fields struct {
		Currency        string
		Applications    []Application
		StorageProfiles []StorageProfile
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"SimpleConfig", fields(cfgGd01)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Currency:        tt.fields.Currency,
				Applications:    tt.fields.Applications,
				StorageProfiles: tt.fields.StorageProfiles,
			}
			c.Print()
		})
	}
}
