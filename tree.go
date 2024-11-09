package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"strings"
)

func createTree() string {
	entries := readIndex()
	content := ""
	for _, entry := range entries {
		content += fmt.Sprintf("%o %s\x00", entry.mode, entry.path)
		content += string(entry.sha1)
		// fmt.Println("entry: ", entry.mode, entry.path, hex.EncodeToString(entry.sha1))
	}
	hash := hashObject(strings.NewReader(content), "tree", len(content))
	writeToObjectFile(strings.NewReader(content), hash, "tree", len(content))
	return hash
}

func readTree(tree string) []string {
	objectType, contents := readObject(tree)
	if objectType != "tree" {
		log.Fatalf("Invalid tree object type: %v", objectType)
	}

	// lines := strings.Split(contents[0], "\x00")
	// fmt.Println([]byte(contents[0]))
	// fmt.Println()
	// fmt.Println(strings.Join(lines, "\n"))

	content := strings.Join(contents, "\n")
	// fmt.Println("contents len of %v: ", tree, len(content))
	var objects []string
	i := 0
	for true {
		if i >= len(content) {
			break
		}

		end := strings.Index(content[i:], "\x00")
		end = i + end

		if end == -1 || end+21 > len(content) {
			break
		}
		// fmt.Println("end: ", end)
		hash := content[end+1 : end+21]
		hash = hex.EncodeToString([]byte(hash))
		objects = append(objects, hash)
		i = end + 21
	}

	return objects
}
