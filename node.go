package buffer

type node struct {
	b    []byte
	next *node
	w    int
	r    int
}

//cap return bytes can write
func (t *node) cap() int {
	return len(t.b) - t.w
}

//len return bytes can read
func (t *node) len() int {
	return t.w - t.r
}

//newNode return a new node instance
func newNode(size int) *node {
	n := new(node)
	n.b = make([]byte, size)

	return n
}
