package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func extractYearlyZip(dir string, file string) {
	zipYearly, err := zip.OpenReader(dir + file)
	if err != nil {
		log.Fatal(err)
	}
	defer zipYearly.Close()

	// create temp directory if not already exist
	if err := os.MkdirAll(dir+"/temp", os.ModeDir); err != nil {
		panic(err)
	}

	for _, file := range zipYearly.File {
		// Create destination path
		filePath := filepath.Join(dir, "/temp/", file.Name)
		// fmt.Println("Extracting file ", filePath)

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

		// Open the source file
		srcFile, err := file.Open()
		if err != nil {
			panic(err)
		}

		// Copy content from source file to destination file
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			panic(err)
		}

		// Close both files
		dstFile.Close()
		srcFile.Close()
	}

}

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

func readDataFile(dir string, file string) string {
	f, err := os.ReadFile(dir + file)
	if err != nil {
		log.Fatal(err)
	}

	fileText := string(f[:])
	textArr := strings.Split(fileText, "\n")
	lineOutput := ""

	for _, line := range textArr {
		if len(line) > 0 {
			switch recType := line[0]; recType {
			case 'B':
				// fmt.Println(lineOutput)
				lineOutput += "\n"
				recItems := strings.Split(line, ";")
				for _, item := range recItems {
					lineOutput += item + " - "
				}
				lineOutput += "\b\b\b"
			case 'C':
				recItems := strings.Split(line, ";")
				lineOutput += recItems[5]
			}
		}
	}

	return lineOutput
}

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

func main() {
	// Extract yearly zip file to "temp" directory
	fmt.Println("Extracting zip...")
	extractYearlyZip(".", "/2024.zip")
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

	// Read all files in "extracted" directory
	fmt.Println("Reading property data:")
	time.Sleep(3 * time.Second)
	entries, err := os.ReadDir("./temp/extracted")
	if err != nil {
		panic(err)
	}

	for i, entry := range entries {
		if i <= 100 {
			continue
		}
		if i >= 110 {
			break
		}
		resultString := readDataFile("./temp/extracted/", entry.Name())
		fmt.Printf("%d - %s\n", i, entry.Name())
		fmt.Print("-----------------------")
		fmt.Print(resultString + "\n\n")
		// fmt.Print("\b\b\b\b\b\b\b\b\b")
		// fmt.Printf("%d/%d", i, len(entries))
	}
	fmt.Print("\b\b\b\b\b\b\b\b\bCompleted\n\n")

	// Clean the "temp" directory
	time.Sleep(3 * time.Second)
	fmt.Println("Cleaning \"temp\" directory")
	os.RemoveAll("./temp")
	fmt.Print("Completed\n\n")
	// readDataFile("./temp/extracted/", "223_SALES_DATA_NNME_01012024.DAT")
}
