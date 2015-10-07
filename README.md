Go Atlas
========

Texture packer written in Go.

### Example Usage

Basic example
```
inFiles := []string{
	"./assets/sprite1.png",
	"./assets/sprite2.png"
}
outputDir := "./assets/spritesheets"
res, err := atlas.Generate(inFiles, outputDir, nil)
```

You can pass params to the Generate function;
```
params := atlas.GenerateParams {
	Name   	   : "atlas" // The base name of the outputted files
	MaxWidth   : 2048 // Maximum width/height of the atlas images
	MaxHeight  : 2048 
	MaxAtlases : 0 // Indicates no maximum
	Packer     : "growing" // The algorithm to use when packing
	Descriptor : "kiwi" // The format of the data file for the atlases
}
res, err := atlas.Generate(inFiles, outputDir, &params)
```