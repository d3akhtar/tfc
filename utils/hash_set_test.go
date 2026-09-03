package utils_test

import (
	"github.com/d3akhtar/tfc/utils"
	"testing"
)

func TestHashSet(t *testing.T) {
	set := utils.NewHashSet[int]()

	expectLength(t, &set, 0)

	set.Add(5)
	expectContains(t, &set, 5)
	expectLength(t, &set, 1)

	set.Add(7)
	expectContains(t, &set, 5)
	expectContains(t, &set, 7)
	expectLength(t, &set, 2)

	set.Remove(4)
	expectContains(t, &set, 5)
	expectContains(t, &set, 7)
	expectLength(t, &set, 2)

	set.Remove(5)
	expectNotContains(t, &set, 5)
	expectContains(t, &set, 7)
	expectLength(t, &set, 1)
}

func expectLength(t *testing.T, set *utils.HashSet[int], expectedLength int) {
	if set.Length() != expectedLength {
		t.Fatalf("set length expected=%v, got=%v", expectedLength, set.Length())
	}
}

func expectContains(t *testing.T, set *utils.HashSet[int], item int) {
	if !set.Contains(item) {
		t.Fatalf("set expected to contain %v", item)
	}
}

func expectNotContains(t *testing.T, set *utils.HashSet[int], item int) {
	if set.Contains(item) {
		t.Fatalf("set shouldn't contain %v", item)
	}
}
