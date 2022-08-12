package buffer

import (
	"errors"
)

var (
	ErrExceedMaximumSize = errors.New("exceed maximum size")
	ErrNoEnoughData      = errors.New("no enough data to read")
	minNodeSize          = 2048
)

type Buffer struct {
	nodes   []*node
	nc      int //node count
	size    int
	maxSize int
}

//#region read logic

func (t *Buffer) ReadBool() (bool, error) {
	res, err := t.ReadUInt8()
	return res != 0, err
}

func (t *Buffer) ReadUInt8() (uint8, error) {
	if err := t.ensureReadable(1); err != nil {
		return 0, err
	}

	defer t.skip(1)
	return t.getUInt8(0), nil
}

func (t *Buffer) ReadInt8() (int8, error) {
	res, err := t.ReadUInt8()
	return int8(res), err
}

func (t *Buffer) ReadByte() (byte, error) {
	return t.ReadUInt8()
}

func (t *Buffer) ReadUInt16() (uint16, error) {
	if err := t.ensureReadable(2); err != nil {
		return 0, err
	}

	defer t.skip(1)
	return t.getUInt16(0), nil
}

func (t *Buffer) ReadInt16() (int16, error) {
	res, err := t.ReadUInt16()
	return int16(res), err
}

func (t *Buffer) ReadUInt32() (uint32, error) {
	if err := t.ensureReadable(4); err != nil {
		return 0, err
	}

	defer t.skip(4)
	return t.getUInt32(0), nil
}

func (t *Buffer) ReadInt32() (int32, error) {
	res, err := t.ReadUInt32()
	return int32(res), err
}

func (t *Buffer) ReadUInt64() (uint64, error) {
	if err := t.ensureReadable(8); err != nil {
		return 0, err
	}

	defer t.skip(8)
	return t.getUInt64(0), nil
}

func (t *Buffer) ReadInt64() (int64, error) {
	res, err := t.ReadUInt64()
	return int64(res), err
}

func (t *Buffer) Len() int {
	return t.size
}

//#endregion

//#region write logic

func (t *Buffer) WriteBytes(b []byte) error {
	if err := t.ensureWriteable(len(b)); err != nil {
		return err
	}

	w := t.writer()
	if w != nil && w.WritableBytes() > 0 {
		n := copy(w.buf[w.w:], b)
		w.w += n
		t.size += n

		if n == len(b) {
			return nil
		}
		b = b[n:]
	}

	s := len(b)
	if s < 2048 {
		s = 2048
	}
	node := newNode(s)
	node.w += copy(node.buf, b)
	t.addNodeToArray(node)
	t.size += node.w
	return nil
}

func (t *Buffer) WriteBool(b bool) error {
	var num byte = 0
	if b {
		num = 1
	}
	return t.WriteByte(num)
}

func (t *Buffer) WriteByte(n byte) error {
	return t.WriteUInt8(n)
}

func (t *Buffer) WriteUInt8(n uint8) error {
	if err := t.ensureWriteable(1); err != nil {
		return err
	}

	t.writeUInt8(n)
	return nil
}

func (t *Buffer) WriteInt8(n int8) error {
	return t.WriteUInt8(uint8(n))
}

func (t *Buffer) WriteUInt16(n uint16) error {
	if err := t.ensureWriteable(2); err != nil {
		return err
	}

	t.writeUInt16(n)
	return nil
}

func (t *Buffer) WriteInt16(n int16) error {
	return t.WriteUInt16(uint16(n))
}

func (t *Buffer) WriteUInt32(n uint32) error {
	if err := t.ensureWriteable(4); err != nil {
		return err
	}

	t.writeUInt32(n)
	return nil
}

func (t *Buffer) WriteInt32(n int32) error {
	return t.WriteUInt32(uint32(n))
}

