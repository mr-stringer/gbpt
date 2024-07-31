package main

/* This file largely contains testing structures for testing the pricing of   */
/* VMs and Disks. These structures get quite large so a dedicated file is     */
/* probably the best place for them.                                          */

import "fmt"

var simpleAppPaygConfig = []Application{
	{
		Name: "Test",
		Environments: []Environment{
			{
				Name:     "Test",
				Location: "uksouth",
				Phase:    0,
				VMs: []Vm{
					{
						Name:           "Test",
						Qty:            1,
						VmSku:          "Standard_D2ds_v5",
						Consumption:    "payg",
						PaygHours:      100,
						StorageProfile: "Test",
					},
				},
			},
		},
	},
}

var simpleAppPaygConfigReturn = []VmPrice{
	{
		VmSku:      "Standard_D2ds_v5",
		Location:   "uksouth",
		Currency:   "GBP",
		PaygHrRate: 0.103775,
	},
}

var simpleApp1YrRiConfig = []Application{
	{
		Name: "Test",
		Environments: []Environment{
			{
				Name:     "Test",
				Location: "uksouth",
				Phase:    0,
				VMs: []Vm{
					{
						Name:           "Test",
						Qty:            1,
						VmSku:          "Standard_D2ds_v5",
						Consumption:    "ri",
						RiTermYears:    1,
						StorageProfile: "Test",
					},
				},
			},
		},
	},
}

var simpleApp1YrRiConfigReturn = []VmPrice{
	{
		VmSku:    "Standard_D2ds_v5",
		Location: "uksouth",
		Currency: "GBP",
		OneYrRi:  500,
	},
}

var simpleApp3YrRiConfig = []Application{
	{
		Name: "Test",
		Environments: []Environment{
			{
				Name:     "Test",
				Location: "uksouth",
				Phase:    0,
				VMs: []Vm{
					{
						Name:           "Test",
						Qty:            1,
						VmSku:          "Standard_D2ds_v5",
						Consumption:    "ri",
						RiTermYears:    3,
						StorageProfile: "Test",
					},
				},
			},
		},
	},
}

var simpleApp3YrRiConfigReturn = []VmPrice{
	{
		VmSku:     "Standard_D2ds_v5",
		Location:  "uksouth",
		Currency:  "GBP",
		ThreeYrRi: 400,
	},
}

var simpleStorageProfilePssd = []StorageProfile{
	{
		Name: "Test",
		Disks: []Disk{
			{
				Name: "Test",
				Type: "pssd",
				Size: 64,
			},
		},
	},
}

var simpleStorageProfilePssdReturn = []DiskPrice{
	{
		DiskType: 0,
		Pssd:     "P6",
		Location: "uksouth",
		Price:    9.783261,
	},
}

var simpleStorageProfilePssdv2 = []StorageProfile{
	{
		Name: "Test",
		Disks: []Disk{
			{
				Name: "Test",
				Type: "pssdv2",
				Size: 100,
				Iops: 4000,
				MBs:  200,
			},
		},
	},
}

var simpleStorageProfilePssdReturnv2 = []DiskPrice{
	{
		DiskType: 1,
		Location: "uksouth",
		GBs:      0.000105,
		Iops:     0.000006,
		MBps:     0.000052,
	},
}

func mockGetterFail(url string) (ApiResponse, error) {
	return ApiResponse{}, fmt.Errorf("an error was found")
}

func mockGetterNoItems(url string) (ApiResponse, error) {
	ap := ApiResponse{
		BillingCurrency:    "GBP",
		CustomerEntityID:   "Default",
		CustomerEntityType: "Retail",
		NextPageLink:       "",
		Count:              0,
		Items:              []Item{},
	}
	return ap, nil
}

