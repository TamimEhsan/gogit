package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func gitStatus() {
	dirIndexes := getDirIndexes()
	indexes := readIndex()
	modified, untracked, deleted := compareIndexes(dirIndexes, indexes)
	printStatus(modified, untracked, deleted)
}

func getDirIndexes() []indexEntry {
	var dirIndexes []indexEntry
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !strings.HasPrefix(path, ".git") {

			file, err := os.Open(path)
			if err != nil {
				log.Fatalf("Failed to open file %s: %v", path, err)
			}
			defer file.Close()

			hash := hashObject(file, "blob", getFileSize(file))
			sha1, _ := hex.DecodeString(hash)
			dirIndexes = append(dirIndexes, indexEntry{path: path, sha1: sha1})
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return dirIndexes
}

func compareIndexes(dirIndexes, indexes []indexEntry) ([]string, []string, []string) {
	sort.Slice(dirIndexes, func(i, j int) bool { return dirIndexes[i].path < dirIndexes[j].path })
	sort.Slice(indexes, func(i, j int) bool { return indexes[i].path < indexes[j].path })

	modified, untracked, deleted := []string{}, []string{}, []string{}
	i, j := 0, 0
	for i < len(dirIndexes) && j < len(indexes) {
		if dirIndexes[i].path == indexes[j].path {
			if hex.EncodeToString(dirIndexes[i].sha1) != hex.EncodeToString(indexes[j].sha1) {
				modified = append(modified, dirIndexes[i].path)
			}
			i++
			j++
		} else if dirIndexes[i].path < indexes[j].path {
			untracked = append(untracked, dirIndexes[i].path)
			i++
		} else {
			deleted = append(deleted, indexes[j].path)
			j++
		}
	}

	for i < len(dirIndexes) {
		untracked = append(untracked, dirIndexes[i].path)
		i++
	}

	for j < len(indexes) {
		deleted = append(deleted, indexes[j].path)
		j++
	}

	return modified, untracked, deleted
}

func printStatus(modified, untracked, deleted []string) {
	const colorRed = "\033[0;31m"
	const colorNone = "\033[0m"

	fmt.Println("Changes not staged for commit:")
	fmt.Println("  (use \"git add/rm <file>...\" to update what will be committed)")
	fmt.Println("  (use \"git restore <file>...\" to discard changes in working directory)")

	fmt.Print(colorRed)
	for _, file := range modified {
		fmt.Println("\tmodified:   ", file)
	}
	for _, file := range deleted {
		fmt.Println("\tdeleted:    ", file)
	}
	fmt.Print(colorNone)
	fmt.Println("Untracked files:")
	fmt.Println("  (use \"git add <file>...\" to include in what will be committed)")
	fmt.Print(colorRed)
	for _, file := range untracked {
		fmt.Println("\t", file)
	}
	fmt.Print(colorNone)
}
