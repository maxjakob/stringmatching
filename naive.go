package main

import (
	// "fmt"
	"sync"
)

type naiveMatcher struct{}

func (b naiveMatcher) fuzzyMatch(pattern, text string, maxDistance int) *Match {
	var result *Match
	for startPos := 0; startPos < len(text)-len(pattern)+1; startPos++ {
		if text[startPos] == pattern[0] {
			match := getMatch(pattern, text, maxDistance, startPos)
			if result == nil || preferSecond(result, match) {
				result = match
			}
		}
	}
	return result
}

// finds highlighting matches from startPos
func getMatch(pattern, text string, maxDistance, startPos int) *Match {
	textPos := startPos
	patternPos := 0
	distance := 0

	matchStart := -1
	var result Match
	for ; distance <= maxDistance; textPos++ {
		if text[textPos] == pattern[patternPos] { // char matches
			if matchStart == -1 { // remember start
				matchStart = textPos
			}
			patternPos++ // consume pattern
		} else { // char does not match
			if matchStart != -1 { // if there was a start set
				result.Positions = append(result.Positions, StartEnd{matchStart, textPos}) // emit the interval
				matchStart = -1                                                            // forget the old start
			}
			distance++ // increase number of non-matches
		}
		if patternPos == len(pattern) {
			break // consumed complete pattern
		}
	}
	result.Distance = distance
	if matchStart != -1 { // end of the string
		result.Positions = append(result.Positions, StartEnd{matchStart, textPos + 1})
	}
	if distance <= maxDistance && len(result.Positions) > 0 {
		return &result
	}
	return nil
}

// prefer old distance, longer, earlier matches
func preferSecond(old, potentialNew *Match) bool {
	if old.Distance > potentialNew.Distance {
		return true
	}
	lengthA := old.Positions[len(old.Positions)-1].End - old.Positions[0].Start
	lengthB := potentialNew.Positions[len(old.Positions)-1].End - potentialNew.Positions[0].Start
	if lengthA < lengthB {
		return true
	}
	if old.Positions[0].Start > potentialNew.Positions[0].Start {
		return true
	}
	return false
}

// experimental
func (b naiveMatcher) fuzzyMatchParallel(pattern, text string, maxDistance, threads int) *Match {
	ingress := make(chan int)
	egress := make(chan *Match)

	// scatter
	go func() {
		for start := 0; start < len(text)-len(pattern); start++ {
			if text[start] == pattern[0] {
				ingress <- int(start)
			}
		}
		close(ingress)
	}()

	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for startPos := range ingress {
				match := getMatch(pattern, text, maxDistance, startPos)
				if match != nil {
					egress <- match
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(egress)
	}()

	// gather
	var result *Match
	for match := range egress {
		if result == nil || preferSecond(result, match) {
			result = match
		}
	}
	return result
}
