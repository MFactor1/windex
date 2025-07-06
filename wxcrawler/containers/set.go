package containers

type Set map[string]struct{}

func NewSet() *Set {
	s := make(Set)
	return &s
}

func (s *Set) Add(item string) {
	(*s)[item] = struct{}{}
}

func (s *Set) Contains(item string) bool {
	_, exists := (*s)[item]
	return exists
}

func (s *Set) Remove(item string) {
	delete(*s, item)
}
