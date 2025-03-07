package main

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Extract zip file to "./temp" directory
func extractYearlyZip(path string) {
	dir := filepath.Dir(path)

	// open the yearly zip file
	zipYearly, err := zip.OpenReader(path)
	if err != nil {
		log.Fatal(err)
	}
	defer zipYearly.Close()

	// create temp directory if not already exist
	if err := os.MkdirAll(filepath.Join(dir, "temp"), os.ModeDir); err != nil {
		panic(err)
	}

	for _, file := range zipYearly.File {
		// Create destination path
		filePath := filepath.Join(dir, "temp", file.Name)

		// Check if file is a directory
		if file.FileInfo().IsDir() {
			// Create the directory
			if err := os.MkdirAll(filepath.Dir(filePath), os.ModeDir); err != nil {
				panic(err)
			}
		}

		// Create empty destination file
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			panic(err)
		}
		defer dstFile.Close()

		// Open the source file
		srcFile, err := file.Open()
		if err != nil {
			panic(err)
		}
		defer srcFile.Close()

		// Copy content from source file to destination file
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			panic(err)
		}
	}

}

// Extract weekly zip file to "./temp/extracted" directory
func extractWeeklyZip(dir string, fileName string) {
	// var wg sync.WaitGroup

	// Skip extracting if file is a directory. Print error if file is invalid
	if info, err := os.Stat(dir + fileName); err != nil || info.IsDir() {
		if err != nil {
			panic(err)
		}
		return
	}

	zipWeekly, err := zip.OpenReader(dir + fileName)
	if err != nil {
		log.Fatal(err)
	}

	// create temp directory if not already exist
	if err := os.MkdirAll(dir+"/extracted", os.ModeDir); err != nil {
		panic(err)
	}

	for _, file := range zipWeekly.File {
		extractFile(file, filepath.Join(dir, "/extracted/"))
		// // Create destination path
		// filePath := filepath.Join(dir, "/extracted/", file.Name)

		// // Check if file is a directory
		// if file.FileInfo().IsDir() {
		// 	// Create the directory
		// 	if err := os.MkdirAll(filepath.Dir(filePath), os.ModeDir); err != nil {
		// 		panic(err)
		// 	}
		// }

		// // Create empty destination file
		// // fmt.Println("Opening file: " + fileName)
		// dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		// if err != nil {
		// 	panic(err)
		// }

		// // Open the source file
		// wg.Add(1)
		// srcFile, err := file.Open()
		// if err != nil {
		// 	panic(err)
		// }

		// // Copy content from source file to destination file
		// if _, err := io.Copy(dstFile, srcFile); err != nil {
		// 	panic(err)
		// }

		// // fmt.Println("Closing: " + fileName)

		// srcFile.Close()
		// dstFile.Close()

		// // fmt.Println("Closed: " + fileName)
		// wg.Done()

		// // Close both files

		// // fmt.Println("Removing: " + fileName)
		// // remove(dir + fileName)
	}

	// wg.Wait()
}

// Extract file to destination folder
func extractFile(src *zip.File, destFolder string) {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Done()

	destPath := filepath.Join(destFolder, src.Name)

	// Ignore files that are directories
	if src.FileInfo().IsDir() {
		// // Create the directory
		// if err := os.MkdirAll(filepath.Dir(destPath), os.ModeDir); err != nil {
		// 	panic(err)
		// }
		return
	}

	// Ignore files that are not .DAT
	if filepath.Ext(src.Name) != ".DAT" {
		return
	}

	// Create empty destination file
	dstFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, src.Mode())
	if err != nil {
		panic(err)
	}
	defer dstFile.Close()

	// Open the source file
	srcFile, err := src.Open()
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	// Copy content from source file to destination file
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		panic(err)
	}
}

