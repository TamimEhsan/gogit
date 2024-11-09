package main

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
)

func gitGetPack(url, userName, password string) []string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.SetBasicAuth(userName, password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get response: %v", err)
	}
	defer resp.Body.Close()

	bufferedReader := bufio.NewReader(resp.Body)
	var lines []string
	for {
		line, err := bufferedReader.ReadString('\n')
		if err != nil {
			break
		}
		lines = append(lines, line)
	}
	return lines
}

func gitPushPack(url, userName, password string, data []byte) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.SetBasicAuth(userName, password)
	req.Header.Set("Content-Type", "application/x-git-receive-pack-request")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get response: %v", err)
	}
	defer resp.Body.Close()
}
