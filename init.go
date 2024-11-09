package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func gitInit(directory string) {
	// create a directory
	if directory == "" || directory == "." {
		directory = "."
	} else {
		createDir(directory)
	}

	// create a .git directory
	createDir(path.Join(directory, ".git"))
	// create subdirectories
	for _, dir := range []string{"branches", "hooks", "info", "objects", "refs"} {
		createDir(path.Join(directory, ".git", dir))
	}
	// create files
	for _, file := range []string{"HEAD", "config", "description" /*, "index", "packed-refs"*/} {
		createFile(path.Join(directory, ".git", file))
	}

	// create objects subdirectories
	for _, dir := range []string{"info", "pack"} {
		createDir(path.Join(directory, ".git", "objects", dir))
	}

	// create refs subdirectories
	for _, dir := range []string{"heads", "tags"} {
		createDir(path.Join(directory, ".git", "refs", dir))
	}

	// create refs subdirectories
	createFile(path.Join(directory, ".git", "info", "exclude"))

	// write to HEAD file
	headFile, _ := os.OpenFile(path.Join(directory, ".git", "HEAD"), os.O_RDWR|os.O_CREATE, 0644)
	defer headFile.Close()
	headFile.Write([]byte("ref: refs/heads/master\n"))

	// write to config file
	configFile, _ := os.OpenFile(path.Join(directory, ".git", "config"), os.O_RDWR|os.O_CREATE, 0644)
	defer configFile.Close()
	configFile.Write([]byte("[core]\n\trepositoryformatversion = 0\n\tfilemode = true\n\tbare = false\n\tlogallrefupdates = true\n"))

	// write to description file
	descriptionFile, _ := os.OpenFile(path.Join(directory, ".git", "description"), os.O_RDWR|os.O_CREATE, 0644)
	defer descriptionFile.Close()
	descriptionFile.Write([]byte("Unnamed repository; edit this file 'description' to name the repository.\n"))

	absDirectory, err := filepath.Abs(directory)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return
	}
	fmt.Println("Initialized empty Git repository in", path.Join(absDirectory, ".git"))
}
