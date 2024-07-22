package main

import (
	"log"
)

func main() {
	log.Print("gbpt is getting the config")
	c1 := GetConfig()
	c1.Print()
	eStrings, err := c1.Validate()
	if err != nil {
		log.Println("One or more configuration errors were discovered")
		for _, v := range eStrings {
			log.Println(v)
		}
		log.Fatal(err)
	} else {
		log.Println("Config is valid")
	}

	err = c1.PriceConfig()
	if err != nil {
		log.Fatal(err)
	}
}
