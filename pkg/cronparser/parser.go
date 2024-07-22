package cronparser

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

const VALID_NUM_OF_CRON_FIELDS = 6

type bound struct {
	min, max int
}

var MinuteBound = bound{0, 59}
var HourBound = bound{0, 23}
var DOMBound = bound{1, 31}
var MonthBound = bound{1, 12}
var DOWBound = bound{0, 6}

var DOW_ABBREVIATIONS = map[string]string{"SUN": "0", "MON": "1", "TUE": "2", "WED": "3", "THU": "4", "FRI": "5", "SAT": "6"}
var MONTH_ABBREVIATIONS = map[string]string{"JAN": "1", "FEB": "2", "MAR": "3", "APR": "4", "MAY": "5", "JUN": "6", "JUL": "7", "AUG": "8", "SEP": "9", "OCT": "10", "NOV": "11", "DEC": "12"}

func PrintCronSchedule(cronExpr string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

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

func validate(cronExpr string) ([]string, error) {
	cronFields := strings.Split(cronExpr, " ")
	if len(cronFields) != VALID_NUM_OF_CRON_FIELDS {
		return nil, errors.New("Validation Error: invalid number of cron fields")
	}

	pattern := `[/*/,/-/\0-9]|(MON|TUE|WED|THU|FRI|SAT|SUN)|(JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC)`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, errors.New("Validation Error: failed to compile Regexp")
	}

	cmdIndex := VALID_NUM_OF_CRON_FIELDS - 1
	for i := 0; i < cmdIndex; i++ {
		cronFields[i] = strings.ToUpper(cronFields[i])
		if !re.MatchString(cronFields[i]) {
			return nil, errors.New("Validation Error: invalid time field")
		}
	}

	return cronFields, nil
}

func parseField(fieldExpr string, bounds bound, abbreviationMap map[string]string) ([]int, error) {
	uniqueValueMap := make(map[int]struct{})

	//handleComma:
	exprs := strings.Split(fieldExpr, ",")
	for _, expr := range exprs {
		valueList, err := handleNonComma(expr, bounds, abbreviationMap)
		if err != nil {
			err := errors.New("Parsing Error: " + err.Error())
			return nil, err
		}

		for _, val := range valueList {
			uniqueValueMap[val] = struct{}{}
		}
	}

	uniqueValueList := make([]int, 0, len(uniqueValueMap))
	for key := range uniqueValueMap {
		uniqueValueList = append(uniqueValueList, key)
	}

	sort.Ints(uniqueValueList)

	return uniqueValueList, nil
}

func handleNonComma(expr string, bounds bound, abbreviationMap map[string]string) ([]int, error) {
	var err error

	cf := NewCronField(expr)

	if err = cf.handleSlash(); err != nil {
		return nil, err
	}

	if err = cf.handleAsterisk(bounds); err != nil {
		return nil, err
	}

	if err = cf.handleSingleValue(); err != nil {
		return nil, err
	}

	if err = cf.handleHyphen(abbreviationMap); err != nil {
		return nil, err
	}

	if err = cf.handleInvalidExpr(bounds, FRInitBounds); err != nil {
		return nil, err
	}

	cf.valueList = buildIntList(cf.min, cf.max, cf.interval)

	return cf.valueList, nil
}
