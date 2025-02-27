package idforge

import (
	"testing"
)

func TestWithAlphabet(t *testing.T) {
	gen := &Generator{}

	WithAlphabet("0123456789")(gen)

	if gen.alphabet != "0123456789" {
		t.Errorf("Expected alphabet to be '0123456789', got '%s'", gen.alphabet)
	}
}

func TestWithAlphabetInvalidLength(t *testing.T) {
	gen := &Generator{alphabet: "abcdefghijklmnopqrstuvwxyz"}

	WithAlphabet("a")(gen)

	if gen.alphabet != "abcdefghijklmnopqrstuvwxyz" {
		t.Errorf("Expected alphabet to be 'abcdefghijklmnopqrstuvwxyz', got '%s'", gen.alphabet)
	}
}

func TestWithSize(t *testing.T) {
	gen := &Generator{}

	WithSize(10)(gen)

	if gen.size != 10 {
		t.Errorf("Expected size to be 10, got %d", gen.size)
	}
}

func TestWithSizeInvalid(t *testing.T) {
	gen := &Generator{size: 20}

	WithSize(-5)(gen)

	if gen.size != 20 {
		t.Errorf("Expected size to be 20, got %d", gen.size)
	}
}

func TestWithMinSize(t *testing.T) {
	gen := &Generator{size: 20}

	WithMinSize(15)(gen)

	if gen.size != 15 {
		t.Errorf("Expected size to be 15, got %d", gen.size)
	}
}

func TestWithMinSizeInvalid(t *testing.T) {
	gen := &Generator{size: 20}

	WithMinSize(25)(gen)

	if gen.size != 20 {
		t.Errorf("Expected size to be 20, got %d", gen.size)
	}
}

func TestWithMaxSize(t *testing.T) {
	gen := &Generator{size: 10}

	WithMaxSize(15)(gen)

	if gen.size != 15 {
		t.Errorf("Expected size to be 15, got %d", gen.size)
	}
}

func TestWithMaxSizeInvalid(t *testing.T) {
	gen := &Generator{size: 20}

	WithMaxSize(15)(gen)

	if gen.size != 20 {
		t.Errorf("Expected size to be 20, got %d", gen.size)
	}
}
