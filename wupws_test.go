package wupws

import "testing"

func TestHpaToInhg(t *testing.T) {
	expectedStr := "29.92"
	result := HpaToInhg(1013.25)
	if result != expectedStr {
		t.Fatalf("Expected %s, got %s", expectedStr, result)
	}
}

func TestDewpointCelcius(t *testing.T) {
	expectedStr := "32"
	result := DewpointCelcius(32, 200)
	if result != expectedStr {
		t.Fatalf("Expected %s, got %s", expectedStr, result)
	}
}

func TestCelciusToFahrenheit(t *testing.T) {
	expectedStr := "32"
	result := CelciusToFahrenheit(0)
	if result != expectedStr {
		t.Fatalf("Expected %s, got %s", expectedStr, result)
	}
}
