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
	"strconv"
	"strings"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
	"github.com/morikuni/aec"
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
	Barcodes  []*BarcodeResult `json:"barcodes"`
	Grey      bool             `json:"grey"`
	Scale     float64          `json:"scale"`
	Contrast  float64          `json:"contrast"`
	Unsharpen string           `json:"unsharpen,omitempty"`
}

func createPreProcessOps(grey bool, scaleFactor float64, unsharpenStr string, contrastFactor float64) []PreProcessOp {
	preProcessOps := make([]PreProcessOp, 0)
	if grey {
		preProcessOps = append(preProcessOps, NewGreyScaleOp())
	}
	if scaleFactor != 1.0 {
		preProcessOps = append(preProcessOps, NewResizeOp(scaleFactor))
	}

	var unsharpen []float64
	unsharpenStrFix := unsharpenStr
	unsharpenStrFix = strings.Trim(unsharpenStrFix, "'\"")
	unsharpenParams := strings.Split(unsharpenStrFix, ",")
	if len(unsharpenParams) == 4 {
		unsharpen = []float64{3, 1.0, 1.0, 0.05}
		for index, param := range unsharpenParams {
			param = strings.TrimSpace(param)
			val, err := strconv.ParseFloat(param, 64)
			if err == nil {
				unsharpen[index] = val
			}
		}
		preProcessOps = append(preProcessOps, NewUnsharpenOp(int(unsharpen[0]), unsharpen[1], unsharpen[2], unsharpen[3]))
	}
	if contrastFactor != 1.0 {
		preProcessOps = append(preProcessOps, NewContrastOp(contrastFactor))
	}
	return preProcessOps
}

func processFile(fileName string, preProcessOps []PreProcessOp) *BarcodeResult {
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

	for _, op := range preProcessOps {
		img = op.Apply(img)
	}

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

const figletStr = `
  _                     _                                 
 | |_ ___ ___ ___ ___ _| |___ ___ ___ ___ ___ ___ ___ ___ 
 | . | .'|  _|  _| . | . | -_|_ -|  _| .'|   |   | -_|  _|
 |___|__,|_| |___|___|___|___|___|___|__,|_|_|_|_|___|_|  											  										   
`

const usageStr = `Usage of barcodescanner:
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

`

func printLogo() {
	logo := aec.CyanF.Apply(figletStr)
	fmt.Println(logo)
}

func usage() {
	fmt.Fprint(os.Stderr, usageStr)
	os.Exit(2)
}

const MaxNumArgs = 100

func main() {
	flag.Usage = usage
	grey := flag.Bool("grey", false, "Convert image to greyscale.")
	scaleFactor := flag.Float64("scale", 1.0, "Factor to resize the image with. Default 1.0 has no effect.")
	contrastFactor := flag.Float64("contrast", 1.0, "Factor to adjust the contrast. Default 1.0 has no effect.")
	unsharpenStr := flag.String("unsharpen", "", "Apply unsharp mask. Four params comma separated: radius, sigma, amount, threshold.")
	prettyJson := flag.Bool("pretty", false, "Pretty-print the json output")
	displayVersion := flag.Bool("version", false, "Display version")
	flag.Parse()

	if *displayVersion {
		printLogo()
		fmt.Printf("barcodescanner v%s (%s)\n", version, buildTime)
		os.Exit(0)
	}
	args := flag.Args()

	if len(args) > MaxNumArgs {
		args = args[0:MaxNumArgs]
	}

	preProcessOps := createPreProcessOps(*grey, *scaleFactor, *unsharpenStr, *contrastFactor)

	processedFiles := make([]*BarcodeResult, len(args))
	for index, fileName := range args {
		processedFiles[index] = processFile(fileName, preProcessOps)
	}
	result := ScannedBarcodes{
		Barcodes:  processedFiles,
		Grey:      *grey,
		Scale:     *scaleFactor,
		Contrast:  *contrastFactor,
		Unsharpen: *unsharpenStr,
	}

	var (
		bytes []byte
		err   error
	)
	if *prettyJson {
		bytes, err = json.MarshalIndent(result, "", "   ")
	} else {
		bytes, err = json.Marshal(result)
	}

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes))
}
