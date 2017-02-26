# ignoreit <img align="center" height="50" src="https://cdn.rawgit.com/whoshuu/assets/940f4518bd4a20cdac9f7bde676a722af8b96753/ignore.svg"> Declarative .gitignore

[![Build Status](https://travis-ci.org/whoshuu/ignoreit.svg?branch=master)](https://travis-ci.org/whoshuu/ignoreit) [![Go Report Card](https://goreportcard.com/badge/github.com/whoshuu/ignoreit)](https://goreportcard.com/report/github.com/whoshuu/ignoreit) [![codecov](https://codecov.io/gh/whoshuu/ignoreit/branch/master/graph/badge.svg)](https://codecov.io/gh/whoshuu/ignoreit)

`ignoreit` is a utility to produce declarative `.gitignore` specifications that can be used to generate actual `.gitignore` files. In fact, this [.gitignore](https://github.com/whoshuu/ignoreit/blob/master/.gitignore) was produced like this:

```bash
ignoreit add Go
ignoreit generate
```

The first command produces a schema `.ignoreit.yml` that abstracts the details of a set of `.gitignore` patterns into a human readable specification. It takes this form:

```yml
sources:
- repo: github/gitignore
  branch: master
  entries:
  - Go
- repo: whoshuu/gitignore
  branch: develop
  entries:
  - C++
  - Python 
custom:
- .custompattern
- .anothercustompattern
schema_version: 1
```

The second command takes this specification and generates a corresponding `.gitignore` file from it.

Both files should be checked into source control. Only the first should be manually edited via `ignoreit`, and the second is simply an artifact of changing the schema.

## Install

You can install the `ignoreit` binary like how you would install any Go program:

```bash
go get github.com/whoshuu/ignoreit
```

Since the dependencies are vendored, this should pull in a single package and binary in your `${GOPATH}/bin` called `ignoreit`.

## Commands

In addition to `ignoreit add` and `ignoreit generate`, there is `ignoreit remove`. This will remove the patterns specified in the argument.

Both `add` and `remove` take an arbitrary number of arguments, so multiple entries can be specified at once:

```
ignoreit add Go CMake C++ Python
ignoreit remove CMake C++
```

These commands take `--repo` and `--branch` flags for specifying the source repository and branch to use for pulling down `.gitignore` entries. By default these are `github/gitignore` and `master` respectively.

Finally, `ignoreit generate` should be run any time changes are made to `.ignoreit.yml`. This command takes no arguments and simply inflates the specification into an appropriate `.gitignore`.