// Read data in provided ".DAT" file and return a 2D slice
func readDataFile(dir string, file string) [][]string {
	f, err := os.ReadFile(dir + file)
	if err != nil {
		log.Fatal(err)
	}

	fileText := string(f[:])
	textArr := strings.Split(fileText, "\n")
	fileOutput := [][]string{}
	lineSlice := []string{}

	for _, line := range textArr {

		if len(line) > 0 {
			switch recType := line[0]; recType {
			case 'B':
				lineSlice = []string{} // Reset lineSlice

				// Adds B record data to lineSlice
				recItems := strings.Split(line, ";")
				lineSlice = append(lineSlice, recItems...)
				lineSlice[len(lineSlice)-1] = "" // Make sure the last string in the slice is empty. This string is used to hold all values of relevant C records
				fileOutput = append(fileOutput, lineSlice)
			case 'C':
				recItems := strings.Split(line, ";")
				if recItems[3] != lineSlice[3] {
					fmt.Printf("Error found in function readDataFile(string, string)")
					fmt.Printf("Property ID (%s) of record type C does not match property ID (%s) of record type B. File: %s \n", recItems[3], lineSlice[3], file)
					continue
				}
				lineSlice[len(lineSlice)-1] += recItems[5]
			}
		}
	}

	// fmt.Println(fileOutput)

	return fileOutput
}

func convertZipToSlice(path string) {
	dir := filepath.Dir(path)
	fileName := filepath.Base(path)

	fmt.Printf("Dir: %s, Path: %s", dir, fileName)
}

func main() {

	// convertZipToSlice("./2024.zip")

	// Extract yearly zip file to "temp" directory
	fmt.Println("Extracting zip...")
	extractYearlyZip("./2024.zip")
	fmt.Print("Completed\n\n")

	time.Sleep(3 * time.Second)

	// Read "temp" directory to get all weekly zip files
	fmt.Println("Extracting inner zip...")
	items, err := os.ReadDir("./temp")
	if err != nil {
		panic(err)
	}

	// Extract all weekly zip files to "temp/extracted" directory
	for i, item := range items {
		extractWeeklyZip("./temp/", item.Name())
		fmt.Printf("\b\b\b\b\b")
		fmt.Printf("%d/%d", i+1, len(items))
	}
	fmt.Print("\b\b\b\b\bCompleted\n\n")

	// Open new csv file to write result to
	csvFile, err := os.Create("./2024.csv")
	if err != nil {
		panic(err)
	}

	// Read all files in "extracted" directory
	fmt.Println("Reading property data...")
	time.Sleep(3 * time.Second)
	resultString := [][]string{
		{
			"Record Type", "District Code", "Property ID", "Sale Counter", "Download Date/Time", "Property Name",
			"Property Unit Number", "Property House Number", "Property Street Name", "Property Locality", "Property Post Code",
			"Area", "Area Type", "Contract Date", "Settlement Date", "Purchase Price", "Zoning", "Nature of Property",
			"Primary Purpose", "Strata Lot Number", "Component Code", "Sale Code", "% Interest of Sale", "Dealing Number", "Property Legal Description",
		},
	}

	entries, err := os.ReadDir("./temp/extracted")
	if err != nil {
		panic(err)
	}

	for i, entry := range entries {
		// if i >= 1 {
		// 	continue
		// }

		resultString = append(resultString, readDataFile("./temp/extracted/", entry.Name())...)

		// fmt.Printf("%d - %s\n", i, entry.Name())
		// fmt.Print("-----------------------\n")
		// fmt.Println(resultString)

		fmt.Print("\b\b\b\b\b\b\b\b\b")
		fmt.Printf("%d/%d", i+1, len(entries))
	}

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()
	writer.WriteAll(resultString)

	fmt.Print("\b\b\b\b\b\b\b\b\bCompleted\n\n")

	// Clean the "temp" directory
	time.Sleep(3 * time.Second)
	fmt.Println("Cleaning \"temp\" directory...")
	removeErr := os.RemoveAll("./temp")
	for removeErr != nil {
		removeErr = os.RemoveAll("./temp")
	}
	fmt.Print("Completed\n\n")
}
