package cronparser

import (
	"errors"
	"strconv"
	"strings"
)

type cronField struct {
	expr      string
	min       int
	max       int
	interval  int
	valueList []int
}

const FRInitBounds = -1  //Initial Bounds of a Cron Field Range
const FRInitInterval = 1 //Initial Interval for a Cron Field Range

func NewCronField(expr string) *cronField {
	return &cronField{expr: expr, min: FRInitBounds, max: FRInitBounds, interval: FRInitInterval}
}

func (cf *cronField) handleSlash() (err error) {
	exprList := strings.Split(cf.expr, "/")
	if len(exprList) == 2 {
		cf.interval, err = strconv.Atoi(exprList[1])
		if err != nil {
			return
		}

		cf.expr = exprList[0]
	}
	return
}

func (cf *cronField) handleAsterisk(bounds bound) (err error) {
	if cf.expr == "*" {
		cf.expr = strconv.Itoa(bounds.min) + "-" + strconv.Itoa(bounds.max)
	}
	return
}

func (cf *cronField) handleSingleValue() (err error) {
	exprList := strings.Split(cf.expr, "-")
	if len(exprList) == 1 {
		cf.expr = exprList[0] + "-" + exprList[0]
	}
	return
}

func (cf *cronField) handleHyphen(abbreviationMap map[string]string) (err error) {
	exprList := strings.Split(cf.expr, "-")
	if len(exprList) == 2 {
		cf.min, err = formatBound(exprList[0], abbreviationMap)
		if err != nil {
			return
		}

		cf.max, err = formatBound(exprList[1], abbreviationMap)
		if err != nil {
			return
		}
	}
	return
}

func (cf cronField) handleInvalidExpr(bounds bound, initBounds int) (err error) {
	if cf.min == initBounds || cf.max == initBounds {
		err = errors.New("invalid cron field")
		return
	}

	if cf.min < bounds.min || cf.max > bounds.max {
		err = errors.New("invalid value, out of bounds")
		return
	}

	if cf.min > cf.max {
		err = errors.New("invalid bounds")
		return
	}

	if cf.interval != 1 {
		_range := bounds.max - bounds.min + 1
		if cf.interval > _range || cf.interval == 0 {
			err = errors.New("invalid interval")
			return
		}
	}

	return nil
}

func formatBound(expr string, abbrMap map[string]string) (val int, err error) {
	//handleabbreviations
	abbrVal, ok := abbrMap[strings.ToUpper(expr)]
	if !ok {
		abbrVal = expr
	}

	val, err = strconv.Atoi(abbrVal)
	return
}