func mockGetterSimpleAppPayg(url string) (ApiResponse, error) {
	return ApiResponse{
		BillingCurrency:    "GBP",
		CustomerEntityID:   "Default",
		CustomerEntityType: "Retail",
		NextPageLink:       "",
		Count:              1,
		Items: []Item{
			{
				CurrencyCode:         "GBP",
				TierMinimumUnits:     0.0,
				RetailPrice:          0.103775,
				UnitPrice:            0.103775,
				ArmRegionName:        "uksouth",
				Location:             "UK South",
				EffectiveStartDate:   "2021-11-01T00:00:00Z",
				MeterID:              "5dcb2502-08d4-5169-8ce1-bb5a1ec26b99",
				MeterName:            "D2ds v5",
				ProductID:            "DZH318Z08MC6",
				SkuID:                "DZH318Z08MC6/0021",
				AvailabilityID:       "",
				ProductName:          "Virtual Machines Ddsv5 Series",
				SkuName:              "Standard_D2ds_v5",
				ServiceName:          "Virtual Machines",
				ServiceID:            "DZH313Z7MMC8",
				ServiceFamily:        "Compute",
				UnitOfMeasure:        "1 Hour",
				Type:                 "Consumption",
				IsPrimaryMeterRegion: true,
				ArmSkuName:           "Standard_D2ds_v5",
				ReservationTerm:      "",
				SavingsPlan:          []SavingsPlan{}},
		},
	}, nil
}

func mockGetterSimpleApp1YrRi(url string) (ApiResponse, error) {
	return ApiResponse{
		BillingCurrency:    "GBP",
		CustomerEntityID:   "Default",
		CustomerEntityType: "Retail",
		NextPageLink:       "",
		Count:              1,
		Items: []Item{
			{
				CurrencyCode:         "GBP",
				TierMinimumUnits:     0.0,
				RetailPrice:          500,
				UnitPrice:            500,
				ArmRegionName:        "uksouth",
				Location:             "UK South",
				EffectiveStartDate:   "2021-11-01T00:00:00Z",
				MeterID:              "5dcb2502-08d4-5169-8ce1-bb5a1ec26b99",
				MeterName:            "D2ds v5",
				ProductID:            "DZH318Z08MC6",
				SkuID:                "DZH318Z08MC6/0021",
				AvailabilityID:       "",
				ProductName:          "Virtual Machines Ddsv5 Series",
				SkuName:              "Standard_D2ds_v5",
				ServiceName:          "Virtual Machines",
				ServiceID:            "DZH313Z7MMC8",
				ServiceFamily:        "Compute",
				UnitOfMeasure:        "1 Hour",
				Type:                 "Reservation",
				IsPrimaryMeterRegion: true,
				ArmSkuName:           "Standard_D2ds_v5",
				ReservationTerm:      "1 Year",
				SavingsPlan:          []SavingsPlan{}},
		},
	}, nil
}

func mockGetterSimpleApp3YrRi(url string) (ApiResponse, error) {
	return ApiResponse{
		BillingCurrency:    "GBP",
		CustomerEntityID:   "Default",
		CustomerEntityType: "Retail",
		NextPageLink:       "",
		Count:              1,
		Items: []Item{
			{
				CurrencyCode:         "GBP",
				TierMinimumUnits:     0.0,
				RetailPrice:          400,
				UnitPrice:            400,
				ArmRegionName:        "uksouth",
				Location:             "UK South",
				EffectiveStartDate:   "2021-11-01T00:00:00Z",
				MeterID:              "5dcb2502-08d4-5169-8ce1-bb5a1ec26b99",
				MeterName:            "D2ds v5",
				ProductID:            "DZH318Z08MC6",
				SkuID:                "DZH318Z08MC6/0021",
				AvailabilityID:       "",
				ProductName:          "Virtual Machines Ddsv5 Series",
				SkuName:              "Standard_D2ds_v5",
				ServiceName:          "Virtual Machines",
				ServiceID:            "DZH313Z7MMC8",
				ServiceFamily:        "Compute",
				UnitOfMeasure:        "1 Hour",
				Type:                 "Reservation",
				IsPrimaryMeterRegion: true,
				ArmSkuName:           "Standard_D2ds_v5",
				ReservationTerm:      "3 Years",
				SavingsPlan:          []SavingsPlan{}},
		},
	}, nil
}

