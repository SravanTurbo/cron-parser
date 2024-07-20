package cronparser

import (
	"reflect"
	"testing"
)

func TestUtitlity(t *testing.T) {
	t.Run("build int list", func(t *testing.T) {
		got := buildIntList(0, 6, 2)
		expected := []int{0, 2, 4, 6}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected %v but got %v", expected, got)
		}
	})

	t.Run("join ints list ", func(t *testing.T) {
		got := intsJoin([]int{1, 2, 3, 4}, ",")
		expected := "1,2,3,4"

		if got != expected {
			t.Errorf("expected %s, but got %s", expected, got)
		}

	})
}
