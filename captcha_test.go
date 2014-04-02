package captcha

import "testing"

func TestVerify(t *testing.T) {
	c := New("secret", "test")
	// shall fail with invalid-site-private-key
	ok, err := c.Verify("localhost", "test", "test")
	if ok || err.Error() != "invalid-site-private-key" {
		t.Error(err.Error())
	}
}
