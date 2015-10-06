package main

type node struct {
	x, y, w, h  int
	right, down *node
}

func (n *node) clone() *node {
	return &node{
		x:     n.x,
		y:     n.y,
		w:     n.w,
		h:     n.h,
		right: n.right,
		down:  n.down,
	}
}

func PackGrowing(files []*File, maxWidth int, maxHeight int) {
	w, h := 0, 0
	if len(files) == 0 {
		return
	}
	w = files[0].Width
	h = files[0].Height
	root := &node{
		x: 0,
		y: 0,
		w: w,
		h: h,
	}
	for _, file := range files {
		if node := root.find(file); node != nil {
			_ = node.split(file)
		} else {
			_ = root.grow(file, maxWidth, maxHeight)
		}
	}
}

func (n *node) find(file *File) *node {
	if n.right != nil || n.down != nil {
		if res := n.right.find(file); res != nil {
			return res
		} else {
			return n.down.find(file)
		}
	} else if file.Width <= n.w && file.Height <= n.h {
		return n
	} else {
		return nil
	}
}

func (n *node) split(file *File) *node {
	file.X = n.x
	file.Y = n.y
	n.right = &node{
		x: n.x + file.Width,
		y: n.y,
		w: n.w - file.Width,
		h: n.h,
	}
	n.down = &node{
		x: n.x,
		y: n.y + file.Height,
		w: n.w,
		h: n.h - file.Height,
	}
	return n
}

func (n *node) grow(file *File, maxWidth int, maxHeight int) *node {
	canGrowRight := file.Height <= n.h && n.w+file.Width < maxWidth
	canGrowDown := file.Width <= n.w && n.h+file.Height < maxHeight
	shouldGrowRight := canGrowRight && n.h >= n.w+file.Width
	shouldGrowDown := canGrowDown && n.w >= n.h+file.Height

	if shouldGrowRight {
		return n.growRight(file)
	} else if shouldGrowDown {
		return n.growDown(file)
	} else if canGrowRight {
		return n.growRight(file)
	} else if canGrowDown {
		return n.growDown(file)
	} else {
		return nil
	}
}

func (n *node) growRight(file *File) *node {
	prev := n.clone()
	n.x, n.y = 0, 0
	n.w = prev.w + file.Width
	n.right = &node{
		x: prev.w,
		y: 0,
		w: file.Width,
		h: prev.h,
	}
	n.down = prev
	if next := n.find(file); next != nil {
		return next.split(file)
	} else {
		return nil
	}
}

func (n *node) growDown(file *File) *node {
	prev := n.clone()
	n.x, n.y = 0, 0
	n.h = prev.h + file.Height
	n.right = prev
	n.down = &node{
		x: 0,
		y: prev.h,
		w: prev.w,
		h: file.Height,
	}
	if next := n.find(file); next != nil {
		return next.split(file)
	} else {
		return nil
	}
}
