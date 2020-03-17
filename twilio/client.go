package twilio

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/grokify/gotilla/net/httputilmore"
	"github.com/grokify/oauth2more"
)

const (
	LookupEndpoint = "https://lookups.twilio.com/v1/PhoneNumbers/"
)

var (
	EnvTwilioAccountSid = "TWILIO_ACCOUNT_SID"
	EnvTwilioAuthToken  = "TWILIO_AUTH_TOKEN"
)

type Client struct {
	accountSid string
	authToken  string
	httpClient *http.Client
}

func NewClient(accountSid, authToken string) (*Client, error) {
	httpClient, err := oauth2more.NewClientBasicAuth(
		accountSid, authToken, false)
	if err != nil {
		return nil, err
	}
	return &Client{httpClient: httpClient}, nil
}

// Params represents optional parameters. Use "Type=carrier"
// for carrier info and "Type=caller-name" for caller name.
type Params struct {
	CountryCode string
	Type        string
}

func (p *Params) Encode() string {
	p.CountryCode = strings.TrimSpace(p.CountryCode)
	p.Type = strings.TrimSpace(p.Type)
	if len(p.CountryCode) == 0 && len(p.Type) == 0 {
		return ""
	}
	qry := url.Values{}
	if len(p.CountryCode) > 0 {
		qry.Add("CountryCode", p.CountryCode)
	}
	if len(p.Type) > 0 {
		qry.Add("Type", p.Type)
	}
	return qry.Encode()
}

func (c *Client) Validate(number string, opts *Params) (NumberInfo, error) {
	ni := NumberInfo{}
	number = strings.TrimSpace(number)
	if len(number) == 0 {
		return ni, errors.New("E_NO_NUMBER")
	}
	if opts == nil {
		opts = &Params{}
	}
	apiURL := LookupEndpoint + strings.TrimSpace(number) +
		"?" + opts.Encode()
	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return ni, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ni, err
	}
	err = json.Unmarshal(body, &ni)
	if err != nil {
		return ni, err
	}
	ni.ApiResponseInfo = httputilmore.ResponseInfo{
		Name:       "twilio",
		URL:        apiURL,
		StatusCode: resp.StatusCode,
		Time:       time.Now(),
		Body:       string(body)}
	return ni, nil
}

type NumberInfo struct {
	CallerName      map[string]string         `json:"caller_name,omitempty"`
	CountryCode     string                    `json:"country_code,omitempty"`
	PhoneNumber     string                    `json:"phone_number,omitempty"`
	NationalFormat  string                    `json:"national_format,omitempty"`
	URL             string                    `json:"url,omitempty"`
	Carrier         Carrier                   `json:"carrier,omitempty"`
	ApiResponseInfo httputilmore.ResponseInfo `json:"api_response_info,omitempty"`
}

type Carrier struct {
	MobileCountryCode string `json:"mobile_country_code,omitempty"`
	MobileNetworkCode string `json:"mobile_network_code,omitempty"`
	Name              string `json:"name,omitempty"`
	Type              string `json:"type,omitempty"`
	ErrorCode         string `json:"error_code,omitempty"`
}
