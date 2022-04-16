package main

import (
	"fmt"

	"github.com/grokify/gophonenumbers"
	"github.com/grokify/mogo/log/logutil"
	geo "github.com/kellydunn/golang-geo"
)

const (
	USNYCAreaCode  = 212
	USNYCLatGoogle = 40.6976684
	USNYCLonGoogle = -74.2605588

	USSFOAreaCode  = 415
	USSFOLatGoogle = 37.7578149
	USSFOLonGoogle = -122.5078121
)

func GcdGoogle() {
	p1 := geo.NewPoint(USNYCLatGoogle, USNYCLonGoogle)
	p2 := geo.NewPoint(USSFOLatGoogle, USSFOLonGoogle)

	dist := p1.GreatCircleDistance(p2)
	fmt.Printf("Great circle distance NYC to SFO: %v\n", dist)
}

func main() {
	GcdGoogle()

	a2g := gophonenumbers.NewAreaCodeToGeo()
	err := a2g.ReadData()
	logutil.FatalErr(err)

	dist, err := a2g.GcdAreaCodes(USNYCAreaCode, USSFOAreaCode)
	logutil.FatalErr(err)

	fmt.Printf("Great circle distance %v to %v: %v\n", USNYCAreaCode, USSFOAreaCode, dist)
	fmt.Println("DONE")
}
