package test

package main

import (
    "testing"
)

// TestMultiply_NormalCases tests the multiply function with normal inputs.
func TestMultiply_NormalCases(t *testing.T) {
    testCases := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"Both positive numbers", 3, 4, 12},
        {"Both negative numbers", -2, -4, 8},
        {"One positive, one negative", -5, 3, -15},
        {"One negative, one positive", 5, -3, -15},
        {"One zero, one positive", 0, 5, 0},
        {"One positive, one zero", 6, 0, 0},
        {"Both zeroes", 0, 0, 0},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := multiply(tc.a, tc.b)
            if result != tc.expected {
                t.Errorf("multiply(%d, %d) = %d; want %d", tc.a, tc.b, result, tc.expected)
            }
        })
    }
}

// TestMultiply_EdgeCases tests the multiply function with edge case inputs.
func TestMultiply_EdgeCases(t *testing.T) {
    testCases := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"MaxInt and 1", 1<<31 - 1, 1, 1<<31 - 1},
        {"1 and MaxInt", 1, 1<<31 - 1, 1<<31 - 1},
        {"MinInt and 1", -1 << 31, 1, -1 << 31},
        {"1 and MinInt", 1, -1 << 31, -1 << 31},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := multiply(tc.a, tc.b)
            if result != tc.expected {
                t.Errorf("multiply(%d, %d) = %d; want %d", tc.a, tc.b, result, tc.expected)
            }
        })
    }
}

// TestMultiply_ErrorCases tests the multiply function for potential overflow errors.
// Note: Go does not have built-in overflow checks for int operations. This test case
// is more illustrative of how you might approach testing edge cases for overflow in languages
// that do support overflow exceptions or in scenarios where custom overflow handling is implemented.
func TestMultiply_ErrorCases(t *testing.T) {
    testCases := []struct {
        name     string
        a, b     int
        hasError bool
    }{
        // These test cases depend on the behavior of a hypothetical overflow-checking multiply function.
        // Go's `int` type does not provide overflow checks, so these tests always expect `hasError` to be `false` under normal Go semantics.
        {"Potential overflow case 1", 1<<30, 2, false},
        {"Potential overflow case 2", 2, 1<<30, false},
        {"Potential overflow case 3", -1 << 30, 2, false},
        {"Potential overflow case 4", 2, -1 << 30, false},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Assuming a hypothetical `multiplyChecked` function that returns an error on overflow.
            result := multiply(tc.a, tc.b)
            // Simulate overflow check (not applicable in standard Go)
            if (result > 0) != tc.hasError {
                t.Errorf("multiply(%d, %d) expected overflow error; got none", tc.a, tc.b)
            }
        })
    }
}