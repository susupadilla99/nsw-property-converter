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
		fmt.Printf("\b\b\b\b\b\b\b")
		fmt.Printf("%d/%d", i+1, len(items))
	}

	fmt.Print("\nCompleted\n\n")
}

// Read all .DAT files in the provided (temp)/extracted path and return a []Property object
func convertFiles(path string) []Property {
	fmt.Println("Converting property data...")
	time.Sleep(3 * time.Second)

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

	fmt.Print("\nCompleted\n\n")

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

// Convert a yearly zip file to []Property
func ConvertYearToSlice(path string) []Property {

	tempPath := extractYearlyZip(path)

	extractWeeklyZips(tempPath)

	convertedData := convertFiles(filepath.Join(tempPath, "extracted"))

	removeTempDir(tempPath)

	return convertedData
}

// Convert the []Property to [][]string and write to a csv file at "path"
func ConvertSliceToCSV(data []Property, path string) {
	// Convert to 2D Slice
	fmt.Println("Converting file to CSV")
	time.Sleep(3 * time.Second)

	properties := [][]string{}
	for i, property := range data {
		properties = append(properties, converters.ConvertPropertyToSlice(property))

		fmt.Printf("\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b")
		fmt.Printf("%d/%d", i+1, len(data))
	}

	fmt.Printf("\nCompleted\n\n")

	// Write to CSV File
	fmt.Println("Writing to CSV file...")
	time.Sleep(3 * time.Second)

	converters.WriteSliceToCSV(properties, path)

	fmt.Printf("Completed\n\n")
}

// Convert the []Property to a JSON string and returns the string
func ConvertSliceToJSON(data []Property) string {
	fmt.Println("Converting to JSON data...")
	time.Sleep(3 * time.Second)

	resultData := converters.ConvertSliceToJSON(data)

	fmt.Printf("Completed\n\n")

	return resultData
}

func main() {

	// Program Parameters
	INPUT_PATH := "./2024.zip"

	dataSlice := ConvertYearToSlice(INPUT_PATH)

	csvPath := strings.TrimSuffix(INPUT_PATH, filepath.Ext(INPUT_PATH)) + ".csv"

	ConvertSliceToCSV(dataSlice, csvPath)

	fmt.Printf("Exported to %s.\n\n", csvPath)

	ConvertSliceToJSON(dataSlice)

	fmt.Printf("Converted to JSON\n\n")
}
