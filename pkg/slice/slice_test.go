package slice

import "testing"

func TestContains(t *testing.T) {
	sl0 := []string{"foo", "bar"}
	sl1 := []int{4, 2}

	if !ContainsString(sl0, "foo") {
		t.Errorf("expected %v to be in slice", "foo")
	}
	if ContainsInt(sl1, 1) {
		t.Errorf("expected %v to not be in slice", 1)
	}
}
