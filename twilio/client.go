package twilio

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/grokify/goauth"
	"github.com/grokify/mogo/net/httputilmore"
)

const (
	LookupEndpoint = "https://lookups.twilio.com/v1/PhoneNumbers/"
)

var (
	EnvTwilioAccountSid = "TWILIO_ACCOUNT_SID"
	EnvTwilioAuthToken  = "TWILIO_AUTH_TOKEN" // #nosec G101
)

type Client struct {
	httpClient *http.Client
}

func NewClient(accountSid, authToken string) (*Client, error) {
	httpClient, err := goauth.NewClientBasicAuth(
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ni, err
	}
	err = json.Unmarshal(body, &ni)
	if err != nil {
		return ni, err
	}
	ni.APIResponseInfo = httputilmore.ResponseInfo{
		Name:       "twilio",
		URL:        apiURL,
		StatusCode: resp.StatusCode,
		Time:       time.Now(),
		Body:       string(body)}
	return ni, nil
}
