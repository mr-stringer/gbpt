package main

import (
	"bufio"
	"os"

	"golang.org/x/exp/slog"
)

// WriteToFile is a function to takes a slice of strings and outputs them in order
// to a file.
func WriteToFile(data *[]string, path string) error {
	f, err := os.Create(path)
	if err != nil {
		slog.Error("Could not create file", "path", path)
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, l := range *data {
		_, err := w.WriteString(l + "\n")
		if err != nil {
			slog.Error("Failed to write file", "path", path)
			return err
		}
	}
	w.Flush()
	slog.Info("Price file written", "path", path)
	return nil
}
