package converters

import (
	"encoding/csv"
	"os"
)

// Convert provided 2D data slice to csv and write to csv file at "path"
func ConvertSliceToCSV(data [][]string, path string) string {

	// Open new csv file to write result to
	csvFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()
	writer.WriteAll(data)

	return path
}
