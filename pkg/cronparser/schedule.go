package cronparser

import (
	"fmt"
)

type Schedule struct {
	minute, hour, dom, month, dow []int
	cmd                           string
}

func (s Schedule) String() string {
	outputFormat := "minute\t\t%s\nhour\t\t%s\nday of month\t%s\nmonth\t\t%s\nday of week\t%s\ncommand\t\t%s"
	minuteString := intsJoin(s.minute, " ")
	hourString := intsJoin(s.hour, " ")
	domString := intsJoin(s.dom, " ")
	monthString := intsJoin(s.month, " ")
	dowString := intsJoin(s.dow, " ")
	cmdString := s.cmd

	return fmt.Sprintf(outputFormat, minuteString, hourString, domString, monthString, dowString, cmdString)
}
