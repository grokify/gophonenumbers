package numverify

import (
	"encoding/json"
	"fmt"
	"net/http"

	hum "github.com/grokify/gotilla/net/httputilmore"
	uu "github.com/grokify/gotilla/net/urlutil"
)

const NumverifyApiUrl = "http://apilayer.net/api/validate"

var EnvNumverifyAccessKey = "NUMVERIFY_ACCESS_KEY"

type NumverifyClient struct {
	AccessKey string
}

// Returns separate objects for API Success and API Error structs because
// Numverify API will return a 200 OK for errors such as auth errors.
func (nc *NumverifyClient) Get(params NumverifyParams) (*NumverifyResponseSuccess, *NumverifyResponseError, *http.Response, error) {
	if len(params.AccessKey) == 0 {
		params.AccessKey = nc.AccessKey
	}
	apiUrl := uu.BuildURLQueryString(NumverifyApiUrl, params)

	resp, respBody, err := hum.GetResponseAndBytes(apiUrl)

	if err != nil {
		return nil, nil, resp, err
	} else if resp.StatusCode >= 300 {
		return nil, nil, resp, fmt.Errorf("Numverify API Error: %v", resp.StatusCode)
	}
	var apiSuccessInfo NumverifyResponseSuccess
	err = json.Unmarshal(respBody, &apiSuccessInfo)

	var apiErrorInfo NumverifyResponseError
	err = json.Unmarshal(respBody, &apiErrorInfo)

	return &apiSuccessInfo, &apiErrorInfo, resp, err
}

// NumverifyParams is the request query parameters for the
// API. AccessKey is added by the client and is not needed
// per-request.
type NumverifyParams struct {
	AccessKey   string `url:"access_key" json:"access_key,omitempty"`
	Number      string `url:"number" json:"number,omitempty"`
	CountryCode string `url:"country_code" json:"country_code,omitempty"`
	Format      int    `url:"format" json:"format,omitempty"`
	Callback    string `url:"callback" json:"callback,omitempty"`
}

type NumverifyResponseSuccess struct {
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

type NumverifyResponseError struct {
	Success bool           `json:"success"`
	Error   NumverifyError `json:"error,omitempty"`
}

type NumverifyError struct {
	Code int    `json:"code,omitempty"`
	Type string `json:"type,omitempty"`
	Info string `json:"info,omitempty"`
}
