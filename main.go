package main

import (
	"flag"
	"fmt"
	"image"
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
	var name, algorithm string
	var maxWidth, maxHeight, maxAtlases int
	flag.StringVar(&name, "name", "atlas", "the base name of the outputted atlas(s)")
	flag.StringVar(&algorithm, "packing", "growing", "the algorthim to use when packing the input files")
	flag.IntVar(&maxWidth, "width", 0, "maximum width of the output image(s)")
	flag.IntVar(&maxHeight, "height", 0, "maximum height of the output image(s)")
	flag.IntVar(&maxAtlases, "maxatlases", 0, "used to limit the number of atlases that can be generated")
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

	packer := GetPackerForAlgorithm(algorithm)
	if packer == nil {
		fmt.Fprintf(os.Stderr, "Unrecognized packing algorithm: %s\n", algorithm)
		os.Exit(2)
	}

	outDir := args[1]
	params := &GenerateParams{
		Name:       name,
		MaxWidth:   maxWidth,
		MaxHeight:  maxHeight,
		MaxAtlases: maxAtlases,
		Packer:     packer,
	}
	_, err = Generate(inFiles, outDir, params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Texture packing failed with error: %s\n", err.Error())
		os.Exit(2)
	}
}

// Includes parameters that can be passed to the Generate function
type GenerateParams struct {
	Name       string
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
	// Apply any default parameters
	if params == nil {
		params = &GenerateParams{}
	}
	if params.Name == "" {
		params.Name = "atlas"
	}
	if params.Packer == nil {
		params.Packer = PackGrowing
	}

	res = &GenerateResult{}
	res.Files = make([]*File, len(files))

	for i, filename := range files {
		// Open the given file
		r, err := os.Open(filename)
		if err != nil {
			return nil, err
		}

		decoded, _, err := image.Decode(r)
		if err != nil && err != image.ErrFormat {
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

	atlas := &Atlas{
		Name:       params.Name,
		Width:      w,
		Height:     h,
		Files:      fit,
		Descriptor: DESC_KIWI,
	}
	err = atlas.Write(outputDir)
	if err != nil {
		return nil, err
	}

	return res, nil
}
