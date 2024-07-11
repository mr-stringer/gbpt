package main

import (
	"reflect"
	"testing"
)

var cfgVmReduce = Config{
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
					},
				},
				{
					Name:     "QA",
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
					},
				},
			},
		},
	},
}

var Development = Environment{
	Name:     "Development",
	Location: "ukwest",
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
	},
}

func TestConfig_ReduceVms(t *testing.T) {
	cfgReduceThreeToTwo := cfgVmReduce
	cfgReduceThreeToTwo.Applications = nil
	cfgReduceThreeToTwo.Applications = append(cfgReduceThreeToTwo.Applications, cfgVmReduce.Applications...)
	cfgReduceThreeToTwo.Applications[0].Environments = append(cfgReduceThreeToTwo.Applications[0].Environments, Development)

	type fields struct {
		Currency        string
		Applications    []Application
		StorageProfiles []StorageProfile
	}
	tests := []struct {
		name   string
		fields fields
		want   []vmReduction
	}{
		{"NoReduction", fields(cfgVmReduce), []vmReduction{{"Standard_ds2dsv5", "uksouth"}}},
		{"ReductionThreeToTwo", fields(cfgReduceThreeToTwo), []vmReduction{{"Standard_ds2dsv5", "uksouth"}, {"Standard_ds2dsv5", "ukwest"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Currency:        tt.fields.Currency,
				Applications:    tt.fields.Applications,
				StorageProfiles: tt.fields.StorageProfiles,
			}
			ret := c.ReduceVms()
			if !reflect.DeepEqual(ret, tt.want) {
				t.Errorf("Config.ReduceVms() = %v, want %v", ret, tt.want)
			}
		})
	}
}
