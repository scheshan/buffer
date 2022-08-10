package buffer

import "sync"

var nodesPool = &sync.Pool{
	New: func() interface{} {
		return &node{}
	},
}

type node struct {
	buf []byte
	r   int
	w   int
	adj int
}

func (t *node) Cap() int {
	return len(t.buf)
}

func (t *node) Len() int {
	return t.w - t.r
}

func (t *node) Available() int {
	return t.Cap() - t.w
}

func (t *node) Release() {
	defaultBytesPool.put(t.buf)
	t.buf = nil

	t.w = 0
	t.r = 0
	t.adj = 0
	nodesPool.Put(t)
}

func newNode(size int) *node {
	n := nodesPool.Get().(*node)

	n.buf = defaultBytesPool.get(size)
	return n
}
