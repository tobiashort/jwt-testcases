package main

import (
	"fmt"
	"time"
)

const ResultStatusPASS = "PASS"
const ResultStatusFAIL = "FAIL"

type TestCase struct {
	ID           string `json:"-"`
	ResultStatus string `json:"resultstatus"`
	Result       string `json:"result"`
	Details      string `json:"details"`
}

func (testCase TestCase) String() string {
	return fmt.Sprintf("%s - %s - %s", testCase.ID, testCase.ResultStatus, testCase.Result)
}

func CheckValidity(originalJWT JWT) TestCase {
	var validFrom float64
	var validTo float64
	exp, expOk := originalJWT.Payload["exp"].(float64)
	nbf, nbfOk := originalJWT.Payload["nbf"].(float64)
	iat, iatOk := originalJWT.Payload["iat"].(float64)
	if expOk {
		validTo = exp
	}
	if nbfOk {
		validFrom = nbf
	} else if iatOk {
		validFrom = iat
	}
	validity := validTo - validFrom
	testCase := TestCase{}
	testCase.ID = "JWT.checkValidity"
	testCase.Details = fmt.Sprintf("%s\n\n%s", originalJWT.Encode(), originalJWT.String())
	if validFrom == 0.0 || validTo == 0.0 {
		testCase.ResultStatus = ResultStatusFAIL
		testCase.Result = `Either the claims "iat", "nbf" or "exp" are missing or invalid.`
		goto done
	}
	testCase.Result = fmt.Sprintf("The JWT is valid for %v.", time.Duration(validity*float64(time.Second)))
	if validity > 3600 {
		testCase.ResultStatus = ResultStatusFAIL
	} else {
		testCase.ResultStatus = ResultStatusPASS
	}
done:
	fmt.Println(testCase)
	return testCase
}
