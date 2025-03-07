package main

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Extract zip file to a temp directory and return the path to that directory
func extractYearlyZip(path string) string {
	dir := filepath.Dir(path)

	// open the yearly zip file
	zipYearly, err := zip.OpenReader(path)
	if err != nil {
		log.Fatal(err)
	}
	defer zipYearly.Close()

	// Check if temp directory already exist, if yes, change name until unique temp folder name found
	tempFolderName := "temp"
	index := 0
	_, tempErr := os.Stat(filepath.Join(dir, tempFolderName))
	if err != nil && !os.IsExist(tempErr) {
		panic(err)
	}

	for tempErr == nil || os.IsExist(tempErr) {
		if len(tempFolderName) > 4 { // If not the first time, remove last number
			tempFolderName = tempFolderName[:len(tempFolderName)-len(strconv.Itoa(index))]
		}
		index++
		tempFolderName += strconv.Itoa(index)
		_, tempErr = os.Stat(filepath.Join(dir, tempFolderName))
	}

	// create temp directory if not already exist
	if err := os.MkdirAll(filepath.Join(dir, tempFolderName), os.ModeDir); err != nil {
		panic(err)
	}

	for _, file := range zipYearly.File {
		// Create destination path
		filePath := filepath.Join(dir, tempFolderName, file.Name)

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

	return filepath.Join(dir, tempFolderName)
}

// Extract weekly zip file to "./temp/extracted" directory
func extractWeeklyZip(path string) {
	dir := filepath.Dir(path)

	// Skip extracting if file is a directory. Print error if file is invalid
	if info, err := os.Stat(path); err != nil || info.IsDir() {
		if err != nil {
			panic(err)
		}
		return
	}

	zipWeekly, err := zip.OpenReader(path)
	if err != nil {
		log.Fatal(err)
	}

	// create temp directory if not already exist
	if err := os.MkdirAll(filepath.Join(dir, "extracted"), os.ModeDir); err != nil {
		panic(err)
	}

	for _, file := range zipWeekly.File {
		extractFile(file, filepath.Join(dir, "extracted"))
	}
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
func readDataFile(path string) [][]string {

	f, err := os.ReadFile(path)
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
					fmt.Printf("Property ID (%s) of record type C does not match property ID (%s) of record type B. File: %s \n", recItems[3], lineSlice[3], filepath.Base(path))
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

	// Program Parameters
	INPUT_PATH := "./2024.zip"

	// convertZipToSlice("./2024.zip")

	// Extract yearly zip file to "temp" directory
	fmt.Println("Extracting zip...")
	tempPath := extractYearlyZip(INPUT_PATH)
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
		extractWeeklyZip(filepath.Join(tempPath, item.Name()))
		fmt.Printf("\b\b\b\b\b")
		fmt.Printf("%d/%d", i+1, len(items))
	}
	fmt.Print("\b\b\b\b\bCompleted\n\n")

	// Open new csv file to write result to
	csvFileName := strings.TrimSuffix(filepath.Base(INPUT_PATH), filepath.Ext(INPUT_PATH)) + ".csv"
	csvPath := filepath.Join(filepath.Dir(INPUT_PATH), csvFileName)
	csvFile, err := os.Create(csvPath)
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

	entries, err := os.ReadDir(filepath.Join(tempPath, "extracted"))
	if err != nil {
		panic(err)
	}

	for i, entry := range entries {
		entryPath := filepath.Join(tempPath, "extracted", entry.Name())
		resultString = append(resultString, readDataFile(entryPath)...)

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
	fmt.Printf("Cleaning %s directory...\n", tempPath)
	removeErr := os.RemoveAll(tempPath)
	for removeErr != nil {
		removeErr = os.RemoveAll(tempPath)
	}
	fmt.Print("Completed\n\n")
	time.Sleep(3 * time.Second)
}
