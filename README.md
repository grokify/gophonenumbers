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
$ go get github.com/grokify/numverify
```

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

 [build-status-svg]: https://api.travis-ci.org/grokify/numverify.svg?branch=master
 [build-status-link]: https://travis-ci.org/grokify/numverify
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/numverify
 [goreport-link]: https://goreportcard.com/report/github.com/grokify/numverify
 [docs-godoc-svg]: https://img.shields.io/badge/docs-godoc-blue.svg
 [docs-godoc-link]: https://godoc.org/github.com/grokify/numverify
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-link]: https://github.com/grokify/numverify/blob/master/LICENSE