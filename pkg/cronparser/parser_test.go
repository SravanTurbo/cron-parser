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
		{name: "extra space in separator", cronExpr: "*/15 0 1, 2 2 1-5 /usr/bin/find", expected: "Validation Error: invalid number of cron fields"},
		{name: "invalid number of fields", cronExpr: "*/15 0 1,15 1-5 /usr/bin/find", expected: "Validation Error: invalid number of cron fields"},
		{name: "invalid special character", cronExpr: "*/15 0 ? 1 1-5 /usr/bin/find", expected: "Validation Error: invalid time field"},
		{name: "invalid special character", cronExpr: "*/15 0 L 2 1-5 /usr/bin/find", expected: "Validation Error: invalid time field"},
	}

	for _, tc := range failureTestCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := validate(tc.cronExpr)
			assertError(t, err, tc.expected)
		})
	}

	t.Run("valid case with abbr", func(t *testing.T) {
		cronExpr := "*/15 0 1 jan Mon /usr/bin/find"
		_, err := validate(cronExpr)
		if err != nil {
			t.Fatal("error is not expected here, but got one: ", err)
		}
	})
}

func TestComputeField(t *testing.T) {
	failureTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		abbr     map[string]string
		expected string
	}{
		{name: "invalid every instant", expr: "2*", bounds: MinuteBound, abbr: map[string]string{}, expected: "strconv.Atoi: parsing \"2*\": invalid syntax"},
		{name: "invalid one instant", expr: "60", bounds: MinuteBound, abbr: map[string]string{}, expected: "invalid value, out of bounds"},
		{name: "invalid regular instants", expr: "*/26", bounds: HourBound, abbr: map[string]string{}, expected: "invalid interval"},
		{name: "invalid bounded instants", expr: "1-32", bounds: DOMBound, abbr: map[string]string{}, expected: "invalid value, out of bounds"},
		{name: "invalid bounded regular instants", expr: "1-12/2", bounds: DOWBound, abbr: DOW_ABBREVIATIONS, expected: "invalid value, out of bounds"},
		{name: "invalid bounded regular instants 2", expr: "1-4/8", bounds: DOWBound, abbr: DOW_ABBREVIATIONS, expected: "invalid interval"},
		{name: "invalid bounded regular instants 2", expr: "Dec-Jan", bounds: MonthBound, abbr: MONTH_ABBREVIATIONS, expected: "invalid bounds"},
	}

	for _, tc := range failureTestCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := computeField(tc.expr, tc.bounds, tc.abbr)
			assertError(t, err, tc.expected)
		})
	}

	successTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		abbr     map[string]string
		expected []int
	}{
		{name: "one instant", expr: "2", bounds: DOWBound, abbr: DOW_ABBREVIATIONS, expected: []int{2}},
		{name: "every instant", expr: "*", bounds: MinuteBound, abbr: map[string]string{}, expected: buildIntList(MinuteBound.min, MinuteBound.max, 1)},
		{name: "regular instants", expr: "*/4", bounds: HourBound, abbr: map[string]string{}, expected: buildIntList(HourBound.min, HourBound.max, 4)},
		{name: "bounded instants", expr: "1-15", bounds: DOMBound, abbr: map[string]string{}, expected: buildIntList(1, 15, 1)},
		{name: "bound regular instants", expr: "1-4/7", bounds: DOWBound, abbr: DOW_ABBREVIATIONS, expected: buildIntList(1, 1, 1)},
		{name: "one abbr instant", expr: "jul", bounds: MonthBound, abbr: MONTH_ABBREVIATIONS, expected: []int{7}},
	}

	for _, tc := range successTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := computeField(tc.expr, tc.bounds, tc.abbr)
			assertSuccess(t, got, tc.expected, err)
		})
	}
}

func TestParseField(t *testing.T) {
	failureTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		abbr     map[string]string
		expected string
	}{
		{name: "out of bounds", expr: "1,4,13", bounds: MonthBound, abbr: MONTH_ABBREVIATIONS, expected: "Parsing Error: invalid value, out of bounds"},
		{name: "chars in expr", expr: "1,a,13", bounds: MonthBound, abbr: MONTH_ABBREVIATIONS, expected: "Parsing Error: strconv.Atoi: parsing \"a\": invalid syntax"},
	}

	for _, tc := range failureTestCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseField(tc.expr, tc.bounds, tc.abbr)
			assertError(t, err, tc.expected)
		})
	}

	successTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		abbr     map[string]string
		expected []int
	}{
		{name: "special chars in expr", expr: "*,4", bounds: MonthBound, abbr: MONTH_ABBREVIATIONS, expected: buildIntList(1, 12, 1)}, //***
		{name: "particular instants with interval", expr: "*/15,4", bounds: DOMBound, abbr: map[string]string{}, expected: []int{1, 4, 16, 31}},
		{name: "particular instants", expr: "1,4,12", bounds: MonthBound, abbr: map[string]string{}, expected: []int{1, 4, 12}},
		{name: "particular instants with single instant", expr: "2,SEP", bounds: MonthBound, abbr: MONTH_ABBREVIATIONS, expected: []int{2, 9}},
		{name: "particular instants with bounded interval", expr: "Mon-Fri,Sun", bounds: DOWBound, abbr: DOW_ABBREVIATIONS, expected: []int{0, 1, 2, 3, 4, 5}},
		{name: "unique values", expr: "Mon-Fri,THU", bounds: DOWBound, abbr: DOW_ABBREVIATIONS, expected: []int{1, 2, 3, 4, 5}},
	}

	for _, tc := range successTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseField(tc.expr, tc.bounds, tc.abbr)
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
