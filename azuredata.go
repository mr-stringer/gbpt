package main

/* Preview version of the API currently needed for pssd_v2 prices*/
const ApiUrl = "https://prices.azure.com/api/retail/prices"
const ApiPreview = "2023-01-01-preview"

var supcur = []string{"USD", "AUD", "BRL",
	"CAD", "CHF", "CNY",
	"DKK", "EUR", "GBP",
	"INR", "JPY", "KRW",
	"NOK", "NZD", "RUB",
	"SEK", "TWD"}

/* list of regions generated with az */
/* az account list-locations --query "[].{Name:name}" -o table */
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

func GetPssdFromSize(sz uint) string {
	switch {
	case sz <= 4:
		return "p1"
	case sz <= 8:
		return "p2"
	case sz <= 16:
		return "p3"
	case sz <= 32:
		return "p4"
	case sz <= 64:
		return "p6"
	case sz <= 128:
		return "p10"
	case sz <= 256:
		return "p15"
	case sz <= 512:
		return "p20"
	case sz <= 1024:
		return "p30"
	case sz <= 2048:
		return "p40"
	case sz <= 4096:
		return "p50"
	case sz <= 8192:
		return "p60"
	case sz <= 16384:
		return "p70"
	case sz <= 32768:
		return "p80"
	default:
		return "error"
	}
}
