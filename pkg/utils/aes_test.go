package utils

import (
	"net/url"
	"testing"
)

func TestDecryptWithBase64(t *testing.T) {
	text := "Ex47@P_i$"
	key := "abcdEFGHIJKLxyzx"
	iv := "9876543210XingYe"

	sign, err := EncryptWithBase64([]byte(key), []byte(iv), []byte(text))

	if err != nil {
		t.Error(err)
		return
	}
	t.Log("sign=", sign)

	plaintext, err := DecryptWithBase64([]byte(key), []byte(iv), sign)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("plaintext=", plaintext)

}

func TestEncryptWithBase64(t *testing.T) {
	str := url.QueryEscape("tb4@d1W2!")
	t.Log("str=", str)
}
