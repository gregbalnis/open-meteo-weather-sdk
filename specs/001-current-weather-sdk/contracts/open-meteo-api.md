# Open Meteo API Contract

**Date**: 2025-12-28
**API Version**: v1
**Documentation**: https://open-meteo.com/en/docs

## Overview

This document defines the external API contract between the SDK and the Open Meteo Weather API. This contract is based on Open Meteo's public documentation and has been validated through research.

## Base URL

```
https://api.open-meteo.com/v1
```

**Protocol**: HTTPS only
**Authentication**: None required for current weather endpoint

---

## Endpoint: Get Forecast (Current Weather)

### Request

**Method**: `GET`

**Path**: `/forecast`

**Query Parameters**:

| Parameter | Type | Required | Description | Example |
|-----------|------|----------|-------------|---------|
| `latitude` | float | Yes | Latitude in decimal degrees | `52.52` |
| `longitude` | float | Yes | Longitude in decimal degrees | `13.41` |
| `current_weather` | boolean | Yes | Enable current weather data | `true` |
| `temperature_unit` | string | No | Temperature unit (celsius/fahrenheit) | `celsius` |
| `windspeed_unit` | string | No | Wind speed unit (kmh/ms/mph/kn) | `ms` |
| `precipitation_unit` | string | No | Precipitation unit (mm/inch) | `mm` |

**SDK Behavior**:
- Always sets `current_weather=true`
- Uses metric units only (per clarifications):
  - `temperature_unit=celsius`
  - `windspeed_unit=ms`
  - `precipitation_unit=mm`

**Example Request**:
```
GET https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&current_weather=true&temperature_unit=celsius&windspeed_unit=ms&precipitation_unit=mm
```

---

### Response

**Status Codes**:
- `200 OK` - Success
- `400 Bad Request` - Invalid parameters
- `500 Internal Server Error` - API error
- `503 Service Unavailable` - API temporarily unavailable

**Content-Type**: `application/json`

**Response Body** (Success):

```json
{
  "latitude": 52.52,
  "longitude": 13.41,
  "generationtime_ms": 0.123,
  "utc_offset_seconds": 0,
  "timezone": "GMT",
  "timezone_abbreviation": "GMT",
  "elevation": 38.0,
  "current_weather": {
    "time": "2025-12-28T10:00",
    "interval": 900,
    "temperature": 15.3,
    "windspeed": 12.5,
    "winddirection": 270,
    "is_day": 1,
    "weathercode": 3
  },
  "current_weather_units": {
    "time": "iso8601",
    "interval": "seconds",
    "temperature": "°C",
    "windspeed": "m/s",
    "winddirection": "°",
    "is_day": "",
    "weathercode": "wmo code"
  }
}
```

**Field Descriptions**:

| Field Path | Type | Nullable | Description |
|------------|------|----------|-------------|
| `latitude` | float | No | Resolved latitude (may differ slightly from request) |
| `longitude` | float | No | Resolved longitude (may differ slightly from request) |
| `current_weather.time` | string | No | ISO 8601 timestamp (UTC) |
| `current_weather.temperature` | float | Yes* | Air temperature at 2m (°C) |
| `current_weather.windspeed` | float | Yes* | Wind speed at 10m (m/s) |
| `current_weather.winddirection` | float | Yes* | Wind direction (degrees, 0=N) |
| `current_weather.weathercode` | int | No | WMO weather code (0-99) |
| `current_weather.is_day` | int | No | 1=day, 0=night |

*Fields marked nullable may return `null` depending on location/conditions

**Note**: The API response shown above is simplified. The actual API may return additional fields in `current_weather`, but the SDK focuses on the 15 parameters specified in the feature spec. Additional fields discovered during implementation:
- `relativehumidity_2m`
- `apparent_temperature`
- `precipitation`
- `rain`
- `showers`
- `snowfall`
- `cloudcover`
- `pressure_msl` (sea level pressure)
- `surface_pressure`
- `wind_gusts_10m`

These are included by adding the following query parameters:
- `current=temperature_2m,relativehumidity_2m,apparent_temperature,is_day,precipitation,rain,showers,snowfall,weathercode,cloudcover,pressure_msl,surface_pressure,windspeed_10m,winddirection_10m,windgusts_10m`

