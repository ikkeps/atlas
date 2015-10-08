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
	Name                string
	Files               []*File
	Width, Height       int
	MaxWidth, MaxHeight int
	Padding, Gutter     int
	Descriptor          DescriptorFormat
}

// Adds a file into the atlas at the given position
// Used by packers to add packed files to the atlas
func (a *Atlas) AddFile(file *File, x, y int) {
	file.Atlas = a
	file.X, file.Y = x, y
	a.Files = append(a.Files, file)
}

// Writes the atlas to the given output directory, this is shorthand
// for calling both WriteImage and WriteDescriptor
func (a *Atlas) Write(outputDir string) error {
	if err := a.WriteImage(outputDir); err == nil {
		return a.WriteDescriptor(outputDir)
	} else {
		return err
	}
}

// Writes the image for this atlas to the given output directory
// Returns an error if any IO operation fails
func (a *Atlas) WriteImage(outputDir string) error {
	// Generate the image data
	im := image.NewRGBA(image.Rect(0, 0, a.Width, a.Height))
	// Set the background colour of the image
	for i, n := 0, len(im.Pix); i < n; i += 4 {
		im.Pix[i] = 0   // Red
		im.Pix[i+1] = 0 // Green
		im.Pix[i+2] = 0 // Blue
		im.Pix[i+3] = 0 // Alpha
	}
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
		dp := image.Pt(file.X+a.Padding+a.Gutter, file.Y+a.Padding+a.Gutter)
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

// Writes the descriptor file for this atlas to the given output directory
// Returns an error if any IO operation fails
func (a *Atlas) WriteDescriptor(outputDir string) error {
	t, err := GetTemplateForFormat(a.Descriptor)
	if err != nil {
		return err
	}
	ext := GetFileExtForFormat(a.Descriptor)
	out, err := os.Create(path.Join(outputDir, fmt.Sprintf("%s.%s", a.Name, ext)))
	if err != nil {
		return err
	}
	return t.Execute(out, a)
}
