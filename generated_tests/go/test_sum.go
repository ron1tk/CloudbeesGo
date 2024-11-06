

package main

import (
    "testing"
)

// TestSum_ValidInputs tests the Sum function with valid inputs.
func TestSum_ValidInputs(t *testing.T) {
    testCases := []struct {
        name string
        a    int
        b    int
        want int
    }{
        {"Positive numbers", 5, 3, 8},
        {"Negative numbers", -2, -4, -6},
        {"Mixed numbers", -1, 2, 1},
        {"Zero and positive", 0, 5, 5},
        {"Positive and zero", 10, 0, 10},
        {"Zero and negative", 0, -3, -3},
        {"Negative and zero", -7, 0, -7},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            got := Sum(tc.a, tc.b)
            if got != tc.want {
                t.Errorf("Sum(%d, %d) = %d; want %d", tc.a, tc.b, got, tc.want)
            }
        })
    }
}

// TestMinus_ValidInputs tests the Minus function with valid inputs.
func TestMinus_ValidInputs(t *testing.T) {
    testCases := []struct {
        name string
        a    int
        b    int
        want int
    }{
        {"Positive numbers", 5, 3, 2},
        {"Negative numbers", -2, -4, 2},
        {"Mixed numbers", -1, 2, -3},
        {"Zero and positive", 0, 5, -5},
        {"Positive and zero", 10, 0, 10},
        {"Zero and negative", 0, -3, 3},
        {"Negative and zero", -7, 0, -7},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            got := Minus(tc.a, tc.b)
            if got != tc.want {
                t.Errorf("Minus(%d, %d) = %d; want %d", tc.a, tc.b, got, tc.want)
            }
        })
    }
}

// TestSum_Overflow tests the Sum function for overflow.
func TestSum_Overflow(t *testing.T) {
    testCases := []struct {
        name string
        a    int
        b    int
        want int
    }{
        // Assuming int is 32-bit for this example. Adjust the test case for your specific environment if needed.
        {"MaxInt32 and positive", 2147483647, 1, -2147483648},
        {"MinInt32 and negative", -2147483648, -1, 2147483647},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            got := Sum(tc.a, tc.b)
            if got != tc.want {
                t.Errorf("Expected overflow: Sum(%d, %d) = %d; want %d", tc.a, tc.b, got, tc.want)
            }
        })
    }
}

// TestMinus_Underflow tests the Minus function for underflow.
func TestMinus_Underflow(t *testing.T) {
    testCases := []struct {
        name string
        a    int
        b    int
        want int
    }{
        // Assuming int is 32-bit for this example. Adjust the test case for your specific environment if needed.
        {"MinInt32 and positive", -2147483648, 1, 2147483647},
        {"MaxInt32 and negative", 2147483647, -1, -2147483648},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            got := Minus(tc.a, tc.b)
            if got != tc.want {
                t.Errorf("Expected underflow: Minus(%d, %d) = %d; want %d", tc.a, tc.b, got, tc.want)
            }
        })
    }
}

