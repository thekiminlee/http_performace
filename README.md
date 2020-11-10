# Description

This is a script that sends a GET requests to specified URL and returns performance summary.

## Dependencies

* go
* github.com/ogier/pflag

To install, run `go get github.com/ogier/pflag` in root directory.

## Usage

parameters:
- --url: URL to call GET request
- --profile: number of requests to make in positive int

### Running the script without performance summary:
```
    go run main.go --url="<url>" 
```
### Running the script with performance summary:
```
    go run main.go --url="<url>" --profile=<num>
```
### Build
```
    go build main.go
```
