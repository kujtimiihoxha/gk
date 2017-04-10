# Go-Kit generator.
### This generator is still a work in progress
Go-kit generator is a cli application that generates boilerplate code for your go-kit services.

## Why?

**Because I'm lazy**, and because it would make it easier for go-kit newcomers to start using it.

## Installation
```bash
go install github.com/kujtimiihoxha/gk
```
## Running the generator
`gk` must be run from a project inside the specified `$GOPATH` for it to work.
When it is run for the first time it will search for `gk.json` configuration file, if it does not find it it will create one
with the default settings.

#### Create a new service
Inside the project run:
```bash
gk new service hello
```
or the shorter command:
```bash
gk n s hello
```
this will create a new service called `HelloService` inside :
```
project
└───hello
│   └───pkg
│   │   └───service
│   │   │    service.go
```
**service.go**
```go
package service

type HelloService interface {
	// Write your interface methods
}
```
Now you need to add the interface methods and initiate your service:
e.x:
```go
package service
import "context"
type HelloService interface {
	// Write your interface methods
	World(ctx context.Context, s string) (rs string, err error)
}
```
than run : 
```bash
gk init hello
```
this will create the service `struct` , `methods`, `endpoints`, `transport` .
The final folder structure is the same as  [addsvc](https://github.com/peterbourgon/go-microservices/tree/master/addsvc) 

## I don't like the folder structure!

The folder structure that the generator is using is following https://github.com/go-kit/kit/issues/70 but 
that can be changed using `gk.json` all the paths are configurable there.

## Cli Help
Every command has the `-h` or `--help` flag this will give you more info on what the command does and how to use it.
e.x 
```bash
gk init -h
```
will return
```bash
Initiates a service

Usage:
  gk init [flags]

Flags:
  -t, --transport string   Specify the transport you want to initiate for the service

Global Flags:
  -d, --debug           If you want to se the debug logs.
      --folder string   If you want to specify the base folder of the project.
  -f, --force           Force overide existing files without asking.
      --testing         If testing the generator.

```
## What is working
The example you see here  https://github.com/go-kit/kit/issues/70

## TODO-s

 - Implement support for other transports besides `http`
 - Implement the update commands, this commands would be used to update an existing service e.x add 
 a new request parameter to an endpoint.
 - Implement middleware generator (service,endpoint).
 - Implement automatic creation of the service main file.
 - Tests tests tests ...
## Warnings

- I only tested this on the mac, should work on other os-s but I have not tested it, I would appreciate feedback on this. 
- Commands may change in the future, this project is still a work in progress.
## Contribute
Thanks a lot for contributing. 

To test your new features/bug-fixes you need a way to run `gk` inside your project this can be done using `test_dir`.

Execute this in your command line :
```bash
export GK_FOLDER="test_dir" 
```
Create a folder in the `gk` repository called `test_dir`, now every time you run `go run main.go [anything]`
`gk` will treat `test_dir` as the project root.
 
