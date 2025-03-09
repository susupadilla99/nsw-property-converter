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

type Property struct {
	Record_Type                string `json:"record_type"`
	District_Code              string `json:"district_code"`
	Property_ID                string `json:"property_id"`
	Sale_Counter               string `json:"sale_counter"`
	Download_Date_Time         string `json:"download_date_time"`
	Property_Name              string `json:"property_name"`
	Property_Unit_Number       string `json:"property_unit_number"`
	Property_House_Number      string `json:"property_house_number"`
	Property_Street_Name       string `json:"property_street_name"`
	Property_Locality          string `json:"property_locality"`
	Property_Post_Code         string `json:"property_post_code"`
	Area                       string `json:"area"`
	Area_Type                  string `json:"area_type"`
	Contract_Date              string `json:"contract_date"`
	Settlement_Date            string `json:"settlement_date"`
	Purchase_Price             string `json:"purchase_price"`
	Zoning                     string `json:"zoning"`
	Nature_Of_Property         string `json:"nature_of_property"`
	Primary_Purpose            string `json:"primary_purpose"`
	Strata_Lot_Number          string `json:"strata_lot_number"`
	Component_Code             string `json:"component_code"`
	Sale_Code                  string `json:"sale_code"`
	Percent_Interest_Of_Sale   string `json:"percent_interest_of_sale"`
	Dealing_Number             string `json:"dealing_number"`
	Property_Legal_Description string `json:"property_legal_description"`
}

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

// Read data in provided ".DAT" file and return a slice of Property objects
func ReadDataFile(path string) []Property {
	// Check if file path is valid
	f, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	fileStr := string(f[:])                 // convert file data to string
	records := strings.Split(fileStr, "\n") // split file string into multiple records
	data := []Property{}                    // whole data extracted from file
	currentRecord := Property{}             // represents current type-B record

	for _, record := range records {

		if len(record) > 0 {
			switch recordType := record[0]; recordType {
			case 'A': // Record A => Do nothing
				continue
			case 'B': // Record B => Add all data to current record
				currentRecord = Property{} // Reset current record

				// Add record data to current record
				recItems := strings.Split(record, ";")
				currentRecord = Property{
					Record_Type:                recItems[0],
					District_Code:              recItems[1],
					Property_ID:                recItems[2],
					Sale_Counter:               recItems[3],
					Download_Date_Time:         recItems[4],
					Property_Name:              recItems[5],
					Property_Unit_Number:       recItems[6],
					Property_House_Number:      recItems[7],
					Property_Street_Name:       recItems[8],
					Property_Locality:          recItems[9],
					Property_Post_Code:         recItems[10],
					Area:                       recItems[11],
					Area_Type:                  recItems[12],
					Contract_Date:              recItems[13],
					Settlement_Date:            recItems[14],
					Purchase_Price:             recItems[15],
					Zoning:                     recItems[16],
					Nature_Of_Property:         recItems[17],
					Primary_Purpose:            recItems[18],
					Strata_Lot_Number:          recItems[19],
					Component_Code:             recItems[20],
					Sale_Code:                  recItems[21],
					Percent_Interest_Of_Sale:   recItems[22],
					Dealing_Number:             recItems[23],
					Property_Legal_Description: recItems[24],
				}
				currentRecord.Property_Legal_Description = "" // Make sure "Property_Legal_Description field is empty. This field is used to hold all values of relevant C records
				data = append(data, currentRecord)
			case 'C': // Record C => Add data to Property_Legal_Description field
				recItems := strings.Split(record, ";")
				if recItems[3] != currentRecord.Property_ID {
					fmt.Printf("Error found in function readDataFile(string, string)")
					fmt.Printf("Property ID (%s) of record type C does not match property ID (%s) of record type B. File: %s \n", recItems[3], currentRecord.Property_ID, filepath.Base(path))
					continue
				}
				currentRecord.Property_Legal_Description += recItems[5]
			case 'D': // Record D => Do nothing
			case 'Z': // Record Z => Do nothing
			}
		}
	}

	// fmt.Println(fileOutput)

	return data
}
