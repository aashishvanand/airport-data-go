package airportdata

import (
	"math"
	"testing"
)

// ---------------------------------------------------------------------------
// Core Search Functions
// ---------------------------------------------------------------------------

func TestGetAirportByIata(t *testing.T) {
	t.Run("should retrieve airport data for a valid IATA code", func(t *testing.T) {
		airports, err := GetAirportByIata("LHR")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airports) == 0 {
			t.Fatal("expected at least one airport")
		}
		if airports[0].IATA != "LHR" {
			t.Errorf("expected IATA 'LHR', got '%s'", airports[0].IATA)
		}
		if !contains(airports[0].Name, "Heathrow") {
			t.Errorf("expected name to contain 'Heathrow', got '%s'", airports[0].Name)
		}
	})

	t.Run("should return error for invalid IATA code", func(t *testing.T) {
		_, err := GetAirportByIata("XYZ")
		if err == nil {
			t.Fatal("expected error for invalid IATA code")
		}
	})
}

func TestGetAirportByIcao(t *testing.T) {
	t.Run("should retrieve airport data for a valid ICAO code", func(t *testing.T) {
		airports, err := GetAirportByIcao("EGLL")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airports) == 0 {
			t.Fatal("expected at least one airport")
		}
		if airports[0].ICAO != "EGLL" {
			t.Errorf("expected ICAO 'EGLL', got '%s'", airports[0].ICAO)
		}
		if !contains(airports[0].Name, "Heathrow") {
			t.Errorf("expected name to contain 'Heathrow', got '%s'", airports[0].Name)
		}
	})
}

func TestSearchByName(t *testing.T) {
	t.Run("should find airports matching name query", func(t *testing.T) {
		results, err := SearchByName("Singapore")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(results) == 0 {
			t.Fatal("expected at least one result")
		}
	})

	t.Run("should reject short queries", func(t *testing.T) {
		_, err := SearchByName("S")
		if err == nil {
			t.Fatal("expected error for short query")
		}
	})
}

// ---------------------------------------------------------------------------
// Filtering Functions
// ---------------------------------------------------------------------------

func TestGetAirportByCountryCode(t *testing.T) {
	t.Run("should retrieve all airports for a given country code", func(t *testing.T) {
		airports, err := GetAirportByCountryCode("US")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airports) <= 100 {
			t.Errorf("expected more than 100 US airports, got %d", len(airports))
		}
		if airports[0].CountryCode != "US" {
			t.Errorf("expected country code 'US', got '%s'", airports[0].CountryCode)
		}
	})
}

func TestGetAirportByContinent(t *testing.T) {
	t.Run("should retrieve all airports for a given continent code", func(t *testing.T) {
		airports, err := GetAirportByContinent("EU")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airports) <= 100 {
			t.Errorf("expected more than 100 EU airports, got %d", len(airports))
		}
		for _, a := range airports {
			if a.Continent != "EU" {
				t.Errorf("expected continent 'EU', got '%s'", a.Continent)
				break
			}
		}
	})
}

