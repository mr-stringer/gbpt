package main

/* Preview version of the API currently needed for pssd_v2 prices*/

// URL for the Azure Retail Price API
const ApiUrl = "https://prices.azure.com/api/retail/prices"

// Preview Version used by gbpt
const ApiPreview = "2023-01-01-preview"

// sapCur is a slice of strings the represent supported currencies
var supCur = []string{"USD", "AUD", "BRL",
	"CAD", "CHF", "CNY",
	"DKK", "EUR", "GBP",
	"INR", "JPY", "KRW",
	"NOK", "NZD", "RUB",
	"SEK", "TWD"}

/* list of regions generated with az */
/* az account list-locations --query "[].{Name:name}" -o table */

// supReg is a slice of strings that represent the supported Azure regions
var supReg = []string{"eastus", "southcentralus", "westus2", "westus3",
	"australiaeast", "southeastasia", "northeurope", "swedencentral",
	"uksouth", "westeurope", "centralus", "southafricanorth",
	"centralindia", "eastasia", "japaneast", "koreacentral",
	"canadacentral", "francecentral", "germanywestcentral", "italynorth",
	"norwayeast", "polandcentral", "spaincentral", "switzerlandnorth",
	"mexicocentral", "uaenorth", "brazilsouth", "israelcentral",
	"qatarcentral", "centralusstage", "eastusstage", "eastus2stage",
	"northcentralusstage", "southcentralusstage", "westusstage", "westus2stage",
	"asia", "asiapacific", "australia", "brazil",
	"canada", "europe", "france", "germany",
	"global", "india", "israel", "italy",
	"japan", "korea", "newzealand", "norway",
	"poland", "qatar", "singapore", "southafrica",
	"sweden", "switzerland", "uae", "uk",
	"unitedstates", "unitedstateseuap", "eastasiastage", "southeastasiastage",
	"brazilus", "eastus2", "eastusstg", "northcentralus",
	"westus", "japanwest", "jioindiawest", "centraluseuap",
	"eastus2euap", "westcentralus", "southafricawest", "australiacentral",
	"australiacentral2", "australiasoutheast", "jioindiacentral", "koreasouth",
	"southindia", "westindia", "canadaeast", "francesouth",
	"germanynorth", "norwaywest", "switzerlandwest", "ukwest",
	"uaecentral", "brazilsoutheast"}

// getPssdFromSize is a function that takes the size of a required premium SSD
// disk and returns the SKU string for pricing. This allows users who do not
// know SKU of the disk sizes to simply supply the disk size in the config.
func getPssdFromSize(sz uint) string {
	switch {
	case sz <= 4:
		return "P1"
	case sz <= 8:
		return "P2"
	case sz <= 16:
		return "P3"
	case sz <= 32:
		return "P4"
	case sz <= 64:
		return "P6"
	case sz <= 128:
		return "P10"
	case sz <= 256:
		return "P15"
	case sz <= 512:
		return "P20"
	case sz <= 1024:
		return "P30"
	case sz <= 2048:
		return "P40"
	case sz <= 4096:
		return "P50"
	case sz <= 8192:
		return "P60"
	case sz <= 16384:
		return "P70"
	case sz <= 32768:
		return "P80"
	default:
		return "error"
	}
}
