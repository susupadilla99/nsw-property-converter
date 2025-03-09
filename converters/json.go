package converters

import (
	"encoding/json"
	"fmt"
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
