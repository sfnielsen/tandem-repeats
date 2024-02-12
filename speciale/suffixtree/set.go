package suffixtree

type Set map[int]struct{}

// Add an element to the Set
func (s Set) Add(element int) {
	s[element] = struct{}{}
}

// Remove an element from the Set
func (s Set) Remove(element int) {
	delete(s, element)
}

// Check if an element is in the Set
func (s Set) Contains(element int) bool {
	_, exists := s[element]
	return exists
}

func (s Set) Size() int {
	return len(s)
}

// Merge leaf lists
func (s Set) AddSet(set Set) {
	if s.Size() > set.Size() {
		for element := range set {
			s[element] = struct{}{}
		}
	}
}
