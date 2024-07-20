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

func assertSuccess(t testing.TB, got, expected interface{}, err error) {
	t.Helper()
	if err != nil {
		t.Fatal("error is not expected here: ", err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v, but got %v", expected, got)
	}
}

func assertError(t testing.TB, got error, expected string) {
	t.Helper()
	if got == nil {
		t.Fatal("expected an error but didn't get one")
	}

	if got.Error() != expected {
		t.Errorf("expected %q, but got %q", expected, got)
	}
}
