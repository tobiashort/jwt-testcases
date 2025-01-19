package main

import "testing"

func TestCheckValidity(t *testing.T) {
	type Example struct {
		EncodedJWT   string
		Result       string
		ResultStatus string
	}
	examples := []Example{
		{
			EncodedJWT:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			Result:       `Either the claims "iat", "nbf" or "exp" are missing or invalid.`,
			ResultStatus: ResultStatusFAIL,
		},
		{
			EncodedJWT:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwibmJmIjoxNTE2MjM5MDIyfQ.cH149pgS9p7gJEOhw46Xx14pZqqBNXBQVF879jOab40",
			Result:       `Either the claims "iat", "nbf" or "exp" are missing or invalid.`,
			ResultStatus: ResultStatusFAIL,
		},
		{
			EncodedJWT:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZXhwIjoxNTE2MjQyNjIzfQ.VScdrNk4twmhYvSegvq1o2jhu9x8Ne8bbbLoPMW-B6E",
			Result:       `Either the claims "iat", "nbf" or "exp" are missing or invalid.`,
			ResultStatus: ResultStatusFAIL,
		},
		{
			EncodedJWT:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwibmJmIjoxNTE2MjM5MDIyLCJleHAiOjE1MTYyNDI2MjN9.jP8TD1xJ4fwomMONu1lA0mWuR_Y4wm1UG4In22B3mM8",
			Result:       "The JWT is valid for 1h0m1s.",
			ResultStatus: ResultStatusFAIL,
		},
		{
			EncodedJWT:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MTYyNDI2MjN9.2eBhGnQTNpCDv6LjP28s3Sv97CXuzkpYikrrZKoVnQQ",
			Result:       "The JWT is valid for 1h0m1s.",
			ResultStatus: ResultStatusFAIL,
		},
		{
			EncodedJWT:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIxLCJuYmYiOjE1MTYyMzkwMjIsImV4cCI6MTUxNjI0MjYyM30.tvunqL5HubJ9nzIiGniFkHZ92_76QlR7LzGFuA-XVZY",
			Result:       "The JWT is valid for 1h0m1s.",
			ResultStatus: ResultStatusFAIL,
		},
		{
			EncodedJWT:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MTYyNDI2MjJ9.kdtHmE7_Cg5Vp_kuDoGp6XnWIfOxeaopNXumatFfDQg",
			Result:       "The JWT is valid for 1h0m0s.",
			ResultStatus: ResultStatusPASS,
		},
		{
			EncodedJWT:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwibmJmIjoxNTE2MjM5MDIyLCJleHAiOjE1MTYyNDI2MjJ9.XnntYuARfvKENAyLGORZeMi8SvAe-eRqY_COaO5eJ1k",
			Result:       "The JWT is valid for 1h0m0s.",
			ResultStatus: ResultStatusPASS,
		},
		{
			EncodedJWT:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIxLCJuYmYiOjE1MTYyMzkwMjIsImV4cCI6MTUxNjI0MjYyMn0.TdqnrW7Y7s7sf_GAF9wxK7cbZAlEm9WDxt-Xrrx7zMY",
			Result:       "The JWT is valid for 1h0m0s.",
			ResultStatus: ResultStatusPASS,
		},
	}
	for _, example := range examples {
		decodedJWT, err := DecodedJWT(example.EncodedJWT)
		if err != nil {
			t.Error(example.EncodedJWT, err)
		}
		testCase := CheckValidity(decodedJWT)
		if testCase.ResultStatus != example.ResultStatus {
			t.Errorf("%s\nExpected '%s', got '%s'", example.EncodedJWT, example.ResultStatus, testCase.ResultStatus)
		}
		if testCase.Result != example.Result {
			t.Errorf("%s\nExpected '%s', got '%s'", example.EncodedJWT, example.Result, testCase.Result)
		}
	}
}
