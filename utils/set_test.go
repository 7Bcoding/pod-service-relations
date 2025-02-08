package utils

import "testing"

func TestSet(t *testing.T) {
	set1 := NewSet()
	set2 := NewSet()

	set1.Add(1)
	if set1.Empty() {
		t.Error("add error")
	} else if set1.Size() != 1 {
		t.Error("add error")
	}

	set1.Remove(1)
	if !set1.Empty() {
		t.Error("remove error")
	} else if set1.Size() == 1 {
		t.Error("remove error")
	}

	set1.Add(1)
	set1.Add(2)
	set1.Add(1)
	if set1.Empty() {
		t.Error("add error")
	} else if set1.Size() != 2 {
		t.Error("add error")
	} else if len(set1.Elements()) != 2 {
		t.Error("add error")
	}

	if set1.Contains(3) {
		t.Error("add error")
	} else if !set1.Contains(1) {
		t.Error("add error")
	}

	set1.Clear()
	if !set1.Empty() {
		t.Error("clear error")
	}

	set1.Add(1)
	set1.Add(2)
	set1.Add(3)
	set2.Add(2)
	set2.Add(4)

	if newSet := set1.Intersect(set2); newSet.Size() != 1 {
		t.Error("intersect error")
	}
	if newSet := set1.Union(set2); newSet.Size() != 4 {
		t.Error("Union error")
	}
	if newSet := set1.Subtract(set2); newSet.Size() != 2 {
		t.Error("subtract error")
	}
	if set2.IsSubset(set1) {
		t.Error("issubset error")
	}
}
