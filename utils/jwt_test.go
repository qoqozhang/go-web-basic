package utils

import (
	"testing"
)

var testSigningStr string

func TestMiddlewareJwt_SignedString(t *testing.T) {
	m := new(MiddlewareJwt).Init()
	data := map[string]interface{}{
		"username": "zhl",
		"role":     "admin",
	}
	var err error
	if testSigningStr, err = m.SignedString(data); err != nil {
		t.Error(err)
		return
	}
	t.Logf("jwt token: %s", testSigningStr)
}

func TestMiddlewareJwt_Parse(t *testing.T) {
	m := new(MiddlewareJwt).Init()
	valid, data, err := m.Parse(testSigningStr)
	if err != nil {
		t.Errorf("unvalid: %v", err)
		return
	}
	t.Logf("jwt valid: %v, username: %s, role: %s", valid, data["username"], data["role"])

}