func mockGetterSimplePssd(url string) (ApiResponse, error) {

	ar := ApiResponse{
		BillingCurrency:    "GBP",
		CustomerEntityID:   "Default",
		CustomerEntityType: "Retail",
		NextPageLink:       "",
		Count:              1,
		Items: []Item{
			{
				CurrencyCode:         "GBP",
				TierMinimumUnits:     0,
				RetailPrice:          9.783261,
				UnitPrice:            9.783261,
				ArmRegionName:        "uksouth",
				Location:             "UK South",
				EffectiveStartDate:   "2018-05-01T00:00:00Z",
				MeterID:              "53c34f93-1bfc-4528-9edc-5579adab17b9",
				MeterName:            "P6 LRS Disk",
				ProductID:            "DZH318Z0BP04",
				SkuID:                "DZH318Z0BP04/0014",
				ProductName:          "Premium SSD Managed Disks",
				SkuName:              "P6 LRS",
				ServiceName:          "Storage",
				ServiceID:            "DZH317F1HKN0",
				ServiceFamily:        "Storage",
				UnitOfMeasure:        "1/Month",
				Type:                 "Consumption",
				IsPrimaryMeterRegion: true,
				ArmSkuName:           "Premium_SSD_Managed_Disk_P6",
			},
		},
	}
	return ar, nil
}

func mockGetterTooManyPssd(url string) (ApiResponse, error) {

	ar := ApiResponse{
		BillingCurrency:    "GBP",
		CustomerEntityID:   "Default",
		CustomerEntityType: "Retail",
		NextPageLink:       "",
		Count:              1,
		Items: []Item{
			{
				CurrencyCode: "GBP",
			},
			{
				CurrencyCode: "GBP",
			},
		},
	}
	return ar, nil
}

