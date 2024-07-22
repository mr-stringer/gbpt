package main

import (
	"bufio"
	"log"
	"os"
)

// WriteToFile is a function to takes a slice of strings and outputs them in order
// to a file.
func WriteToFile(data *[]string, path string) error {
	f, err := os.Create(path)
	if err != nil {
		log.Printf("Could not create file %s\n", path)
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, l := range *data {
		_, err := w.WriteString(l + "\n")
		if err != nil {
			log.Printf("Failed to write file")
			return err
		}
	}
	w.Flush()
	return nil
}
