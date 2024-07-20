package main

import (
	"os"

	"github.com/SravanTurbo/cron-parser/pkg/cronparser"
)

func main() {
	cronExpr := os.Args[1]
	cronparser.PrintCronSchedule(cronExpr)
}
