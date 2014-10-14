package perfect

import (
	"log"
	"math/big"
	"testing"
)

func TestNewPrivateKey(t *testing.T) {
	key_types := []int{EC_P521, EC_P384}

	for _, kt := range key_types {
		key := NewPrivateKey(kt)

		if key.Type != kt {
			t.Errorf("key.Type = %v, expected %v", key.Type, kt)
		}
	}
}

func BenchmarkNewPrivateKey_P521(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewPrivateKey(EC_P521)
	}
}

func BenchmarkNewPrivateKey_P384(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewPrivateKey(EC_P384)
	}
}

func TestGeneratePrivateKey(t *testing.T) {
	key, err := GeneratePrivateKey(EC_P521)

	if err != nil {
		t.Errorf("err = %v", err)
	}

	if key == nil {
		t.Errorf("key is nil, expected non-nil")
	}

	if len(key.Id) == 0 {
		t.Errorf("len(key.Id) == %v, expected non-zero", len(key.Id))
	}

	//only EC_P521 is supported
	if key.Type != EC_P521 {
		t.Errorf("key.Type = %v, expected %v (EC_P521)", key.Type, EC_P521)
	}

	//test that the private key is not nil
	if key.PrivateKey == nil {
		t.Errorf("key.PrivateKey = %v, expected non-nil", key.PrivateKey)
	}
}

func BenchmarkGeneratePrivateKey_P521(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GeneratePrivateKey(EC_P521)
	}
}

func BenchmarkGeneratePrivateKey_P384(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GeneratePrivateKey(EC_P384)
	}
}

func TestGenerateKeyId(t *testing.T) {
	id, err := GenerateKeyId()

	if err != nil {
		t.Errorf("err = %v", err)
	}

	if len(id) == 0 {
		t.Errorf("len(id) == %v, expected non-zero", len(id))
	}

	id2, err := GenerateKeyId()

	if err != nil {
		t.Errorf("err2 = %v", err)
	}

	if len(id2) == 0 {
		t.Errorf("len(id2) == %v, expected non-zero", len(id2))
	}

	if id == id2 {
		t.Errorf("id1 == id2, expected them to be different\n\tid1=%v\n\tid2=%v\n", id, id2)
	}
}

func BenchmarkGenerateKeyId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateKeyId()
	}
}

func TestMarshalJSON(t *testing.T) {
	key, err := GeneratePrivateKey(EC_P521)

	if err != nil {
		t.Errorf("err = %v", err)
	}

	if key == nil {
		t.Errorf("key is nil, expected non-nil")
	}

	data, err := key.MarshalJSON()

	if err != nil {
		t.Errorf("err = %v", err)
	}

	if len(data) == 0 {
		t.Errorf("len(data) == %v, expected %v", len(data), 0)
	}
}

func BenchmarkMarshalJSON(b *testing.B) {
	key, err := GeneratePrivateKey(EC_P521)

	if err != nil {
		b.Errorf("err = %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = key.MarshalJSON()
	}
}

func TestUnmarshalJSON(t *testing.T) {
	key, err := GeneratePrivateKey(EC_P521)
	if err != nil {
		t.Errorf("err = %v", err)
	}

	keydata, err := key.MarshalJSON()
	if err != nil {
		t.Errorf("err = %v", err)
	}

	key2 := NewPrivateKey(EC_P521)
	err = key2.UnmarshalJSON(keydata)

	if key2.Id != key.Id {
		t.Errorf("key2.Id = %v, expected %v", key2.Id, key.Id)
	}

	if key2.Type != key.Type {
		t.Errorf("key2.Type = %v, expected %v", key2.Type, key.Type)
	}

	if key2.D.Cmp(key.D) != 0 {
		t.Errorf("key2.PrivateKey.D = %v, expected %v", key2.D, key.D)
	}

	if key2.X.Cmp(key.X) != 0 {
		t.Errorf("key2.PrivateKey.X = %v, expected %v", key2.X, key.X)
	}

	if key2.Y.Cmp(key.Y) != 0 {
		t.Errorf("key2.PrivateKey.Y = %v, expected %v", key2.Y, key.Y)
	}
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	key, err := GeneratePrivateKey(EC_P521)
	if err != nil {
		b.Errorf("err = %v", err)
	}

	keydata, err := key.MarshalJSON()
	if err != nil {
		b.Errorf("err = %v", err)
	}

	key2 := NewPrivateKey(EC_P521)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = key2.UnmarshalJSON(keydata)
	}
}

func TestMD5Sum(t *testing.T) {
	s := "The quick brown fox jumped over the lazy dog"
	expected := "08a008a01d498c404b0c30852b39d3b8"
	actual := MD5Sum(s)
	if actual != expected {
		t.Fatalf("actual md5 is '%v', expected '%v' for string: '%v'", actual, expected, s)
	}
}

func BenchmarkPrivateKey_MarshalBigIntInt64(b *testing.B) {
	var n *big.Int = big.NewInt(1000)
	var i64 int64 = 0

	for i := 0; i < b.N; i++ {
		i64 = n.Int64()
	}

	log.Println(i64)
}

func TestMarshalBigInt(t *testing.T) {
	var expected *big.Int = big.NewInt(1000)
	i64 := expected.Int64()

	actual := big.NewInt(i64)

	if expected.Cmp(actual) != 0 {
		t.Fatalf("big.Int numbers are not equal. expected: %v, actual %v", expected.String(), actual.String())
	}
}

func TestPrivateKey_Equals(t *testing.T) {
	expected, err := GeneratePrivateKey(EC_P521)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	actual := expected

	if !expected.Equals(actual) {
		t.Fatalf("keys are not equal:\nexpected: %v\nactual: %v", expected, actual)
	}
}

func BenchmarkPrivateKey_MarshalBigIntString(b *testing.B) {
	var n *big.Int = big.NewInt(1000)
	var sn string = ""

	for i := 0; i < b.N; i++ {
		sn = n.String()
	}

	log.Println(sn)
}

func TestPrivateKey_SerializeBSON(t *testing.T) {
	expected, err := GeneratePrivateKey(EC_P521)
	if err != nil {
		t.Fatalf("err = %v")
	}

	serialized, err := expected.MarshalBSON()
	if err != nil {
		t.Fatalf("err = %v")
	}

	actual := &PrivateKey{}
	err = actual.UnmarshalBSON(serialized)
	if err != nil {
		t.Fatalf("err = %v")
	}

	if !expected.Equals(actual) {
		t.Fatalf("keys are not equal\nexpected: %#v\n- ecdsa key: %#v\nactual: %#v\n- ecdsa key: %#v", expected, *expected.PrivateKey, actual, *actual.PrivateKey)
	}
}
