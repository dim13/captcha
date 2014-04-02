package captcha

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Captcha struct {
	private string
	Public  string
	Server  string
	answer  []string
}

var errorcodes = map[string]string{
	"invalid-site-private-key": "We weren't able to verify the private key.",
	"invalid-request-cookie":   "The challenge parameter of the verify script was incorrect.",
	"incorrect-captcha-sol":    "The CAPTCHA solution was incorrect.",
	"captcha-timeout":          "The solution was received after the CAPTCHA timed out.",
	"recaptcha-not-reachable":  "reCAPTCHA never returns this error code.",
}

const apiServer = "http://www.google.com/recaptcha/api"

func New(private, public string) *Captcha {
	return &Captcha{private: private, Public: public, Server: apiServer}
}

func (c *Captcha) Verify(remoteip, challenge, response string) (bool, error) {
	values := url.Values{
		"privatekey": {c.private},
		"remoteip":   {remoteip},
		"challenge":  {challenge},
		"response":   {response},
	}
	resp, err := http.PostForm(c.Server+"/verify", values)
	if err != nil {
		return false, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return false, err
	}
	c.answer = strings.Split(string(body), "\n")
	return c.answer[0] == "true", c
}

func (c *Captcha) Error() string {
	return errorcodes[c.answer[1]]
}
