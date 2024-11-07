package main

import (
	"math"
	"testing"
)

// multiply is a simple function to be tested.
func multiply(a, b int) int {
	return a * b
}

func TestMultiplyBasicCases(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"Positive numbers", 3, 4, 12},
		{"One positive, one negative", -1, 5, -5},
		{"Two negative numbers", -3, -4, 12},
		{"One zero value", 0, 5, 0},
		{"Both zero values", 0, 0, 0},
		{"Large numbers", 100000, 100000, 10000000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := multiply(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("multiply(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMultiplyEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		a, b        int
		expectError bool
	}{
		// Assuming we modify multiply to return an error on overflow
		{"Overflow edge case", math.MaxInt32, 2, true},
		{"Underflow edge case", math.MinInt32, 2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := multiply(tt.a, tt.b)
			if (err != nil) != tt.expectError {
				t.Errorf("multiply(%d, %d) expected error: %v, got: %v", tt.a, tt.b, tt.expectError, err != nil)
			}
		})
	}
}

func TestMultiplyOverflow(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"MaxInt overflow positive", math.MaxInt32, 2, -2}, // Overflow, but specific behavior depends on implementation
		{"MaxInt overflow negative", math.MinInt32, 2, 0},  // Overflow, behavior example
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := multiply(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Potential overflow detected in multiply(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}