import sys
import os
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../..")))

package main

import (
    "testing"
)

// TestSum_PositiveNumbers tests the Sum function with positive numbers.
func TestSum_PositiveNumbers(t *testing.T) {
    a := 5
    b := 3
    expected := 8
    result := Sum(a, b)
    if result != expected {
        t.Errorf("Sum of %d and %d was incorrect, got: %d, want: %d.", a, b, result, expected)
    }
}

// TestSum_NegativeNumbers tests the Sum function with negative numbers.
func TestSum_NegativeNumbers(t *testing.T) {
    a := -5
    b := -3
    expected := -8
    result := Sum(a, b)
    if result != expected {
        t.Errorf("Sum of %d and %d was incorrect, got: %d, want: %d.", a, b, result, expected)
    }
}

// TestSum_Zero tests the Sum function with zero values.
func TestSum_Zero(t *testing.T) {
    a := 0
    b := 0
    expected := 0
    result := Sum(a, b)
    if result != expected {
        t.Errorf("Sum of %d and %d was incorrect, got: %d, want: %d.", a, b, result, expected)
    }
}

// TestMinus_PositiveNumbers tests the Minus function with positive numbers.
func TestMinus_PositiveNumbers(t *testing.T) {
    a := 5
    b := 3
    expected := 2
    result := Minus(a, b)
    if result != expected {
        t.Errorf("Minus of %d and %d was incorrect, got: %d, want: %d.", a, b, result, expected)
    }
}

// TestMinus_NegativeNumbers tests the Minus function with negative numbers.
func TestMinus_NegativeNumbers(t *testing.T) {
    a := -5
    b := -3
    expected := -2
    result := Minus(a, b)
    if result != expected {
        t.Errorf("Minus of %d and %d was incorrect, got: %d, want: %d.", a, b, result, expected)
    }
}

// TestMinus_Zero tests the Minus function with zero values.
func TestMinus_Zero(t *testing.T) {
    a := 0
    b := 0
    expected := 0
    result := Minus(a, b)
    if result != expected {
        t.Errorf("Minus of %d and %d was incorrect, got: %d, want: %d.", a, b, result, expected)
    }
}

// TestMinus_PositiveAndNegativeNumbers tests the Minus function with a mix of positive and negative numbers.
func TestMinus_PositiveAndNegativeNumbers(t *testing.T) {
    a := 5
    b := -3
    expected := 8
    result := Minus(a, b)
    if result != expected {
        t.Errorf("Minus of %d and %d was incorrect, got: %d, want: %d.", a, b, result, expected)
    }
}

// TestSum_LargeNumbers tests the Sum function with large numbers.
func TestSum_LargeNumbers(t *testing.T) {
    a := 1000000
    b := 2000000
    expected := 3000000
    result := Sum(a, b)
    if result != expected {
        t.Errorf("Sum of %d and %d was incorrect, got: %d, want: %d.", a, b, result, expected)
    }
}