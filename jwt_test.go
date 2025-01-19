package main

import "testing"

func TestDecodeJWT(t *testing.T) {
	encodedJWTs := []string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
	}
	for _, encodedJWT := range encodedJWTs {
		_, err := DecodedJWT(encodedJWT)
		if err != nil {
			t.Error(encodedJWT, err)
		}
	}
}
