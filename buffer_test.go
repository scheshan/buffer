package buffer

import (
	"math/rand"
	"testing"
)

func TestBuffer_Len(t *testing.T) {
	buf := New(1000000)

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
	buf := New(1000000)

	if err := buf.WriteUInt8(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadUInt8(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadUInt8(); err != ErrNoEnoughData {
		t.Fail()
	}
}

func TestBuffer_WriteInt8(t *testing.T) {
	buf := New(1000000)

	if err := buf.WriteInt8(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadInt8(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadInt8(); err != ErrNoEnoughData {
		t.Fail()
	}
}

func TestBuffer_WriteBool(t *testing.T) {
	buf := New(1000000)

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

	if _, err := buf.ReadBool(); err != ErrNoEnoughData {
		t.Fail()
	}
}

func TestBuffer_ReadByte(t *testing.T) {
	buf := New(1000000)

	if err := buf.WriteByte(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadByte(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadByte(); err != ErrNoEnoughData {
		t.Fail()
	}
}

func TestBuffer_WriteUInt16(t *testing.T) {
	buf := New(1000000)

	num := uint16(rand.Int())

	if err := buf.WriteUInt16(num); err != nil {
		t.Fail()
	}

	if n := buf.getUInt16(0); n != num {
		t.Fail()
	}

	n1, _ := buf.ReadUInt8()
	n2, _ := buf.ReadUInt8()

	num2 := (uint16(n1) << 8) | uint16(n2)
	if num != num2 {
		t.Fail()
	}
}

func TestBuffer_WriteInt16(t *testing.T) {
	buf := New(1000000)

	if err := buf.WriteInt16(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadInt16(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadInt16(); err != ErrNoEnoughData {
		t.Fail()
	}
}

func TestBuffer_WriteUInt32(t *testing.T) {
	buf := New(1000000)

	if err := buf.WriteUInt32(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadUInt32(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadUInt32(); err != ErrNoEnoughData {
		t.Fail()
	}
}

func TestBuffer_WriteInt32(t *testing.T) {
	buf := New(1000000)

	if err := buf.WriteInt32(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadInt32(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadInt32(); err != ErrNoEnoughData {
		t.Fail()
	}
}

func TestBuffer_WriteUInt64(t *testing.T) {
	buf := New(1000000)

	if err := buf.WriteUInt64(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadUInt64(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadUInt64(); err != ErrNoEnoughData {
		t.Fail()
	}
}

func TestBuffer_WriteInt64(t *testing.T) {
	buf := New(1000000)

	if err := buf.WriteInt64(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadInt64(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadInt64(); err != ErrNoEnoughData {
		t.Fail()
	}
}

func TestBuffer_WriteUInt(t *testing.T) {
	buf := New(1000000)

	if err := buf.WriteUInt(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadUInt(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadUInt(); err != ErrNoEnoughData {
		t.Fail()
	}
}

func TestBuffer_WriteInt(t *testing.T) {
	buf := New(1000000)

	if err := buf.WriteInt(1); err != nil {
		t.Fail()
	}

	if n, err := buf.ReadInt(); err != nil || n != 1 {
		t.Fail()
	}

	if _, err := buf.ReadInt(); err != ErrNoEnoughData {
		t.Fail()
	}
}

//func TestBuffer_HalfReadWrite(t *testing.T) {
//	buf := NewBufferSize(1)
//
//	if err := buf.WriteInt64(1); err != nil {
//		t.Fail()
//	}
//
//	if n, err := buf.ReadInt64(); err != nil || n != 1 {
//		t.Fail()
//	}
//
//	if err := buf.WriteUInt64(1); err != nil {
//		t.Fail()
//	}
//
//	if n, err := buf.ReadUInt64(); err != nil || n != 1 {
//		t.Fail()
//	}
//
//	if _, err := buf.ReadInt64(); err != ErrBufferNotEnough {
//		t.Fail()
//	}
//}

//func TestBuffer_Release(t *testing.T) {
//	buf := New(1000000)
//
//	buf.WriteInt64(1)
//
//	if buf.Ref() != 1 {
//		t.Fail()
//	}
//
//	buf.IncrRef()
//	if buf.Ref() != 2 {
//		t.Fail()
//	}
//
//	buf.Release()
//	buf.Release()
//	if buf.Ref() != 0 {
//		t.Fail()
//	}
//
//	if err := buf.Release(); err != ErrBufferReleased || buf.Ref() != 0 {
//		t.Fail()
//	}
//	if err := buf.WriteUInt8(1); err != ErrBufferReleased {
//		t.Fail()
//	}
//	if err := buf.WriteUInt16(1); err != ErrBufferReleased {
//		t.Fail()
//	}
//	if err := buf.WriteUInt32(1); err != ErrBufferReleased {
//		t.Fail()
//	}
//	if err := buf.WriteUInt64(1); err != ErrBufferReleased {
//		t.Fail()
//	}
//	if err := buf.WriteBytes([]byte{1}); err != ErrBufferReleased {
//		t.Fail()
//	}
//	if _, err := buf.CopyFromFile(3); err != ErrBufferReleased {
//		t.Fail()
//	}
//	if _, err := buf.CopyToFile(3); err != ErrBufferReleased {
//		t.Fail()
//	}
//}
//
//func TestBuffer_CopyToFile(t *testing.T) {
//	path := "/Users/heshan/tmp/test"
//	os.Remove(path)
//
//	f, err := os.Create(path)
//	if err != nil {
//		t.FailNow()
//	}
//	fd := int(f.Fd())
//
//	buf := New(1000000)
//	buf.WriteInt64(1)
//
//	if n, err := buf.CopyToFile(fd); err != nil || n != 8 {
//		t.FailNow()
//	}
//	if _, err := buf.CopyToFile(fd); err != ErrBufferNotEnough {
//		t.FailNow()
//	}
//	unix.Close(fd)
//
//	f, err = os.Open(path)
//	if err != nil {
//		t.FailNow()
//	}
//	fd = int(f.Fd())
//	buf = New(1000000)
//	if n, err := buf.CopyFromFile(fd); err != nil || n != 8 {
//		t.FailNow()
//	}
//
//	if n, err := buf.ReadInt64(); err != nil || n != 1 {
//		t.FailNow()
//	}
//
//	os.Remove(path)
//}
//
//func TestBuffer_WriteBytes(t *testing.T) {
//	buf := NewBufferSize(4)
//	buf.WriteByte(1)
//
//	data := []byte{1, 2, 3, 4}
//	buf.WriteBytes(data)
//
//	if _, err := buf.ReadByte(); err != nil {
//		t.Fail()
//	}
//	d2, err := buf.ReadBytes(4)
//	if err != nil {
//		t.Fail()
//	}
//	if len(d2) != 4 || d2[0] != 1 || d2[1] != 2 || d2[2] != 3 || d2[3] != 4 {
//		t.Fail()
//	}
//
//	if _, err := buf.ReadBytes(1); err != ErrBufferNotEnough {
//		t.Fail()
//	}
//
//	if err := buf.WriteBytes(nil); err != nil {
//		t.Fail()
//	}
//}
//
//func TestBuffer_WriteString(t *testing.T) {
//	str := "hello world"
//
//	buf := NewBufferSize(1)
//	if err := buf.WriteString(str); err != nil {
//		t.Fail()
//	}
//
//	if str2, err := buf.ReadString(len(str)); err != nil || str2 != str {
//		t.Fail()
//	}
//
//	if _, err := buf.ReadString(1); err != ErrBufferNotEnough {
//		t.Fail()
//	}
//}
//
func TestBuffer_Skip(t *testing.T) {
	buf := New(1000000)

	if err := buf.Skip(1); err != ErrNoEnoughData {
		t.FailNow()
	}

	buf.WriteByte(0)
	buf.WriteUInt32(100)

	if err := buf.Skip(1); err != nil {
		t.FailNow()
	}

	if n, err := buf.ReadUInt32(); err != nil || n != 100 {
		t.FailNow()
	}
}

func Test_BufferFindByte(t *testing.T) {
	buf := New(1000000)

	buf.WriteBytes([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9})

	if ind, ok, err := buf.FindByte(0, 5); !ok || err != nil || ind != 4 {
		t.Fail()
	}

	if ind, ok, err := buf.FindByte(8, 5); ok || err != nil || ind != -1 {
		t.Fail()
	}
}
