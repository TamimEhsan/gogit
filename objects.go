package main

import (
	"bufio"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strings"
)

func getMissingObjects(localObjects, remoteObjects []string) []string {
	i, j := 0, 0
	objects := []string{}
	for i < len(localObjects) && j < len(remoteObjects) {
		if localObjects[i] == remoteObjects[j] {
			i++
			j++
		} else if localObjects[i] < remoteObjects[j] {
			// push the object
			objects = append(objects, localObjects[i])
			i++
		} else {
			// pull the object
			// fmt.Println("pulling: ", remoteObjects[j])
			j++
		}
	}

	for i < len(localObjects) {
		objects = append(objects, localObjects[i])
		i++
	}

	for j < len(remoteObjects) {
		// fmt.Println("pulling: ", remoteObjects[j])
		j++
	}

	return objects
}

func uniqueObjects(objects []string) []string {
	sort.Strings(objects)
	// keep only unique objects
	uniqueObjects := []string{}
	for i := 0; i < len(objects); i++ {
		if i == 0 || objects[i] != objects[i-1] {
			uniqueObjects = append(uniqueObjects, objects[i])
		}
	}
	return uniqueObjects

}

func getObjects(commit, until string) []string {

	objects := []string{}
	if commit == "" || commit == "0000000000000000000000000000000000000000" {
		return objects
	}
	objects = append(objects, commit)

	tree, parent := getMetaObjectsOfCommit(commit)

	objects = append(objects, tree)
	objects = append(objects, readTree(tree)...)

	if parent != "" && parent != "0000000000000000000000000000000000000000" {
		objects = append(objects, getObjects(parent, until)...)
	}
	return objects

}

func readObject(object string) (string, []string) {
	objectFile, err := os.Open(path.Join(".git", "objects", object[:2], object[2:]))

	if err != nil {
		log.Fatalf("Failed to open object file: %v", err)
	}
	defer objectFile.Close()
	zlibReader, _ := zlib.NewReader(objectFile)
	defer zlibReader.Close()

	bufferedReader := bufio.NewReader(zlibReader)
	var contents []string
	for {
		line, err := bufferedReader.ReadString('\n')
		line = strings.TrimSuffix(line, "\n")
		contents = append(contents, line)

		if err != nil {
			break
		}
	}

	nullIndex := strings.Index(contents[0], "\x00")
	if nullIndex == -1 {
		log.Fatalf("Invalid object format")
	}
	objectType := contents[0][:nullIndex]
	contents[0] = contents[0][nullIndex+1:]
	objectType = strings.Split(objectType, " ")[0]
	return objectType, contents
}

func getMetaObjectsOfCommit(commit string) (string, string) {

	objectType, contents := readObject(commit)
	if objectType != "commit" {
		log.Fatalf("Invalid commit object type %v: %v", commit, objectType)
	}

	tree, parent := "", ""
	for line := range contents {
		if strings.HasPrefix(contents[line], "tree") {
			tree = strings.Split(contents[line], " ")[1]
		} else if strings.HasPrefix(contents[line], "parent") {
			parent = strings.Split(contents[line], " ")[1]
		}
	}

	return tree, parent
}

/**
 * writeToObject is a function that takes a reader, a file hash, an object type and a size
 * and writes the object to the .git/objects directory
 * The object is written in the form "<type> <size>\x00\<content>"
 */

func writeToObjectFile(reader io.Reader, fileHash string, objectType string, sz int) {

	createDir(path.Join(".git", "objects", fileHash[:2]))
	objectFile, err := os.Create(path.Join(".git", "objects", fileHash[:2], fileHash[2:]))
	if err != nil {
		log.Fatalf("Failed to create object file: %v", err)
	}
	defer objectFile.Close()

	zlibWriter, _ := zlib.NewWriterLevel(objectFile, zlib.DefaultCompression)
	defer zlibWriter.Close()

	header := fmt.Sprintf("%s %d\x00", objectType, sz)
	zlibWriter.Write([]byte(header))
	io.Copy(zlibWriter, reader)
	zlibWriter.Flush()
}

func gitHashObject(filename string, objectType string) {

	var file *os.File
	err := error(nil)
	if filename == "" || filename == "-" {
		file, err = os.CreateTemp("", "buffered-content-")
		if err != nil {
			log.Fatalf("failed to create temporary file: %v", err)
		}
		defer os.Remove(file.Name())
		io.Copy(file, os.Stdin)
		file.Seek(0, 0)
	} else {
		file, err = os.Open(filename)
		if err != nil {
			log.Fatalf("Failed to open file %s: %v", filename, err)
		}
		defer file.Close()
	}

	sz := getFileSize(file)
	hash := hashObject(file, objectType, sz)
	fmt.Println(hash)
}

/**
 * hashObject is a function that takes a file and an object type and
 * returns the sha1 hash of the object in form "<type> <size>\x00\<content>""
 */
func hashObject(reader io.Reader, objectType string, sz int) string {

	header := fmt.Sprintf("%s %d\x00", objectType, sz)
	hasher := sha1.New()
	hasher.Write([]byte(header))

	io.Copy(hasher, reader)

	return hex.EncodeToString(hasher.Sum(nil))

}
