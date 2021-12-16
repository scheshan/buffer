package buffer

import (
	"golang.org/x/sys/unix"
	"os"
	"testing"
)

func TestBuffer_Cap(t *testing.T) {
	buf := NewBufferCap(0, 1)
	if buf.Cap() != 1 {
		t.Fail()
	}

	if err := buf.WriteByte(1); err != nil {
		t.Fail()
	}

	if err := buf.WriteUInt8(1); err != ErrBufferOverflow {
		t.Fail()
	}

	if err := buf.WriteUInt16(1); err != ErrBufferOverflow {
		t.Fail()
	}

	if err := buf.WriteUInt32(1); err != ErrBufferOverflow {
		t.Fail()
	}

	if err := buf.WriteUInt64(1); err != ErrBufferOverflow {
		t.Fail()
	}

	buf = NewBuffer()
	if buf.Cap() != -1 {
		t.Fail()
	}
}

func TestBuffer_Len(t *testing.T) {
	buf := NewBuffer()

	buf.WriteInt64(1)
	if buf.Len() != 8 {
		t.Fail()
	}

	buf.WriteInt32(1)
	if buf.Len() != 12 {
		t.Fail()
	}
}

func TestBuffer_WriteUInt8(t *testing.T) {
	buf := NewBuffer()

	if err := buf.WriteUInt8(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadUInt8(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadUInt8(); err != ErrBufferNotEnough {
		t.Fail()
	}
}

func TestBuffer_WriteInt8(t *testing.T) {
	buf := NewBuffer()

	if err := buf.WriteInt8(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadInt8(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadInt8(); err != ErrBufferNotEnough {
		t.Fail()
	}
}

func TestBuffer_WriteBool(t *testing.T) {
	buf := NewBuffer()

	if err := buf.WriteBool(true); err != nil {
		t.Fail()
	}

	if err := buf.WriteBool(false); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadBool(); err != nil || !n {
		t.Fail()
	}

	if n, err := buf.ReadBool(); err != nil || n {
		t.Fail()
	}

	if _, err := buf.ReadBool(); err != ErrBufferNotEnough {
		t.Fail()
	}
}

func TestBuffer_ReadByte(t *testing.T) {
	buf := NewBuffer()

	if err := buf.WriteByte(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadByte(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadByte(); err != ErrBufferNotEnough {
		t.Fail()
	}
}

func TestBuffer_WriteUInt16(t *testing.T) {
	buf := NewBuffer()

	if err := buf.WriteUInt16(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadUInt16(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadUInt16(); err != ErrBufferNotEnough {
		t.Fail()
	}
}

func TestBuffer_WriteInt16(t *testing.T) {
	buf := NewBuffer()

	if err := buf.WriteInt16(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadInt16(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadInt16(); err != ErrBufferNotEnough {
		t.Fail()
	}
}

func TestBuffer_WriteUInt32(t *testing.T) {
	buf := NewBuffer()

	if err := buf.WriteUInt32(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadUInt32(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadUInt32(); err != ErrBufferNotEnough {
		t.Fail()
	}
}

func TestBuffer_WriteInt32(t *testing.T) {
	buf := NewBuffer()

	if err := buf.WriteInt32(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadInt32(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadInt32(); err != ErrBufferNotEnough {
		t.Fail()
	}
}

func TestBuffer_WriteUInt64(t *testing.T) {
	buf := NewBuffer()

	if err := buf.WriteUInt64(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadUInt64(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadUInt64(); err != ErrBufferNotEnough {
		t.Fail()
	}
}

func TestBuffer_WriteInt64(t *testing.T) {
	buf := NewBuffer()

	if err := buf.WriteInt64(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadInt64(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadInt64(); err != ErrBufferNotEnough {
		t.Fail()
	}
}

func TestBuffer_WriteUInt(t *testing.T) {
	buf := NewBuffer()

	if err := buf.WriteUInt(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadUInt(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadUInt(); err != ErrBufferNotEnough {
		t.Fail()
	}
}

func TestBuffer_WriteInt(t *testing.T) {
	buf := NewBuffer()

	if err := buf.WriteInt(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadInt(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadInt(); err != ErrBufferNotEnough {
		t.Fail()
	}
}

func TestBuffer_HalfReadWrite(t *testing.T) {
	buf := NewBufferSize(1)

	if err := buf.WriteInt64(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadInt64(); err != nil || n != 1 {
		t.Fail()
	}

	if err := buf.WriteUInt64(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadUInt64(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadInt64(); err != ErrBufferNotEnough {
		t.Fail()
	}
}

func TestBuffer_Release(t *testing.T) {
	buf := NewBuffer()

	buf.WriteInt64(1)

	if buf.Ref() != 1 {
		t.Fail()
	}

	buf.IncrRef()
	if buf.Ref() != 2 {
		t.Fail()
	}

	buf.Release()
	buf.Release()
	if buf.Ref() != 0 {
		t.Fail()
	}

	buf.Release()
	if buf.Ref() != 0 {
		t.Fail()
	}
}

func TestBuffer_CopyToFile(t *testing.T) {
	path := "/Users/heshan/tmp/test"
	os.Remove(path)

	f, err := os.Create(path)
	if err != nil {
		t.Fail()
	}
	fd := int(f.Fd())

	buf := NewBuffer()
	buf.WriteInt64(1)

	if n, err, complete := buf.CopyToFile(fd); err != nil || !complete || n != 8 {
		t.Fail()
	}
	if _, err, _ := buf.CopyToFile(fd); err != ErrBufferNotEnough {
		t.Fail()
	}
	unix.Close(fd)

	f, err = os.Open(path)
	if err != nil {
		t.Fail()
	}
	fd = int(f.Fd())
	buf = NewBuffer()
	if n, err, complete := buf.CopyFromFile(fd); err != nil || !complete || n != 8 {
		t.Fail()
	}

	if n, err := buf.ReadInt64(); err != nil || n != 1 {
		t.Fail()
	}

	os.Remove(path)
}
