package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/susupadilla99/nsw-property-converter/extractors"
)

// Convert a yearly zip file to a 2d string slice
func convertYearToSlice(path string) [][]string {
	// Extract yearly zip file to "temp" directory
	fmt.Println("Extracting zip...")
	tempPath := extractors.ExtractYearlyZip(path)
	fmt.Print("Completed\n\n")

	time.Sleep(3 * time.Second)

	// Read "temp" directory to get all weekly zip files
	fmt.Println("Extracting inner zip...")
	items, err := os.ReadDir(tempPath)
	if err != nil {
		panic(err)
	}

	// Extract all weekly zip files to "temp/extracted" directory
	for i, item := range items {
		extractors.ExtractWeeklyZip(filepath.Join(tempPath, item.Name()))
		fmt.Printf("\b\b\b\b\b")
		fmt.Printf("%d/%d", i+1, len(items))
	}
	fmt.Print("\b\b\b\b\bCompleted\n\n")

	// Read all files in "extracted" directory
	fmt.Println("Converting property data...")
	time.Sleep(3 * time.Second)
	convertedData := [][]string{
		{
			"Record Type", "District Code", "Property ID", "Sale Counter", "Download Date/Time", "Property Name",
			"Property Unit Number", "Property House Number", "Property Street Name", "Property Locality", "Property Post Code",
			"Area", "Area Type", "Contract Date", "Settlement Date", "Purchase Price", "Zoning", "Nature of Property",
			"Primary Purpose", "Strata Lot Number", "Component Code", "Sale Code", "% Interest of Sale", "Dealing Number", "Property Legal Description",
		},
	}

	entries, err := os.ReadDir(filepath.Join(tempPath, "extracted"))
	if err != nil {
		panic(err)
	}

	for i, entry := range entries {
		entryPath := filepath.Join(tempPath, "extracted", entry.Name())
		convertedData = append(convertedData, extractors.ReadDataFile(entryPath)...)

		fmt.Print("\b\b\b\b\b\b\b\b\b")
		fmt.Printf("%d/%d", i+1, len(entries))
	}
	fmt.Print("\b\b\b\b\b\b\b\b\bCompleted\n\n")

	// Clean the "temp" directory
	fmt.Printf("Cleaning %s directory...\n", tempPath)
	removeErr := os.RemoveAll(tempPath)
	for removeErr != nil {
		removeErr = os.RemoveAll(tempPath)
	}
	fmt.Print("Completed\n\n")
	time.Sleep(3 * time.Second)

	return convertedData
}

func writeToCSV(data [][]string, path string) {
	fmt.Println("Writing to CSV file...")

	// Open new csv file to write result to
	csvFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()
	writer.WriteAll(data)

	time.Sleep(3 * time.Second)
	fmt.Println("Completed")
}

func main() {

	// Program Parameters
	INPUT_PATH := "./2024.zip"

	dataSlice := convertYearToSlice(INPUT_PATH)

	csvPath := strings.TrimSuffix(INPUT_PATH, filepath.Ext(INPUT_PATH)) + ".csv"

	writeToCSV(dataSlice, csvPath)
}
