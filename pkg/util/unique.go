package util

import "sort"

// SortUnique sorts and deduplicates a slice of any type.
// It first sorts the slice using sort.Slice, based on the order defined by the cmp function.
// cmp is a comparison function used to determine the order of two elements in the slice.
// If cmp(i, j) < 0, it indicates that element i should come before element j.
// Once the slice is sorted, SortUnique then iterates over the sorted slice,
// adding only the unique elements to the result slice.
// If cmp(i, j) == 0, it is considered that elements i and j are equal,
// This function returns a new slice containing the unique elements of the original slice,
// ordered according to the order defined by the cmp function.
//
// Example:
//
//	slice := []int{3, 1, 2, 3, 4, 1}
//	uniqueSlice := SortUnique(slice, func(i, j int) bool {
//	    return slice[i] < slice[j]
//	})
//	fmt.Println("Unique slice:", uniqueSlice)
//
// This example will output: Unique slice: [1 2 3 4]
// Note: This function modifies the contents of the input slice.
func SortUnique[T any](slice []T, cmp func(i, j int) int) []T {
	if len(slice) == 0 {
		return []T{}
	}
	sort.Slice(slice, func(i, j int) bool {
		return cmp(i, j) < 0
	})
	result := []T{slice[0]}
	for i := 1; i < len(slice); i++ {
		if cmp(i-1, i) != 0 {
			result = append(result, slice[i])
		}
	}
	return result
}

// Unique removes duplicate elements from a slice.
// It iterates over the given slice and uses a custom equals function to determine if an element is a duplicate.
// The equals function takes two integer parameters i and j, which are indices of elements in the slice.
// If equals(i, j) returns true, it indicates that the elements at indices i and j are equal.
// For each element, Unique checks whether it already exists in the result slice.
// The element is added to the result slice only if it is not a duplicate.
// Ultimately, this function returns a new slice containing only the unique elements from the original slice.
//
// Example:
//
//	slice := []string{"apple", "banana", "apple", "orange"}
//	uniqueSlice := Unique(slice, func(i, j int) bool {
//	    return slice[i] == slice[j]
//	})
//	fmt.Println("Unique slice:", uniqueSlice)
//
// This example will output: Unique slice: ["apple" "banana" "orange"]
// Note that this function does not modify the contents of the original slice.
func Unique[T any](slice []T, equals func(i, j int) bool) []T {
	var result []T
	for i := 0; i < len(slice); i++ {
		duplicate := false
		for j := 0; j < len(result); j++ {
			if equals(i, j) {
				duplicate = true
				break
			}
		}
		if !duplicate {
			result = append(result, slice[i])
		}
	}
	return result
}
