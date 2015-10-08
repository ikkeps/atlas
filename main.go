package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"math"
	"os"
)

// Shows the help info for command line usage
func showHelp() {
	fmt.Fprintf(os.Stderr, "usage: %s [-params] <infiles> <outdir>\n", os.Args[0])
	os.Exit(2)
}

// Handles command line usage
func main() {
	var name, algorithm, sorter string
	var maxWidth, maxHeight, maxAtlases int
	flag.StringVar(&name, "name", "atlas", "the base name of the outputted atlas(s)")
	flag.StringVar(&algorithm, "packing", "growing", "the algorthim to use when packing the input files")
	flag.StringVar(&sorter, "sort", SORT_DEFAULT, "the sorting method to use when ordering the files")
	flag.IntVar(&maxWidth, "width", 0, "maximum width of the output image(s)")
	flag.IntVar(&maxHeight, "height", 0, "maximum height of the output image(s)")
	flag.IntVar(&maxAtlases, "maxatlases", 0, "used to limit the number of atlases that can be generated")
	flag.Parse()

	args := flag.Args()
	inFiles := args[:len(args)-1]
	outDir := args[len(args)-1]

	if len(args) < 2 {
		showHelp()
	}

	packer := GetPackerForAlgorithm(algorithm)
	if packer == nil {
		fmt.Fprintf(os.Stderr, "Unrecognized packing algorithm: %s\n", algorithm)
		os.Exit(2)
	}

	params := &GenerateParams{
		Name:       name,
		MaxWidth:   maxWidth,
		MaxHeight:  maxHeight,
		MaxAtlases: maxAtlases,
		Packer:     packer,
		Sorter:     GetSorterFromString(sorter),
	}
	_, err := Generate(inFiles, outDir, params)
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
	Sorter     Sorter
	Descriptor DescriptorFormat
}

// Includes details of the result of a texture atlas Generate request
type GenerateResult struct {
	Files   []*File
	Atlases []*Atlas
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
	if params.Descriptor == DESC_INVALID {
		params.Descriptor = DESC_KIWI
	}
	if params.MaxWidth == 0 {
		params.MaxWidth = math.MaxInt32
	}
	if params.MaxHeight == 0 {
		params.MaxHeight = math.MaxInt32
	}
	if params.Sorter == nil {
		params.Sorter = GetSorterFromString(SORT_DEFAULT)
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
			if size.X > params.MaxWidth || size.Y > params.MaxHeight {
				return nil, errors.New(fmt.Sprintf("File %s exceeds maximum size of atlas (%dx%d)",
					filename, size.X, size.Y))
			}

			res.Files[i] = &File{
				FileName: filename,
				Width:    size.X,
				Height:   size.Y,
			}
		} else {
			fmt.Printf("Incorrect format for file: %s\n", filename)
		}
	}

	if len(res.Files) == 0 {
		fmt.Printf("No files to pack\n")
		return res, nil
	}

	res.Atlases = make([]*Atlas, 0)

	pending := params.Sorter(res.Files)
	for i := 0; len(pending) > 0; i++ {
		atlas := &Atlas{
			Name:       fmt.Sprintf("%s-%d", params.Name, (i + 1)),
			MaxWidth:   params.MaxWidth,
			MaxHeight:  params.MaxHeight,
			Descriptor: DESC_KIWI,
		}
		res.Atlases = append(res.Atlases, atlas)
		pending = params.Packer(atlas, pending)
		fmt.Printf("Writing atlas named %s to %s\n", atlas.Name, outputDir)
		err = atlas.Write(outputDir)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
