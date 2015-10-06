package main

// Represents a single atlas to be outputted
type Atlas struct {
	Name          string
	Files         []*File
	Width, Height int
}
