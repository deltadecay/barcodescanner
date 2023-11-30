# A simple barcode scanner 

This command line tool reads images (gif, jpg, png) from the command line and outputs a json result 
with information whether a barcode was found or if there were errors.

## Go

To build it requires go. It has been tested with go 1.21. Perform **go mod tidy** to get dependencies. 


## Run and build

You can run without building a binary
```
go run main.go yourimagefile.jpg
```

If you want to build binaries use the Makefile and then run the executable. 
```
./build/barcodescanner.darwin.arm64 yourimagefile.jpg
```
