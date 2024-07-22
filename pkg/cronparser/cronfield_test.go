package cronparser

import (
	"testing"
)

func TestSlashHandler(t *testing.T) {
	slashIntervalTestCases := []struct {
		name     string
		expr     string
		expected int
	}{
		{name: "every instant", expr: "*", expected: FRInitInterval},
		{name: "particular instant", expr: "0", expected: FRInitInterval},
		{name: "bounded instants", expr: "0-55", expected: FRInitInterval},
		{name: "regular instants", expr: "*/44", expected: 44},
		{name: "bounded regular instants", expr: "1-5/0", expected: 0},
		{name: "multiple abbreviations", expr: "MON-Fri", expected: 1},
		{name: "single abbreviation", expr: "Jan", expected: 1},
	}

	for _, tc := range slashIntervalTestCases {
		t.Run(tc.name, func(t *testing.T) {
			cf := NewCronField(tc.expr)
			err := cf.handleSlash()
			assertSuccess(t, cf.interval, tc.expected, err)
		})
	}

	slashExprTestCases := []struct {
		name     string
		expr     string
		expected string
	}{
		{name: "every instant", expr: "*", expected: "*"},
		{name: "particular instant", expr: "0", expected: "0"},
		{name: "bounded instants", expr: "0-55", expected: "0-55"},
		{name: "regular instants", expr: "*/44", expected: "*"},
		{name: "bounded regular instants", expr: "1-5/0", expected: "1-5"},
		{name: "multiple abbreviations", expr: "MON-Fri", expected: "MON-Fri"},
		{name: "single abbreviation", expr: "Jan", expected: "Jan"},
	}

	for _, tc := range slashExprTestCases {
		t.Run(tc.name, func(t *testing.T) {
			cf := NewCronField(tc.expr)
			err := cf.handleSlash()
			assertSuccess(t, cf.expr, tc.expected, err)
		})
	}
}

func TestAsteriskHandler(t *testing.T) {
	asteriskTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		expected string
	}{
		{name: "every instant", expr: "*", bounds: MinuteBound, expected: "0-59"},
		{name: "particular instant", expr: "?", bounds: HourBound, expected: "?"},
		{name: "bounded instants", expr: "0-55", bounds: DOMBound, expected: "0-55"},
		{name: "regular instants", expr: "*/44", expected: "*/44"},
		{name: "bounded regular instants", expr: "1-5/0", expected: "1-5/0"},
		{name: "multiple abbreviations", expr: "MON-Fri", bounds: DOMBound, expected: "MON-Fri"},
		{name: "single abbreviation", expr: "Jan", expected: "Jan"},
	}

	for _, tc := range asteriskTestCases {
		t.Run(tc.name, func(t *testing.T) {
			cf := NewCronField(tc.expr)
			err := cf.handleAsterisk(tc.bounds)
			assertSuccess(t, cf.expr, tc.expected, err)
		})
	}
}

func TestSingleValueHandler(t *testing.T) {
	svTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		expected string
	}{
		{name: "particular instant", expr: "0", bounds: HourBound, expected: "0-0"},
		{name: "invalid particular instant", expr: "55", bounds: DOMBound, expected: "55-55"},
		{name: "invalid special char", expr: "?", bounds: MinuteBound, expected: "?-?"},
		{name: "regular instants", expr: "*/44", expected: "*/44-*/44"},
		{name: "bounded regular instants", expr: "1-5/0", expected: "1-5/0"},
		{name: "multiple abbreviation", expr: "Mon-Fri", expected: "Mon-Fri"},
		{name: "single abbreviation", expr: "Janu", expected: "Janu-Janu"},
	}

	for _, tc := range svTestCases {
		t.Run(tc.name, func(t *testing.T) {
			cf := NewCronField(tc.expr)
			err := cf.handleSingleValue()
			assertSuccess(t, cf.expr, tc.expected, err)
		})
	}
}

