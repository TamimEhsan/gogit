package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func createCommit(msg string) {

	now := time.Now()
	timestamp := now.Unix()
	timezone := now.Format("-0700")

	treeHash := createTree()
	commitContent := ""
	commitContent += fmt.Sprintf("tree %s\n", treeHash)
	currentCommit := getLocalMasterCommit()
	if currentCommit != "" {
		commitContent += fmt.Sprintf("parent %s\n", currentCommit)
	}

	name := getUserName()
	email := getEmail()

	commitContent += fmt.Sprintf("author %s <%s> %d %s\n", name, email, timestamp, timezone)
	commitContent += fmt.Sprintf("committer %s <%s> %d %s\n", name, email, timestamp, timezone)
	commitContent += "\n"
	commitContent += msg
	commitContent += "\n"

	hash := hashObject(strings.NewReader(commitContent), "commit", len(commitContent))
	writeToObjectFile(strings.NewReader(commitContent), hash, "commit", len(commitContent))
	writeCommit(hash)
}

func getLocalMasterCommit() string {
	file, err := os.Open(path.Join(".git", "refs", "heads", "master"))
	if err != nil {
		// return 0000... if the file does not exist
		return ""
	}
	defer file.Close()
	hash := make([]byte, 40)
	file.Read(hash)
	return string(hash)
}

func getRemoteMasterCommit(url, userName, password string) string {
	url = url + "/info/refs?service=git-receive-pack"

	lines := gitGetPack(url, userName, password)

	hash := lines[1][8:48]
	return hash
}

func writeCommit(commit string) {
	file, err := os.OpenFile(path.Join(".git", "refs", "heads", "master"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open commit file: %v", err)
	}
	defer file.Close()
	file.Truncate(0)
	file.Write([]byte(commit))
}
