# Contributing

Welcome, if you're looking to help, this document is a great place to start! 


# Finding Things That Need Help

If you're looking to help, this is a great place to start. 
- Read through the [MDM Protocol Reference](https://developer.apple.com/library/content/documentation/Miscellaneous/Reference/MobileDeviceManagementProtocolRef/3-MDM_Protocol/MDM_Protocol.html) on the Apple website. Having a deeper understanding of MDM can help with designing features and uncovering bugs. 
- Follow the [Quickstart](https://github.com/micromdm/micromdm/wiki/Quickstart) guide and make edits if something doesn't look right. 
- If you run into a problem that you're not sure how to fix, file a bug. 
- Browse through the open issues. We try to tag issues as [**beginner friendly**](https://github.com/micromdm/micromdm/issues?q=is%3Aissue+is%3Aopen+label%3Abeginner-friendly) where appropriate.

# Building the project

## If you have Go

1. `go get github.com/micromdm/micromdm`
2. `cd $GOPATH/src/github.com/micromdm/micromdm`
3. `dep ensure` install the necessary dependencies into /vendor folder
4. `go build` or `go install`

## If you're new to Go

Go is a bit different from other languages in its requirements for how it expects its programmers to organize Go code on a system.
First, Go expects you to choose a folder, called a workspace (you can name it anything you'd like). The path to this folder must always be set in an environment variable - `GOPATH` (example: `GOPATH=/Users/groob/code/go`)  
Your `GOPATH` must have thee subfolders - `bin`, `pkg` and `src`, and any code you create must live inside the `src` folder. It's also helpful to add `$GOPATH/bin` to your environment's `PATH` as that is where `go install` will place go binaries that you build.

A few helpful resources for getting started with Go.

* [Writing, building, installing, and testing Go code](https://www.youtube.com/watch?v=XCsL89YtqCs)
* [Resources for new Go programmers](http://dave.cheney.net/resources-for-new-go-programmers)
* [How I start](https://howistart.org/posts/go/1)
* [How to write Go code](https://golang.org/doc/code.html)
* [GOPATH - go wiki page](https://github.com/golang/go/wiki/GOPATH)

To build MicroMDM you will need to:  

1. Download and install [`Go`](https://golang.org/dl/)  
2. Install [`dep`](https://github.com/golang/dep)
3. Set the `GOPATH` as explained above.
4. `mkdir -p $GOPATH/src/github.com/micromdm`
5. `git clone` the project into the above folder.  
The repo must always be in the folder `$GOPATH/src/github.com/micromdm/micromdm` even if you forked the project. Add a git remote to your fork.  
6. `dep ensure` The `dep` command will download and install all necessary dependencies for the project to compile.
7. `go build` or `go install`
8. File an issue or a pull request if the instructions were unclear.
