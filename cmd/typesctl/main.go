package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/RHEnVision/provisioning-backend/cmd/typesctl/providers"
	"github.com/RHEnVision/provisioning-backend/internal/config"
)

func main() {
	config.Initialize()

	validProviders := make([]string, 0)
	for p := range providers.TypeProviders {
		validProviders = append(validProviders, p)
	}
	helpProviders := strings.Join(validProviders, ",")

	var providerFlag = flag.String("provider", "", fmt.Sprintf("provider: [%s] (required)", helpProviders))
	var printAllFlag = flag.Bool("all", false, "print everything (long output)")
	var printTypeFlag = flag.String("type", "", "print specific instance type detail (or 'all')")
	var printRegionFlag = flag.String("region", "", "print instance type names for a region (or 'all')")
	var printZoneFlag = flag.String("zone", "", "print instance type names for a zone (region is needed too)")
	var generateFlag = flag.Bool("generate", false, "generate new type information")
	flag.Parse()

	provider, ok := providers.TypeProviders[strings.ToLower(*providerFlag)]
	if !ok {
		fmt.Println("Unknown or unspecified provider, use -provider")
		flag.Usage()
		return
	}

	if *printAllFlag {
		provider.PrintRegisteredTypes("")
		provider.PrintRegionalAvailability("", "")
	} else if *printTypeFlag == "all" {
		provider.PrintRegisteredTypes("")
	} else if *printTypeFlag != "" {
		provider.PrintRegisteredTypes(*printTypeFlag)
	} else if (*printRegionFlag != "" && *printZoneFlag != "") ||
		(*printRegionFlag != "" && *printZoneFlag == "") ||
		(*printRegionFlag == "all" && *printZoneFlag == "") {
		provider.PrintRegionalAvailability(*printRegionFlag, *printZoneFlag)
	} else if *generateFlag {
		err := provider.GenerateTypes()
		if err != nil {
			panic(err)
		}
	} else {
		flag.Usage()
	}
}
