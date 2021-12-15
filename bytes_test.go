package buffer

import "testing"

func Test_bytesPool_getBytes(t *testing.T) {
	pool := newBytesPoolSize(3)

	data := pool.get(4)
	if len(data) != 4 {
		t.Fail()
	}
	pool.put(data)

	data = getBytes(5)
	if len(data) != 8 {
		t.Fail()
	}
	pool.put(data[:5])

	pool.put(make([]byte, 16))
}

func Test_newBytesPoolSize(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fail()
		}
	}()

	pool := newBytesPoolSize(1)
	if pool == nil {
		t.Fail()
	}

	newBytesPoolSize(-1)
}

func Test_getBytes(t *testing.T) {
	data := getBytes(4)
	if len(data) != 4 {
		t.Fail()
	}
	releaseBytes(data)
}

func Test_getBytes2(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fail()
		}
	}()

	getBytes(-1)
}

func Test_releaseBytes(t *testing.T) {
	data := make([]byte, 3)
	releaseBytes(data)

	data = nil
	releaseBytes(data)
}
