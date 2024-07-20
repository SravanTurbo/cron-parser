package cronparser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
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

func parseField(expr string, bounds bound) ([]int, error) {
	min, max, interval := 0, -1, 1
	var result []int
	var err error

	value, _ := strconv.Atoi(expr)
	if value > 0 || expr == "0" {
		result = append(result, value)
		return result, nil
	}

	if strings.Contains(expr, ",") {
		values := strings.Split(expr, ",")
		for i := 0; i < len(values); i++ {
			value, err := strconv.Atoi(strings.TrimSpace(values[i]))
			if err != nil || bounds.min > value || bounds.max < value {
				return nil, errors.New("invalid value provided")
			}
			result = append(result, value)
		}

		return result, nil
	}

	if strings.Contains(expr, "/") {
		exprList := strings.Split(expr, "/")

		interval, err = strconv.Atoi(exprList[1])
		if err != nil {
			return nil, errors.New("invalid interval provided")
		}

		expr = exprList[0]
	}

	if strings.Contains(expr, "*") {
		if len(expr) > 1 {
			return nil, errors.New("invalid cron expression")
		}

		max = bounds.max
	}

	if strings.Contains(expr, "-") {
		exprList := strings.Split(expr, "-")
		if len(exprList) > 2 {
			return nil, errors.New("invalid range provided")
		}

		min, err = strconv.Atoi(exprList[0])
		if err != nil || min < bounds.min {
			return nil, errors.New("invalid minimum value of range")
		}

		max, err = strconv.Atoi(exprList[1])
		if err != nil || max > bounds.max {
			return nil, errors.New("invalid maximum value of range")
		}
	}

	if interval > max && max != -1 {
		return nil, errors.New("interval out of bounds")
	}

	result = buildIntList(min, max, interval)

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
