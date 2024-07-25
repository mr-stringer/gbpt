package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

var (
	Client HTTPClient
)

func main() {
	/* Create logger */
	/*Do something here to get the log level from command line! */
	lvlSet, printConfig, err := handleFlags()
	if err != nil {
		fmt.Print(err)
		os.Exit(osExitIncorrectFlagConfig)
	}

	lvl := new(slog.LevelVar)
	lvl.Set(lvlSet)

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	}))

	slog.SetDefault(logger)

	c1, err := GetConfig()
	if err != nil {
		os.Exit(osExitLoadConfig)
	}
	eStrings, err := c1.Validate()
	if err != nil {
		slog.Error("One or more configuration errors were discovered")
		for _, v := range eStrings {
			slog.Error(v)
		}
		slog.Error("Failed to validate error ", "error", err.Error())
		os.Exit(osExitValidateConfig)

	} else {
		slog.Info("Config is valid")
	}

	if printConfig {
		c1.Print()
		os.Exit(0)
	}

	/*Init the client */
	Client = &http.Client{}
	err = c1.PriceConfig()
	if err != nil {
		slog.Info("Failed to price config ", "error", err.Error())
		os.Exit(osExitPriceConfig)
	}
	slog.Info("Process Complete")
}
