package sliceutil

import (
	"slices"
	"strconv"
)

// PruneEmptyStrings from the slice
func PruneEmptyStrings(v []string) []string {
	return PruneEqual(v, "")
}

// PruneEqual removes items from the slice equal to the specified value
func PruneEqual[T comparable](inputSlice []T, equalTo T) (r []T) {
	for i := range inputSlice {
		if inputSlice[i] != equalTo {
			r = append(r, inputSlice[i])
		}
	}

	return
}

// Dedupe removes duplicates from a slice of elements preserving the order
func Dedupe[T comparable](inputSlice []T) (result []T) {
	seen := make(map[T]struct{})
	for _, inputValue := range inputSlice {
		if _, ok := seen[inputValue]; !ok {
			seen[inputValue] = struct{}{}

			result = append(result, inputValue)
		}
	}

	return
}

// Contains if a slice contains an element
func Contains[T comparable](inputSlice []T, element T) bool {
	return slices.Contains(inputSlice, element)
}

// ContainsItems checks if s1 contains s2
func ContainsItems[T comparable](s1 []T, s2 []T) bool {
	return slices.ContainsFunc(s1, func(e T) bool {
		return Contains(s2, e)
	})
}

// ToInt converts a slice of strings to a slice of ints
func ToInt(s []string) ([]int, error) {
	var ns []int

	for _, ss := range s {
		n, err := strconv.Atoi(ss)
		if err != nil {
			return nil, err
		}

		ns = append(ns, n)
	}

	return ns, nil
}

// Equal checks if the items of two slices are equal respecting the order
func Equal[T comparable](s1, s2 []T) bool {
	return slices.Equal(s1, s2)
}

// IsEmpty checks if the slice has length zero
func IsEmpty[V comparable](s []V) bool {
	return len(s) == 0
}

// ElementsMatch asserts that the specified listA(array, slice...) is equal to specified
// listB(array, slice...) ignoring the order of the elements. If there are duplicate elements,
// the number of appearances of each of them in both lists should match
func ElementsMatch[V comparable](s1, s2 []V) bool {
	if IsEmpty(s1) && IsEmpty(s2) {
		return true
	}

	extraS1, extrS2 := Diff(s1, s2)

	return IsEmpty(extraS1) && IsEmpty(extrS2)
}

// Diff calculates the extra elements between two sequences
func Diff[V comparable](s1, s2 []V) (extraS1, extraS2 []V) {
	s1Len := len(s1)
	s2Len := len(s2)

	visited := make([]bool, s2Len)

	for i := 0; i < s1Len; i++ {
		element := s1[i]
		found := false

		for j := 0; j < s2Len; j++ {
			if visited[j] {
				continue
			}

			if s2[j] == element {
				visited[j] = true
				found = true

				break
			}
		}

		if !found {
			extraS1 = append(extraS1, element)
		}
	}

	for j := 0; j < s2Len; j++ {
		if visited[j] {
			continue
		}

		extraS2 = append(extraS2, s2[j])
	}

	return
}

// Merge and dedupe multiple items
func Merge[V comparable](ss ...[]V) []V {
	var final []V

	for _, s := range ss {
		final = append(final, s...)
	}

	return Dedupe(final)
}

// MergeItems takes in multiple items of a comparable type and merges them into a single slice
// while removing any duplicates. It then returns the resulting slice with duplicate elements removed
func MergeItems[V comparable](items ...V) []V {
	return Dedupe(items)
}

// FirstNonZero function takes a slice of comparable type inputs, and returns
// the first non-zero element in the slice along with a boolean value indicating if a non-zero element was found or not
func FirstNonZero[T comparable](inputs []T) (T, bool) {
	var zero T

	for _, v := range inputs {
		if v != zero {
			return v, true
		}
	}

	return zero, false
}

// Clone a slice through built-in copy
func Clone[T comparable](t []T) []T {
	return slices.Clone(t)
}
