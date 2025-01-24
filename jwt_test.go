package main

import "testing"

func TestDecodeJWT(t *testing.T) {
	encodedJWTs := []string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
	}
	for _, encodedJWT := range encodedJWTs {
		_, err := DecodeJWT(encodedJWT)
		if err != nil {
			t.Error(encodedJWT, err)
		}
	}
}

func TestEncodeJWT(t *testing.T) {
	encodedJWTs := []string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
	}
	for _, encodedJWT := range encodedJWTs {
		decodedJWT, err := DecodeJWT(encodedJWT)
		if err != nil {
			t.Error(encodedJWT, err)
		}
		reencodedJWT := decodedJWT.Encode()
		if reencodedJWT != encodedJWT {
			t.Errorf("\n   %s\n-> %s", encodedJWT, reencodedJWT)
		}
	}
}
