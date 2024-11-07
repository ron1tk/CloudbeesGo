package main

import (
    "testing"
)

// TestMultiply tests the multiply function with various inputs.
func TestMultiply(t *testing.T) {
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

// TestMultiplyOverflow tests the multiply function for potential integer overflow.
func TestMultiplyOverflow(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int // Note: This simplistic approach doesn't properly handle overflow expectation.
    }{
        {"MaxInt overflow positive", 2147483647, 2, -2},
        {"MaxInt overflow negative", -2147483648, 2, 0}, // Overflow behavior in Go results in wrapping around
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