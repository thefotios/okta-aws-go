package okta

import (
	"net/url"

	"github.com/davecgh/go-spew/spew"

	"gopkg.in/resty.v0"
)

type Okta struct {
	client resty.Client
}

func New(hostname string) *Okta {
	url := url.URL{
		Scheme: "https",
		Host:   hostname,
		Path:   "/api/v1",
	}

	o := Okta{}

	o.client = *resty.New()
	o.client.SetHostURL(url.String())

	o.client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		spew.Dump(resp)
		return nil
	})

	return &o
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
		Post("/authn")

	if err != nil {
		return &resty.Response{}, err
	}

	return resp, nil
}
