package perfect

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"time"
)

const (
	EC_P521 = 521
	EC_P384 = 384
)

type PrivateKey struct {
	Id   string
	Type int
	*ecdsa.PrivateKey
}

type serializableKey struct {
	Id   string `json:"id" bson:"id"`
	Type int    `json:"type" bson:"type"`
	D    string `json:"secret" bson:"secret"`
	X    string `json:"x" bson:"x"`
	Y    string `json:"y" bson:"y"`
}

func NewPrivateKey(key_type int) *PrivateKey {
	return &PrivateKey{
		Type: key_type,
	}
}

//generate a new Elliptical private key using the P521 curve and /dev/urandom
func GeneratePrivateKey(key_type int) (private_key *PrivateKey, err error) {
	curve, err := getCurve(key_type)
	if err != nil {
		return
	}

	ec_key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}

	private_key = &PrivateKey{
		PrivateKey: ec_key,
		Type:       key_type,
	}

	private_key.Id, err = GenerateKeyId()
	if err != nil {
		return nil, err
	}

	return
}

//generates a new Id by concatenating the current time in nanoseconds with 16 random bytes
func GenerateKeyId() (id string, err error) {
	//get the number of nanoseconds since the Epoch
	t := time.Now().UnixNano()
	b := make([]byte, 16)

	//read random bytes
	_, err = rand.Read(b)
	if err != nil {
		return
	}

	//convert the timestamp to string, then to byte array, and append it to the random bytes
	b = []byte(string(t) + string(b))

	//generate a sha512 hash of the bytes
	hashed := sha512.New().Sum(b)

	//return the hex representation of the hash
	id = hex.EncodeToString(hashed)

	return
}

//JSON.stringify
func (key *PrivateKey) MarshalJSON() (data []byte, err error) {
	skey := &serializableKey{
		Id:   key.Id,
		Type: key.Type,
		D:    base64.StdEncoding.EncodeToString(key.D.Bytes()),
		X:    base64.StdEncoding.EncodeToString(key.PublicKey.X.Bytes()),
		Y:    base64.StdEncoding.EncodeToString(key.PublicKey.Y.Bytes()),
	}

	data, err = json.Marshal(skey)
	return
}

//JSON.parse
func (key *PrivateKey) UnmarshalJSON(data []byte) (err error) {
	skey := &serializableKey{}
	err = json.Unmarshal(data, skey)
	if err != nil {
		return
	}

	key.Id = skey.Id
	key.Type = skey.Type
	key.PrivateKey = &ecdsa.PrivateKey{
		D: big.NewInt(0),
	}

	curve, err := getCurve(key.Type)
	if err != nil {
		return
	}

	//.PublicKey is the embedded .PrivateKey.PublicKey
	key.PublicKey = ecdsa.PublicKey{
		Curve: curve,
		X:     big.NewInt(0),
		Y:     big.NewInt(0),
	}

	//decode X
	bytes, err := base64.StdEncoding.DecodeString(skey.X)
	if err != nil {
		return
	}
	key.X.SetBytes(bytes)

	//decode Y
	bytes, err = base64.StdEncoding.DecodeString(skey.Y)
	if err != nil {
		return
	}
	key.Y.SetBytes(bytes)

	//decode D
	bytes, err = base64.StdEncoding.DecodeString(skey.D)
	if err != nil {
		return
	}
	key.D.SetBytes(bytes)

	return
}

func MD5Sum(s string) string {
	hash := md5.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

func getCurve(key_type int) (curve elliptic.Curve, err error) {
	switch key_type {
	case EC_P521:
		curve = elliptic.P521()
	case EC_P384:
		curve = elliptic.P384()
	default:
		err = fmt.Errorf("%v: unknown key type", key_type)
	}

	return
}
