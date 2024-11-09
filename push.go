package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"strings"
)

func gitPush(remote, userName, password string) {

	if userName == "" {
		userName = getUserName()
	}
	if password == "" {
		password = getPassword()
	}

	if remote == "" {
		panic("Remote repository not specified")
	}

	remoteHash := getRemoteMasterCommit(remote, userName, password)
	localHash := getLocalMasterCommit()

	localObjects := getObjects(localHash, remoteHash)
	localObjects = uniqueObjects(localObjects)

	remoteObjects := getObjects(remoteHash, "")
	remoteObjects = uniqueObjects(remoteObjects)

	missingObjects := getMissingObjects(localObjects, remoteObjects)
	// fmt.Println("missing objects: ", missingObjects)
	// os.Exit(0)
	line := fmt.Sprintf("%s %s refs/heads/master\x00 report-status", remoteHash, localHash)
	line = fmt.Sprintf("%04x%s\n0000", len(line)+5, line)

	pack := createPack(missingObjects)

	data := append([]byte(line), pack...)
	// send a post request to the remote repository
	url := remote + "/git-receive-pack"

	gitPushPack(url, userName, password, data)

}

func createPack(objects []string) []byte {

	header := []byte("PACK")
	header = append(header, paddInteger(2, 4)...)
	header = append(header, paddInteger(len(objects), 4)...)

	body := []byte{}
	for _, object := range objects {
		body = append(body, encodePack(object)...)
	}

	contents := append(header, []byte(body)...)

	hasher := sha1.New()
	hasher.Write(contents)
	contents = append(contents, hasher.Sum(nil)...)

	return contents
}

func encodePack(object string) []byte {
	objectType, contents := readObject(object)
	data := []byte(strings.Join(contents, "\n"))
	header := []byte{}

	enum := 0
	if objectType == "commit" {
		enum = 1
	} else if objectType == "tree" {
		enum = 2
	} else if objectType == "blob" {
		enum = 3
	}
	size := len(data)
	byt := byte(enum<<4 | size&0x0f)
	size >>= 4
	for size > 0 {
		header = append(header, byt|0x80)
		byt = byte(size & 0x7f)
		size >>= 7
	}
	header = append(header, byt)

	var buf bytes.Buffer
	writer, _ := zlib.NewWriterLevel(&buf, zlib.BestSpeed)
	_, err := writer.Write(data)
	if err != nil {
		return nil
	}
	writer.Close()

	compressedData := buf.Bytes()
	header = append(header, compressedData...)
	return header
}
