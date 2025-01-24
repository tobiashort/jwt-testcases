package main

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"text/tabwriter"
	"time"
)

const (
	ResultStatusPASS   = "PASS"
	ResultStatusFAIL   = "FAIL"
	ResultStatusManual = "MANUAL"
)

type TestCase struct {
	ID             string
	Description    string
	ExpectedResult string
	ActualResult   string
	ResultStatus   string
	Details        string
}

func (testCase TestCase) String() string {
	columnWidth := 30

	applyMaxWidthToText := func(text string) string {
		if len(text) < columnWidth {
			return text
		}
		builder := strings.Builder{}
		line := strings.Builder{}
		for _, word := range strings.Fields(text) {
			if len(line.String())+len(word) < columnWidth {
				if line.String() != "" {
					line.WriteString(" ")
				}
				line.WriteString(word)
			} else {
				line.WriteString("\n")
				builder.WriteString(line.String())
				line = strings.Builder{}
				line.WriteString(word)
			}
		}
		builder.WriteString(line.String())

		return builder.String()
	}

	description := applyMaxWidthToText(testCase.Description)
	descriptionCol := []string{"Desciption", "-------------------------------"}
	for _, line := range strings.Split(description, "\n") {
		descriptionCol = append(descriptionCol, line)
	}

	expectedResult := applyMaxWidthToText(testCase.ExpectedResult)
	expectedResultCol := []string{"Expected Result", "-------------------------------"}
	for _, line := range strings.Split(expectedResult, "\n") {
		expectedResultCol = append(expectedResultCol, line)
	}

	actualResult := applyMaxWidthToText(testCase.ActualResult)
	actualResultCol := []string{"Actual Result", "-------------------------------"}
	for _, line := range strings.Split(actualResult, "\n") {
		actualResultCol = append(actualResultCol, line)
	}

	resultStatus := applyMaxWidthToText(testCase.ResultStatus)
	resultStatusCol := []string{"PASS/FAIL", "-------------------------------", resultStatus}

	maxRow := math.Max(float64(len(descriptionCol)), float64(len(expectedResultCol)))
	maxRow = math.Max(maxRow, float64(len(actualResultCol)))
	maxRow = math.Max(maxRow, float64(len(resultStatusCol)))

	table := make([][]string, 4)

	for idx := range table {
		table[idx] = make([]string, int(maxRow))
	}

	for idx := range descriptionCol {
		table[0][idx] = descriptionCol[idx]
	}

	for idx := range expectedResultCol {
		table[1][idx] = expectedResultCol[idx]
	}

	for idx := range actualResultCol {
		table[2][idx] = actualResultCol[idx]
	}

	for idx := range resultStatusCol {
		table[3][idx] = resultStatusCol[idx]
	}

	builder := strings.Builder{}
	for row := 0; row < int(maxRow); row++ {
		for col := 0; col < 4; col++ {
			builder.WriteString(table[col][row])
			builder.WriteString("\t")
		}
		builder.WriteString("\n")
	}

	buffer := new(bytes.Buffer)
	writer := tabwriter.NewWriter(buffer, columnWidth+4, 0, 4, ' ', 0)
	fmt.Fprint(writer, builder.String())
	err := writer.Flush()
	AssertNil(err)
	return buffer.String()
}

func CheckValidity(originalJWT JWT) TestCase {
	testCase := TestCase{}
	testCase.ID = "JWT.checkValidity"
	testCase.Description = "How long is the token valid?"
	testCase.ExpectedResult = "At maximum 1 hour."
	var validFrom float64
	var validTo float64
	exp, expOk := originalJWT.Payload.Get("exp")
	nbf, nbfOk := originalJWT.Payload.Get("nbf")
	iat, iatOk := originalJWT.Payload.Get("iat")
	if !expOk {
		testCase.ActualResult = `The JWT does not have an expiration time, as the "exp" claim is missing.`
		testCase.ResultStatus = ResultStatusFAIL
		return testCase
	} else {
		validTo = exp.(float64)
	}
	if nbfOk {
		validFrom = nbf.(float64)
	} else if iatOk {
		validFrom = iat.(float64)
	}
	if validFrom == 0 && validTo > 0 {
		testCase.ActualResult = `It was not possible to determine the validity period of the JWT because the
"iat" or "nbf" claim is missing, although the "exp" claim is present. The calculation must be performed
manually by extracting the creation time from the "exp" claim.`
		testCase.ResultStatus = ResultStatusManual
		return testCase
	}
	validity := validTo - validFrom
	testCase.ActualResult = fmt.Sprintf("The JWT is valid for %v.", time.Duration(validity*float64(time.Second)))
	if validity > 3600 {
		testCase.ResultStatus = ResultStatusFAIL
	} else {
		testCase.ResultStatus = ResultStatusPASS
	}
	return testCase
}

func CheckSignatureExclusionAttack(originalCurlCommand, originalCurlCommandOutput, originalEncodedJWT, canary string) TestCase {
	testCase := TestCase{}
	testCase.ID = "JWT.checkSignatureExclusionAttack"
	testCase.Description = "Is it possible to use tokens without a signature (signature exclusion attack)?"
	testCase.ExpectedResult = "No."
	jwt, err := DecodeJWT(originalEncodedJWT)
	AssertNil(err)
	jwt.Signature = ""
	curlCommand := strings.ReplaceAll(originalCurlCommand, originalEncodedJWT, jwt.Encode())
	curlCommandOutput := ExecuteCurlCommand(curlCommand)
	similarity := CosineSimilarity(originalCurlCommandOutput, curlCommandOutput)
	hasCanary := canary != ""
	containsCanary := strings.Contains(curlCommandOutput, canary)
	if hasCanary && containsCanary {
		testCase.ActualResult = fmt.Sprintf("Yes. (similarity %f, response contains canary)", similarity)
		testCase.ResultStatus = ResultStatusFAIL
		return testCase
	} else if hasCanary && !containsCanary {
		testCase.ActualResult = fmt.Sprintf("Yes. (similarity %f, response does not contain canary)", similarity)
		testCase.ResultStatus = ResultStatusPASS
		return testCase
	} else if !hasCanary && similarity > 0.9 {
		testCase.ActualResult = fmt.Sprintf("Yes. (similarity %f)", similarity)
		testCase.ResultStatus = ResultStatusFAIL
		return testCase
	} else {
		testCase.ActualResult = fmt.Sprintf("Yes. (similarity %f)", similarity)
		testCase.ResultStatus = ResultStatusPASS
		return testCase
	}
}
