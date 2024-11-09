package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
)

func main() {

	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	hashObjectCmd := flag.NewFlagSet("hash-object", flag.ExitOnError)
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	catFilesCmd := flag.NewFlagSet("cat-file", flag.ExitOnError)
	pushCmd := flag.NewFlagSet("push", flag.ExitOnError)

	password := pushCmd.String("p", "", "The password for the remote repository")
	userName := pushCmd.String("u", "", "The username for the remote repository")
	remote := pushCmd.String("r", "", "The remote repository")

	// statusCmd := flag.NewFlagSet("status", flag.ExitOnError)

	commitCmd := flag.NewFlagSet("commit", flag.ExitOnError)

	// diffCmd := flag.NewFlagSet("diff", flag.ExitOnError)

	objectType := hashObjectCmd.String("t", "blob", "The type of the object")

	switch os.Args[1] {
	case "init":
		initCmd.Parse(os.Args[2:])
		directory := initCmd.Arg(0)
		gitInit(directory)
	case "hash-object":
		hashObjectCmd.Parse(os.Args[2:])
		file := hashObjectCmd.Arg(0)
		gitHashObject(file, *objectType)
	case "add":
		addCmd.Parse(os.Args[2:])
		files := addCmd.Args()
		gitAdd(files)
	case "ls-files":
		indexes := readIndex()
		for _, entry := range indexes {
			fmt.Println(entry.path, " ", entry.size, " ", hex.EncodeToString(entry.sha1))
		}
	case "cat-file":
		catFilesCmd.Parse(os.Args[2:])
		file := catFilesCmd.Arg(0)
		gitCatFile(file)
	case "status":
		gitStatus()
	case "tree":
		createTree()
	case "commit":
		commitCmd.Parse(os.Args[2:])
		msg := commitCmd.Arg(0)
		createCommit(msg)
	case "push":
		pushCmd.Parse(os.Args[2:])
		gitPush(*remote, *userName, *password)
	case "version":
		fmt.Println("gogit version 0.0.1")
	}

}
