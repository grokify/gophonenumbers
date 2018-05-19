package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/grokify/gotilla/config"
	"github.com/grokify/gotilla/fmt/fmtutil"
	nv "github.com/grokify/numverify/numverify"
	"github.com/jessevdk/go-flags"
)

type cliOptions struct {
	EnvFile   string `short:"e" long:"env" description:"Env filepath"`
	Number    string `short:"n" long:"number" description:"Number to verify"`
	Verbose   []bool `short:"v" long:"verbose" description:"Verbose"`
	Countries []bool `short:"c" long:"countries" description:"List Countries"`
	//Token   string `short:"t" long:"token" description:"Access token"`
}

func showNumber(client nv.NumverifyClient, number string) {
	// Returns separate objects for API Success and API Error
	// because Numverify API returns a 200 OK on errors like
	// auth errors.
	apiSuccessInfo, apiErrorInfo, resp, err := client.Validate(
		nv.NumverifyParams{Number: number})
	if err != nil {
		panic(err)
	}
	showResponse(apiSuccessInfo, apiErrorInfo, resp)
}

func showCountries(client nv.NumverifyClient) {
	// Returns separate objects for API Success and API Error
	// because Numverify API returns a 200 OK on errors like
	// auth errors.
	countries, apiErrorInfo, resp, err := client.Countries()
	if err != nil {
		panic(err)
	}
	showResponse(countries, apiErrorInfo, resp)
}

func showResponse(apiSuccessInfo interface{}, apiErrorInfo *nv.NumverifyResponseError, resp *http.Response) {
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

	err = config.LoadDotEnvSkipFirst(opts.EnvFile, os.Getenv("ENV_PATH"), "./.env")
	if err != nil {
		log.Fatal(err)
	}

	client := nv.NumverifyClient{
		AccessKey: os.Getenv(nv.EnvNumverifyAccessKey),
	}

	if len(opts.Number) > 0 {
		showNumber(client, opts.Number)
	} else if len(opts.Countries) > 0 {
		showCountries(client)
	}

	fmt.Println("DONE")
}
