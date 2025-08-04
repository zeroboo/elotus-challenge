package elotuschallenge

import "testing"

func TestFindLength_Correct(t *testing.T) {
	nums1 := []int{1, 2, 3, 2, 1}
	nums2 := []int{3, 2, 1, 4, 7}
	expected := 3 // The longest common subarray is [2, 1] or [1, 2]
	result := findLength(nums1, nums2)
	if result != expected {
		t.Errorf("Expected %d but got %d", expected, result)
	}
}

func TestFindLength_Example1_Correct(t *testing.T) {
	nums1 := []int{1, 2, 3, 2, 1}
	nums2 := []int{3, 2, 1, 4, 7}
	expected := 3 // The longest common subarray is [3,2,1]
	result := findLength(nums1, nums2)
	if result != expected {
		t.Errorf("Expected %d but got %d", expected, result)
	}
}

func TestFindLength_AllZeroSlice_Correct(t *testing.T) {
	nums1 := []int{0, 0, 0, 0, 0}
	nums2 := []int{0, 0, 0, 0, 0}
	expected := 5 // The longest common subarray is [0,0,0,0,0]
	result := findLength(nums1, nums2)
	if result != expected {
		t.Errorf("Expected %d but got %d", expected, result)
	}
}