func TestHyphenHandler(t *testing.T) {
	hyphenFailureTestCases := []struct {
		name     string
		expr     string
		abbr     map[string]string
		expected string
	}{
		{name: "invalid special char", expr: "?-?", abbr: map[string]string{}, expected: "strconv.Atoi: parsing \"?\": invalid syntax"},
		{name: "invalid abbreviation", expr: "Janu-Janu", abbr: MONTH_ABBREVIATIONS, expected: "strconv.Atoi: parsing \"Janu\": invalid syntax"},
		{name: "with interval", expr: "1-5/0", abbr: DOW_ABBREVIATIONS, expected: "strconv.Atoi: parsing \"5/0\": invalid syntax"},
	}

	for _, tc := range hyphenFailureTestCases {
		t.Run(tc.name, func(t *testing.T) {
			cf := NewCronField(tc.expr)
			err := cf.handleHyphen(tc.abbr)
			assertError(t, err, tc.expected)
		})
	}

	hyphenMaxBoundTestCases := []struct {
		name     string
		expr     string
		abbr     map[string]string
		expected int
	}{
		{name: "maxBound: particular instant", expr: "0-0", abbr: map[string]string{}, expected: 0},
		{name: "maxBound: bounded instants", expr: "0-55", abbr: map[string]string{}, expected: 55},
		{name: "maxBound: single abbreviation", expr: "Jan-Jan", abbr: MONTH_ABBREVIATIONS, expected: 1},
		{name: "maxBound: multiple abbreviation", expr: "Fri-Mon", abbr: DOW_ABBREVIATIONS, expected: 1},
	}

	for _, tc := range hyphenMaxBoundTestCases {
		t.Run(tc.name, func(t *testing.T) {
			cf := NewCronField(tc.expr)
			err := cf.handleHyphen(tc.abbr)
			assertSuccess(t, cf.max, tc.expected, err)
		})
	}

	hyphenMinBoundTestCases := []struct {
		name     string
		expr     string
		abbr     map[string]string
		expected int
	}{
		{name: "minBound: particular instant", expr: "0-0", abbr: map[string]string{}, expected: 0},
		{name: "minBound: bounded instants", expr: "0-55", abbr: map[string]string{}, expected: 0},
		{name: "minBound: single abbreviation", expr: "Jan-Jan", abbr: MONTH_ABBREVIATIONS, expected: 1},
		{name: "minBound: multiple abbreviations", expr: "Fri-MoN", abbr: DOW_ABBREVIATIONS, expected: 5},
	}

	for _, tc := range hyphenMinBoundTestCases {
		t.Run(tc.name, func(t *testing.T) {
			cf := NewCronField(tc.expr)
			err := cf.handleHyphen(tc.abbr)
			assertSuccess(t, cf.min, tc.expected, err)
		})
	}
}

func TestInvalidExprHandler(t *testing.T) {
	t.Run("invalid max interval", func(t *testing.T) {
		expr := "1-5/44"
		cf := NewCronField(expr)
		cf.interval = 44 ////from TestSlashHandler.bounded_regular_instant
		cf.min = 1
		cf.max = 5
		err := cf.handleInvalidExpr(DOWBound, FRInitBounds)
		expected := "invalid interval"
		assertError(t, err, expected)
	})

	t.Run("invalid min interval", func(t *testing.T) {
		expr := "*/0"
		cf := NewCronField(expr)
		cf.interval = 0 ////from TestSlashHandler.bounded_regular_instant
		cf.min = 0
		cf.max = 6
		err := cf.handleInvalidExpr(DOWBound, FRInitBounds)
		expected := "invalid interval"
		assertError(t, err, expected)
	})

	t.Run("invalid Max bound", func(t *testing.T) {
		expr := "0-55"
		cf := NewCronField(expr)
		cf.max = 55 //from hyphenMaxBoundTestCases.bounded_instants
		cf.min = 0
		err := cf.handleInvalidExpr(DOMBound, FRInitBounds)
		expected := "invalid value, out of bounds"
		assertError(t, err, expected)
	})

	t.Run("invalid Min bound", func(t *testing.T) {
		expr := "0-55"
		cf := NewCronField(expr)
		cf.min = 0 //from hyphenMinBoundTestCases.bounded_instants
		cf.max = 55
		err := cf.handleInvalidExpr(DOMBound, FRInitBounds)
		expected := "invalid value, out of bounds"
		assertError(t, err, expected)
	})

	t.Run("invalid abbr bounds", func(t *testing.T) {
		expr := "Fri-Mon"
		cf := NewCronField(expr) //from hyphenMinBoundTestCases.bounded_instants
		cf.interval = 1
		cf.min = 5
		cf.max = 1
		err := cf.handleInvalidExpr(DOWBound, FRInitBounds)
		expected := "invalid bounds"
		assertError(t, err, expected)
	})

	t.Run("valid single Value", func(t *testing.T) {
		expr := "2"
		cf := NewCronField(expr)
		cf.interval = 1
		cf.min = 2
		cf.max = 2
		err := cf.handleInvalidExpr(DOWBound, FRInitBounds)
		if err != nil {
			t.Fatal("error is not expected here: ", err)
		}
	})
}

func TestBoundFormat(t *testing.T) {
	failureTestCases := []struct {
		name     string
		val      string
		abbr     map[string]string
		expected string
	}{
		{name: "abbreviation", val: "january", abbr: MONTH_ABBREVIATIONS, expected: "strconv.Atoi: parsing \"january\": invalid syntax"},
		{name: "abbreviation", val: "L", abbr: map[string]string{}, expected: "strconv.Atoi: parsing \"L\": invalid syntax"},
	}

	for _, tc := range failureTestCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := formatBound(tc.val, tc.abbr)
			assertError(t, err, tc.expected)
		})
	}

	successTestCases := []struct {
		name     string
		val      string
		abbr     map[string]string
		expected int
	}{
		{name: "abbreviation", val: "jan", abbr: MONTH_ABBREVIATIONS, expected: 1},
		{name: "abbreviation", val: "2134", abbr: map[string]string{}, expected: 2134},
	}

	for _, tc := range successTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formatBound(tc.val, tc.abbr)
			assertSuccess(t, got, tc.expected, err)
		})
	}
}
