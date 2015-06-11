package main

type Match struct {
	Start    int
	End      int
	Distance int
}

type StringMatcher interface {
	fuzzyMatch(pattern, text string, distance int) []Match
}

type ByStartEndDistance []Match

func (m ByStartEndDistance) Len() int      { return len(m) }
func (m ByStartEndDistance) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m ByStartEndDistance) Less(i, j int) bool {
	if m[i].Start != m[j].Start {
		return m[i].Start < m[j].Start
	}
	if m[i].End != m[j].End {
		return m[i].End < m[j].End
	}
	return m[i].Distance < m[j].Distance
}

func main() {

}
