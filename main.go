package main

import (
	"log/slog"
	"os"
)

func main() {
	/* Create logger */
	/*Do something here to get the log level from command line! */
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)

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
		slog.Info("One or more configuration errors were discovered")
		for _, v := range eStrings {
			slog.Info(v)
		}
		slog.Error("Failed to validate error ", "error", err.Error())
		os.Exit(osExitValidateConfig)

	} else {
		slog.Info("Config is valid")
		c1.Print()
	}

	err = c1.PriceConfig()
	if err != nil {
		slog.Info("Failed to price config ", "error", err.Error())
		os.Exit(osExitPriceConfig)
	}
	slog.Info("Done")
}
