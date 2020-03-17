package numlookup

import "github.com/grokify/gotilla/net/httputilmore"

const (
	CarrierATT     = "att.com"
	CarrierSprint  = "sprint.com"
	CarrierTMobile = "t-mobile.com"
	CarrierVerizon = "verizon.com"
)

type NumberInfo struct {
	Number            string
	Carrier           string
	LineType          string
	CallerName        string
	MobileCountryCode int
	MobileNetworkCode int
	ApiHistory        map[string]httputilmore.ResponseInfo
}
