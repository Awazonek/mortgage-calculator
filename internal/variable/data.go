package variable

import (
	"encoding/json"
	"fmt"
	"os"
)

// Structs to match the JSON structure
type Data struct {
	X      XData                 `json:"x"`
	Series map[string]SeriesData `json:"series"`
}

type XData struct {
	Data []string `json:"data"`
}

type SeriesData struct {
	Data []float64 `json:"data"`
}

func ReadData(filename string) (*Data, error) {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var data Data
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func MapDatesToRates(data *Data) (map[string]float64, error) {
	if len(data.X.Data) != len(data.Series["best-rates.5y-variable"].Data) {
		return nil, fmt.Errorf("mismatched number of dates and rates")
	}

	dateToRate := make(map[string]float64)
	for i := range data.X.Data {
		dateStr := data.X.Data[i]
		rate := data.Series["best-rates.5y-variable"].Data[i]
		dateToRate[dateStr] = rate
	}

	return dateToRate, nil
}
