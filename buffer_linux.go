package buffer

import (
	"golang.org/x/sys/unix"
	"sync"
)

var iovsPool = &sync.Pool{
	New: func() interface{} {
		return make([][]byte, 0)
	},
}

func (t *Buffer) CopyToFile(fd int) (n int, err error) {
	iovs := iovsPool.Get().([][]byte)

	h := t.head
	for h != nil {
		iovs = append(iovs, h.b[h.r:h.w])

		h = h.next
	}

	n, err = unix.Writev(fd, iovs)
	if n > 0 {
		_ = t.Skip(n)
	}

	iovsPool.Put(iovs)

	return n, err
}
