package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/hobochili/rhcc/config"
	owm "github.com/hobochili/rhcc/openweathermap"
)

const (
	defaultTimeLocation = "America/Chicago"
)

func main() {
	var outputJSON bool

	app := &cli.App{
		Name:  "rhcc",
		Usage: "Output optimal contact method based on OpenWeatherMap 5-day forecast",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "json",
				Usage:       "Output in JSON Format",
				Destination: &outputJSON,
			},
		},
		Action: func(c *cli.Context) error {
			cfg, err := config.LoadConfig()
			if err != nil {
				log.Error(err)
				cli.ShowCommandHelp(c, c.Args().First())
				os.Exit(1)
			}

			level, err := log.ParseLevel(cfg.LogLevel)
			if err != nil {
				log.Error(err)
				cli.ShowCommandHelp(c, c.Args().First())
				os.Exit(1)
			}

			log.SetLevel(level)

			cfg.OutputJSON = outputJSON

			return run(cfg)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func run(c *config.Config) error {
	// Create an OpenWeatherMap 5-day forecast API client
	f, err := owm.NewForecast(c.OWMKey)
	if err != nil {
		return err
	}

	// Get the 5-day forecast
	if err = f.Get5DayForecast(); err != nil {
		return err
	}

	if c.OutputJSON {
		return outputJSON(f)
	}

	return outputTable(f)
}

// Output results in JSON format to stdout.
func outputJSON(f *owm.Forecast) error {
	type data struct {
		Timestamp     string `json:"timestamp"`
		ContactMethod string `json:"contactMethod"`
	}

	d := make([]data, len(f.FiveDay))

	for i, forecast := range f.FiveDay {
		d[i] = data{
			Timestamp:     forecast.Time.Format(time.RFC3339),
			ContactMethod: contactMethod(forecast, f.Sunrise, f.Sunset),
		}
	}

	result, err := json.Marshal(d)
	if err != nil {
		return err
	}

	os.Stdout.Write(result)

	return nil
}

// Output results in table format to stdout.
func outputTable(f *owm.Forecast) error {
	table := tablewriter.NewWriter(os.Stdout)

	weekday := ""
	timeZone := ""
	first := true

	for _, forecast := range f.FiveDay {
		local := local(forecast.Time)
		dateString := ""

		// Print a row separator and the date if this result marks a new day
		if local.Weekday().String() != weekday {
			weekday = local.Weekday().String()

			if !first {
				// HACK: Use table.AddSeparator() if this gets merged:
				// => https://github.com/olekukonko/tablewriter/pull/158
				table.Append([]string{
					"-----------------------------",
					"------------",
					"-------",
				})
			}

			dateString = fmt.Sprintf("%v, %v %v, %v",
				local.Weekday().String(),
				local.Month().String(),
				local.Day(),
				local.Year(),
			)
		}

		table.Append([]string{
			dateString,
			fmt.Sprintf("%v ", local.Format("03:04:05 PM")),
			contactMethod(forecast, f.Sunrise, f.Sunset),
		})

		// Time zone should be constant across all forecast results, so just pull
		// it from the first result.
		if first {
			timeZone = local.Format("MST")
			first = false
		}
	}

	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{
		"Date",
		fmt.Sprintf("Time (%v)", timeZone),
		"Contact\nMethod"},
	)

	table.Render()

	return nil
}

// Return the best contact method based on the given forecast.
func contactMethod(rf *owm.RelevantForecast, sunrise, sunset time.Time) string {
	if isColdOrRainy(rf) {
		return "phone"
	}

	if isWarmAndSunny(rf, sunrise, sunset) {
		return "text"
	}

	if isNice(rf) {
		return "email"
	}

	return "default"
}

// Return true if the forecasted temperature is less than 55 degrees Fahrenheit
// or if there is any rain in the forecast.
func isColdOrRainy(rf *owm.RelevantForecast) bool {
	if rf.Temp < 55 || rf.Rain > 0 {
		return true
	}

	return false
}

// Return true when all of the following conditions are met:
//   - temperature is greater than 75 degress Fahrenheit
//   - cloud coverage is not greater than 50%
//   - forecast time is between sunrise and sunset
//
// Note: for simplicity we consider today's sunrise/sunset timestamps as
// 			 constant over the next 5 days.
func isWarmAndSunny(rf *owm.RelevantForecast, sunrise, sunset time.Time) bool {
	if rf.Temp > 75 &&
		rf.CloudPercentage <= 50 &&
		rf.Time.After(sunrise) &&
		rf.Time.Before(sunset) {
		return true
	}

	return false
}

// Return true if the forecasted temperature is between 55 and 75 degrees
// Fahrenheit.
func isNice(rf *owm.RelevantForecast) bool {
	if rf.Temp >= 55 && rf.Temp <= 75 {
		return true
	}

	return false
}

func local(t time.Time) time.Time {
	location, _ := time.LoadLocation(defaultTimeLocation)
	return t.In(location)
}
