package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func DoAssertErrNil(err error) {
	if err != nil {
		panic(err)
	}
}

func RetrieveCurlCommand() string {
	writer := bufio.NewWriter(os.Stdout)
	writer.WriteString("Enter CURL command (Accept with Ctrl-D): ")
	writer.Flush()
	data, err := io.ReadAll(os.Stdin)
	DoAssertErrNil(err)
	return string(data)
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
	DoAssertErrNil(err)
	err = cmd.Start()
	DoAssertErrNil(err)
	_, err = stdinPipe.Write([]byte(output))
	DoAssertErrNil(err)
	err = stdinPipe.Close()
	DoAssertErrNil(err)
	cmd.Wait()
ask:
	writer := bufio.NewWriter(os.Stdout)
	writer.WriteString("Is the output as expected (y/n)? ")
	writer.Flush()
	reader := bufio.NewReader(os.Stdin)
	data, _, err := reader.ReadLine()
	DoAssertErrNil(err)
	answer := strings.TrimSpace(string(data))
	answer = strings.ToLower(answer)
	if answer != "y" && answer != "n" {
		goto ask
	}
	return answer == "y"
}

func main() {
initialCurlCommand:
	curlCommand := RetrieveCurlCommand()
	output := ExecuteCurlCommand(curlCommand)
	outputOk := VerifyOutput(output)
	if !outputOk {
		goto initialCurlCommand
	}
}
