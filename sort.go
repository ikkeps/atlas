package main

import "sort"

type fileSorter struct {
	files []*File
	by    func(f1, f2 *File) bool
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
	return s.by(s.files[i], s.files[j])
}

// Sorts the files by the maximum size of their sides
func SortMaxSide(files []*File) {
	s := fileSorter{
		files: files,
		by:    maxSide,
	}
	sort.Sort(s)
}

func maxSide(f1, f2 *File) bool {
	// TODO implement max side sort
	return false
}
