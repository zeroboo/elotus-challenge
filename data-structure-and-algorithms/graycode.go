package elotuschallenge

import (
	"math"
)

// grayCode generates an n-bit gray code sequence.
// A gray code is a sequence of binary numbers where two successive values differ in only one bit.
func grayCode(n int) []int {
	// The number of integers in an n-bit gray code sequence is 2^n.
	// We use math.Pow to calculate this. It returns a float64, so we cast it to int.
	sequenceSize := int(math.Pow(2, float64(n)))

	// Create a slice to hold the resulting sequence.
	result := make([]int, sequenceSize)

	// from 0 up to 2^n - 1.
	for i := 0; i < sequenceSize; i++ {
		// The formula to generate the i-th gray code is i XOR (i >> 1).
		// i >> 1 is a bitwise right shift, which is equivalent to integer division by 2.
		// The XOR operation with its right-shifted self produces the reflective property of gray codes.
		result[i] = i ^ (i >> 1)
	}

	return result
}

func CountDifferentBits(a, b int) int {
	count := 0
	xor := a ^ b
	for xor > 0 {
		count += xor & 1
		xor >>= 1
	}
	return count
}

// FindPossibleGrayCodes returns possible gray codes which differ by one bit from the given number value
func FindPossibleGrayCodes(value int, max int) []int {
	var grayCodes []int
	for i := 0; i <= max; i++ {
		if i != value {
			if CountDifferentBits(value, i) == 1 {
				grayCodes = append(grayCodes, i)
			}
		}

	}
	return grayCodes
}

func NotInSlice[T comparable](slice []T, val T) bool {
	for _, item := range slice {
		if item == val {
			return false // Found the value, so it IS in the slice
		}
	}
	return true // Did not find the value
}

// RemoveFromSlice removes all occurrences of a specific value from a slice.
// It is a generic function that works with any comparable type.
// It returns a new slice without the specified value and does not modify the original.
func RemoveFromSlice[T comparable](slice []T, valueToRemove T) []T {
	// Create a new slice to store the result.
	// We initialize it with a capacity of the original slice for efficiency,
	// as it's the maximum possible size.
	result := make([]T, 0, len(slice))

	// Iterate over the input slice.
	for _, item := range slice {
		// If the current item does not match the value to remove,
		// append it to our result slice.
		if item != valueToRemove {
			result = append(result, item)
		}
	}

	// Return the new slice that contains only the desired elements.
	return result
}

// FindFirstDuplicate finds the first value that appears more than once in a slice.
// It is a generic function that works for any comparable type (int, string, etc.).
// It returns the duplicated value and a boolean 'true' if a duplicate is found.
// If no duplicates exist, it returns the zero value for the type and 'false'.
func FindFirstDuplicate[T comparable](slice []T) (T, bool) {
	// 'seen' will function as a set to store the elements we have encountered.
	// Using a map with an empty struct `struct{}` is a common Go idiom
	// because the empty struct consumes no memory.
	seen := make(map[T]struct{})

	// Iterate over each item in the input slice.
	for _, item := range slice {
		// The "comma ok" idiom checks if the key (the item) exists in the map.
		if _, ok := seen[item]; ok {
			// If 'ok' is true, the item is already in our map, which means
			// we have found the first duplicate. Return it immediately.
			return item, true
		}
		// If the item is not in the map, add it to our set so we can
		// detect it if it appears again later in the slice.
		seen[item] = struct{}{}
	}

	// If the loop completes without finding any duplicates, we return the
	// zero value for the slice's type (e.g., 0 for int, "" for string)
	// and false to indicate that no duplicate was found.
	var zeroValue T
	return zeroValue, false
}
