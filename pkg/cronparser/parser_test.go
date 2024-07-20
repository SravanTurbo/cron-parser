package cronparser

import (
	"reflect"
	"testing"
)

func TestValidator(t *testing.T) {
	failureTestCases := []struct {
		name     string
		cronExpr string
		expected string
	}{
		{name: "invalid number of fields", cronExpr: "*/15 0 1,15 1-5 /usr/bin/find", expected: "Invalid cron expression, check README"},
		{name: "invalid special character", cronExpr: "*/15 0 ? 2 1-5 /usr/bin/find", expected: "Invalid cron field"},
		{name: "extra space in separator", cronExpr: "*/15 0 1  2 1-5 /usr/bin/find", expected: "Invalid cron expression, check README"},
	}

	for _, tc := range failureTestCases {
		_, err := validate(tc.cronExpr)
		assertError(t, err, tc.expected)
	}
}

func TestParseField(t *testing.T) {
	successTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		expected []int
	}{
		{name: "one instant", expr: "2", bounds: DOWBound, expected: []int{2}},
		{name: "every instant", expr: "*", bounds: MinuteBound, expected: buildIntList(0, 59, 1)},
		{name: "regular instants", expr: "*/4", bounds: HourBound, expected: buildIntList(0, 23, 4)},
		{name: "bounded instants", expr: "1-15", bounds: DOMBound, expected: buildIntList(1, 15, 1)},
		{name: "regular instants", expr: "1, 4, 12", bounds: MonthBound, expected: []int{1, 4, 12}},
		{name: "bound regular instants", expr: "1-6/2", bounds: DOWBound, expected: []int{1, 3, 5}},
	}

	failureTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		expected string
	}{
		{name: "every instant", expr: "2*", bounds: MinuteBound, expected: "invalid cron expression"},
		{name: "regular instants", expr: "*/24", bounds: HourBound, expected: "interval out of bounds"},
		{name: "bounded instants", expr: "1-32", bounds: DOMBound, expected: "invalid maximum value of range"},
		{name: "regular instants", expr: "1, 4, 13", bounds: MonthBound, expected: "invalid value provided"},
		{name: "bound regular instants", expr: "1-12/2", bounds: DOWBound, expected: "invalid maximum value of range"},
		{name: "bound regular instants", expr: "1-4/7", bounds: DOWBound, expected: "interval out of bounds"},
	}

	for _, tc := range successTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseField(tc.expr, tc.bounds)
			if err != nil {
				t.Fatal("error is not expected here: ", err)
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("expected %v, but got %v", tc.expected, got)
			}
		})
	}

	for _, tc := range failureTestCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseField(tc.expr, tc.bounds)
			assertError(t, err, tc.expected)
		})
	}
}

func TestParse(t *testing.T) {
	cronExpr := "*/15 0 1,15 * 1-5 /usr/bin/find"
	got, err := parse(cronExpr)
	expected := "minute\t\t0 15 30 45\nhour\t\t0\nday of month\t1 15\nmonth\t\t0 1 2 3 4 5 6 7 8 9 10 11 12\nday of week\t1 2 3 4 5\ncommand\t\t/usr/bin/find"
	if err != nil {
		t.Fatal("error not expected here: ", err)
	}

	if got.String() != expected {
		t.Errorf("expected \n%s, but got \n%s", expected, got)
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
