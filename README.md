Numverify Go Client SDK and CLI app
===================================

[![Build Status][build-status-svg]][build-status-link]
[![Go Report Card][goreport-svg]][goreport-link]
[![Docs][docs-godoc-svg]][docs-godoc-link]
[![License][license-svg]][license-link]

This is a Go client and CLI app for the Numverify API:

https://numverify.com/documentation

## Installation

| Install | Command |
|---------|---------|
| SDK only | `$ go get github.com/grokify/numlookup` |
| CLI only | `$ go get github.com/grokify/numlookup/apps/numverify` |
| SDK and CLI | `$ go get github.com/grokify/numlookup/...` |

## Usage

### CLI App

| Options | Long | Short | Example |
|---------|------|-------|---------|
| `.env` File | `--env` | `-e` | `-e=/path/to/.env` |
| Acces Token | `--token` | `-t` | `-t=<myToken>` |
| Validate Number | `--number` | `-n` | `-n=<number>` |
| List Countries | `--countries` | `-c` | `-c` |

#### Example Commands

```
$ numverify -e=/path/to/.env -n=+16505550100
$ numverify -t=<myToken> -n=+16505550100
$ numverify -e=/path/to/.env -c
$ numverify -t=<myToken> -c
```

### SDK

#### Validate Number

```go
import(
	nv "github.com/grokify/numverify"
)

func main() {
	client := nv.NumverifyClient{AccessKey: "myAccessKey"}

	apiSuccessInfo, apiErrorInfo, resp, err := client.Validate(
		nv.NumverifyParams{Number: number})

	[...]
}
```

#### Get Countries

```go
import(
	nv "github.com/grokify/numverify"
)

func main() {
	client := nv.NumverifyClient{AccessKey: "myAccessKey"}

	countries, apiErrorInfo, resp, err := client.Countries()

	[...]
}
```

 [build-status-svg]: https://api.travis-ci.org/grokify/numlookup.svg?branch=master
 [build-status-link]: https://travis-ci.org/grokify/numlookup
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/numlookup
 [goreport-link]: https://goreportcard.com/report/github.com/grokify/numlookup
 [docs-godoc-svg]: https://img.shields.io/badge/docs-godoc-blue.svg
 [docs-godoc-link]: https://godoc.org/github.com/grokify/numlookup
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-link]: https://github.com/grokify/numlookup/blob/master/LICENSE