package main

// The packer type represents a packing alogrithm that can be used to
// modify file positions, sorting them into a series of atlases
type Packer func(files []*File, maxWidth int, maxHeight int)
