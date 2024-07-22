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

func TestNonCommaHandler(t *testing.T) {
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
		{name: "invalid special char", expr: "L", bounds: MonthBound, abbr: MONTH_ABBREVIATIONS, expected: "strconv.Atoi: parsing \"L\": invalid syntax"},
	}

	for _, tc := range failureTestCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := handleNonComma(tc.expr, tc.bounds, tc.abbr)
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
		{name: "one abbr instant", expr: "MOn", bounds: DOWBound, abbr: DOW_ABBREVIATIONS, expected: []int{1}},
	}

	for _, tc := range successTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := handleNonComma(tc.expr, tc.bounds, tc.abbr)
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
		{name: "FC: out of bounds", expr: "1,4,13", bounds: MonthBound, abbr: MONTH_ABBREVIATIONS, expected: "Parsing Error: invalid value, out of bounds"},
		{name: "FC: chars in expr", expr: "1,a,13", bounds: MonthBound, abbr: MONTH_ABBREVIATIONS, expected: "Parsing Error: strconv.Atoi: parsing \"a\": invalid syntax"},
		{name: "FC: invalid special char", expr: "L", bounds: MonthBound, abbr: MONTH_ABBREVIATIONS, expected: "Parsing Error: strconv.Atoi: parsing \"L\": invalid syntax"},
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
		{name: "SC: special chars in expr", expr: "*,4", bounds: MonthBound, abbr: MONTH_ABBREVIATIONS, expected: buildIntList(1, 12, 1)}, //***
		{name: "SC: particular instants with interval", expr: "*/15,4", bounds: DOMBound, abbr: map[string]string{}, expected: []int{1, 4, 16, 31}},
		{name: "SC: particular instants", expr: "1,4,12", bounds: MonthBound, abbr: map[string]string{}, expected: []int{1, 4, 12}},
		{name: "SC: particular instants with single instant", expr: "2,SEP", bounds: MonthBound, abbr: MONTH_ABBREVIATIONS, expected: []int{2, 9}},
		{name: "SC: particular instants with bounded interval", expr: "Mon-Fri,Sun", bounds: DOWBound, abbr: DOW_ABBREVIATIONS, expected: []int{0, 1, 2, 3, 4, 5}},
		{name: "SC: unique values", expr: "Mon-Fri,THU", bounds: DOWBound, abbr: DOW_ABBREVIATIONS, expected: []int{1, 2, 3, 4, 5}},
	}

	for _, tc := range successTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseField(tc.expr, tc.bounds, tc.abbr)
			assertSuccess(t, got, tc.expected, err)
		})
	}
}

func TestParse(t *testing.T) {
	parseFailureTestCases := []struct {
		name     string
		cronExpr string
		expected string
	}{
		{name: "FC: invalid number of cron fields", cronExpr: "* * * * /usr/bin/find", expected: "Validation Error: invalid number of cron fields"},
		{name: "FC: space in comma separated field", cronExpr: "* * * 1, 12 * /usr/bin/find", expected: "Validation Error: invalid number of cron fields"},
		{name: "FC: invalid cron field", cronExpr: "abc * * * * /usr/bin/find", expected: "Validation Error: invalid time field"},
		{name: "FC: invalid special char", cronExpr: "* * * * ? /usr/bin/find", expected: "Validation Error: invalid time field"},

		{name: "FC: invalid minute cron field", cronExpr: "2* * * * * /usr/bin/find", expected: "Parsing Error: strconv.Atoi: parsing \"2*\": invalid syntax"},
		{name: "FC: invalid month cron field", cronExpr: "* * * * Monday /usr/bin/find", expected: "Parsing Error: strconv.Atoi: parsing \"MONDAY\": invalid syntax"},

		{name: "FC: invalid regular instants 1", cronExpr: "*/100 * * * * /usr/bin/find", expected: "Parsing Error: invalid interval"},
		{name: "FC: invalid regular instants 2", cronExpr: "* * * * */0 /usr/bin/find", expected: "Parsing Error: invalid interval"},

		{name: "FC: invalid bound val", cronExpr: "* * ?-? * * /usr/bin/find", expected: "Parsing Error: strconv.Atoi: parsing \"?\": invalid syntax"},
		{name: "FC: invalid bounds", cronExpr: "* * 0-32 * * /usr/bin/find", expected: "Parsing Error: invalid value, out of bounds"},

		{name: "FC: invalid abbr interval", cronExpr: "* * * * Mon-Fri/8 /usr/bin/find", expected: "Parsing Error: invalid interval"},
		{name: "FC: invalid comma separated field", cronExpr: "* * * 1,march * /usr/bin/find", expected: "Parsing Error: strconv.Atoi: parsing \"MARCH\": invalid syntax"},
	}

	for _, tc := range parseFailureTestCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.cronExpr)
			assertError(t, err, tc.expected)
		})
	}

	parseSuccessTestCases := []struct {
		name     string
		cronExpr string
		expected string
	}{
		{name: "SC: every minute", cronExpr: "* * * * * cmd", expected: "minute\t\t0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55 56 57 58 59\nhour\t\t0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23\nday of month\t1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31\nmonth\t\t1 2 3 4 5 6 7 8 9 10 11 12\nday of week\t0 1 2 3 4 5 6\ncommand\t\tcmd"},
		{name: "SC: 4:30UTC on 1st day of every quarter", cronExpr: "30 4 1 */3 * cmd", expected: "minute\t\t30\nhour\t\t4\nday of month\t1\nmonth\t\t1 4 7 10\nday of week\t0 1 2 3 4 5 6\ncommand\t\tcmd"},
		{name: "SC: Christmas and New Year Greetings", cronExpr: "59 23 24,31 12 * greetings", expected: "minute\t\t59\nhour\t\t23\nday of month\t24 31\nmonth\t\t12\nday of week\t0 1 2 3 4 5 6\ncommand\t\tgreetings"},

		{name: "SC: days via hyphen", cronExpr: "0 0 1-5 1 0 cmd", expected: "minute\t\t0\nhour\t\t0\nday of month\t1 2 3 4 5\nmonth\t\t1\nday of week\t0\ncommand\t\tcmd"},

		{name: "SC: assigment example", cronExpr: "*/15 0 1,15 * 1-5 /usr/bin/find", expected: "minute\t\t0 15 30 45\nhour\t\t0\nday of month\t1 15\nmonth\t\t1 2 3 4 5 6 7 8 9 10 11 12\nday of week\t1 2 3 4 5\ncommand\t\t/usr/bin/find"},
		{name: "SC: assigment example abbr", cronExpr: "*/15 0 1,15 Jan-Dec Mon-Fri /usr/bin/find", expected: "minute\t\t0 15 30 45\nhour\t\t0\nday of month\t1 15\nmonth\t\t1 2 3 4 5 6 7 8 9 10 11 12\nday of week\t1 2 3 4 5\ncommand\t\t/usr/bin/find"},

		{name: "SC: combination of abbr and int", cronExpr: "30 4 */15,4 2,SEP */2,Mon,5 cmd", expected: "minute\t\t30\nhour\t\t4\nday of month\t1 4 16 31\nmonth\t\t2 9\nday of week\t0 1 2 4 5 6\ncommand\t\tcmd"},
	}

	for _, tc := range parseSuccessTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Parse(tc.cronExpr)
			assertSuccess(t, got.String(), tc.expected, err)
		})
	}
}
