package declarations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/tool"
)

type WeatherArgs struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type weatherResponse struct {
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
		WindSpeed   float64 `json:"windspeed"`
		Time        string  `json:"time"`
		WeatherCode int     `json:"weathercode"`
	} `json:"current_weather"`
}

func callWeatherAPI(lat, lon float64) (*weatherResponse, error) {
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

	var data weatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

var get_weather = tool.DeclareTool(
	"get_weather",
	"Get weather information by coordinates.",
	func(ctx domains.AgentContext, args WeatherArgs) (any, error) {

		data, err := callWeatherAPI(args.Lat, args.Lon)
		if err != nil {
			return nil, err
		}

		return map[string]any{
			"temperature":  data.CurrentWeather.Temperature,
			"wind_speed":   data.CurrentWeather.WindSpeed,
			"time":         data.CurrentWeather.Time,
			"weather_code": data.CurrentWeather.WeatherCode,
		}, nil
	},
	domains.STATIC_TOOL,
	domains.NO_GROUP,
)
