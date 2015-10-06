package main

import "fmt"

type Node struct {
	X     int
	Y     int
	W     int
	H     int
	Used  bool
	Right *Node
	Down  *Node
}

func PackGrowing(files []*File, maxWidth int, maxHeight int) {
	w, h := 0, 0
	if len(files) > 0 {
		w = files[0].Width
		h = files[0].Height
	}
	root := &Node{
		X: 0,
		Y: 0,
		W: w,
		H: h,
	}
	for _, file := range files {
		if node := findNode(root, file); node != nil {
			_ = splitNode(node, file)
		} else {
			_ = growNode(root, file, maxWidth, maxHeight)
		}
	}
}

func findNode(root *Node, file *File) *Node {
	fmt.Println("Find node", root, file)
	if root.Used {
		if res := findNode(root.Right, file); res != nil {
			return res
		} else {
			return findNode(root.Down, file)
		}
	} else if file.Width <= root.W && file.Height <= root.H {
		return root
	} else {
		return nil
	}
}

func splitNode(node *Node, file *File) *Node {
	fmt.Println("Split node", node, file)
	node.Used = true
	file.X = node.X
	file.Y = node.Y
	fmt.Println("Changed the file", file)
	node.Right = &Node{
		X: node.X + file.Width,
		Y: node.Y,
		W: node.W - file.Width,
		H: node.H + file.Height,
	}
	node.Down = &Node{
		X: node.X,
		Y: node.Y + file.Height,
		W: node.W,
		H: node.H - file.Height,
	}
	return node
}

func growNode(root *Node, file *File, maxWidth int, maxHeight int) *Node {
	fmt.Println("Grow node", root, file)
	canGrowRight := file.Height <= root.H && root.W+file.Width < maxWidth
	canGrowDown := file.Width <= root.W && root.H+file.Height < maxHeight
	shouldGrowRight := canGrowRight && root.H >= root.W+file.Width
	shouldGrowDown := canGrowDown && root.W >= root.H+file.Height

	if shouldGrowRight {
		return growRight(root, file)
	} else if shouldGrowDown {
		return growDown(root, file)
	} else if canGrowRight {
		return growRight(root, file)
	} else if canGrowDown {
		return growDown(root, file)
	} else {
		return nil
	}
}

func growRight(root *Node, file *File) *Node {
	fmt.Println("Grow right", root, file)
	oldRoot := &Node{
		X:     root.X,
		Y:     root.Y,
		W:     root.W,
		H:     root.H,
		Used:  root.Used,
		Right: root.Right,
		Down:  root.Down,
	}
	root.X = 0
	root.Y = 0
	root.W = oldRoot.W + file.Width
	root.H = oldRoot.H
	root.Right = &Node{
		X: oldRoot.W,
		Y: 0,
		W: file.Width,
		H: oldRoot.H,
	}
	root.Down = oldRoot
	if node := findNode(root, file); root != nil {
		return splitNode(node, file)
	} else {
		return nil
	}
}

func growDown(root *Node, file *File) *Node {
	fmt.Println("Grow down", root, file)
	oldRoot := &Node{
		X:     root.X,
		Y:     root.Y,
		W:     root.W,
		H:     root.H,
		Used:  root.Used,
		Right: root.Right,
		Down:  root.Down,
	}
	root.X = 0
	root.Y = 0
	root.W = oldRoot.W
	root.H = oldRoot.H + file.Height
	root.Right = oldRoot
	root.Down = &Node{
		X: 0,
		Y: oldRoot.H,
		W: oldRoot.W,
		H: file.Height,
	}
	if node := findNode(root, file); root != nil {
		return splitNode(node, file)
	} else {
		return nil
	}
}
