package main

import (
	"fmt"
	"testing"
)

func TestPackGrowing(t *testing.T) {
	FILES := []*File{
		&File{Width: 200, Height: 200},
		&File{Width: 100, Height: 100},
		&File{Width: 50, Height: 50},
	}

	PackGrowing(FILES, 400, 400)
	for _, file := range FILES {
		fmt.Println(*file)
	}
}
