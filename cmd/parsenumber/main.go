package main

import (
	"fmt"
	"strings"

	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	flags "github.com/jessevdk/go-flags"
	"github.com/nyaruka/phonenumbers"
)

const ExampleBadNumber = "+16611106327"

type Options struct {
	Number  string `short:"N" long:"number" description:"Phone Number"`
	Country string `short:"C" long:"country" description:"Phone Number Country Hint"`
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)
	logutil.FatalErr(err)

	opts.Number = strings.TrimSpace(opts.Number)
	if len(opts.Number) == 0 {
		fmt.Println("no number provided, using ExampleBadNumber = " + ExampleBadNumber)
		opts.Number = ExampleBadNumber
	}

	num, err := phonenumbers.Parse(opts.Number, opts.Country)
	logutil.FatalErr(err)
	fmtutil.MustPrintJSON(num)
}
