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

To build MicroMDM from source, you will need [Go 1.11](https://golang.org/dl/) or later installed.

```
git clone git@github.com:micromdm/micromdm && cd micromdm
make deps
make

# run
./build/darwin/micromdm -h
./build/darwin/mdmctl -h
```

## Git workflow
```
username=groob
# add your remote/upstream
git remote add $username git@github.com:groob/micromdm.git

# update from origin/master
git pull --rebase

# create a branch
git checkout -b my_feature

# push changes from my_feature to your fork.
#    -u, --set-upstream    set upstream for git pull/status
git push -u $username
```


## Go Resources

A few helpful resources for getting started with Go:

* [Writing, building, installing, and testing Go code](https://www.youtube.com/watch?v=XCsL89YtqCs)
* [Resources for new Go programmers](http://dave.cheney.net/resources-for-new-go-programmers)
* [How I start](https://howistart.org/posts/go/1)
* [How to write Go code](https://golang.org/doc/code.html)

## Important libraries and frameworks

MicroMDM is built using a few popular Go packages outside the standard libraries. It might be worth checking them out.

- [Go Kit](https://github.com/go-kit/kit#go-kit------) is a set of Go libraries used by MicroMDM to provide [logging](https://github.com/go-kit/kit/tree/master/log), and abstractions for building HTTP services. Its [examples](https://gokit.io/examples/) page is a good place to start.
- [BoltDB](https://github.com/boltdb/bolt#getting-started) is a key/value database used to provide persistant storage for many components of MicroMDM.
- [gorilla/mux](http://www.gorillatoolkit.org/pkg/mux) is used to provide routing for http handlers.

## Other resources

Also see [Contributing wiki page](https://github.com/micromdm/micromdm/wiki/Contributing) which has some additional notes on running, troubleshooting, and developing with MicroMDM.
