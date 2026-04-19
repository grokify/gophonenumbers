# GoPhoneNumbers

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
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/gophonenumbers/blob/master/LICENSE

A Go library for phone number parsing, formatting, validation, and geolocation. Includes clients for Numverify and Twilio Lookup APIs.

## Features

- **Phone Number Parsing**: Parse E.164, national, and international formats
- **NANP Components**: Extract area code, exchange code, and line number from North American numbers
- **Multiple Formats**: Get E.164, national, international, and RFC3966 formats
- **Geolocation**: US area code to geographic coordinates mapping (embedded data, no external files)
- **Distance Calculation**: Great circle distance between area codes
- **Fictitious Numbers**: Generate valid US fictitious phone numbers (555-01XX range)
- **API Clients**: Numverify and Twilio Lookup API integrations
- **CLI Tools**: Command-line apps for number validation and lookup

## Installation

```bash
go get github.com/grokify/gophonenumbers
```

## Quick Start

### Parse and Format Phone Numbers

```go
package main

import (
    "fmt"
    "github.com/grokify/gophonenumbers"
)

func main() {
    // Parse an E.164 number into multiple formats
    formats, err := gophonenumbers.FormatsParse("+16505551234", "US")
    if err != nil {
        panic(err)
    }

    fmt.Println("E.164:", formats.E164)           // +16505551234
    fmt.Println("National:", formats.National)    // (650) 555-1234
    fmt.Println("International:", formats.International) // +1 650-555-1234
    fmt.Println("RFC3966:", formats.RFC3966)      // tel:+1-650-555-1234
}
```

### Extract NANP Components

```go
package main

import (
    "fmt"
    "github.com/grokify/gophonenumbers"
)

func main() {
    num := gophonenumbers.Number{E164Number: "+16505551234"}
    comp, err := num.NANPComponents()
    if err != nil {
        panic(err)
    }

    fmt.Println("Country Code:", comp.CountryCode)     // 1
    fmt.Println("Area Code:", comp.NANPAreaCode)       // 650
    fmt.Println("Exchange Code:", comp.NANPExchangeCode) // 555
    fmt.Println("Line Number:", comp.NANPLineNumber)   // 1234
}
```

### US Area Code Geolocation

The library includes embedded US area code geolocation data - no external files required.

```go
package main

import (
    "fmt"
    "github.com/grokify/gophonenumbers"
)

func main() {
    // Load embedded area code data
    a2g := gophonenumbers.NewAreaCodeToGeo()
    if err := a2g.ReadData(); err != nil {
        panic(err)
    }

    // Get area code info
    aci := a2g.AreaCodeInfos[650] // San Francisco Bay Area
    fmt.Printf("Area Code 650: Lat %.4f, Lon %.4f\n",
        aci.Point.Lat(), aci.Point.Lng())

    // Calculate distance between area codes
    dist, _ := a2g.GcdAreaCodes(650, 212) // SF to NYC
    fmt.Printf("Distance 650 to 212: %.0f km\n", dist)
}
```

### Generate Fictitious Phone Numbers

Generate valid US fictitious numbers in the 555-01XX range for testing:

```go
package main

import (
    "fmt"
    "github.com/grokify/gophonenumbers"
)

func main() {
    a2g := gophonenumbers.NewAreaCodeToGeo()
    if err := a2g.ReadData(); err != nil {
        panic(err)
    }

    fng := gophonenumbers.NewFakeNumberGenerator(a2g.AreaCodes())

    // Generate a random fictitious number
    num, _ := fng.RandomLocalNumberUS()
    fmt.Printf("Fictitious number: +%d\n", num)

    // Generate unique numbers
    set := map[uint64]int8{}
    for i := 0; i < 5; i++ {
        num, set, _ = fng.RandomLocalNumberUSUnique(set)
        fmt.Printf("Unique number %d: +%d\n", i+1, num)
    }
}
```

## API Clients

### Numverify

Validate phone numbers using the [Numverify API](https://numverify.com/documentation):

```go
package main

import (
    "fmt"
    nv "github.com/grokify/gophonenumbers/numverify"
)

func main() {
    client := nv.Client{AccessKey: "your-access-key"}

    resp, _, _, err := client.Validate(nv.Params{Number: "+16505551234"})
    if err != nil {
        panic(err)
    }

    if resp.Success != nil {
        fmt.Println("Valid:", resp.Success.Valid)
        fmt.Println("Carrier:", resp.Success.Carrier)
        fmt.Println("Line Type:", resp.Success.LineType)
    }
}
```

### Twilio Lookup

Look up phone numbers using the [Twilio Lookup API](https://www.twilio.com/docs/lookup/api):

```go
package main

import (
    "fmt"
    "github.com/grokify/gophonenumbers/twilio"
)

func main() {
    client := twilio.NewClient("account-sid", "auth-token")

    info, err := client.Validate("+16505551234", &twilio.Params{Type: "carrier"})
    if err != nil {
        panic(err)
    }

    fmt.Println("Phone Number:", info.PhoneNumber)
    fmt.Println("Carrier:", info.Carrier.Name)
}
```

## CLI Tools

### Numverify CLI

```bash
# Install
go install github.com/grokify/gophonenumbers/cmd/numverify@latest

# Validate a number
numverify -t=<access-key> -n=+16505551234

# Using .env file
numverify -e=/path/to/.env -n=+16505551234

# List supported countries
numverify -t=<access-key> -c
```

### Area Code Distance

```bash
# Install
go install github.com/grokify/gophonenumbers/cmd/areacode_distance@latest

# Calculate distance between area codes
areacode_distance 650 212
```

## Data Reference

### Embedded Data

The library embeds US area code geolocation data from the [Area Code Geolocation Database](https://github.com/ravisorg/Area-Code-Geolocation-Database). No external files or GOPATH configuration required.

### Number Components

| Component | Description | Example |
|-----------|-------------|---------|
| Country Code | ITU-T E.164 country code | 1 (US/Canada) |
| Area Code (NPA) | Numbering Plan Area code | 650 |
| Exchange Code (NXX) | Central office code | 555 |
| Line Number | Subscriber number | 1234 |

### Fictitious Numbers

The library generates numbers in the reserved 555-01XX range per [NANPA guidelines](https://www.nanpa.com/number_resource_info/555_service_numbers.html). These numbers are safe for testing and will never conflict with real numbers.

## Contributing

Contributions are welcome. Please open an issue or submit a pull request.

## License

MIT License - see [LICENSE](LICENSE) for details.
