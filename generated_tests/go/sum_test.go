package main

import (
	"testing"
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

// TestSum_MixedNumbers tests the Sum function with one positive and one negative number.
func TestSum_MixedNumbers(t *testing.T) {
    a, b := -5, 3
    expected := -2
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

// TestMinus_MixedNumbers tests the Minus function with one positive and one negative number.
func TestMinus_MixedNumbers(t *testing.T) {
    a, b := 5, -3
    expected := 8
    result := Minus(a, b)
    if result != expected {
        t.Errorf("Minus(%d, %d) = %d; want %d", a, b, result, expected)
    }
}

// TestMinus_ZeroResult tests the Minus function with a result of zero.
func TestMinus_ZeroResult(t *testing.T) {
    a, b := 3, 3
    expected := 0
    result := Minus(a, b)
    if result != expected {
        t.Errorf("Minus(%d, %d) = %d; want %d", a, b, result, expected)
    }
}