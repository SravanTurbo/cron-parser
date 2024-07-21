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

var DOW_ABBREVIATIONS = map[string]string{"SUN": "0", "MON": "1", "TUE": "2", "WED": "3", "THU": "4", "FRI": "5", "SAT": "6"}
var MONTH_ABBREVIATIONS = map[string]string{"JAN": "1", "FEB": "2", "MAR": "3", "APR": "4", "MAY": "5", "JUN": "6", "JUL": "7", "AUG": "8", "SEP": "9", "OCT": "10", "NOV": "11", "DEC": "11"}

func PrintCronSchedule(cronExpr string) {
	schedule, err := Parse(cronExpr)
	if err != nil {
		panic(err)
	}

	fmt.Println(schedule)
}

func Parse(cronExpr string) (*Schedule, error) {
	cronFields, err := validate(cronExpr)
	if err != nil {
		return nil, err
	}

	minute, err := parseField(cronFields[0], MinuteBound, map[string]string{})
	if err != nil {
		return nil, err
	}

	hour, err := parseField(cronFields[1], HourBound, map[string]string{})
	if err != nil {
		return nil, err
	}

	dom, err := parseField(cronFields[2], DOMBound, map[string]string{})
	if err != nil {
		return nil, err
	}

	month, err := parseField(cronFields[3], MonthBound, MONTH_ABBREVIATIONS)
	if err != nil {
		return nil, err
	}

	dow, err := parseField(cronFields[4], DOWBound, DOW_ABBREVIATIONS)
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

func parseField(fieldExpr string, bounds bound, abbreviationMap map[string]string) ([]int, error) {
	var values []int

	exprs := strings.Split(fieldExpr, ",")
	for _, expr := range exprs {
		value, err := computeField(expr, bounds, abbreviationMap)
		if err != nil {
			err := errors.New("invalid cron: " + err.Error())
			return nil, err
		}

		if expr == "*" {
			values = value
			break
		}

		values = append(values, value...)
	}

	return values, nil
}

func computeField(expr string, bounds bound, abbreviationMap map[string]string) ([]int, error) {
	fR := NewFieldRange(expr)
	var result []int
	var err error

	if err = fR.handleSlash(); err != nil {
		return nil, err
	}

	if err = fR.handleAsterisk(bounds); err != nil {
		return nil, err
	}

	if err = fR.handleSingleValue(); err != nil {
		return nil, err
	}

	if err = fR.handleHyphen(abbreviationMap); err != nil {
		return nil, err
	}

	if err = fR.handleInvalidExpr(bounds, FRInitBounds); err != nil {
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
