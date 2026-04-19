package gophonenumbers

import (
	"testing"
)

func TestAreaCodeGeoCSV(t *testing.T) {
	data := AreaCodeGeoCSV()
	if len(data) == 0 {
		t.Error("embedded CSV data is empty")
	}
}

func TestAreaCodeToGeoReadData(t *testing.T) {
	a2g := NewAreaCodeToGeo()
	if err := a2g.ReadData(); err != nil {
		t.Fatalf("ReadData() error: %v", err)
	}

	// Verify we have area codes loaded
	acs := a2g.AreaCodes()
	if len(acs) == 0 {
		t.Error("no area codes loaded")
	}

	// Check a well-known area code exists (201 is New Jersey)
	if _, ok := a2g.AreaCodeInfos[201]; !ok {
		t.Error("expected area code 201 to exist")
	}
}

func TestAreaCodeToGeoDistanceMatrix(t *testing.T) {
	a2g := NewAreaCodeToGeo()
	if err := a2g.ReadData(); err != nil {
		t.Fatalf("ReadData() error: %v", err)
	}

	// Distance matrix should be populated
	if len(a2g.DistanceMatrix) == 0 {
		t.Error("distance matrix is empty")
	}

	// Test distance between two area codes
	dist, err := a2g.GcdAreaCodes(201, 212) // NJ to NYC
	if err != nil {
		t.Fatalf("GcdAreaCodes() error: %v", err)
	}
	if dist <= 0 {
		t.Errorf("expected positive distance, got %v", dist)
	}
}

func TestFakeNumberGenerator(t *testing.T) {
	a2g := NewAreaCodeToGeo()
	if err := a2g.ReadData(); err != nil {
		t.Fatalf("ReadData() error: %v", err)
	}

	acs := a2g.AreaCodes()
	if len(acs) == 0 {
		t.Fatal("no area codes available")
	}

	fng := NewFakeNumberGenerator(acs)

	// Generate a fake number
	num, err := fng.RandomLocalNumberUS()
	if err != nil {
		t.Fatalf("RandomLocalNumberUS() error: %v", err)
	}

	// US numbers should be 11 digits (1 + 10)
	if num < 10000000000 || num > 19999999999 {
		t.Errorf("generated number %d is not a valid US number", num)
	}
}
