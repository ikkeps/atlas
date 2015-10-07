package main

// Available agorithms for packing
const (
	PACK_GROWING = "growing"
)

// The packer type represents a packing alogrithm that can be used to
// modify file positions, sorting them into a series of atlases
type Packer func(atlas *Atlas, files []*File) (unfit []*File)

// Returns the packer function for the given alorithm
// Will return nil if the algorithm is not recognised
func GetPackerForAlgorithm(algorithm string) Packer {
	switch algorithm {
	case PACK_GROWING:
		return PackGrowing
	default:
		return nil
	}
}
