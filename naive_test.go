package main

import (
	"testing"
)

func check(gotP, expectedP *Match, t *testing.T) {
	if expectedP == nil && gotP == nil {
		return
	}
	if expectedP == nil && gotP != nil {
		t.Fatalf("expected nil got: %v", *gotP)
	}
	if expectedP != nil && gotP == nil {
		t.Fatalf("got nil expected: %v", *expectedP)
	}
	got := *gotP
	expected := *expectedP
	if got.Distance != expected.Distance {
		t.Fatalf("incorrect distance! expected: %d, got: %d", expected.Distance, got.Distance)
	}
	if len(got.Positions) != len(expected.Positions) {
		t.Fatalf("number of matches incorrect! expected: %d (%v), got: %d (%v)",
			len(expected.Positions), expected.Positions, len(got.Positions), got.Positions)
	}
	for i := 0; i < len(got.Positions); i++ {
		if got.Positions[i] != expected.Positions[i] {
			t.Fatalf("match at index %d incorrect! expected: %v, got: %v", i, expected.Positions[i], got.Positions[i])
		}
	}
}

var naive = naiveMatcher{}

func TestExactMatchNaive(t *testing.T) {
	distance := 0
	pattern := "abc"

	text := "XXXabcYYY"
	expected := &Match{[]StartEnd{StartEnd{3, 6}}, 0}
	check(naive.fuzzyMatch(pattern, text, distance), expected, t)

	text = "XXXacYYY"
	expected = nil
	check(naive.fuzzyMatch(pattern, text, distance), expected, t)
}

func TestExactOnlyFirstNaive(t *testing.T) {
	distance := 0

	pattern := "abc"
	text := "XXXabcYYYabcZZZ"
	expected := &Match{[]StartEnd{StartEnd{3, 6}}, 0}
	check(naive.fuzzyMatch(pattern, text, distance), expected, t)
}

func TestDistanceOneNaive(t *testing.T) {
	distance := 1
	pattern := "abc"

	text := "XXXabbcYYY"
	expected := &Match{[]StartEnd{StartEnd{3, 5}, StartEnd{6, 7}}, 1}
	check(naive.fuzzyMatch(pattern, text, distance), expected, t)

	text = "XXXaZZbZcYYY"
	expected = nil
	check(naive.fuzzyMatch(pattern, text, distance), expected, t)
}

func TestDistanceOneNotAllConsumedNaive(t *testing.T) {
	distance := 1
	pattern := "abc"

	text := "XXXabYYY"
	expected := &Match{[]StartEnd{StartEnd{3, 5}}, 1}
	check(naive.fuzzyMatch(pattern, text, distance), expected, t)
}

func TestDistanceTwoNaive(t *testing.T) {
	distance := 2
	pattern := "abc"

	text := "XXXaZZbcYYY"
	expected := &Match{[]StartEnd{StartEnd{3, 4}, StartEnd{6, 8}}, 2}
	check(naive.fuzzyMatch(pattern, text, distance), expected, t)

	text = "XXXaZZZbcYYY"
	expected = nil
	check(naive.fuzzyMatch(pattern, text, distance), expected, t)
}
