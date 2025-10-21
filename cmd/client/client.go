package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	endpoint := "http://localhost:8080/"

	fmt.Println("Enter long URL")

	reader := bufio.NewReader(os.Stdin)

	long, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}
	long = strings.TrimSuffix(long, "\n")

	client := &http.Client{}

	request, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(long))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	request.Header.Add("Content-Type", "text/plain")

	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	defer response.Body.Close()
	fmt.Println("Status code: ", response.Status)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Println(string(body))
}
