package main

import "sort"

// A function that will compare two files and return the difference
// A negative number will be returned if f2 is should be sorted before
// file 1, a zero represents equality and a positive number is returned
// if file 1 should appear before file 2
type sortFunc func(f1, f2 *File) int

// Sorter interface
type fileSorter struct {
	files []*File
	by    sortFunc
}

// Len is part of sort.Interface.
func (s fileSorter) Len() int {
	return len(s.files)
}

// Swap is part of sort.Interface.
func (s fileSorter) Swap(i, j int) {
	s.files[i], s.files[j] = s.files[j], s.files[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s fileSorter) Less(i, j int) bool {
	return s.by(s.files[i], s.files[j]) > 0
}

// Compares the width of the two files and returns the difference
func width(f1, f2 *File) int {
	return f1.Width - f2.Width
}

// Compares the height of the two files and returns the difference
func height(f1, f2 *File) int {
	return f1.Height - f2.Height
}

// Compares the shorter dimensions of the two files and returns the difference
func minSide(f1, f2 *File) int {
	f1min, f2min := f1.Width, f2.Width
	if f1.Height < f1.Width {
		f1min = f1.Height
	}
	if f2.Height < f2.Width {
		f2min = f2.Height
	}
	return f1min - f2min
}

// Compares the longer dimensions of the two files and returns the difference
func maxSide(f1, f2 *File) int {
	f1max, f2max := f1.Width, f2.Width
	if f1.Height > f1.Width {
		f1max = f1.Height
	}
	if f2.Height > f2.Width {
		f2max = f2.Height
	}
	return f1max - f2max
}

// Compares two files using a number of sorting methods. If the first sorting
// method returns 0 (the two files are equal) then the next sorting method
// will be used. Continues until a sorting method returns a non-zero value.
func multiSort(f1, f2 *File, methods ...sortFunc) int {
	for _, method := range methods {
		if res := method(f1, f2); res != 0 {
			return res
		}
	}
	return 0
}

// Compares the two files using a combination of other sorts
func maxSideMinSideHeightWidth(f1, f2 *File) int {
	return multiSort(f1, f2, maxSide, minSide, height, width)
}

// Sorts the files by the maximum size of their sides
func SortMaxSide(files []*File) (sorted []*File) {
	s := fileSorter{
		files: files,
		by:    maxSideMinSideHeightWidth,
	}
	sort.Sort(s)
	return s.files
}
