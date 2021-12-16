package buffer

import (
	"errors"
	"golang.org/x/sys/unix"
	"sync"
)

var (
	ErrBufferOverflow  = errors.New("buffer size overflow")
	ErrBufferNotEnough = errors.New("buffer has no enough data")
	defaultBufferSize  = 8192
	bufferPool         = &sync.Pool{
		New: func() interface{} {
			return new(Buffer)
		},
	}
)

type Buffer struct {
	head    *node
	tail    *node
	r       *node
	size    int
	cap     int
	ref     int
	minSize int
}

//#region read methods

func (t *Buffer) ReadByte() (byte, error) {
	return t.ReadUInt8()
}

func (t *Buffer) ReadBool() (bool, error) {
	n, err := t.ReadUInt8()
	return n != 0, err
}

func (t *Buffer) ReadUInt8() (uint8, error) {
	if err := t.checkLen(1); err != nil {
		return 0, err
	}

	return t.readUInt8(), nil
}

func (t *Buffer) ReadInt8() (int8, error) {
	n, err := t.ReadUInt8()
	return int8(n), err
}

func (t *Buffer) ReadUInt16() (uint16, error) {
	if err := t.checkLen(2); err != nil {
		return 0, err
	}

	return t.readUInt16(), nil
}

func (t *Buffer) ReadInt16() (int16, error) {
	n, err := t.ReadUInt16()
	return int16(n), err
}

func (t *Buffer) ReadUInt32() (uint32, error) {
	if err := t.checkLen(4); err != nil {
		return 0, err
	}

	return t.readUInt32(), nil
}

func (t *Buffer) ReadInt32() (int32, error) {
	n, err := t.ReadUInt32()
	return int32(n), err
}

func (t *Buffer) ReadUInt64() (uint64, error) {
	if err := t.checkLen(8); err != nil {
		return 0, err
	}

	return t.readUInt64(), nil
}

func (t *Buffer) ReadInt64() (int64, error) {
	n, err := t.ReadUInt64()
	return int64(n), err
}

func (t *Buffer) ReadInt() (int, error) {
	n, err := t.ReadInt64()
	return int(n), err
}

func (t *Buffer) ReadUInt() (int, error) {
	n, err := t.ReadUInt64()
	return int(n), err
}

func (t *Buffer) CopyToFile(fd int) (n int, err error, complete bool) {
	if err := t.checkLen(1); err != nil {
		return 0, err, false
	}

	n, err = unix.Write(fd, t.tail.b[t.tail.r:t.tail.w])
	if n > 0 {
		t.tail.r += n
		t.size -= n

		t.releaseHead()
	}

	complete = t.Len() == 0

	return
}

//#endregion

//#region write methods

func (t *Buffer) WriteByte(n byte) error {
	return t.WriteUInt8(n)
}

func (t *Buffer) WriteBool(n bool) error {
	if n {
		return t.WriteByte(1)
	} else {
		return t.WriteByte(0)
	}
}

func (t *Buffer) WriteUInt8(n uint8) error {
	if err := t.checkCap(1); err != nil {
		return err
	}

	t.writeUInt8(n)
	return nil
}

func (t *Buffer) WriteInt8(n int8) error {
	return t.WriteByte(byte(n))
}

func (t *Buffer) WriteUInt16(n uint16) error {
	if err := t.checkCap(2); err != nil {
		return err
	}

	t.writeUInt16(n)
	return nil
}

func (t *Buffer) WriteInt16(n int16) error {
	return t.WriteUInt16(uint16(n))
}

func (t *Buffer) WriteUInt32(n uint32) error {
	if err := t.checkCap(4); err != nil {
		return err
	}

	t.writeUInt32(n)
	return nil
}

func (t *Buffer) WriteInt32(n int32) error {
	return t.WriteUInt32(uint32(n))
}

func (t *Buffer) WriteUInt64(n uint64) error {
	if err := t.checkCap(8); err != nil {
		return err
	}

	t.writeUInt64(n)
	return nil
}

func (t *Buffer) WriteInt64(n int64) error {
	return t.WriteUInt64(uint64(n))
}

func (t *Buffer) WriteInt(n int) error {
	return t.WriteInt64(int64(n))
}

func (t *Buffer) WriteUInt(n uint) error {
	return t.WriteUInt64(uint64(n))
}

func (t *Buffer) CopyFromFile(fd int) (n int, err error, complete bool) {
	if t.tail == nil || t.tail.cap() == 0 {
		t.addNode(t.minSize)
	}

	n, err = unix.Read(fd, t.tail.b[t.tail.w:])
	if n > 0 {
		t.tail.w += n
		t.size += n
	}

	complete = t.tail.cap() > 0

	return
}

//#endregion

//#region public methods

func (t *Buffer) Len() int {
	return t.size
}

func (t *Buffer) Cap() int {
	if t.cap >= 0 {
		return t.cap
	}

	return -1
}

func (t *Buffer) Ref() int {
	return t.ref
}

func (t *Buffer) IncrRef() {
	t.ref++
}

func (t *Buffer) Release() {
	if t.ref <= 0 {
		return
	}

	t.ref--
	if t.ref > 0 {
		return
	}

	for t.head != nil {
		h := t.head
		t.head = t.head.next

		h.release()
	}
	t.tail = nil
	t.r = nil
	t.size = 0
	t.cap = 0
	t.minSize = 0
	bufferPool.Put(t)
}

//#endregion

//#region private methods

func (t *Buffer) checkLen(n int) error {
	if t.Len() < n {
		return ErrBufferNotEnough
	}

	return nil
}

