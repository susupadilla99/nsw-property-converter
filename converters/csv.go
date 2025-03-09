package converters

import (
	"encoding/csv"
	"os"
)

// Convert a single Property object to a corresponding []string
func ConvertPropertyToSlice(data Property) []string {
	res := []string{
		data.Record_Type,
		data.District_Code,
		data.Property_ID,
		data.Sale_Counter,
		data.Download_Date_Time,
		data.Property_Name,
		data.Property_Unit_Number,
		data.Property_House_Number,
		data.Property_Street_Name,
		data.Property_Locality,
		data.Property_Post_Code,
		data.Area,
		data.Area_Type,
		data.Contract_Date,
		data.Settlement_Date,
		data.Purchase_Price,
		data.Zoning,
		data.Nature_Of_Property,
		data.Primary_Purpose,
		data.Strata_Lot_Number,
		data.Component_Code,
		data.Sale_Code,
		data.Percent_Interest_Of_Sale,
		data.Dealing_Number,
		data.Property_Legal_Description,
	}

	return res
}

// Convert a []Property to a corresponding []string
func ConvertPropertiesToSlices(data []Property) [][]string {
	res := [][]string{}

	for _, record := range data {
		res = append(res, ConvertPropertyToSlice(record))
	}

	return res
}

// Add a header "row" at the top of the [][]string
func AddHeader(data [][]string) [][]string {
	res := [][]string{
		{
			"Record Type", "District Code", "Property ID", "Sale Counter", "Download Date/Time", "Property Name",
			"Property Unit Number", "Property House Number", "Property Street Name", "Property Locality", "Property Post Code",
			"Area", "Area Type", "Contract Date", "Settlement Date", "Purchase Price", "Zoning", "Nature of Property",
			"Primary Purpose", "Strata Lot Number", "Component Code", "Sale Code", "% Interest of Sale", "Dealing Number", "Property Legal Description",
		},
	}

	res = append(res, data...)

	return res
}

// Write provided 2d slice to csv to csv file at "path"
func WriteSliceToCSV(data [][]string, path string) {
	// Open new csv file to write result to
	csvFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	// Write data to new csv file
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()
	writer.WriteAll(data)
}
