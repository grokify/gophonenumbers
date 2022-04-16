package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/grokify/gophonenumbers/numverify"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/type/stringsutil"
	flags "github.com/jessevdk/go-flags"
)

type cliOptions struct {
	EnvFile    string   `short:"e" long:"env" description:"Env filepath"`
	Token      string   `short:"t" long:"token" description:"Access token"`
	Numbers    []string `short:"n" long:"numbers" description:"Numbers to verify"`
	Verbose    []bool   `short:"v" long:"verbose" description:"Verbose"`
	Countries  []bool   `short:"c" long:"countries" description:"List Countries"`
	WriteFiles []bool   `short:"w" long:"writefiles" description:"WriteFiles"`
	Fileroot   string   `long:"fileroot" description:"File At line count"`
	FileAt     uint     `long:"fileat" description:"File At line count"`
	LogAt      uint     `long:"logat" description:"Log At line count"`
}

func (opts *cliOptions) Inflate() {
	newNumbers := []string{}
	for _, num := range opts.Numbers {
		newNumbers = append(newNumbers, strings.Split(num, ",")...)
	}
	opts.Numbers = stringsutil.SliceCondenseSpace(newNumbers, true, true)
}

func showNumber(client *numverify.Client, number string) {
	// Returns separate objects for API Success and API Error
	// because Numverify API returns a 200 OK on errors like
	// auth errors.
	info, resp, err := client.Validate(
		numverify.Params{Number: number})
	if err != nil {
		panic(err)
	}
	showResponse(info.Success, info.Failure, resp)
}

func showCountries(client *numverify.Client) {
	// Returns separate objects for API Success and API Error
	// because Numverify API returns a 200 OK on errors like
	// auth errors.
	countries, apiErrorInfo, resp, err := client.Countries()
	if err != nil {
		panic(err)
	}
	showResponse(countries, apiErrorInfo, resp)
}

func showResponse(apiSuccessInfo interface{}, apiErrorInfo *numverify.ResponseError, resp *http.Response) {
	fmt.Printf("API_RESPONSE_STATUS: [%v]\n", resp.StatusCode)
	fmtutil.PrintJSON(apiErrorInfo)
	fmtutil.PrintJSON(apiSuccessInfo)
}

// main: Usage: numverify --number=+16505550100 -n=+16505550101 --verbose
func main() {
	opts := cliOptions{}

	_, err := flags.Parse(&opts)
	logutil.FatalErr(err)
	opts.Inflate()
	if len(opts.Verbose) > 0 {
		fmtutil.MustPrintJSON(opts)
	}

	err = config.LoadDotEnvFirst(opts.EnvFile, os.Getenv("ENV_PATH"), "./.env")
	logutil.FatalErr(err)

	numverifyAccessToken := strings.TrimSpace(opts.Token)
	if len(numverifyAccessToken) == 0 {
		numverifyAccessToken = os.Getenv(numverify.EnvNumverifyAccessKey)
	}

	client := &numverify.Client{AccessKey: numverifyAccessToken}

	if len(opts.WriteFiles) > 0 {
		res := numverify.GetWriteValidationMulti(client,
			opts.Numbers, []string{}, opts.Fileroot, opts.LogAt, opts.FileAt)
		fmtutil.PrintJSON(res)
	} else {
		if len(opts.Numbers) > 0 {
			for _, number := range opts.Numbers {
				showNumber(client, number)
			}
		} else if len(opts.Countries) > 0 {
			showCountries(client)
		}
	}

	fmt.Println("DONE")
}
