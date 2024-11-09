package main

import (
	"log"
	"os"
	"path"
	"sort"
)

func gitAdd(files []string) {
	indexEntries := make([]indexEntry, 0)
	filemap := make(map[string]bool)
	for _, filename := range files {
		filemap[filename] = true
	}
	indexes := readIndex()
	for _, index := range indexes {
		if _, ok := filemap[index.path]; !ok {
			indexEntries = append(indexEntries, index)
		}
	}
	for _, filename := range files {

		// create a file
		filename = path.Join(filename)
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Failed to open file %s: %v", filename, err)
		}
		defer file.Close()

		sz := getFileSize(file)
		fileHash := hashObject(file, "blob", sz)
		file.Seek(0, 0)
		writeToObjectFile(file, fileHash, "blob", sz)

		indexEntry := createIndexEntry(filename, fileHash)

		indexEntries = append(indexEntries, indexEntry)

	}
	sort.Slice(indexEntries, func(i, j int) bool { return indexEntries[i].path < indexEntries[j].path })
	writeToIndexFile(indexEntries)
}
