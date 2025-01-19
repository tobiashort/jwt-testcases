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
	writer := bufio.NewWriter(os.Stdout)
	writer.WriteString("Enter cURL command (accept with Ctrl-D): ")
	writer.Flush()
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
	writer.WriteString("Is the output as expected (y/n)? ")
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

func main() {
	fmt.Println("To begin, you must provide a cURL command for a valid request.")
originalCurlCommand:
	originalCurlCommand := RetrieveCurlCommand()
	originalEncodedJWTs := ExtractJWTs(originalCurlCommand)
	if len(originalEncodedJWTs) == 0 {
		fmt.Println("No JWTs found in cURL command.")
		goto originalCurlCommand
	} else if len(originalEncodedJWTs) > 1 {
		fmt.Println("Multiple JWTs found in cURL command.")
		fmt.Println("This is currently not supported.")
		fmt.Println("Continuing with the first one.")
	}
	originalJWT, err := DecodedJWT(originalEncodedJWTs[0])
	AssertNil(err)
	originalOutput := ExecuteCurlCommand(originalCurlCommand)
	originalOutputOk := VerifyOutput(originalOutput)
	if !originalOutputOk {
		goto originalCurlCommand
	}
	CheckValidity(originalJWT)
}
