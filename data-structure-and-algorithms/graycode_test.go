package elotuschallenge

import (
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestInput_0_Correct_GrayCode(t *testing.T) {
	n := 0
	expected := []int{0}
	result := grayCode(n)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("For input %d, expected %v but got %v", n, expected, result)
	}
}

func TestInput_1_Correct_GrayCode(t *testing.T) {
	n := 1
	expected := []int{0b00, 0b01}
	result := grayCode(n)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("For input %d, expected %v but got %v", n, expected, result)
	}
}

func TestInput_2_Correct_GrayCode(t *testing.T) {
	n := 2
	first := []int{0b00, 0b01, 0b11, 0b10}
	second := []int{0b00, 0b10, 0b01, 0b11}
	result := grayCode(n)
	if !reflect.DeepEqual(result, first) && !reflect.DeepEqual(result, second) {
		t.Errorf("For input %d, expected %v or %v but got %v", n, first, second, result)
	}
}

func TestInput_10_Correct_GrayCode(t *testing.T) {
	n := 10
	result := grayCode(n)
	value, dup := FindFirstDuplicate(grayCode(n))
	t.Logf("Gray code for n=10: %v", result)
	if dup {
		t.Errorf("For input %d, found duplicate value %d in gray code sequence", n, value)
	} else {
		t.Logf("No duplicates found in gray code sequence for n=%d", n)
	}

}
