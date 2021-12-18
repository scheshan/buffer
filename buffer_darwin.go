package buffer

import (
	"golang.org/x/sys/unix"
	"sync"
	"unsafe"
)

var iovecPool = &sync.Pool{
	New: func() interface{} {
		return make([]unix.Iovec, 0)
	},
}

func (t *Buffer) CopyToFile(fd int) (int, error) {
	if err := t.checkRead(1); err != nil {
		return 0, err
	}

	iovecs := iovecPool.Get().([]unix.Iovec)

	h := t.head
	for h != nil {
		if h.len() > 0 {
			iovec := unix.Iovec{}
			iovec.SetLen(h.len())
			iovec.Base = &h.b[h.r]

			iovecs = append(iovecs, iovec)
		}
		h = h.next
	}

	wn, _, err := unix.RawSyscall(unix.SYS_WRITEV, uintptr(fd), uintptr(unsafe.Pointer(&iovecs[0])), uintptr(len(iovecs)))

	if err != 0 {
		return 0, err
	}

	n := int(wn)
	if n > 0 {
		_ = t.Skip(n)
	}

	return n, nil
}
