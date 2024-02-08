package sitter

type IterMode int

const (
	DFSMode IterMode = iota
	BFSMode
)

// treeIterator for a tree of nodes
type treeIterator struct {
	named bool
	mode  IterMode

	nodesToVisit []*Node
}

// TreeIterator takes a node and mode (DFS/BFS) and returns iterator over children of the node.
// named determines whether to iterate over only named children
func TreeIterator(n *Node, mode IterMode, named bool) func(func(*Node) bool) {
	iter := &treeIterator{
		named:        named,
		mode:         mode,
		nodesToVisit: []*Node{n},
	}
	return iter.Next
}

func (iter *treeIterator) Next(yield func(*Node) bool) {
	for len(iter.nodesToVisit) > 0 {
		current := iter.nodesToVisit[0]
		iter.nodesToVisit = iter.nodesToVisit[1:]
		if !yield(current) {
			return
		}
		var children []*Node
		if iter.named {
			for i := 0; i < int(current.NamedChildCount()); i++ {
				children = append(children, current.NamedChild(i))
			}
		} else {
			for i := 0; i < int(current.ChildCount()); i++ {
				children = append(children, current.Child(i))
			}
		}

		switch iter.mode {
		case DFSMode:
			iter.nodesToVisit = append(children, iter.nodesToVisit...)
		case BFSMode:
			iter.nodesToVisit = append(iter.nodesToVisit, children...)
		default:
			panic("not implemented")
		}
	}
}

type queryIterator struct {
	root     *Node
	captures map[uint32]string
	cursor   *QueryCursor
}

func QueryIterator(root *Node, query *Query) func(func(map[string]*Node) bool) {
	cursor := NewQueryCursor()
	cursor.Exec(query, root)
	captures := make(map[uint32]string)
	for i := uint32(0); i < query.CaptureCount(); i++ {
		captures[i] = query.CaptureNameForId(i)
	}
	iter := &queryIterator{
		root:     root,
		captures: captures,
		cursor:   cursor,
	}
	return iter.Next
}

func (iter *queryIterator) Next(yield func(map[string]*Node) bool) {
	for {
		match, found := iter.cursor.NextMatch()
		if !found {
			return
		}
		results := make(map[string]*Node)

		for _, capture := range match.Captures {
			results[iter.captures[capture.Index]] = capture.Node
		}

		if !yield(results) {
			return
		}
	}
}
