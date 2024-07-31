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

func TestConfig_PriceVms(t *testing.T) {
	Applications := cfgGd01.Applications
	StorageProfiles := cfgGd01.StorageProfiles
	type fields struct {
		Currency        string
		Applications    []Application
		StorageProfiles []StorageProfile
	}
	type args struct {
		in0 apiGetter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []VmPrice
		wantErr bool
	}{
		{"ApiCallFailed", fields{"GBP", Applications, StorageProfiles}, args{mockGetterFail}, []VmPrice{}, true},
		{"NoItemsInResponse", fields{"GBP", Applications, StorageProfiles}, args{mockGetterNoItems}, []VmPrice{}, true},
		{"ApiSimplePaygVm", fields{"GBP", simpleAppPaygConfig, simpleStorageProfilePssd}, args{mockGetterSimpleAppPayg}, simpleAppPaygConfigReturn, false},
		{"ApiSimple1YrRiVm", fields{"GBP", simpleApp1YrRiConfig, simpleStorageProfilePssd}, args{mockGetterSimpleApp1YrRi}, simpleApp1YrRiConfigReturn, false},
		{"ApiSimple3YrRiVm", fields{"GBP", simpleApp3YrRiConfig, simpleStorageProfilePssd}, args{mockGetterSimpleApp3YrRi}, simpleApp3YrRiConfigReturn, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Currency:        tt.fields.Currency,
				Applications:    tt.fields.Applications,
				StorageProfiles: tt.fields.StorageProfiles,
			}
			got, err := c.PriceVms(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.PriceVms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.PriceVms() = %v, want %v", got, tt.want)
			}
		})
	}
}
