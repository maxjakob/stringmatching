package main

import "fmt"

type bitapMatcher struct{}

func (b bitapMatcher) fuzzyMatch(pattern, text string, distance int) []Match {
	matches := []Match{}

	var (
		m                = uint(len(pattern))
		maxVal      uint = 255 //
		patternMask      = make([]int, maxVal+1)
		R                = make([]int, maxVal+1)
		i, d        uint
	)

	if m == 0 {
		return nil
	}

	if m > 31 {
		fmt.Println("pattern too long")
		return nil
	}
	fmt.Println("pattern", pattern)

	// Init bit array
	for i = 0; i <= uint(distance); i++ {
		R[i] = ^1
	}

	// Init pattern bitmasks
	for i = 0; i <= maxVal; i++ {
		patternMask[i] = ^0
	}
	for i = 0; i < m; i++ {
		patternMask[pattern[i]] &= ^(1 << i)
		if patternMask[pattern[i]] != -1 {
			fmt.Println("spot", i, pattern[i], patternMask[pattern[i]])
		}
	}

	fmt.Println("m", m)
	fmt.Println("patternMask", patternMask)
	fmt.Println("R", R)

	// Match
	for i = 0; i < uint(len(text)); i++ {
		old_Rd1 := R[0]

		R[0] |= patternMask[text[i]]
		R[0] <<= 1

		for d = 1; d <= i; d++ {
			tmp := R[d]
			// Substitution is all we care about. // TODO(mj) not so !
			R[d] = (old_Rd1 & (R[d] | patternMask[text[i]])) << 1
			old_Rd1 = tmp
		}

		if (R[distance] & (1 << m)) == 0 {
			fmt.Println("R", R)
			matches = append(matches, Match{int(i - m + 1), len(text), -1}) // TODO(mj) real end and distance
		}
	}

	return matches
}