func mockGetterPssdv2(url string) (ApiResponse, error) {
	ar := ApiResponse{

		BillingCurrency:    "GBP",
		CustomerEntityID:   "Default",
		CustomerEntityType: "Retail",
		NextPageLink:       "",
		Count:              5,
		Items: []Item{
			{
				CurrencyCode:         "GBP",
				TierMinimumUnits:     0,
				RetailPrice:          0,
				UnitPrice:            0,
				ArmRegionName:        "ukwest",
				Location:             "UK West",
				EffectiveStartDate:   "2023-12-01T00:00:00Z",
				MeterID:              "01223859-c2d3-54a8-899b-f8e519a2c96a",
				MeterName:            "Premium LRS Provisioned Throughput (MBps)",
				ProductID:            "DZH318Z09V07",
				SkuID:                "DZH318Z09V07/000R",
				ProductName:          "Azure Premium SSD v2",
				SkuName:              "Premium LRS",
				ServiceName:          "Storage",
				ServiceID:            "DZH317F1HKN0",
				ServiceFamily:        "Storage",
				UnitOfMeasure:        "1/Hour",
				Type:                 "Consumption",
				IsPrimaryMeterRegion: true,
				ArmSkuName:           "",
			},
			{
				CurrencyCode:         "GBP",
				TierMinimumUnits:     125,
				RetailPrice:          0.000052,
				UnitPrice:            0.000052,
				ArmRegionName:        "ukwest",
				Location:             "UK West",
				EffectiveStartDate:   "2023-12-01T00:00:00Z",
				MeterID:              "01223859-c2d3-54a8-899b-f8e519a2c96a",
				MeterName:            "Premium LRS Provisioned Throughput (MBps)",
				ProductID:            "DZH318Z09V07",
				SkuID:                "DZH318Z09V07/000R",
				ProductName:          "Azure Premium SSD v2",
				SkuName:              "Premium LRS",
				ServiceName:          "Storage",
				ServiceID:            "DZH317F1HKN0",
				ServiceFamily:        "Storage",
				UnitOfMeasure:        "1/Hour",
				Type:                 "Consumption",
				IsPrimaryMeterRegion: true,
				ArmSkuName:           "",
			},
			{
				CurrencyCode:         "GBP",
				TierMinimumUnits:     0,
				RetailPrice:          0,
				UnitPrice:            0,
				ArmRegionName:        "ukwest",
				Location:             "UK West",
				EffectiveStartDate:   "2023-12-01T00:00:00Z",
				MeterID:              "91f08375-fa20-5d41-a8e6-c84d692a8c7d",
				MeterName:            "Premium LRS Provisioned IOPS",
				ProductID:            "DZH318Z09V07",
				SkuID:                "DZH318Z09V07/000R",
				ProductName:          "Azure Premium SSD v2",
				SkuName:              "Premium LRS",
				ServiceName:          "Storage",
				ServiceID:            "DZH317F1HKN0",
				ServiceFamily:        "Storage",
				UnitOfMeasure:        "1/Hour",
				Type:                 "Consumption",
				IsPrimaryMeterRegion: true,
				ArmSkuName:           "",
			},
			{
				CurrencyCode:         "GBP",
				TierMinimumUnits:     3000,
				RetailPrice:          0.000006,
				UnitPrice:            0.000006,
				ArmRegionName:        "ukwest",
				Location:             "UK West",
				EffectiveStartDate:   "2023-12-01T00:00:00Z",
				MeterID:              "91f08375-fa20-5d41-a8e6-c84d692a8c7d",
				MeterName:            "Premium LRS Provisioned IOPS",
				ProductID:            "DZH318Z09V07",
				SkuID:                "DZH318Z09V07/000R",
				ProductName:          "Azure Premium SSD v2",
				SkuName:              "Premium LRS",
				ServiceName:          "Storage",
				ServiceID:            "DZH317F1HKN0",
				ServiceFamily:        "Storage",
				UnitOfMeasure:        "1/Hour",
				Type:                 "Consumption",
				IsPrimaryMeterRegion: true,
				ArmSkuName:           "",
			},
			{
				CurrencyCode:         "GBP",
				TierMinimumUnits:     0,
				RetailPrice:          0.000105,
				UnitPrice:            0.000105,
				ArmRegionName:        "ukwest",
				Location:             "UK West",
				EffectiveStartDate:   "2023-08-01T00:00:00Z",
				MeterID:              "f7edead9-c427-56b1-af63-3626b41308d5",
				MeterName:            "Premium LRS Provisioned Capacity",
				ProductID:            "DZH318Z09V07",
				SkuID:                "DZH318Z09V07/000R",
				ProductName:          "Azure Premium SSD v2",
				SkuName:              "Premium LRS",
				ServiceName:          "Storage",
				ServiceID:            "DZH317F1HKN0",
				ServiceFamily:        "Storage",
				UnitOfMeasure:        "1 GiB/Hour",
				Type:                 "Consumption",
				IsPrimaryMeterRegion: true,
				ArmSkuName:           "",
			},
		},
	}
	return ar, nil
}

func mockGetterPssdv2NoResult(url string) (ApiResponse, error) {
	ar := ApiResponse{

		BillingCurrency:    "GBP",
		CustomerEntityID:   "Default",
		CustomerEntityType: "Retail",
		NextPageLink:       "",
		Count:              0,
		Items:              []Item{},
	}
	return ar, nil
}

func mockGetterPssdv2ApiCallFailed(url string) (ApiResponse, error) {
	return ApiResponse{}, fmt.Errorf("ApiCallFailed")
}

