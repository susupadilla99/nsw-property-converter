package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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
		fmt.Println("Extracting file ", filePath)

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

func main() {
	extractYearlyZip(".", "/2024.zip")
	time.Sleep(10 * time.Second)
	os.RemoveAll("./temp/")
}