func TestGetAirportsByType(t *testing.T) {
	t.Run("should retrieve all large airports", func(t *testing.T) {
		airports, err := GetAirportsByType("large_airport")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airports) <= 10 {
			t.Errorf("expected more than 10 large airports, got %d", len(airports))
		}
		for _, a := range airports {
			if a.Type != "large_airport" {
				t.Errorf("expected type 'large_airport', got '%s'", a.Type)
				break
			}
		}
	})

	t.Run("should retrieve all medium airports", func(t *testing.T) {
		airports, err := GetAirportsByType("medium_airport")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airports) <= 10 {
			t.Errorf("expected more than 10 medium airports, got %d", len(airports))
		}
		for _, a := range airports {
			if a.Type != "medium_airport" {
				t.Errorf("expected type 'medium_airport', got '%s'", a.Type)
				break
			}
		}
	})

	t.Run("should retrieve all airports when searching for 'airport'", func(t *testing.T) {
		airports, err := GetAirportsByType("airport")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airports) <= 50 {
			t.Errorf("expected more than 50 airports, got %d", len(airports))
		}
		for _, a := range airports {
			if !contains(a.Type, "airport") {
				t.Errorf("expected type to contain 'airport', got '%s'", a.Type)
				break
			}
		}
	})

	t.Run("should handle different airport types", func(t *testing.T) {
		heliports, err := GetAirportsByType("heliport")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, a := range heliports {
			if a.Type != "heliport" {
				t.Errorf("expected type 'heliport', got '%s'", a.Type)
				break
			}
		}

		seaplaneBases, err := GetAirportsByType("seaplane_base")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, a := range seaplaneBases {
			if a.Type != "seaplane_base" {
				t.Errorf("expected type 'seaplane_base', got '%s'", a.Type)
				break
			}
		}
	})

	t.Run("should handle case-insensitive searches", func(t *testing.T) {
		upper, err := GetAirportsByType("LARGE_AIRPORT")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		lower, err := GetAirportsByType("large_airport")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(upper) != len(lower) {
			t.Errorf("case-insensitive search returned different counts: %d vs %d", len(upper), len(lower))
		}
		if len(upper) == 0 {
			t.Error("expected results for large_airport")
		}
	})

	t.Run("should return empty array for non-existent type", func(t *testing.T) {
		airports, err := GetAirportsByType("nonexistent_type")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airports) != 0 {
			t.Errorf("expected 0 airports for non-existent type, got %d", len(airports))
		}
	})
}

func TestGetAirportsByTimezone(t *testing.T) {
	t.Run("should find all airports within a specific timezone", func(t *testing.T) {
		airports, err := GetAirportsByTimezone("Europe/London")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airports) <= 10 {
			t.Errorf("expected more than 10 airports in Europe/London timezone, got %d", len(airports))
		}
		for _, a := range airports {
			if a.Timezone != "Europe/London" {
				t.Errorf("expected timezone 'Europe/London', got '%s'", a.Timezone)
				break
			}
		}
	})
}

