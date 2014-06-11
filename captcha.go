// reCAPTCHA Go API
//
// https://developers.google.com/recaptcha/intro
package captcha

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Captcha struct {
	private string
	Public  string
	Server  string
}

const apiServer = "https://www.google.com/recaptcha/api"

func New(private, public string) Captcha {
	return Captcha{private: private, Public: public, Server: apiServer}
}

func (c Captcha) Verify(remoteip, challenge, response string) (bool, error) {
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
	answer := strings.Split(string(body), "\n")
	return answer[0] == "true", errors.New(answer[1])
}

// Error codes:
// invalid-site-private-key: We weren't able to verify the private key.
// invalid-request-cookie:   The challenge parameter of the verify script was incorrect.
// incorrect-captcha-sol:    The CAPTCHA solution was incorrect.
// captcha-timeout:          The solution was received after the CAPTCHA timed out.
// recaptcha-not-reachable:  reCAPTCHA never returns this error code.
