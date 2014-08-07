package orm

import (
	"testing"
	"time"
)

func TestString(t *testing.T) {
	var expected string = "abcd"

	actual := String(expected)

	if actual == nil {
		t.Fatalf("string pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Fatalf("string value is '%v', expected '%v'", *actual, expected)
	}
}

func TestInt(t *testing.T) {
	var expected int = 1

	actual := Int(expected)

	if actual == nil {
		t.Fatalf("int pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Fatalf("int value is %v, expected %v", *actual, expected)
	}
}

func TestInt8(t *testing.T) {
	var expected int8 = 1

	actual := Int8(expected)

	if actual == nil {
		t.Fatalf("int8 pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Fatalf("int8 value is %v, expected %v", *actual, expected)
	}
}

func TestInt16(t *testing.T) {
	var expected int16 = 1

	actual := Int16(expected)

	if actual == nil {
		t.Fatalf("int16 pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Fatalf("int16 value is %v, expected %v", *actual, expected)
	}
}

func TestInt32(t *testing.T) {
	var expected int32 = 1

	actual := Int32(expected)

	if actual == nil {
		t.Fatalf("int32 pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Fatalf("int32 value is %v, expected %v", *actual, expected)
	}
}

func TestInt64(t *testing.T) {
	var expected int64 = 1

	actual := Int64(expected)

	if actual == nil {
		t.Fatalf("int64 pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Fatalf("int64 value is %v, expected %v", *actual, expected)
	}
}

func TestBool(t *testing.T) {
	var expected bool = true

	actual := Bool(expected)

	if actual == nil {
		t.Fatalf("bool pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Fatalf("bool value is %v, expected %v", *actual, expected)
	}
}

func TestByte(t *testing.T) {
	var expected byte = 1

	actual := Byte(expected)

	if actual == nil {
		t.Fatalf("byte pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Fatalf("byte value is %v, expected %v", *actual, expected)
	}
}

func TestFloat32(t *testing.T) {
	var expected float32 = 1.0

	actual := Float32(expected)

	if actual == nil {
		t.Fatalf("float32 pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Fatalf("float32 value is %v, expected %v", *actual, expected)
	}
}

func TestFloat64(t *testing.T) {
	var expected float64 = 1.0

	actual := Float64(expected)

	if actual == nil {
		t.Fatalf("float64 pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Fatalf("float64 value is %v, expected %v", *actual, expected)
	}
}

func TestRune(t *testing.T) {
	var expected rune = 1.0

	actual := Rune(expected)

	if actual == nil {
		t.Fatalf("rune pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Fatalf("rune value is %v, expected %v", *actual, expected)
	}
}

func TestUint(t *testing.T) {
	var expected uint = 1

	actual := Uint(expected)

	if actual == nil {
		t.Fatalf("uint pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Fatalf("uint value is %v, expected %v", *actual, expected)
	}
}

func TestUint8(t *testing.T) {
	var expected uint8 = 1

	actual := Uint8(expected)

	if actual == nil {
		t.Fatalf("uint8 pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Fatalf("uint8 value is %v, expected %v", *actual, expected)
	}
}

func TestUint16(t *testing.T) {
	var expected uint16 = 1

	actual := Uint16(expected)

	if actual == nil {
		t.Fatalf("uint16 pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Fatalf("uint16 value is %v, expected %v", *actual, expected)
	}
}

func TestUint32(t *testing.T) {
	var expected uint32 = 1

	actual := Uint32(expected)

	if actual == nil {
		t.Fatalf("uint32 pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Fatalf("uint32 value is %v, expected %v", *actual, expected)
	}
}

func TestUint64(t *testing.T) {
	var expected uint64 = 1

	actual := Uint64(expected)

	if actual == nil {
		t.Fatalf("uint64 pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Fatalf("uint64 value is %v, expected %v", *actual, expected)
	}
}

func TestIs(t *testing.T) {

	type testCase struct {
		Value    *bool
		Expected bool
	}

	test_cases := []testCase{
		{Value: Bool(true), Expected: true},
		{Value: Bool(false), Expected: false},
		{Value: nil, Expected: false},
	}

	var actual bool

	for i, test_case := range test_cases {
		actual = Is(test_case.Value)
		if actual != test_case.Expected {
			t.Fatalf("Test case %v: orm.Is() returned '%v', expected '%v'", i, actual, test_case.Expected)
		}
	}
}

func TestTime(t *testing.T) {
	var expected = time.Now()
	actual := Time(expected)

	if actual == nil {
		t.Fatalf("time.Time pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Fatalf("time.Time value is %vm expected %v", *actual, expected)
	}
}

func TestDuration(t *testing.T) {
	expected, err := time.ParseDuration("1000ms")
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	actual := Duration(expected)

	if actual == nil {
		t.Fatalf("time.Duration pointer is nil, expected non-nil")
	}

	if *actual != expected {
		t.Fatalf("time.Duration value is %vm expected %v", *actual, expected)
	}
}
