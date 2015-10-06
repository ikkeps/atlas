package main

import "testing"

type TestParams struct {
	Files  []string
	Params *GenerateParams
}

type TestWant struct {
	NumFiles int
}

func TestGenerate(t *testing.T) {
	OUTPUT_DIR := "./output"

	BUTTONS := []string{
		"./fixtures/button.png",
		"./fixtures/button_active.png",
		"./fixtures/button_hover.png",
	}

	cases := []struct {
		in   TestParams
		want TestWant
	}{
		{TestParams{
			Files:  BUTTONS,
			Params: nil,
		}, TestWant{
			NumFiles: len(BUTTONS),
		}},
	}
	for _, c := range cases {
		got, err := Generate(c.in.Files, OUTPUT_DIR, c.in.Params)
		if err != nil {
			t.Errorf("Generate threw an error: %s", err.Error())
		}
		if len(got.Files) != c.want.NumFiles {
			t.Errorf("Generate did not use all files: want %v files, got %v", c.want.NumFiles, got.Files)
		}
	}
}
