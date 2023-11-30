package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
	"github.com/nfnt/resize"
	"hawx.me/code/img/greyscale"
)

var (
	buildTime string = "1970-01-01T00:00:00UTC"
	version   string = "0.0dev"
)

type BarcodeResult struct {
	FileName string `json:"file"`
	Format   string `json:"format,omitempty"`
	Data     string `json:"data,omitempty"`
	Country  string `json:"country,omitempty"`
	Error    string `json:"error,omitempty"`
}

type ScannedBarcodes struct {
	Barcodes []*BarcodeResult `json:"barcodes"`
}

func processFile(fileName string) *BarcodeResult {
	output := &BarcodeResult{
		FileName: fileName,
	}

	file, err := os.Open(fileName)
	if err != nil {
		output.Error = err.Error()
		return output
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		output.Error = err.Error()
		return output
	}

	img = greyscale.Greyscale(img)
	//img = sharpen.UnsharpMask(img, 4, 1.0, 1.0, 0.05)
	//img = contrast.Linear(img, 1.5)

	width := uint(img.Bounds().Max.X - img.Bounds().Min.X)
	img = resize.Resize(width*2, 0, img, resize.Lanczos3)

	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	hints := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_POSSIBLE_FORMATS: []gozxing.BarcodeFormat{
			gozxing.BarcodeFormat_EAN_13,
			gozxing.BarcodeFormat_UPC_A,
		},
		gozxing.DecodeHintType_TRY_HARDER:    true,
		gozxing.DecodeHintType_ALSO_INVERTED: true,
	}
	reader := oned.NewMultiFormatUPCEANReader(hints)
	result, err := reader.Decode(bmp, hints)
	if err != nil {
		output.Error = err.Error()
		return output
	}

	output.Format = result.GetBarcodeFormat().String()
	output.Data = result.GetText()

	metaData := result.GetResultMetadata()
	if val, found := metaData[gozxing.ResultMetadataType_POSSIBLE_COUNTRY]; found {
		possibleCountry := val.(string)
		output.Country = possibleCountry
	}
	return output
}

const usageStr = `Usage of barcodescanner:
barcodescanner [flags] file...

This tool scans for barcodes (EAN-13 and UPC-A) in the specified files.
The argument file... is one or more image files to scan. Supported image formats
are: bmp, gif, jpeg, png, tiff, webp.

Optional flags:
  --pretty	
		Pretty-print the json output
  --version
		Display version
  -h, --help
		Display this help

`

func usage() {
	fmt.Fprint(os.Stderr, usageStr)
	os.Exit(2)
}

const MaxNumArgs = 100

func main() {
	flag.Usage = usage
	prettyJson := flag.Bool("pretty", false, "Pretty-print the json output")
	displayVersion := flag.Bool("version", false, "Display version")
	flag.Parse()

	if *displayVersion {
		fmt.Printf("barcodescanner v%s (%s)\n", version, buildTime)
		os.Exit(0)
	}
	//args := os.Args[1:]
	args := flag.Args()

	if len(args) > MaxNumArgs {
		args = args[0:MaxNumArgs]
	}

	processedFiles := make([]*BarcodeResult, len(args))
	for index, arg := range args {
		fileName := arg
		result := processFile(fileName)
		processedFiles[index] = result
	}
	scannedFiles := ScannedBarcodes{Barcodes: processedFiles}

	var (
		bytes []byte
		err   error
	)
	if *prettyJson {
		bytes, err = json.MarshalIndent(scannedFiles, "", "   ")
	} else {
		bytes, err = json.Marshal(scannedFiles)
	}

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes))
}
