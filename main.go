package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func AssertNil(val any) {
	if val != nil {
		panic(val)
	}
}

func RetrieveCurlCommand() string {
	fmt.Println("Please paste your cURL command below and press [Ctrl-D] to confirm:")
	data, err := io.ReadAll(os.Stdin)
	AssertNil(err)
	return string(data)
}

func ExtractJWTs(input string) []string {
	jwtPattern := `ey[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+`
	re := regexp.MustCompile(jwtPattern)
	matches := re.FindAllString(input, -1)
	return matches
}

func ExecuteCurlCommand(curlCommand string) string {
	cmd := exec.Command("bash", "-c", curlCommand)
	data, err := cmd.CombinedOutput()
	output := string(data)
	if err != nil {
		fmt.Print(output)
		panic(err)
	}
	return output
}

func VerifyOutput(output string) bool {
	fmt.Print(`
Your cURL command has been executed successfully.
Please review the response to confirm whether it
indicates success, meaning you received an
authorized response.

Press [Enter] to continue.
`)
	_, _, err := bufio.NewReader(os.Stdin).ReadLine()
	AssertNil(err)
	cmd := exec.Command("more")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	stdinPipe, err := cmd.StdinPipe()
	AssertNil(err)
	err = cmd.Start()
	AssertNil(err)
	_, err = stdinPipe.Write([]byte(output))
	AssertNil(err)
	err = stdinPipe.Close()
	AssertNil(err)
	cmd.Wait()
ask:
	writer := bufio.NewWriter(os.Stdout)
	writer.WriteString("\nIs the response as expected (y/n)? ")
	writer.Flush()
	reader := bufio.NewReader(os.Stdin)
	data, _, err := reader.ReadLine()
	AssertNil(err)
	answer := strings.TrimSpace(string(data))
	answer = strings.ToLower(answer)
	if answer != "y" && answer != "n" {
		goto ask
	}
	return answer == "y"
}

func RetrieveCanaryToken() string {
	fmt.Print(`
Enter your canary token and press [Enter] to confirm,
or skip by simply pressing [Enter]: `)
	canary, err := bufio.NewReader(os.Stdin).ReadString('\n')
	AssertNil(err)
	// strip newline at the end
	canary = canary[:len(canary)-1]
	return canary
}

func main() {
	fmt.Print(`Welcome to jwt-testcases!

This tool is designed to assist you with automated
attacks on JSON Web Tokens (JWT).

To get started, provide a cURL command for a valid
request that includes a JWT. The response to this request
should indicate success, meaning you must receive an
authorized response. For example, you can copy the cURL
command directly from Burp Suite by following these steps:
Right-click on the request →  Select "Copy as curl command (bash)."

`)

originalCurlCommand:
	originalCurlCommand := RetrieveCurlCommand()
	originalEncodedJWTs := ExtractJWTs(originalCurlCommand)
	if len(originalEncodedJWTs) == 0 {
		fmt.Println("The provided cURL command does not contain any JWT.")
		goto originalCurlCommand
	} else if len(originalEncodedJWTs) > 1 {
		fmt.Print(`
The provided cURL command contains multiple JWTs.
This is currently not supported.
We will continue with the first one found.
`)
	}
	originalEncodedJWT := originalEncodedJWTs[0]
	originalJWT, err := DecodeJWT(originalEncodedJWT)
	AssertNil(err)
	originalCurlCommandOutput := ExecuteCurlCommand(originalCurlCommand)
	originalCurlCommandOutputOk := VerifyOutput(originalCurlCommandOutput)
	if !originalCurlCommandOutputOk {
		goto originalCurlCommand
	}
	fmt.Print(`
Would you like to provide a canary token?
This is a unique string found only in authorized responses,
which helps minimize false positives.
Providing it is optional.
`)
canary:
	canary := RetrieveCanaryToken()
	if canary != "" && !strings.Contains(originalCurlCommandOutput, canary) {
		fmt.Print(`
The provided canary token is not contained the response of
the cURL command you have entered. Perhaps you have a typo.
Please retry.
`)
		goto canary
	}

	fmt.Print(`
Performing test cases...
`)

	testCases := []TestCase{
		CheckValidity(originalJWT),
		CheckSignatureExclusionAttack(originalCurlCommand, originalCurlCommandOutput, originalEncodedJWT, canary),
	}
	for _, testCase := range testCases {
		fmt.Printf("\n%s\n", testCase.String())
	}
}
