package converters

import (
	"encoding/csv"
	"os"
	"path/filepath"
)

func ConvertSliceToCSV(data [][]string, path string) string {

	// Open new csv file to write result to
	csvFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()
	writer.WriteAll(data)

	return filepath.Join(filepath.Dir(path), csvFile.Name())
}
