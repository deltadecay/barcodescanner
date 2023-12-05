[Build](https://github.com/deltadecay/barcodescanner/actions/workflows/go.yml/badge.svg)

# A simple barcode scanner 

This command line tool reads images (bmp, gif, jpeg, png, tiff, webp) from the command line and outputs a json result with information whether a barcode was found or if there were errors.

## Go

To build it requires go. It has been tested with go 1.21. Perform **go mod tidy** to get dependencies. 


## Run and build

You can run without building a binary
```sh
go run . yourimagefile.jpg
```

If you want to build binaries use the Makefile and then run the executable. 
```sh
./build/barcodescanner.darwin.arm64 yourimagefile.jpg
```


## Usage
```
Usage of barcodescanner:
barcodescanner [flags] file...

This tool scans for barcodes (EAN-13 and UPC-A) in the specified files.
The argument file... is one or more image files to scan. Supported image formats
are: bmp, gif, jpeg, png, tiff, webp.

Optional flags:
  --grey
                Convert image to greyscale. Applied first.
  --scale
                Factor to resize the image with. Default 1.0 has no effect. Applied second.
  --unsharpen
                Apply unsharp mask. Four params comma separated: radius, sigma, amount, threshold. Applied third.
  --contrast
                Factor to adjust the contrast. Default 1.0 has no effect. Applied last.
  --pretty
                Pretty-print the json output
  --version
                Display version
  -h, --help
                Display this help

```


## Example

```sh
./barcodescanner.darwin.arm64 --pretty upc1.jpg dvd1_back.jpg
```
Scanning two images, one with a barcode and one where it cannot be detected generates this json:
```json
{
   "barcodes": [
      {
         "file": "upc1.jpg",
         "format": "UPC_A",
         "data": "883929215768",
         "country": "US/CA"
      },
      {
         "file": "dvd1_back.jpg",
         "error": "NotFoundException"
      }
   ]
}
```
