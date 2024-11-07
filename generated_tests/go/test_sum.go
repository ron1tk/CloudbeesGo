import sys
import os
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../..")))

package main

import (
	"testing"
)

// TestSum_PositiveNumbersZero tests the Sum function with a positive number and zero.
func TestSum_PositiveNumbersZero(t *testing.T) {
	a, b := 5, 0
	expected := 5
	result := Sum(a, b)
	if result != expected {
		t.Errorf("Sum(%d, %d) = %d; want %d", a, b, result, expected)
	}
}

// TestSum_NegativeNumbersZero tests the Sum function with a negative number and zero.
func TestSum_NegativeNumbersZero(t *testing.T) {
	a, b := -5, 0
	expected := -5
	result := Sum(a, b)
	if result != expected {
		t.Errorf("Sum(%d, %d) = %d; want %d", a, b, result, expected)
	}
}

// TestSum_LargeNumbers tests the Sum function with large numbers.
func TestSum_LargeNumbers(t *testing.T) {
	a, b := 1000000, 1000000
	expected := 2000000
	result := Sum(a, b)
	if result != expected {
		t.Errorf("Sum(%d, %d) = %d; want %d", a, b, result, expected)
	}
}

// TestMinus_PostiveAndNegativeNumbers tests the Minus function with a positive and a negative number.
func TestMinus_PostiveAndNegativeNumbers(t *testing.T) {
	a, b := 5, -3
	expected := 8
	result := Minus(a, b)
	if result != expected {
		t.Errorf("Minus(%d, %d) = %d; want %d", a, b, result, expected)
	}
}

// TestMinus_NegativeAndPositiveNumbers tests the Minus function with a negative and a positive number.
func TestMinus_NegativeAndPositiveNumbers(t *testing.T) {
	a, b := -5, 3
	expected := -8
	result := Minus(a, b)
	if result != expected {
		t.Errorf("Minus(%d, %d) = %d; want %d", a, b, result, expected)
	}
}

// TestMinus_LargeNumbers tests the Minus function with large numbers.
func TestMinus_LargeNumbers(t *testing.T) {
	a, b := 1000000, 500000
	expected := 500000
	result := Minus(a, b)
	if result != expected {
		t.Errorf("Minus(%d, %d) = %d; want %d", a, b, result, expected)
	}
}

// TestMinus_SameNumbers tests the Minus function with two equal numbers.
func TestMinus_SameNumbers(t *testing.T) {
	a, b := 5, 5
	expected := 0
	result := Minus(a, b)
	if result != expected {
		t.Errorf("Minus(%d, %d) = %d; want %d", a, b, result, expected)
	}
}