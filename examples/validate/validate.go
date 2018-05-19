package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	nv "github.com/grokify/go-numverify"
	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/joho/godotenv"
)

// main: Usage: verify -number +16505550100
func main() {
	if len(os.Getenv("ENV_PATH")) > 0 {
		err := godotenv.Load(os.Getenv("ENV_PATH"))
		if err != nil {
			panic(err)
		}
	}

	var numToVerify string
	flag.StringVar(&numToVerify, "number", "", "Number to verify")
	flag.Parse()

	if len(strings.TrimSpace(numToVerify)) == 0 {
		fmt.Println("Usage: validate.go --number=+16505550200")
		return
	}

	client := nv.NumverifyClient{
		AccessKey: os.Getenv(nv.EnvNumverifyAccessKey),
	}

	p := nv.NumverifyParams{
		Number: numToVerify,
	}

	// Returns separate objects for API Success and API Error
	// because Numverify API returns a 200 OK on errors like
	// auth errors.
	apiSuccessInfo, apiErrorInfo, resp, err := client.Validate(p)
	if err != nil {
		panic(err)
	}
	fmt.Printf("STATUS %v\n", resp.StatusCode)
	fmtutil.PrintJSON(apiErrorInfo)
	fmtutil.PrintJSON(apiSuccessInfo)

	fmt.Println("DONE")
}
