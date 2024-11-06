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
        t.Errorf("Sum(%d, %d) = %d; want %d", a, b, result, expected)
    }
}

// TestSum_NegativeNumbers tests the Sum function with negative numbers.
func TestSum_NegativeNumbers(t *testing.T) {
    a := -5
    b := -3
    expected := -8

    result := Sum(a, b)

    if result != expected {
        t.Errorf("Sum(%d, %d) = %d; want %d", a, b, result, expected)
    }
}

// TestSum_MixedNumbers tests the Sum function with mixed numbers.
func TestSum_MixedNumbers(t *testing.T) {
    a := -5
    b := 3
    expected := -2

    result := Sum(a, b)

    if result != expected {
        t.Errorf("Sum(%d, %d) = %d; want %d", a, b, result, expected)
    }
}

// TestMinus_PositiveNumbers tests the Minus function with positive numbers.
func TestMinus_PositiveNumbers(t *testing.T) {
    a := 5
    b := 3
    expected := 2

    result := Minus(a, b)

    if result != expected {
        t.Errorf("Minus(%d, %d) = %d; want %d", a, b, result, expected)
    }
}

// TestMinus_NegativeNumbers tests the Minus function with negative numbers.
func TestMinus_NegativeNumbers(t *testing.T) {
    a := -5
    b := -3
    expected := -2

    result := Minus(a, b)

    if result != expected {
        t.Errorf("Minus(%d, %d) = %d; want %d", a, b, result, expected)
    }
}

// TestMinus_MixedNumbers tests the Minus function with mixed numbers.
func TestMinus_MixedNumbers(t *testing.T) {
    a := 5
    b := -3
    expected := 8

    result := Minus(a, b)

    if result != expected {
        t.Errorf("Minus(%d, %d) = %d; want %d", a, b, result, expected)
    }
}