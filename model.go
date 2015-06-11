package main

type Match struct {
	Positions []StartEnd
	Distance  int
}

type StartEnd struct {
	Start int
	End   int
}

// find the first match; fuzzyness can split it up into multiple
// to caller must normalize both pattern and text
type StringMatcher interface {
	fuzzyMatch(pattern, text string, maxDistance int) *Match
}

func main() {

}