func (t *Buffer) checkCap(n int) error {
	if t.cap >= 0 && t.Len()+n > t.cap {
		return ErrBufferOverflow
	}

	return nil
}

func (t *Buffer) addTail(node *node) {
	if t.tail == nil {
		t.head = node
		t.tail = node
		t.r = node
	} else {
		t.tail.next = node
		t.tail = node
	}
}

func (t *Buffer) addNode(size int) {
	node := newNode(size)
	t.addTail(node)
}

func (t *Buffer) writeUInt8(n uint8) {
	if t.tail == nil || t.tail.cap() == 0 {
		t.addNode(t.minSize)
	}

	t.tail.b[t.tail.w] = n
	t.tail.w++
	t.size++
}

func (t *Buffer) writeUInt16(n uint16) {
	if t.tail == nil || t.tail.cap() == 0 {
		t.addNode(t.minSize)
	}

	if t.tail.cap() < 2 {
		t.writeUInt8(uint8(n >> 8))
		t.writeUInt8(uint8(n))
	} else {
		t.tail.b[t.tail.w] = uint8(n >> 8)
		t.tail.b[t.tail.w+1] = uint8(n)
		t.tail.w += 2
		t.size += 2
	}
}

func (t *Buffer) writeUInt32(n uint32) {
	if t.tail == nil || t.tail.cap() == 0 {
		t.addNode(t.minSize)
	}

	if t.tail.cap() < 4 {
		t.writeUInt16(uint16(n >> 16))
		t.writeUInt16(uint16(n))
	} else {
		t.tail.b[t.tail.w] = uint8(n >> 24)
		t.tail.b[t.tail.w+1] = uint8(n >> 16)
		t.tail.b[t.tail.w+2] = uint8(n >> 8)
		t.tail.b[t.tail.w+3] = uint8(n)
		t.tail.w += 4
		t.size += 4
	}
}

func (t *Buffer) writeUInt64(n uint64) {
	if t.tail == nil || t.tail.cap() == 0 {
		t.addNode(t.minSize)
	}

	if t.tail.cap() < 4 {
		t.writeUInt32(uint32(n >> 32))
		t.writeUInt32(uint32(n))
	} else {
		t.tail.b[t.tail.w] = uint8(n >> 56)
		t.tail.b[t.tail.w+1] = uint8(n >> 48)
		t.tail.b[t.tail.w+2] = uint8(n >> 40)
		t.tail.b[t.tail.w+3] = uint8(n >> 32)
		t.tail.b[t.tail.w+4] = uint8(n >> 24)
		t.tail.b[t.tail.w+5] = uint8(n >> 16)
		t.tail.b[t.tail.w+6] = uint8(n >> 8)
		t.tail.b[t.tail.w+7] = uint8(n)
		t.tail.w += 8
		t.size += 8
	}
}

func (t *Buffer) readUInt8() (n uint8) {
	n = t.r.b[t.r.r]
	t.r.r++
	t.size--

	t.releaseHead()

	return
}

func (t *Buffer) readUInt16() (n uint16) {
	if t.r.len() >= 2 {
		n = uint16(t.r.b[t.r.r]) << 8
		n |= uint16(t.r.b[t.r.r+1])
		t.r.r += 2
		t.size -= 2

		t.releaseHead()
	} else {
		n = uint16(t.readUInt8()) << 8
		n |= uint16(t.readUInt8())
	}

	return
}

func (t *Buffer) readUInt32() (n uint32) {
	if t.r.len() >= 4 {
		n = uint32(t.r.b[t.r.r]) << 24
		n |= uint32(t.r.b[t.r.r+1]) << 16
		n |= uint32(t.r.b[t.r.r+2]) << 8
		n |= uint32(t.r.b[t.r.r+3])
		t.r.r += 4
		t.size -= 4

		t.releaseHead()
	} else {
		n = uint32(t.readUInt16()) << 16
		n |= uint32(t.readUInt16())
	}

	return
}

func (t *Buffer) readUInt64() (n uint64) {
	if t.r.len() >= 8 {
		n = uint64(t.r.b[t.r.r]) << 56
		n |= uint64(t.r.b[t.r.r+1]) << 48
		n |= uint64(t.r.b[t.r.r+2]) << 40
		n |= uint64(t.r.b[t.r.r+3]) << 32
		n |= uint64(t.r.b[t.r.r+4]) << 24
		n |= uint64(t.r.b[t.r.r+5]) << 16
		n |= uint64(t.r.b[t.r.r+6]) << 8
		n |= uint64(t.r.b[t.r.r+7])
		t.r.r += 8
		t.size -= 8

		t.releaseHead()
	} else {
		n = uint64(t.readUInt32()) << 32
		n |= uint64(t.readUInt32())
	}

	return
}

//releaseHead tries to release empty nodes.
//After every single read, releaseHead should be called.
func (t *Buffer) releaseHead() {
	for t.head != nil && t.head.len() == 0 && t.head.cap() == 0 {
		h := t.head
		t.head = t.head.next
		if t.r == h {
			t.r = t.head
		}

		h.release()
	}

	if t.head == nil {
		t.r = nil
		t.tail = nil
	}
}

//#endregion

func NewBuffer() *Buffer {
	return NewBufferSize(0)
}

func NewBufferSize(minSize int) *Buffer {
	return NewBufferCap(minSize, -1)
}

func NewBufferCap(minSize int, cap int) *Buffer {
	if minSize <= 0 {
		minSize = defaultBufferSize
	}

	b := bufferPool.Get().(*Buffer)
	b.cap = cap
	b.minSize = minSize
	b.ref = 1

	return b
}
