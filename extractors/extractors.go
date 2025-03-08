package extractors

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// Extract zip file to a temp directory and return the path to that directory
func ExtractYearlyZip(path string) string {
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
func ExtractWeeklyZip(path string) {
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
func ReadDataFile(path string) [][]string {

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
