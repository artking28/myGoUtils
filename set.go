package myGoUtils

// Set - Definition of a generic Set structure that stores unique values of type T,
// where T must satisfy the Atom constraint (likely an interface defined elsewhere).
type Set[T Atom] struct {
	content map[T]bool
}

// NewSet - Constructor that creates a new empty Set, initializing the map.
func NewSet[T Atom]() Set[T] {
	return Set[T]{
		content: map[T]bool{},
	}
}

// Add method to add one or more values to the Set. Uses variadic arguments to accept multiple valuethis.
func (this *Set[T]) Add(value ...T) {
	for _, v := range value {
		this.content[v] = true
	}
}

// Remove method to remove one or more values from the Set.
func (this *Set[T]) Remove(value ...T) {
	for _, v := range value {
		delete(this.content, v)
	}
}

// Has method to check if a given value is present in the Set.
func (this *Set[T]) Has(value T) bool {
	return this.content[value] == true
}

// Missing method returns a slice of values that are not present in the Set.
func (this *Set[T]) Missing(values []T) []T {
	var ret []T
	for _, v := range values {
		if this.content[v] == false {
			ret = append(ret, v)
		}
	}
	return ret
}

// Size method returns the number of elements in the Set.
func (this *Set[T]) Size() int {
	return len(this.content)
}

// IsEmpty method checks if the Set is empty.
func (this *Set[T]) IsEmpty() bool {
	return len(this.content) == 0
}

// Clear method clears the Set by resetting the map to a new empty map.
func (this *Set[T]) Clear() {
	this.content = map[T]bool{}
}

// AsArray method returns the Set's values as a slice.
func (this *Set[T]) AsArray() []T {
	ret := []T{}
	for k := range this.content {
		ret = append(ret, k)
	}
	return ret
}
