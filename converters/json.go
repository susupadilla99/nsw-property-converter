package converters

import (
	"encoding/json"
	"fmt"

	"github.com/susupadilla99/nsw-property-converter/extractors"
)

type Property extractors.Property

func ConvertSliceToJSON(data [][]string) string {
	jsonData := []Property{}

	fmt.Println(" I am here!")
	fmt.Println(data[0])

	for i, item := range data {
		// Skip first header item
		if i == 0 {
			continue
		}

		jsonData = append(jsonData, Property{
			Record_Type:                item[0],
			District_Code:              item[1],
			Property_ID:                item[2],
			Sale_Counter:               item[3],
			Download_Date_Time:         item[4],
			Property_Name:              item[5],
			Property_Unit_Number:       item[6],
			Property_House_Number:      item[7],
			Property_Street_Name:       item[8],
			Property_Locality:          item[9],
			Property_Post_Code:         item[10],
			Area:                       item[11],
			Area_Type:                  item[12],
			Contract_Date:              item[13],
			Settlement_Date:            item[14],
			Purchase_Price:             item[15],
			Zoning:                     item[16],
			Nature_Of_Property:         item[17],
			Primary_Purpose:            item[18],
			Strata_Lot_Number:          item[19],
			Component_Code:             item[20],
			Sale_Code:                  item[21],
			Percent_Interest_Of_Sale:   item[22],
			Dealing_Number:             item[23],
			Property_Legal_Description: item[24],
		})
	}

	b, err := json.MarshalIndent(jsonData[0:10], "", "  ")
	if err != nil {
		panic(err)
	}

	return string(b)
}
