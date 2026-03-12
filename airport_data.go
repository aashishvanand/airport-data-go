// Package airportdata provides a comprehensive airport database with worldwide coverage.
// It supports search by IATA/ICAO codes, geographic proximity, filtering by country/continent/type,
// statistical analysis, bulk operations, and validation utilities.
//
// The airport data is embedded in the binary at compile time, so no external files are needed at runtime.
package airportdata

import (
	"embed"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

//go:embed data/airports.json
var airportsJSON embed.FS

// Airport represents a single airport entry with all available metadata.
type Airport struct {
	IATA             string  `json:"iata"`
	ICAO             string  `json:"icao"`
	Timezone         string  `json:"time"`
	UTC              float64 `json:"utc"`
	CountryCode      string  `json:"country_code"`
	Continent        string  `json:"continent"`
	Name             string  `json:"airport"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	ElevationFt      int     `json:"elevation_ft"`
	Type             string  `json:"type"`
	ScheduledService string  `json:"scheduled_service"`
	Wikipedia        string  `json:"wikipedia"`
	Website          string  `json:"website"`
	RunwayLength     int     `json:"runway_length"`
	Flightradar24URL string  `json:"flightradar24_url"`
	RadarboxURL      string  `json:"radarbox_url"`
	FlightawareURL   string  `json:"flightaware_url"`
}

// rawAirport is used for flexible JSON unmarshaling where fields may have mixed types.
type rawAirport struct {
	IATA             string          `json:"iata"`
	ICAO             json.RawMessage `json:"icao"`
	Timezone         string          `json:"time"`
	UTC              json.RawMessage `json:"utc"`
	CountryCode      string          `json:"country_code"`
	Continent        string          `json:"continent"`
	Name             string          `json:"airport"`
	Latitude         json.RawMessage `json:"latitude"`
	Longitude        json.RawMessage `json:"longitude"`
	ElevationFt      json.RawMessage `json:"elevation_ft"`
	Type             string          `json:"type"`
	ScheduledService string          `json:"scheduled_service"`
	Wikipedia        string          `json:"wikipedia"`
	Website          string          `json:"website"`
	RunwayLength     json.RawMessage `json:"runway_length"`
	Flightradar24URL string          `json:"flightradar24_url"`
	RadarboxURL      string          `json:"radarbox_url"`
	FlightawareURL   string          `json:"flightaware_url"`
}

func rawToString(raw json.RawMessage) string {
	if raw == nil || string(raw) == "null" {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	// If it's a number, convert to string
	var f float64
	if err := json.Unmarshal(raw, &f); err == nil {
		return strconv.FormatFloat(f, 'f', -1, 64)
	}
	return ""
}

func rawToFloat64(raw json.RawMessage) float64 {
	if raw == nil || string(raw) == "null" {
		return 0
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err == nil {
		return f
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		if s == "" {
			return 0
		}
		v, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return v
		}
	}
	return 0
}

func rawToInt(raw json.RawMessage) int {
	if raw == nil || string(raw) == "null" {
		return 0
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err == nil {
		return int(f)
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		if s == "" {
			return 0
		}
		v, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return int(v)
		}
	}
	return 0
}

func convertRawAirport(r rawAirport) Airport {
	return Airport{
		IATA:             r.IATA,
		ICAO:             rawToString(r.ICAO),
		Timezone:         r.Timezone,
		UTC:              rawToFloat64(r.UTC),
		CountryCode:      r.CountryCode,
		Continent:        r.Continent,
		Name:             r.Name,
		Latitude:         rawToFloat64(r.Latitude),
		Longitude:        rawToFloat64(r.Longitude),
		ElevationFt:      rawToInt(r.ElevationFt),
		Type:             r.Type,
		ScheduledService: r.ScheduledService,
		Wikipedia:        r.Wikipedia,
		Website:          r.Website,
		RunwayLength:     rawToInt(r.RunwayLength),
		Flightradar24URL: r.Flightradar24URL,
		RadarboxURL:      r.RadarboxURL,
		FlightawareURL:   r.FlightawareURL,
	}
}

// NearbyAirport is an Airport with an additional Distance field (in km).
type NearbyAirport struct {
	Airport
	Distance float64 `json:"distance"`
}

// AirportLinks contains external links for an airport.
type AirportLinks struct {
	Website       string `json:"website,omitempty"`
	Wikipedia     string `json:"wikipedia,omitempty"`
	Flightradar24 string `json:"flightradar24,omitempty"`
	Radarbox      string `json:"radarbox,omitempty"`
	Flightaware   string `json:"flightaware,omitempty"`
}

// AirportStats contains statistics about airports in a country.
type AirportStats struct {
	Total                int            `json:"total"`
	ByType               map[string]int `json:"byType"`
	WithScheduledService int            `json:"withScheduledService"`
	AverageRunwayLength  float64        `json:"averageRunwayLength"`
	AverageElevation     float64        `json:"averageElevation"`
	Timezones            []string       `json:"timezones"`
}

// ContinentStats contains statistics about airports on a continent.
type ContinentStats struct {
	Total                int            `json:"total"`
	ByType               map[string]int `json:"byType"`
	ByCountry            map[string]int `json:"byCountry"`
	WithScheduledService int            `json:"withScheduledService"`
	AverageRunwayLength  float64        `json:"averageRunwayLength"`
	AverageElevation     float64        `json:"averageElevation"`
	Timezones            []string       `json:"timezones"`
}

// DistanceMatrixAirport holds basic info for a distance matrix entry.
type DistanceMatrixAirport struct {
	Code string `json:"code"`
	Name string `json:"name"`
	IATA string `json:"iata"`
	ICAO string `json:"icao"`
}

// DistanceMatrix contains the result of a distance matrix calculation.
type DistanceMatrix struct {
	Airports  []DistanceMatrixAirport       `json:"airports"`
	Distances map[string]map[string]float64 `json:"distances"`
}

// FindAirportsFilter defines criteria for advanced airport filtering.
type FindAirportsFilter struct {
	CountryCode         string
	Continent           string
	Type                string
	HasScheduledService *bool
	MinRunwayFt         int
}

var (
	airports     []Airport
	iataIndex    map[string][]int
	icaoIndex    map[string][]int
	countryIndex map[string][]int
	continentIdx map[string][]int
	timezoneIdx  map[string][]int
	typeIndex    map[string][]int
	loadOnce     sync.Once
	loadErr      error
)

func ensureLoaded() error {
	loadOnce.Do(func() {
		data, err := airportsJSON.ReadFile("data/airports.json")
		if err != nil {
			loadErr = fmt.Errorf("failed to read embedded airports data: %w", err)
			return
		}

		var rawAirports []rawAirport
		if err := json.Unmarshal(data, &rawAirports); err != nil {
			loadErr = fmt.Errorf("failed to parse airports data: %w", err)
			return
		}

		airports = make([]Airport, len(rawAirports))
		for i, r := range rawAirports {
			airports[i] = convertRawAirport(r)
		}

		// Build indexes
		iataIndex = make(map[string][]int)
		icaoIndex = make(map[string][]int)
		countryIndex = make(map[string][]int)
		continentIdx = make(map[string][]int)
		timezoneIdx = make(map[string][]int)
		typeIndex = make(map[string][]int)

		for i, a := range airports {
			if a.IATA != "" {
				key := strings.ToUpper(a.IATA)
				iataIndex[key] = append(iataIndex[key], i)
			}
			if a.ICAO != "" {
				key := strings.ToUpper(a.ICAO)
				icaoIndex[key] = append(icaoIndex[key], i)
			}
			if a.CountryCode != "" {
				key := strings.ToUpper(a.CountryCode)
				countryIndex[key] = append(countryIndex[key], i)
			}
			if a.Continent != "" {
				key := strings.ToUpper(a.Continent)
				continentIdx[key] = append(continentIdx[key], i)
			}
			if a.Timezone != "" {
				timezoneIdx[a.Timezone] = append(timezoneIdx[a.Timezone], i)
			}
			if a.Type != "" {
				key := strings.ToLower(a.Type)
				typeIndex[key] = append(typeIndex[key], i)
			}
		}
	})
	return loadErr
}

// ---------------------------------------------------------------------------
// Core Search Functions
// ---------------------------------------------------------------------------

// GetAirportByIata finds airports by their 3-letter IATA code.
// Returns an error if no airport is found.
func GetAirportByIata(iataCode string) ([]Airport, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	code := strings.ToUpper(strings.TrimSpace(iataCode))
	indices, ok := iataIndex[code]
	if !ok || len(indices) == 0 {
		return nil, fmt.Errorf("no data found for IATA code: %s", iataCode)
	}
	result := make([]Airport, len(indices))
	for i, idx := range indices {
		result[i] = airports[idx]
	}
	return result, nil
}

// GetAirportByIcao finds airports by their 4-character ICAO code.
// Returns an error if no airport is found.
func GetAirportByIcao(icaoCode string) ([]Airport, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	code := strings.ToUpper(strings.TrimSpace(icaoCode))
	indices, ok := icaoIndex[code]
	if !ok || len(indices) == 0 {
		return nil, fmt.Errorf("no data found for ICAO code: %s", icaoCode)
	}
	result := make([]Airport, len(indices))
	for i, idx := range indices {
		result[i] = airports[idx]
	}
	return result, nil
}

// SearchByName searches for airports by name (case-insensitive, minimum 2 characters).
func SearchByName(query string) ([]Airport, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	query = strings.TrimSpace(query)
	if len(query) < 2 {
		return nil, fmt.Errorf("search query must be at least 2 characters")
	}
	lowerQuery := strings.ToLower(query)
	var result []Airport
	for _, a := range airports {
		if strings.Contains(strings.ToLower(a.Name), lowerQuery) {
			result = append(result, a)
		}
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// Geographic Functions
// ---------------------------------------------------------------------------

const earthRadiusKm = 6371.0

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// haversineDistance calculates the great-circle distance in km between two points.
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := degToRad(lat2 - lat1)
	dLon := degToRad(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degToRad(lat1))*math.Cos(degToRad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

// FindNearbyAirports finds airports within a specified radius (km) of given coordinates.
func FindNearbyAirports(lat, lon, radiusKm float64) ([]NearbyAirport, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	var result []NearbyAirport
	for _, a := range airports {
		d := haversineDistance(lat, lon, a.Latitude, a.Longitude)
		if d <= radiusKm {
			result = append(result, NearbyAirport{Airport: a, Distance: math.Round(d*100) / 100})
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Distance < result[j].Distance
	})
	return result, nil
}

// CalculateDistance calculates the great-circle distance in km between two airports.
// Accepts IATA or ICAO codes.
func CalculateDistance(code1, code2 string) (float64, error) {
	if err := ensureLoaded(); err != nil {
		return 0, err
	}
	a1, err := resolveAirport(code1)
	if err != nil {
		return 0, err
	}
	a2, err := resolveAirport(code2)
	if err != nil {
		return 0, err
	}
	return haversineDistance(a1.Latitude, a1.Longitude, a2.Latitude, a2.Longitude), nil
}

// FindNearestAirport finds the single nearest airport to given coordinates,
// optionally applying filters. Returns NearbyAirport with distance in km.
func FindNearestAirport(lat, lon float64, filter *FindAirportsFilter) (*NearbyAirport, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}

	var nearest *NearbyAirport
	minDist := math.MaxFloat64

	for _, a := range airports {
		if filter != nil && !matchesFilter(a, filter) {
			continue
		}
		d := haversineDistance(lat, lon, a.Latitude, a.Longitude)
		if d < minDist {
			minDist = d
			nearest = &NearbyAirport{Airport: a, Distance: math.Round(d*100) / 100}
		}
	}

	if nearest == nil {
		return nil, fmt.Errorf("no airport found matching the given criteria")
	}
	return nearest, nil
}

// ---------------------------------------------------------------------------
// Filtering Functions
// ---------------------------------------------------------------------------

// GetAirportByCountryCode finds all airports in a specific country.
func GetAirportByCountryCode(countryCode string) ([]Airport, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	code := strings.ToUpper(strings.TrimSpace(countryCode))
	indices, ok := countryIndex[code]
	if !ok || len(indices) == 0 {
		return nil, fmt.Errorf("no airports found for country code: %s", countryCode)
	}
	result := make([]Airport, len(indices))
	for i, idx := range indices {
		result[i] = airports[idx]
	}
	return result, nil
}

// GetAirportByContinent finds all airports on a specific continent.
// Continent codes: AS, EU, NA, SA, AF, OC, AN.
func GetAirportByContinent(continentCode string) ([]Airport, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	code := strings.ToUpper(strings.TrimSpace(continentCode))
	indices, ok := continentIdx[code]
	if !ok || len(indices) == 0 {
		return nil, fmt.Errorf("no airports found for continent code: %s", continentCode)
	}
	result := make([]Airport, len(indices))
	for i, idx := range indices {
		result[i] = airports[idx]
	}
	return result, nil
}

// GetAirportsByType finds airports by their type.
// Valid types: large_airport, medium_airport, small_airport, heliport, seaplane_base.
// The special type "airport" matches large_airport, medium_airport, and small_airport.
func GetAirportsByType(airportType string) ([]Airport, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	t := strings.ToLower(strings.TrimSpace(airportType))

	if t == "airport" {
		var result []Airport
		for _, key := range []string{"large_airport", "medium_airport", "small_airport"} {
			if indices, ok := typeIndex[key]; ok {
				for _, idx := range indices {
					result = append(result, airports[idx])
				}
			}
		}
		return result, nil
	}

	indices, ok := typeIndex[t]
	if !ok {
		return []Airport{}, nil
	}
	result := make([]Airport, len(indices))
	for i, idx := range indices {
		result[i] = airports[idx]
	}
	return result, nil
}

// GetAirportsByTimezone finds all airports within a specific timezone.
func GetAirportsByTimezone(timezone string) ([]Airport, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	tz := strings.TrimSpace(timezone)
	indices, ok := timezoneIdx[tz]
	if !ok || len(indices) == 0 {
		return nil, fmt.Errorf("no airports found for timezone: %s", timezone)
	}
	result := make([]Airport, len(indices))
	for i, idx := range indices {
		result[i] = airports[idx]
	}
	return result, nil
}

// FindAirports finds airports matching multiple criteria.
func FindAirports(filter FindAirportsFilter) ([]Airport, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	var result []Airport
	for _, a := range airports {
		if matchesFilter(a, &filter) {
			result = append(result, a)
		}
	}
	return result, nil
}

func isScheduledService(a Airport) bool {
	return strings.EqualFold(a.ScheduledService, "TRUE") || strings.EqualFold(a.ScheduledService, "yes")
}

func matchesFilter(a Airport, f *FindAirportsFilter) bool {
	if f.CountryCode != "" && !strings.EqualFold(a.CountryCode, f.CountryCode) {
		return false
	}
	if f.Continent != "" && !strings.EqualFold(a.Continent, f.Continent) {
		return false
	}
	if f.Type != "" {
		ft := strings.ToLower(f.Type)
		at := strings.ToLower(a.Type)
		if ft == "airport" {
			if at != "large_airport" && at != "medium_airport" && at != "small_airport" {
				return false
			}
		} else if at != ft {
			return false
		}
	}
	if f.HasScheduledService != nil {
		scheduled := isScheduledService(a)
		if *f.HasScheduledService != scheduled {
			return false
		}
	}
	if f.MinRunwayFt > 0 && a.RunwayLength < f.MinRunwayFt {
		return false
	}
	return true
}

// ---------------------------------------------------------------------------
// Advanced Functions
// ---------------------------------------------------------------------------

// GetAutocompleteSuggestions provides autocomplete suggestions for search interfaces.
// Returns at most 10 results matching by name or IATA code.
func GetAutocompleteSuggestions(query string) ([]Airport, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	query = strings.TrimSpace(query)
	if len(query) < 1 {
		return nil, fmt.Errorf("query must be at least 1 character")
	}
	lowerQuery := strings.ToLower(query)
	var result []Airport
	for _, a := range airports {
		if strings.Contains(strings.ToLower(a.Name), lowerQuery) ||
			strings.Contains(strings.ToLower(a.IATA), lowerQuery) {
			result = append(result, a)
			if len(result) >= 10 {
				break
			}
		}
	}
	return result, nil
}

// GetAirportLinks gets external links for an airport using IATA or ICAO code.
func GetAirportLinks(code string) (*AirportLinks, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	a, err := resolveAirport(code)
	if err != nil {
		return nil, err
	}
	links := &AirportLinks{
		Wikipedia:     a.Wikipedia,
		Flightradar24: a.Flightradar24URL,
		Radarbox:      a.RadarboxURL,
		Flightaware:   a.FlightawareURL,
	}
	if a.Website != "" {
		links.Website = a.Website
	}
	return links, nil
}

// ---------------------------------------------------------------------------
// Statistical & Analytical Functions
// ---------------------------------------------------------------------------

// GetAirportStatsByCountry gets comprehensive statistics about airports in a specific country.
func GetAirportStatsByCountry(countryCode string) (*AirportStats, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	cc := strings.ToUpper(strings.TrimSpace(countryCode))
	indices, ok := countryIndex[cc]
	if !ok || len(indices) == 0 {
		return nil, fmt.Errorf("no airports found for country code: %s", countryCode)
	}

	stats := &AirportStats{
		ByType: make(map[string]int),
	}
	tzSet := make(map[string]bool)
	var totalRunway, totalElev float64
	var runwayCount, elevCount int

	for _, idx := range indices {
		a := airports[idx]
		stats.Total++
		stats.ByType[a.Type]++
		if isScheduledService(a) {
			stats.WithScheduledService++
		}
		if a.RunwayLength > 0 {
			totalRunway += float64(a.RunwayLength)
			runwayCount++
		}
		if a.ElevationFt != 0 {
			totalElev += float64(a.ElevationFt)
			elevCount++
		}
		if a.Timezone != "" {
			tzSet[a.Timezone] = true
		}
	}

	if runwayCount > 0 {
		stats.AverageRunwayLength = math.Round(totalRunway / float64(runwayCount))
	}
	if elevCount > 0 {
		stats.AverageElevation = math.Round(totalElev / float64(elevCount))
	}

	for tz := range tzSet {
		stats.Timezones = append(stats.Timezones, tz)
	}
	sort.Strings(stats.Timezones)

	return stats, nil
}

// GetAirportStatsByContinent gets comprehensive statistics about airports on a specific continent.
func GetAirportStatsByContinent(continentCode string) (*ContinentStats, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	cc := strings.ToUpper(strings.TrimSpace(continentCode))
	indices, ok := continentIdx[cc]
	if !ok || len(indices) == 0 {
		return nil, fmt.Errorf("no airports found for continent code: %s", continentCode)
	}

	stats := &ContinentStats{
		ByType:    make(map[string]int),
		ByCountry: make(map[string]int),
	}
	tzSet := make(map[string]bool)
	var totalRunway, totalElev float64
	var runwayCount, elevCount int

	for _, idx := range indices {
		a := airports[idx]
		stats.Total++
		stats.ByType[a.Type]++
		stats.ByCountry[a.CountryCode]++
		if isScheduledService(a) {
			stats.WithScheduledService++
		}
		if a.RunwayLength > 0 {
			totalRunway += float64(a.RunwayLength)
			runwayCount++
		}
		if a.ElevationFt != 0 {
			totalElev += float64(a.ElevationFt)
			elevCount++
		}
		if a.Timezone != "" {
			tzSet[a.Timezone] = true
		}
	}

	if runwayCount > 0 {
		stats.AverageRunwayLength = math.Round(totalRunway / float64(runwayCount))
	}
	if elevCount > 0 {
		stats.AverageElevation = math.Round(totalElev / float64(elevCount))
	}

	for tz := range tzSet {
		stats.Timezones = append(stats.Timezones, tz)
	}
	sort.Strings(stats.Timezones)

	return stats, nil
}

// GetLargestAirportsByContinent gets the largest airports on a continent sorted by
// runway length or elevation. The sortBy parameter can be "runway" or "elevation".
// Default limit is 10 if limit <= 0.
func GetLargestAirportsByContinent(continentCode string, limit int, sortBy string) ([]Airport, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 10
	}
	if sortBy == "" {
		sortBy = "runway"
	}

	cc := strings.ToUpper(strings.TrimSpace(continentCode))
	indices, ok := continentIdx[cc]
	if !ok || len(indices) == 0 {
		return nil, fmt.Errorf("no airports found for continent code: %s", continentCode)
	}

	result := make([]Airport, len(indices))
	for i, idx := range indices {
		result[i] = airports[idx]
	}

	switch strings.ToLower(sortBy) {
	case "elevation":
		sort.Slice(result, func(i, j int) bool {
			return result[i].ElevationFt > result[j].ElevationFt
		})
	default: // runway
		sort.Slice(result, func(i, j int) bool {
			return result[i].RunwayLength > result[j].RunwayLength
		})
	}

	if limit > len(result) {
		limit = len(result)
	}
	return result[:limit], nil
}

// ---------------------------------------------------------------------------
// Bulk Operations
// ---------------------------------------------------------------------------

// GetMultipleAirports fetches multiple airports by their IATA or ICAO codes.
// Returns nil in the result slice for codes that are not found.
func GetMultipleAirports(codes []string) ([]*Airport, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	result := make([]*Airport, len(codes))
	for i, code := range codes {
		a, err := resolveAirport(code)
		if err != nil {
			result[i] = nil
		} else {
			cp := *a
			result[i] = &cp
		}
	}
	return result, nil
}

// CalculateDistanceMatrix calculates distances between all pairs of airports in a list.
// Requires at least 2 valid codes. Returns an error if fewer than 2 codes or any code is invalid.
func CalculateDistanceMatrix(codes []string) (*DistanceMatrix, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	if len(codes) < 2 {
		return nil, fmt.Errorf("at least 2 airport codes are required for a distance matrix")
	}

	resolvedAirports := make([]Airport, len(codes))
	matrixAirports := make([]DistanceMatrixAirport, len(codes))
	resolvedCodes := make([]string, len(codes))

	for i, code := range codes {
		a, err := resolveAirport(code)
		if err != nil {
			return nil, fmt.Errorf("invalid airport code: %s", code)
		}
		resolvedAirports[i] = *a
		resolvedCodes[i] = strings.ToUpper(strings.TrimSpace(code))
		matrixAirports[i] = DistanceMatrixAirport{
			Code: resolvedCodes[i],
			Name: a.Name,
			IATA: a.IATA,
			ICAO: a.ICAO,
		}
	}

	distances := make(map[string]map[string]float64)
	for i, c1 := range resolvedCodes {
		distances[c1] = make(map[string]float64)
		for j, c2 := range resolvedCodes {
			if i == j {
				distances[c1][c2] = 0
			} else {
				d := haversineDistance(
					resolvedAirports[i].Latitude, resolvedAirports[i].Longitude,
					resolvedAirports[j].Latitude, resolvedAirports[j].Longitude,
				)
				distances[c1][c2] = math.Round(d)
			}
		}
	}

	return &DistanceMatrix{
		Airports:  matrixAirports,
		Distances: distances,
	}, nil
}

// ---------------------------------------------------------------------------
// Validation & Utilities
// ---------------------------------------------------------------------------

var (
	iataRegexp = regexp.MustCompile(`^[A-Z]{3}$`)
	icaoRegexp = regexp.MustCompile(`^[A-Z0-9]{4}$`)
)

// ValidateIataCode validates if an IATA code exists in the database.
// Returns false for codes that don't match the 3-letter uppercase format.
func ValidateIataCode(code string) (bool, error) {
	if err := ensureLoaded(); err != nil {
		return false, err
	}
	c := strings.TrimSpace(code)
	if !iataRegexp.MatchString(c) {
		return false, nil
	}
	_, ok := iataIndex[c]
	return ok, nil
}

// ValidateIcaoCode validates if an ICAO code exists in the database.
// Returns false for codes that don't match the 4-character uppercase format.
func ValidateIcaoCode(code string) (bool, error) {
	if err := ensureLoaded(); err != nil {
		return false, err
	}
	c := strings.TrimSpace(code)
	if !icaoRegexp.MatchString(c) {
		return false, nil
	}
	_, ok := icaoIndex[c]
	return ok, nil
}

// GetAirportCount gets the count of airports matching the given filters.
// Pass nil or zero-value filter to get the total count.
func GetAirportCount(filter *FindAirportsFilter) (int, error) {
	if err := ensureLoaded(); err != nil {
		return 0, err
	}
	if filter == nil {
		return len(airports), nil
	}
	count := 0
	for _, a := range airports {
		if matchesFilter(a, filter) {
			count++
		}
	}
	return count, nil
}

// IsAirportOperational checks if an airport has scheduled commercial service.
// Accepts IATA or ICAO codes.
func IsAirportOperational(code string) (bool, error) {
	if err := ensureLoaded(); err != nil {
		return false, err
	}
	a, err := resolveAirport(code)
	if err != nil {
		return false, fmt.Errorf("invalid airport code: %s", code)
	}
	return isScheduledService(*a), nil
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// resolveAirport looks up an airport by IATA code first, then ICAO code.
func resolveAirport(code string) (*Airport, error) {
	c := strings.ToUpper(strings.TrimSpace(code))

	// Try IATA first
	if indices, ok := iataIndex[c]; ok && len(indices) > 0 {
		return &airports[indices[0]], nil
	}
	// Try ICAO
	if indices, ok := icaoIndex[c]; ok && len(indices) > 0 {
		return &airports[indices[0]], nil
	}
	return nil, fmt.Errorf("no airport found for code: %s", code)
}
