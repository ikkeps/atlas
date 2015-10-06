package main

import (
	"flag"
	"fmt"
	"image"
	//"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

// Shows the help info for command line usage
func showHelp() {
	fmt.Fprintf(os.Stderr, "usage: %s [-params] <infiles> <outdir>\n", os.Args[0])
	os.Exit(2)
}

// Handles command line usage
func main() {
	var width, height int
	flag.IntVar(&width, "width", 2048, "maximum width of the output image(s)")
	flag.IntVar(&height, "height", 2048, "maximum height of the output image(s)")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		showHelp()
	}

	inGlob := args[0]
	inFiles, err := filepath.Glob(inGlob)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(2)
	}

	outDir := args[1]
	_, err = Generate(inFiles, outDir, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Texture packing failed with error: %s\n", err.Error())
		os.Exit(2)
	}
}

// Includes parameters that can be passed to the Generate function
type GenerateParams struct {
	MaxWidth   int
	MaxHeight  int
	MaxAtlases int
	Packer     Packer
}

// Includes details of the result of a texture atlas Generate request
type GenerateResult struct {
	Files []*File
}

// Generates a series of texture atlases using the given files as input
// and outputting to the given directory with the given parameters.
// Will generate an error if any IO operations fail or if the GenerateParams
// represent an invalid configuration
func Generate(files []string, outputDir string, params *GenerateParams) (res *GenerateResult, err error) {
	fmt.Printf("Files: %v\n", files)

	// Apply any default parameters
	if params == nil {
		params = &GenerateParams{}
	}
	if params.Packer == nil {
		params.Packer = PackGrowing
	}

	res = &GenerateResult{}
	res.Files = make([]*File, len(files))

	for i, filename := range files {
		fmt.Printf("Found file: %s\n", filename)

		// Open the given file
		r, err := os.Open(filename)
		if err != nil {
			return nil, err
		}

		decoded, _, err := image.Decode(r)
		if err != nil && err != image.ErrFormat {
			fmt.Printf("Failed to open file")
			return nil, err
		}

		if err != image.ErrFormat {
			size := decoded.Bounds().Size()
			res.Files[i] = &File{
				FileName: filename,
				Width:    size.X,
				Height:   size.Y,
			}
		} else {
			fmt.Printf("Incorrect format for file: %s\n", filename)
		}
	}

	fit, w, h := params.Packer(res.Files, params.MaxWidth, params.MaxHeight)
	res.Files = fit

	fmt.Printf("%s\n", fit)
	fmt.Printf("%d,%d\n", w, h)
	return res, nil
}
