package sitter

import (
	"fmt"
)

type IterMode int

const (
	DFSMode IterMode = iota
	BFSMode
)

type (
	TreeIter  func(yield func(*Node) bool)
	QueryIter func(yield func(map[string]*Node) bool)
)

// TreeIterator takes a node and mode (DFS/BFS) and returns iterator over children of the node.
// named determines whether to iterate over only named children
func TreeIterator(n *Node, mode IterMode, named bool) TreeIter {
	return func(yield func(*Node) bool) {
		nodesToVisit := []*Node{n}
		var current *Node
		for len(nodesToVisit) > 0 {
			switch mode {
			case DFSMode:
				current, nodesToVisit = nodesToVisit[len(nodesToVisit)-1], nodesToVisit[:len(nodesToVisit)-1]
			case BFSMode:
				current, nodesToVisit = nodesToVisit[0], nodesToVisit[1:]
			default:
				panic(fmt.Errorf("unsupported iteration mode: %v", mode))
			}
			if !yield(current) {
				return
			}
			if named {
				for i := 0; i < int(current.NamedChildCount()); i++ {
					nodesToVisit = append(nodesToVisit, current.NamedChild(i))
				}
			} else {
				for i := 0; i < int(current.ChildCount()); i++ {
					nodesToVisit = append(nodesToVisit, current.Child(i))
				}
			}
		}
	}
}

func QueryIterator(root *Node, query *Query) QueryIter {
	return func(yield func(map[string]*Node) bool) {
		cursor := NewQueryCursor()
		cursor.Exec(query, root)
		captures := make(map[uint32]string)
		for i := uint32(0); i < query.CaptureCount(); i++ {
			captures[i] = query.CaptureNameForId(i)
		}

		for {
			match, found := cursor.NextMatch()
			if !found {
				return
			}
			results := make(map[string]*Node, len(match.Captures))

			for _, capture := range match.Captures {
				results[captures[capture.Index]] = capture.Node
			}

			if !yield(results) {
				return
			}
		}
	}
}
