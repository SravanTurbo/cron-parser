package cronparser

import "testing"

func TestString(t *testing.T) {
	schedule := &Schedule{minute: []int{30}, hour: []int{4}, dom: []int{1}, month: []int{1}, dow: []int{0}, cmd: "cmd"}
	got := schedule.String()
	expected := "minute\t\t30\nhour\t\t4\nday of month\t1\nmonth\t\t1\nday of week\t0\ncommand\t\tcmd"

	if got != expected {
		t.Errorf("expected %s, but got %s", expected, got)
	}
}
