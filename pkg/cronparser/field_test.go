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
		{name: "bounded instant", expr: "0-55", expected: FRInitInterval},
		{name: "regular instants", expr: "*/44", expected: 44},
		{name: "bounded regular instant", expr: "1-5/0", expected: 0},
		{name: "multiple abbreviations", expr: "MON-Fri", expected: 1},
		{name: "single abbreviation", expr: "Jan", expected: 1},
	}

	for _, tc := range slashIntervalTestCases {
		t.Run(tc.name, func(t *testing.T) {
			fR := NewFieldRange(tc.expr)
			err := fR.handleSlash()
			assertSuccess(t, fR.interval, tc.expected, err)
		})
	}

	slashExprTestCases := []struct {
		name     string
		expr     string
		expected string
	}{
		{name: "every instant", expr: "*", expected: "*"},
		{name: "particular instant", expr: "0", expected: "0"},
		{name: "bounded instant", expr: "0-55", expected: "0-55"},
		{name: "regular instants", expr: "*/44", expected: "*"},
		{name: "bounded regular instant", expr: "1-5/0", expected: "1-5"},
		{name: "multiple abbreviations", expr: "MON-Fri", expected: "MON-Fri"},
		{name: "single abbreviation", expr: "Jan", expected: "Jan"},
	}

	for _, tc := range slashExprTestCases {
		t.Run(tc.name, func(t *testing.T) {
			fR := NewFieldRange(tc.expr)
			err := fR.handleSlash()
			assertSuccess(t, fR.expr, tc.expected, err)
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
		{name: "particular instant", expr: "0", bounds: HourBound, expected: "0"},
		{name: "bounded instants", expr: "0-55", bounds: DOMBound, expected: "0-55"},
		{name: "multiple abbreviations", expr: "MON-Fri", bounds: DOMBound, expected: "MON-Fri"},
		{name: "single abbreviation", expr: "Jan", expected: "Jan"},
	}

	for _, tc := range asteriskTestCases {
		t.Run(tc.name, func(t *testing.T) {
			fR := NewFieldRange(tc.expr)
			err := fR.handleAsterisk(tc.bounds)
			assertSuccess(t, fR.expr, tc.expected, err)
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
		{name: "multiple abbreviation", expr: "Mon-Fri", expected: "Mon-Fri"},
		{name: "single abbreviation", expr: "Janu", expected: "Janu-Janu"},
	}

	for _, tc := range svTestCases {
		t.Run(tc.name, func(t *testing.T) {
			fR := NewFieldRange(tc.expr)
			err := fR.handleSingleValue()
			assertSuccess(t, fR.expr, tc.expected, err)
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
	}

	for _, tc := range hyphenFailureTestCases {
		t.Run(tc.name, func(t *testing.T) {
			fR := NewFieldRange(tc.expr)
			err := fR.handleHyphen(tc.abbr)
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
		{name: "maxBound: multiple abbreviation", expr: "Mon-Fri", abbr: DOW_ABBREVIATIONS, expected: 5},
	}

	for _, tc := range hyphenMaxBoundTestCases {
		t.Run(tc.name, func(t *testing.T) {
			fR := NewFieldRange(tc.expr)
			err := fR.handleHyphen(tc.abbr)
			assertSuccess(t, fR.max, tc.expected, err)
		})
	}

	hyphenMinBoundTestCases := []struct {
		name     string
		expr     string
		abbr     map[string]string
		expected int
	}{
		{name: "particular instant", expr: "0-0", abbr: map[string]string{}, expected: 0},
		{name: "bounded instants", expr: "0-55", abbr: map[string]string{}, expected: 0},
		{name: "single abbreviation", expr: "Jan-Jan", abbr: MONTH_ABBREVIATIONS, expected: 1},
		{name: "multiple abbreviations", expr: "MON-Fri", abbr: DOW_ABBREVIATIONS, expected: 1},
	}

	for _, tc := range hyphenMinBoundTestCases {
		t.Run(tc.name, func(t *testing.T) {
			fR := NewFieldRange(tc.expr)
			err := fR.handleHyphen(tc.abbr)
			assertSuccess(t, fR.min, tc.expected, err)
		})
	}
}

func TestInvalidExprHandler(t *testing.T) {
	t.Run("invalid max interval", func(t *testing.T) {
		expr := "1-5/44"
		fR := NewFieldRange(expr)
		fR.interval = 44 ////from TestSlashHandler.bounded_regular_instant
		fR.min = 1
		fR.max = 5
		err := fR.handleInvalidExpr(DOWBound, FRInitBounds)
		expected := "invalid interval"
		assertError(t, err, expected)
	})

	t.Run("invalid min interval", func(t *testing.T) {
		expr := "*/0"
		fR := NewFieldRange(expr)
		fR.interval = 0 ////from TestSlashHandler.bounded_regular_instant
		fR.min = 0
		fR.max = 6
		err := fR.handleInvalidExpr(DOWBound, FRInitBounds)
		expected := "invalid interval"
		assertError(t, err, expected)
	})

	t.Run("invalid Max bound", func(t *testing.T) {
		expr := "0-55"
		fR := NewFieldRange(expr)
		fR.max = 55 //from hyphenMaxBoundTestCases.bounded_instants
		fR.min = 0
		err := fR.handleInvalidExpr(DOMBound, FRInitBounds)
		expected := "invalid value, out of bounds"
		assertError(t, err, expected)
	})

	t.Run("invalid Min bound", func(t *testing.T) {
		expr := "0-55"
		fR := NewFieldRange(expr)
		fR.min = 0 //from hyphenMinBoundTestCases.bounded_instants
		fR.max = 55
		err := fR.handleInvalidExpr(DOMBound, FRInitBounds)
		expected := "invalid value, out of bounds"
		assertError(t, err, expected)
	})

	t.Run("invalid bounds", func(t *testing.T) {
		expr := "4-2"
		fR := NewFieldRange(expr)
		fR.min = 4 //from hyphenMinBoundTestCases.bounded_instants
		fR.max = 2
		err := fR.handleInvalidExpr(DOWBound, FRInitBounds)
		expected := "invalid bounds"
		assertError(t, err, expected)
	})

	t.Run("valid single Value", func(t *testing.T) {
		expr := "2"
		fR := NewFieldRange(expr)
		fR.interval = 1
		fR.min = 2
		fR.max = 2
		err := fR.handleInvalidExpr(DOWBound, FRInitBounds)
		if err != nil {
			t.Fatal("error is not expected here: ", err)
		}
	})
}
