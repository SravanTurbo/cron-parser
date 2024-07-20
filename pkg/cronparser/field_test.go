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
		{name: "regular instants", expr: "*/4", expected: 4},
		{name: "bounded regular instant", expr: "1-5/44", expected: 44},
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
		{name: "regular instants", expr: "*/4", expected: "*"},
		{name: "bounded regular instant", expr: "1-5/44", expected: "1-5"},
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
	asteriskMaxBoundTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		expected int
	}{
		{name: "every instant", expr: "*", bounds: MinuteBound, expected: MinuteBound.max},
		{name: "particular instant", expr: "0", bounds: HourBound, expected: FRInitBounds},
		{name: "bounded instants", expr: "0-55", bounds: DOMBound, expected: FRInitBounds},
	}

	for _, tc := range asteriskMaxBoundTestCases {
		t.Run(tc.name, func(t *testing.T) {
			fR := NewFieldRange(tc.expr)
			err := fR.handleAsterisk(tc.bounds)
			assertSuccess(t, fR.max, tc.expected, err)
		})
	}

	asteriskMinBoundTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		expected int
	}{
		{name: "every instant", expr: "*", bounds: MinuteBound, expected: MinuteBound.min},
		{name: "particular instant", expr: "0", bounds: HourBound, expected: FRInitBounds},
		{name: "bounded instants", expr: "0-55", bounds: DOMBound, expected: FRInitBounds},
	}

	for _, tc := range asteriskMinBoundTestCases {
		t.Run(tc.name, func(t *testing.T) {
			fR := NewFieldRange(tc.expr)
			err := fR.handleAsterisk(tc.bounds)
			assertSuccess(t, fR.min, tc.expected, err)
		})
	}

	asteriskExprTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		expected string
	}{
		{name: "every instant", expr: "*", bounds: MinuteBound, expected: ""},
		{name: "particular instant", expr: "0", bounds: HourBound, expected: "0"},
		{name: "bounded instants", expr: "0-55", bounds: DOMBound, expected: "0-55"},
	}

	for _, tc := range asteriskExprTestCases {
		t.Run(tc.name, func(t *testing.T) {
			fR := NewFieldRange(tc.expr)
			err := fR.handleAsterisk(tc.bounds)
			assertSuccess(t, fR.expr, tc.expected, err)
		})
	}
}

func TestHyphenHandler(t *testing.T) {
	hyphenMaxBoundTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		expected int
	}{
		{name: "particular instant", expr: "0", bounds: HourBound, expected: FRInitBounds},
		{name: "bounded instants", expr: "0-55", bounds: DOMBound, expected: 55},
	}

	for _, tc := range hyphenMaxBoundTestCases {
		t.Run(tc.name, func(t *testing.T) {
			fR := NewFieldRange(tc.expr)
			err := fR.handleHyphen()
			assertSuccess(t, fR.max, tc.expected, err)
		})
	}

	hyphenMinBoundTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		expected int
	}{
		{name: "particular instant", expr: "0", bounds: HourBound, expected: FRInitBounds},
		{name: "bounded instants", expr: "0-55", bounds: DOMBound, expected: 0},
	}

	for _, tc := range hyphenMinBoundTestCases {
		t.Run(tc.name, func(t *testing.T) {
			fR := NewFieldRange(tc.expr)
			err := fR.handleHyphen()
			assertSuccess(t, fR.min, tc.expected, err)
		})
	}

	hyphenExprBoundTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		expected string
	}{
		{name: "particular instant", expr: "0", bounds: HourBound, expected: "0"},
		{name: "bounded instant", expr: "1-5", bounds: DOMBound, expected: ""},
	}

	for _, tc := range hyphenExprBoundTestCases {
		t.Run(tc.name, func(t *testing.T) {
			fR := NewFieldRange(tc.expr)
			err := fR.handleHyphen()
			assertSuccess(t, fR.expr, tc.expected, err)
		})
	}
}

func TestSingleValueHandler(t *testing.T) {
	svFailureTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		expected string
	}{
		{name: "bounded instant", expr: "?", bounds: MinuteBound, expected: "strconv.Atoi: parsing \"?\": invalid syntax"},
		{name: "bounded instant", expr: "2*", bounds: DOWBound, expected: "strconv.Atoi: parsing \"2*\": invalid syntax"},
	}

	for _, tc := range svFailureTestCases {
		t.Run(tc.name, func(t *testing.T) {
			fR := NewFieldRange(tc.expr)
			err := fR.handleSingleValue()
			assertError(t, err, tc.expected)
		})
	}

	svMaxBoundTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		expected int
	}{
		{name: "particular instant", expr: "0", bounds: HourBound, expected: 0},
		{name: "bounded instant", expr: "55", bounds: DOMBound, expected: 55},
	}

	for _, tc := range svMaxBoundTestCases {
		t.Run(tc.name, func(t *testing.T) {
			fR := NewFieldRange(tc.expr)
			err := fR.handleSingleValue()
			assertSuccess(t, fR.max, tc.expected, err)
		})
	}

	svMinBoundTestCases := []struct {
		name     string
		expr     string
		bounds   bound
		expected int
	}{
		{name: "particular instant", expr: "0", bounds: HourBound, expected: 0},
		{name: "bounded instant", expr: "55", bounds: DOMBound, expected: 55},
	}

	for _, tc := range svMinBoundTestCases {
		t.Run(tc.name, func(t *testing.T) {
			fR := NewFieldRange(tc.expr)
			err := fR.handleSingleValue()
			assertSuccess(t, fR.min, tc.expected, err)
		})
	}
}

func TestInvalidExprHandler(t *testing.T) {
	t.Run("invalid interval", func(t *testing.T) {
		expr := "1-5/44"
		fR := NewFieldRange(expr)
		fR.interval = 44 ////from TestSlashHandler.bounded_regular_instant
		fR.min = 1
		fR.max = 5
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
