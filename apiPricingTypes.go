package main

// ApiResponse is a struct that encapsulates the response object from the Azure
// Retail Pricing API
type ApiResponse struct {
	BillingCurrency    string `json:"BillingCurrency"`
	CustomerEntityId   string `json:"CustomerEntityId"`
	CustomerEntityType string `json:"CustomerEntityType"`
	NextPageLink       string `json:"NextPageLink"`
	Count              uint   `json:"Count"`
	Items              []Item `json:"Items"`
}

// Item is a struct that encapsulates the Item object found in the Azure Retail
// Pricing API response
type Item struct {
	CurrencyCode         string        `json:"currencyCode"`
	TierMinimumUnits     float32       `json:"tierMinimumUnits"`
	RetailPrice          float32       `json:"retailPrice"`
	UnitPrice            float32       `json:"unitPrice"`
	ArmRegionName        string        `json:"armRegionName"`
	Location             string        `json:"location"`
	EffectiveStartDate   string        `json:"effectiveStartDate"`
	MeterID              string        `json:"0084b086-37bf-4bee-b27f-6eb0f9ee4954"`
	MeterName            string        `json:"meterName"`
	ProductID            string        `json:"productId"`
	SkuId                string        `json:"skuId"`
	AvailabilityId       string        `json:"availabilityId"`
	ProductName          string        `json:"productName"`
	SkuName              string        `json:"skuName"`
	ServiceName          string        `json:"serviceName"`
	ServiceID            string        `json:"serviceId"`
	ServiceFamily        string        `json:"serviceFamily"`
	UnitOfMeasure        string        `json:"unitOfMeasure"`
	Type                 string        `json:"type"`
	IsPrimaryMeterRegion bool          `json:"isPrimaryMeterRegion"`
	ArmSkuName           string        `json:"armSkuName"`
	ReservationTerm      string        `json:"reservationTerm"`
	SavingsPlan          []SavingsPlan `json:"savingsPlan"`
}

// SavingsPlan is a struct that encapsulates the SavingsPlan object found in the
// Items object found in the Azure Retail Pricing API
type SavingsPlan struct {
	UnitPrice   float32 `json:"unitPrice"`
	RetailPrice float32 `json:"retailPrice"`
	Term        string  `json:"term"`
}
