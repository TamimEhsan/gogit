# gogit

A simple subset of `git` made with Go to understand the internal mechanism of git. This implements the bare minimum required to initialize an empty git repository, add new files to index, commit changes and push to remote(yet to be build). 

This is supposed to be for educational purposes and in no way should be used as an alternative to the original `git`.

> This repo was created, contents added and commited using gogit

## Installation

The current implementation is linux based. So, use it from linux or wsl. Make sure you have Go installed in your system. 

At first clone this repository,
```bash
git clone https://github.com/TamimEhsan/gogit.git
```
Move inside the folder and install
```bash
cd gogit
go install
```
It should be installed inside bin folder of GOPATH. To find yours run
```bash
go env GOPATH
```
You can follow [this guideline](https://www.digitalocean.com/community/tutorials/how-to-install-go-and-set-up-a-local-programming-environment-on-ubuntu-18-04) to learn more about setting up `go install`

To check succesful installation, run
```bash
gogit version
```
If you see something like
```bash
gogit version 0.0.1
```
then you have successfully installed gogit in your system.

## Available Commands
```
These are the supported GoGit commands used in various situations:

start a working area 
   init      Create an empty Git repository or reinitialize an existing one

work on the current change 
   add       Add file contents to the index

examine the history and state 
   status    Show the working tree status

grow, mark and tweak your common history
   commit    Record changes to the repository

other utils
  cat-file   Provide content for repository objects
  ls-files   Show information about files in the index and the working tree
  version    Display version information about Gogit
```

## References
The references used to build this can be found in
- [git-scm](https://git-scm.com/)
- [pygit](https://benhoyt.com/writings/pygit/)
