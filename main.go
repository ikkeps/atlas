package atlas

import (
	"errors"
	"fmt"
	"image"
	"math"
	"os"
)

// Includes parameters that can be passed to the Generate function
type GenerateParams struct {
	Name                string
	Descriptor          DescriptorFormat
	Packer              Packer
	Sorter              Sorter
	MaxWidth, MaxHeight int
	MaxAtlases          int
	Padding, Gutter     int
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
	if params.Descriptor == DESC_INVALID {
		params.Descriptor = DESC_KIWI
	}
	if params.Packer == nil {
		params.Packer = PackGrowing
	}
	if params.Sorter == nil {
		params.Sorter = GetSorterFromString(SORT_DEFAULT)
	}
	if params.MaxWidth == 0 {
		params.MaxWidth = math.MaxInt32
	}
	if params.MaxHeight == 0 {
		params.MaxHeight = math.MaxInt32
	}

	res = &GenerateResult{}
	res.Files = make([]*File, len(files))

	// The amount that will be added to the files width/height
	// by padding and gutter (we *2 to include both sides ie. top & bottom)
	border := params.Padding*2 + params.Gutter*2
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
			// Here we use padding*2 as if there is only one image it will still need
			// padding on both sides left & right in the atlas
			if size.X+border > params.MaxWidth ||
				size.Y+border > params.MaxHeight {
				return nil, errors.New(fmt.Sprintf("File %s exceeds maximum size of atlas (%dx%d)",
					filename, size.X, size.Y))
			}
			// Here we only add padding to the width and height once because otherwise
			// we will end up with double gaps between images
			res.Files[i] = &File{
				FileName: filename,
				Width:    size.X + border,
				Height:   size.Y + border,
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
			Padding:    params.Padding,
			Gutter:     params.Gutter,
		}
		res.Atlases = append(res.Atlases, atlas)
		params.Packer(atlas, pending)
		pending = getRemainingFiles(pending)
		fmt.Printf("Writing atlas named %s to %s\n", atlas.Name, outputDir)
		err = atlas.Write(outputDir)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func getRemainingFiles(files []*File) (remaining []*File) {
	remaining = make([]*File, 0)
	for _, file := range files {
		if file.Atlas == nil {
			remaining = append(remaining, file)
		}
	}
	return remaining
}