**Revised Example Request**:
```
GET https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&current=temperature_2m,relativehumidity_2m,apparent_temperature,is_day,precipitation,rain,showers,snowfall,weathercode,cloudcover,pressure_msl,surface_pressure,windspeed_10m,winddirection_10m,windgusts_10m&temperature_unit=celsius&windspeed_unit=ms&precipitation_unit=mm
```

**Revised Response Body**:
```json
{
  "latitude": 52.52,
  "longitude": 13.41,
  "current": {
    "time": "2025-12-28T10:00",
    "temperature_2m": 15.3,
    "relativehumidity_2m": 65.0,
    "apparent_temperature": 14.1,
    "is_day": 1,
    "precipitation": 0.0,
    "rain": 0.0,
    "showers": 0.0,
    "snowfall": 0.0,
    "weathercode": 3,
    "cloudcover": 75.0,
    "pressure_msl": 1013.25,
    "surface_pressure": 1010.0,
    "windspeed_10m": 12.5,
    "winddirection_10m": 270.0,
    "windgusts_10m": 18.0
  }
}
```

---

### Error Responses

**400 Bad Request**:
```json
{
  "error": true,
  "reason": "Latitude must be in range of -90 to 90°. Given: 999.0"
}
```

**500 Internal Server Error**:
```json
{
  "error": true,
  "reason": "Internal server error"
}
```

**SDK Behavior**:
- Parse `error` field to detect error responses
- Extract `reason` field for error message
- Wrap in `ErrorTypeAPI` error

---

## Data Types

### WMO Weather Code

Integer codes 0-99 representing weather conditions:

| Code | Description |
|------|-------------|
| 0 | Clear sky |
| 1-3 | Mainly clear, partly cloudy, overcast |
| 45, 48 | Fog and depositing rime fog |
| 51-55 | Drizzle (light to heavy) |
| 61-65 | Rain (slight to heavy) |
| 71-75 | Snow fall (slight to heavy) |
| 77 | Snow grains |
| 80-82 | Rain showers (slight to violent) |
| 85, 86 | Snow showers (slight and heavy) |
| 95 | Thunderstorm |
| 96, 99 | Thunderstorm with hail |

**Reference**: https://www.nodc.noaa.gov/archive/arc0021/0002199/1.1/data/0-data/HTML/WMO-CODE/WMO4677.HTM

---

## Rate Limits

**Free Tier**:
- 10,000 API calls per day
- No burst limit documented
- No API key required

**SDK Behavior**:
- Does not implement rate limiting (per clarifications)
- Applications should implement their own rate limiting if needed
- Concurrent request limit (10) is for resource protection, not API compliance

---

## Caching Headers

The API may return standard HTTP caching headers:
- `Cache-Control`
- `ETag`
- `Last-Modified`

**SDK Behavior**:
- Does not implement caching (per clarifications/out of scope)
- Applications can implement custom HTTP clients with caching if needed

---

## Reliability

**Expected Availability**: 99.9%+ (based on free tier)

**Typical Response Time**: 100-500ms

**Failure Modes**:
- Network timeouts (handled by SDK with timeout)
- DNS failures (wrapped as `ErrorTypeNetwork`)
- HTTP 5xx errors (wrapped as `ErrorTypeAPI`)
- Malformed JSON (wrapped as `ErrorTypeAPI`)

---

## Breaking Changes

**API Versioning**: Path includes `/v1/` for versioning

**Stability**: Open Meteo maintains backward compatibility within major versions

**SDK Strategy**: 
- Pin to `/v1/` endpoint
- If API changes, increment SDK major version
- Monitor Open Meteo changelog for breaking changes

---

## Testing

**Mock Responses**: For unit tests, use `httptest.Server` with sample JSON

**Integration Tests**: Optional tests against real API (build tag: `integration`)

**Contract Tests**: Validate JSON schema matches expected structure

---

## Dependencies

**None**: SDK uses standard library only
- `net/http` for requests
- `encoding/json` for parsing
- No third-party libraries

---

## Security

**TLS**: All requests use HTTPS (TLS 1.2+)

**Input Validation**: Coordinates validated before request to prevent injection

**No Secrets**: No API keys or authentication required
