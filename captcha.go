// reCAPTCHA 1.0 Go API
//
// https://developers.google.com/recaptcha/old/intro
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

var (
	ErrInvalidPrivateKey = errors.New("We weren't able to verify the private key")
	ErrInvalidChallenge  = errors.New("The challenge parameter of the verify script was incorrect")
	ErrInvalidSolution   = errors.New("The CAPTCHA solution was incorrect")
	ErrTimeout           = errors.New("The solution was received after the CAPTCHA timed out")
	ErrNotReachable      = errors.New("reCAPTCHA never returns this error code")
)

var errMap = map[string]error{
	"success":                  nil,
	"invalid-site-private-key": ErrInvalidPrivateKey,
	"invalid-request-cookie":   ErrInvalidChallenge,
	"incorrect-captcha-sol":    ErrInvalidSolution,
	"captcha-timeout":          ErrTimeout,
	"recaptcha-not-reachable":  ErrNotReachable,
}

func New(private, public string) Captcha {
	return Captcha{private: private, Public: public, Server: apiServer}
}

func remoteip(r *http.Request) string {
	ra := r.RemoteAddr
	if i := strings.LastIndex(ra, ":"); i > 0 {
		ra = ra[:i]
	}
	return ra
}

func (c Captcha) Verify(r *http.Request) (bool, error) {
	values := url.Values{
		"privatekey": {c.private},
		"remoteip":   {remoteip(r)},
		"challenge":  {r.PostFormValue("recaptcha_challenge_field")},
		"response":   {r.PostFormValue("recaptcha_response_field")},
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
	if len(answer) != 2 {
		return false, ErrNotReachable
	}
	err, ok := errMap[answer[1]]
	if !ok {
		err = ErrNotReachable
	}
	return answer[0] == "true", err
}
