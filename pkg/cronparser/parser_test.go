package cronparser

import (
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
		{name: "extra space in separator", cronExpr: "*/15 0 1, 2 2 1-5 /usr/bin/find", expected: "Invalid cron expression, check README"},
	}

	for _, tc := range failureTestCases {
		_, err := validate(tc.cronExpr)
		assertError(t, err, tc.expected)
	}
}

func TestComputeField(t *testing.T) {
	failureTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		expected string
	}{
		{name: "invalid every instant", expr: "2*", bounds: MinuteBound, expected: "strconv.Atoi: parsing \"2*\": invalid syntax"},
		{name: "every instant", expr: "60", bounds: MinuteBound, expected: "invalid value, out of bounds"},
		{name: "invalid regular instants", expr: "*/26", bounds: HourBound, expected: "invalid interval"},
		{name: "invalid bounded instants", expr: "1-32", bounds: DOMBound, expected: "invalid value, out of bounds"},
		{name: "invalid bounded regular instants", expr: "1-12/2", bounds: DOWBound, expected: "invalid value, out of bounds"},
		{name: "invalid bounded regular instants 2", expr: "1-4/8", bounds: DOWBound, expected: "invalid interval"},
	}

	for _, tc := range failureTestCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := computeField(tc.expr, tc.bounds)
			assertError(t, err, tc.expected)
		})
	}

	successTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		expected []int
	}{
		{name: "one instant", expr: "2", bounds: DOWBound, expected: []int{2}},
		{name: "every instant", expr: "*", bounds: MinuteBound, expected: buildIntList(MinuteBound.min, MinuteBound.max, 1)},
		{name: "regular instants", expr: "*/4", bounds: HourBound, expected: buildIntList(HourBound.min, HourBound.max, 4)},
		{name: "bounded instants", expr: "1-15", bounds: DOMBound, expected: buildIntList(1, 15, 1)},
		{name: "bound regular instants", expr: "1-4/7", bounds: DOWBound, expected: buildIntList(1, 1, 1)},
	}

	for _, tc := range successTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := computeField(tc.expr, tc.bounds)
			assertSuccess(t, got, tc.expected, err)
		})
	}
}

func TestParseField(t *testing.T) {
	failureTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		expected string
	}{
		{name: "out of bounds", expr: "1,4,13", bounds: MonthBound, expected: "invalid cron: invalid value, out of bounds"},
		{name: "chars in expr", expr: "1,a,13", bounds: MonthBound, expected: "invalid cron: strconv.Atoi: parsing \"a\": invalid syntax"},
	}

	for _, tc := range failureTestCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseField(tc.expr, tc.bounds)
			assertError(t, err, tc.expected)
		})
	}

	successTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		expected []int
	}{
		{name: "special chars in expr", expr: "*,4", bounds: MonthBound, expected: buildIntList(1, 12, 1)}, //***
		{name: "special chars in expr", expr: "*/2,4", bounds: DOMBound, expected: append(buildIntList(1, 31, 2), []int{4}...)},
		{name: "particular instants", expr: "1,4,12", bounds: MonthBound, expected: []int{1, 4, 12}},
	}

	for _, tc := range successTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseField(tc.expr, tc.bounds)
			assertSuccess(t, got, tc.expected, err)
		})
	}
}

func TestParse(t *testing.T) {
	cronExpr := "*/15 0 1,15 * 1-5 /usr/bin/find"
	got, err := Parse(cronExpr)
	expected := "minute\t\t0 15 30 45\nhour\t\t0\nday of month\t1 15\nmonth\t\t1 2 3 4 5 6 7 8 9 10 11 12\nday of week\t1 2 3 4 5\ncommand\t\t/usr/bin/find"
	if err != nil {
		t.Fatal("error not expected here: ", err)
	}

	if got.String() != expected {
		t.Errorf("expected \n%s, but got \n%s", expected, got)
	}
}
