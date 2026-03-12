# Airport Data Go

A comprehensive Go library for retrieving airport information by IATA codes, ICAO codes, and various other criteria. This library provides easy access to a large dataset of airports worldwide with detailed information including coordinates, timezone, type, and external links.

## Installation

```bash
go get github.com/aashishvanand/airport-data-go
```

## Features

- Comprehensive airport database with worldwide coverage
- Search by IATA codes, ICAO codes, country, continent, and more
- Geographic proximity search with customizable radius
- External links to Wikipedia, airport websites, and flight tracking services
- Distance calculation between airports
- Filter by airport type (large_airport, medium_airport, small_airport, heliport, seaplane_base)
- Timezone-based airport lookup
- Autocomplete suggestions for search interfaces
- Advanced multi-criteria filtering
- Statistical analysis by country and continent
- Bulk operations for multiple airports
- Code validation utilities
- Airport ranking by runway length and elevation
- Zero external dependencies
- Data embedded at compile time via `//go:embed`

## Airport Data Structure

Each airport is represented by the `Airport` struct:

```go
type Airport struct {
    IATA             string  `json:"iata"`              // 3-letter IATA code
    ICAO             string  `json:"icao"`              // 4-letter ICAO code
    Timezone         string  `json:"time"`              // Timezone identifier
    UTC              float64 `json:"utc"`               // UTC offset
    CountryCode      string  `json:"country_code"`      // 2-letter country code
    Continent        string  `json:"continent"`         // 2-letter continent code (AS, EU, NA, SA, AF, OC, AN)
    Name             string  `json:"airport"`           // Airport name
    Latitude         float64 `json:"latitude"`          // Latitude coordinate
    Longitude        float64 `json:"longitude"`         // Longitude coordinate
    ElevationFt      int     `json:"elevation_ft"`      // Elevation in feet
    Type             string  `json:"type"`              // Airport type
    ScheduledService string  `json:"scheduled_service"` // Has scheduled commercial service
    Wikipedia        string  `json:"wikipedia"`         // Wikipedia URL
    Website          string  `json:"website"`           // Official website URL
    RunwayLength     int     `json:"runway_length"`     // Longest runway in feet
    Flightradar24URL string  `json:"flightradar24_url"` // Flightradar24 URL
    RadarboxURL      string  `json:"radarbox_url"`      // Radarbox URL
    FlightawareURL   string  `json:"flightaware_url"`   // FlightAware URL
}
```

## Basic Usage

```go
package main

import (
    "fmt"
    "log"

    airportdata "github.com/aashishvanand/airport-data-go"
)

func main() {
    // Get airport by IATA code
    airports, err := airportdata.GetAirportByIata("SIN")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(airports[0].Name) // "Singapore Changi Airport"

    // Get airport by ICAO code
    airports, err = airportdata.GetAirportByIcao("WSSS")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(airports[0].CountryCode) // "SG"

    // Search airports by name
    airports, err = airportdata.SearchByName("Singapore")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d airports matching 'Singapore'\n", len(airports))

    // Find nearby airports (within 50km of coordinates)
    nearby, err := airportdata.FindNearbyAirports(1.35019, 103.994003, 50)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d airports near Singapore Changi\n", len(nearby))
}
```

## API Reference

### Core Search Functions

#### `GetAirportByIata(iataCode string) ([]Airport, error)`

Finds airports by their 3-letter IATA code.

```go
airports, err := airportdata.GetAirportByIata("LHR")
```

#### `GetAirportByIcao(icaoCode string) ([]Airport, error)`

Finds airports by their 4-character ICAO code.

```go
airports, err := airportdata.GetAirportByIcao("EGLL")
```

#### `SearchByName(query string) ([]Airport, error)`

Searches for airports by name (case-insensitive, minimum 2 characters).

```go
airports, err := airportdata.SearchByName("Heathrow")
```

### Geographic Functions

#### `FindNearbyAirports(lat, lon, radiusKm float64) ([]NearbyAirport, error)`

Finds airports within a specified radius of given coordinates.

```go
nearby, err := airportdata.FindNearbyAirports(51.5074, -0.1278, 100)
// Returns airports within 100km of London coordinates
```

#### `CalculateDistance(code1, code2 string) (float64, error)`

Calculates the great-circle distance between two airports using IATA or ICAO codes.

```go
distance, err := airportdata.CalculateDistance("LHR", "JFK")
// Returns distance in kilometers (approximately 5540)
```

#### `FindNearestAirport(lat, lon float64, filter *FindAirportsFilter) (*NearbyAirport, error)`

Finds the single nearest airport to given coordinates, optionally with filters.

```go
// Find nearest airport to coordinates
nearest, err := airportdata.FindNearestAirport(1.35019, 103.994003, nil)

// Find nearest large airport with scheduled service
nearest, err := airportdata.FindNearestAirport(1.35019, 103.994003, &airportdata.FindAirportsFilter{
    Type:                "large_airport",
    HasScheduledService: true,
})
```

### Filtering Functions

#### `GetAirportByCountryCode(countryCode string) ([]Airport, error)`

Finds all airports in a specific country.

```go
usAirports, err := airportdata.GetAirportByCountryCode("US")
```

#### `GetAirportByContinent(continentCode string) ([]Airport, error)`

Finds all airports on a specific continent.

```go
asianAirports, err := airportdata.GetAirportByContinent("AS")
// Continent codes: AS, EU, NA, SA, AF, OC, AN
```

#### `GetAirportsByType(airportType string) ([]Airport, error)`

Finds airports by their type.

