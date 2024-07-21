package cronparser

import (
	"errors"
	"strconv"
	"strings"
)

type fieldRange struct {
	expr     string
	min      int
	max      int
	interval int
}

const FRInitBounds = -1
const FRInitInterval = 1

func NewFieldRange(expr string) *fieldRange {
	return &fieldRange{expr: expr, min: FRInitBounds, max: FRInitBounds, interval: FRInitInterval}
}

func (fR *fieldRange) handleSlash() (err error) {
	exprList := strings.Split(fR.expr, "/")
	if len(exprList) == 2 {
		fR.interval, err = strconv.Atoi(exprList[1])
		if err != nil {
			return
		}

		fR.expr = exprList[0]
	}
	return
}

func (fR *fieldRange) handleAsterisk(bounds bound) (err error) {
	if fR.expr == "*" {
		fR.expr = strconv.Itoa(bounds.min) + "-" + strconv.Itoa(bounds.max)
	}
	return
}

func (fR *fieldRange) handleSingleValue() (err error) {
	exprList := strings.Split(fR.expr, "-")
	if len(exprList) == 1 {
		fR.expr = exprList[0] + "-" + exprList[0]
	}
	return
}

func (fR *fieldRange) handleHyphen(abbreviationMap map[string]string) (err error) {
	exprList := strings.Split(fR.expr, "-")
	if len(exprList) == 2 {
		fR.min, err = computeValue(exprList[0], abbreviationMap)
		if err != nil {
			return
		}

		fR.max, err = computeValue(exprList[1], abbreviationMap)
		if err != nil {
			return
		}
	}
	return
}

func (fR fieldRange) handleInvalidExpr(bounds bound, initBounds int) error {
	if fR.min == initBounds || fR.max == initBounds {
		return errors.New("invalid cron field")
	}

	if fR.min < bounds.min || fR.max > bounds.max {
		return errors.New("invalid value, out of bounds")
	}

	if fR.min > fR.max {
		return errors.New("invalid bounds")
	}

	if fR.interval != 1 {
		_range := bounds.max - bounds.min + 1
		if fR.interval > _range || fR.interval == 0 {
			return errors.New("invalid interval")
		}
	}

	return nil
}

func computeValue(expr string, abbrMap map[string]string) (val int, err error) {
	abbrVal, ok := abbrMap[strings.ToUpper(expr)]
	if !ok {
		abbrVal = expr
	}

	val, err = strconv.Atoi(abbrVal)
	return
}
