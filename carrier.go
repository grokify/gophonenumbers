package gophonenumbers

type CarrierNumberInfo struct {
	E164Number        string `json:"e164Number"`
	MobileCountryCode string `json:"mobileCountryCode,omitempty"`
	MobileNetworkCode string `json:"mobileNetworkCode,omitempty"`
	Name              string `json:"name,omitempty"`
	LineType          string `json:"lineType,omitempty"`
	ErrorCode         string `json:"errorCode,omitempty"`
}