func (t *Buffer) WriteUInt64(n uint64) error {
	if err := t.ensureWriteable(8); err != nil {
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

//#endregion

func (t *Buffer) addNodeToArray(n *node) {
	t.expand()

	t.nodes[t.nc] = n
	t.nc++
}

func (t *Buffer) expand() {
	if t.nodes == nil {
		t.nodes = make([]*node, 1)
		return
	}

	if t.nc == len(t.nodes) {
		s := len(t.nodes) << 1
		nodes := make([]*node, s)

		copy(nodes, t.nodes)
		t.nodes = nodes
		return
	}
}

func (t *Buffer) shrink() {
	if t.nodes == nil || t.nodes[0].ReadableBytes() > 0 {
		return
	}

	l := 0
	r := 0
	n := t.nc
	for r < n {
		if t.nodes[r].ReadableBytes() <= 0 {
			t.nc--
		} else {
			t.nodes[l] = t.nodes[r]
			l++
		}
		r++
	}

	for l < len(t.nodes) {
		t.nodes[l] = nil
		l++
	}
}

func (t *Buffer) adjust() {
	if t.nc == 0 {
		return
	}

	t.nodes[0].adj = 0
	for i := 1; i < t.nc; i++ {
		t.nodes[i].adj = t.nodes[i-1].adj + t.nodes[i-1].ReadableBytes()
	}
}

func (t *Buffer) ensureWriteable(size int) error {
	if t.maxSize > 0 && t.maxSize-t.size < size {
		return ErrExceedMaximumSize
	}

	return nil
}

func (t *Buffer) writer() *node {
	if t.nc == 0 {
		return nil
	}

	return t.nodes[t.nc-1]
}

func (t *Buffer) reader() *node {
	if t.nodes == nil {
		return nil
	}

	return t.nodes[0]
}

func (t *Buffer) ensureReadable(size int) error {
	if t.size < size {
		return ErrNoEnoughData
	}
	return nil
}

func (t *Buffer) writeUInt8(n uint8) {
	if t.writer() == nil || t.writer().WritableBytes() < 1 {
		t.addNodeToArray(newNode(minNodeSize))
	}

	t.writer().buf[t.writer().w] = n
	t.writer().w++
	t.size++
}

func (t *Buffer) writeUInt16(n uint16) {
	if t.writer() != nil && t.writer().WritableBytes() >= 2 {
		t.writer().buf[t.writer().w] = uint8(n >> 8)
		t.writer().buf[t.writer().w+1] = uint8(n)

		t.writer().w += 2
		t.size += 2
	} else {
		t.writeUInt8(uint8(n >> 8))
		t.writeUInt8(uint8(n))
	}
}

func (t *Buffer) writeUInt32(n uint32) {
	if t.writer() != nil && t.writer().WritableBytes() >= 4 {
		t.writer().buf[t.writer().w] = uint8(n >> 24)
		t.writer().buf[t.writer().w+1] = uint8(n >> 16)
		t.writer().buf[t.writer().w+2] = uint8(n >> 8)
		t.writer().buf[t.writer().w+3] = uint8(n)

		t.writer().w += 4
		t.size += 4
	} else {
		t.writeUInt16(uint16(n >> 16))
		t.writeUInt16(uint16(n))
	}
}

func (t *Buffer) writeUInt64(n uint64) {
	if t.writer() != nil && t.writer().WritableBytes() >= 8 {
		t.writer().buf[t.writer().w] = uint8(n >> 56)
		t.writer().buf[t.writer().w+1] = uint8(n >> 48)
		t.writer().buf[t.writer().w+2] = uint8(n >> 40)
		t.writer().buf[t.writer().w+3] = uint8(n >> 32)
		t.writer().buf[t.writer().w+4] = uint8(n >> 24)
		t.writer().buf[t.writer().w+5] = uint8(n >> 16)
		t.writer().buf[t.writer().w+6] = uint8(n >> 8)
		t.writer().buf[t.writer().w+7] = uint8(n)

		t.writer().w += 8
		t.size += 8
	} else {
		t.writeUInt32(uint32(n >> 32))
		t.writeUInt32(uint32(n))
	}
}

func (t *Buffer) skip(n int) {
	t.size -= n

	i := 0
	var no *node
	for n > 0 {
		no = t.nodes[i]
		avail := no.ReadableBytes()
		if avail > n {
			no.r += n
			n = 0
		} else {
			no.r = no.w
			n -= avail
		}
	}

	t.shrink()
	t.adjust()
}

func (t *Buffer) getUInt8(idx int) uint8 {
	n, i := t.getNode(idx)
	return n.buf[i]
}

func (t *Buffer) getUInt16(idx int) uint16 {
	n, i := t.getNode(idx)
	if i <= n.Cap()-2 {
		return (uint16(n.buf[i]) << 8) | uint16(n.buf[i+1])
	} else {
		return (uint16(t.getUInt8(idx)) << 8) | uint16(t.getUInt8(idx+1))
	}
}

func (t *Buffer) getUInt32(idx int) uint32 {
	n, i := t.getNode(idx)
	if i <= n.Cap()-4 {
		return (uint32(n.buf[i]) << 24) |
			(uint32(n.buf[i+1]) << 16) |
			(uint32(n.buf[i+2]) << 8) |
			uint32(n.buf[i+3])
	} else {
		return (uint32(t.getUInt16(idx)) << 16) | uint32(t.getUInt16(idx+2))
	}
}

func (t *Buffer) getUInt64(idx int) uint64 {
	n, i := t.getNode(idx)
	if i <= n.Cap()-8 {
		return (uint64(n.buf[i]) << 56) |
			(uint64(n.buf[i+1]) << 48) |
			(uint64(n.buf[i+2]) << 40) |
			(uint64(n.buf[i+3]) << 32) |
			(uint64(n.buf[i+4]) << 24) |
			(uint64(n.buf[i+5]) << 16) |
			(uint64(n.buf[i+6]) << 8) |
			uint64(n.buf[i+7])
	} else {
		return (uint64(t.getUInt32(idx)) << 32) | uint64(t.getUInt16(idx+4))
	}
}

func (t *Buffer) getNode(idx int) (*node, int) {
	l, r := 0, t.nc
	var m int
	for l < r {
		m = l + (r >> 1)

		n := t.nodes[m]
		if n.adj > idx {
			r = r >> 1
		} else if n.adj+n.ReadableBytes() < idx {
			l = l << 1
		} else {
			return n, idx - n.adj + n.r
		}
	}

	return nil, -1
}

func New(maxSize int) *Buffer {
	buf := &Buffer{
		maxSize: maxSize,
	}
	return buf
}
