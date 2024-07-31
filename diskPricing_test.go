package main

import (
	"reflect"
	"testing"
)

func TestConfig_ReduceDisks(t *testing.T) {
	rdSimple := []DiskPrice{
		{0, "P10", "uksouth", 0, 0, 0, 0, 0},
		{1, "", "uksouth", 0, 0, 0, 0, 0},
	}
	rdComplex := []DiskPrice{
		{0, "P10", "uksouth", 0, 0, 0, 0, 0},
		{1, "", "uksouth", 0, 0, 0, 0, 0},
		{0, "P6", "ukwest", 0, 0, 0, 0, 0},
		{0, "P20", "ukwest", 0, 0, 0, 0, 0},
	}

	type fields struct {
		Currency        string
		Applications    []Application
		StorageProfiles []StorageProfile
	}
	tests := []struct {
		name   string
		fields fields
		want   []DiskPrice
	}{
		{"GoodSimple", fields{"GBP", cfgGd01.Applications, cfgGd01.StorageProfiles}, rdSimple},
		{"GoodMultiSite", fields{"GBP", cfgGd02.Applications, cfgGd02.StorageProfiles}, rdComplex},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Currency:        tt.fields.Currency,
				Applications:    tt.fields.Applications,
				StorageProfiles: tt.fields.StorageProfiles,
			}
			if got := c.ReduceDisks(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.ReduceDisks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApiPssdPriceString(t *testing.T) {
	type args struct {
		c string
		l string
		p string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Basic", args{"GBP", "uksouth", "P6"}, "https://prices.azure.com/api/retail/prices?api-version=2023-01-01-preview&CurrencyCode=GBP&$filter=armRegionName%20eq%20'uksouth'%20and%20serviceFamily%20eq%20'Storage'%20and%20skuName%20eq%20'P6%20LRS'%20and%20productName%20eq%20'Premium%20SSD%20Managed%20Disks'%20and%20meterName%20eq%20'P6%20LRS%20Disk'%20and%20priceType%20eq%20'Consumption'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ApiPssdPriceString(tt.args.c, tt.args.l, tt.args.p); got != tt.want {
				t.Errorf("ApiPssdPriceString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApiPssdv2PriceString(t *testing.T) {
	type args struct {
		c string
		l string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Basic", args{"GBP", "ukwest"}, "https://prices.azure.com/api/retail/prices?api-version=2023-01-01-preview&CurrencyCode=GBP&$filter=armRegionName%20eq%20'ukwest'%20and%20serviceFamily%20eq%20'Storage'%20and%20priceType%20eq%20'Consumption'%20and%20productName%20eq%20'Azure%20Premium%20SSD%20v2'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ApiPssdv2PriceString(tt.args.c, tt.args.l); got != tt.want {
				t.Errorf("ApiPssdv2PriceString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_PriceDisks(t *testing.T) {
	type fields struct {
		Currency        string
		Applications    []Application
		StorageProfiles []StorageProfile
	}
	type args struct {
		apg apiGetter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []DiskPrice
		wantErr bool
	}{
		{"SimplePssd", fields{"GBP", simpleAppPaygConfig, simpleStorageProfilePssd}, args{mockGetterSimplePssd}, simpleStorageProfilePssdReturn, false},
		{"ApiCallFailsPssd", fields{"GBP", simpleAppPaygConfig, simpleStorageProfilePssd}, args{mockGetterFail}, []DiskPrice{}, true},
		{"ApiNoResultsPssd", fields{"GBP", simpleAppPaygConfig, simpleStorageProfilePssd}, args{mockGetterNoItems}, []DiskPrice{}, true},
		{"TooManyPssdReturned", fields{"GBP", simpleAppPaygConfig, simpleStorageProfilePssd}, args{mockGetterTooManyPssd}, []DiskPrice{}, true},
		{"SimplePssdv2", fields{"GBP", simpleAppPaygConfig, simpleStorageProfilePssdv2}, args{mockGetterPssdv2}, simpleStorageProfilePssdReturnv2, false},
		{"ApiNoResultsPssdv2", fields{"GBP", simpleAppPaygConfig, simpleStorageProfilePssdv2}, args{mockGetterPssdv2NoResult}, []DiskPrice{}, true},
		{"ApiCallFailsPssdv2", fields{"GBP", simpleAppPaygConfig, simpleStorageProfilePssdv2}, args{mockGetterPssdv2ApiCallFailed}, []DiskPrice{}, true},
		{"ZeroPricePssdv2", fields{"GBP", simpleAppPaygConfig, simpleStorageProfilePssdv2}, args{mockGetterPssdv2ZeroPrice}, []DiskPrice{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Currency:        tt.fields.Currency,
				Applications:    tt.fields.Applications,
				StorageProfiles: tt.fields.StorageProfiles,
			}
			got, err := c.PriceDisks(tt.args.apg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.PriceDisks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.PriceDisks() = %v, want %v", got, tt.want)
			}
		})
	}
}