func mockGetterPssdv2ZeroPrice(url string) (ApiResponse, error) {
	ar := ApiResponse{

		BillingCurrency:    "GBP",
		CustomerEntityID:   "Default",
		CustomerEntityType: "Retail",
		NextPageLink:       "",
		Count:              5,
		Items: []Item{
			{
				CurrencyCode:         "GBP",
				TierMinimumUnits:     0,
				RetailPrice:          0,
				UnitPrice:            0,
				ArmRegionName:        "ukwest",
				Location:             "UK West",
				EffectiveStartDate:   "2023-12-01T00:00:00Z",
				MeterID:              "01223859-c2d3-54a8-899b-f8e519a2c96a",
				MeterName:            "Premium LRS Provisioned Throughput (MBps)",
				ProductID:            "DZH318Z09V07",
				SkuID:                "DZH318Z09V07/000R",
				ProductName:          "Azure Premium SSD v2",
				SkuName:              "Premium LRS",
				ServiceName:          "Storage",
				ServiceID:            "DZH317F1HKN0",
				ServiceFamily:        "Storage",
				UnitOfMeasure:        "1/Hour",
				Type:                 "Consumption",
				IsPrimaryMeterRegion: true,
				ArmSkuName:           "",
			},
			{
				CurrencyCode:         "GBP",
				TierMinimumUnits:     125,
				RetailPrice:          0.000000,
				UnitPrice:            0.000000,
				ArmRegionName:        "ukwest",
				Location:             "UK West",
				EffectiveStartDate:   "2023-12-01T00:00:00Z",
				MeterID:              "01223859-c2d3-54a8-899b-f8e519a2c96a",
				MeterName:            "Premium LRS Provisioned Throughput (MBps)",
				ProductID:            "DZH318Z09V07",
				SkuID:                "DZH318Z09V07/000R",
				ProductName:          "Azure Premium SSD v2",
				SkuName:              "Premium LRS",
				ServiceName:          "Storage",
				ServiceID:            "DZH317F1HKN0",
				ServiceFamily:        "Storage",
				UnitOfMeasure:        "1/Hour",
				Type:                 "Consumption",
				IsPrimaryMeterRegion: true,
				ArmSkuName:           "",
			},
			{
				CurrencyCode:         "GBP",
				TierMinimumUnits:     0,
				RetailPrice:          0,
				UnitPrice:            0,
				ArmRegionName:        "ukwest",
				Location:             "UK West",
				EffectiveStartDate:   "2023-12-01T00:00:00Z",
				MeterID:              "91f08375-fa20-5d41-a8e6-c84d692a8c7d",
				MeterName:            "Premium LRS Provisioned IOPS",
				ProductID:            "DZH318Z09V07",
				SkuID:                "DZH318Z09V07/000R",
				ProductName:          "Azure Premium SSD v2",
				SkuName:              "Premium LRS",
				ServiceName:          "Storage",
				ServiceID:            "DZH317F1HKN0",
				ServiceFamily:        "Storage",
				UnitOfMeasure:        "1/Hour",
				Type:                 "Consumption",
				IsPrimaryMeterRegion: true,
				ArmSkuName:           "",
			},
			{
				CurrencyCode:         "GBP",
				TierMinimumUnits:     3000,
				RetailPrice:          0.000006,
				UnitPrice:            0.000006,
				ArmRegionName:        "ukwest",
				Location:             "UK West",
				EffectiveStartDate:   "2023-12-01T00:00:00Z",
				MeterID:              "91f08375-fa20-5d41-a8e6-c84d692a8c7d",
				MeterName:            "Premium LRS Provisioned IOPS",
				ProductID:            "DZH318Z09V07",
				SkuID:                "DZH318Z09V07/000R",
				ProductName:          "Azure Premium SSD v2",
				SkuName:              "Premium LRS",
				ServiceName:          "Storage",
				ServiceID:            "DZH317F1HKN0",
				ServiceFamily:        "Storage",
				UnitOfMeasure:        "1/Hour",
				Type:                 "Consumption",
				IsPrimaryMeterRegion: true,
				ArmSkuName:           "",
			},
			{
				CurrencyCode:         "GBP",
				TierMinimumUnits:     0,
				RetailPrice:          0.000105,
				UnitPrice:            0.000105,
				ArmRegionName:        "ukwest",
				Location:             "UK West",
				EffectiveStartDate:   "2023-08-01T00:00:00Z",
				MeterID:              "f7edead9-c427-56b1-af63-3626b41308d5",
				MeterName:            "Premium LRS Provisioned Capacity",
				ProductID:            "DZH318Z09V07",
				SkuID:                "DZH318Z09V07/000R",
				ProductName:          "Azure Premium SSD v2",
				SkuName:              "Premium LRS",
				ServiceName:          "Storage",
				ServiceID:            "DZH317F1HKN0",
				ServiceFamily:        "Storage",
				UnitOfMeasure:        "1 GiB/Hour",
				Type:                 "Consumption",
				IsPrimaryMeterRegion: true,
				ArmSkuName:           "",
			},
		},
	}
	return ar, nil
}
