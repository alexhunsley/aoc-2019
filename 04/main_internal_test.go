package main

import (
	"testing"
)

func TestIsValidNumberPart1(t *testing.T) {
	if !isValidNumberPart1("123455") {
		t.Error("Expected to find a valid number, 1")
	}
	if !isValidNumberPart1("123445") {
		t.Error("Expected to find a valid number, 2")
	}
	if !isValidNumberPart1("123355") {
		t.Error("Expected to find a valid number, 3")
	}
	if !isValidNumberPart1("122456") {
		t.Error("Expected to find a valid number, 4")
	}
	if !isValidNumberPart1("002456") {
		t.Error("Expected to find a valid number, 5")
	}

	if isValidNumberPart1("654321") {
		t.Error("Expected to find an invalid number, 6")
	}
	if isValidNumberPart1("123454") {
		t.Error("Expected to find an invalid number, 7")
	}
	if isValidNumberPart1("111110") {
		t.Error("Expected to find an invalid number, 8")
	}
	if isValidNumberPart1("225506") {
		t.Error("Expected to find an invalid number, 9")
	}

	// wrong length code
	if isValidNumberPart1("") {
		t.Error("Expected to find an invalid number, 10")
	}
	if isValidNumberPart1("1") {
		t.Error("Expected to find an invalid number, 11")
	}
	if isValidNumberPart1("21") {
		t.Error("Expected to find an invalid number, 12")
	}
	if isValidNumberPart1("12345") {
		t.Error("Expected to find an invalid number, 13")
	}
	if isValidNumberPart1("1234567") {
		t.Error("Expected to find an invalid number, 14")
	}
	if isValidNumberPart1("1234567891234567890") {
		t.Error("Expected to find an invalid number, 15")
	}
}

func TestIsValidNumberPart2(t *testing.T) {
	if !isValidNumberPart2("111122") {
		t.Error("Expected to find a valid number pt2, 1")
	}
	if !isValidNumberPart2("112222") {
		t.Error("Expected to find a valid number pt2, 2")
	}
	if !isValidNumberPart2("112233") {
		t.Error("Expected to find a valid number pt2, 3")
	}


	if isValidNumberPart2("111111") {
		t.Error("Expected to find an invalid number pt2, 4")
	}
	if isValidNumberPart2("111333") {
		t.Error("Expected to find an invalid number pt2, 5")
	}
	if isValidNumberPart2("1") {
		t.Error("Expected to find an invalid number pt2, 6")
	}
	if isValidNumberPart2("11") {
		t.Error("Expected to find an invalid number pt2, 7")
	}
	if isValidNumberPart2("223454") {
		t.Error("Expected to find an invalid number pt2, 8")
	}
	if isValidNumberPart2("553333") {
		t.Error("Expected to find an invalid number pt2, 9")
	}
}
