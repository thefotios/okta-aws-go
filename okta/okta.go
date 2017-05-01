package okta

import (
	"net/url"
	"time"

	"github.com/davecgh/go-spew/spew"

	"gopkg.in/resty.v0"
)

type Okta struct {
	client resty.Client
	debug  bool
}

func New(hostname string, debug bool) *Okta {
	url := url.URL{
		Scheme: "https",
		Host:   hostname,
		Path:   "/api/v1",
	}

	o := Okta{
		debug: debug,
	}

	o.client = *resty.New()
	o.client.SetHostURL(url.String())

	o.client.SetDebug(debug)

	return &o
}

type AuthenticationTransactionModel struct {
	StateToken   string    `json:"stateToken"`
	SessionToken string    `json:"sessionToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	Status       string    `json:"status"`
	RelayState   string    `json:"relayState"`
	FactorResult string    `json:"factorResult"`
}

// https://github.com/oktadeveloper/okta-aws-cli-assume-role/blob/master/src/main/java/com/okta/tools/awscli.java#L179
func (okta Okta) PasswordLogin(username, password string) (*resty.Response, error) {
	resp, err := okta.client.R().
		SetHeaders(map[string]string{
			"Accept":        "application/json",
			"Content-Type":  "application/json",
			"Cache-Control": "no-cache",
		}).
		SetBody(map[string]string{
			"username": username,
			"password": password,
		}).
		SetResult(&AuthenticationTransactionModel{}).
		Post("/authn")

	if err != nil {
		return &resty.Response{}, err
	}

	reply, ok := resp.Result().(*AuthenticationTransactionModel)
	if !ok {
		panic("Unable to convert")
	}

	if okta.debug {
		spew.Dump(reply)
	}

	return resp, nil
}
