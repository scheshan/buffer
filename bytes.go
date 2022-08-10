package buffer

import "sync"

var (
	defaultBytesPool = newBytesPoolSize(63)
)

type bytesPool struct {
	pools   []*sync.Pool
	maxSize int
}

func (t *bytesPool) isPowerOf2(size int) bool {
	return size&(size-1) == 0
}

func (t *bytesPool) ind(size int) int {
	ind := 0
	if !t.isPowerOf2(size) {
		ind = 1
	}

	for size > 0 {
		size = size >> 1
		ind++
	}

	return ind - 1
}

func (t *bytesPool) get(size int) []byte {
	if size <= 0 || size > t.maxSize {
		panic("invalid size")
	}

	ind := t.ind(size)
	b := t.pools[ind].Get().([]byte)

	return b
}

func (t *bytesPool) put(data []byte) {
	if data == nil {
		return
	}

	if !t.isPowerOf2(len(data)) {
		return
	}

	ind := t.ind(len(data))
	if ind >= len(t.pools) {
		return
	}

	t.pools[ind].Put(data)
}

func newBytesPoolSize(size int) *bytesPool {
	if size <= 0 || size > 64 {
		panic("invalid pool size")
	}

	p := new(bytesPool)
	p.maxSize = 1 << (size - 1)
	p.pools = make([]*sync.Pool, size, size)
	for i := 0; i < size; i++ {
		bytes := 1 << i
		p.pools[i] = &sync.Pool{
			New: func() interface{} {
				buf := make([]byte, bytes)

				return buf
			},
		}
	}

	return p
}

func getBytes(size int) []byte {
	return defaultBytesPool.get(size)
}

func releaseBytes(data []byte) {
	defaultBytesPool.put(data)
}
