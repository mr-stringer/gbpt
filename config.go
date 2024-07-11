package main

import (
	"log"

	"github.com/spf13/viper"
)

func GetConfig() Config {

	viper.AddConfigPath("./")
	viper.SetConfigFile("gbpt.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Print("Failed to read in config")
		log.Fatal(err)
	}

	c1 := Config{}
	err = viper.Unmarshal(&c1)
	if err != nil {
		log.Print("Failed ro unmarshal config")
		log.Fatal(err)
	}
	return c1
}