func TestFindAirports(t *testing.T) {
	t.Run("should find airports with multiple matching criteria", func(t *testing.T) {
		airports, err := FindAirports(FindAirportsFilter{
			CountryCode: "GB",
			Type:        "airport",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, a := range airports {
			if a.CountryCode != "GB" {
				t.Errorf("expected country code 'GB', got '%s'", a.CountryCode)
				break
			}
			if !contains(a.Type, "airport") {
				t.Errorf("expected type to contain 'airport', got '%s'", a.Type)
				break
			}
		}
	})

	t.Run("should filter by scheduled service availability", func(t *testing.T) {
		trueVal := true
		falseVal := false

		withService, err := FindAirports(FindAirportsFilter{HasScheduledService: &trueVal})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		withoutService, err := FindAirports(FindAirportsFilter{HasScheduledService: &falseVal})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(withService)+len(withoutService) == 0 {
			t.Error("expected at least some airports with or without scheduled service")
		}

		for _, a := range withService {
			if !isScheduledService(a) {
				t.Errorf("expected scheduled service for airport %s", a.IATA)
				break
			}
		}

		for _, a := range withoutService {
			if isScheduledService(a) {
				t.Errorf("expected no scheduled service for airport %s", a.IATA)
				break
			}
		}
	})
}

// ---------------------------------------------------------------------------
// Geographic Functions
// ---------------------------------------------------------------------------

func TestFindNearbyAirports(t *testing.T) {
	t.Run("should find airports within a given radius", func(t *testing.T) {
		lat := 51.5074
		lon := -0.1278
		airports, err := FindNearbyAirports(lat, lon, 50)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airports) < 1 {
			t.Fatal("expected at least 1 airport near London")
		}
		found := false
		for _, a := range airports {
			if a.IATA == "LHR" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected LHR to be in nearby airports")
		}
	})
}

func TestCalculateDistance(t *testing.T) {
	t.Run("should calculate the distance between two airports using IATA codes", func(t *testing.T) {
		distance, err := CalculateDistance("LHR", "JFK")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Approx 5541 km
		if math.Abs(distance-5541) > 100 {
			t.Errorf("expected distance ~5541 km, got %.0f", distance)
		}
	})
}

// ---------------------------------------------------------------------------
// Advanced Functions
// ---------------------------------------------------------------------------

func TestGetAutocompleteSuggestions(t *testing.T) {
	t.Run("should return suggestions based on airport name", func(t *testing.T) {
		suggestions, err := GetAutocompleteSuggestions("London")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(suggestions) == 0 {
			t.Fatal("expected at least one suggestion")
		}
		if len(suggestions) > 10 {
			t.Errorf("expected at most 10 suggestions, got %d", len(suggestions))
		}
		found := false
		for _, a := range suggestions {
			if a.IATA == "LHR" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected LHR to be in suggestions")
		}
	})
}

func TestGetAirportLinks(t *testing.T) {
	t.Run("should retrieve a map of all available external links", func(t *testing.T) {
		links, err := GetAirportLinks("LHR")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !contains(links.Wikipedia, "Heathrow_Airport") {
			t.Errorf("expected Wikipedia link to contain 'Heathrow_Airport', got '%s'", links.Wikipedia)
		}
		if links.Website == "" {
			t.Error("expected website link to exist")
		}
	})

	t.Run("should handle airports with missing links gracefully", func(t *testing.T) {
		links, err := GetAirportLinks("HND")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !contains(links.Wikipedia, "Tokyo_International_Airport") {
			t.Errorf("expected Wikipedia link to contain 'Tokyo_International_Airport', got '%s'", links.Wikipedia)
		}
		if links.Website == "" {
			t.Error("expected website link to exist")
		}
	})
}

// ---------------------------------------------------------------------------
// Statistical & Analytical Functions
// ---------------------------------------------------------------------------

func TestGetAirportStatsByCountry(t *testing.T) {
	t.Run("should return comprehensive statistics for a country", func(t *testing.T) {
		stats, err := GetAirportStatsByCountry("SG")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if stats.Total == 0 {
			t.Error("expected total > 0")
		}
		if stats.ByType == nil {
			t.Error("expected byType to be non-nil")
		}
		if stats.Timezones == nil {
			t.Error("expected timezones to be non-nil")
		}
	})

	t.Run("should calculate correct statistics for US airports", func(t *testing.T) {
		stats, err := GetAirportStatsByCountry("US")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if stats.Total <= 1000 {
			t.Errorf("expected total > 1000 for US, got %d", stats.Total)
		}
		if stats.ByType["large_airport"] == 0 {
			t.Error("expected large_airport count > 0 for US")
		}
	})

	t.Run("should return error for invalid country code", func(t *testing.T) {
		_, err := GetAirportStatsByCountry("XYZ")
		if err == nil {
			t.Fatal("expected error for invalid country code")
		}
	})
}

func TestGetAirportStatsByContinent(t *testing.T) {
	t.Run("should return comprehensive statistics for a continent", func(t *testing.T) {
		stats, err := GetAirportStatsByContinent("AS")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if stats.Total <= 100 {
			t.Errorf("expected total > 100 for Asia, got %d", stats.Total)
		}
		if len(stats.ByCountry) <= 10 {
			t.Errorf("expected more than 10 countries in Asia, got %d", len(stats.ByCountry))
		}
	})

	t.Run("should include country breakdown", func(t *testing.T) {
		stats, err := GetAirportStatsByContinent("EU")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if stats.ByCountry["GB"] == 0 {
			t.Error("expected GB in country breakdown")
		}
		if stats.ByCountry["FR"] == 0 {
			t.Error("expected FR in country breakdown")
		}
		if stats.ByCountry["DE"] == 0 {
			t.Error("expected DE in country breakdown")
		}
	})
}

func TestGetLargestAirportsByContinent(t *testing.T) {
	t.Run("should return top airports by runway length", func(t *testing.T) {
		airports, err := GetLargestAirportsByContinent("AS", 5, "runway")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airports) == 0 {
			t.Fatal("expected at least one airport")
		}
		if len(airports) > 5 {
			t.Errorf("expected at most 5 airports, got %d", len(airports))
		}
		// Check sorted by runway length descending
		for i := 0; i < len(airports)-1; i++ {
			if airports[i].RunwayLength < airports[i+1].RunwayLength {
				t.Errorf("airports not sorted by runway length: %d < %d",
					airports[i].RunwayLength, airports[i+1].RunwayLength)
			}
		}
	})

	t.Run("should return top airports by elevation", func(t *testing.T) {
		airports, err := GetLargestAirportsByContinent("SA", 5, "elevation")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airports) > 5 {
			t.Errorf("expected at most 5 airports, got %d", len(airports))
		}
		// Check sorted by elevation descending
		for i := 0; i < len(airports)-1; i++ {
			if airports[i].ElevationFt < airports[i+1].ElevationFt {
				t.Errorf("airports not sorted by elevation: %d < %d",
					airports[i].ElevationFt, airports[i+1].ElevationFt)
			}
		}
	})

	t.Run("should respect the limit parameter", func(t *testing.T) {
		airports, err := GetLargestAirportsByContinent("EU", 3, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airports) > 3 {
			t.Errorf("expected at most 3 airports, got %d", len(airports))
		}
	})
}

