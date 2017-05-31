# Contributing

Welcome! If you're looking to help, this document is a great place to start! 


## Finding things that need help

Here's a few places to get started and find out what's outstanding.

- Read through the [MDM Protocol Reference](https://developer.apple.com/library/content/documentation/Miscellaneous/Reference/MobileDeviceManagementProtocolRef/3-MDM_Protocol/MDM_Protocol.html) on the Apple website. Having a deeper understanding of MDM can help with designing features and uncovering bugs. 
- Follow the [Quickstart](https://github.com/micromdm/micromdm/wiki/Quickstart) guide and make edits if something doesn't look or work right. 
- If you run into a problem that you're not sure how to fix, file a bug in [the issue tracker](https://github.com/micromdm/micromdm/issues)
- Browse through the open issues in [the issue tracker](https://github.com/micromdm/micromdm/issues). We try to tag issues as [**beginner friendly**](https://github.com/micromdm/micromdm/issues?q=is%3Aissue+is%3Aopen+label%3Abeginner-friendly) where appropriate.
- See something that others might benefit from? Considering updating or writing a [wiki page](https://github.com/micromdm/micromdm/wiki).

## Building the project

To build MicroMDM from source, you will need [Go 1.8](https://golang.org/dl/) or later installed.

## Once you have Go

MicroMDM uses the go lang `dep` tool for vendor management. 
Use `which dep` to verify you have it installed and in your PATH.
If `dep` is not installed please review the [If you're new to Go](#if-youre-new-to-go) section for install steps.

1. `go get github.com/micromdm/micromdm`
2. `cd $GOPATH/src/github.com/micromdm/micromdm`
3. `dep ensure` install the necessary dependencies into  the `/vendor` folder
4. `go build` or `go install`

## If you're new to Go

Go is a bit different from other languages in its requirements for how it expects its programmers to organize Go code files in directories.
First, Go requires a folder, called a workspace (you can name it anything you'd like) to exist for go source, dependencies, etc. Before Go 1.8 the path to this folder must always be set in the environment variable `GOPATH` (example: `export GOPATH=/Users/groob/code/go`). As of Go 1.8 the default `GOPATH` is set to `$HOME/go` but you can still set it to whatever you like.
Your `GOPATH` must have thee subfolders: `bin`, `pkg`, and `src`. Any code you create must live inside the `src` folder. It's also helpful to add `$GOPATH/bin` to your environment's `PATH` as that is where `go install` will place go binaries that you build. This makes it so that binaries that are insalled can just be invoked by name rather than their full page.

A few helpful resources for getting started with Go:

* [Writing, building, installing, and testing Go code](https://www.youtube.com/watch?v=XCsL89YtqCs)
* [Resources for new Go programmers](http://dave.cheney.net/resources-for-new-go-programmers)
* [How I start](https://howistart.org/posts/go/1)
* [How to write Go code](https://golang.org/doc/code.html)
* [GOPATH on the go wiki](https://github.com/golang/go/wiki/GOPATH)

To build MicroMDM you will need to:  

1. Download and install [`Go`](https://golang.org/dl/)  
2. Make a workspace directory and set the `GOPATH` as explained above.
3. Install [`dep`](https://github.com/golang/dep) via the command `go get -u github.com/golang/dep/...`
Note that `dep` is a very new project itself. If you're running trouble with the `dep ensure` command, ping @groob in the #micromdm channel on Slack.
4. `mkdir -p $GOPATH/src/github.com/micromdm`
5. `git clone git@github.com:micromdm/micromdm.git` the project into the above folder.
The repo must always be in the folder `$GOPATH/src/github.com/micromdm/micromdm` even if you forked the project. Add a git remote to your fork to track upstream.
6. `dep ensure` The `dep` command will download and install all necessary dependencies for the project to compile.
7. `go build` or `go install`
8. File an issue or a pull request if the instructions are unclear or problems pop-up for you.

## Important libraries and frameworks

MicroMDM is built using a few popular Go packages outside the standard libraries. It might be worth checking them out. 

- [Go Kit](https://github.com/go-kit/kit#go-kit------) is a set of Go libraries used by MicroMDM to provide [logging](https://github.com/go-kit/kit/tree/master/log), and abstractions for building HTTP services. Its [examples](https://gokit.io/examples/) page is a good place to start.  
- [BoltDB](https://github.com/boltdb/bolt#getting-started) is a key/value database used to provide persistant storage for many components of MicroMDM.
- [gorilla/mux](http://www.gorillatoolkit.org/pkg/mux) is used to provide routing for http handlers. 

## Other resources

Also see [Contributing wiki page](https://github.com/micromdm/micromdm/wiki/Contributing) which has some additional notes on running, troubleshooting, and developing with MicroMDM.
