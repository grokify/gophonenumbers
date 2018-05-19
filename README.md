Numverify Go Client SDK and CLI app
===================================

[![Build Status][build-status-svg]][build-status-link]
[![Go Report Card][goreport-svg]][goreport-link]
[![Docs][docs-godoc-svg]][docs-godoc-link]
[![License][license-svg]][license-link]

This is a Go client and CLI app for the Numverify API:

https://numverify.com/documentation

## Installation

```
$ go get github.com/grokify/go-numverify
```

## Usage

### CLI App

```
$ numverify -e=/path/to/.env -n=+16505550100
```

### SDK

#### Validate Number

```go
import(
	"os"

	nv "github.com/grokify/numverify/numverify"
)

func main() {
	client := nv.NumverifyClient{
		AccessKey: os.Getenv("NUMVERIFY_ACCESS_KEY"),
	}

	apiSuccessInfo, apiErrorInfo, resp, err := client.Validate(
		nv.NumverifyParams{Number: number})

	[...]
}
```

#### Get Countries

```go
import(
	"os"

	nv "github.com/grokify/numverify/numverify"
)

func main() {
	client := nv.NumverifyClient{
		AccessKey: os.Getenv("NUMVERIFY_ACCESS_KEY"),
	}

	countries, apiErrorInfo, resp, err := client.Countries()

	[...]
}
```

 [build-status-svg]: https://api.travis-ci.org/grokify/go-numverify.svg?branch=master
 [build-status-link]: https://travis-ci.org/grokify/go-numverify
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/go-numverify
 [goreport-link]: https://goreportcard.com/report/github.com/grokify/go-numverify
 [docs-godoc-svg]: https://img.shields.io/badge/docs-godoc-blue.svg
 [docs-godoc-link]: https://godoc.org/github.com/grokify/go-numverify
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-link]: https://github.com/grokify/go-numverify/blob/master/LICENSE