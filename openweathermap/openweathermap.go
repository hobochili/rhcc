package openweathermap

import (
	"fmt"
	"time"

	owm "github.com/briandowns/openweathermap"
)

const (
	defaultForecastType    = "5"
	defaultTemperatureUnit = "F"
	defaultLanguageCode    = "EN"
	defaultLocationName    = "Minneapolis,US"
)

// Forecast holds the config for an OWM Forecast client
type Forecast struct {
	client   *owm.ForecastWeatherData
	current  *owm.CurrentWeatherData
	location string

	FiveDay []*RelevantForecast
	Sunrise time.Time
	Sunset  time.Time
}

// RelevantForecast holds the relevant data from a forecast list
type RelevantForecast struct {
	Time            time.Time
	Temp            float64
	CloudPercentage int
	Rain            float64
}

// NewForecast returns a new Forecast client
func NewForecast(apiKey string) (*Forecast, error) {
	client, err := owm.NewForecast(
		defaultForecastType,
		defaultTemperatureUnit,
		defaultLanguageCode,
		apiKey,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create forecast client: %v", err)
	}

	// Get the current weather to extract sunrise and sunset timestamps.

	current, err := owm.NewCurrent(defaultTemperatureUnit, defaultLanguageCode, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create current forecast client: %v", err)
	}

	if err := current.CurrentByName(defaultLocationName); err != nil {
		return nil, fmt.Errorf("failed to fetch sunrise/sunset data: %v", err)
	}

	return &Forecast{
		client:   client,
		current:  current,
		location: defaultLocationName,

		FiveDay: nil,
		Sunrise: time.Unix(int64(current.Sys.Sunrise), 0),
		Sunset:  time.Unix(int64(current.Sys.Sunset), 0),
	}, nil
}

// Get5DayForecast returns forecast data points for the next 5 days
func (f *Forecast) Get5DayForecast() error {

	// Fetch 5 days worth of data, each day containing 8 data points
	if err := f.client.DailyByName(f.location, 5*8); err != nil {
		return err
	}

	data := f.client.ForecastWeatherJson.(*owm.Forecast5WeatherData)

	for _, s := range data.List {
		f.FiveDay = append(f.FiveDay, &RelevantForecast{
			Time:            s.DtTxt.Add(0 * time.Second),
			Temp:            s.Main.Temp,
			CloudPercentage: s.Clouds.All,
			Rain:            s.Rain.ThreeH,
		})
	}

	return nil
}
