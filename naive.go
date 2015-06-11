package main

import (
	// "fmt"
	"sort"
	"sync"
)

type naiveMatcher struct{}

// to caller must normalize both pattern and text
func (b naiveMatcher) fuzzyMatch(pattern, text string, maxDistance int) []Match {
	ingress := make(chan int)
	egress := make(chan Match)
	concurrency := 1

	//TODO is the parallelism worth it?

	// scatter
	go func() {
		for start := 0; start < len(text)-len(pattern); start++ {
			if text[start] == pattern[0] {
				ingress <- int(start)
				// concurrency++
			}
		}
		close(ingress)
	}()

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for startPos := range ingress {
				d := int(0)
				patternPos := 0
				for textPos := startPos; d <= maxDistance; textPos++ {
					if patternPos == len(pattern) {
						egress <- Match{startPos, textPos, d}
						break // pattern completely consumed
					}
					if patternPos+d == len(pattern) {
						egress <- Match{startPos, textPos, d}
						// there might be more matches
					}
					if text[textPos] == pattern[patternPos] {
						patternPos++
					} else {
						d++
					}

				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(egress)
	}()

	// gather
	matches := []Match{}
	for match := range egress {
		matches = append(matches, match)
	}
	if len(matches) == 0 {
		return matches
	}
	sort.Sort(ByStartEndDistance(matches))

	return resolveOverlaps(matches)
}

// assumes sorted input
func resolveOverlaps(matches []Match) []Match {
	result := []Match{}
	last := matches[0]
	for j := 1; j < len(matches); j++ {
		this := matches[j]
		if overlap(last, this) {
			this = preference(last, this)
		} else {
			result = append(result, last)
		}
		last = this
	}
	result = append(result, last)
	return result
}

func overlap(a, b Match) bool {
	if a.Start == b.Start || a.End == b.End {
		return true
	}
	// start of b in a
	if a.Start < b.Start && a.End > b.Start {
		return true
	}
	// start of a in b
	if a.Start > b.Start && a.End < b.Start {
		return true
	}
	return false
}

// prefer less distance, longer, earlier matches
func preference(a, b Match) Match {
	if a.Distance < b.Distance {
		return a
	}
	if a.Distance > b.Distance {
		return b
	}
	lengthA := a.End - a.Start
	lengthB := b.End - b.Start
	if lengthA > lengthB {
		return a
	}
	if lengthA < lengthB {
		return b
	}
	if a.Start < b.Start {
		return a
	}
	if a.Start > b.Start {
		return b
	}
	return a // same same
}
