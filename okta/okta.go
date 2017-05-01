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

// http://developer.okta.com/docs/api/resources/authn.html#authentication-transaction-modelhttp://developer.okta.com/docs/api/resources/authn.html#authentication-transaction-model
type AuthenticationTransactionModel struct {
	StateToken   string    `json:"stateToken"`
	SessionToken string    `json:"sessionToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	// TODO: Status can be an Enum with custom marshalling: https://golang.org/pkg/encoding/json/#example__customMarshalJSON
	Status       string `json:"status"`
	RelayState   string `json:"relayState"`
	FactorResult string `json:"factorResult"`
}

// http://developer.okta.com/docs/api/resources/authn.html#options-object
type AuthnOptions struct {
	MultiOptionalFactorEnroll bool `json:"multiOptionalFactorEnroll"`
	WarnBeforePasswordExpired bool `json:"warnBeforePasswordExpired"`
}

// http://developer.okta.com/docs/api/resources/authn.html#context-object
type AuthnContext struct {
	DeviceToken string `json:"deviceToken"`
}

// http://developer.okta.com/docs/api/resources/authn.html#request-parameters
type AuthnRequest struct {
	Username   string       `json:"username"`
	Password   string       `json:"password"`
	RelayState string       `json:"relayState"`
	Token      string       `json:"token"`
	Options    AuthnOptions `json:"options"`
	Context    AuthnContext `json:"context"`
}

// https://github.com/oktadeveloper/okta-aws-cli-assume-role/blob/master/src/main/java/com/okta/tools/awscli.java#L179
func (okta Okta) PasswordLogin(username, password string) (*resty.Response, error) {
	params := AuthnRequest{
		Username: username,
		Password: password,
	}
	resp, err := okta.client.R().
		SetHeaders(map[string]string{
			"Accept":        "application/json",
			"Content-Type":  "application/json",
			"Cache-Control": "no-cache",
		}).
		SetBody(params).
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
