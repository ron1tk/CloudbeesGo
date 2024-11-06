package main

import (
    "testing"
)

// TestSum_EmptyInput tests the Sum function with no input.
func TestSum_EmptyInput(t *testing.T) {
    got := Sum(0, 0)
    want := 0
    if got != want {
        t.Errorf("Sum(0, 0) = %d; want %d", got, want)
    }
}

// TestMinus_EmptyInput tests the Minus function with no input.
func TestMinus_EmptyInput(t *testing.T) {
    got := Minus(0, 0)
    want := 0
    if got != want {
        t.Errorf("Minus(0, 0) = %d; want %d", got, want)
    }
}

// TestSum_InvalidInput tests the Sum function with invalid (non-integer) inputs.
// This test is not applicable as the function signature only allows integers. It's included for completeness.
func TestSum_InvalidInput(t *testing.T) {
    // This space intentionally left blank. In Go, the type system prevents non-integer inputs.
}

// TestMinus_InvalidInput tests the Minus function with invalid (non-integer) inputs.
// This test is not applicable as the function signature only allows integers. It's included for completeness.
func TestMinus_InvalidInput(t *testing.T) {
    // This space intentionally left blank. In Go, the type system prevents non-integer inputs.
}

// TestSum_ExtremeValues tests the Sum function with extreme values of int.
func TestSum_ExtremeValues(t *testing.T) {
    testCases := []struct {
        name string
        a    int
        b    int
        want int
    }{
        {"MaxInt and MaxInt", 2147483647, 2147483647, -2}, // Overflow
        {"MinInt and MinInt", -2147483648, -2147483648, 0}, // Underflow
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            got := Sum(tc.a, tc.b)
            if got != tc.want {
                t.Errorf("Potential overflow/underflow: Sum(%d, %d) = %d; want %d", tc.a, tc.b, got, tc.want)
            }
        })
    }
}

// TestMinus_ExtremeValues tests the Minus function with extreme values of int.
func TestMinus_ExtremeValues(t *testing.T) {
    testCases := []struct {
        name string
        a    int
        b    int
        want int
    }{
        {"MaxInt and MinInt", 2147483647, -2147483648, -1}, // Underflow
        {"MinInt and MaxInt", -2147483648, 2147483647, 1},  // Overflow
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            got := Minus(tc.a, tc.b)
            if got != tc.want {
                t.Errorf("Potential overflow/underflow: Minus(%d, %d) = %d; want %d", tc.a, tc.b, got, tc.want)
            }
        })
    }
}