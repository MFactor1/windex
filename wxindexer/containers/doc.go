package containers

type Doc struct {
	Body string
	Links []string
}

type WordFrequencies struct {
	Words map[string]int
}
