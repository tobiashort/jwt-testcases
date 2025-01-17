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

func AssertErrNil(err error) {
	if err != nil {
		panic(err)
	}
}

func RetrieveCurlCommand() string {
	writer := bufio.NewWriter(os.Stdout)
	writer.WriteString("Enter cURL command (accept with Ctrl-D): ")
	writer.Flush()
	data, err := io.ReadAll(os.Stdin)
	AssertErrNil(err)
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
	cmd := exec.Command("more")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	stdinPipe, err := cmd.StdinPipe()
	AssertErrNil(err)
	err = cmd.Start()
	AssertErrNil(err)
	_, err = stdinPipe.Write([]byte(output))
	AssertErrNil(err)
	err = stdinPipe.Close()
	AssertErrNil(err)
	cmd.Wait()
ask:
	writer := bufio.NewWriter(os.Stdout)
	writer.WriteString("Is the output as expected (y/n)? ")
	writer.Flush()
	reader := bufio.NewReader(os.Stdin)
	data, _, err := reader.ReadLine()
	AssertErrNil(err)
	answer := strings.TrimSpace(string(data))
	answer = strings.ToLower(answer)
	if answer != "y" && answer != "n" {
		goto ask
	}
	return answer == "y"
}

func main() {
	fmt.Println("To begin, you must provide a cURL command for a valid request.")
initialCurlCommand:
	curlCommand := RetrieveCurlCommand()
	encodedJWTs := ExtractJWTs(curlCommand)
	if len(encodedJWTs) == 0 {
		fmt.Println("No JWTs found in cURL command.")
		goto initialCurlCommand
	} else if len(encodedJWTs) > 1 {
		fmt.Println("Multiple JWTs found in cURL command.")
		fmt.Println("This is currently not supported.")
		fmt.Println("Continuing with the first one.")
	}
	jwt, err := ParseEncodedJWT(encodedJWTs[0])
	AssertErrNil(err)
	output := ExecuteCurlCommand(curlCommand)
	outputOk := VerifyOutput(output)
	if !outputOk {
		goto initialCurlCommand
	}
}
