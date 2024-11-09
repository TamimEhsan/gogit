package main

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path"
	"syscall"
)

type indexEntry struct {
	ctimeSec  int
	ctimeNsec int
	mtimeSec  int
	mtimeNsec int
	dev       int
	ino       int
	mode      int
	uid       int
	gid       int
	size      int
	sha1      []byte
	flags     int
	path      string
}

func createIndexEntry(filename, fileHash string) indexEntry {
	indexEntry := indexEntry{}
	fileInfo, err := os.Stat(filename)
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}
	stat := fileInfo.Sys().(*syscall.Stat_t)
	indexEntry.ctimeSec = int(stat.Ctim.Sec)
	indexEntry.ctimeNsec = int(stat.Ctim.Nsec)
	indexEntry.mtimeSec = int(stat.Mtim.Sec)
	indexEntry.mtimeNsec = int(stat.Mtim.Nsec)
	indexEntry.dev = int(stat.Dev)
	indexEntry.ino = int(stat.Ino)
	indexEntry.mode = int(stat.Mode)
	indexEntry.mode = setCorrectMode(indexEntry.mode)
	indexEntry.uid = int(stat.Uid)
	indexEntry.gid = int(stat.Gid)
	indexEntry.size = int(stat.Size)
	indexEntry.path = filename
	sha1, _ := hex.DecodeString(fileHash)
	indexEntry.sha1 = sha1
	indexEntry.flags = len(filename)
	return indexEntry
}

func readIndex() []indexEntry {
	file, err := os.Open(path.Join(".git", "index"))
	if err != nil {
		if os.IsNotExist(err) {
			return []indexEntry{}
		}
		log.Fatalf("Failed to open index file: %v", err)
	}
	defer file.Close()
	validateFile(file)
	file.Seek(0, 0)

	bytes := make([]byte, 4)
	file.Read(bytes)

	file.Read(bytes)
	_ = binary.BigEndian.Uint32(bytes)

	file.Read(bytes)
	entryCount := binary.BigEndian.Uint32(bytes)
	// fmt.Println("entry count: ", entryCount)

	indexes := make([]indexEntry, 0)
	for i := 0; i < int(entryCount); i++ {
		entry := indexEntry{}
		bytes = make([]byte, 4)
		file.Read(bytes)
		entry.ctimeSec = int(binary.BigEndian.Uint32(bytes))

		file.Read(bytes)
		entry.ctimeNsec = int(binary.BigEndian.Uint32(bytes))

		file.Read(bytes)
		entry.mtimeSec = int(binary.BigEndian.Uint32(bytes))

		file.Read(bytes)
		entry.mtimeNsec = int(binary.BigEndian.Uint32(bytes))

		file.Read(bytes)
		entry.dev = int(binary.BigEndian.Uint32(bytes))

		file.Read(bytes)
		entry.ino = int(binary.BigEndian.Uint32(bytes))

		file.Read(bytes)
		entry.mode = int(binary.BigEndian.Uint32(bytes))

		file.Read(bytes)
		entry.uid = int(binary.BigEndian.Uint32(bytes))

		file.Read(bytes)
		entry.gid = int(binary.BigEndian.Uint32(bytes))

		file.Read(bytes)
		entry.size = int(binary.BigEndian.Uint32(bytes))

		entry.sha1 = make([]byte, 20)
		file.Read(entry.sha1)

		bytes = make([]byte, 2)
		file.Read(bytes)
		entry.flags = int(binary.BigEndian.Uint16(bytes))

		path := make([]byte, entry.flags)
		file.Read(path)
		entry.path = string(path)
		file.Seek(1, 1)
		pad := (8 - ((62 + len(entry.path) + 1) % 8)) % 8
		// fmt.Println("pad: ", hex.EncodeToString(entry.sha1), " ", pad)
		file.Seek(int64(pad), 1)

		indexes = append(indexes, entry)
	}
	return indexes

}

func writeToIndexFile(indexEntries []indexEntry) {
	file, err := os.OpenFile(path.Join(".git", "index"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open index file: %v", err)
	}
	defer file.Close()
	file.Truncate(0)
	header := []byte("DIRC")
	header = append(header, paddInteger(2, 4)...)
	headerBytes := append([]byte(header), paddInteger(len(indexEntries), 4)...)
	file.Write(headerBytes)
	for _, entry := range indexEntries {

		file.Write(paddInteger(entry.ctimeSec, 4))
		file.Write(paddInteger(entry.ctimeNsec, 4))
		file.Write(paddInteger(entry.mtimeSec, 4))
		file.Write(paddInteger(entry.mtimeNsec, 4))
		file.Write(paddInteger(entry.dev, 4))
		file.Write(paddInteger(entry.ino, 4))
		// the mode is a 4 byte integer
		// where the last 9 bit can be only of two type 111101101 or 110100100

		file.Write(paddInteger(entry.mode, 4))
		file.Write(paddInteger(entry.uid, 4))
		file.Write(paddInteger(entry.gid, 4))
		file.Write(paddInteger(entry.size, 4))
		file.Write(entry.sha1)
		file.Write(paddInteger(entry.flags, 2))
		file.Write([]byte(entry.path))
		file.Write([]byte{0})

		pad := (8 - ((62 + len(entry.path) + 1) % 8)) % 8
		file.Write(make([]byte, pad))

	}
	file.Seek(0, 0)

	hasher := sha1.New()
	io.Copy(hasher, file)
	file.Write(hasher.Sum(nil))

}
