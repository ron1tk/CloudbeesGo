package main

import (
    "testing"
)

// TestMultiply_Success tests the multiply function with normal input values.
func TestMultiply_Success(t *testing.T) {
    tests := []struct {
        name     string
        a        int
        b        int
        expected int
    }{
        {"Positive numbers", 3, 4, 12},
        {"Negative and positive number", -2, 5, -10},
        {"Two negative numbers", -3, -3, 9},
        {"Zero and positive number", 0, 5, 0},
        {"Positive number and zero", 5, 0, 0},
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

// TestMultiply_EdgeCases tests the multiply function with edge case input values.
func TestMultiply_EdgeCases(t *testing.T) {
    tests := []struct {
        name     string
        a        int
        b        int
        expected int
    }{
        {"MaxInt and 1", 1<<31 - 1, 1, 1<<31 - 1},
        {"1 and MaxInt", 1, 1<<31 - 1, 1<<31 - 1},
        {"MinInt and 1", -1 << 31, 1, -1 << 31},
        {"1 and MinInt", 1, -1 << 31, -1 << 31},
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

// TestMultiply_ErrorCases could be implemented if there were error-producing inputs to consider.
// Since the multiply function does not return an error, this is not applicable in this scenario.