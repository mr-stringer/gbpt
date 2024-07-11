# GBPT

Go Based Pricing Tool

https://prices.azure.com/api/retail/prices?api-version=2023-01-01-preview&currencyCode='GBP'&$filter=armRegionName eq 'uksouth' and priceType eq 'Consumption' and serviceFamily eq 'Compute' and serviceName eq 'Virtual Machines' and skuName eq 'Standard_D2ds_v5'

https://prices.azure.com/api/retail/prices?api-version=2023-01-01-preview&currencyCode='GBP'&$filter=armRegionName eq 'uksouth' and serviceFamily eq 'Compute' and serviceName eq 'Virtual Machines' and skuName eq 'Standard_D2ds_v5'

s2 = fmt.Sprintf("%scurrencyCode='%s'&$filter=armRegionName eq '%s' and SkuName eq '%s' and priceType eq 'Consumption' and serviceFamily eq 'Compute' and serviceName eq 'Virtual Machines'", BaseUrl, r[0], r[1], r[2])

% az vm list-sizes --location uksouth --query "[?numberOfCores ==\`2\` && memoryInMb <=\`9000\`  && memoryInMb >=\`8000\`].{Name:name,  Cores:numberOfCores, MemoryMiB:memoryInMb}" --output table
Name                   Cores    MemoryMiB
---------------------  -------  -----------
Standard_B2ms          2        8192
Standard_B2s           2        4096
...