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
		fR.min = bounds.min
		fR.max = bounds.max
		fR.expr = ""
	}

	return nil
}

func (fR *fieldRange) handleHyphen() (err error) {
	exprList := strings.Split(fR.expr, "-")
	if len(exprList) == 2 {
		fR.min, err = strconv.Atoi(exprList[0])
		if err != nil {
			return
		}

		fR.max, err = strconv.Atoi(exprList[1])
		if err != nil {
			return
		}

		fR.expr = ""
	}
	return
}

func (fR *fieldRange) handleSingleValue() (err error) {
	if fR.expr != "" {
		fR.min, err = strconv.Atoi(fR.expr)
		if err != nil {
			return err
		}

		if fR.min >= 0 {
			fR.max = fR.min
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
		if fR.interval > _range {
			return errors.New("invalid interval")
		}
	}

	return nil
}
