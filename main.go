package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/susupadilla99/nsw-property-converter/converters"
	"github.com/susupadilla99/nsw-property-converter/extractors"
)

type Property = extractors.Property

// Extract all yearly zip file to a temp directory and return the path to that directory
func extractYearlyZip(path string) string {
	fmt.Println("Extracting zip...")
	time.Sleep(3 * time.Second)

	resultTempPath := extractors.ExtractYearlyZip(path)

	fmt.Print("Completed\n\n")

	return resultTempPath
}

// Read provided (temp) directory and extract all weekly zip files to (temp)/extracted folder
func extractWeeklyZips(path string) {
	fmt.Println("Extracting inner zip...")
	time.Sleep(3 * time.Second)

	items, err := os.ReadDir(path) // Read temp directory
	if err != nil {
		panic(err)
	}

	// Extract all weekly zip files to "temp/extracted" directory
	for i, item := range items {
		extractors.ExtractWeeklyZip(filepath.Join(path, item.Name()))
		fmt.Printf("\b\b\b\b\b")
		fmt.Printf("%d/%d", i+1, len(items))
	}

	fmt.Print("\b\b\b\b\bCompleted\n\n")
}

// Read all .DAT files in the provided (temp)/extracted path and return a 2D slice containing all those data with header
func convertFilesToSlice(path string) []Property {
	fmt.Println("Converting property data...")
	time.Sleep(3 * time.Second)

	// resultData := [][]string{
	// 	{
	// 		"Record Type", "District Code", "Property ID", "Sale Counter", "Download Date/Time", "Property Name",
	// 		"Property Unit Number", "Property House Number", "Property Street Name", "Property Locality", "Property Post Code",
	// 		"Area", "Area Type", "Contract Date", "Settlement Date", "Purchase Price", "Zoning", "Nature of Property",
	// 		"Primary Purpose", "Strata Lot Number", "Component Code", "Sale Code", "% Interest of Sale", "Dealing Number", "Property Legal Description",
	// 	},
	// }

	resultData := []Property{}

	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for i, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		resultData = append(resultData, extractors.ReadDataFile(entryPath)...)

		fmt.Print("\b\b\b\b\b\b\b\b\b")
		fmt.Printf("%d/%d", i+1, len(entries))
	}

	fmt.Print("\b\b\b\b\b\b\b\b\bCompleted\n\n")

	return resultData
}

// Remove the provided (temp) directory
func removeTempDir(path string) {
	fmt.Printf("Cleaning %s directory...\n", path)
	time.Sleep(3 * time.Second)

	removeErr := os.RemoveAll(path)
	for removeErr != nil {
		removeErr = os.RemoveAll(path)
	}

	fmt.Print("Completed\n\n")
}

// Convert a yearly zip file to a 2D string slice
func ConvertYearToSlice(path string) []Property {

	tempPath := extractYearlyZip(path)

	extractWeeklyZips(tempPath)

	convertedData := convertFilesToSlice(filepath.Join(tempPath, "extracted"))

	removeTempDir(tempPath)

	return convertedData
}

func ConvertSliceToCSV(data []Property, path string) string {
	fmt.Println("Writing to CSV file...")
	time.Sleep(3 * time.Second)

	resultFilePath := converters.ConvertSliceToCSV(data, path)

	fmt.Printf("Completed\n\n")

	return resultFilePath
}

func ConvertSliceToJSON(data []Property) {
	fmt.Println("Converting to JSON data...")
	time.Sleep(3 * time.Second)

	fmt.Println("I am here first")
	fmt.Println(data[0])

	resultData := converters.ConvertSliceToJSON(data)

	fmt.Println(resultData)

	fmt.Printf("Completed\n\n")
}

func main() {

	// Program Parameters
	INPUT_PATH := "./2024.zip"

	dataSlice := ConvertYearToSlice(INPUT_PATH)

	csvPath := strings.TrimSuffix(INPUT_PATH, filepath.Ext(INPUT_PATH)) + ".csv"

	csvFileFullPath := ConvertSliceToCSV(dataSlice, csvPath)

	fmt.Printf("Exported to %s.\n\n", csvFileFullPath)

	ConvertSliceToJSON(dataSlice)

}
