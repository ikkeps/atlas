package main

import "testing"

func TestPackGrowing(t *testing.T) {
	FILES := []*File{
		&File{Width: 200, Height: 200},
		&File{Width: 100, Height: 100},
		&File{Width: 50, Height: 50},
	}

	EXPECT := []*File{
		&File{X: 0, Y: 0},
		&File{X: 200, Y: 0},
		&File{X: 200, Y: 100},
	}

	ATLAS := &Atlas{
		MaxWidth:  400,
		MaxHeight: 400,
	}
	unfit := PackGrowing(ATLAS, FILES)
	if len(unfit) > 0 {
		t.Errorf("Failed to fit %d file(s)", len(unfit))
	}
	if ATLAS.Width != 300 || ATLAS.Height != 200 {
		t.Errorf("Unexpected atlas size: %d,%d", ATLAS.Width, ATLAS.Height)
	}
	for i, file := range FILES {
		expect := EXPECT[i]
		if file.X != expect.X || file.Y != expect.Y {
			t.Errorf("File position %d,%d does not match expected position %d,%d",
				file.X, file.Y, expect.X, expect.Y)
		}
	}
}
