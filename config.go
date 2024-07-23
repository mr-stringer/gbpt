package main

import (
	"log/slog"

	"github.com/spf13/viper"
)

const (
	osExitIncorrectFlagConfig = iota + 1
	osExitLoadConfig
	osExitValidateConfig
	osExitPriceConfig
)

func GetConfig() (Config, error) {
	slog.Info("Attempting to read config")
	viper.AddConfigPath("./")
	viper.SetConfigFile("gbpt.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		slog.Error("Failed to read in config")
		return Config{}, err
	}

	c1 := Config{}
	err = viper.Unmarshal(&c1)
	if err != nil {
		slog.Error("Failed ro unmarshal config")
		return Config{}, err
	}
	return c1, nil
}
