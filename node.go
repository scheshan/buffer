package buffer

import "sync"

var (
	nodePool = &sync.Pool{
		New: func() interface{} {
			return &node{}
		},
	}
)

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

//release free the resource which node uses
func (t *node) release() {
	releaseBytes(t.b)

	t.b = nil
	t.next = nil
	t.w = 0
	t.r = 0
	nodePool.Put(t)
}

//newNode return a new node instance
func newNode(size int) *node {
	n := nodePool.Get().(*node)
	n.b = getBytes(size)

	return n
}
