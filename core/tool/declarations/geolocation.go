package declarations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/smtdfc/nagare/core/tool"
)

type GeoArgs struct {
	Query string `json:"query" jsonschema:"description=Place name, for example: Hanoi"`
}

type geoResult struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"latitude"`
	Lon     float64 `json:"longitude"`
	Country string  `json:"country"`
}

type geoAPIResponse struct {
	Results []geoResult `json:"results"`
}

func callGeoAPI(query string) (*geoResult, error) {
	url := fmt.Sprintf(
		"https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1",
		query,
	)

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data geoAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if len(data.Results) == 0 {
		return nil, fmt.Errorf("location not found")
	}

	return &data.Results[0], nil
}

var geolocation = tool.DeclareTool(
	"geolocation",
	"Convert location names to coordinates.",
	func(ctx context.Context, args GeoArgs) (any, error) {

		r, err := callGeoAPI(args.Query)
		if err != nil {
			return nil, err
		}

		return map[string]any{
			"name":    r.Name,
			"lat":     r.Lat,
			"lon":     r.Lon,
			"country": r.Country,
		}, nil
	},
)
