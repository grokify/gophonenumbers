package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/grokify/gophonenumbers/twilio"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	EnvFile      string `short:"e" long:"env" description:"Env filepath"`
	PhoneNumbers string `short:"n" long:"number" description:"Phone number" required:"true"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	loaded, err := config.LoadDotEnv(
		[]string{os.Getenv("ENV_PATH"), "./.env", opts.EnvFile}, -1)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.MustPrintJSON(loaded)

	client, err := twilio.NewClient(
		os.Getenv(twilio.EnvTwilioAccountSid),
		os.Getenv(twilio.EnvTwilioAuthToken))
	if err != nil {
		log.Fatal(err)
	}

	phoneNumbers := stringsutil.SliceCondenseSpace(
		strings.Split(opts.PhoneNumbers, ","), true, false)

	if len(phoneNumbers) == 1 {
		info, err := client.Validate(
			opts.PhoneNumbers,
			&twilio.Params{Type: "carrier"})
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.MustPrintJSON(info)
	} else {
		mr := twilio.GetWriteValidationMulti(
			client, phoneNumbers, []string{}, "_twilio_multi", uint(2), uint(2))
		fmtutil.MustPrintJSON(mr)
	}

	fmt.Println("DONE")
}
