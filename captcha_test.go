package captcha

import (
	"net/http"
	"testing"
)

func TestVerify(t *testing.T) {
	c := New("", "")
	// shall fail with invalid-site-private-key error
	ok, err := c.Verify(&http.Request{})
	if ok || err != ErrInvalidPrivateKey {
		t.Error(err.Error())
	}
}
