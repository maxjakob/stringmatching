package main

import (
	"testing"
)

type namedMatcher struct {
	matcher StringMatcher
	name    string
}

var matchers = []namedMatcher{
	namedMatcher{naiveMatcher{}, "naive"},
	namedMatcher{bitapMatcher{}, "bitap"},
}

func check(pattern, text string, distance int, expected []Match, t *testing.T) {
	for _, nm := range matchers {
		got := nm.matcher.fuzzyMatch(pattern, text, distance)
		if len(got) != len(expected) {
			t.Fatalf("%s matcher: number of matches incorrect! expected: %d (%v), got: %d (%v)", nm.name, len(expected), expected, len(got), got)
		}
		for i := 0; i < len(got); i++ {
			if got[i] != expected[i] {
				t.Fatalf("%s matcher: match at index %d incorrect! expected: %v, got: %v", nm.name, i, expected[i], got[i])
			}
		}
	}
}

func TestExactMatch(t *testing.T) {
	distance := 0
	pattern := "abc"
	text := "XXXabcYYYabcZZZ"
	expected := []Match{Match{3, 6, 0}, Match{9, 12, 0}}
	check(pattern, text, distance, expected, t)
}

func TestDistanceOne(t *testing.T) {
	distance := 1
	pattern := "abc"
	text := "XXXabbcYYYabaZZZ"
	expected := []Match{Match{3, 7, 1}, Match{10, 13, 1}}
	check(pattern, text, distance, expected, t)
}

func TestDistanceTwo(t *testing.T) {
	distance := 2
	pattern := "abc"
	text := "XXXabbcYYYabaZZZaBCc"
	expected := []Match{Match{3, 7, 1}, Match{10, 13, 1}, Match{16, 19, 2}}
	check(pattern, text, distance, expected, t)
}

// func TestDistanceOneSpacesNaive(t *testing.T) {
// 	distance := 1
// 	pattern := "abc"
// 	text := "XXXab YYY"
// 	expected := []Match{Match{3, 5, 1}}
//  check(pattern, text, distance, expected, t)
// }
