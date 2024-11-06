import sys
import os
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../..")))

package main

import (
	"testing"
)

func TestSum_ZeroValues(t *testing.T) {
	a, b := 0, 0
	expected := 0

	result := Sum(a, b)
	if result != expected {
		t.Errorf("Sum(%d, %d) = %d; want %d", a, b, result, expected)
	}
}

func TestMinus_ZeroValues(t *testing.T) {
	a, b := 0, 0
	expected := 0

	result := Minus(a, b)
	if result != expected {
		t.Errorf("Minus(%d, %d) = %d; want %d", a, b, result, expected)
	}
}

func TestSum_Overflow(t *testing.T) {
	a, b := int64(9223372036854775807), int64(1)
	expected := int64(-9223372036854775808) // Overflow

	result := Sum(a, b)
	if result != expected {
		t.Errorf("Sum(%d, %d) = %d; want %d due to overflow", a, b, result, expected)
	}
}

func TestMinus_Underflow(t *testing.T) {
	a, b := int64(-9223372036854775808), int64(1)
	expected := int64(9223372036854775807) // Underflow

	result := Minus(a, b)
	if result != expected {
		t.Errorf("Minus(%d, %d) = %d; want %d due to underflow", a, b, result, expected)
	}
}

func TestMinus_LargeNegativeResult(t *testing.T) {
	a, b := int64(-1000000), int64(1000000)
	expected := int64(-2000000)

	result := Minus(a, b)
	if result != expected {
		t.Errorf("Minus(%d, %d) = %d; want %d", a, b, result, expected)
	}
}

func TestSum_LargePositiveResult(t *testing.T) {
	a, b := int64(1000000), int64(1000000)
	expected := int64(2000000)

	result := Sum(a, b)
	if result != expected {
		t.Errorf("Sum(%d, %d) = %d; want %d", a, b, result, expected)
	}
}