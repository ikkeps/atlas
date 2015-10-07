package main

import "testing"

var (
	file_50x50   = &File{FileName: "50x50", Width: 50, Height: 50}
	file_100x100 = &File{FileName: "100x100", Width: 100, Height: 100}
	file_200x200 = &File{FileName: "200x200", Width: 200, Height: 200}
)

func TestSortMaxSide(t *testing.T) {

	cases := []struct {
		in   []*File
		want []*File
	}{
		{
			in: []*File{
				file_50x50,
				file_100x100,
				file_200x200,
			},
			want: []*File{
				file_200x200,
				file_100x100,
				file_50x50,
			},
		},
	}

	for _, c := range cases {
		SortMaxSide(c.in)
		for i, file := range c.want {
			if c.in[i] != file {
				t.Errorf("File not sorted correctly: want %s, got %s", file.FileName, c.in[i].FileName)
			}
		}
	}
}
