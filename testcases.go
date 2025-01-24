package main

import (
	"fmt"
	"strings"
	"time"
)

const ResultStatusPASS = "PASS"
const ResultStatusFAIL = "FAIL"

type TestCase struct {
	ID           string `json:"id"`
	Description  string `json:"description"`
	Result       string `json:"result"`
	ResultStatus string `json:"resultstatus"`
	Details      string `json:"details"`
}

func (testCase TestCase) String() string {
	return fmt.Sprintf("%s - %s - %s", testCase.ID, testCase.ResultStatus, testCase.Result)
}

func CheckValidity(originalJWT JWT) TestCase {
	testCase := TestCase{}
	testCase.ID = "JWT.checkValidity"
	testCase.Description = "How long is the token valid?"
	var validFrom float64
	var validTo float64
	exp, expOk := originalJWT.Payload.Get("exp")
	nbf, nbfOk := originalJWT.Payload.Get("nbf")
	iat, iatOk := originalJWT.Payload.Get("iat")
	if expOk {
		validTo = exp.(float64)
	}
	if nbfOk {
		validFrom = nbf.(float64)
	} else if iatOk {
		validFrom = iat.(float64)
	}
	validity := validTo - validFrom
	details := strings.Builder{}
	details.WriteString("JWT\n")
	details.WriteString(originalJWT.Encode())
	details.WriteString("\n\nDecoded\n")
	details.WriteString(originalJWT.String())
	if validFrom == 0.0 || validTo == 0.0 {
		testCase.ResultStatus = ResultStatusFAIL
		testCase.Result = `Either the claims "iat", "nbf" or "exp" are missing or invalid.`
		testCase.Details = details.String()
		fmt.Println(testCase)
		return testCase
	}
	details.WriteString("\n\nValidity\n")
	details.WriteString(fmt.Sprintf("%.0f - %.0f = %.0f", validTo, validFrom, validity))
	testCase.Result = fmt.Sprintf("The JWT is valid for %v.", time.Duration(validity*float64(time.Second)))
	testCase.Details = details.String()
	if validity > 3600 {
		testCase.ResultStatus = ResultStatusFAIL
	} else {
		testCase.ResultStatus = ResultStatusPASS
	}
	fmt.Println(testCase)
	return testCase
}

func CheckSignatureExclusionAttack(originalCurlCommand, originalCurlCommandOutput string, originalEncodedJWT string) TestCase {
	testCase := TestCase{}
	testCase.ID = "JWT.checkSignatureExclusionAttack"
	testCase.Description = "Is it possible to use tokens without a signature (signature exclusion attack)?"
	jwt, err := DecodeJWT(originalEncodedJWT)
	AssertNil(err)
	jwt.Signature = ""
	curlCommand := strings.ReplaceAll(originalCurlCommand, originalEncodedJWT, jwt.Encode())
	curlCommandOutput := ExecuteCurlCommand(curlCommand)
	similarity := CosineSimilarity(originalCurlCommandOutput, curlCommandOutput)
	details := strings.Builder{}
	details.WriteString("$ ")
	details.WriteString(curlCommand)
	details.WriteString("\n")
	details.WriteString(curlCommandOutput)
	testCase.Details = details.String()
	testCase.Result = ""
	if similarity > 0.9 {
		testCase.Result = fmt.Sprintf("Yes. (similarity %f)", similarity)
		testCase.ResultStatus = ResultStatusFAIL
	} else {
		testCase.Result = fmt.Sprintf("No. (similarity %f)", similarity)
		testCase.ResultStatus = ResultStatusPASS
	}
	fmt.Println(testCase)
	return testCase
}
