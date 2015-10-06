package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"path"
)

// Represents a single atlas to be outputted
type Atlas struct {
	Name          string
	Files         []*File
	Width, Height int
}

// Writes the image for this atlas to the given output directory
// Returns an error if any IO operation fails
func (a *Atlas) WriteImage(outputDir string) error {
	// Generate the image data
	im := image.NewRGBA(image.Rect(0, 0, a.Width, a.Height))
	for _, file := range a.Files {
		// Open the given file
		r, err := os.Open(file.FileName)
		if err != nil {
			return err
		}
		// Decode the image
		cim, _, err := image.Decode(r)
		if err != nil {
			return err
		}
		dp := image.Pt(file.X, file.Y)
		draw.Draw(im, image.Rectangle{dp, dp.Add(cim.Bounds().Size())}, cim, cim.Bounds().Min, draw.Src)
	}

	out, err := os.Create(path.Join(outputDir, fmt.Sprintf("%s.png", a.Name)))
	if err != nil {
		return err
	}

	err = png.Encode(out, im)
	if err != nil {
		return err
	}
	return nil
}
