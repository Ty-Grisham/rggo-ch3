package main

import (
	"bytes"
	"testing"
)

// TestCountWords tests the count function set to words
func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3 word4\n")
	assertInt(t, 4, count(b, false, false))
}

// TestCountLines tests the count function set to count lines
func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("word1 word2\nline2\nline3 word1")
	assertInt(t, 3, count(b, true, false))
}

func TestCountBytes(t *testing.T) {
	b := bytes.NewBufferString("I don't know how many bytes this is")
	assertInt(t, 35, count(b, false, true))
}

// assertInt tests whther the expected int matched the resultant int
func assertInt(t *testing.T, exp, res int) {
	t.Helper()

	if res != exp {
		t.Errorf("expected %d, got %d instead", exp, res)
	}
}
