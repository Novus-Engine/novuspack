package testhelpers

import (
	"testing"
)

func TestContains(t *testing.T) {
	t.Run("finds substring", func(t *testing.T) {
		if !Contains("hello world", "world") {
			t.Error("Contains should find 'world' in 'hello world'")
		}
	})

	t.Run("case sensitive", func(t *testing.T) {
		if Contains("hello world", "WORLD") {
			t.Error("Contains should be case-sensitive")
		}
	})

	t.Run("empty substring", func(t *testing.T) {
		if !Contains("hello", "") {
			t.Error("Contains should return true for empty substring")
		}
	})

	t.Run("substring not found", func(t *testing.T) {
		if Contains("hello", "goodbye") {
			t.Error("Contains should return false when substring not found")
		}
	})

	t.Run("exact match", func(t *testing.T) {
		if !Contains("test", "test") {
			t.Error("Contains should return true for exact match")
		}
	})
}

func TestContainsIgnoreCase(t *testing.T) {
	t.Run("finds substring case insensitive", func(t *testing.T) {
		if !ContainsIgnoreCase("hello world", "WORLD") {
			t.Error("ContainsIgnoreCase should find 'WORLD' in 'hello world'")
		}
	})

	t.Run("mixed case", func(t *testing.T) {
		if !ContainsIgnoreCase("Hello World", "hElLo") {
			t.Error("ContainsIgnoreCase should find 'hElLo' in 'Hello World'")
		}
	})

	t.Run("lowercase in uppercase", func(t *testing.T) {
		if !ContainsIgnoreCase("HELLO WORLD", "world") {
			t.Error("ContainsIgnoreCase should find 'world' in 'HELLO WORLD'")
		}
	})

	t.Run("uppercase in lowercase", func(t *testing.T) {
		if !ContainsIgnoreCase("hello world", "HELLO") {
			t.Error("ContainsIgnoreCase should find 'HELLO' in 'hello world'")
		}
	})

	t.Run("empty substring", func(t *testing.T) {
		if !ContainsIgnoreCase("hello", "") {
			t.Error("ContainsIgnoreCase should return true for empty substring")
		}
	})

	t.Run("substring not found", func(t *testing.T) {
		if ContainsIgnoreCase("hello", "goodbye") {
			t.Error("ContainsIgnoreCase should return false when substring not found")
		}
	})

	t.Run("exact match different case", func(t *testing.T) {
		if !ContainsIgnoreCase("test", "TEST") {
			t.Error("ContainsIgnoreCase should return true for exact match with different case")
		}
	})
}

func TestIndexIgnoreCase(t *testing.T) {
	t.Run("finds index case insensitive", func(t *testing.T) {
		idx := IndexIgnoreCase("hello world", "WORLD")
		if idx != 6 {
			t.Errorf("IndexIgnoreCase should return 6 for 'WORLD' in 'hello world', got %d", idx)
		}
	})

	t.Run("finds index at start", func(t *testing.T) {
		idx := IndexIgnoreCase("Hello World", "hElLo")
		if idx != 0 {
			t.Errorf("IndexIgnoreCase should return 0 for 'hElLo' at start, got %d", idx)
		}
	})

	t.Run("finds lowercase in uppercase", func(t *testing.T) {
		idx := IndexIgnoreCase("HELLO WORLD", "world")
		if idx != 6 {
			t.Errorf("IndexIgnoreCase should return 6 for 'world' in 'HELLO WORLD', got %d", idx)
		}
	})

	t.Run("finds uppercase in lowercase", func(t *testing.T) {
		idx := IndexIgnoreCase("hello world", "HELLO")
		if idx != 0 {
			t.Errorf("IndexIgnoreCase should return 0 for 'HELLO' in 'hello world', got %d", idx)
		}
	})

	t.Run("substring not found", func(t *testing.T) {
		idx := IndexIgnoreCase("hello", "goodbye")
		if idx != -1 {
			t.Errorf("IndexIgnoreCase should return -1 when substring not found, got %d", idx)
		}
	})

	t.Run("empty substring", func(t *testing.T) {
		idx := IndexIgnoreCase("hello", "")
		if idx != 0 {
			t.Errorf("IndexIgnoreCase should return 0 for empty substring, got %d", idx)
		}
	})

	t.Run("exact match different case", func(t *testing.T) {
		idx := IndexIgnoreCase("test", "TEST")
		if idx != 0 {
			t.Errorf("IndexIgnoreCase should return 0 for exact match with different case, got %d", idx)
		}
	})

	t.Run("multiple occurrences returns first", func(t *testing.T) {
		idx := IndexIgnoreCase("hello hello", "HELLO")
		if idx != 0 {
			t.Errorf("IndexIgnoreCase should return first occurrence index 0, got %d", idx)
		}
	})
}
