package cronparser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type bound struct {
	min, max int
}

var MinuteBound = bound{0, 59}
var HourBound = bound{0, 23}
var DOMBound = bound{1, 31}
var MonthBound = bound{1, 12}
var DOWBound = bound{0, 6}

func PrintCronSchedule(cronExpr string) {
	schedule, err := parse(cronExpr)
	if err != nil {
		panic(err)
	}

	fmt.Println(schedule)
}

func parse(cronExpr string) (*Schedule, error) {
	cronFields, err := validate(cronExpr)
	if err != nil {
		return nil, err
	}

	minute, err := parseField(cronFields[0], MinuteBound)
	if err != nil {
		return nil, err
	}

	hour, err := parseField(cronFields[1], HourBound)
	if err != nil {
		return nil, err
	}

	dom, err := parseField(cronFields[2], DOMBound)
	if err != nil {
		return nil, err
	}

	month, err := parseField(cronFields[3], MonthBound)
	if err != nil {
		return nil, err
	}

	dow, err := parseField(cronFields[4], DOWBound)
	if err != nil {
		return nil, err
	}

	return &Schedule{
		minute: minute,
		hour:   hour,
		dom:    dom,
		month:  month,
		dow:    dow,
		cmd:    cronFields[5]}, nil
}

func parseField(fieldExpr string, bounds bound) ([]int, error) {
	var values []int

	exprs := strings.Split(fieldExpr, ",")
	for _, expr := range exprs {
		value, err := computeField(expr, bounds)
		if err != nil {
			return nil, err
		}

		values = append(values, value...)
	}

	return values, nil
}

func computeField(expr string, bounds bound) ([]int, error) {
	fR := NewFieldRange(expr)
	var result []int

	if err := fR.handleSlash(); err != nil {
		return nil, err
	}

	if err := fR.handleAsterisk(bounds); err != nil {
		return nil, err
	}

	if err := fR.handleHyphen(); err != nil {
		return nil, err
	}

	if err := fR.handleSingleValue(); err != nil {
		return nil, err
	}

	if err := fR.handleInvalidExpr(bounds, FRInitBounds); err != nil {
		return nil, err
	}

	result = buildIntList(fR.min, fR.max, fR.interval)

	return result, nil

}

func validate(cronExpr string) ([]string, error) {
	cronFields := strings.Split(cronExpr, " ")
	if len(cronFields) != 6 {
		return nil, errors.New("Invalid cron expression, check README")
	}

	pattern := `[,*-/0-9]`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, errors.New("Failed to compile Regexp")
	}

	for i := 0; i < 6; i++ {
		if !re.MatchString(cronFields[i]) {
			return nil, errors.New("Invalid cron field")
		}
	}

	return cronFields, nil
}
