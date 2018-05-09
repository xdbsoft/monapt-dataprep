package main

import (
	"encoding/csv"
	"os"
)

func readCsv(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(f)
	reader.Comma = ';'

	all, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return all, nil
}