```go
largeAirports, err := airportdata.GetAirportsByType("large_airport")
// Available types: large_airport, medium_airport, small_airport, heliport, seaplane_base

// Convenience search for all airports
allAirports, err := airportdata.GetAirportsByType("airport")
// Returns large_airport, medium_airport, and small_airport
```

#### `GetAirportsByTimezone(timezone string) ([]Airport, error)`

Finds all airports within a specific timezone.

```go
londonAirports, err := airportdata.GetAirportsByTimezone("Europe/London")
```

#### `FindAirports(filter FindAirportsFilter) ([]Airport, error)`

Finds airports matching multiple criteria.

```go
// Find large airports in Great Britain with scheduled service
airports, err := airportdata.FindAirports(airportdata.FindAirportsFilter{
    CountryCode:         "GB",
    Type:                "large_airport",
    HasScheduledService: true,
})

// Find airports with minimum runway length
longRunwayAirports, err := airportdata.FindAirports(airportdata.FindAirportsFilter{
    MinRunwayFt: 10000,
})
```

### Advanced Functions

#### `GetAutocompleteSuggestions(query string) ([]Airport, error)`

Provides autocomplete suggestions for search interfaces (returns max 10 results).

```go
suggestions, err := airportdata.GetAutocompleteSuggestions("Lon")
```

#### `GetAirportLinks(code string) (*AirportLinks, error)`

Gets external links for an airport using IATA or ICAO code.

```go
links, err := airportdata.GetAirportLinks("SIN")
// Returns:
// AirportLinks{
//     Website:      "https://www.changiairport.com",
//     Wikipedia:    "https://en.wikipedia.org/wiki/Singapore_Changi_Airport",
//     Flightradar24: "https://www.flightradar24.com/airport/SIN",
//     Radarbox:     "https://www.radarbox.com/airport/WSSS",
//     Flightaware:  "https://www.flightaware.com/live/airport/WSSS",
// }
```

### Statistical & Analytical Functions

#### `GetAirportStatsByCountry(countryCode string) (*AirportStats, error)`

Gets comprehensive statistics about airports in a specific country.

```go
stats, err := airportdata.GetAirportStatsByCountry("US")
// stats.Total, stats.ByType, stats.WithScheduledService, etc.
```

#### `GetAirportStatsByContinent(continentCode string) (*ContinentStats, error)`

Gets comprehensive statistics about airports on a specific continent.

```go
stats, err := airportdata.GetAirportStatsByContinent("AS")
```

#### `GetLargestAirportsByContinent(continentCode string, limit int, sortBy string) ([]Airport, error)`

Gets the largest airports on a continent by runway length or elevation.

```go
// Get top 5 airports in Asia by runway length
airports, err := airportdata.GetLargestAirportsByContinent("AS", 5, "runway")

// Get top 10 airports in South America by elevation
highAltitude, err := airportdata.GetLargestAirportsByContinent("SA", 10, "elevation")
```

### Bulk Operations

#### `GetMultipleAirports(codes []string) ([]*Airport, error)`

Fetches multiple airports by their IATA or ICAO codes in one call.

```go
airports, err := airportdata.GetMultipleAirports([]string{"SIN", "LHR", "JFK", "WSSS"})
// Returns slice of *Airport (nil for codes not found)
```

#### `CalculateDistanceMatrix(codes []string) (*DistanceMatrix, error)`

Calculates distances between all pairs of airports in a list.

```go
matrix, err := airportdata.CalculateDistanceMatrix([]string{"SIN", "LHR", "JFK"})
// matrix.Distances["SIN"]["LHR"] => distance in km
// matrix.Distances["LHR"]["JFK"] => distance in km
```

### Validation & Utilities

#### `ValidateIataCode(code string) (bool, error)`

Validates if an IATA code exists in the database.

```go
valid, err := airportdata.ValidateIataCode("SIN") // true, nil
valid, err = airportdata.ValidateIataCode("XYZ")  // false, nil
```

#### `ValidateIcaoCode(code string) (bool, error)`

Validates if an ICAO code exists in the database.

```go
valid, err := airportdata.ValidateIcaoCode("WSSS") // true, nil
```

#### `GetAirportCount(filter *FindAirportsFilter) (int, error)`

Gets the count of airports matching the given filters.

```go
// Get total airport count
total, err := airportdata.GetAirportCount(nil)

// Get count of large airports in the US
count, err := airportdata.GetAirportCount(&airportdata.FindAirportsFilter{
    CountryCode: "US",
    Type:        "large_airport",
})
```

#### `IsAirportOperational(code string) (bool, error)`

Checks if an airport has scheduled commercial service.

```go
operational, err := airportdata.IsAirportOperational("SIN") // true, nil
```

## Error Handling

All functions return an error as the last return value. Errors are returned for invalid input or when no data is found.

```go
airports, err := airportdata.GetAirportByIata("XYZ")
if err != nil {
    fmt.Println(err) // "no data found for IATA code: XYZ"
}
```

## Data Source

This library uses a comprehensive dataset of worldwide airports with regular updates to ensure accuracy and completeness.

## Related Libraries

- [airport-data-js](https://github.com/aashishvanand/airport-data-js) - JavaScript/Node.js
- [airport-data-swift](https://github.com/aashishvanand/airport-data-swift) - Swift
- [airport-data-dart](https://github.com/aashishvanand/airport-data-dart) - Dart/Flutter
- [airport-data-kotlin](https://github.com/aashishvanand/airport-data-kotlin) - Kotlin
- [airport-data-rust](https://github.com/aashishvanand/airport-data-rust) - Rust
- [airports-py](https://github.com/aashishvanand/airports-py) - Python

## License

This project is licensed under the Creative Commons Attribution 4.0 International (CC BY 4.0) - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the [issues page](https://github.com/aashishvanand/airport-data-go/issues).
