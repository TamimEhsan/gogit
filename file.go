package main

import (
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path"
)

func gitCatFile(file string) {
	// check if the file exists
	if _, err := os.Stat(path.Join(".git", "objects", file[:2], file[2:])); os.IsNotExist(err) {
		log.Fatalf("File not found: %v", file)
	}
	objectFile, err := os.Open(path.Join(".git", "objects", file[:2], file[2:]))
	if err != nil {
		log.Fatalf("Failed to open object file: %v", err)
	}
	defer objectFile.Close()

	zlibReader, _ := zlib.NewReader(objectFile)
	defer zlibReader.Close()

	io.Copy(os.Stdout, zlibReader)

}

func createDir(path string) {
	// check if the directory exists
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return
	}
	if err := os.Mkdir(path, 0755); err != nil {
		log.Fatalf("Failed to create directory %s: %v", path, err)
	}
}

func createFile(path string) {
	if _, err := os.Create(path); err != nil {
		log.Fatalf("Failed to create file %s: %v", path, err)
	}
}

func getFileSize(file *os.File) int {
	fileStat, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to get file stat: %v", err)
	}
	return int(fileStat.Size())
}

func validateFile(file *os.File) {
	// check if the file is a regular file
	fileStat, _ := file.Stat()
	sz := fileStat.Size()
	hasher := sha1.New()
	b := make([]byte, 1)
	for i := 0; i < int(sz)-20; i++ {
		file.Read(b)
		hasher.Write(b)
	}
	hash := hasher.Sum(nil)
	checksum := make([]byte, 20)
	file.Read(checksum)
	if hex.EncodeToString(hash) != hex.EncodeToString(checksum) {
		log.Fatalf("The file is corrupted")
	}

}
