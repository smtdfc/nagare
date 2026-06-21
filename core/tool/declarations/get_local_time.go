package declarations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/smtdfc/nagare/core/tool"
)

type LocalTimeArgs struct {
	City string `json:"city" jsonschema:"description=City name"`
}

type geoResp struct {
	Results []struct {
		Name string  `json:"name"`
		Lat  float64 `json:"latitude"`
		Lon  float64 `json:"longitude"`
	} `json:"results"`
}

func fetchGeo(city string) (*geoResp, error) {
	url := fmt.Sprintf(
		"https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1",
		city,
	)

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data geoResp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func getTimezoneByOffset(lat, lon float64) (*time.Location, error) {
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true",
		lat, lon,
	)

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, err
	}

	return loc, nil
}

var get_local_time = tool.DeclareTool(
	"get_local_time",
	"Get the current time in the city.",
	func(ctx context.Context, args LocalTimeArgs) (any, error) {

		geo, err := fetchGeo(args.City)
		if err != nil {
			return nil, err
		}

		if len(geo.Results) == 0 {
			return nil, fmt.Errorf("city not found")
		}

		r := geo.Results[0]

		loc, err := getTimezoneByOffset(r.Lat, r.Lon)
		if err != nil {
			return nil, err
		}

		now := time.Now().In(loc)

		return map[string]any{
			"city": args.City,
			"lat":  r.Lat,
			"lon":  r.Lon,
			"time": now.Format("2006-01-02 15:04:05"),
		}, nil
	},
)
