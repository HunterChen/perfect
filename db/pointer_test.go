package pointer

import (
	"testing"
)

func TestString(t *testing.T) {
	var expected string = "abcd"

	actual := String(expected)

	if actual == nil {
		t.Errorf("string pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Errorf("string value is '%v', expected '%v'", *actual, expected)
	}
}

func TestInt(t *testing.T) {
	var expected int = 1

	actual := Int(expected)

	if actual == nil {
		t.Errorf("int pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Errorf("int value is %v, expected %v", *actual, expected)
	}
}

func TestInt8(t *testing.T) {
	var expected int8 = 1

	actual := Int8(expected)

	if actual == nil {
		t.Errorf("int8 pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Errorf("int8 value is %v, expected %v", *actual, expected)
	}
}

func TestInt16(t *testing.T) {
	var expected int16 = 1

	actual := Int16(expected)

	if actual == nil {
		t.Errorf("int16 pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Errorf("int16 value is %v, expected %v", *actual, expected)
	}
}

func TestInt32(t *testing.T) {
	var expected int32 = 1

	actual := Int32(expected)

	if actual == nil {
		t.Errorf("int32 pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Errorf("int32 value is %v, expected %v", *actual, expected)
	}
}

func TestInt64(t *testing.T) {
	var expected int64 = 1

	actual := Int64(expected)

	if actual == nil {
		t.Errorf("int64 pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Errorf("int64 value is %v, expected %v", *actual, expected)
	}
}

func TestBool(t *testing.T) {
	var expected bool = true

	actual := Bool(expected)

	if actual == nil {
		t.Errorf("bool pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Errorf("bool value is %v, expected %v", *actual, expected)
	}
}

func TestByte(t *testing.T) {
	var expected byte = 1

	actual := Byte(expected)

	if actual == nil {
		t.Errorf("byte pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Errorf("byte value is %v, expected %v", *actual, expected)
	}
}

func TestFloat32(t *testing.T) {
	var expected float32 = 1.0

	actual := Float32(expected)

	if actual == nil {
		t.Errorf("float32 pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Errorf("float32 value is %v, expected %v", *actual, expected)
	}
}

func TestFloat64(t *testing.T) {
	var expected float64 = 1.0

	actual := Float64(expected)

	if actual == nil {
		t.Errorf("float64 pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Errorf("float64 value is %v, expected %v", *actual, expected)
	}
}

func TestRune(t *testing.T) {
	var expected rune = 1.0

	actual := Rune(expected)

	if actual == nil {
		t.Errorf("rune pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Errorf("rune value is %v, expected %v", *actual, expected)
	}
}

func TestUint(t *testing.T) {
	var expected uint = 1

	actual := Uint(expected)

	if actual == nil {
		t.Errorf("uint pouinter is nil, expected non-nil")
	}

	if *actual != expected {
		t.Errorf("uint value is %v, expected %v", *actual, expected)
	}
}

func TestUint8(t *testing.T) {
	var expected uint8 = 1

	actual := Uint8(expected)

	if actual == nil {
		t.Errorf("uint8 pouinter is nil, expected non-nil")
	}

	if *actual != expected {
		t.Errorf("uint8 value is %v, expected %v", *actual, expected)
	}
}

func TestUint16(t *testing.T) {
	var expected uint16 = 1

	actual := Uint16(expected)

	if actual == nil {
		t.Errorf("uint16 pouinter is nil, expected non-nil")
	}

	if *actual != expected {
		t.Errorf("uint16 value is %v, expected %v", *actual, expected)
	}
}

func TestUint32(t *testing.T) {
	var expected uint32 = 1

	actual := Uint32(expected)

	if actual == nil {
		t.Errorf("uint32 pouinter is nil, expected non-nil")
	}

	if *actual != expected {
		t.Errorf("uint32 value is %v, expected %v", *actual, expected)
	}
}

func TestUint64(t *testing.T) {
	var expected uint64 = 1

	actual := Uint64(expected)

	if actual == nil {
		t.Errorf("uint64 pouinter is nil, expected non-nil")
	}

	if *actual != expected {
		t.Errorf("uint64 value is %v, expected %v", *actual, expected)
	}
}
