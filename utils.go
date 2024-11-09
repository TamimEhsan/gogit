package main

import (
	"encoding/binary"
	"os"
)

func getUserName() string {
	name := os.Getenv("GOGIT_USERNAME")
	if name == "" {
		panic("Username not set")
	}
	return name
}

func getPassword() string {
	pass := os.Getenv("GOGIT_PASSWORD")
	if pass == "" {
		panic("Password not set")
	}
	return pass
}

func getEmail() string {
	email := os.Getenv("GOGIT_EMAIL")
	if email == "" {
		panic("Email not set")
	}
	return email
}

func paddInteger(n int, size int) []byte {
	if size == 2 {
		b := make([]byte, size)
		binary.BigEndian.PutUint16(b, uint16(n))
		return b
	}

	b := make([]byte, size)
	binary.BigEndian.PutUint32(b, uint32(n))
	return b
}

func setCorrectMode(mode int) int {
	const mask = 0x1FF // Mask to extract the last 9 bits
	const validMode1 = 0x1ED
	const validMode2 = 0x1A4
	mode = mode &^ mask
	return mode | validMode1
}
