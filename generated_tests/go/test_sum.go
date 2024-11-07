import sys
import os
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../..")))


import (
    "testing"
    "gocode"
)

// TestSum_PositiveNumbers tests the Sum function with positive numbers.
func TestSum_PositiveNumbers(t *testing.T) {
    a, b := 5, 3
    expected := 8
    result := Sum(a, b)
    if result != expected {
        t.Errorf("Sum(%d, %d) = %d; want %d", a, b, result, expected)
    }
}

// TestSum_NegativeNumbers tests the Sum function with negative numbers.
func TestSum_NegativeNumbers(t *testing.T) {
    a, b := -5, -3
    expected := -8
    result := Sum(a, b)
    if result != expected {
        t.Errorf("Sum(%d, %d) = %d; want %d", a, b, result, expected)
    }
}

// TestSum_Zero tests the Sum function with zero.
func TestSum_Zero(t *testing.T) {
    a, b := 0, 0
    expected := 0
    result := Sum(a, b)
    if result != expected {
        t.Errorf("Sum(%d, %d) = %d; want %d", a, b, result, expected)
    }
}

// TestMinus_PositiveNumbers tests the Minus function with positive numbers.
func TestMinus_PositiveNumbers(t *testing.T) {
    a, b := 5, 3
    expected := 2
    result := Minus(a, b)
    if result != expected {
        t.Errorf("Minus(%d, %d) = %d; want %d", a, b, result, expected)
    }
}

// TestMinus_NegativeNumbers tests the Minus function with negative numbers.
func TestMinus_NegativeNumbers(t *testing.T) {
    a, b := -5, -3
    expected := -2
    result := Minus(a, b)
    if result != expected {
        t.Errorf("Minus(%d, %d) = %d; want %d", a, b, result, expected)
    }
}

// TestMinus_Zero tests the Minus function with zero.
func TestMinus_Zero(t *testing.T) {
    a, b := 0, 0
    expected := 0
    result := Minus(a, b)
    if result != expected {
        t.Errorf("Minus(%d, %d) = %d; want %d", a, b, result, expected)
    }
}

// TestMinus_LargerSubtrahend tests the Minus function with a larger subtrahend.
func TestMinus_LargerSubtrahend(t *testing.T) {
    a, b := 3, 5
    expected := -2
    result := Minus(a, b)
    if result != expected {
        t.Errorf("Minus(%d, %d) = %d; want %d", a, b, result, expected)
    }
}