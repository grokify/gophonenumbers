package numverify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/grokify/gophonenumbers"
	hum "github.com/grokify/gotilla/net/httputilmore"
	uu "github.com/grokify/gotilla/net/urlutil"
)

const (
	CountriesEndpoint = "http://apilayer.net/api/countries"
	ValidateEndpoint  = "http://apilayer.net/api/validate"
)

var EnvNumverifyAccessKey = "NUMVERIFY_ACCESS_KEY"

type Client struct {
	AccessKey string
}

func NewClient(accessKey string) Client {
	return Client{AccessKey: strings.TrimSpace(accessKey)}
}

// Returns separate objects for API Success and API Error structs because
// Numverify API will return a 200 OK for errors such as auth errors.
func (nc *Client) Validate(params Params) (*Response, *http.Response, error) {
	params.AccessKey = strings.TrimSpace(params.AccessKey)
	if len(params.AccessKey) == 0 {
		params.AccessKey = nc.AccessKey
	}
	apiURL := uu.BuildURLQueryString(ValidateEndpoint, params)

	resp, respBody, err := hum.GetResponseAndBytes(apiURL)

	nvResp := Response{
		StatusCode: resp.StatusCode,
		ClientErr:  err,
		Body:       string(respBody),
		Time:       time.Now()}
	if err != nil {
		return &nvResp, resp, err
	} else if resp.StatusCode >= 300 {
		return &nvResp, resp, fmt.Errorf("Numverify API Error: %v", resp.StatusCode)
	}

	// Try both success and response. Will
	// error for one.
	var apiSuccessInfo ResponseSuccess
	var apiErrorInfo ResponseError

	err = json.Unmarshal(respBody, &apiSuccessInfo)
	if err != nil {
		err = json.Unmarshal(respBody, &apiErrorInfo)
	}
	nvResp.Success = &apiSuccessInfo
	nvResp.Failure = &apiErrorInfo

	return &nvResp, resp, err
	//return &apiSuccessInfo, &apiErrorInfo, resp, err
}

// Countries returns separate objects for API Success and API Error structs because
// Numverify API will return a 200 OK for errors such as auth errors.
func (nc *Client) Countries() (map[string]Country, *ResponseError, *http.Response, error) {
	apiURL := CountriesEndpoint + "?access_key=" + nc.AccessKey
	resp, respBody, err := hum.GetResponseAndBytes(apiURL)

	countries := map[string]Country{}

	if err != nil {
		return countries, nil, resp, err
	} else if resp.StatusCode >= 300 {
		return countries, nil, resp, fmt.Errorf("Numverify API Error: %v", resp.StatusCode)
	}

	err = json.Unmarshal(respBody, &countries)

	var apiErrorInfo ResponseError
	err = json.Unmarshal(respBody, &apiErrorInfo)

	return countries, &apiErrorInfo, resp, err
}

// Params is the request query parameters for the
// API. AccessKey is added by the client and is not needed
// per-request.
type Params struct {
	AccessKey   string `url:"access_key" json:"access_key,omitempty"`
	Number      string `url:"number" json:"number,omitempty"`
	CountryCode string `url:"country_code" json:"country_code,omitempty"`
	Format      int    `url:"format" json:"format,omitempty"`
	Callback    string `url:"callback" json:"callback,omitempty"`
}

type ResponseSuccess struct {
	Valid               bool   `json:"valid,omitempty"`
	Number              string `json:"number,omitempty"`
	LocalFormat         string `json:"local_format,omitempty"`
	InternationalFormat string `json:"international_format,omitempty"`
	CountryPrefix       string `json:"country_prefix,omitempty"`
	CountryCode         string `json:"country_code,omitempty"`
	CountryName         string `json:"country_name,omitempty"`
	Location            string `json:"location,omitempty"`
	Carrier             string `json:"carrier,omitempty"`
	LineType            string `json:"line_type,omitempty"`
}

func (rs *ResponseSuccess) NumberInfoLookup() gophonenumbers.NumberInfoLookup {
	cni := gophonenumbers.NumberInfoLookup{
		NumberE164: rs.InternationalFormat,
		Carrier: gophonenumbers.Carrier{
			Name:     rs.Carrier,
			LineType: rs.LineType},
	}
	return cni
}

type ResponseError struct {
	Success bool  `json:"success"`
	Error   Error `json:"error,omitempty"`
}

type Error struct {
	Code int    `json:"code,omitempty"`
	Type string `json:"type,omitempty"`
	Info string `json:"info,omitempty"`
}

type Country struct {
	CountryName string `json:"country_name,omitempty"`
	DialingCode string `json:"dialling_code,omitempty"`
}

type Response struct {
	StatusCode int
	Body       string
	ClientErr  error
	Success    *ResponseSuccess
	Failure    *ResponseError
	Time       time.Time
}
