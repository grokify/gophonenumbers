GoPhoneNumbers
==============

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/grokify/gophonenumbers/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/grokify/gophonenumbers/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/grokify/gophonenumbers/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/grokify/gophonenumbers/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/grokify/gophonenumbers/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/grokify/gophonenumbers/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/gophonenumbers
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/gophonenumbers
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/gophonenumbers
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/gophonenumbers
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fgophonenumbers
 [loc-svg]: https://tokei.rs/b1/github/grokify/gophonenumbers
 [repo-url]: https://github.com/grokify/gophonenumbers
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/gophonenumbers/blob/master/LICENSE

This library provides phone number information functionality including ability to access number look up APIs. It includes a Go client and CLI app for the Numverify API, and Twilio API.

https://numverify.com/documentation

## Installation

| Install | Command |
|---------|---------|
| SDK only | `$ go get github.com/grokify/gophonenumbers` |
| CLI only | `$ go get github.com/grokify/gophonenumbers/apps/numverify` |
| SDK and CLI | `$ go get github.com/grokify/gophonenumbers/...` |

## Usage

### CLI App

| Options | Long | Short | Example |
|---------|------|-------|---------|
| `.env` File | `--env` | `-e` | `-e=/path/to/.env` |
| Access Token | `--token` | `-t` | `-t=<myToken>` |
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
	nv "github.com/grokify/gophonenumbers/numverify"
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
	nv "github.com/grokify/gophonenumbers/numverify"
)

func main() {
	client := nv.NumverifyClient{AccessKey: "myAccessKey"}

	countries, apiErrorInfo, resp, err := client.Countries()

	[...]
}
```
