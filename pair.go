package myGoUtils

import "fmt"

// Pair represents a pair of values of possibly different types.
type Pair[T0, T1 any] struct {
	Left  T0
	Right T1
}

// NewPair creates a new Pair with the specified values.
func NewPair[T0, T1 any](left T0, right T1) Pair[T0, T1] {
	return Pair[T0, T1]{Left: left, Right: right}
}

// Swap returns a new Pair with Left and Right values swapped.
func (this Pair[T0, T1]) Swap() Pair[T1, T0] {
	return Pair[T1, T0]{Left: this.Right, Right: this.Left}
}

// String returns a string representation of the Pair.
func (this Pair[T0, T1]) String() string {
	return fmt.Sprintf("(%v, %v)", this.Left, this.Right)
}
