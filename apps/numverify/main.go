package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	nv "github.com/grokify/gophonenumbers/numverify"
	"github.com/grokify/simplego/config"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/jessevdk/go-flags"
)

type cliOptions struct {
	EnvFile   string `short:"e" long:"env" description:"Env filepath"`
	Token     string `short:"t" long:"token" description:"Access token"`
	Number    string `short:"n" long:"number" description:"Number to verify"`
	Verbose   []bool `short:"v" long:"verbose" description:"Verbose"`
	Countries []bool `short:"c" long:"countries" description:"List Countries"`
}

func showNumber(client nv.Client, number string) {
	// Returns separate objects for API Success and API Error
	// because Numverify API returns a 200 OK on errors like
	// auth errors.
	info, resp, err := client.Validate(
		nv.Params{Number: number})
	if err != nil {
		panic(err)
	}
	showResponse(info.Success, info.Failure, resp)
}

func showCountries(client nv.Client) {
	// Returns separate objects for API Success and API Error
	// because Numverify API returns a 200 OK on errors like
	// auth errors.
	countries, apiErrorInfo, resp, err := client.Countries()
	if err != nil {
		panic(err)
	}
	showResponse(countries, apiErrorInfo, resp)
}

func showResponse(apiSuccessInfo interface{}, apiErrorInfo *nv.ResponseError, resp *http.Response) {
	fmt.Printf("API_RESPONSE_STATUS: [%v]\n", resp.StatusCode)
	fmtutil.PrintJSON(apiErrorInfo)
	fmtutil.PrintJSON(apiSuccessInfo)
}

// main: Usage: numverify --number=+16505550100 --verbose
func main() {
	opts := cliOptions{}

	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}
	if len(opts.Verbose) > 0 {
		fmtutil.PrintJSON(opts)
	}

	err = config.LoadDotEnvFirst(opts.EnvFile, os.Getenv("ENV_PATH"), "./.env")
	if err != nil {
		log.Fatal(err)
	}

	numverifyAccessToken := strings.TrimSpace(opts.Token)
	if len(numverifyAccessToken) == 0 {
		numverifyAccessToken = os.Getenv(nv.EnvNumverifyAccessKey)
	}

	client := nv.Client{AccessKey: numverifyAccessToken}

	if len(opts.Number) > 0 {
		showNumber(client, opts.Number)
	} else if len(opts.Countries) > 0 {
		showCountries(client)
	}

	fmt.Println("DONE")
}
