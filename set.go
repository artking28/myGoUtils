package myGoUtils

type Set[T Atom] struct {
	content *map[T]bool
}

func NewSet[T Atom]() Set[T] {
	return Set[T]{
		content: &map[T]bool{},
	}
}

func (this *Set[T]) Add(value ...T) {
	for _, v := range value {
		(*this.content)[v] = true
	}
}

func (this *Set[T]) Remove(value ...T) {
	for _, v := range value {
		delete(*this.content, v)
	}
}

func (this *Set[T]) Has(value T) bool {
	return (*this.content)[value] == true
}

func (this *Set[T]) Missing(values []T) []T {
	var ret []T
	for _, v := range values {
		if (*this.content)[v] == false {
			ret = append(ret, v)
		}
	}
	return ret
}

func (this *Set[T]) Size() int {
	return len(*this.content)
}

func (this *Set[T]) IsEmpty() bool {
	return len(*this.content) == 0
}

func (this *Set[T]) Clear() {
	this.content = &map[T]bool{}
}

func (this *Set[T]) AsArray() []T {
	ret := []T{}
	for k := range *this.content {
		ret = append(ret, k)
	}
	return ret
}