// ---------------------------------------------------------------------------
// Bulk Operations
// ---------------------------------------------------------------------------

func TestGetMultipleAirports(t *testing.T) {
	t.Run("should fetch multiple airports by IATA codes", func(t *testing.T) {
		airports, err := GetMultipleAirports([]string{"SIN", "LHR", "JFK"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airports) != 3 {
			t.Fatalf("expected 3 airports, got %d", len(airports))
		}
		if airports[0].IATA != "SIN" {
			t.Errorf("expected IATA 'SIN', got '%s'", airports[0].IATA)
		}
		if airports[1].IATA != "LHR" {
			t.Errorf("expected IATA 'LHR', got '%s'", airports[1].IATA)
		}
		if airports[2].IATA != "JFK" {
			t.Errorf("expected IATA 'JFK', got '%s'", airports[2].IATA)
		}
	})

	t.Run("should handle mix of IATA and ICAO codes", func(t *testing.T) {
		airports, err := GetMultipleAirports([]string{"SIN", "EGLL", "JFK"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airports) != 3 {
			t.Fatalf("expected 3 airports, got %d", len(airports))
		}
		for _, a := range airports {
			if a == nil {
				t.Error("expected all airports to be non-nil")
			}
		}
	})

	t.Run("should return nil for invalid codes", func(t *testing.T) {
		airports, err := GetMultipleAirports([]string{"SIN", "INVALID", "LHR"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airports) != 3 {
			t.Fatalf("expected 3 results, got %d", len(airports))
		}
		if airports[0] == nil {
			t.Error("expected first airport to be non-nil")
		}
		if airports[1] != nil {
			t.Error("expected second airport to be nil for invalid code")
		}
		if airports[2] == nil {
			t.Error("expected third airport to be non-nil")
		}
	})

	t.Run("should handle empty array", func(t *testing.T) {
		airports, err := GetMultipleAirports([]string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(airports) != 0 {
			t.Errorf("expected 0 airports, got %d", len(airports))
		}
	})
}

func TestCalculateDistanceMatrix(t *testing.T) {
	t.Run("should calculate distance matrix for multiple airports", func(t *testing.T) {
		matrix, err := CalculateDistanceMatrix([]string{"SIN", "LHR", "JFK"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(matrix.Airports) != 3 {
			t.Fatalf("expected 3 airports in matrix, got %d", len(matrix.Airports))
		}

		// Check diagonal is zero
		if matrix.Distances["SIN"]["SIN"] != 0 {
			t.Error("expected diagonal to be 0")
		}
		if matrix.Distances["LHR"]["LHR"] != 0 {
			t.Error("expected diagonal to be 0")
		}
		if matrix.Distances["JFK"]["JFK"] != 0 {
			t.Error("expected diagonal to be 0")
		}

		// Check symmetry
		if matrix.Distances["SIN"]["LHR"] != matrix.Distances["LHR"]["SIN"] {
			t.Error("expected SIN-LHR distance to be symmetric")
		}
		if matrix.Distances["SIN"]["JFK"] != matrix.Distances["JFK"]["SIN"] {
			t.Error("expected SIN-JFK distance to be symmetric")
		}

		// Check reasonable distances
		if matrix.Distances["SIN"]["LHR"] <= 5000 {
			t.Errorf("expected SIN-LHR > 5000 km, got %.0f", matrix.Distances["SIN"]["LHR"])
		}
		if matrix.Distances["LHR"]["JFK"] <= 3000 {
			t.Errorf("expected LHR-JFK > 3000 km, got %.0f", matrix.Distances["LHR"]["JFK"])
		}
	})

	t.Run("should return error for less than 2 airports", func(t *testing.T) {
		_, err := CalculateDistanceMatrix([]string{"SIN"})
		if err == nil {
			t.Fatal("expected error for less than 2 airports")
		}
	})

	t.Run("should return error for invalid codes", func(t *testing.T) {
		_, err := CalculateDistanceMatrix([]string{"SIN", "INVALID"})
		if err == nil {
			t.Fatal("expected error for invalid codes")
		}
	})
}

func TestFindNearestAirport(t *testing.T) {
	t.Run("should find nearest airport to coordinates", func(t *testing.T) {
		nearest, err := FindNearestAirport(1.35019, 103.994003, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if nearest.IATA != "SIN" {
			t.Errorf("expected nearest airport to be SIN, got '%s'", nearest.IATA)
		}
		if nearest.Distance >= 2 {
			t.Errorf("expected distance < 2 km, got %.2f", nearest.Distance)
		}
	})

	t.Run("should find nearest airport with type filter", func(t *testing.T) {
		nearest, err := FindNearestAirport(51.5074, -0.1278, &FindAirportsFilter{
			Type: "large_airport",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if nearest == nil {
			t.Fatal("expected a result")
		}
		if nearest.Type != "large_airport" {
			t.Errorf("expected type 'large_airport', got '%s'", nearest.Type)
		}
	})

	t.Run("should find nearest airport with type and country filters", func(t *testing.T) {
		nearest, err := FindNearestAirport(40.7128, -74.0060, &FindAirportsFilter{
			Type:        "large_airport",
			CountryCode: "US",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if nearest == nil {
			t.Fatal("expected a result")
		}
		if nearest.Type != "large_airport" {
			t.Errorf("expected type 'large_airport', got '%s'", nearest.Type)
		}
		if nearest.CountryCode != "US" {
			t.Errorf("expected country 'US', got '%s'", nearest.CountryCode)
		}
	})
}

// ---------------------------------------------------------------------------
// Validation & Utilities
// ---------------------------------------------------------------------------

func TestValidateIataCode(t *testing.T) {
	t.Run("should return true for valid IATA codes", func(t *testing.T) {
		for _, code := range []string{"SIN", "LHR", "JFK"} {
			valid, err := ValidateIataCode(code)
			if err != nil {
				t.Fatalf("unexpected error for %s: %v", code, err)
			}
			if !valid {
				t.Errorf("expected %s to be valid", code)
			}
		}
	})

	t.Run("should return false for invalid IATA codes", func(t *testing.T) {
		for _, code := range []string{"XYZ", "ZZZ"} {
			valid, err := ValidateIataCode(code)
			if err != nil {
				t.Fatalf("unexpected error for %s: %v", code, err)
			}
			if valid {
				t.Errorf("expected %s to be invalid", code)
			}
		}
	})

	t.Run("should return false for incorrect format", func(t *testing.T) {
		for _, code := range []string{"ABCD", "AB", "abc", ""} {
			valid, err := ValidateIataCode(code)
			if err != nil {
				t.Fatalf("unexpected error for '%s': %v", code, err)
			}
			if valid {
				t.Errorf("expected '%s' to be invalid", code)
			}
		}
	})
}

func TestValidateIcaoCode(t *testing.T) {
	t.Run("should return true for valid ICAO codes", func(t *testing.T) {
		for _, code := range []string{"WSSS", "EGLL", "KJFK"} {
			valid, err := ValidateIcaoCode(code)
			if err != nil {
				t.Fatalf("unexpected error for %s: %v", code, err)
			}
			if !valid {
				t.Errorf("expected %s to be valid", code)
			}
		}
	})

	t.Run("should return false for invalid ICAO codes", func(t *testing.T) {
		for _, code := range []string{"XXXX", "ZZZ0"} {
			valid, err := ValidateIcaoCode(code)
			if err != nil {
				t.Fatalf("unexpected error for %s: %v", code, err)
			}
			if valid {
				t.Errorf("expected %s to be invalid", code)
			}
		}
	})

	t.Run("should return false for incorrect format", func(t *testing.T) {
		for _, code := range []string{"ABC", "ABCDE", "abcd", ""} {
			valid, err := ValidateIcaoCode(code)
			if err != nil {
				t.Fatalf("unexpected error for '%s': %v", code, err)
			}
			if valid {
				t.Errorf("expected '%s' to be invalid", code)
			}
		}
	})
}

func TestGetAirportCount(t *testing.T) {
	t.Run("should return total count of all airports", func(t *testing.T) {
		count, err := GetAirportCount(nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count <= 5000 {
			t.Errorf("expected count > 5000, got %d", count)
		}
	})

	t.Run("should return count with type filter", func(t *testing.T) {
		largeCount, err := GetAirportCount(&FindAirportsFilter{Type: "large_airport"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		totalCount, err := GetAirportCount(nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if largeCount == 0 {
			t.Error("expected large airport count > 0")
		}
		if largeCount >= totalCount {
			t.Error("expected large airport count < total count")
		}
	})

	t.Run("should return count with country filter", func(t *testing.T) {
		usCount, err := GetAirportCount(&FindAirportsFilter{CountryCode: "US"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if usCount <= 1000 {
			t.Errorf("expected US count > 1000, got %d", usCount)
		}
	})

	t.Run("should return count with multiple filters", func(t *testing.T) {
		count, err := GetAirportCount(&FindAirportsFilter{
			CountryCode: "US",
			Type:        "large_airport",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count == 0 {
			t.Error("expected count > 0")
		}
		if count >= 200 {
			t.Errorf("expected count < 200, got %d", count)
		}
	})
}

func TestIsAirportOperational(t *testing.T) {
	t.Run("should return true for operational airports", func(t *testing.T) {
		for _, code := range []string{"SIN", "LHR", "JFK"} {
			op, err := IsAirportOperational(code)
			if err != nil {
				t.Fatalf("unexpected error for %s: %v", code, err)
			}
			if !op {
				t.Errorf("expected %s to be operational", code)
			}
		}
	})

	t.Run("should work with both IATA and ICAO codes", func(t *testing.T) {
		op1, err := IsAirportOperational("SIN")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		op2, err := IsAirportOperational("WSSS")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !op1 || !op2 {
			t.Error("expected both SIN and WSSS to be operational")
		}
	})

	t.Run("should return error for invalid airport code", func(t *testing.T) {
		_, err := IsAirportOperational("INVALID")
		if err == nil {
			t.Fatal("expected error for invalid airport code")
		}
	})
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func contains(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && indexOfLower(s, substr) >= 0)
}

func indexOfLower(s, substr string) int {
	sl := toLower(s)
	sl2 := toLower(substr)
	for i := 0; i <= len(sl)-len(sl2); i++ {
		if sl[i:i+len(sl2)] == sl2 {
			return i
		}
	}
	return -1
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		b[i] = c
	}
	return string(b)
}
