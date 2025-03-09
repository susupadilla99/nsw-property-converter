package converters

import (
	"encoding/json"

	"github.com/susupadilla99/nsw-property-converter/extractors"
)

type Property = extractors.Property

func ConvertSliceToJSON(data []Property) string {
	//Convert Property slice to indented Json
	bytes, err := json.MarshalIndent(data[0:10], "", "  ")
	if err != nil {
		panic(err)
	}

	return string(bytes)
}
