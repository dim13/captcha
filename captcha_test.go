package captcha

import "testing"

func TestVerify(t *testing.T) {
	c := New("secret", "test")
	// shall fail with invalid-site-private-key
	ok, err := c.Verify("localhost", "test", "test")
	if ok || err != ErrInvalidPrivateKey {
		t.Error(err.Error())
	}
}
