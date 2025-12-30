package openmeteo

import "time"

// CurrentWeather represents a complete snapshot of current weather conditions at a specific location.
// All weather parameter fields use metric units (°C, m/s, mm, hPa, %).
// Zero values indicate the absence of data from the API or that the measurement is zero (e.g., 0mm precipitation).
type CurrentWeather struct {
	// Latitude of the weather observation location in degrees (-90 to 90)
	Latitude float64

	// Longitude of the weather observation location in degrees (-180 to 180)
	Longitude float64

	// Time of the weather observation in UTC
	Time time.Time

	// Temperature is the air temperature at 2 meters height in degrees Celsius
	Temperature float64

	// RelativeHumidity is the relative humidity at 2 meters height in percent (0-100)
	RelativeHumidity float64

	// ApparentTemperature is the perceived "feels like" temperature in degrees Celsius
	ApparentTemperature float64

	// IsDay indicates whether it is currently daytime (true) or nighttime (false)
	IsDay bool

	// Precipitation is the total precipitation (rain + snow) in millimeters
	Precipitation float64

	// Rain is the liquid rain amount in millimeters
	Rain float64

	// Showers is the shower precipitation amount in millimeters
	Showers float64

	// Snowfall is the snowfall amount in millimeters (water equivalent)
	Snowfall float64

	// WeatherCode is the WMO weather code (0-99) indicating general weather conditions
	WeatherCode int

	// CloudCover is the total cloud cover in percent (0-100)
	CloudCover float64

	// PressureMSL is the atmospheric pressure reduced to sea level in hectopascals
	PressureMSL float64

	// SurfacePressure is the atmospheric pressure at surface level in hectopascals
	SurfacePressure float64

	// WindSpeed is the wind speed at 10 meters height in meters per second
	WindSpeed float64

	// WindDirection is the wind direction at 10 meters height in degrees (0-360)
	WindDirection float64

	// WindGusts is the maximum wind gust speed at 10 meters height in meters per second
	WindGusts float64
}

// weatherResponse is an internal structure for unmarshaling JSON responses from the Open Meteo API.
// It uses pointer types to detect null values from the API.
type weatherResponse struct {
	Latitude       float64                `json:"latitude"`
	Longitude      float64                `json:"longitude"`
	CurrentWeather currentWeatherResponse `json:"current"`
}

// currentWeatherResponse is an internal structure for unmarshaling the current_weather object
// from the Open Meteo API JSON response. Pointer types allow detection of null/missing values.
type currentWeatherResponse struct {
	Time                *string  `json:"time"`
	Temperature         *float64 `json:"temperature_2m"`
	Windspeed           *float64 `json:"wind_speed_10m"`
	Winddirection       *float64 `json:"wind_direction_10m"`
	Weathercode         *int     `json:"weather_code"`
	IsDay               *int     `json:"is_day"`
	RelativeHumidity    *float64 `json:"relative_humidity_2m"`
	ApparentTemperature *float64 `json:"apparent_temperature"`
	Precipitation       *float64 `json:"precipitation"`
	Rain                *float64 `json:"rain"`
	Showers             *float64 `json:"showers"`
	Snowfall            *float64 `json:"snowfall"`
	CloudCover          *float64 `json:"cloud_cover"`
	PressureMSL         *float64 `json:"pressure_msl"`
	SurfacePressure     *float64 `json:"surface_pressure"`
	WindGusts           *float64 `json:"wind_gusts_10m"`
}
