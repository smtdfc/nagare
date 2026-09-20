package declarations

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type OpenMeteoApiResponse struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationTimeMs     float64 `json:"generationtime_ms"`
	UtcOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`
	CurrentUnits         struct {
		Time               string `json:"time"`
		Interval           string `json:"interval"`
		Temperature2M      string `json:"temperature_2m"`
		RelativeHumidity2M string `json:"relative_humidity_2m"`
		WindSpeed10M       string `json:"wind_speed_10m"`
		WeatherCode        string `json:"weathercode"`
	} `json:"current_units"`
	Current struct {
		Time               string  `json:"time"`
		Interval           int     `json:"interval"`
		Temperature2M      float64 `json:"temperature_2m"`
		RelativeHumidity2M int     `json:"relative_humidity_2m"`
		WindSpeed10M       float64 `json:"wind_speed_10m"`
		WeatherCode        int     `json:"weathercode"`
	} `json:"current"`
}

type WeatherToolInput struct {
	Lat float32
	Lng float32
}

type WeatherValue[T any] struct {
	Value T      `json:"value"`
	Unit  string `json:"unit"`
}

type WeatherToolOutput struct {
	Temperature      *WeatherValue[float64] `json:"temperature"`
	WindSpeed        *WeatherValue[float64] `json:"wind_speed"`
	RelativeHumidity *WeatherValue[int]     `json:"relative_humidity"`
	WeatherCode      *WeatherValue[int]     `json:"weather_code"`
}

var WeatherTool = tool.DefineTool(
	"weather_tool",
	"Get weather",
	func(ctx *context.ExecuteContext, args *WeatherToolInput, _ tool.Bindings) (*WeatherToolOutput, error) {
		apiURL := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,relative_humidity_2m,wind_speed_10m,weathercode",
			args.Lat, args.Lng)
		req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
		if err != nil {
			return nil, err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var weatherData OpenMeteoApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
			return nil, fmt.Errorf("failed to decode weather response: %w", err)
		}

		return &WeatherToolOutput{
			Temperature: &WeatherValue[float64]{
				Value: weatherData.Current.Temperature2M,
				Unit:  weatherData.CurrentUnits.Temperature2M,
			},
			WindSpeed: &WeatherValue[float64]{
				Value: weatherData.Current.WindSpeed10M,
				Unit:  weatherData.CurrentUnits.WindSpeed10M,
			},
			RelativeHumidity: &WeatherValue[int]{
				Value: weatherData.Current.RelativeHumidity2M,
				Unit:  weatherData.CurrentUnits.RelativeHumidity2M,
			},
			WeatherCode: &WeatherValue[int]{
				Value: weatherData.Current.WeatherCode,
				Unit:  weatherData.CurrentUnits.WeatherCode,
			},
		}, nil
	},
)
