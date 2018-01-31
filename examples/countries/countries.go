package main

import (
	"fmt"
	"os"

	nv "github.com/grokify/go-numverify"
	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/joho/godotenv"
)

func main() {
	if len(os.Getenv("ENV_PATH")) > 0 {
		err := godotenv.Load(os.Getenv("ENV_PATH"))
		if err != nil {
			panic(err)
		}
	}

	client := nv.NumverifyClient{
		AccessKey: os.Getenv(nv.EnvNumverifyAccessKey),
	}

	// Returns separate objects for API Success and API Error
	// because Numverify API returns a 200 OK on errors like
	// auth errors.
	apiSuccessInfo, apiErrorInfo, resp, err := client.Countries()
	if err != nil {
		panic(err)
	}
	fmt.Printf("STATUS %v\n", resp.StatusCode)
	fmtutil.PrintJSON(apiErrorInfo)
	fmtutil.PrintJSON(apiSuccessInfo)

	fmt.Println("DONE")
}
