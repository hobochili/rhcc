# RHCC

A command line utility to output the optimal customer contact method based on [OpenWeatherMap's](https://openweathermap.org/) 5-day
Minneapolis, MN forecast.

## Build

Building this tool requires a Go 1.11 or later environment. It was developed and tested with Go 1.13. Official instructions for setting up this environment can be found [here](https://golang.org/doc/install), but on a Linux system you would do the following:

```BASH
wget go1.13.10.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.14.2.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

Then, clone and build `rhcc`:

```BASH
git clone github.com/hobochili/rhcc.git
cd rhcc
go build
```

## Run

In order to fetch forecast results you must first acquire an OpenWeatherMap API key and store it
in the `RHCC_OWM_API_KEY` environment variable. To obtain the key:

1. Sign up for an account at https://home.openweathermap.org/users/sign_up
1. Collect your API key (or create a new one) from https://home.openweathermap.org/api_keys

Once you have obtained your key you can set the required environment variable and execute:

```BASH
export RHCC_OWM_API_KEY=<OWM_KEY>
./rhcc
```

## Usage

By default, `rfcc` displays results in table format to stdout.

To output the results in JSON format, run: `./rfcc --json`.

## Assumptions

This utility assumes that 5-day forecast results come in spans of 3 hours, with eight 3-hour spans comprising a single day.

### When is it considered sunny?

The optimal contact method is dependent on sunniness. We consider it to be sunny weather if the sky condition can be described as sunny or mostly sunny. It is never considered sunny if the sky condition can be described as mostly cloudy or overcast.

Based on [this](https://www.weather.gov/media/pah/ServiceGuide/A-forecast.pdf) outline provided by [https://weather.gov](https://weather.gov) we consider it to be sunny if there is not more than 50% cloud coverage.

### When is it considered rainy?

We consider it to be rainy if the forecasted rainfall is greater than 0 mm.

## TODO

- Write tests
- Make location configurable (only results for Minneapolis, MN are currently supported)
- Add optional filters for display results (ex: don't show results for forecasts outside of business hours)
